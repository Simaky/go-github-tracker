// Package storage owns the database handle and the mapping between rows and Go
// structs. Business rules stay in the app layer.
//
// The connection is opened through Ent's SQL dialect driver. Once you add your
// first schema under ent/schema and run `make ent-generate`, wrap this driver
// in the generated *ent.Client to get typed queries plus auto-migration:
//
//	import "github.com/Simaky/go-github-tracker/backend/ent"
//
//	client := ent.NewClient(ent.Driver(s.Driver()))
//	if err := client.Schema.Create(ctx); err != nil { // auto-migrate at startup
//		return nil, fmt.Errorf("auto-migrate: %w", err)
//	}
//
// Until the first schema exists there are no tables to migrate; the connection
// is still live and the service reports healthy.
package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver
)

const reconnectPause = 2 * time.Second

// Storage wraps the database connection.
type Storage struct {
	drv *entsql.Driver
}

// New opens the database, retrying with a fixed back-off until it is reachable
// so the service survives infrastructure that is still coming up.
func New(ctx context.Context, dsn string) (*Storage, error) {
	drv := connect(dsn)

	if err := drv.DB().PingContext(ctx); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	return &Storage{drv: drv}, nil
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

// Driver exposes the underlying Ent SQL driver, for constructing a generated
// *ent.Client once schemas exist.
func (s *Storage) Driver() *entsql.Driver { return s.drv }

// Ping verifies the database is reachable.
func (s *Storage) Ping(ctx context.Context) error {
	if err := s.drv.DB().PingContext(ctx); err != nil {
		return fmt.Errorf("pinging database: %w", err)
	}
	return nil
}

// Close releases the database connection.
func (s *Storage) Close() error {
	if err := s.drv.Close(); err != nil {
		return fmt.Errorf("closing database: %w", err)
	}
	return nil
}
