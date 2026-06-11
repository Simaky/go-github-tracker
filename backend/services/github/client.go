// Package github is the outbound client for the GitHub REST API. It owns its
// own types (see types.go) and knows nothing about the app layer; the app
// imports it, calls FetchRepo, and maps the result/errors into its domain.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	defaultBaseURL = "https://api.github.com"
	apiVersion     = "2026-03-10"
)

// Client talks to the GitHub REST API over an *http.Client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// Option customises a Client.
type Option func(*Client)

// WithBaseURL overrides the API base URL (used in tests against a mock server).
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = u }
}

// New constructs a Client. An empty token means unauthenticated requests (lower
// rate limit). A nil httpClient falls back to http.DefaultClient.
func New(httpClient *http.Client, token string, opts ...Option) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{httpClient: httpClient, baseURL: defaultBaseURL, token: token}
	for _, o := range opts {
		o(c)
	}
	return c
}

// repoResponse is the subset of GitHub's repository payload we read.
type repoResponse struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Language    string `json:"language"`
	HTMLURL     string `json:"html_url"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}

// FetchRepo returns metadata for owner/name from GET /repos/{owner}/{name}.
// It reports ErrNotFound on 404 and ErrRateLimited on 403/429; other failures
// are wrapped plain errors.
func (c *Client) FetchRepo(ctx context.Context, owner, name string) (*Repo, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build github request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
		// proceed
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusForbidden, http.StatusTooManyRequests:
		return nil, fmt.Errorf("%w (status %d)", ErrRateLimited, resp.StatusCode)
	default:
		return nil, fmt.Errorf("github: unexpected status %d", resp.StatusCode)
	}

	var body repoResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decoding github response: %w", err)
	}

	return &Repo{
		Owner:       body.Owner.Login,
		Name:        body.Name,
		FullName:    body.FullName,
		Description: body.Description,
		Stars:       body.Stars,
		Language:    body.Language,
		HTMLURL:     body.HTMLURL,
	}, nil
}
