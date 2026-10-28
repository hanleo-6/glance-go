package watch

import (
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func File(targetFile string, onChange func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	dir := filepath.Dir(targetFile)
	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		return err
	}

	var timer *time.Timer
	debounceDuration := 50 * time.Millisecond

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filepath.Clean(event.Name) == filepath.Clean(targetFile) &&
					(event.Has(fsnotify.Write) ||
						event.Has(fsnotify.Create) ||
						event.Has(fsnotify.Rename)) {

					log.Println("event:", event.Op)

					// trigger the change after debouncing
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(debounceDuration, func() {
						log.Println("changed!")
						onChange()
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watch error:", err)
			}
		}
	}()

	return nil
}
