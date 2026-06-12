// Package app holds the business/domain layer. It knows nothing about HTTP or
// SQL syntax — the server layer calls into it, and it calls down to interfaces
// it declares itself (see interfaces.go).
package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/Simaky/go-github-tracker/backend/services/github"
)

// App is the domain orchestrator. All dependencies arrive through New as
// interfaces this package declares.
type App struct {
	store  Storager
	github GitHubClient
}

// New constructs the App with its dependencies.
func New(store Storager, github GitHubClient) *App {
	return &App{store: store, github: github}
}

// Health reports whether the service's backing dependencies are reachable.
func (a *App) Health(ctx context.Context) error {
	if err := a.store.Ping(ctx); err != nil {
		return fmt.Errorf("storage health: %w", err)
	}
	return nil
}

// TrackRepo fetches a repository's metadata from GitHub and persists it.
// A duplicate (already-tracked) repo surfaces as a conflict.
func (a *App) TrackRepo(ctx context.Context, req CreateRepoRequest) (*Repo, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	gh, err := a.github.FetchRepo(ctx, req.Owner, req.Name)
	if err != nil {
		return nil, mapGitHubError(err)
	}
	return a.store.CreateRepo(ctx, metadataFromGitHub(gh))
}

// ListRepos returns all tracked repos, optionally filtered by language.
func (a *App) ListRepos(ctx context.Context, language string) ([]*Repo, error) {
	return a.store.ListRepos(ctx, language)
}

// GetRepo returns a single tracked repo by id.
func (a *App) GetRepo(ctx context.Context, id int) (*Repo, error) {
	return a.store.GetRepo(ctx, id)
}

// UpdateNotes updates the user-editable notes on a tracked repo.
func (a *App) UpdateNotes(ctx context.Context, id int, req UpdateNotesRequest) (*Repo, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return a.store.UpdateNotes(ctx, id, req.Notes)
}

// DeleteRepo removes a repo from the watchlist.
func (a *App) DeleteRepo(ctx context.Context, id int) error {
	return a.store.DeleteRepo(ctx, id)
}

// RefreshRepo re-fetches a tracked repo's metadata from GitHub and updates the
// stored fields. Notes are preserved.
func (a *App) RefreshRepo(ctx context.Context, id int) (*Repo, error) {
	existing, err := a.store.GetRepo(ctx, id)
	if err != nil {
		return nil, err
	}
	gh, err := a.github.FetchRepo(ctx, existing.Owner, existing.Name)
	if err != nil {
		return nil, mapGitHubError(err)
	}
	return a.store.RefreshRepo(ctx, id, metadataFromGitHub(gh))
}

// metadataFromGitHub maps the github service's type into the app's storage
// contract.
func metadataFromGitHub(r *github.RepoResponse) RepoMetadata {
	return RepoMetadata{
		Owner:       r.Owner.Login,
		Name:        r.Name,
		FullName:    r.FullName,
		Description: r.Description,
		Stars:       r.Stars,
		Language:    r.Language,
		HTMLURL:     r.HTMLURL,
	}
}

// mapGitHubError translates the github service's sentinel errors into domain
// errors; anything else is treated as an upstream failure.
func mapGitHubError(err error) error {
	switch {
	case errors.Is(err, github.ErrNotFound):
		return NewNotFoundError("repository")
	case errors.Is(err, github.ErrRateLimited):
		return NewUpstreamError("github rate limit exceeded", err)
	default:
		return NewUpstreamError("github request failed", err)
	}
}
