package router

import (
	"github.com/google/wire"
	"github.com/nbonair/currency-exchange-server/internal/router/rate"
)

var WireSet = wire.NewSet(
	rate.WireSet,
	NewAppRouter,
	NewRouter,
)
