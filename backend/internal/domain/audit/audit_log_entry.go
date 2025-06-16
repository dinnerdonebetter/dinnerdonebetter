package auditlogentries

import (
	"time"
)

const (
	AuditLogEventTypeOther    = "other"
	AuditLogEventTypeCreated  = "created"
	AuditLogEventTypeUpdated  = "updated"
	AuditLogEventTypeArchived = "archived"
	AuditLogEventTypeDeleted  = "deleted"
)

type (
	ChangeLog struct {
		OldValue string `json:"oldValue"`
		NewValue string `json:"newValue"`
	}

	AuditLogEntry struct {
		_ struct{} `json:"-"`

		CreatedAt        time.Time            `json:"createdAt"`
		Changes          map[string]ChangeLog `json:"changes"`
		BelongsToAccount *string              `json:"belongsToAccount"`
		ID               string               `json:"id"`
		ResourceType     string               `json:"resourceType"`
		RelevantID       string               `json:"relevantID"`
		EventType        string               `json:"eventType"`
		BelongsToUser    string               `json:"belongsToUser"`
	}

	AuditLogEntryDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		Changes          map[string]ChangeLog `json:"-"`
		BelongsToAccount *string              `json:"-"`
		ID               string               `json:"-"`
		ResourceType     string               `json:"-"`
		RelevantID       string               `json:"-"`
		EventType        string               `json:"-"`
		BelongsToUser    string               `json:"-"`
	}
)
