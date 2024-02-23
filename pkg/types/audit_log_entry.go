package types

import (
	"time"
)

const (
	AuditLogEventTypeOther   = "other"
	AuditLogEventTypeCreated = "created"
	AuditLogEventTypeUpdated = "updated"
	AuditLogEventTypeArchive = "archived"
)

type (
	AuditLogEventType string

	ChangeLog struct {
		OldValue string `json:"oldValue"`
		NewValue string `json:"newValue"`
	}

	AuditLogEntry struct {
		_                  struct{}             `json:"-"`
		CreatedAt          time.Time            `json:"createdAt"`
		Changes            map[string]ChangeLog `json:"changes"`
		BelongsToHousehold *string              `json:"belongsToHousehold"`
		ID                 string               `json:"id"`
		ResourceType       string               `json:"resourceType"`
		RelevantID         string               `json:"relevantID"`
		EventType          AuditLogEventType    `json:"eventType"`
		BelongsToUser      string               `json:"belongsToUser"`
	}
)
