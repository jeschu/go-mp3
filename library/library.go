package library

import (
	"github.com/dhowden/tag"
	t "github.com/tcolgate/mp3"
	"go-mp3/eta"
	"golang.org/x/text/message"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Entry struct {
	File        string         `json:"file"`
	Title       string         `json:"title"`
	Album       string         `json:"album"`
	Artist      string         `json:"artist"`
	AlbumArtist string         `json:"albumArtist"`
	Composer    string         `json:"composer"`
	Year        int            `json:"year"`
	Genre       string         `json:"genre"`
	Format      tag.Format     `json:"format"`
	Comment     string         `json:"comment"`
	Track       int            `json:"track"`
	Tracks      int            `json:"tracks"`
	Disc        int            `json:"disc"`
	Discs       int            `json:"discs"`
	FileType    tag.FileType   `json:"fileType"`
	Lyrics      string         `json:"lyrics"`
	Picture     *tag.Picture   `json:"picture"`
	Duration    *time.Duration `json:"duration"`
}

type Library struct {
	entries      []*Entry
	albums       map[string][]*Entry
	artists      map[string][]*Entry
	albumArtists map[string][]*Entry
	composers    map[string][]*Entry
}

func NewLibrary() *Library {
	return &Library{
		entries:      make([]*Entry, 0),
		albums:       make(map[string][]*Entry),
		artists:      make(map[string][]*Entry),
		albumArtists: make(map[string][]*Entry),
		composers:    make(map[string][]*Entry),
	}
}

func (lib *Library) Add(file string, md tag.Metadata) {
	track, tracks := md.Track()
	disc, discs := md.Disc()
	entry := &Entry{
		File:        file,
		Title:       md.Title(),
		Album:       md.Album(),
		Genre:       md.Genre(),
		Artist:      md.Artist(),
		AlbumArtist: md.AlbumArtist(),
		Format:      md.Format(),
		Comment:     md.Comment(),
		Composer:    md.Composer(),
		Disc:        disc,
		Discs:       discs,
		Track:       track,
		Tracks:      tracks,
		FileType:    md.FileType(),
		Lyrics:      md.Lyrics(),
		Year:        md.Year(),
		Picture:     md.Picture(),
		Duration:    nil,
	}
	lib.entries = append(lib.entries, entry)
	appendEntry(lib.albums, entry.Album, entry)
	appendEntry(lib.artists, entry.Artist, entry)
	appendEntry(lib.albumArtists, entry.AlbumArtist, entry)
	appendEntry(lib.composers, entry.Composer, entry)
}

func appendEntry(m map[string][]*Entry, key string, entry *Entry) {
	var (
		entries []*Entry
		ok      bool
	)
	if entries, ok = m[key]; !ok {
		entries = make([]*Entry, 0)
	}
	m[key] = append(entries, entry)
}

func (lib *Library) TitleCount() int { return len(lib.entries) }

func (lib *Library) Stats() string {
	s := "Library Stats:\n"
	s += p.Sprintf("  %10d Titles\n", len(lib.entries))
	s += p.Sprintf("  %10d AlbumCount\n", len(lib.albums))
	s += p.Sprintf("  %10d ArtistCount\n", len(lib.artists))
	s += p.Sprintf("  %10d Album-ArtistCount\n", len(lib.albumArtists))
	s += p.Sprintf("  %10d ComposerCount\n", len(lib.composers))
	return s
}

func (lib *Library) Entries() []*Entry { return lib.entries }
func EntriesWithoutLyricsAndPictures(libEntries []*Entry) []Entry {
	entries := make([]Entry, len(libEntries))
	for _, libEntry := range libEntries {
		entries = append(entries,
			Entry{
				File:        libEntry.File,
				Title:       libEntry.Title,
				Album:       libEntry.Album,
				Artist:      libEntry.Artist,
				AlbumArtist: libEntry.AlbumArtist,
				Composer:    libEntry.Composer,
				Year:        libEntry.Year,
				Genre:       libEntry.Genre,
				Format:      libEntry.Format,
				Comment:     libEntry.Comment,
				Track:       libEntry.Track,
				Tracks:      libEntry.Tracks,
				Disc:        libEntry.Disc,
				Discs:       libEntry.Discs,
				FileType:    libEntry.FileType,
				Lyrics:      "",
				Picture:     nil,
				Duration:    libEntry.Duration,
			})
	}
	return entries

}

func (lib *Library) TitlesWithoutLyricsAndPictures() []Entry {
	return EntriesWithoutLyricsAndPictures(lib.entries)
}

func (lib *Library) Albums() map[string][]*Entry         { return lib.albums }
func (lib *Library) Artists() map[string][]*Entry        { return lib.artists }
func (lib *Library) AlbumArtists() map[string][]*Entry   { return lib.albumArtists }
func (lib *Library) Composers() map[string][]*Entry      { return lib.composers }
func (lib *Library) AlbumCount() map[string]uint64       { return countingMap(lib.artists) }
func (lib *Library) ArtistCount() map[string]uint64      { return countingMap(lib.artists) }
func (lib *Library) AlbumArtistCount() map[string]uint64 { return countingMap(lib.albumArtists) }
func (lib *Library) ComposerCount() map[string]uint64    { return countingMap(lib.composers) }
func countingMap(in map[string][]*Entry) map[string]uint64 {
	out := make(map[string]uint64, len(in))
	for key, entries := range in {
		out[key] = uint64(len(entries))
	}
	return out
}

var p = message.NewPrinter(message.MatchLanguage("de_DE"))

func (lib *Library) TotalDuration() (totalDuration time.Duration, complete bool) {
	totalDuration = time.Duration(0)
	complete = true
	for _, entry := range lib.entries {
		if entry.Duration != nil {
			totalDuration += *entry.Duration
		} else {
			complete = false
		}
	}
	return
}

func (lib *Library) UpdateDurations() (*eta.Eta, chan bool) {
	entries := make(chan *Entry, len(lib.entries))
	total := uint64(0)
	for _, e := range lib.entries {
		if e.Duration == nil {
			entries <- e
			total++
		}
	}
	close(entries)
	et := eta.NewEta(total)
	done := make(chan bool, 1)
	go func() {
		wg := sync.WaitGroup{}
		for i := 0; i < runtime.NumCPU(); i++ {
			wg.Add(1)
			go func() {
				for entry := range entries {
					entry.Duration = calcDuration(entry.File)
					et.IncCount()
				}
				wg.Done()
			}()
		}
		wg.Wait()
		done <- true
		close(done)
	}()
	return &et, done
}

func calcDuration(path string) *time.Duration {
	duration := time.Duration(-1)
	f, err := os.Open(path)
	if err != nil {
		return &duration
	}
	defer func(f *os.File) { _ = f.Close() }(f)
	dec := t.NewDecoder(f)
	var frame t.Frame
	skipped := 0
	for {
		if err := dec.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			return &duration
		}
		duration = duration + frame.Duration()
	}
	return &duration
}
