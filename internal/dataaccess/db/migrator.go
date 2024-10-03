package db

import (
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

// MigrateUp applies all pending migrations.
func (db *Database) MigrateUp() error {
	sqlDB := stdlib.OpenDB(*db.Pool.Config().ConnConfig)
	defer sqlDB.Close()

	goose.SetBaseFS(embedMigrations)
	goose.SetDialect("postgres")

	if err := goose.Up(sqlDB, "sql/schema"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// MigrateDown rolls back the latest migration.
func (db *Database) MigrateDown() error {
	sqlDB := stdlib.OpenDBFromPool(db.Pool)
	defer sqlDB.Close()

	goose.SetBaseFS(embedMigrations)
	goose.SetDialect("postgres")

	if err := goose.Down(sqlDB, "./sql/schema"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}
