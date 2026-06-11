// Package services is the home for outbound callers to external systems — one
// sub-package per upstream (e.g. services/github/ wrapping the GitHub API).
// Each sub-package owns its own Client type over *http.Client, its wire types,
// and its sentinel errors. It is empty in the bare skeleton.
package services
