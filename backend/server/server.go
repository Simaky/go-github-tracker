// Package server receives traffic from the outside world (HTTP), decodes it,
// calls a method on the *app.App, and serialises the result. It holds no
// business logic.
package server

import (
	"github.com/Simaky/go-github-tracker/backend/app"
	"github.com/Simaky/go-github-tracker/backend/app/config"
)

type server struct {
	app *app.App
	cfg config.Config
}

// New constructs the server with the application and loaded config.
func New(appInst *app.App, cfg config.Config) *server {
	return &server{app: appInst, cfg: cfg}
}

// Run starts the HTTP transport and blocks until it stops.
func (s *server) Run(version string) error {
	return s.runHTTP(version)
}

// addressOf returns the bind address; config always supplies a default.
func (s *server) addressOf() string {
	return s.cfg.Listen
}
