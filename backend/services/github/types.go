package github

import "errors"

// Repo is the repository metadata returned by the GitHub API. It is the github
// service's own type (a view of a 3rd-party payload), not an app type — the app
// imports it and maps it into its domain.
type Repo struct {
	Owner       string
	Name        string
	FullName    string
	Description string
	Stars       int
	Language    string
	HTMLURL     string
}

// Sentinel errors callers match with errors.Is; the app translates them into
// domain errors.
var (
	// ErrNotFound means GitHub has no such repository (404).
	ErrNotFound = errors.New("github: repository not found")
	// ErrRateLimited means the API rate limit was hit (403/429).
	ErrRateLimited = errors.New("github: rate limit exceeded")
)
