package main

import (
	"sync"
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

//QbitPendingAdds is a FIFO of timestamps, one per qbitAddTorrent call ProcessFile has made that qbitMonitor hasn’t yet seen confirmed by a tag query - qBittorrent can take a few seconds to index a freshly added torrent, so a query run right after an add can legitimately find nothing yet.
//Each entry ages out independently (oldest first), so one add that never shows up (e.g. silently rejected) can’t be masked by unrelated add/confirm activity elsewhere in the queue, and a fresh add can’t inherit an older one’s accrued staleness.
var (
	QbitPendingAddsMu sync.Mutex
	QbitPendingAdds   []time.Time
)

//QbitPendingAddStaleAfter is how long an entry can sit at the front of QbitPendingAdds unconfirmed before qbitMonitor gives up on it
const QbitPendingAddStaleAfter = 20 * time.Second

//QbitAddTorrent adds torrentPath to qBittorrent (saving to savePath, tagged with tag) and registers it in QbitPendingAdds so qbitMonitor knows to expect it, then wakes qbitMonitor.
//The entry is registered before the underlying add call, not after it succeeds - otherwise qbitMonitor’s background poll could see the torrent (qBittorrent having already indexed it) and count it as confirmed before this function ever recorded that it was expecting it.
//If the call fails the entry is removed (nothing will ever confirm an add that was never actually made); if it succeeds the entry’s timestamp is reset to now, so staleness is measured from “qBittorrent was actually asked” rather than from before this (possibly slow) call was even made.
func QbitAddTorrent(torrentPath, savePath, tag string) error {
	AddedAt := time.Now()
	QbitPendingAddsMu.Lock()
	QbitPendingAdds = append(QbitPendingAdds, AddedAt)
	QbitPendingAddsMu.Unlock()

	err := qbitAddTorrent(torrentPath, savePath, tag)

	QbitPendingAddsMu.Lock()
	for i, T := range QbitPendingAdds {
		if T.Equal(AddedAt) {
			if err != nil {
				QbitPendingAdds = append(QbitPendingAdds[:i], QbitPendingAdds[i+1:]...)
			} else {
				QbitPendingAdds[i] = time.Now()
			}
			break
		}
	}
	QbitPendingAddsMu.Unlock()

	if err != nil {
		return err
	}
	select {
	case QbitWake <- struct{}{}:
	default:
	}
	return nil
}

//qbitMonitor is the sole background goroutine that talks to qBittorrent outside of the initial add in ProcessFile.
//It keeps the WebUI session alive, and while at least one torrent is pending, polls progress for all of them in a single batched request every 2 seconds.
//It idles (blocked on QbitWake) whenever nothing is pending, rather than polling an always-empty list.
func qbitMonitor() {
	Known := map[string]string{} //hash -> name, for torrents we’ve seen with progress<1; local to this goroutine, no lock needed

	HaveTorrents := false
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
			QbitPendingAddsMu.Lock()
			if NewlySeen > len(QbitPendingAdds) {
				NewlySeen = len(QbitPendingAdds)
			}
			QbitPendingAdds = QbitPendingAdds[NewlySeen:] //the NewlySeen oldest entries are assumed confirmed; which exact add matches which confirmation doesn’t matter, only that the count and remaining ages stay right
			for len(QbitPendingAdds) > 0 && time.Since(QbitPendingAdds[0]) > QbitPendingAddStaleAfter {
				QbitPendingAdds = QbitPendingAdds[1:]
			}
			HaveTorrents = len(Known) > 0 || len(QbitPendingAdds) > 0
			QbitPendingAddsMu.Unlock()
		}

		select {
		case <-QbitWake:
			HaveTorrents = true
		case <-time.After(2 * time.Second):
		}
	}
}