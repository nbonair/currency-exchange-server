//go:build wireinject

package wiring

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/nbonair/currency-exchange-server/configs"
	"github.com/nbonair/currency-exchange-server/internal/dataaccess"
	"github.com/nbonair/currency-exchange-server/internal/handler"
	"github.com/nbonair/currency-exchange-server/internal/lib"
	"github.com/nbonair/currency-exchange-server/internal/repo"
	"github.com/nbonair/currency-exchange-server/internal/router"
	"github.com/nbonair/currency-exchange-server/internal/service"
)

var RepositorySet = wire.NewSet(
	dataaccess.WireSet,
	lib.WireSet,
	repo.WireSet,
)

var ServiceSet = wire.NewSet(
	service.NewExchangeRateService,
)

var HandlerSet = wire.NewSet(
	handler.NewExchangeRateHandler,
)

var ApplicationSet = wire.NewSet(
	RepositorySet,
	ServiceSet,
	HandlerSet,
)

func InitializeRouter(cfgDb configs.DatabaseConfig, cfgApi configs.APIsConfig, cfgCache configs.CacheConfig) (*gin.Engine, func(), error) {
	wire.Build(
		ApplicationSet,
		router.WireSet,
	)
	return nil, nil, nil
}
