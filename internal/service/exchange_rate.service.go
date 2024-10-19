package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nbonair/currency-exchange-server/internal/dataaccess/cache"
	"github.com/nbonair/currency-exchange-server/internal/lib/openexchangerates"
	"github.com/nbonair/currency-exchange-server/internal/repo"
)

type ExchangeRateService interface {
	FetchLatestRates(ctx context.Context, baseCurrency string, targetCurrencies []string) (map[string]float64, error)
}

type exchangeRateService struct {
	repoExchangeRate  repo.ExchangeRateRepository
	repoHistory       repo.ExchangeRateHistoryRepository
	apiClient         openexchangerates.OpenExchangeRateClient
	exchangeRateCache cache.ExchangeRateCache
	cacheTTL          time.Duration
	staleThreshold    time.Duration
}

func NewExchangeRateService(repoExchangeRate repo.ExchangeRateRepository, repoHistory repo.ExchangeRateHistoryRepository, apiClient openexchangerates.OpenExchangeRateClient, exchangeRateCache cache.ExchangeRateCache) ExchangeRateService {
	return &exchangeRateService{
		repoExchangeRate:  repoExchangeRate,
		repoHistory:       repoHistory,
		apiClient:         apiClient,
		exchangeRateCache: exchangeRateCache,
		cacheTTL:          time.Hour * 2,
		staleThreshold:    time.Hour * 4,
	}
}

func (es *exchangeRateService) FetchLatestRates(ctx context.Context, baseCurrency string, targetCurrencies []string) (map[string]float64, error) {
	cachedRates, missingCurrencies, err := es.getRatesFromCache(ctx, baseCurrency, targetCurrencies)
	if err != nil {
		return nil, err
	}

	if len(missingCurrencies) == 0 {
		return cachedRates, nil
	}

	// Get rate from database
	dbRates, staleCurrencies, err := es.getRatesFromDb(ctx, baseCurrency, missingCurrencies)
	if err != nil {
		return nil, err
	}

	// Update cache with fresh rates from DB
	for targetCurrency, rate := range dbRates {
		cachedRates[targetCurrency] = rate.ExchangeRate
		if err := es.exchangeRateCache.Set(ctx, baseCurrency, targetCurrency, rate.ExchangeRate, es.cacheTTL); err != nil {
			fmt.Printf("failed to set cache for %s-%s\n", baseCurrency, targetCurrency)
		}
	}

	if len(staleCurrencies) > 0 {
		apiRates, err := es.fetchRatesFromAPI(baseCurrency, staleCurrencies)
		if err != nil {
			return nil, err
		}

		es.repoExchangeRate.UpdateExchangeRates(ctx, baseCurrency, apiRates)

		for targetCurrency, rate := range apiRates {
			cachedRates[targetCurrency] = rate
			if err := es.exchangeRateCache.Set(ctx, baseCurrency, targetCurrency, rate, es.cacheTTL); err != nil {
				fmt.Printf("failed to set cache for %s-%s\n", baseCurrency, targetCurrency)
			}
			if err := es.repoHistory.InsertRateHistory(ctx, baseCurrency, targetCurrency, rate); err != nil {
				fmt.Printf("failed to insert history rate for %s-%s\n", baseCurrency, targetCurrency)
			}
		}
	}

	return cachedRates, nil
}

func (es *exchangeRateService) getRatesFromCache(ctx context.Context, baseCurrency string, targetCurrencies []string) (map[string]float64, []string, error) {
	rates := make(map[string]float64)
	missingCurrencies := []string{}

	for _, targetCurrency := range targetCurrencies {
		rate, err := es.exchangeRateCache.Get(ctx, baseCurrency, targetCurrency)
		if err == nil {
			rates[targetCurrency] = rate
		} else if errors.Is(err, cache.ErrCacheMiss) {
			missingCurrencies = append(missingCurrencies, targetCurrency)
		} else {
			fmt.Printf("Cache error for %s-%s: %v\n", baseCurrency, targetCurrency, err)
			missingCurrencies = append(missingCurrencies, targetCurrency)
		}
	}
	return rates, missingCurrencies, nil
}

func (es *exchangeRateService) getRatesFromDb(ctx context.Context, baseCurrency string, targetCurrencies []string) (map[string]repo.ExchangeRateDTO, []string, error) {
	dbRates := make(map[string]repo.ExchangeRateDTO)
	staleCurrencies := []string{}
	for _, targetCurrency := range targetCurrencies {
		rate, err := es.repoExchangeRate.GetExchangeRate(ctx, baseCurrency, targetCurrency)
		if err != nil {
			fmt.Printf("database error for %s-%s: %v\n", baseCurrency, targetCurrency, err)
			staleCurrencies = append(staleCurrencies, targetCurrency)
			continue
		}

		if time.Since(rate.UpdatedAt) > es.staleThreshold {
			staleCurrencies = append(staleCurrencies, targetCurrency)
		} else {
			dbRates[targetCurrency] = *rate
		}
	}
	return dbRates, staleCurrencies, nil
}

func (es *exchangeRateService) fetchRatesFromAPI(baseCurrency string, targetCurrencies []string) (map[string]float64, error) {
	rates, err := es.apiClient.FetchLatestRate(baseCurrency, targetCurrencies)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rate from API %w", err)
	}
	return rates, nil
}
