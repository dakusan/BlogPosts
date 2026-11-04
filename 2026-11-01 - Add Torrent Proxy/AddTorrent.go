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
}