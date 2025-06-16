package auditlogentries

import (
	"net/http"
)

// AuditLogEntryDataService describes a structure capable of serving traffic related to audit log entries.
type AuditLogEntryDataService interface {
	ReadAuditLogEntryHandler(http.ResponseWriter, *http.Request)
	ListUserAuditLogEntriesHandler(http.ResponseWriter, *http.Request)
	ListAccountAuditLogEntriesHandler(http.ResponseWriter, *http.Request)
}
