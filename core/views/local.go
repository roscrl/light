package views

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func watchLocalTemplates(views *Views, changed chan<- struct{}) {
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

	err = addWatchers(PathViews)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("watching templates")

	fileEndingsToReloadOn := map[string]struct{}{
		".tmpl":  {},
		".css":   {},
		".js":    {},
		".svg":   {},
		".png":   {},
		".jpg":   {},
		".jpeg":  {},
		".gif":   {},
		".ico":   {},
		".woff":  {},
		".woff2": {},
		".ttf":   {},
		".eot":   {},
		".otf":   {},
	}

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
					if _, ok := fileEndingsToReloadOn[filepath.Ext(event.Name)]; !ok {
						continue
					}

					log.Printf("%s changed %s, reloading~", event.Name, event.Op)

					templates := findAndParseTemplates(os.DirFS(PathTemplates), views.funcMap)

					views.templates = templates // not thread safe but only used for local development

					changed <- struct{}{}
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
