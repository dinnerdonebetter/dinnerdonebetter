package manager

import "github.com/google/wire"

var (
	DataPrivacyManagerProviders = wire.NewSet(
		NewDataPrivacyManager,
	)
)
