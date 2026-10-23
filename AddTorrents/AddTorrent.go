package main

import (
	"bytes"
	"context"
	"github.com/fsnotify/fsnotify"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const root = "/home/plex/Needed torrents/"
const newRoot = "/home/plex/library/"
const logFile = "/home/plex/.var/app/org.qbittorrent.qBittorrent/data/qBittorrent/logs/qbittorrent.log"
const listenOn = ":5151"
const doneDir = root+"DONE/"

type LogWriter struct {
	Buffer bytes.Buffer
}

func (lw *LogWriter) Write(p []byte) (int, error) {
	//Write to internal buffer
	n, err := lw.Buffer.Write(p)
	if err != nil {
		return n, err
	}

	//Also write to stderr
	return os.Stderr.Write(p)
}

func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}
	defer watcher.Close()

	//Send log to a buffer
	lw := LogWriter{}
	log.SetOutput(&lw)

	//Recursive function to add watches to directories
	var addWatch func(string) error
	addWatch = func(dir string) error {
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
				if err := addWatch(subdir); err != nil {
					log.Println("Error adding watch to subdirectory:", err)
				}
			}
		}
		return nil
	}

	//Start watching from the root directory
	if err := addWatch(root); err != nil {
		log.Fatal("Error watching root directory: " + err.Error())
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
//				if event.Op&(fsnotify.Chmod|fsnotify.Write) != 0 { //Ignore chmod and write
//					return
//				}
				log.Println("Event:", event)

				//Only handle create events
				if event.Op&(fsnotify.Create|fsnotify.Rename) == 0 {
					break
				}

				//Stat the new object
				info, err := os.Stat(event.Name)
				if err != nil {
					log.Println("Error statting:", err)
					break
				}

				//Handle watching new directory
				if info.IsDir() {
					if err := addWatch(event.Name); err != nil {
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
	}()

	//Listen for SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	//Http handler func
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/favicon.ico" {
			return
		}
		if r.RequestURI != "/?AlreadyLogged" {
			log.Println("Log request: " + r.RemoteAddr)
		}
		w.Write([]byte("<html><head><meta http-equiv=refresh content=5;url=./?AlreadyLogged></head><body>"))
		w.Write([]byte(strings.ReplaceAll(html.EscapeString(lw.Buffer.String()), "\n", "<br>\n")))
		w.Write([]byte("\n<div id=End></div><script>window.addEventListener('DOMContentLoaded', () => setTimeout(() => window.scrollTo(0, document.getElementById('End').offsetTop), 100));</script></body></html>"))
	}

	//Start the http server for log output
	srv := &http.Server{
		Addr:         listenOn,
		Handler:      http.HandlerFunc(handler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		log.Println("Starting server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Server closed ("+err.Error()+")")
		}
		log.Println()
	}()

	//Close the server
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Attempt graceful shutdown
	log.Println("Shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

func ProcessFile(eventPath string) string {
	//Get the path parts of the new file
	newFilePath := strings.Split(eventPath[len(root):], "/")

	//If not a valid path, return an error
	if !slices.Contains([]string{"Dir1", "Dir2"}, newFilePath[0]) {
		return "Ignoring invalid top level directory: " + newFilePath[0]
	}
	if len(newFilePath) < 3 {
		return "Must have at least 3 parts in the path"
	}

	//Ignore files named DoNotMove or that contain .crdownload or end with .tmp
	newFile := newFilePath[len(newFilePath)-1]
	if strings.ToLower(newFile) == "donotmove" || strings.Contains(newFile, ".crdownload") || strings.HasSuffix(newFile, ".tmp") {
		return "Ignoring file: " + newFile
	}

	//Create the new directory
	newDir := newRoot + strings.Join(newFilePath[0:len(newFilePath)-1], "/") + "/"
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

	//Remove the file and parent directories (only warn if it cannot be removed)
	if err := os.Rename(eventPath, doneDir+newFile); err != nil {
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

	//Send the torrent to qbittorrent
	stdOut, err := os.CreateTemp("", "Stdout")
	if err != nil {
		return "Error creating temp stdOut: " + err.Error()
	}
	stdErr, err := os.CreateTemp("", "stdErr")
	if err != nil {
		return "Error creating temp stdErr: " + err.Error()
	}
	params := []string{
		"/usr/bin/flatpak",
		"run", "--branch=stable", "--arch=x86_64", "--command=qbittorrent",
		"org.qbittorrent.qBittorrent",
		newDir + newFile,
		"--save-path=" + newDir,
		"--skip-dialog=true",
	}
	Proc, err := os.StartProcess(
		"/usr/bin/flatpak",
		params,
		&os.ProcAttr{
			"",
			nil,
			[]*os.File{nil, stdOut, stdErr},
			nil,
		},
	)
	if err != nil {
		return "Error from qbittorrent: " + err.Error()
	}

	//Output the results
	log.Println("Running: ", strings.Join(quoteParams(params), " "))
	Proc.Wait()
	log.Println("STDOUT: ", returnFirst(io.ReadAll(stdOut)))
	log.Println("STDERR: ", returnFirst(io.ReadAll(stdErr)))
	stdOut.Close()
	stdErr.Close()

	//Output the last 5 lines from the log file
	go func() {
		time.Sleep(5 * time.Second) //Pause for 5 seconds to make sure the action has completed
		log.Println("--------------------------------------------------------------------------")
		if output, err := exec.Command("tail", "-n", "5", logFile).Output(); err != nil {
			log.Println("Error executing tail log command: ", err)
		} else {
			log.Println(string(output))
		}
		log.Println("--------------------------------------------------------------------------")
	}()

	//Check every 5 seconds for a non torrent file
	go func() {
		for {
			files, err := os.ReadDir(newDir)
			if err != nil {
				log.Println("Error while waiting for file: ", err)
				return
			}
			for _, file := range files {
				if !strings.HasSuffix(strings.ToLower(file.Name()), ".torrent") {
					log.Println("Found file: ", file.Name())
					return
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()

	return "Success"
}

func quoteParams(params []string) []string {
	quoted := make([]string, len(params))
	for i, param := range params {
		quoted[i] = strconv.Quote(param)
	}
	return quoted
}

func returnFirst[T any](a T, _ interface{}) T {
	return a
}