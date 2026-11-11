package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

//qbitAddTorrent uploads the .torrent file at torrentPath to qBittorrent via the WebUI API, saving to savePath and tagging it with tag
func qbitAddTorrent(torrentPath, savePath, tag string) error {
	body, err := qbitApiCall(nil, "POST", "add", map[string]string{"savepath": savePath, "autoTMM": "false", "tags": tag, "*torrents": torrentPath}, nil)
	if err != nil {
		return err
	}
	if strings.TrimSpace(body) != "Ok." {
		return fmt.Errorf("add rejected: %s", body)
	}
	return nil
}

//qbitTorrentsByTag returns every torrent currently carrying tag, in a single request
func qbitTorrentsByTag(tag string) ([]qbitTorrentInfo, error) {
	body, err := qbitApiCall(nil, "GET", "info?tag="+url.QueryEscape(tag), nil, nil)
	if err != nil {
		return nil, err
	}

	var infos []qbitTorrentInfo
	if err := json.Unmarshal([]byte(body), &infos); err != nil {
		return nil, err
	}
	return infos, nil
}

//qbitRemoveTag removes tag from the torrent identified by hash
func qbitRemoveTag(hash, tag string) error {
	_, err := qbitApiCall(nil, "POST", "removeTags", map[string]string{"hashes": hash, "tags": tag}, nil)
	return err
}

//qbitDeleteTag deletes tag entirely - from every torrent carrying it, and from qBittorrent’s tag list itself - rather than untagging one specific torrent. Used for the one-off pending tags: since each is unique to a single QbitPrepareAddTorrent call, there’s no torrent to preserve the tag on, and leaving the tag definition behind would accumulate a never-cleaned-up entry per add.
func qbitDeleteTag(tag string) error {
	_, err := qbitApiCall(nil, "POST", "deleteTags", map[string]string{"tags": tag}, nil)
	return err
}