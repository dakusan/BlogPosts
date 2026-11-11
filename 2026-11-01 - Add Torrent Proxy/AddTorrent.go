package main

import "log"

func main() {
	maxLogLines := parseConfigAndFlags()

	//Send log to the ring buffer (also mirrored to stderr)
	LogBuf := &LogRingBuffer{Lines: make([]string, maxLogLines)}
	log.SetOutput(LogBuf)

	//Log into the qBittorrent WebUI, then hand off session keep-alive and download progress monitoring to qbitMonitor
	qbitLogin()
	go qbitMonitor()

	watcher := startFileWatcher()
	defer watcher.Close()

	runHttpServer(LogBuf)

	//On a clean shutdown, delete the managed-tag
	if err := qbitDeleteTag(QbitManagedTag); err != nil {
		log.Println("Warning: could not delete “"+QbitManagedTag+"” tag on shutdown: ", err)
	}
}