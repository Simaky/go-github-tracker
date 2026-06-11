// Package ent is the Ent ORM layer: the schema lives in ./schema and the typed
// client is generated from it by the directive below.
package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
