package watch

import (
	"log"
	"path/filepath"

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

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// check if this is event for targetfile
				// distinguish between write, create and rename
				// if both, log event and call onChange
				if filepath.Clean(event.Name) == filepath.Clean(targetFile) &&
			      (event.Has(fsnotify.Write) ||
				   event.Has(fsnotify.Create) ||
				   event.Has(fsnotify.Rename)) {
					log.Println("event:", event.Op)
					onChange()
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

