//go:build tools

// Package tools pins the code-generation tools the build depends on so that
// `go mod tidy` keeps them in go.mod. It is never compiled into the binary.
package tools

import _ "entgo.io/ent/cmd/ent"
