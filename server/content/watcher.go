package content

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ContentWatcher watches content files for changes and reloads them
type ContentWatcher struct {
	manager *Manager
	watcher *fsnotify.Watcher
	stop    chan bool
}

// NewContentWatcher creates a new content file watcher
func NewContentWatcher(mgr *Manager) (*ContentWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	return &ContentWatcher{
		manager: mgr,
		watcher: watcher,
		stop:    make(chan bool),
	}, nil
}

// Start begins watching content directories
func (cw *ContentWatcher) Start() error {
	// Watch all content subdirectories
	paths := []string{
		filepath.Join(cw.manager.basePath, "default", "skills"),
		filepath.Join(cw.manager.basePath, "default", "npcs", "templates"),
		filepath.Join(cw.manager.basePath, "default", "items"),
		filepath.Join(cw.manager.basePath, "default", "rooms"),
		filepath.Join(cw.manager.basePath, "default", "quests"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if err := cw.addWatchRecursive(path); err != nil {
				log.Printf("Warning: Could not watch %s: %v", path, err)
			}
		}
	}

	go cw.watchLoop()
	log.Println("Content watcher started - hot-reload enabled")
	return nil
}

// addWatchRecursive adds a directory and all subdirectories to the watcher
func (cw *ContentWatcher) addWatchRecursive(path string) error {
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := cw.watcher.Add(walkPath); err != nil {
				return fmt.Errorf("failed to watch %s: %w", walkPath, err)
			}
			log.Printf("Watching: %s", walkPath)
		}
		return nil
	})
}

// watchLoop handles file change events
func (cw *ContentWatcher) watchLoop() {
	debounce := make(map[string]time.Time)
	debounceInterval := 500 * time.Millisecond

	for {
		select {
		case event, ok := <-cw.watcher.Events:
			if !ok {
				return
			}

			// Only process write events for YAML files
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				if !strings.HasSuffix(event.Name, ".yaml") && !strings.HasSuffix(event.Name, ".yml") {
					continue
				}

				// Debounce - avoid reloading same file multiple times
				if lastTime, exists := debounce[event.Name]; exists && time.Since(lastTime) < debounceInterval {
					continue
				}
				debounce[event.Name] = time.Now()

				log.Printf("Content changed: %s", event.Name)
				if err := cw.reloadFile(event.Name); err != nil {
					log.Printf("Error reloading %s: %v", event.Name, err)
				}
			}

		case err, ok := <-cw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)

		case <-cw.stop:
			return
		}
	}
}

// reloadFile reloads a single content file
func (cw *ContentWatcher) reloadFile(path string) error {
	// Determine content type from path
	contentType := cw.detectContentType(path)
	
	switch contentType {
	case "skill":
		return cw.reloadSkill(path)
	case "npc":
		return cw.reloadNPC(path)
	case "item":
		return cw.reloadItem(path)
	case "room":
		return cw.reloadRoom(path)
	case "quest":
		return cw.reloadQuest(path)
	default:
		return fmt.Errorf("unknown content type for %s", path)
	}
}

// detectContentType determines what type of content a file contains
func (cw *ContentWatcher) detectContentType(path string) string {
	path = strings.ToLower(path)
	if strings.Contains(path, "/skills/") {
		return "skill"
	}
	if strings.Contains(path, "/npcs/") {
		return "npc"
	}
	if strings.Contains(path, "/items/") {
		return "item"
	}
	if strings.Contains(path, "/rooms/") {
		return "room"
	}
	if strings.Contains(path, "/quests/") {
		return "quest"
	}
	return "unknown"
}

// reloadSkill reloads a single skill file
func (cw *ContentWatcher) reloadSkill(path string) error {
	return cw.manager.ReloadSkillFile(path)
}

// reloadNPC reloads a single NPC file
func (cw *ContentWatcher) reloadNPC(path string) error {
	return cw.manager.ReloadNPCFile(path)
}

// reloadItem reloads a single item file
func (cw *ContentWatcher) reloadItem(path string) error {
	return cw.manager.ReloadItemFile(path)
}

// reloadRoom reloads a single room file
func (cw *ContentWatcher) reloadRoom(path string) error {
	return cw.manager.ReloadRoomFile(path)
}

// reloadQuest reloads a single quest file
func (cw *ContentWatcher) reloadQuest(path string) error {
	return cw.manager.ReloadQuestFile(path)
}

// Stop stops the content watcher
func (cw *ContentWatcher) Stop() {
	close(cw.stop)
	cw.watcher.Close()
}
