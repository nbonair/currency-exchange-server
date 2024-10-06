package openexchangerates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OpenExchangeRatesClient struct {
	appId   string
	baseURL string
	Client  *http.Client
}

func NewOpenExchangeRatesClient(appId string) (*OpenExchangeRatesClient, error) {
	return &OpenExchangeRatesClient{
		appId:   appId,
		baseURL: "https://openexchangerates.org/api/",
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *OpenExchangeRatesClient) FetchLatestRate(baseCurrency string, symbols []string) (map[string]float64, error) {
	if baseCurrency != "USD" && baseCurrency != "" {
		return nil, fmt.Errorf("only accept usd as base currency")
	}

	endpoint := fmt.Sprintf("%s/latest.json?app_id=%s", c.baseURL, c.appId)
	if len(symbols) > 0 {
		endpoint += "&symbol=" + strings.Join(symbols, ",")
	}

	res, err := c.Client.Get(endpoint)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error code %d", res.StatusCode)
	}

	//Response parser
	var data struct {
		Rates map[string]float64 `json:"rates"`
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return data.Rates, nil

}
