package datachangemessagehandler

import (
	"github.com/google/wire"
	"github.com/verygoodsoftwarenotvirus/platform/routing/chi"
)

var ProvidersMiscellaneous = wire.NewSet(
	chi.NewRouteParamManager,
)
