package filewalker

import (
	"go-mp3/eta"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

func WalkDirs(callback func(file string), roots ...string) {
	files := make(chan string, 100000)
	total := atomic.Uint64{}
	et := eta.NewEta(0)
	go func() {
		for _, root := range roots {
			_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
				if !d.IsDir() {
					files <- path
					et.SetTotal(total.Add(1))
				}
				return err
			})
		}
		close(files)
	}()
	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			for file := range files {
				callback(file)
				et.IncCount()
			}
			wg.Done()
		}()
	}
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				log.Printf("%s\n", et.String())
			case <-done:
				return
			}
		}
	}()
	wg.Wait()
	done <- true
}
