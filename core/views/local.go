package views

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func watchLocalTemplates(views *Views) {
	watcher, err := fsnotify.NewWatcher() // leaks but only used for local development
	if err != nil {
		log.Fatal(err)
	}

	addWatchers := func(path string) error {
		err := watcher.Add(path)
		if err != nil {
			return fmt.Errorf("error adding watcher for %s: %w", path, err)
		}

		// Walk the directory and add watchers for subdirectories
		err = filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking %s: %w", subpath, err)
			}
			if info.IsDir() {
				err = watcher.Add(subpath)
				if err != nil {
					return fmt.Errorf("adding watcher for %s: %w", subpath, err)
				}
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("walking %s: %w", path, err)
		}

		return nil
	}

	err = addWatchers("./" + PathViews)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("watching templates")

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Chmod) {
					continue
				}

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					if !strings.HasSuffix(event.Name, ".tmpl") {
						continue
					}

					log.Printf("%s changed %s, reloading~", event.Name, event.Op)

					templates := findAndParseTemplates(os.DirFS(PathTemplates), views.funcMap)

					views.templates = templates // not thread safe but only used for local development
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Fatal(err)
			}
		}
	}()
}
