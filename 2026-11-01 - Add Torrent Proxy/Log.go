package main

import (
	"os"
	"strings"
	"sync"
)

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