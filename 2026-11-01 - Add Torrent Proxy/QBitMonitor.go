package main

import (
	"time"
)

//qbitTorrentInfo mirrors the fields we need from the qBittorrent WebUI’s /api/v2/torrents/info response
type qbitTorrentInfo struct {
	Hash          string  `json:"hash"`
	Name          string  `json:"name"`
	Progress      float64 `json:"progress"`
	Eta           int64   `json:"eta"`
	Completed     int64   `json:"completed"`
	Size          int64   `json:"size"`
	NumSeeds      int     `json:"num_seeds"`
	NumComplete   int     `json:"num_complete"`
	NumLeechs     int     `json:"num_leechs"`
	NumIncomplete int     `json:"num_incomplete"`
	DlSpeed       int64   `json:"dlspeed"`
}

//QbitManagedTag is applied to every torrent this tool adds to qBittorrent for as long as it’s still downloading
const QbitManagedTag = "TorrentProxyManaged"

//QbitWake lets ProcessFile wake qbitMonitor immediately after adding a torrent
var QbitWake = make(chan struct{}, 1)

//qbitMonitor is the sole background goroutine that talks to qBittorrent outside of the initial add in ProcessFile.
//It keeps the WebUI session alive, and while at least one torrent is pending, polls progress for all of them in a single batched request every 2 seconds.
//It idles (blocked on QbitWake) whenever nothing is pending, rather than polling an always-empty list.
func qbitMonitor() {
	Known := map[string]string{} //hash -> name, for torrents we’ve seen with progress<1; local to this goroutine, no lock needed
	var LastLoginAttempt time.Time
	for {
		if len(Known) == 0 {
			<-QbitWake //nothing pending - block until ProcessFile adds something
		}

		qbitStateMu.Lock()
		connected := qbitConnected
		qbitStateMu.Unlock()

		if !connected {
			if time.Since(LastLoginAttempt) >= 5*time.Second {
				LastLoginAttempt = time.Now()
				qbitLogin()
			}
		} else {
			refreshDownloads(Known)
		}

		select {
		case <-QbitWake:
		case <-time.After(2 * time.Second):
		}
	}
}