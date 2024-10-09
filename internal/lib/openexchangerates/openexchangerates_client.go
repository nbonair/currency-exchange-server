package openexchangerates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OpenExchangeRateClient interface {
	FetchLatestRate(baseCurrency string, targetCurrencies []string) (map[string]float64, error)
}

type openExchangeRateClient struct {
	appID   string
	baseURL string
	client  *http.Client
}

func NewOpenExchangeRateClient(appID string) (OpenExchangeRateClient, error) {
	if appID == "" {
		return nil, fmt.Errorf("OpenExchangeRates API key cannot be empty")
	}

	return &openExchangeRateClient{
		appID:   appID,
		baseURL: "https://openexchangerates.org/api/",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *openExchangeRateClient) FetchLatestRate(baseCurrency string, targetCurrencies []string) (map[string]float64, error) {
	if baseCurrency == "" {
		baseCurrency = "USD"
	} else if baseCurrency != "USD" {
		return nil, fmt.Errorf("only USD is supported as base currency")
	}

	endpoint := fmt.Sprintf("%slatest.json?app_id=%s", c.baseURL, c.appID)
	if len(targetCurrencies) > 0 {
		endpoint += "&symbols=" + strings.Join(targetCurrencies, ",")
	}

	res, err := c.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Response parser
	var data struct {
		Rates map[string]float64 `json:"rates"`
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return data.Rates, nil
}
