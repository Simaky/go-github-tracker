// Package storage owns the database handle and the mapping between rows and Go
// structs. Business rules stay in the app layer.
//
// The connection is opened through Ent's SQL dialect driver and wrapped in the
// generated *ent.Client, which gives typed queries plus auto-migration. The
// schema is created/upgraded on startup via client.Schema.Create.
package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver

	"github.com/Simaky/go-github-tracker/backend/ent"
)

const reconnectPause = 2 * time.Second

// Storage wraps the Ent client and the underlying SQL driver.
type Storage struct {
	client *ent.Client
	drv    *entsql.Driver
}

// New opens the database (retrying until reachable so the service survives
// infrastructure that is still coming up), wraps it in the Ent client, and runs
// auto-migration to create/upgrade the schema.
func New(ctx context.Context, dsn string) (*Storage, error) {
	drv := connect(dsn)

	if err := drv.DB().PingContext(ctx); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	client := ent.NewClient(ent.Driver(drv))
	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("auto-migrate: %w", err)
	}
	return &Storage{client: client, drv: drv}, nil
}

func connect(dsn string) *entsql.Driver {
	for {
		drv, err := entsql.Open(dialect.Postgres, dsn)
		if err != nil {
			log.Printf("db open: %s", err)
			time.Sleep(reconnectPause)
			continue
		}
		if err := drv.DB().Ping(); err != nil {
			log.Printf("db ping: %s", err)
			_ = drv.Close()
			time.Sleep(reconnectPause)
			continue
		}
		return drv
	}
}

// Ping verifies the database is reachable.
func (s *Storage) Ping(ctx context.Context) error {
	if err := s.drv.DB().PingContext(ctx); err != nil {
		return fmt.Errorf("pinging database: %w", err)
	}
	return nil
}

// Close releases the database connection.
func (s *Storage) Close() error {
	if err := s.client.Close(); err != nil {
		return fmt.Errorf("closing database: %w", err)
	}
	return nil
}
