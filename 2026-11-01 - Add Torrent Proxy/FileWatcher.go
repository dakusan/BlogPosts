package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

//startFileWatcher creates the fsnotify watcher, starts watching cfg.Root recursively, and kicks off the event-loop goroutine. The caller is responsible for closing the returned watcher.
func startFileWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}

	//Start watching from the root directory
	if err := addWatch(watcher, cfg.Root); err != nil {
		log.Fatal("Error watching root directory: " + err.Error())
	}

	go watchEvents(watcher)

	return watcher
}

//Recursive function to add watches to directories
func addWatch(watcher *fsnotify.Watcher, dir string) error {
	if cfg.DoneDir != "" && isWithinRoot(dir, cfg.DoneDir) {
		return nil //Never watch doneDir: files moved there must not generate events and get reprocessed
	}
	if slices.Contains(watcher.WatchList(), dir) {
		return nil //Already watching this directory
	}

	if err := watcher.Add(dir); err != nil {
		return err
	}
	log.Println("Started watching:", dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			if err := addWatch(watcher, subdir); err != nil {
				log.Println("Error adding watch to subdirectory:", err)
			}
		}
	}
	return nil
}

func watchEvents(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
//			if event.Op&(fsnotify.Chmod|fsnotify.Write) != 0 { //Ignore chmod and write
//				return
//			}
			log.Println("Event:", event)

			//Only handle create events
			if event.Op&(fsnotify.Create|fsnotify.Rename) == 0 {
				break
			}

			//Stat the new object. Lstat (not Stat) is used deliberately so a symlink is never followed: this keeps the watched/moved tree from being able to escape root via a symlink planted inside it.
			info, err := os.Lstat(event.Name)
			if err != nil {
				log.Println("Error statting:", err)
				break
			}

			//Handle watching new directory
			if info.IsDir() {
				if err := addWatch(watcher, event.Name); err != nil {
					log.Println("Error adding watch to new directory:", err)
				}
				break
			}

			//Make sure the new file is a regular file
			if !info.Mode().IsRegular() {
				break
			}

			//Process the new file
			log.Printf("Return for “%s”: %s\n", event.Name, ProcessFile(event.Name))

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

func ProcessFile(eventPath string) string {
	//Get the path of the new file, relative to root, and require it to match fileRegex (and not fileRegexBlacklist).
	//Any match error (e.g. a timeout from a pathological pattern/input) is treated as “ignore the file”, not “allow it”.
	relPath := eventPath[len(cfg.Root):]
	if matched, err := fileRegex.MatchString(relPath); err != nil {
		return "Error evaluating fileRegex, ignoring: " + err.Error()
	} else if !matched {
		return "File does not match regex, ignoring: " + relPath
	}
	if fileRegexBlacklist != nil {
		if matched, err := fileRegexBlacklist.MatchString(relPath); err != nil {
			return "Error evaluating fileRegexBlacklist, ignoring: " + err.Error()
		} else if matched {
			return "File matches blacklist regex, ignoring: " + relPath
		}
	}
	newFilePath := strings.Split(relPath, "/")
	newFile := newFilePath[len(newFilePath)-1]

	//Create the new directory (isWithinRoot guards against a crafted/misconfigured path ever landing outside newRoot)
	newDir := filepath.Clean(cfg.NewRoot+strings.Join(newFilePath[0:len(newFilePath)-1], "/")) + "/"
	if !isWithinRoot(newDir, cfg.NewRoot) {
		return "Refusing to write outside newRoot: " + newDir
	}
	if err := os.MkdirAll(newDir, 0775); err != nil {
		return "Could not create new path: " + err.Error()
	}

	//Copy the file
	time.Sleep(3 * time.Second) //Pause for 3 seconds to make sure the file finishes copying
	source, err := os.Open(eventPath)
	if err != nil {
		return "Could not open file for moving: " + err.Error()
	}
	defer source.Close()
	destination, err := os.Create(newDir + newFile)
	if err != nil {
		return "Could not create new file: " + err.Error()
	}
	defer destination.Close()
	if _, err := io.Copy(destination, source); err != nil {
		return "Could not finish copying file: " + err.Error()
	}
	source.Close()
	log.Println("File copied to: ", newDir+newFile)

	//Remove the original file: delete it if no doneDir is configured, otherwise move it there (only warn if either fails)
	if cfg.DoneDir == "" {
		if err := os.Remove(eventPath); err != nil {
			log.Printf("Warning: Could not delete “%s”: %s\n", eventPath, err.Error())
		}
	} else if err := os.MkdirAll(cfg.DoneDir, 0775); err != nil {
		log.Printf("Warning: Could not create DONE directory “%s”: %s\n", cfg.DoneDir, err.Error())
	} else if err := os.Rename(eventPath, filepath.Clean(cfg.DoneDir)+"/"+newFile); err != nil {
		log.Printf("Warning: Could not move “%s” to DONE directory: %s\n", eventPath, err.Error())
	}

	/*
	for i := len(newFilePath); i > 1; i-- {
		pathToRemove := root + strings.Join(newFilePath[0:i], "/")
		if err := os.Remove(pathToRemove); err != nil {
			log.Printf("Warning: Could not remove “%s”: %s\n", pathToRemove, err.Error())
			break
		} else {
			log.Println("Removed: " + pathToRemove)
		}
		break //I decided to only remove the file to move and leave the directories in tact
	}
	*/

	//If the file is a torrent, then start it in qbittorrent
	if !strings.HasSuffix(strings.ToLower(newFile), ".torrent") {
		return "Ignoring non-torrent file"
	}

	//Send the torrent to qbittorrent via the WebUI API; qbitMonitor picks up progress/cleanup from here
	if err := QbitAddTorrent(newDir+newFile, newDir, QbitManagedTag); err != nil {
		return "Error adding torrent to qbittorrent: " + err.Error()
	}
	log.Println("Added to qbittorrent: ", newDir+newFile)

	return "Success"
}