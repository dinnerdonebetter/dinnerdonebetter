package manager

import "github.com/google/wire"

var (
	AuditManagerProviders = wire.NewSet(
		NewAuditDataManager,
	)
)
