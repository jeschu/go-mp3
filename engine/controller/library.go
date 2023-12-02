package controller

import (
	"github.com/gin-gonic/gin"
	"go-mp3/filewalker"
	"go-mp3/library"
	"go-mp3/mp3scanner"
	"log"
	"net/http"
	"net/url"
	"time"
)

type LibraryController struct {
	library *library.Library
}

func NewLibraryController() *LibraryController {
	c := LibraryController{}
	go func() {
		c.library = loadLibrary()
	}()
	return &c
}

func loadLibrary() *library.Library {
	scanner := mp3scanner.NewScanner()
	start := time.Now()
	filewalker.WalkDirs(scanner.Fn(), "/Users/jens/Music/Music/Media.localized/Music")
	lib := scanner.Library()
	log.Printf("scanned %d files in %s\n", lib.TitleCount(), time.Now().Sub(start))
	errors := scanner.Errors()
	if len(errors) > 0 {
		log.Printf("  got %d errors\n", len(errors))
		for _, file := range scanner.Errors() {
			log.Printf("    %s\n", file)
		}
	}
	return lib
}

func (controller *LibraryController) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/lib/titles/count",
		controller.handleWithLibrary(func(c *gin.Context, lib *library.Library) {
			c.PureJSON(http.StatusOK, gin.H{"count": lib.TitleCount()})
		}))
	engine.GET("/lib/titles",
		controller.handleWithLibrary(func(c *gin.Context, lib *library.Library) {
			c.PureJSON(http.StatusOK, lib.TitlesWithoutLyricsAndPictures())
		}))
	engine.GET("/lib/album/:album",
		controller.handleWithLibrary(func(c *gin.Context, lib *library.Library) {
			album, _ := url.QueryUnescape(c.Param("album"))
			c.PureJSON(http.StatusOK,
				library.EntriesWithoutLyricsAndPictures(lib.Albums()[album]))
		}))
	engine.GET("/lib/albums",
		controller.handleWithLibrary(func(c *gin.Context, lib *library.Library) {
			c.PureJSON(http.StatusOK, lib.AlbumCount())
		}))
	engine.GET("/lib/artists",
		controller.handleWithLibrary(func(c *gin.Context, lib *library.Library) {
			c.PureJSON(http.StatusOK, lib.ArtistCount())
		}))
	engine.GET("/lib/artist/:artist",
		controller.handleWithLibrary(func(c *gin.Context, lib *library.Library) {
			artist, _ := url.QueryUnescape(c.Param("artist"))
			c.PureJSON(http.StatusOK,
				library.EntriesWithoutLyricsAndPictures(
					lib.Artists()[artist],
				))
		}))

	engine.GET("/lib/albumartists",
		controller.handleWithLibrary(func(c *gin.Context, lib *library.Library) {
			c.PureJSON(http.StatusOK, lib.AlbumArtistCount())
		}))
	engine.GET("/lib/composers",
		controller.handleWithLibrary(func(c *gin.Context, lib *library.Library) {
			c.PureJSON(http.StatusOK, lib.ComposerCount())
		}))
}

func (controller *LibraryController) handleWithLibrary(handler withLibrary) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		if controller.library == nil {
			AbortWithErrorResponse(c, http.StatusServiceUnavailable, "library not loaded yet")
		} else {
			handler(c, controller.library)
		}
	}
}

type withLibrary = func(c *gin.Context, lib *library.Library)
type EntryCount map[string]uint64
