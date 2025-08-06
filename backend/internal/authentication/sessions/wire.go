package sessions

import (
	"context"
	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideContextDataFetcherFromContext,
	)
)

func ProvideContextDataFetcherFromContext() func(context.Context) (*ContextData, error) {
	return FetchContextDataFromContext
}
