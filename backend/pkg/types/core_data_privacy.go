package types

type (
	CoreUserDataCollection struct {
		_ struct{} `json:"-"`

		Webhooks                         map[string][]Webhook                     `json:"webhooks"`
		ServiceSettingConfigurations     map[string][]ServiceSettingConfiguration `json:"serviceSettingConfigurations"`
		AuditLogEntries                  map[string][]AuditLogEntry               `json:"auditLogEntries"`
		ReceivedInvites                  []HouseholdInvitation                    `json:"receivedInvites"`
		SentInvites                      []HouseholdInvitation                    `json:"sentInvites"`
		UserServiceSettingConfigurations []ServiceSettingConfiguration            `json:"userServiceSettingConfigurations"`
		UserAuditLogEntries              []AuditLogEntry                          `json:"userAuditLogEntries"`
		Households                       []Household                              `json:"households"`
	}
)
