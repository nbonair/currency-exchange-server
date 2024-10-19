package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/nbonair/currency-exchange-server/internal/dataaccess/db"
	"github.com/nbonair/currency-exchange-server/internal/database"
)

type ExchangeRateHistoryRepository interface {
	InsertRateHistory(ctx context.Context, baseCurrency string, targetCurrency string, exchangeRate float64) error
}

type exchangeRateHistoryRepository struct {
	db      *db.Database
	queries *database.Queries
}

func NewExchangeRateHistoryRepository(db *db.Database) ExchangeRateHistoryRepository {
	return &exchangeRateHistoryRepository{
		db:      db,
		queries: database.New(db.Pool),
	}
}

func (rh *exchangeRateHistoryRepository) InsertRateHistory(ctx context.Context, baseCurrency string, targetCurrency string, exchangeRate float64) error {
	tx, err := rh.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := rh.queries.WithTx(tx)

	insertParam := database.InsertExchangeRateHistoryParams{
		BaseCurrencyCode:   baseCurrency,
		TargetCurrencyCode: targetCurrency,
		ExchangeRate:       exchangeRate,
	}

	qtx.InsertExchangeRateHistory(ctx, insertParam)

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
