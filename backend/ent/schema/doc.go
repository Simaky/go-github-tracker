// Package schema holds the Ent entity schemas for go-github-tracker.
//
// This is a bare skeleton: there are no entities yet. To add the first one:
//
//  1. Create a schema file here, e.g. trackedrepo.go, declaring a type that
//     embeds ent.Schema (see https://entgo.io/docs/schema-def).
//  2. Run `make ent-generate` (or `go generate ./ent`) to regenerate the typed
//     client into backend/ent/.
//  3. The storage layer's auto-migrate (storage.New → client.Schema.Create) will
//     then create/upgrade the table on the next startup.
package schema
