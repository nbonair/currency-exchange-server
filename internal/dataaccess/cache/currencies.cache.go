package cache

import (
	"context"
	"fmt"
)

const (
	setKeyNameSupportedCurrencies = "supported_currencies_set"
)

type CurrenciesCache interface {
	Has(ctx context.Context, currency string) (bool, error)
	Add(ctx context.Context, currency string) error
	HasMultiple(ctx context.Context, currencies []string) (map[string]bool, error)
}

type currenciesCache struct {
	client Client
}

func NewCurrenciesCache(client Client) CurrenciesCache {
	return &currenciesCache{
		client: client,
	}
}

func (c *currenciesCache) Add(ctx context.Context, currency string) error {
	if err := c.client.AddToSet(ctx, setKeyNameSupportedCurrencies, currency); err != nil {
		return fmt.Errorf("failed to add currency to supported set in cache. Error: %w", err)
	}

	return nil
}

func (c *currenciesCache) Has(ctx context.Context, currency string) (bool, error) {
	result, err := c.client.IsDataInSet(ctx, setKeyNameSupportedCurrencies, currency)
	if err != nil {
		return false, fmt.Errorf("failed to check supported currency. Error: %w", err)
	}

	return result, nil
}

func (c *currenciesCache) HasMultiple(ctx context.Context, currencies []string) (map[string]bool, error) {
	data := make([]any, len(currencies))
	for i, currency := range currencies {
		data[i] = currency
	}
	result, err := c.client.IsDataInSetMultiple(ctx, setKeyNameSupportedCurrencies, data)
	if err != nil {
		return nil, fmt.Errorf("failed to check supported currencies. Error: %w", err)
	}

	supportedMap := make(map[string]bool, len(currencies))
	for i, currency := range currencies {
		supportedMap[currency] = result[i]
	}

	return supportedMap, nil
}
