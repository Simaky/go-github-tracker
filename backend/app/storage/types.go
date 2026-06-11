package storage

import "errors"

// ErrNotFound is the sentinel returned when a queried row does not exist. The
// app layer translates it into a domain not-found error. Add row structs (with
// `db`/ent tags) and more sentinels here as the schema grows.
var ErrNotFound = errors.New("not found")
