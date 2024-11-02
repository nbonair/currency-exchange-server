package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nbonair/currency-exchange-server/internal/dataaccess/db"
	"github.com/nbonair/currency-exchange-server/internal/database"
)

type ExchangeRateRepository interface {
	GetSupportedCurrencies(ctx context.Context) ([]string, error)
	GetExchangeRate(ctx context.Context, baseCurrency string, targetCurrent string) (*ExchangeRateDTO, error)
	UpdateExchangeRates(ctx context.Context, baseCurrency string, rates map[string]float64) error
}

type exchangeRateRepository struct {
	db      *db.Database
	queries *database.Queries
}

type ExchangeRateDTO struct {
	BaseCurrencyCode   string    `json:"base_currency_code"`
	TargetCurrencyCode string    `json:"target_currency_code"`
	ExchangeRate       float64   `json:"exchange_rate"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func NewExchangeRateRepository(db *db.Database) ExchangeRateRepository {
	return &exchangeRateRepository{
		db:      db,
		queries: database.New(db.Pool),
	}
}

func (er *exchangeRateRepository) GetSupportedCurrencies(ctx context.Context) ([]string, error) {
	codes, err := er.queries.GetSupportedCurrencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get supported currencies list. Error: %w", err)
	}

	var supportedCurrencies []string
	supportedCurrencies = append(supportedCurrencies, codes...)

	return supportedCurrencies, nil
}

func (er *exchangeRateRepository) GetExchangeRate(ctx context.Context, baseCurrency string, targetCurrent string) (*ExchangeRateDTO, error) {
	row, err := er.queries.GetExchangeRate(ctx, database.GetExchangeRateParams{
		BaseCurrencyCode:   baseCurrency,
		TargetCurrencyCode: targetCurrent,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate %w", err)
	}
	exchangeRateDTO := &ExchangeRateDTO{
		BaseCurrencyCode:   row.BaseCurrencyCode,
		TargetCurrencyCode: row.TargetCurrencyCode,
		ExchangeRate:       row.ExchangeRate,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
	return exchangeRateDTO, nil
}

func (er *exchangeRateRepository) UpdateExchangeRates(ctx context.Context, baseCurrency string, rates map[string]float64) error {

	tx, err := er.db.Pool.BeginTx(ctx, pgx.TxOptions{})
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
