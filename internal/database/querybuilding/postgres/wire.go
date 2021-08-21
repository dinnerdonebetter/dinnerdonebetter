package postgres

import (
	"github.com/google/wire"
)

var (
	// Providers are what we provide for dependency injection.
	Providers = wire.NewSet(
		ProvidePostgresDB,
		ProvidePostgres,
	)
)
