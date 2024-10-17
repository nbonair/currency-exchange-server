package dataaccess

import (
	"github.com/google/wire"
	"github.com/nbonair/currency-exchange-server/internal/dataaccess/cache"
	"github.com/nbonair/currency-exchange-server/internal/dataaccess/db"
)

var WireSet = wire.NewSet(
	db.WireSet,
	cache.WireSet,
)
