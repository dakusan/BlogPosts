package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

//qbitAddTorrent uploads the .torrent file at torrentPath to qBittorrent via the WebUI API, saving to savePath and tagging it with tag
func qbitAddTorrent(torrentPath, savePath, tag string) error {
	resp, err := qbitDo(func() (*http.Request, error) {
		file, err := os.Open(torrentPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		part, err := mw.CreateFormFile("torrents", filepath.Base(torrentPath))
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(part, file); err != nil {
			return nil, err
		}
		mw.WriteField("savepath", savePath)
		mw.WriteField("autoTMM", "false")
		mw.WriteField("tags", tag)
		if err := mw.Close(); err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", cfg.QbitURL+"/api/v2/torrents/add", &buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mw.FormDataContentType())
		return req, nil
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || strings.TrimSpace(string(body)) != "Ok." {
		return fmt.Errorf("add rejected (status %d): %s", resp.StatusCode, string(body))
	}
	return nil
}

//qbitTorrentsByTag returns every torrent currently carrying tag, in a single request
func qbitTorrentsByTag(tag string) ([]qbitTorrentInfo, error) {
	resp, err := qbitDo(func() (*http.Request, error) {
		return http.NewRequest("GET", cfg.QbitURL+"/api/v2/torrents/info?tag="+url.QueryEscape(tag), nil)
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var infos []qbitTorrentInfo
	if err := json.NewDecoder(resp.Body).Decode(&infos); err != nil {
		return nil, err
	}
	return infos, nil
}

//qbitRemoveTag removes tag from the torrent identified by hash
func qbitRemoveTag(hash, tag string) error {
	resp, err := qbitDo(func() (*http.Request, error) {
		req, err := http.NewRequest("POST", cfg.QbitURL+"/api/v2/torrents/removeTags", strings.NewReader(url.Values{"hashes": {hash}, "tags": {tag}}.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return req, nil
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("removeTags rejected (status %d): %s", resp.StatusCode, string(body))
	}
	return nil
}