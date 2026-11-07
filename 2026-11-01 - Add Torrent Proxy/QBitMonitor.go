package main

import (
	"sync/atomic"
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

//QbitPendingAdds counts torrents ProcessFile has told qBittorrent to add but qbitMonitor hasn’t yet seen come back from a tag query - qBittorrent can take a few seconds to index a freshly added torrent, so a query run right after an add can legitimately find nothing yet. Incremented by ProcessFile, decremented by qbitMonitor.
var QbitPendingAdds atomic.Int64

//qbitMonitor is the sole background goroutine that talks to qBittorrent outside of the initial add in ProcessFile.
//It keeps the WebUI session alive, and while at least one torrent is pending, polls progress for all of them in a single batched request every 2 seconds.
//It idles (blocked on QbitWake) whenever nothing is pending, rather than polling an always-empty list.
func qbitMonitor() {
	Known := map[string]string{} //hash -> name, for torrents we’ve seen with progress<1; local to this goroutine, no lock needed

	HaveTorrents := false
	StaleCycles := 0 //consecutive confirmed-successful refreshes that found no newly-seen torrent while QbitPendingAdds was still >0; capped so an add that never shows up (e.g. silently rejected) can’t block idling forever
	var LastLoginAttempt time.Time
	for {
		if !HaveTorrents {
			<-QbitWake //nothing pending - block until ProcessFile adds something
			HaveTorrents = true
		}

		qbitStateMu.Lock()
		connected := qbitConnected
		qbitStateMu.Unlock()

		if !connected {
			if time.Since(LastLoginAttempt) >= 5*time.Second {
				LastLoginAttempt = time.Now()
				qbitLogin()
			}
		} else if Ok, NewlySeen := refreshDownloads(Known); Ok {
			if NewlySeen > 0 {
				if QbitPendingAdds.Add(-int64(NewlySeen)) < 0 {
					QbitPendingAdds.Store(0)
				}
				StaleCycles = 0
			} else if QbitPendingAdds.Load() > 0 {
				if StaleCycles++; StaleCycles >= 10 {
					QbitPendingAdds.Store(0)
					StaleCycles = 0
				}
			}
			HaveTorrents = len(Known) > 0 || QbitPendingAdds.Load() > 0
		}

		select {
		case <-QbitWake:
			HaveTorrents = true
		case <-time.After(2 * time.Second):
		}
	}
}