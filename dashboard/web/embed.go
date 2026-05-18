// Package web bundles the dashboard's static UI into the binary via embed
// so the runtime image only needs to ship the Go binary.
package web

import "embed"

//go:embed all:templates
var FS embed.FS
