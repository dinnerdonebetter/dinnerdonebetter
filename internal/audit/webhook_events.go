package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// WebhookAssignmentKey is the key we use to indicate that an audit log entry is associated with a webhook.
	WebhookAssignmentKey = "webhook_id"

	// WebhookCreationEvent events indicate a user created a webhook.
	WebhookCreationEvent = "webhook_created"
	// WebhookUpdateEvent events indicate a user updated a webhook.
	WebhookUpdateEvent = "webhook_updated"
	// WebhookArchiveEvent events indicate a user deleted a webhook.
	WebhookArchiveEvent = "webhook_archived"
)

// BuildWebhookCreationEventEntry builds an entry creation input for when a webhook is created.
func BuildWebhookCreationEventEntry(webhook *types.Webhook, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: WebhookCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:     createdByUser,
			CreationAssignmentKey:  webhook,
			WebhookAssignmentKey:   webhook.ID,
			HouseholdAssignmentKey: webhook.BelongsToHousehold,
		},
	}
}

// BuildWebhookUpdateEventEntry builds an entry creation input for when a webhook is updated.
func BuildWebhookUpdateEventEntry(changedByUser, householdID, webhookID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: WebhookUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:     changedByUser,
			HouseholdAssignmentKey: householdID,
			WebhookAssignmentKey:   webhookID,
			ChangesAssignmentKey:   changes,
		},
	}
}

// BuildWebhookArchiveEventEntry builds an entry creation input for when a webhook is archived.
func BuildWebhookArchiveEventEntry(archivedByUser, householdID, webhookID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: WebhookArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:     archivedByUser,
			HouseholdAssignmentKey: householdID,
			WebhookAssignmentKey:   webhookID,
		},
	}
}
