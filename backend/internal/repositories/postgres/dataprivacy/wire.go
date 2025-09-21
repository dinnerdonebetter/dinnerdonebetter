package dataprivacy

import "github.com/google/wire"

var (
	DataPrivRepoProviders = wire.NewSet(
		ProvideDataPrivacyRepository,
	)
)
