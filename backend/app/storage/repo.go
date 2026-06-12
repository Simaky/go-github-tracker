package storage

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

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
		SetForksCount(m.ForksCount).
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

// CountRepos returns how many repos are currently tracked.
func (s *Storage) CountRepos(ctx context.Context) (int, error) {
	n, err := s.client.Repo.Query().Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count repos: %w", err)
	}
	return n, nil
}

// TotalStars returns the sum of the stars field across all tracked repos.
// Returns 0 when no repos are tracked.
func (s *Storage) TotalStars(ctx context.Context) (int, error) {
	var agg []struct {
		Sum int `json:"sum"`
	}
	err := s.client.Repo.Query().
		Aggregate(ent.Sum(repo.FieldStars)).
		Scan(ctx, &agg)
	if err != nil {
		return 0, fmt.Errorf("total stars: %w", err)
	}
	// No rows → empty slice; SUM over no rows is NULL, decoded as the zero value.
	if len(agg) == 0 {
		return 0, nil
	}
	return agg[0].Sum, nil
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
		SetForksCount(m.ForksCount).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, app.NewNotFoundError("repository")
		}
		return nil, fmt.Errorf("refresh repo: %w", err)
	}
	return toAppRepo(row), nil
}

// MostUsedLanguage returns the language tracked across the most repos, along
// with how many repos use it. Repos with no language are ignored. When no repo
// has a language set, it returns ("", 0, nil).
//
// The whole computation runs in the database: GROUP BY language, COUNT, then
// ORDER BY that count DESC LIMIT 1 — so a single row comes back, no Go-side
// tallying. The ordering/limit are pushed onto the selector from inside the
// aggregate function, which also emits the COUNT(*) AS count column.
func (s *Storage) MostUsedLanguage(ctx context.Context) (string, int, error) {
	var rows []struct {
		Language string `json:"language"`
		Count    int    `json:"count"`
	}
	err := s.client.Repo.Query().
		Where(repo.LanguageNEQ("")).
		GroupBy(repo.FieldLanguage).
		Aggregate(func(sel *sql.Selector) string {
			sel.OrderBy(sql.Desc("count")).Limit(1)
			return sql.As(sql.Count("*"), "count")
		}).
		Scan(ctx, &rows)
	if err != nil {
		return "", 0, fmt.Errorf("most used language: %w", err)
	}
	if len(rows) == 0 {
		return "", 0, nil
	}
	return rows[0].Language, rows[0].Count, nil
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
