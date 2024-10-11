package openexchangerates

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewOpenExchangeRateClient,
)
