package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/fsnotify/fsnotify"
	"github.com/tidwall/jsonc"
	"html"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"
)

const configFile = "config.json"

//Config holds the settings loaded from configFile at startup
type Config struct {
	Root               string `json:"root"`
	NewRoot            string `json:"newRoot"`
	ListenOn           string `json:"listenOn"`
	QbitURL            string `json:"qbitURL"`
	QbitUser           string `json:"qbitUser"`
	QbitPass           string `json:"qbitPass"`
	FileRegex          string `json:"fileRegex"`
	FileRegexBlacklist string `json:"fileRegexBlacklist"` //optional; empty means no blacklist filtering
}

var cfg Config
var doneDir string
var fileRegex *regexp2.Regexp
var fileRegexBlacklist *regexp2.Regexp //nil when cfg.FileRegexBlacklist is empty, i.e. no blacklist

//regexMatchTimeout bounds how long a single match attempt may run, since regexp2 (unlike stdlib regexp/RE2) backtracks and a pathological pattern+input pair could otherwise hang.
const regexMatchTimeout = 2 * time.Second

//compileRegex compiles pattern with regexMatchTimeout applied
func compileRegex(pattern string) (*regexp2.Regexp, error) {
	re, err := regexp2.Compile(pattern, regexp2.None)
	if err != nil {
		return nil, err
	}
	re.MatchTimeout = regexMatchTimeout
	return re, nil
}

//loadConfig reads and parses the JSONC (JSON with // and /* */ comments allowed) config file at path, erroring if it is missing or invalid
func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := json.Unmarshal(jsonc.ToJSON(data), &c); err != nil {
		return Config{}, err
	}
	return c, nil
}

//validateConfig errors if any required field (after config file + flag overrides are merged) is still empty
func validateConfig(c Config) error {
	required := []struct{ name, value string }{
		{"root", c.Root}, {"newRoot", c.NewRoot}, {"listenOn", c.ListenOn},
		{"qbitURL", c.QbitURL}, {"qbitUser", c.QbitUser}, {"qbitPass", c.QbitPass},
		{"fileRegex", c.FileRegex},
	}
	for _, f := range required {
		if f.value == "" {
			return fmt.Errorf("missing required field %q", f.name)
		}
	}
	return nil
}

//isWithinRoot reports whether path, once cleaned, is root itself or a descendant of it.
//Used as defense-in-depth so a crafted/misconfigured path can never make a copy/move land outside the intended tree.
func isWithinRoot(path, root string) bool {
	cleanRoot := filepath.Clean(root)
	cleanPath := filepath.Clean(path)
	return cleanPath == cleanRoot || strings.HasPrefix(cleanPath, cleanRoot+string(filepath.Separator))
}

type LogWriter struct {
	Buffer bytes.Buffer
}

func (lw *LogWriter) Write(p []byte) (int, error) {
	//Write to internal buffer
	n, err := lw.Buffer.Write(p)
	if err != nil {
		return n, err
	}

	//Also write to stderr
	return os.Stderr.Write(p)
}

//qbitClient holds the authenticated session (cookie jar) used to talk to the qBittorrent WebUI API; guarded by qbitClientMu since it's read/written from multiple goroutines
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

