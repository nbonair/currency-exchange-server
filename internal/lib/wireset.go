package lib

import (
	"github.com/google/wire"
	"github.com/nbonair/currency-exchange-server/internal/lib/openexchangerates"
)

var WireSet = wire.NewSet(
	openexchangerates.WireSet,
)
