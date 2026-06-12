package app

import (
	"context"

	"github.com/Simaky/go-github-tracker/backend/services/github"
)

// Storager is what the app needs from durable storage. The concrete
// *storage.Storage (constructed in main) satisfies it. Methods return the
// app-layer Repo and translate storage-specific failures into domain errors
// (not-found, conflict) so the app layer stays persistence-agnostic.
type Storager interface {
	Ping(ctx context.Context) error
	CreateRepo(ctx context.Context, m RepoMetadata) (*Repo, error)
	ListRepos(ctx context.Context, language string) ([]*Repo, error)
	GetRepo(ctx context.Context, id int) (*Repo, error)
	UpdateNotes(ctx context.Context, id int, notes string) (*Repo, error)
	RefreshRepo(ctx context.Context, id int, m RepoMetadata) (*Repo, error)
	DeleteRepo(ctx context.Context, id int) error
	CountRepos(ctx context.Context) (int, error)
	TotalStars(ctx context.Context) (int, error)
	MostUsedLanguage(ctx context.Context) (string, int, error)
}

// GitHubClient is what the app needs from the GitHub API. The concrete client
// in services/github satisfies it. It returns the service's own *github.Repo;
// the app maps that (and the service's sentinel errors) into its domain.
type GitHubClient interface {
	FetchRepo(ctx context.Context, owner, name string) (*github.RepoResponse, error)
}
