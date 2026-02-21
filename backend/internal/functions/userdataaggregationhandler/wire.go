package userdataaggregationhandler

import (
	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		NewUserDataAggregationHandler,
	)
)
