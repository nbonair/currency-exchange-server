package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nbonair/currency-exchange-server/internal/router/rate"
)

type Router interface {
	InitRoutes(*gin.Engine)
}

type AppRouter struct {
	ExchangeRateRouter *rate.ExchangeRateRouter
}

func NewAppRouter(exchangeRateRouter *rate.ExchangeRateRouter) *AppRouter {
	return &AppRouter{
		ExchangeRateRouter: exchangeRateRouter,
	}
}

func (ar *AppRouter) InitRoutes(r *gin.Engine) {
	MainGroup := r.Group("/api/v1/2024")
	{
		MainGroup.GET("/check_status", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "OK"})
		})
	}
	{
		ar.ExchangeRateRouter.InitExchangeRateRouter(MainGroup)
	}
}

func NewRouter(appRouter *AppRouter) *gin.Engine {
	r := gin.Default()
	appRouter.InitRoutes(r)
	return r
}
