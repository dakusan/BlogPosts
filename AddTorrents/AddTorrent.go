package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/fsnotify/fsnotify"
	"github.com/tidwall/jsonc"
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
	DoneDir            string `json:"doneDir"` //optional; empty means delete the original file instead of moving it here
	ListenOn           string `json:"listenOn"`
	QbitURL            string `json:"qbitURL"`
	QbitUser           string `json:"qbitUser"`
	QbitPass           string `json:"qbitPass"`
	FileRegex          string `json:"fileRegex"`
	FileRegexBlacklist string `json:"fileRegexBlacklist"` //optional; empty means no blacklist filtering
}

var cfg Config
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

//DefaultMaxLogLines is the default capacity of a LogRingBuffer (override via -maxLogLines).
const DefaultMaxLogLines = 1000

//LogRingBuffer is a log.SetOutput target that keeps the last MaxLogLines log lines, each numbered by a never-wrapping, ever-increasing line number (0-indexed; -1 is the client-facing sentinel for “nothing yet”).
//A line’s storage slot is just lineNumber % MaxLogLines.
type LogRingBuffer struct {
	LockMu         sync.Mutex
	Lines          []string
	NextLineNumber int //count of lines written so far; also the line number the next Write will receive
}

//Write implements io.Writer. Each call corresponds to exactly one log.Println/Printf/etc statement
//The log package’s Output always does exactly one Write per statement, even if the message contains embedded newlines), so no line-splitting is needed: one Write is one ring-buffer entry.
func (lb *LogRingBuffer) Write(P []byte) (int, error) {
	lb.LockMu.Lock()
	lb.Lines[lb.NextLineNumber%len(lb.Lines)] = strings.TrimSuffix(string(P), "\n")
	lb.NextLineNumber++
	lb.LockMu.Unlock()

	//Also write to stderr
	return os.Stderr.Write(P)
}

//LinesSince returns every retained log line after AfterLine (AfterLine == -1 means “from the oldest retained line”) as plain text
//The client is responsible for safely rendering it (e.g. via textContent) plus the line number the caller should remember as its new AfterLine for the next call.
//Callers are responsible for passing -1 instead of a stale AfterLine when they’ve detected (via ServerStartTime) that it’s no longer valid.
func (lb *LogRingBuffer) LinesSince(AfterLine int) (Lines []string, NewAfterLine int) {
	lb.LockMu.Lock()
	defer lb.LockMu.Unlock()

	MaxLogLines := len(lb.Lines)
	OldestAvailable := max(0, lb.NextLineNumber-MaxLogLines)
	Start := max(AfterLine+1, OldestAvailable)

	//Return nil if no lines to send
	Count := lb.NextLineNumber - Start
	if Count <= 0 {
		return nil, lb.NextLineNumber - 1
	}

	//The result is built with at most two slice copies (one if the wanted range doesn’t wrap around the end of the backing array, two if it does) rather than growing it one line at a time.
	Lines = make([]string, Count)
	StartIdx := Start % MaxLogLines
	if FirstPart := MaxLogLines - StartIdx; FirstPart < Count {
		//The wanted range wraps around the end of the backing array
		copy(Lines, lb.Lines[StartIdx:])
		copy(Lines[FirstPart:], lb.Lines[:Count-FirstPart])
	} else {
		copy(Lines, lb.Lines[StartIdx:StartIdx+Count])
	}
	return Lines, lb.NextLineNumber - 1
}

//ServerStartTime identifies this particular process run (set once, at startup).
//The client echoes back whatever value it was last given; if it doesn’t match, the client’s LastLogLine is from a since-gone process
var ServerStartTime = time.Now().Unix()

//UpdatesRequest is the JSON body the /Updates endpoint expects. ServerStartTime is 0 for a client that has never gotten a real value from the server yet (e.g. a fresh page load) - that’s a normal case, not a mismatch, since there’s nothing to have gone stale.
type UpdatesRequest struct {
	LastLogLine     int   `json:"LastLogLine"`
	ServerStartTime int64 `json:"ServerStartTime"`
}

//UpdatesResponse is the JSON the /Updates endpoint returns on success
type UpdatesResponse struct {
	NewLogLines      []string `json:"NewLogLines"`
	CurrentDownloads []string `json:"CurrentDownloads"`
	NewLastLogLine   int      `json:"NewLastLogLine"`
	ServerStartTime  int64    `json:"ServerStartTime"`
}

//UpdatesErrorResponse is the JSON the /Updates endpoint returns if the request itself was invalid
type UpdatesErrorResponse struct {
	Error string `json:"Error"`
}

//LogPageHTML is the static page served at “/”, embedded into the binary at compile time from logPage.html.
//go:embed logPage.html
var LogPageHTML string

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

