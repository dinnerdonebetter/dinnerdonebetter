package types

import (
	"context"
	"net/http"
	"time"
)

const (
	AuditLogEventTypeOther    = "other"
	AuditLogEventTypeCreated  = "created"
	AuditLogEventTypeUpdated  = "updated"
	AuditLogEventTypeArchived = "archived"

	AuditLogResourceTypesQueryParamKey = "resources"
)

type (
	AuditLogEntryEventType string

	ChangeLog struct {
		OldValue string `json:"oldValue"`
		NewValue string `json:"newValue"`
	}

	AuditLogEntry struct {
		_ struct{} `json:"-"`

		CreatedAt          time.Time              `json:"createdAt"`
		Changes            map[string]ChangeLog   `json:"changes"`
		BelongsToHousehold *string                `json:"belongsToHousehold"`
		ID                 string                 `json:"id"`
		ResourceType       string                 `json:"resourceType"`
		RelevantID         string                 `json:"relevantID"`
		EventType          AuditLogEntryEventType `json:"eventType"`
		BelongsToUser      string                 `json:"belongsToUser"`
	}

	AuditLogEntryDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		Changes            map[string]ChangeLog   `json:"changes"`
		BelongsToHousehold *string                `json:"belongsToHousehold"`
		ID                 string                 `json:"id"`
		ResourceType       string                 `json:"resourceType"`
		RelevantID         string                 `json:"relevantID"`
		EventType          AuditLogEntryEventType `json:"eventType"`
		BelongsToUser      string                 `json:"belongsToUser"`
	}

	// AuditLogEntryDataManager describes a structure capable of storing audit log entries.
	AuditLogEntryDataManager interface {
		GetAuditLogEntry(ctx context.Context, auditLogID string) (*AuditLogEntry, error)
		GetAuditLogEntriesForUser(ctx context.Context, userID string, filter *QueryFilter) (*QueryFilteredResult[AuditLogEntry], error)
		GetAuditLogEntriesForUserAndResourceType(ctx context.Context, userID string, resourceTypes []string, filter *QueryFilter) (*QueryFilteredResult[AuditLogEntry], error)
		GetAuditLogEntriesForHousehold(ctx context.Context, householdID string, filter *QueryFilter) (*QueryFilteredResult[AuditLogEntry], error)
		GetAuditLogEntriesForHouseholdAndResourceType(ctx context.Context, householdID string, resourceTypes []string, filter *QueryFilter) (*QueryFilteredResult[AuditLogEntry], error)
	}

	// AuditLogEntryDataService describes a structure capable of serving traffic related to audit log entries.
	AuditLogEntryDataService interface {
		ReadAuditLogEntryHandler(http.ResponseWriter, *http.Request)
		ListUserAuditLogEntriesHandler(http.ResponseWriter, *http.Request)
		ListHouseholdAuditLogEntriesHandler(http.ResponseWriter, *http.Request)
	}
)
