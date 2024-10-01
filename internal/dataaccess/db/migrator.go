package db

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// MigrateUp applies all pending migrations.
func (db *Database) MigrateUp() error {
	goose.SetBaseFS(embedMigrations)
	goose.SetDialect("postgres")

	if err := goose.Up(db.SQLDB, "migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// MigrateDown rolls back the latest migration.
func (db *Database) MigrateDown() error {
	goose.SetBaseFS(embedMigrations)
	goose.SetDialect("postgres")

	if err := goose.Down(db.SQLDB, "migrations"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}
