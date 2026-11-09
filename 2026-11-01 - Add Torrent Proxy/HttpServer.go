package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
	NewLogLines      []string      `json:"NewLogLines"`
	CurrentDownloads []DownloadRow `json:"CurrentDownloads"`
	NewLastLogLine   int           `json:"NewLastLogLine"`
	ServerStartTime  int64         `json:"ServerStartTime"`
}

//UpdatesErrorResponse is the JSON the /Updates endpoint returns if the request itself was invalid
type UpdatesErrorResponse struct {
	Error string `json:"Error"`
}

//LogPageHTML is the static page served at “/”, embedded into the binary at compile time from logPage.html.
//go:embed logPage.html
var LogPageHTML string

//runHttpServer serves the log page and the /Updates polling endpoint until SIGINT/SIGTERM, then shuts down gracefully
func runHttpServer(LogBuf *LogRingBuffer) {
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
		var NeedsReset = Req.ServerStartTime != 0 && Req.ServerStartTime != ServerStartTime //The client’s ServerStartTime doesn’t match ours so reset the current AfterLine and send an extra message to the client
		if NeedsReset {
			AfterLine = -1
		}
		Lines, NewLastLogLine := LogBuf.LinesSince(AfterLine)
		if NeedsReset {
			Lines = append([]string{time.Now().Format("2006/01/02 15:04:05") + " Server was restarted; resyncing log"}, Lines...)
		}

		DownloadsMu.Lock()
		Downloads := CurrentDownloads
		DownloadsMu.Unlock()

		json.NewEncoder(w).Encode(UpdatesResponse{
			NewLogLines:      Lines,
			CurrentDownloads: Downloads,
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