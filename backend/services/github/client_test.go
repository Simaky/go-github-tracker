package github_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Simaky/go-github-tracker/backend/services/github"
)

const repoJSON = `{
	"name": "go",
	"full_name": "golang/go",
	"description": "The Go programming language",
	"stargazers_count": 120000,
	"language": "Go",
	"html_url": "https://github.com/golang/go",
	"owner": {"login": "golang"}
}`

func TestFetchRepo_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/golang/go" {
			t.Errorf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer tok" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer tok")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(repoJSON))
	}))
	defer srv.Close()

	c := github.New(srv.Client(), "tok", github.WithBaseURL(srv.URL))
	got, err := c.FetchRepo(context.Background(), "golang", "go")
	if err != nil {
		t.Fatalf("FetchRepo: unexpected error: %v", err)
	}

	want := &github.Repo{
		Owner:       "golang",
		Name:        "go",
		FullName:    "golang/go",
		Description: "The Go programming language",
		Stars:       120000,
		Language:    "Go",
		HTMLURL:     "https://github.com/golang/go",
	}
	if *got != *want {
		t.Fatalf("repo = %+v, want %+v", *got, *want)
	}
}

func TestFetchRepo_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := github.New(srv.Client(), "", github.WithBaseURL(srv.URL))
	_, err := c.FetchRepo(context.Background(), "missing", "repo")
	if !errors.Is(err, github.ErrNotFound) {
		t.Fatalf("error = %v, want ErrNotFound", err)
	}
}

func TestFetchRepo_RateLimited(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	c := github.New(srv.Client(), "", github.WithBaseURL(srv.URL))
	_, err := c.FetchRepo(context.Background(), "o", "n")
	if !errors.Is(err, github.ErrRateLimited) {
		t.Fatalf("error = %v, want ErrRateLimited", err)
	}
}

func TestFetchRepo_UnauthenticatedOmitsAuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Authorization"]; ok {
			t.Error("Authorization header should be absent without a token")
		}
		_, _ = w.Write([]byte(repoJSON))
	}))
	defer srv.Close()

	c := github.New(srv.Client(), "", github.WithBaseURL(srv.URL))
	if _, err := c.FetchRepo(context.Background(), "golang", "go"); err != nil {
		t.Fatalf("FetchRepo: unexpected error: %v", err)
	}
}
