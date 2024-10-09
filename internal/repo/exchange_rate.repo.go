package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nbonair/currency-exchange-server/internal/database"
)

type ExchangeRateRepository interface {
	UpdateExchangeRates(ctx context.Context, baseCurrency string, rates map[string]float64) error
}

type exchangeRateRepository struct {
	pool    *pgxpool.Pool
	queries *database.Queries
}

func NewExchangeRateRepository(pool *pgxpool.Pool) ExchangeRateRepository {
	return &exchangeRateRepository{
		pool:    pool,
		queries: database.New(pool),
	}
}

func (er *exchangeRateRepository) UpdateExchangeRates(ctx context.Context, baseCurrency string, rates map[string]float64) error {
	conn, err := er.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := er.queries.WithTx(tx)

	baseCurrencyId, err := er.queries.GetCurrencyIdByCode(ctx, baseCurrency)
	if err != nil {
		return fmt.Errorf("failed to get base currency ID: %w", err)
	}

	var upsertParams []database.UpsertExchangeRateParams
	for targetCurrency, rate := range rates {
		targetCurrencyId, err := er.queries.GetCurrencyIdByCode(ctx, targetCurrency)
		if err != nil {
			return fmt.Errorf("failed to get currency ID: %w", err)
		}

		upsertParam := database.UpsertExchangeRateParams{
			BaseCurrencyID: baseCurrencyId,
			CurrencyID:     targetCurrencyId,
			ExchangeRate:   rate,
			UpdatedAt:      time.Now(),
		}
		upsertParams = append(upsertParams, upsertParam)
	}

	// Batch Upsert
	for _, params := range upsertParams {
		if err := qtx.UpsertExchangeRate(ctx, params); err != nil {
			return fmt.Errorf("failed to insert into exchange_rate_history: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
