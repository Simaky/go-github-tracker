package app_test

import (
	"context"
	"testing"

	"github.com/Simaky/go-github-tracker/backend/app"
	"github.com/Simaky/go-github-tracker/backend/services/github"
)

// stubStore is an in-memory Storager; only the methods a test needs are set.
type stubStore struct {
	createFn func(context.Context, app.RepoMetadata) (*app.Repo, error)
}

func (s *stubStore) Ping(context.Context) error { return nil }
func (s *stubStore) CreateRepo(ctx context.Context, m app.RepoMetadata) (*app.Repo, error) {
	return s.createFn(ctx, m)
}
func (s *stubStore) ListRepos(context.Context, string) ([]*app.Repo, error)      { return nil, nil }
func (s *stubStore) GetRepo(context.Context, int) (*app.Repo, error)             { return nil, nil }
func (s *stubStore) UpdateNotes(context.Context, int, string) (*app.Repo, error) { return nil, nil }
func (s *stubStore) RefreshRepo(context.Context, int, app.RepoMetadata) (*app.Repo, error) {
	return nil, nil
}
func (s *stubStore) DeleteRepo(context.Context, int) error   { return nil }
func (s *stubStore) CountRepos(context.Context) (int, error) { return 0, nil }
func (s *stubStore) TotalStars(context.Context) (int, error) { return 0, nil }
func (s *stubStore) MostUsedLanguage(context.Context) (string, int, error) {
	return "", 0, nil
}

// stubGitHub is a configurable GitHubClient.
type stubGitHub struct {
	fetchFn func(context.Context, string, string) (*github.RepoResponse, error)
}

func (g *stubGitHub) FetchRepo(ctx context.Context, owner, name string) (*github.RepoResponse, error) {
	return g.fetchFn(ctx, owner, name)
}

// TestTrackRepo_Success exercises the core business flow end to end: validate
// the request, fetch metadata from GitHub, map it into a RepoMetadata, and hand
// that to storage. It asserts both the returned repo and the value the app maps
// for persistence.
func TestTrackRepo_Success(t *testing.T) {
	gh := &github.RepoResponse{
		Owner:    github.Owner{Login: "golang"},
		Name:     "go",
		FullName: "golang/go",
		Stars:    5,
	}

	var gotMeta app.RepoMetadata
	store := &stubStore{
		createFn: func(_ context.Context, m app.RepoMetadata) (*app.Repo, error) {
			gotMeta = m
			return &app.Repo{ID: 1, FullName: m.FullName}, nil
		},
	}
	client := &stubGitHub{fetchFn: func(_ context.Context, _, _ string) (*github.RepoResponse, error) {
		return gh, nil
	}}

	got, err := app.New(store, client).TrackRepo(context.Background(),
		app.CreateRepoRequest{Owner: "golang", Name: "go"})
	if err != nil {
		t.Fatalf("TrackRepo: unexpected error: %v", err)
	}
	if got.ID != 1 || got.FullName != "golang/go" {
		t.Fatalf("repo = %+v, want id 1 golang/go", got)
	}

	want := app.RepoMetadata{Owner: "golang", Name: "go", FullName: "golang/go", Stars: 5}
	if gotMeta != want {
		t.Fatalf("store received %+v, want %+v", gotMeta, want)
	}
}
