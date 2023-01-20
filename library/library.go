package library

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/goccy/go-json"
	"log"
	"os"
	"path"
	"sync"
)

type Library struct {
	path     string
	librarys map[string][]string
	lock     sync.RWMutex
}

func New(path string) *Library {
	return &Library{
		path: path,
	}
}

func (i *Library) IsExist(id string) bool {
	i.lock.RLock()
	defer i.lock.RUnlock()
	if _, e := i.librarys[id]; e {
		return true
	}
	return false
}

func (i *Library) ListItem(libId string) []string {
	i.lock.RLock()
	defer i.lock.RUnlock()
	if _, e := i.librarys[libId]; e {
		return i.librarys[libId]
	}
	return nil
}

func (i *Library) Save() error {
	i.lock.RLock()
	defer i.lock.RUnlock()
	f, err := os.OpenFile(i.path,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open librarys file error: %s", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(&i.librarys)
	if err != nil {
		return fmt.Errorf("encode librarys error: %s", err)
	}
	return nil
}

func (i *Library) Load() error {
	i.lock.Lock()
	defer i.lock.Unlock()
	f, err := os.Open(i.path)
	if err != nil {
		return fmt.Errorf("open librarys file error: %s", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&i.librarys)
	if err != nil {
		return fmt.Errorf("decode librarys error: %s", err)
	}
	return nil
}

func (i *Library) StartWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("new watcher error: %s", err)
	}
	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				if event.Has(fsnotify.Write) {
					if event.Name != path.Base(i.path) {
						continue
					}
					log.Println("lib file changed, reloading...")
					err := i.Load()
					if err != nil {
						log.Printf("reload error: %s", err)
					}
					log.Println("reload success")
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.Panicf("watcher error: %s", err)
				}
			}
		}
	}()
	err = watcher.Add(path.Dir(i.path))
	if err != nil {
		return fmt.Errorf("watch file error: %s", err)
	}
	return nil
}
