// Package schema holds the Ent entity schemas — the hand-written source of
// truth from which the typed client under ent/ is generated (make ent-generate).
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// TimeMixin adds Ent-managed created_at / updated_at timestamps to a schema.
type TimeMixin struct {
	mixin.Schema
}

// Fields of the TimeMixin.
func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Repo is a GitHub repository the user has chosen to track. Its metadata mirrors
// the subset of the GitHub API we persist; notes are user-editable.
type Repo struct {
	ent.Schema
}

// Mixin attaches the shared timestamp fields.
func (Repo) Mixin() []ent.Mixin {
	return []ent.Mixin{TimeMixin{}}
}

// Fields of the Repo.
func (Repo) Fields() []ent.Field {
	return []ent.Field{
		field.String("owner").NotEmpty(),
		field.String("name").NotEmpty(),
		// full_name ("owner/name") is the natural key: one row per tracked repo.
		field.String("full_name").NotEmpty().Unique(),
		field.String("description").Optional(),
		field.Int("stars").NonNegative().Default(0),
		field.String("language").Optional(),
		field.String("html_url").NotEmpty(),
		// notes is the only user-editable field; everything else comes from GitHub.
		field.String("notes").Optional().Default(""),
		// fetched_at records when we last pulled metadata from GitHub.
		field.Time("fetched_at"),
		field.Int("forks_count").NonNegative().Default(0),
	}
}

// Indexes of the Repo.
func (Repo) Indexes() []ent.Index {
	return []ent.Index{
		// The list endpoint filters by language.
		index.Fields("language"),
	}
}
