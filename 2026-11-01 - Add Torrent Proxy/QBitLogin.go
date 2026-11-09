package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
)

//qbitClient holds the authenticated session (cookie jar) used to talk to the qBittorrent WebUI API; guarded by qbitClientMu since it’s read/written from multiple goroutines
var (
	qbitClientMu sync.Mutex
	qbitClient   *http.Client
)

//qbitConnState is qbitState's type. qbitStateUnknown (the zero value) is distinct from qbitStateDisconnected so the very first connection attempt's outcome gets logged too, not just later transitions - otherwise a connection that's never once succeeded (bad qbitURL/credentials, qBittorrent never reachable) would log nothing at all, forever.
type qbitConnState int

const (
	qbitStateUnknown qbitConnState = iota
	qbitStateConnected
	qbitStateDisconnected
)

//qbitState tracks whether the last login/request succeeded, so transitions (lost/regained) can be logged without spamming on every retry
var (
	qbitStateMu sync.Mutex
	qbitState   qbitConnState
)

//qbitMarkConnected flips the connection state to connected, logging only if it wasn’t already connected
func qbitMarkConnected() {
	qbitStateMu.Lock()
	WasConnected := qbitState == qbitStateConnected
	qbitState = qbitStateConnected
	qbitStateMu.Unlock()
	if !WasConnected {
		log.Println("Connected to qbittorrent WebUI")
	}
}

//qbitMarkDisconnected flips the connection state to disconnected, logging only if it wasn’t already disconnected
func qbitMarkDisconnected(err error) {
	qbitStateMu.Lock()
	WasDisconnected := qbitState == qbitStateDisconnected
	qbitState = qbitStateDisconnected
	qbitStateMu.Unlock()
	if !WasDisconnected {
		log.Println("Lost connection to qbittorrent WebUI: " + err.Error())
	}
}

//qbitLogin authenticates against the qBittorrent WebUI and stores the resulting session cookie in qbitClient
func qbitLogin() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		qbitMarkDisconnected(err)
		return err
	}
	client := &http.Client{Jar: jar}

	body, err := qbitApiCall(client, "POST", "/auth/login", map[string]string{"username": cfg.QbitUser, "password": cfg.QbitPass}, []string{"password"})
	if err != nil {
		return err
	}
	if strings.TrimSpace(body) != "Ok." {
		err := fmt.Errorf("login rejected: %s", body)
		qbitMarkDisconnected(err)
		return err
	}

	qbitClientMu.Lock()
	qbitClient = client
	qbitClientMu.Unlock()
	qbitMarkConnected()
	return nil
}

//qbitDo performs an authenticated request built by makeReq, transparently logging back in and retrying once if the session expired (403)
func qbitDo(makeReq func() (*http.Request, error)) (*http.Response, error) {
	qbitClientMu.Lock()
	client := qbitClient
	qbitClientMu.Unlock()
	if client == nil {
		return nil, fmt.Errorf("not connected to qbittorrent WebUI")
	}

	req, err := makeReq()
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		qbitMarkDisconnected(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusForbidden {
		qbitMarkConnected()
		return resp, nil
	}
	resp.Body.Close()

	//Session expired; log back in and retry once
	if err := qbitLogin(); err != nil {
		return nil, err
	}
	qbitClientMu.Lock()
	client = qbitClient
	qbitClientMu.Unlock()

	req, err = makeReq()
	if err != nil {
		return nil, err
	}
	resp, err = client.Do(req)
	if err != nil {
		qbitMarkDisconnected(err)
		return nil, err
	}
	qbitMarkConnected()
	return resp, nil
}