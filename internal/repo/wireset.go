package repo

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewExchangeRateRepository,
	NewExchangeRateHistoryRepository,
)
