package cache

import (
	"context"
	"fmt"
	"time"
)

type ExchangeRateCache interface {
	Get(ctx context.Context, baseCurrency string, targerCurrency string) (float64, error)
	Set(ctx context.Context, baseCurrency string, targerCurrency string, rate float64) error
}

type exchangeRateCache struct {
	client Client
	ttl    time.Duration
}

func NewExchangeRateCache(client Client, ttl time.Duration) ExchangeRateCache {
	return &exchangeRateCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *exchangeRateCache) getExchangeRateCacheKey(baseCurrency string, targetCurrency string) string {
	return fmt.Sprintf("exchange_rate:%s:%s", baseCurrency, targetCurrency)
}

// Get implements ExchangeRateCache.
func (c *exchangeRateCache) Get(ctx context.Context, baseCurrency string, targerCurrency string) (float64, error) {
	cacheKey := c.getExchangeRateCacheKey(baseCurrency, targerCurrency)
	cacheEntry, err := c.client.Get(ctx, cacheKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get exchange rate cache: %s:%s. Error: %s", baseCurrency, targerCurrency, err)

	}

	if cacheEntry == nil {
		return 0, ErrCacheMiss
	}

	rate, ok := cacheEntry.(float64)
	if !ok {
		return 0, fmt.Errorf("cache entry is not of type float")
	}

	return rate, nil
}

// Set implements ExchangeRateCache.
func (c *exchangeRateCache) Set(ctx context.Context, baseCurrency string, targerCurrency string, rate float64) error {
	cacheKey := c.getExchangeRateCacheKey(baseCurrency, targerCurrency)
	if err := c.client.Set(ctx, cacheKey, rate, c.ttl); err != nil {
		return fmt.Errorf("failed to insert new %s:%s rate into cache", baseCurrency, targerCurrency)
	}
	return nil
}