//qbitTorrentInfo mirrors the fields we need from the qBittorrent WebUI’s /api/v2/torrents/info response
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

//qbitDeleteTag removes tag from qBittorrent entirely. Each tag is unique to the one torrent it was
//created for, so this doesn’t affect anything else; it just keeps qBittorrent’s tag list from
//accumulating one orphaned entry per torrent ever added forever.
func qbitDeleteTag(tag string) error {
	resp, err := qbitDo(func() (*http.Request, error) {
		req, err := http.NewRequest("POST", cfg.QbitURL+"/api/v2/torrents/deleteTags", strings.NewReader(url.Values{"tags": {tag}}.Encode()))
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
		return fmt.Errorf("deleteTags rejected (status %d): %s", resp.StatusCode, string(body))
	}
	return nil
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

//formatEta renders a qBittorrent eta (seconds, or 8640000 for “unknown”) as a human-readable duration
func formatEta(seconds int64) string {
	if seconds < 0 || seconds >= 8640000 {
		return "unknown"
	}
	return (time.Duration(seconds) * time.Second).Round(time.Second).String()
}

//pendingUpdate is the outcome of checking one pending torrent’s progress, applied back to the live list by tag
type pendingUpdate struct {
	remove bool
	misses int
}

//pendingStatusLines queries qBittorrent for the live progress/ETA of every pending torrent, dropping any that have finished or gone stale
func pendingStatusLines() []string {
	pendingMu.Lock()
	snapshot := append([]pendingTorrent(nil), pending...)
	pendingMu.Unlock()

	var lines []string
	updates := make(map[string]pendingUpdate, len(snapshot))
	for _, p := range snapshot {
		info, err := qbitTorrentByTag(p.tag)
		if err != nil || info == nil {
			misses := p.misses + 1
			if misses < 5 {
				updates[p.tag] = pendingUpdate{misses: misses}
				lines = append(lines, fmt.Sprintf("Waiting: %s - status unknown", p.name))
			} else {
				updates[p.tag] = pendingUpdate{remove: true}
				log.Println("Gave up tracking progress for: ", p.name)
				if err := qbitDeleteTag(p.tag); err != nil {
					log.Println("Warning: could not delete tag for “"+p.name+"”: ", err)
				}
			}
			continue
		}
		if info.Progress >= 1 {
			updates[p.tag] = pendingUpdate{remove: true}
			log.Println("Finished downloading: ", p.name)
			if err := qbitDeleteTag(p.tag); err != nil {
				log.Println("Warning: could not delete tag for “"+p.name+"”: ", err)
			}
			continue //finished downloading, drop it
		}
		updates[p.tag] = pendingUpdate{misses: 0}
		lines = append(lines, fmt.Sprintf("Waiting: %s - %.1f%% (ETA %s)", p.name, info.Progress*100, formatEta(info.Eta)))
	}

	//Apply updates to the live list by tag rather than overwriting it with the stale snapshot, so a
	//torrent added concurrently (while the qBittorrent lookups above were in flight) isn’t lost.
	pendingMu.Lock()
	filtered := make([]pendingTorrent, 0, len(pending))
	for _, p := range pending {
		if u, ok := updates[p.tag]; ok {
			if u.remove {
				continue
			}
			p.misses = u.misses
		}
		filtered = append(filtered, p)
	}
	pending = filtered
	pendingMu.Unlock()
	return lines
}

func main() {

	configPath := flag.String("config", configFile, "Path to the JSON config file")
	rootFlag := flag.String("root", "", "Override: root directory to watch for new files")
	newRootFlag := flag.String("newRoot", "", "Override: the directory files are copied into, and torrents are downloaded into, at the same relative path as under “root”")
	doneDirFlag := flag.String("doneDir", "", "Override: where the original file is moved after being copied. If empty, the original is deleted instead.")
	listenOnFlag := flag.String("listenOn", "", "Override: the “address:port” for the http log web UI")
	qbitURLFlag := flag.String("qbitURL", "", "Override: the qBittorrent WebUI base URL")
	qbitUserFlag := flag.String("qbitUser", "", "Override: the qBittorrent WebUI username")
	qbitPassFlag := flag.String("qbitPass", "", "Override: the qBittorrent WebUI password (WARNING: visible to other local users via ps/proc; prefer setting this in the config file instead)")
	fileRegexFlag := flag.String("fileRegex", "", "Override: a regex that a new file’s path, relative to root, must match to trigger actions.")
	fileRegexBlacklistFlag := flag.String("fileRegexBlacklist", "", "Override: a regex that, if matched, excludes an otherwise-matching file (empty disables the blacklist)")
	maxLogLinesFlag := flag.Int("maxLogLines", DefaultMaxLogLines, "Number of most-recent log lines to retain in memory; not part of config.json, mainly for testing with a small buffer")
	flag.Parse()

	if *maxLogLinesFlag <= 0 {
		log.Fatal("Error: -maxLogLines must be positive")
	}

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
		{rootFlag, &cfg.Root}, {newRootFlag, &cfg.NewRoot}, {doneDirFlag, &cfg.DoneDir}, {listenOnFlag, &cfg.ListenOn},
		{qbitURLFlag, &cfg.QbitURL}, {qbitUserFlag, &cfg.QbitUser}, {qbitPassFlag, &cfg.QbitPass},
		{fileRegexFlag, &cfg.FileRegex}, {fileRegexBlacklistFlag, &cfg.FileRegexBlacklist},
	} {
		if *o.flag != "" {
			*o.dest = *o.flag
		}
	}

	//root/newRoot are relied on as path prefixes throughout (e.g. stripping len(cfg.Root) off an event path), which silently breaks if either is missing its trailing slash; normalize it here rather than requiring every config/flag author to remember it.
	if cfg.Root != "" && !strings.HasSuffix(cfg.Root, "/") {
		cfg.Root += "/"
	}
	if cfg.NewRoot != "" && !strings.HasSuffix(cfg.NewRoot, "/") {
		cfg.NewRoot += "/"
	}

	if err := validateConfig(cfg); err != nil {
		log.Fatal("Error in config: " + err.Error())
	}

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

	//Send log to the ring buffer (also mirrored to stderr)
	LogBuf := &LogRingBuffer{Lines: make([]string, *maxLogLinesFlag)}
	log.SetOutput(LogBuf)

	//Log into the qBittorrent WebUI, then keep retrying every 15 seconds in the background if disconnected
	qbitLogin()
	go qbitMonitor()

	//Recursive function to add watches to directories
	var addWatch func(string) error
	addWatch = func(dir string) error {
		if cfg.DoneDir != "" && isWithinRoot(dir, cfg.DoneDir) {
			return nil //Never watch doneDir: files moved there must not generate events and get reprocessed
		}
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
		//Ignore favicon
		if r.URL.Path == "/favicon.ico" {
			return
		}

		//Anything not “/Updates” serves the static log page (and is logged, as a real page load)
		if r.URL.Path != "/Updates" {
			log.Println("Log request: " + r.RemoteAddr)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(LogPageHTML))
			return
		}

		//“/Updates” serves the AJAX JSON polling endpoint
		w.Header().Set("Content-Type", "application/json")
		var Req UpdatesRequest
		if Err := json.NewDecoder(r.Body).Decode(&Req); Err != nil || Req.LastLogLine < -1 {
			json.NewEncoder(w).Encode(UpdatesErrorResponse{Error: "Invalid request"})
			return
		}

		//Compile the lines
		AfterLine := Req.LastLogLine
		var NeedsReset=Req.ServerStartTime != 0 && Req.ServerStartTime != ServerStartTime //The client’s ServerStartTime doesn’t match ours so reset the current AfterLine and send an extra message to the client
		if NeedsReset {
			AfterLine = -1
		}
		Lines, NewLastLogLine := LogBuf.LinesSince(AfterLine)
		if NeedsReset {
			Lines = append([]string{time.Now().Format("2006/01/02 15:04:05") + " Server was restarted; resyncing log"}, Lines...)
		}

		json.NewEncoder(w).Encode(UpdatesResponse{
			NewLogLines:      Lines,
			CurrentDownloads: pendingStatusLines(),
			NewLastLogLine:   NewLastLogLine,
			ServerStartTime:  ServerStartTime,
		})
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
 		//ListenAndServe always returns a non-nil error; http.ErrServerClosed means Shutdown()
		//was called deliberately (e.g. SIGTERM), so that case must not be treated as fatal.
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server closed (" + err.Error() + ")")
		}
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
	//Any match error (e.g. a timeout from a pathological pattern/input) is treated as “ignore the file”, not “allow it”.
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

	//Remove the original file: delete it if no doneDir is configured, otherwise move it there (only warn if either fails)
	if cfg.DoneDir == "" {
		if err := os.Remove(eventPath); err != nil {
			log.Printf("Warning: Could not delete “%s”: %s\n", eventPath, err.Error())
		}
	} else if err := os.MkdirAll(cfg.DoneDir, 0775); err != nil {
		log.Printf("Warning: Could not create DONE directory “%s”: %s\n", cfg.DoneDir, err.Error())
	} else if err := os.Rename(eventPath, filepath.Clean(cfg.DoneDir)+"/"+newFile); err != nil {
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