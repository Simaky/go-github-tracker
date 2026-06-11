package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Simaky/go-github-tracker/backend/app"
	"github.com/Simaky/go-github-tracker/backend/ent"
	"github.com/Simaky/go-github-tracker/backend/ent/repo"
)

// CreateRepo persists newly-fetched repository metadata. A duplicate full_name
// (the repo is already tracked) surfaces as a domain conflict.
func (s *Storage) CreateRepo(ctx context.Context, m app.RepoMetadata) (*app.Repo, error) {
	row, err := s.client.Repo.Create().
		SetOwner(m.Owner).
		SetName(m.Name).
		SetFullName(m.FullName).
		SetDescription(m.Description).
		SetStars(m.Stars).
		SetLanguage(m.Language).
		SetHTMLURL(m.HTMLURL).
		SetFetchedAt(time.Now()).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, app.NewConflictError(m.FullName + " is already tracked")
		}
		return nil, fmt.Errorf("create repo: %w", err)
	}
	return toAppRepo(row), nil
}

// ListRepos returns tracked repos ordered by id, optionally filtered by language.
func (s *Storage) ListRepos(ctx context.Context, language string) ([]*app.Repo, error) {
	q := s.client.Repo.Query()
	if language != "" {
		q = q.Where(repo.Language(language))
	}
	rows, err := q.Order(ent.Asc(repo.FieldID)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list repos: %w", err)
	}
	out := make([]*app.Repo, len(rows))
	for i, row := range rows {
		out[i] = toAppRepo(row)
	}
	return out, nil
}

// GetRepo returns a single tracked repo by id.
func (s *Storage) GetRepo(ctx context.Context, id int) (*app.Repo, error) {
	row, err := s.client.Repo.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, app.NewNotFoundError("repository")
		}
		return nil, fmt.Errorf("get repo: %w", err)
	}
	return toAppRepo(row), nil
}

// UpdateNotes sets the user-editable notes on a tracked repo.
func (s *Storage) UpdateNotes(ctx context.Context, id int, notes string) (*app.Repo, error) {
	row, err := s.client.Repo.UpdateOneID(id).SetNotes(notes).Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, app.NewNotFoundError("repository")
		}
		return nil, fmt.Errorf("update notes: %w", err)
	}
	return toAppRepo(row), nil
}

// RefreshRepo overwrites the GitHub-sourced fields with freshly-fetched metadata
// and bumps fetched_at. Notes (user-owned) are left untouched.
func (s *Storage) RefreshRepo(ctx context.Context, id int, m app.RepoMetadata) (*app.Repo, error) {
	row, err := s.client.Repo.UpdateOneID(id).
		SetOwner(m.Owner).
		SetName(m.Name).
		SetFullName(m.FullName).
		SetDescription(m.Description).
		SetStars(m.Stars).
		SetLanguage(m.Language).
		SetHTMLURL(m.HTMLURL).
		SetFetchedAt(time.Now()).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, app.NewNotFoundError("repository")
		}
		return nil, fmt.Errorf("refresh repo: %w", err)
	}
	return toAppRepo(row), nil
}

// DeleteRepo removes a tracked repo by id.
func (s *Storage) DeleteRepo(ctx context.Context, id int) error {
	if err := s.client.Repo.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return app.NewNotFoundError("repository")
		}
		return fmt.Errorf("delete repo: %w", err)
	}
	return nil
}
