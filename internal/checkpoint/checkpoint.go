// Package checkpoint persists the last-read byte offset for each tailed file
// so that logdrift can resume from where it left off after a restart.
package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
)

// Store holds per-file byte offsets and flushes them to a JSON file.
type Store struct {
	mu      sync.Mutex
	path    string
	offsets map[string]int64
}

// New loads an existing checkpoint file at path, or starts with an empty
// offset map if the file does not yet exist.
func New(path string) (*Store, error) {
	s := &Store{path: path, offsets: make(map[string]int64)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.offsets); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns the saved offset for the given file path, or 0 if unknown.
func (s *Store) Get(file string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.offsets[file]
}

// Set updates the in-memory offset for the given file.
func (s *Store) Set(file string, offset int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offsets[file] = offset
}

// Flush writes all offsets to disk atomically.
func (s *Store) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := json.MarshalIndent(s.offsets, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
