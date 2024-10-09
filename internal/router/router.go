package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nbonair/currency-exchange-server/internal/handler"
)

func SetupRouter(exchangeRateHandler handler.ExchangeRateHandler) *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api")
	{
		v1Group := apiGroup.Group("/v1")
		{
			v1Group.GET("/exchange-rates", exchangeRateHandler.GetExchangeRates)
		}
	}

	return router
}
