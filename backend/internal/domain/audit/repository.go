package auditlogentries

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
)

// AuditLogEntryDataManager describes a structure capable of storing audit log entries.
type AuditLogEntryDataManager interface {
	GetAuditLogEntry(ctx context.Context, auditLogID string) (*AuditLogEntry, error)
	GetAuditLogEntriesForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[AuditLogEntry], error)
	GetAuditLogEntriesForUserAndResourceTypes(ctx context.Context, userID string, resourceTypes []string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[AuditLogEntry], error)
	GetAuditLogEntriesForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[AuditLogEntry], error)
	GetAuditLogEntriesForAccountAndResourceTypes(ctx context.Context, accountID string, resourceTypes []string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[AuditLogEntry], error)
}
