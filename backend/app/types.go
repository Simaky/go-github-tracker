package app

import (
	"strings"
	"time"
)

// Repo is the API/domain representation of a tracked GitHub repository. It is
// what the app layer returns to the server and, in turn, what clients see.
type Repo struct {
	ID          int       `json:"id"`
	Owner       string    `json:"owner"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Stars       int       `json:"stars"`
	Language    string    `json:"language"`
	HTMLURL     string    `json:"html_url"`
	Notes       string    `json:"notes"`
	FetchedAt   time.Time `json:"fetched_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RepoMetadata is the subset of a repository's data sourced from GitHub. The
// GitHubClient produces it; the Storager persists it on create and refresh.
type RepoMetadata struct {
	Owner       string
	Name        string
	FullName    string
	Description string
	Stars       int
	Language    string
	HTMLURL     string
}

// CreateRepoRequest is the body of POST /api/repos.
type CreateRepoRequest struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

// Validate checks that owner and name are present and are single path segments
// (they are interpolated into the GitHub API path owner/name).
func (r CreateRepoRequest) Validate() error {
	if strings.TrimSpace(r.Owner) == "" {
		return NewValidationError("owner is required")
	}
	if strings.TrimSpace(r.Name) == "" {
		return NewValidationError("name is required")
	}
	if strings.ContainsAny(r.Owner, "/ ") || strings.ContainsAny(r.Name, "/ ") {
		return NewValidationError("owner and name must not contain spaces or slashes")
	}
	return nil
}

// maxNotesLen caps the user-editable notes field to keep payloads sane.
const maxNotesLen = 10_000

// UpdateNotesRequest is the body of PATCH /api/repos/:id.
type UpdateNotesRequest struct {
	Notes string `json:"notes"`
}

// Validate checks the notes length.
func (r UpdateNotesRequest) Validate() error {
	if len(r.Notes) > maxNotesLen {
		return NewValidationError("notes must be at most 10000 characters")
	}
	return nil
}
