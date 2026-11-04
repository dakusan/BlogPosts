package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/tidwall/jsonc"
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

//parseConfigAndFlags loads config.json (or -config's target), merges in any command-line flag overrides, validates the result, and compiles the regex fields - fatally exiting on any error.
//Returns the maxLogLines value, since that's the one flag main() still needs directly.
func parseConfigAndFlags() int {
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

	return *maxLogLinesFlag
}