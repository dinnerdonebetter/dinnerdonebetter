package datachangemessagehandler

import (
	"github.com/dinnerdonebetter/backend/internal/platform/routing/chi"

	"github.com/google/wire"
)

var ProvidersMiscellaneous = wire.NewSet(
	chi.NewRouteParamManager,
)
