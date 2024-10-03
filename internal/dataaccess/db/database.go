package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nbonair/currency-exchange-server/configs"
)

type Database struct {
	Pool *pgxpool.Pool
}

func InitializeDB(cfg configs.DatabaseConfig) (*Database, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Printf("error connecting to the database: %+v\n", err)
		return nil, nil, err
	}

	db := &Database{
		Pool: pool,
	}

	if err := db.MigrateUp(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	cleanup := func() {
		db.Close()
	}

	log.Println("Database connection established and migrations applied.")
	return db, cleanup, nil
}

// Close closes the database connections.
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
	log.Println("Database connections closed.")
}
