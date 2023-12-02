package mp3scanner

import (
	"github.com/dhowden/tag"
	"go-mp3/library"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var acceptedExtensions = map[string]any{
	"mp3": nil,
	"m4a": nil,
	"wav": nil,
}

type Scanner struct {
	errors  []string
	fn      func(file string)
	library *library.Library
}

func NewScanner() *Scanner {
	mu := sync.Mutex{}
	scanner := new(Scanner)
	scanner.errors = make([]string, 0)
	scanner.library = library.NewLibrary()
	scanner.fn = func(file string) {
		ext := strings.ToLower(filepath.Ext(file))[1:]
		if _, ok := acceptedExtensions[ext]; ok {
			f, err := os.Open(file)
			if err == nil {
				defer func(f *os.File) { _ = f.Close() }(f)
				md, err := tag.ReadFrom(f)
				if err == nil {
					mu.Lock()
					scanner.library.Add(file, md)
					mu.Unlock()
				}
			} else {
				mu.Lock()
				scanner.errors = append(scanner.errors, file)
				mu.Unlock()
			}
		}
	}
	return scanner
}

func (scanner *Scanner) Fn() func(file string)     { return scanner.fn }
func (scanner *Scanner) ErrCount() int             { return len(scanner.errors) }
func (scanner *Scanner) Errors() []string          { return scanner.errors }
func (scanner *Scanner) Library() *library.Library { return scanner.library }
