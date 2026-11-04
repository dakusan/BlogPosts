package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

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

//isWithinRoot reports whether path, once cleaned, is root itself or a descendant of it.
//Used as defense-in-depth so a crafted/misconfigured path can never make a copy/move land outside the intended tree.
func isWithinRoot(path, root string) bool {
	cleanRoot := filepath.Clean(root)
	cleanPath := filepath.Clean(path)
	return cleanPath == cleanRoot || strings.HasPrefix(cleanPath, cleanRoot+string(filepath.Separator))
}

//formatEta renders a qBittorrent eta (seconds, or 8640000 for “unknown”) as a human-readable duration
func formatEta(seconds int64) string {
	if seconds < 0 || seconds >= 8640000 {
		return "unknown"
	}
	return (time.Duration(seconds) * time.Second).Round(time.Second).String()
}

//formatSize renders a byte count as “#.##LL” (2 decimals, binary units: B/KB/MB/GB/TB/PB)
func formatSize(bytes int64) string {
	units := [...]string{"B", "KB", "MB", "GB", "TB", "PB"}
	size := float64(bytes)
	unit := 0
	for size >= 1024 && unit < len(units)-1 {
		size /= 1024
		unit++
	}
	return fmt.Sprintf("%.2f%s", size, units[unit])
}