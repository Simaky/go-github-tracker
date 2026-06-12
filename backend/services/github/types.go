package github

import "errors"

type (
	// RepoResponse is the subset of GitHub's repository payload we read.
	RepoResponse struct {
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		Stars       int    `json:"stargazers_count"`
		Language    string `json:"language"`
		HTMLURL     string `json:"html_url"`
		Owner       Owner  `json:"owner"`
	}

	Owner struct {
		Login string `json:"login"`
	}
)

// Sentinel errors callers match with errors.Is; the app translates them into
// domain errors.
var (
	// ErrNotFound means GitHub has no such repository (404).
	ErrNotFound = errors.New("github: repository not found")
	// ErrRateLimited means the API rate limit was hit (403/429).
	ErrRateLimited = errors.New("github: rate limit exceeded")
)
