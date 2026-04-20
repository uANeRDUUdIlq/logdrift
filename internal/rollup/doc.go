// Package rollup provides time-windowed aggregation of structured log lines
// grouped by a specified JSON field.
//
// Lines are ingested via Record; after each window the Rollup automatically
// calls Flush, which emits Entry values onto the Entries channel.  Consumers
// can also call Flush manually when window is zero (useful in tests or
// pipeline shutdown).
//
// Typical usage:
//
//	r := rollup.New("level", 10*time.Second)
//	defer r.Stop()
//	go func() {
//		for e := range r.Entries() {
//			fmt.Printf("%s %s=%s count=%d\n", e.Service, e.Field, e.Value, e.Count)
//		}
//	}()
//	r.Record("api", line)
package rollup
