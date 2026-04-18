// Package config loads and validates logdrift YAML configuration files.
//
// A minimal configuration file looks like:
//
//	poll_interval: 250ms
//	filters:
//	  - 'level == "error"'
//	services:
//	  - name: api
//	    path: /var/log/api.log
//	    format: pretty
//
// poll_interval controls how often log files are polled for new content.
// filters are CEL expressions evaluated against each JSON log line.
// Each service entry requires at minimum a name and a path.
package config
