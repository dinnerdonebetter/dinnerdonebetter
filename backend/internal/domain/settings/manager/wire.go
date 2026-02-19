package manager

import (
	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	settingsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/settings"

	"github.com/google/wire"
)

var (
	SettingsManagerProviders = wire.NewSet(
		NewSettingsDataManager,
		settingsrepo.ProvideSettingsRepository,
		wire.Bind(new(settingsRepo), new(*settingsrepo.Repository)),
		wire.Bind(new(settings.Repository), new(SettingsDataManager)),
	)
)
