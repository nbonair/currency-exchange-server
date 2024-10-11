package rate

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewExchangeRateRouter,
)
