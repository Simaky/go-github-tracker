package storage

import (
	"github.com/Simaky/go-github-tracker/backend/app"
	"github.com/Simaky/go-github-tracker/backend/ent"
)

// toAppRepo maps a persisted *ent.Repo row into the app-layer Repo. Keeping the
// mapping here means the app layer never depends on the ORM's row type.
func toAppRepo(e *ent.Repo) *app.Repo {
	return &app.Repo{
		ID:          e.ID,
		Owner:       e.Owner,
		Name:        e.Name,
		FullName:    e.FullName,
		Description: e.Description,
		Stars:       e.Stars,
		Language:    e.Language,
		HTMLURL:     e.HTMLURL,
		Notes:       e.Notes,
		ForksCount:  e.ForksCount,
		FetchedAt:   e.FetchedAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
