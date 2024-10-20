package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nbonair/currency-exchange-server/internal/dataaccess/db"
	"github.com/nbonair/currency-exchange-server/internal/database"
)

type ExchangeRateHistoryDTO struct {
	BaseCurrencyCode   string    `json:"base_currency_code"`
	TargetCurrencyCode string    `json:"target_currency_code"`
	Rate               float64   `json:"exchange_rate"`
	Timestamp          time.Time `json:"timestamp"`
}

type ExchangeRateHistoryRepository interface {
	InsertRateHistory(ctx context.Context, baseCurrency string, targetCurrency string, exchangeRate float64) error
	GetRateHistory(ctx context.Context, baseCurrency string, targetCurrency string, startTime time.Time, endTime time.Time) ([]ExchangeRateHistoryDTO, error)
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

func (rh *exchangeRateHistoryRepository) GetRateHistory(ctx context.Context, baseCurrency string, targetCurrency string, startTime time.Time, endTime time.Time) ([]ExchangeRateHistoryDTO, error) {

	rows, err := rh.queries.GetExchangeRateHistory(ctx, database.GetExchangeRateHistoryParams{
		BaseCurrencyCode:   baseCurrency,
		TargetCurrencyCode: targetCurrency,
		StartTime:          startTime,
		EndTime:            endTime,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate %w", err)
	}
	var rateHistory []ExchangeRateHistoryDTO

	for _, row := range rows {
		rateHistory = append(rateHistory, ExchangeRateHistoryDTO{
			BaseCurrencyCode:   row.BaseCurrencyCode,
			TargetCurrencyCode: row.TargetCurrencyCode,
			Rate:               row.Rate,
			Timestamp:          row.Timestamp,
		})
	}
	return rateHistory, nil
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
