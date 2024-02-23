package types

import (
	"context"
	"net/http"
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
		_ struct{} `json:"-"`

		CreatedAt          time.Time            `json:"createdAt"`
		Changes            map[string]ChangeLog `json:"changes"`
		BelongsToHousehold *string              `json:"belongsToHousehold"`
		ID                 string               `json:"id"`
		ResourceType       string               `json:"resourceType"`
		RelevantID         string               `json:"relevantID"`
		EventType          AuditLogEventType    `json:"eventType"`
		BelongsToUser      string               `json:"belongsToUser"`
	}

	AuditLogEntryDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		Changes            map[string]ChangeLog `json:"changes"`
		BelongsToHousehold *string              `json:"belongsToHousehold"`
		ID                 string               `json:"id"`
		ResourceType       string               `json:"resourceType"`
		RelevantID         string               `json:"relevantID"`
		EventType          AuditLogEventType    `json:"eventType"`
		BelongsToUser      string               `json:"belongsToUser"`
	}

	// AuditLogEntryDataManager describes a structure capable of storing audit log entries.
	AuditLogEntryDataManager interface {
		GetAuditLogEntry(ctx context.Context, auditLogID, householdID string) (*AuditLogEntry, error)
		GetAuditLogEntriesForUser(ctx context.Context, householdID string, filter *QueryFilter) (*QueryFilteredResult[AuditLogEntry], error)
		GetAuditLogEntriesForUserAndResourceType(ctx context.Context, user, resourceType string, filter *QueryFilter) (*QueryFilteredResult[AuditLogEntry], error)
		GetAuditLogEntriesForHousehold(ctx context.Context, householdID string, filter *QueryFilter) (*QueryFilteredResult[AuditLogEntry], error)
		GetAuditLogEntriesForHouseholdAndResourceType(ctx context.Context, householdID, resourceType string, filter *QueryFilter) (*QueryFilteredResult[AuditLogEntry], error)
		CreateAuditLogEntry(ctx context.Context, input *AuditLogEntryDatabaseCreationInput) (*AuditLogEntry, error)
	}

	// AuditLogEntryDataService describes a structure capable of serving traffic related to audit log entries.
	AuditLogEntryDataService interface {
		ListUserAuditLogEntriesHandler(http.ResponseWriter, *http.Request)
		ListUserAuditLogEntriesForResourceTypeHandler(http.ResponseWriter, *http.Request)
		ListHouseholdAuditLogEntriesHandler(http.ResponseWriter, *http.Request)
		ListHouseholdAuditLogEntriesForResourceTypeHandler(http.ResponseWriter, *http.Request)
		CreateAuditLogEntryHandler(http.ResponseWriter, *http.Request)
		ReadAuditLogEntryHandler(http.ResponseWriter, *http.Request)
		ArchiveAuditLogEntryHandler(http.ResponseWriter, *http.Request)
	}
)
