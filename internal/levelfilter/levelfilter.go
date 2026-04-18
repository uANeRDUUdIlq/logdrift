// Package levelfilter provides log-level based filtering for structured JSON logs.
package levelfilter

import (
	"encoding/json"
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"warning": LevelWarn,
	"error": LevelError,
	"err":   LevelError,
	"fatal": LevelError,
}

// Filter drops log lines whose level is below the configured minimum.
type Filter struct {
	min      Level
	levelKey string
}

// New creates a Filter that passes lines at or above minLevel.
// levelKey is the JSON field name to inspect (e.g. "level").
// If minLevel is empty, all lines are passed.
func New(minLevel, levelKey string) *Filter {
	if levelKey == "" {
		levelKey = "level"
	}
	min, ok := levelNames[strings.ToLower(minLevel)]
	if !ok {
		min = LevelDebug
	}
	return &Filter{min: min, levelKey: levelKey}
}

// Allow returns true if the line should be forwarded.
func (f *Filter) Allow(line string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		// Non-JSON lines always pass.
		return true
	}
	raw, ok := obj[f.levelKey]
	if !ok {
		return true
	}
	s, ok := raw.(string)
	if !ok {
		return true
	}
	lvl, ok := levelNames[strings.ToLower(s)]
	if !ok {
		return true
	}
	return lvl >= f.min
}
