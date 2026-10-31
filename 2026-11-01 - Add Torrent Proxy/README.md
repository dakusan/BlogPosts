# Add Torrent Proxy

This monitors a shared directory in which an unprivileged user can save a torrent file, triggering a torrent download into an otherwise-inaccessible destination with the same relative path. Includes a web UI to monitor progress.
- Watches a directory tree (<kbd>root</kbd>) for newly added files
- Copies the files to the same relative location under <kbd>newRoot</kbd>
	- Files are restricted by whitelist/blacklist rules
- Optionally moves the file into a unified <kbd>doneDir</kbd> (deletes if blank)
- Any `.torrent` file among them is also added to qBittorrent for downloading into the relative directory under <kbd>newRoot</kbd>
- Has a web UI to monitor logs and torrent download status
- Tested on Linux and would probably work on Windows

## Compile

```bash
if [ ! -f config.json ]; then
	cp config.default.json config.json
	${EDITOR:-nano} config.json
fi

go build AddTorrent.go
```

## Run

Make sure the config file is set up first (See compile)

If the executable is already compiled: `./AddTorrent`
Otherwise: `go run AddTorrent.go`

## Settings

All settings are loaded from a JSONC config file (`config.json` by default, override with -config) at startup. Each setting can also be overridden individually via a same-named command-line flag (1 or 2 dashes before the name, e.g. `-root` or `--root`), which takes precedence over the value in the config file.

The program exits immediately (logging why) if the config file is missing/invalid, or if any setting is still empty after merging in the flag overrides. See `config.default.json` for the expected fields/format, or run with -h to see the flags.

| [Config] Name                             | Description |
| ----------------------------------------- | ----------- |
| <kbd>root</kbd>                           | Root directory to watch for new files |
| <kbd>newRoot</kbd>                        | The directory files are copied into, and torrents are downloaded into, at the same relative path as under <kbd>root</kbd> |
| <kbd>doneDir</kbd>                        | Where the original file is moved after being copied. If empty, the original is deleted instead. |
| <kbd>listenOn</kbd>                       | The “address:port” for the http log web UI |
| <kbd>qbitURL</kbd>                        | The qBittorrent WebUI base URL |
| <kbd>qbitUser</kbd>                       | The qBittorrent WebUI username |
| <kbd>qbitPass</kbd>                       | The qBittorrent WebUI password |
| <kbd>fileRegex</kbd><sup>1</sup>          | A regex that a new file’s path, relative to <kbd>root</kbd>, must match to trigger actions. |
| <kbd>fileRegexBlacklist</kbd><sup>1</sup> | A regex that, if matched, excludes an otherwise-matching file (empty disables the blacklist) |

* **<sup>1</sup>** Both regex fields use [`dlclark/regexp2`](https://github.com/dlclark/regexp2) (a pure-Go, no-cgo Perl/.NET-style engine), which is mostly PCRE compatible.
	- Matches are capped at a 2-second timeout, treated as “no match”
	- Matches are unanchored (not a full-string match) — add `^`/`$` yourself if you need one.

## What it does

- Logs into the qBittorrent WebUI at startup; if that fails (or later requests start failing), a background goroutine retries the login every 15 seconds.
	- “Connected”/“Lost connection” is only logged on the actual lost/regained transition, not on every retry.
- Recursively watches <kbd>root</kbd> for filesystem events using fsnotify, adding watches to any new subdirectories as they appear.
	- Symlinks are never followed
	- <kbd>doneDir</kbd> (if it’s inside <kbd>root</kbd>) is never watched, so files moved there can’t generate events and get reprocessed.
- On a new file creation/rename event, calls ProcessFile on the new path:
	- The path, relative to <kbd>root</kbd>, must match <kbd>fileRegex</kbd> and must not match <kbd>fileRegexBlacklist</kbd> (if set), or it’s ignored and logged accordingly.
	- Mirrors the directory structure under <kbd>newRoot</kbd>, creating directories as needed. The resulting path is re-checked to confirm it still falls under <kbd>newRoot</kbd> before anything is written.
	- Waits 3 seconds (to let the source file finish being written), then copies the file to its new location under <kbd>newRoot</kbd>.
	- Once copied, the original file is moved (renamed) into <kbd>doneDir</kbd>, or deleted if <kbd>doneDir</kbd> is empty.
	- If the file is a “.torrent” file, adds it to qBittorrent via its WebUI API (<kbd>qbitURL</kbd>) with the save path set to the new directory and a unique tag so its progress can be tracked, and logs whether the add succeeded or failed.
- Web client
	- Served over HTTP on <kbd>listenOn</kbd>
	- `GET /` serves a static page with two panels:
		- A log panel
			- Log lines are only ever appended (never cleared)
			- Auto-scrolls to the bottom on new lines only if it was already scrolled to the bottom; otherwise its scroll position is left alone
			- Takes up the rest of the screen that the downloads panel is not using
		- A downloads panel
			- Shows the live percent/ETA (queried from the qBittorrent WebUI API) of any torrents that were added and have not finished downloading yet
			- Grows to fit its content but always leaves room for at least ~3 log lines.
		- These requests are logged as normal page loads
	- The page polls `POST /Updates` every 5 seconds with the highest log line number it has
		- The server returns:
			- The new lines; the last 1000 log lines are kept in memory in a ring buffer (also mirrored to stderr)
			- The current downloading list
			- A server-start identifier which lets the server detect a restart, in which case it sends the client a notice and the new log lines.
- Listens for SIGINT/SIGTERM and performs a graceful HTTP server shutdown (5 second timeout) on receipt.

## Test

`testing/test.sh` (or `python3 testing/run_test.py` directly) runs an end-to-end test against a temporary workspace, with no real qBittorrent needed:

- Builds `AddTorrent` and `testing/fakeqbit.go` (a stand-in qBittorrent WebUI implementing just the 4 endpoints actually used).
- Starts both, then drops several fake `.torrent` files (each containing a random token) into the test root, staggered over time.
- Polls `/Updates` throughout the run, the same way a browser would (needed since “Finished downloading” is only logged/cleaned up in response to a poll).
- Waits for everything to finish, then checks that the original files ended up in the right places with the right content, and that the captured log contains every expected line in the correct order.
- Prints `PASS`/`FAIL` and exits 0/1 accordingly. Prints the URL to open in a browser partway through, so you can watch it live while it runs.

## Blog post mirror

[Add Torrent Proxy](https://www.castledragmire.com/Posts/Add_Torrent_Proxy)<br>
Date: 2026-11-01<br>
Tags: Servers, Proxy, Torrent, GoLang, JavaScript, Permissions