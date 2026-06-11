package app_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Simaky/go-github-tracker/backend/app"
	"github.com/Simaky/go-github-tracker/backend/services/github"
)

// stubStore is a configurable in-memory Storager for app-level tests.
type stubStore struct {
	pingErr   error
	createFn  func(context.Context, app.RepoMetadata) (*app.Repo, error)
	listFn    func(context.Context, string) ([]*app.Repo, error)
	getFn     func(context.Context, int) (*app.Repo, error)
	updateFn  func(context.Context, int, string) (*app.Repo, error)
	refreshFn func(context.Context, int, app.RepoMetadata) (*app.Repo, error)
	deleteFn  func(context.Context, int) error
}

func (s *stubStore) Ping(context.Context) error { return s.pingErr }
func (s *stubStore) CreateRepo(ctx context.Context, m app.RepoMetadata) (*app.Repo, error) {
	return s.createFn(ctx, m)
}
func (s *stubStore) ListRepos(ctx context.Context, lang string) ([]*app.Repo, error) {
	return s.listFn(ctx, lang)
}
func (s *stubStore) GetRepo(ctx context.Context, id int) (*app.Repo, error) {
	return s.getFn(ctx, id)
}
func (s *stubStore) UpdateNotes(ctx context.Context, id int, notes string) (*app.Repo, error) {
	return s.updateFn(ctx, id, notes)
}
func (s *stubStore) RefreshRepo(ctx context.Context, id int, m app.RepoMetadata) (*app.Repo, error) {
	return s.refreshFn(ctx, id, m)
}
func (s *stubStore) DeleteRepo(ctx context.Context, id int) error { return s.deleteFn(ctx, id) }

// stubGitHub is a configurable GitHubClient.
type stubGitHub struct {
	fetchFn func(context.Context, string, string) (*github.Repo, error)
}

func (g *stubGitHub) FetchRepo(ctx context.Context, owner, name string) (*github.Repo, error) {
	return g.fetchFn(ctx, owner, name)
}

func domainCode(t *testing.T, err error) string {
	t.Helper()
	var de *app.DomainError
	if !errors.As(err, &de) {
		t.Fatalf("error %v is not a *app.DomainError", err)
	}
	return de.Code
}

func TestTrackRepo_Success(t *testing.T) {
	gh := &github.Repo{Owner: "golang", Name: "go", FullName: "golang/go", Stars: 5}
	var gotMeta app.RepoMetadata
	store := &stubStore{
		createFn: func(_ context.Context, m app.RepoMetadata) (*app.Repo, error) {
			gotMeta = m
			return &app.Repo{ID: 1, FullName: m.FullName}, nil
		},
	}
	client := &stubGitHub{fetchFn: func(_ context.Context, _, _ string) (*github.Repo, error) {
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

func TestTrackRepo_ValidationSkipsGitHub(t *testing.T) {
	client := &stubGitHub{fetchFn: func(_ context.Context, _, _ string) (*github.Repo, error) {
		t.Fatal("GitHub should not be called when validation fails")
		return nil, nil
	}}
	store := &stubStore{} // no funcs set; must not be called

	_, err := app.New(store, client).TrackRepo(context.Background(),
		app.CreateRepoRequest{Owner: "", Name: "go"})
	if code := domainCode(t, err); code != app.CodeValidation {
		t.Fatalf("code = %s, want %s", code, app.CodeValidation)
	}
}

func TestTrackRepo_GitHubNotFoundMapped(t *testing.T) {
	client := &stubGitHub{fetchFn: func(_ context.Context, _, _ string) (*github.Repo, error) {
		return nil, github.ErrNotFound
	}}
	store := &stubStore{
		createFn: func(context.Context, app.RepoMetadata) (*app.Repo, error) {
			t.Fatal("storage should not be called when GitHub fails")
			return nil, nil
		},
	}

	_, err := app.New(store, client).TrackRepo(context.Background(),
		app.CreateRepoRequest{Owner: "missing", Name: "repo"})
	if code := domainCode(t, err); code != app.CodeNotFound {
		t.Fatalf("code = %s, want %s", code, app.CodeNotFound)
	}
}

func TestTrackRepo_RateLimitMappedToUpstream(t *testing.T) {
	client := &stubGitHub{fetchFn: func(_ context.Context, _, _ string) (*github.Repo, error) {
		return nil, github.ErrRateLimited
	}}
	store := &stubStore{}

	_, err := app.New(store, client).TrackRepo(context.Background(),
		app.CreateRepoRequest{Owner: "golang", Name: "go"})
	if code := domainCode(t, err); code != app.CodeUpstream {
		t.Fatalf("code = %s, want %s", code, app.CodeUpstream)
	}
}

func TestTrackRepo_DuplicateConflict(t *testing.T) {
	client := &stubGitHub{fetchFn: func(_ context.Context, _, _ string) (*github.Repo, error) {
		return &github.Repo{FullName: "golang/go"}, nil
	}}
	store := &stubStore{
		createFn: func(context.Context, app.RepoMetadata) (*app.Repo, error) {
			return nil, app.NewConflictError("golang/go is already tracked")
		},
	}

	_, err := app.New(store, client).TrackRepo(context.Background(),
		app.CreateRepoRequest{Owner: "golang", Name: "go"})
	if code := domainCode(t, err); code != app.CodeConflict {
		t.Fatalf("code = %s, want %s", code, app.CodeConflict)
	}
}

func TestRefreshRepo_FetchesStoredOwnerName(t *testing.T) {
	const id = 7
	var fetchedOwner, fetchedName string
	var refreshedID int
	store := &stubStore{
		getFn: func(_ context.Context, gotID int) (*app.Repo, error) {
			if gotID != id {
				t.Errorf("GetRepo id = %d, want %d", gotID, id)
			}
			return &app.Repo{ID: id, Owner: "golang", Name: "go"}, nil
		},
		refreshFn: func(_ context.Context, gotID int, _ app.RepoMetadata) (*app.Repo, error) {
			refreshedID = gotID
			return &app.Repo{ID: id, Stars: 99}, nil
		},
	}
	client := &stubGitHub{fetchFn: func(_ context.Context, owner, name string) (*github.Repo, error) {
		fetchedOwner, fetchedName = owner, name
		return &github.Repo{FullName: "golang/go", Stars: 99}, nil
	}}

	got, err := app.New(store, client).RefreshRepo(context.Background(), id)
	if err != nil {
		t.Fatalf("RefreshRepo: unexpected error: %v", err)
	}
	if fetchedOwner != "golang" || fetchedName != "go" {
		t.Fatalf("fetched %s/%s, want golang/go", fetchedOwner, fetchedName)
	}
	if refreshedID != id || got.Stars != 99 {
		t.Fatalf("refresh id=%d stars=%d, want id=%d stars=99", refreshedID, got.Stars, id)
	}
}

func TestUpdateNotes_TooLong(t *testing.T) {
	store := &stubStore{
		updateFn: func(context.Context, int, string) (*app.Repo, error) {
			t.Fatal("storage should not be called when validation fails")
			return nil, nil
		},
	}
	a := app.New(store, &stubGitHub{})

	_, err := a.UpdateNotes(context.Background(), 1,
		app.UpdateNotesRequest{Notes: strings.Repeat("x", 10_001)})
	if code := domainCode(t, err); code != app.CodeValidation {
		t.Fatalf("code = %s, want %s", code, app.CodeValidation)
	}
}
