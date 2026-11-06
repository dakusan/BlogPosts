package main

import (
	"fmt"
	"log"
	"sync"
)

//DownloadRow is one row of the Downloads table shown on the log page
type DownloadRow struct {
	Name     string `json:"Name"`
	Progress string `json:"Progress"`
	Eta      string `json:"Eta"`
	Size     string `json:"Size"`
	Conns    string `json:"Conns"`
	Speed    string `json:"Speed"`
}

//CurrentDownloads is read by the /Updates handler and written only by refreshDownloads (from qbitMonitor), so
//handlers never block on - or trigger - a qBittorrent request themselves.
var (
	DownloadsMu      sync.Mutex
	CurrentDownloads []DownloadRow
)

//refreshDownloads makes one batched query for every torrent carrying QbitManagedTag and updates Known (hash -> name, mutated in place) to reflect what’s still pending, and publishes the resulting rows for /Updates to read.
//Ok is false only if the query itself failed (Known/CurrentDownloads are left untouched in that case); the caller (qbitMonitor) relies on this to tell “confirmed nothing pending” apart from “couldn’t check”.
func refreshDownloads(Known map[string]string) (Ok bool) {
	Infos, err := qbitTorrentsByTag(QbitManagedTag)
	if err != nil {
		log.Println("Warning: could not query qbittorrent for managed torrents: ", err)
		return false //leave Known/CurrentDownloads untouched; retried on the next tick
	}

	Seen := make(map[string]bool, len(Infos))
	var Rows []DownloadRow
	for _, Info := range Infos {
		Seen[Info.Hash] = true
		if Info.Progress >= 1 {
			log.Println("Finished downloading: ", Info.Name)
			if err := qbitRemoveTag(Info.Hash, QbitManagedTag); err != nil {
				log.Println("Warning: could not remove tag for “"+Info.Name+"”: ", err)
			}
			delete(Known, Info.Hash)
			continue //finished downloading, drop it
		}
		Known[Info.Hash] = Info.Name
		Rows = append(Rows, DownloadRow{
			Name:     Info.Name,
			Progress: fmt.Sprintf("%.1f%%", Info.Progress*100),
			Eta:      formatEta(Info.Eta),
			Size:     formatSize(Info.Completed) + "/" + formatSize(Info.Size),
			Conns:    fmt.Sprintf("%d(%d)/%d(%d)", Info.NumSeeds, Info.NumComplete, Info.NumLeechs, Info.NumIncomplete),
			Speed:    formatSize(Info.DlSpeed) + "/s",
		})
	}

	//Anything we were tracking that didn’t come back this time was removed from qBittorrent by something other than us
	for Hash, Name := range Known {
		if !Seen[Hash] {
			log.Println("Gave up tracking (no longer found in qbittorrent): ", Name)
			delete(Known, Hash)
		}
	}

	DownloadsMu.Lock()
	CurrentDownloads = Rows
	DownloadsMu.Unlock()
	return true
}