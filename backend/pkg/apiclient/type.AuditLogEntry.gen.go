// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	AuditLogEntry struct {
		BelongsToHousehold string    `json:"belongsToHousehold"`
		BelongsToUser      string    `json:"belongsToUser"`
		Changes            ChangeLog `json:"changes"`
		CreatedAt          string    `json:"createdAt"`
		EventType          string    `json:"eventType"`
		ID                 string    `json:"id"`
		RelevantID         string    `json:"relevantID"`
		ResourceType       string    `json:"resourceType"`
	}
)
