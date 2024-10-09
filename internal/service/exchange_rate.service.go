package service

import (
	"context"
	"fmt"

	"github.com/nbonair/currency-exchange-server/internal/lib/openexchangerates"
	"github.com/nbonair/currency-exchange-server/internal/repo"
)

type ExchangeRateService interface {
	FetchLatestRates(ctx context.Context, baseCurrency string, targetCurrencies []string) (map[string]float64, error)
}

type exchangeRateService struct {
	repo      repo.ExchangeRateRepository
	apiClient openexchangerates.OpenExchangeRateClient
}

func NewExchangeRateService(repo repo.ExchangeRateRepository, apiClient openexchangerates.OpenExchangeRateClient) ExchangeRateService {
	return &exchangeRateService{
		repo:      repo,
		apiClient: apiClient,
	}
}

// FetchLatestRates implements ExchangeRateService.
func (es *exchangeRateService) FetchLatestRates(ctx context.Context, baseCurrency string, targetCurrencies []string) (map[string]float64, error) {
	rates, err := es.apiClient.FetchLatestRate(baseCurrency, targetCurrencies)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rate from API %w", err)
	}

	if err := es.repo.UpdateExchangeRates(ctx, baseCurrency, rates); err != nil {
		return nil, fmt.Errorf("failed to update rates into database: %w", err)
	}

	return rates, nil
}
