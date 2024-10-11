package rate

import (
	"github.com/gin-gonic/gin"
	"github.com/nbonair/currency-exchange-server/internal/handler"
)

type ExchangeRateRouter struct {
	Handler handler.ExchangeRateHandler
}

func NewExchangeRateRouter(handler handler.ExchangeRateHandler) *ExchangeRateRouter {
	return &ExchangeRateRouter{
		Handler: handler,
	}
}

func (er *ExchangeRateRouter) InitExchangeRateRouter(Router *gin.RouterGroup) {
	exchangeRateRouterPublic := Router.Group("/exchange-rate")
	{
		exchangeRateRouterPublic.GET("/all", er.Handler.GetExchangeRates)
	}
}
