package settings

import "github.com/google/wire"

var (
	SettingRepoProviders = wire.NewSet(
		ProvideSettingsRepository,
	)
)
