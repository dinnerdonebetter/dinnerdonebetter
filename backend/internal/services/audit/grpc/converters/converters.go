package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
)

func ConvertAuditLogEntryToGRPCAuditLogEntry(entry *audit.AuditLogEntry) *auditsvc.AuditLogEntry {
	changes := map[string]*auditsvc.ChangeLog{}
	for k, v := range entry.Changes {
		changes[k] = &auditsvc.ChangeLog{
			OldValue: v.OldValue,
			NewValue: v.NewValue,
		}
	}

	return &auditsvc.AuditLogEntry{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(entry.CreatedAt),
		Changes:          changes,
		BelongsToAccount: pointer.Dereference(entry.BelongsToAccount),
		ID:               entry.ID,
		ResourceType:     entry.ResourceType,
		RelevantID:       entry.RelevantID,
		EventType:        entry.EventType,
		BelongsToUser:    entry.BelongsToUser,
	}
}
