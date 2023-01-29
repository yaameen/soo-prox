package utils

import (
	"log"
	"sooprox/types"

	"github.com/fsnotify/fsnotify"
)

func WatchConfig(configChanged chan types.Config, files ...string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					configChanged <- ReadConfig(files...)
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	for _, v := range files {
		err = watcher.Add(v)
		if err != nil {
			log.Fatal(err)
		}
	}
	<-done
}
