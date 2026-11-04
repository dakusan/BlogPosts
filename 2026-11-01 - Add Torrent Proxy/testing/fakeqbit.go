package main

//fakeqbit is a throwaway stand-in for the qBittorrent WebUI, implementing just the endpoints
//AddTorrent.go actually calls. It's used by run_test.py to drive AddTorrent.go end-to-end without
//needing a real qBittorrent instance.
//
//Behavior:
//  - /api/v2/auth/login always succeeds ("Ok.") - no credential checking, this is a disposable test double.
//  - /api/v2/torrents/add reads the uploaded "torrents" file's entire content as an opaque Token
//    (the real AddTorrent.go doesn't care what's inside a .torrent file, so the test harness just
//    writes a random string there instead of real bencode data), generates a fake Hash for it (mirroring
//    a real info-hash), records it keyed by that hash, and schedules a timer that - after DownloadSeconds -
//    writes a fake "downloaded" file at SavePath/OutName containing Token, simulating what a real download
//    would eventually produce. Name is just the uploaded filename (e.g. "T1.torrent"), standing in for
//    whatever title a real torrent's own metadata would report.
//  - /api/v2/torrents/info?tag=X reports every torrent currently carrying tag X (no "tag" param means
//    "any tag", matching real qBittorrent), with live progress/eta computed from elapsed wall-clock time
//    since the add. size/completed/dlspeed are derived from SimulatedTotalSize; num_seeds/num_complete/
//    num_leechs/num_incomplete are fixed placeholder values - just enough to exercise the Downloads table.
//  - /api/v2/torrents/removeTags removes the given tags from the given hashes only (not a global tag
//    deletion), mirroring real qBittorrent's per-torrent untagging.

import (
	"crypto/rand"
	"encoding/hex"
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
	Hash      string
	Name      string
	SavePath  string
	OutName   string
	Token     string
	StartTime time.Time
	Tags      map[string]bool
}

var (
	TorrentsMu sync.Mutex
	Torrents   = map[string]*FakeTorrent{} //keyed by Hash
)

var DownloadSeconds float64

//SimulatedTotalSize is the fake torrent size (bytes) used to synthesize the size/completed/dlspeed fields
const SimulatedTotalSize int64 = 500 * 1024 * 1024

func main() {
	ListenOn := flag.String("listenOn", ":18080", "Address to listen on")
	flag.Float64Var(&DownloadSeconds, "downloadSeconds", 20, "Simulated seconds a fake download takes to go from 0% to 100%")
	flag.Parse()

	http.HandleFunc("/api/v2/auth/login", HandleLogin)
	http.HandleFunc("/api/v2/torrents/add", HandleAdd)
	http.HandleFunc("/api/v2/torrents/info", HandleInfo)
	http.HandleFunc("/api/v2/torrents/removeTags", HandleRemoveTags)

	log.Println("fakeqbit listening on", *ListenOn, "- simulated download time:", DownloadSeconds, "seconds")
	log.Fatal(http.ListenAndServe(*ListenOn, nil))
}

func HandleLogin(W http.ResponseWriter, R *http.Request) {
	fmt.Fprint(W, "Ok.")
}

//NewFakeHash generates a fake 40-hex-char info-hash, the same shape a real torrent's SHA1 hash would have
func NewFakeHash() string {
	Bytes := make([]byte, 20)
	rand.Read(Bytes)
	return hex.EncodeToString(Bytes)
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
	Tags := map[string]bool{}
	for _, Tag := range strings.Split(R.FormValue("tags"), ",") {
		if Tag = strings.TrimSpace(Tag); Tag != "" {
			Tags[Tag] = true
		}
	}
	BaseName := strings.TrimSuffix(Header.Filename, filepath.Ext(Header.Filename))
	Torrent := &FakeTorrent{
		Hash:      NewFakeHash(),
		Name:      Header.Filename,
		SavePath:  SavePath,
		OutName:   BaseName + ".mp4",
		Token:     strings.TrimSpace(string(TokenBytes)),
		StartTime: time.Now(),
		Tags:      Tags,
	}

	TorrentsMu.Lock()
	Torrents[Torrent.Hash] = Torrent
	TorrentsMu.Unlock()

	log.Printf("Added: hash=%s tags=%s savepath=%s out=%s token=%s", Torrent.Hash, R.FormValue("tags"), SavePath, Torrent.OutName, Torrent.Token)

	time.AfterFunc(time.Duration(DownloadSeconds*float64(time.Second)), func() {
		FinishDownload(Torrent)
	})

	fmt.Fprint(W, "Ok.")
}

//FinishDownload writes the fake "downloaded" output file. Runs on its own timer goroutine, independent
//of whether/when anything ever polls /torrents/info for this torrent.
func FinishDownload(Torrent *FakeTorrent) {
	if Err := os.MkdirAll(Torrent.SavePath, 0775); Err != nil {
		log.Println("Error creating savepath for", Torrent.Hash, ":", Err)
		return
	}
	OutPath := filepath.Join(Torrent.SavePath, Torrent.OutName)
	if Err := os.WriteFile(OutPath, []byte(Torrent.Token), 0664); Err != nil {
		log.Println("Error writing fake download for", Torrent.Hash, ":", Err)
		return
	}
	log.Println("Finished fake download:", OutPath)
}

func HandleInfo(W http.ResponseWriter, R *http.Request) {
	Tag := R.URL.Query().Get("tag") //empty means "any tag", matching real qBittorrent

	TorrentsMu.Lock()
	Snapshot := make([]*FakeTorrent, 0, len(Torrents))
	for _, Torrent := range Torrents {
		if Tag == "" || Torrent.Tags[Tag] {
			Snapshot = append(Snapshot, Torrent)
		}
	}
	TorrentsMu.Unlock()

	Results := make([]map[string]any, 0, len(Snapshot))
	for _, Torrent := range Snapshot {
		Elapsed := time.Since(Torrent.StartTime).Seconds()
		Progress := min(1, Elapsed/DownloadSeconds)
		Eta := int64(max(0, DownloadSeconds-Elapsed))
		Results = append(Results, map[string]any{
			"hash":           Torrent.Hash,
			"name":           Torrent.Name,
			"progress":       Progress,
			"eta":            Eta,
			"completed":      int64(Progress * float64(SimulatedTotalSize)),
			"size":           SimulatedTotalSize,
			"num_seeds":      3,
			"num_complete":   10,
			"num_leechs":     2,
			"num_incomplete": 5,
			"dlspeed":        int64(float64(SimulatedTotalSize) / DownloadSeconds),
		})
	}

	W.Header().Set("Content-Type", "application/json")
	json.NewEncoder(W).Encode(Results)
}

func HandleRemoveTags(W http.ResponseWriter, R *http.Request) {
	if Err := R.ParseForm(); Err != nil {
		http.Error(W, Err.Error(), http.StatusBadRequest)
		return
	}
	Hashes := strings.Split(R.FormValue("hashes"), ",")
	TagsToRemove := strings.Split(R.FormValue("tags"), ",")

	TorrentsMu.Lock()
	for _, Hash := range Hashes {
		if Torrent, Found := Torrents[Hash]; Found {
			for _, Tag := range TagsToRemove {
				delete(Torrent.Tags, strings.TrimSpace(Tag))
			}
		}
	}
	TorrentsMu.Unlock()

	log.Println("Removed tags", TagsToRemove, "from", Hashes)
	fmt.Fprint(W, "Ok.")
}
