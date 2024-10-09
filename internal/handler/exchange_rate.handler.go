package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nbonair/currency-exchange-server/internal/service"
)

type ExchangeRateHandler interface {
	GetExchangeRates(c *gin.Context)
}

type exchangeRateHandler struct {
	exchangeRateService service.ExchangeRateService
}

func NewExchangeRateHandler(exchangeRateService service.ExchangeRateService) ExchangeRateHandler {
	return &exchangeRateHandler{
		exchangeRateService: exchangeRateService,
	}
}

func (eh *exchangeRateHandler) GetExchangeRates(c *gin.Context) {
	baseCurrency := c.Query("base")
	targetCurrenciesParam := c.Query("targets")
	targetCurrencies := strings.Split(targetCurrenciesParam, ",")

	if baseCurrency == "" || targetCurrenciesParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing input"})
		return
	}

	rates, err := eh.exchangeRateService.FetchLatestRates(c.Request.Context(), baseCurrency, targetCurrencies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"base_currency": baseCurrency,
		"rates":         rates,
	})
}
