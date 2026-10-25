# Add Torrent Proxy

This watches a directory tree (<kbd>root</kbd>) for newly added torrent files and automatically feeds them into qBittorrent, moving files to the same location under <kbd>newRoot</kbd>. Only `.torrent` files are sent to qBittorrent.

## Settings

All settings are loaded from a JSON config file (`config.json` by default, override with -config) at startup. Each setting can also be overridden individually via a same-named command-line flag (1 or 2 dashes before the name, e.g. `-root` or `--root`), which takes precedence over the value in the config file.

The program exits immediately (logging why) if the config file is missing/invalid, or if any setting is still empty after merging in the flag overrides. See `config.default.json` for the expected fields/format, or run with -h to see the flags.

| [Config] Name         | Description |
| --------------------- | ----------- |
| <kbd>root</kbd>       | The directory to watch |
| <kbd>newRoot</kbd>    | The directory to move files to |
| <kbd>listenOn</kbd>   | The address:port for the http log interface |
| <kbd>qbitURL</kbd>    | The url to access the qBittorrent WebUI |
| <kbd>qbitUser</kbd>   | The username to access the qBittorrent WebUI |
| <kbd>qbitPass</kbd>   | The password to access the qBittorrent WebUI |

## What it does

- Logs into the qBittorrent WebUI at startup; if that fails (or later requests start failing), a background goroutine retries the login every 15 seconds.
	- “Connected”/“Lost connection” is only logged on the actual lost/regained transition, not on every retry.
- Recursively watches <kbd>root</kbd> for filesystem events using fsnotify, adding watches to any new subdirectories as they appear.
- On a new file creation/rename event, calls ProcessFile on the new path:
	- Validates the path is at least 3 levels deep under <kbd>root</kbd> and that its top-level directory is one of: Dir1 or Dir2 (Legacy code to be removed)
	- Ignores files named “DoNotMove”, containing “.crdownload”, or ending in “.tmp”.
	- Mirrors the directory structure under <kbd>newRoot</kbd>, creating directories as needed.
	- Waits 3 seconds (to let the source file finish being written), then copies the file to its new location under <kbd>newRoot</kbd>.
	- Moves (renames) the original file into <kbd>doneDir</kbd> (<kbd>root</kbd>+“DONE/”) once copied.
	- If the file is a “.torrent” file, adds it to qBittorrent via its WebUI API (<kbd>qbitURL</kbd>) with the save path set to the new directory and a unique tag so its progress can be tracked, and logs whether the add succeeded or failed.
- All log output from the current session is captured in memory (via LogWriter, which also echoes to stderr) and served over HTTP on <kbd>listenOn</kbd>. Any request returns an HTML page with:
	- The accumulated log buffer
	- The last line(s) of that page show the live percent/ETA (queried from the qBittorrent WebUI API) of any torrents that were added and have not finished downloading yet
	- Auto-refreshes every 5 seconds and auto-scrolls to the bottom
- Listens for SIGINT/SIGTERM and performs a graceful HTTP server shutdown (5 second timeout) on receipt.