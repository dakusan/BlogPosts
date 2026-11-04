package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
)

//qbitClient holds the authenticated session (cookie jar) used to talk to the qBittorrent WebUI API; guarded by qbitClientMu since it’s read/written from multiple goroutines
var (
	qbitClientMu sync.Mutex
	qbitClient   *http.Client
)

//qbitConnected tracks whether the last login/request succeeded, so transitions (lost/regained) can be logged without spamming on every retry
var (
	qbitStateMu   sync.Mutex
	qbitConnected bool
)

//qbitMarkConnected flips the connection state to connected, logging only if it was previously down
func qbitMarkConnected() {
	qbitStateMu.Lock()
	wasConnected := qbitConnected
	qbitConnected = true
	qbitStateMu.Unlock()
	if !wasConnected {
		log.Println("Connected to qbittorrent WebUI")
	}
}

//qbitMarkDisconnected flips the connection state to disconnected, logging only if it was previously up
func qbitMarkDisconnected(err error) {
	qbitStateMu.Lock()
	wasConnected := qbitConnected
	qbitConnected = false
	qbitStateMu.Unlock()
	if wasConnected {
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

	resp, err := client.PostForm(cfg.QbitURL+"/api/v2/auth/login", url.Values{
		"username": {cfg.QbitUser},
		"password": {cfg.QbitPass},
	})
	if err != nil {
		qbitMarkDisconnected(err)
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if strings.TrimSpace(string(body)) != "Ok." {
		err := fmt.Errorf("login rejected (status %d): %s", resp.StatusCode, string(body))
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