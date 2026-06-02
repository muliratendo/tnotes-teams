package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"

	"github.com/tendo-mulira/tnotes-teams/internal/config"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
)

// Connect establishes connections to PostgreSQL for both Ent and raw SQL (sqlc).
// Returns the raw *sql.DB (for sqlc queries) and the Ent client.
func Connect(cfg *config.Config) (*sql.DB, *ent.Client, error) {
	// Open raw database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")

	// Create Ent client using the same connection pool
	drv := entsql.OpenDB(dialect.Postgres, db)
	entClient := ent.NewClient(ent.Driver(drv))

	return db, entClient, nil
}

// Migrate runs Ent auto-migrations (development only).
func Migrate(ctx context.Context, client *ent.Client) error {
	return client.Schema.Create(ctx)
}
