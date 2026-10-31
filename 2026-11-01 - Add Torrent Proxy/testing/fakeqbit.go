package main

//fakeqbit is a throwaway stand-in for the qBittorrent WebUI, implementing just the 4 endpoints
//AddTorrent.go actually calls. It's used by run_test.py to drive AddTorrent.go end-to-end without
//needing a real qBittorrent instance.
//
//Behavior:
//  - /api/v2/auth/login always succeeds ("Ok.") - no credential checking, this is a disposable test double.
//  - /api/v2/torrents/add reads the uploaded "torrents" file's entire content as an opaque Token
//    (the real AddTorrent.go doesn't care what's inside a .torrent file, so the test harness just
//    writes a random string there instead of real bencode data), records {SavePath, OutName, Token,
//    StartTime} keyed by the "tags" field, and schedules a timer that - after DownloadSeconds - writes
//    a fake "downloaded" file at SavePath/OutName containing Token, simulating what a real download
//    would eventually produce.
//  - /api/v2/torrents/info?tag=X reports live progress/eta computed from elapsed wall-clock time since
//    the add, so AddTorrent.go's normal polling sees a realistic ramp from 0 to 1 over DownloadSeconds.
//  - /api/v2/torrents/deleteTags removes the record (mirrors real behavior; AddTorrent.go calls this
//    once a torrent is finished or given up on).

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//FakeTorrent tracks one torrent "added" to this fake server
type FakeTorrent struct {
	SavePath  string
	OutName   string
	Token     string
	StartTime time.Time
}

var (
	TorrentsMu sync.Mutex
	Torrents   = map[string]*FakeTorrent{}
)

var DownloadSeconds float64

func main() {
	ListenOn := flag.String("listenOn", ":18080", "Address to listen on")
	flag.Float64Var(&DownloadSeconds, "downloadSeconds", 20, "Simulated seconds a fake download takes to go from 0% to 100%")
	flag.Parse()

	http.HandleFunc("/api/v2/auth/login", HandleLogin)
	http.HandleFunc("/api/v2/torrents/add", HandleAdd)
	http.HandleFunc("/api/v2/torrents/info", HandleInfo)
	http.HandleFunc("/api/v2/torrents/deleteTags", HandleDeleteTags)

	log.Println("fakeqbit listening on", *ListenOn, "- simulated download time:", DownloadSeconds, "seconds")
	log.Fatal(http.ListenAndServe(*ListenOn, nil))
}

func HandleLogin(W http.ResponseWriter, R *http.Request) {
	fmt.Fprint(W, "Ok.")
}

func HandleAdd(W http.ResponseWriter, R *http.Request) {
	if Err := R.ParseMultipartForm(10 << 20); Err != nil {
		http.Error(W, Err.Error(), http.StatusBadRequest)
		return
	}
	File, Header, Err := R.FormFile("torrents")
	if Err != nil {
		http.Error(W, Err.Error(), http.StatusBadRequest)
		return
	}
	defer File.Close()
	TokenBytes, Err := io.ReadAll(File)
	if Err != nil {
		http.Error(W, Err.Error(), http.StatusBadRequest)
		return
	}

	SavePath := R.FormValue("savepath")
	Tag := R.FormValue("tags")
	BaseName := strings.TrimSuffix(Header.Filename, filepath.Ext(Header.Filename))
	Torrent := &FakeTorrent{
		SavePath:  SavePath,
		OutName:   BaseName + ".mp4",
		Token:     strings.TrimSpace(string(TokenBytes)),
		StartTime: time.Now(),
	}

	TorrentsMu.Lock()
	Torrents[Tag] = Torrent
	TorrentsMu.Unlock()

	log.Printf("Added: tag=%s savepath=%s out=%s token=%s", Tag, SavePath, Torrent.OutName, Torrent.Token)

	time.AfterFunc(time.Duration(DownloadSeconds*float64(time.Second)), func() {
		FinishDownload(Tag, Torrent)
	})

	fmt.Fprint(W, "Ok.")
}

//FinishDownload writes the fake "downloaded" output file. Runs on its own timer goroutine, independent
//of whether/when anything ever polls /torrents/info for this tag.
func FinishDownload(Tag string, Torrent *FakeTorrent) {
	if Err := os.MkdirAll(Torrent.SavePath, 0775); Err != nil {
		log.Println("Error creating savepath for", Tag, ":", Err)
		return
	}
	OutPath := filepath.Join(Torrent.SavePath, Torrent.OutName)
	if Err := os.WriteFile(OutPath, []byte(Torrent.Token), 0664); Err != nil {
		log.Println("Error writing fake download for", Tag, ":", Err)
		return
	}
	log.Println("Finished fake download:", OutPath)
}

func HandleInfo(W http.ResponseWriter, R *http.Request) {
	Tag := R.URL.Query().Get("tag")

	TorrentsMu.Lock()
	Torrent, Found := Torrents[Tag]
	TorrentsMu.Unlock()

	W.Header().Set("Content-Type", "application/json")
	if !Found {
		json.NewEncoder(W).Encode([]any{})
		return
	}

	Elapsed := time.Since(Torrent.StartTime).Seconds()
	Progress := min(1, Elapsed/DownloadSeconds)
	Eta := int64(max(0, DownloadSeconds-Elapsed))
	json.NewEncoder(W).Encode([]map[string]any{{"progress": Progress, "eta": Eta}})
}

func HandleDeleteTags(W http.ResponseWriter, R *http.Request) {
	if Err := R.ParseForm(); Err != nil {
		http.Error(W, Err.Error(), http.StatusBadRequest)
		return
	}
	Tag := R.FormValue("tags")

	TorrentsMu.Lock()
	delete(Torrents, Tag)
	TorrentsMu.Unlock()

	log.Println("Deleted tag:", Tag)
	fmt.Fprint(W, "Ok.")
}