//qbitMonitor keeps the qBittorrent WebUI session alive, retrying the login every 15 seconds while disconnected
func qbitMonitor() {
	for {
		qbitStateMu.Lock()
		connected := qbitConnected
		qbitStateMu.Unlock()

		if !connected {
			qbitLogin()
		}
		time.Sleep(15 * time.Second)
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

//qbitTorrentInfo mirrors the fields we need from the qBittorrent WebUI's /api/v2/torrents/info response
type qbitTorrentInfo struct {
	Progress float64 `json:"progress"`
	Eta      int64   `json:"eta"`
}

//qbitTorrentByTag looks up the (single) torrent with the given unique tag, returning nil if none is found
func qbitTorrentByTag(tag string) (*qbitTorrentInfo, error) {
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
	if len(infos) == 0 {
		return nil, nil
	}
	return &infos[0], nil
}

//pendingTorrent tracks a torrent added to qBittorrent whose progress is still shown on the log page
type pendingTorrent struct {
	tag    string
	name   string
	misses int //consecutive failed/not-found lookups, used to eventually give up on stale entries
}

var (
	pendingMu sync.Mutex
	pending   []pendingTorrent
)

//addPending registers a newly-added torrent so its progress is shown on the log page until it finishes
func addPending(tag, name string) {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	pending = append(pending, pendingTorrent{tag: tag, name: name})
}

//formatEta renders a qBittorrent eta (seconds, or 8640000 for "unknown") as a human-readable duration
func formatEta(seconds int64) string {
	if seconds < 0 || seconds >= 8640000 {
		return "unknown"
	}
	return (time.Duration(seconds) * time.Second).Round(time.Second).String()
}

//pendingStatusLines queries qBittorrent for the live progress/ETA of every pending torrent, dropping any that have finished or gone stale
func pendingStatusLines() []string {
	pendingMu.Lock()
	snapshot := append([]pendingTorrent(nil), pending...)
	pendingMu.Unlock()

	var lines []string
	var stillPending []pendingTorrent
	for _, p := range snapshot {
		info, err := qbitTorrentByTag(p.tag)
		if err != nil || info == nil {
			p.misses++
			if p.misses < 5 {
				stillPending = append(stillPending, p)
				lines = append(lines, fmt.Sprintf("Waiting: %s - status unknown", p.name))
			} else {
				log.Println("Gave up tracking progress for: ", p.name)
			}
			continue
		}
		if info.Progress >= 1 {
			log.Println("Finished downloading: ", p.name)
			continue //finished downloading, drop it
		}
		p.misses = 0
		stillPending = append(stillPending, p)
		lines = append(lines, fmt.Sprintf("Waiting: %s - %.1f%% (ETA %s)", p.name, info.Progress*100, formatEta(info.Eta)))
	}

	pendingMu.Lock()
	pending = stillPending
	pendingMu.Unlock()
	return lines
}

func main() {

	configPath := flag.String("config", configFile, "Path to the JSON config file")
	rootFlag := flag.String("root", "", "Override: root directory to watch for new torrent files")
	newRootFlag := flag.String("newRoot", "", "Override: directory to copy completed files into")
	listenOnFlag := flag.String("listenOn", "", "Override: address the log HTTP server listens on")
	qbitURLFlag := flag.String("qbitURL", "", "Override: qBittorrent WebUI base URL")
	qbitUserFlag := flag.String("qbitUser", "", "Override: qBittorrent WebUI username")
	qbitPassFlag := flag.String("qbitPass", "", "Override: qBittorrent WebUI password")
	fileRegexFlag := flag.String("fileRegex", "", "Override: regex that a new file's path (relative to root) must match to be copied/moved")
	fileRegexBlacklistFlag := flag.String("fileRegexBlacklist", "", "Override: regex that, if matched, excludes an otherwise-matching file (empty disables the blacklist)")
	flag.Parse()

	loadedCfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatal("Error loading " + *configPath + ": " + err.Error())
	}
	cfg = loadedCfg

	//Apply command-line overrides
	for _, o := range []struct {
		flag *string
		dest *string
	}{
		{rootFlag, &cfg.Root}, {newRootFlag, &cfg.NewRoot}, {listenOnFlag, &cfg.ListenOn},
		{qbitURLFlag, &cfg.QbitURL}, {qbitUserFlag, &cfg.QbitUser}, {qbitPassFlag, &cfg.QbitPass},
		{fileRegexFlag, &cfg.FileRegex}, {fileRegexBlacklistFlag, &cfg.FileRegexBlacklist},
	} {
		if *o.flag != "" {
			*o.dest = *o.flag
		}
	}

	if err := validateConfig(cfg); err != nil {
		log.Fatal("Error in config: " + err.Error())
	}
	doneDir = cfg.Root + "DONE/"

	fileRegex, err = compileRegex(cfg.FileRegex)
	if err != nil {
		log.Fatal("Error in config: invalid fileRegex: " + err.Error())
	}
	if cfg.FileRegexBlacklist != "" {
		fileRegexBlacklist, err = compileRegex(cfg.FileRegexBlacklist)
		if err != nil {
			log.Fatal("Error in config: invalid fileRegexBlacklist: " + err.Error())
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}
	defer watcher.Close()

	//Send log to a buffer
	lw := LogWriter{}
	log.SetOutput(&lw)

	//Log into the qBittorrent WebUI, then keep retrying every 15 seconds in the background if disconnected
	qbitLogin()
	go qbitMonitor()

	//Recursive function to add watches to directories
	var addWatch func(string) error
	addWatch = func(dir string) error {
		if slices.Contains(watcher.WatchList(), dir) {
			return nil //Already watching this directory
		}

		if err := watcher.Add(dir); err != nil {
			return err
		}
		log.Println("Started watching:", dir)

		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				subdir := filepath.Join(dir, entry.Name())
				if err := addWatch(subdir); err != nil {
					log.Println("Error adding watch to subdirectory:", err)
				}
			}
		}
		return nil
	}

	//Start watching from the root directory
	if err := addWatch(cfg.Root); err != nil {
		log.Fatal("Error watching root directory: " + err.Error())
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
//				if event.Op&(fsnotify.Chmod|fsnotify.Write) != 0 { //Ignore chmod and write
//					return
//				}
				log.Println("Event:", event)

				//Only handle create events
				if event.Op&(fsnotify.Create|fsnotify.Rename) == 0 {
					break
				}

				//Stat the new object. Lstat (not Stat) is used deliberately so a symlink is never followed: this keeps the watched/moved tree from being able to escape root via a symlink planted inside it.
				info, err := os.Lstat(event.Name)
				if err != nil {
					log.Println("Error statting:", err)
					break
				}

				//Handle watching new directory
				if info.IsDir() {
					if err := addWatch(event.Name); err != nil {
						log.Println("Error adding watch to new directory:", err)
					}
					break
				}

				//Make sure the new file is a regular file
				if !info.Mode().IsRegular() {
					break
				}

				//Process the new file
				log.Printf("Return for “%s”: %s\n", event.Name, ProcessFile(event.Name))

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	//Listen for SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	//Http handler func
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/favicon.ico" {
			return
		}
		if r.RequestURI != "/?AlreadyLogged" {
			log.Println("Log request: " + r.RemoteAddr)
		}
		w.Write([]byte("<html><head><meta http-equiv=refresh content=5;url=./?AlreadyLogged></head><body>"))
		w.Write([]byte(strings.ReplaceAll(html.EscapeString(lw.Buffer.String()), "\n", "<br>\n")))
		for _, line := range pendingStatusLines() {
			w.Write([]byte("<br>\n" + html.EscapeString(line)))
		}
		w.Write([]byte("\n<div id=End></div><script>window.addEventListener('DOMContentLoaded', () => setTimeout(() => window.scrollTo(0, document.getElementById('End').offsetTop), 100));</script></body></html>"))
	}

	//Start the http server for log output
	srv := &http.Server{
		Addr:         cfg.ListenOn,
		Handler:      http.HandlerFunc(handler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		log.Println("Starting server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Server closed ("+err.Error()+")")
		}
		log.Println()
	}()

	//Close the server
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Attempt graceful shutdown
	log.Println("Shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

func ProcessFile(eventPath string) string {
	//Get the path of the new file, relative to root, and require it to match fileRegex (and not fileRegexBlacklist).
	//Any match error (e.g. a timeout from a pathological pattern/input) is treated as "ignore the file", not "allow it".
	relPath := eventPath[len(cfg.Root):]
	if matched, err := fileRegex.MatchString(relPath); err != nil {
		return "Error evaluating fileRegex, ignoring: " + err.Error()
	} else if !matched {
		return "File does not match regex, ignoring: " + relPath
	}
	if fileRegexBlacklist != nil {
		if matched, err := fileRegexBlacklist.MatchString(relPath); err != nil {
			return "Error evaluating fileRegexBlacklist, ignoring: " + err.Error()
		} else if matched {
			return "File matches blacklist regex, ignoring: " + relPath
		}
	}
	newFilePath := strings.Split(relPath, "/")
	newFile := newFilePath[len(newFilePath)-1]

	//Create the new directory (isWithinRoot guards against a crafted/misconfigured path ever landing outside newRoot)
	newDir := filepath.Clean(cfg.NewRoot+strings.Join(newFilePath[0:len(newFilePath)-1], "/")) + "/"
	if !isWithinRoot(newDir, cfg.NewRoot) {
		return "Refusing to write outside newRoot: " + newDir
	}
	if err := os.MkdirAll(newDir, 0775); err != nil {
		return "Could not create new path: " + err.Error()
	}

	//Copy the file
	time.Sleep(3 * time.Second) //Pause for 3 seconds to make sure the file finishes copying
	source, err := os.Open(eventPath)
	if err != nil {
		return "Could not open file for moving: " + err.Error()
	}
	defer source.Close()
	destination, err := os.Create(newDir + newFile)
	if err != nil {
		return "Could not create new file: " + err.Error()
	}
	defer destination.Close()
	if _, err := io.Copy(destination, source); err != nil {
		return "Could not finish copying file: " + err.Error()
	}
	source.Close()
	log.Println("File copied to: ", newDir+newFile)

	//Remove the file and parent directories (only warn if it cannot be removed)
	if err := os.Rename(eventPath, doneDir+newFile); err != nil {
		log.Printf("Warning: Could not move “%s” to DONE directory: %s\n", eventPath, err.Error())
	}

	/*
	for i := len(newFilePath); i > 1; i-- {
		pathToRemove := root + strings.Join(newFilePath[0:i], "/")
		if err := os.Remove(pathToRemove); err != nil {
			log.Printf("Warning: Could not remove “%s”: %s\n", pathToRemove, err.Error())
			break
		} else {
			log.Println("Removed: " + pathToRemove)
		}
		break //I decided to only remove the file to move and leave the directories in tact
	}
	*/

	//If the file is a torrent, then start it in qbittorrent
	if !strings.HasSuffix(strings.ToLower(newFile), ".torrent") {
		return "Ignoring non-torrent file"
	}

	//Send the torrent to qbittorrent via the WebUI API
	tag := fmt.Sprintf("AddTorrent-%d", time.Now().UnixNano())
	if err := qbitAddTorrent(newDir+newFile, newDir, tag); err != nil {
		return "Error adding torrent to qbittorrent: " + err.Error()
	}
	log.Println("Added to qbittorrent: ", newDir+newFile)
	addPending(tag, newFile)

	return "Success"
}