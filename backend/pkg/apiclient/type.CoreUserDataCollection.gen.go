// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	CoreUserDataCollection struct {
		AuditLogEntries                  map[string]any                `json:"auditLogEntries"`
		ServiceSettingConfigurations     map[string]any                `json:"serviceSettingConfigurations"`
		Webhooks                         map[string]any                `json:"webhooks"`
		Households                       []Household                   `json:"households"`
		ReceivedInvites                  []HouseholdInvitation         `json:"receivedInvites"`
		SentInvites                      []HouseholdInvitation         `json:"sentInvites"`
		UserAuditLogEntries              []AuditLogEntry               `json:"userAuditLogEntries"`
		UserServiceSettingConfigurations []ServiceSettingConfiguration `json:"userServiceSettingConfigurations"`
	}
)
