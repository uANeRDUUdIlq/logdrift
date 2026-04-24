// Package jsonexpand provides a pipeline stage that expands dot-notation keys
// found inside structured JSON log lines into proper nested JSON objects.
//
// # Motivation
//
// Some log producers emit flat keys like "http.method" or "db.query.duration"
// instead of a proper nested structure. jsonexpand promotes those keys into
// real nested maps so that downstream stages (jsonpath, jsonflatten, etc.) can
// reason about the hierarchy correctly.
//
// # Usage
//
//	e := jsonexpand.New(nil)          // expand every dot-key automatically
//	e := jsonexpand.New([]string{"a.b", "meta.env"}) // expand named keys only
//	output := e.Apply(inputLine)
package jsonexpand
