// Package app holds the business/domain layer. It knows nothing about HTTP or
// SQL syntax — the server layer calls into it, and it calls down to interfaces
// it declares itself (see interfaces.go).
package app

import (
	"context"
	"fmt"
)

// App is the domain orchestrator. All dependencies arrive through New as
// interfaces this package declares.
type App struct {
	store Storager
}

// New constructs the App with its dependencies.
func New(store Storager) *App {
	return &App{store: store}
}

// Health reports whether the service's backing dependencies are reachable.
func (a *App) Health(ctx context.Context) error {
	if err := a.store.Ping(ctx); err != nil {
		return fmt.Errorf("storage health: %w", err)
	}
	return nil
}
