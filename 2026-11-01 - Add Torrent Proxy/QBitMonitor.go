package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

//qbitTorrentInfo mirrors the fields we need from the qBittorrent WebUI’s /api/v2/torrents/info response
type qbitTorrentInfo struct {
	Hash          string  `json:"hash"`
	Name          string  `json:"name"`
	Tags          string  `json:"tags"` //comma-separated
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

//QbitPendingTagMaxTries is how many refreshDownloads cycles (one per qbitMonitor tick, ~2s apart) a pending tag can go unmatched before it’s given up on
const QbitPendingTagMaxTries = 10

//QbitPendingTags maps each one-off pending tag currently awaiting confirmation to how many refreshDownloads cycles it’s gone unmatched so far. qbitResolvePendingTags (the only reader/resolver) removes an entry once it’s matched to a hash, or once it’s been given up on; 
//qbitPendingTagCounter generates each entry’s unique tag suffix; it’s a plain int rather than an atomic, since every place that touches it already holds QbitPendingTagsMu.
var (
	QbitPendingTagsMu     sync.Mutex
	QbitPendingTags       = map[string]int{} //tag -> tries so far
	qbitPendingTagCounter int
)

//QbitPrepareAddTorrent adds torrentPath to qBittorrent (saving to savePath), tagged with both QbitManagedTag and a one-off tag unique to this call, registers that tag for qbitMonitor’s regular polling to confirm later, then wakes qbitMonitor.
func QbitPrepareAddTorrent(torrentPath, savePath string) error {
	QbitPendingTagsMu.Lock()
	qbitPendingTagCounter++
	PendingTag := fmt.Sprintf("TorrentProxyPending_%d", qbitPendingTagCounter)
	QbitPendingTags[PendingTag] = 0
	QbitPendingTagsMu.Unlock()

	if err := qbitAddTorrent(torrentPath, savePath, QbitManagedTag+","+PendingTag); err != nil {
		QbitPendingTagsMu.Lock()
		delete(QbitPendingTags, PendingTag)
		QbitPendingTagsMu.Unlock()
		return err
	}

	select {
	case QbitWake <- struct{}{}:
	default:
	}
	return nil
}

//qbitResolvePendingTags cross-references Infos - the same /info response refreshDownloads just fetched for QbitManagedTag - against QbitPendingTags. A pending tag found on a torrent is resolved: its one-off tag is removed, leaving QbitManagedTag intact for normal tracking. A pending tag not found this round moves one try closer to being given up on (logged, since nothing is waiting around to be told directly).
func qbitResolvePendingTags(Infos []qbitTorrentInfo) {
	QbitPendingTagsMu.Lock()
	if len(QbitPendingTags) == 0 {
		QbitPendingTagsMu.Unlock()
		return
	}

	type resolvedTag struct{ Tag, Hash string }
	var Resolved []resolvedTag
	for _, Info := range Infos {
		for _, Tag := range strings.Split(Info.Tags, ",") {
			if Tag = strings.TrimSpace(Tag); Tag != QbitManagedTag { //every managed torrent carries this tag, so it’s never worth a QbitPendingTags lookup
				if _, Pending := QbitPendingTags[Tag]; Pending {
					Resolved = append(Resolved, resolvedTag{Tag, Info.Hash})
					delete(QbitPendingTags, Tag)
				}
			}
		}
	}

	var GivenUp []string
	for Tag, Tries := range QbitPendingTags {
		if Tries++; Tries >= QbitPendingTagMaxTries {
			GivenUp = append(GivenUp, Tag)
			delete(QbitPendingTags, Tag)
		} else {
			QbitPendingTags[Tag] = Tries
		}
	}
	QbitPendingTagsMu.Unlock()

	for _, R := range Resolved {
		go func() { //fire-and-forget: qbitMonitor’s own cycle shouldn’t block on this network call
			if err := qbitRemoveTag(R.Hash, R.Tag); err != nil {
				log.Println("Warning: could not remove pending tag “"+R.Tag+"”: ", err)
			}
		}()
	}
	for _, Tag := range GivenUp {
		log.Printf("Warning: qBittorrent never confirmed the torrent for pending tag “%s” after %d tries, giving up on it\n", Tag, QbitPendingTagMaxTries)
	}
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
		connected := qbitState == qbitStateConnected
		qbitStateMu.Unlock()

		if !connected {
			if time.Since(LastLoginAttempt) >= 5*time.Second {
				LastLoginAttempt = time.Now()
				qbitLogin()
			}
		} else if Ok := refreshDownloads(Known); Ok {
			QbitPendingTagsMu.Lock()
			PendingCount := len(QbitPendingTags)
			QbitPendingTagsMu.Unlock()
			HaveTorrents = len(Known) > 0 || PendingCount > 0 //keep polling until every pending tag has been matched to a hash or given up on - otherwise one could be left stranded with nothing left to ever call refreshDownloads again
		}

		select {
		case <-QbitWake:
			HaveTorrents = true
		case <-time.After(2 * time.Second):
		}
	}
}