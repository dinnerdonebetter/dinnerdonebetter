package settings

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		ProvideSettingsRepository,
	)
)
