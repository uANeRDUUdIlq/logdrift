// Package checkpoint provides a persistent store for file read offsets.
//
// When logdrift restarts it can resume tailing each log file from the last
// committed position rather than re-processing lines already seen.
//
// Usage:
//
//	store, err := checkpoint.New(".logdrift.checkpoint")
//	offset := store.Get("/var/log/app.log")
//	// ... seek tailer to offset ...
//	store.Set("/var/log/app.log", newOffset)
//	_ = store.Flush()
package checkpoint
