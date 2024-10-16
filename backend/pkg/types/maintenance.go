package types

import "context"

type (
	MaintenanceDataManager interface {
		DeleteExpiredOAuth2ClientTokens(context.Context) error
	}
)
