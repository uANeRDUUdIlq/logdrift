// Package multiplexer provides fan-in functionality for logdrift.
//
// It accepts a map of named line channels (one per service) and merges them
// into a single stream of [Entry] values, each tagged with its originating
// service name. Consumers read from the single Out() channel and can pass
// entries to the formatter or filter packages without caring about which
// underlying tailer produced them.
package multiplexer
