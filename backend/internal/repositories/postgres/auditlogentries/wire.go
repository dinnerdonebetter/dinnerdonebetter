package auditlogentries

import "github.com/google/wire"

var (
	Providers = wire.NewSet(
		ProvideAuditLogRepository,
	)
)
