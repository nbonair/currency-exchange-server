package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type ExchangeRateCache interface {
	Get(ctx context.Context, baseCurrency string, targetCurrency string) (float64, error)
	Set(ctx context.Context, baseCurrency string, targetCurrency string, rate float64, ttl time.Duration) error
}

type exchangeRateCache struct {
	client Client
}

func NewExchangeRateCache(client Client) ExchangeRateCache {
	return &exchangeRateCache{
		client: client,
	}
}

func (c *exchangeRateCache) getExchangeRateCacheKey(baseCurrency string, targetCurrency string) string {
	return fmt.Sprintf("exchange_rate:%s:%s", baseCurrency, targetCurrency)
}

// Get implements ExchangeRateCache.
func (c *exchangeRateCache) Get(ctx context.Context, baseCurrency string, targetCurrency string) (float64, error) {
	cacheKey := c.getExchangeRateCacheKey(baseCurrency, targetCurrency)
	cacheEntry, err := c.client.Get(ctx, cacheKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get exchange rate cache: %s:%s. Error: %s", baseCurrency, targetCurrency, err)

	}

	if cacheEntry == "" {
		return 0, ErrCacheMiss
	}
	rate, err := strconv.ParseFloat(cacheEntry, 64)
	if err != nil {
		return 0, fmt.Errorf("cache entry is not of type float")
	}

	return rate, nil
}

// Set implements ExchangeRateCache.
func (c *exchangeRateCache) Set(ctx context.Context, baseCurrency string, targetCurrency string, rate float64, ttl time.Duration) error {
	cacheKey := c.getExchangeRateCacheKey(baseCurrency, targetCurrency)
	if err := c.client.Set(ctx, cacheKey, rate, ttl); err != nil {
		return fmt.Errorf("failed to insert new %s:%s rate into cache", baseCurrency, targetCurrency)
	}
	return nil
}
