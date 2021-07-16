package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// ReportAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	ReportAssignmentKey = "report_id"
	// ReportCreationEvent is the event type used to indicate an item was created.
	ReportCreationEvent = "report_created"
	// ReportUpdateEvent is the event type used to indicate an item was updated.
	ReportUpdateEvent = "report_updated"
	// ReportArchiveEvent is the event type used to indicate an item was archived.
	ReportArchiveEvent = "report_archived"
)

// BuildReportCreationEventEntry builds an entry creation input for when a report is created.
func BuildReportCreationEventEntry(report *types.Report, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ReportCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:    createdByUser,
			ReportAssignmentKey:   report.ID,
			CreationAssignmentKey: report,
			AccountAssignmentKey:  report.BelongsToAccount,
		},
	}
}

// BuildReportUpdateEventEntry builds an entry creation input for when a report is updated.
func BuildReportUpdateEventEntry(changedByUser, reportID, accountID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ReportUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:   changedByUser,
			AccountAssignmentKey: accountID,
			ReportAssignmentKey:  reportID,
			ChangesAssignmentKey: changes,
		},
	}
}

// BuildReportArchiveEventEntry builds an entry creation input for when a report is archived.
func BuildReportArchiveEventEntry(archivedByUser, accountID, reportID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ReportArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:   archivedByUser,
			AccountAssignmentKey: accountID,
			ReportAssignmentKey:  reportID,
		},
	}
}
