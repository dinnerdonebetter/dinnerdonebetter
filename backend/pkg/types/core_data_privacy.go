package types

type (
	CoreUserDataCollection struct {
		_ struct{} `json:"-"`

		Webhooks                         map[string][]Webhook                     `json:"webhooks"`
		ServiceSettingConfigurations     map[string][]ServiceSettingConfiguration `json:"serviceSettingConfigurations"`
		AuditLogEntries                  map[string][]AuditLogEntry               `json:"auditLogEntries"`
		ReceivedInvites                  []AccountInvitation                      `json:"receivedInvites"`
		SentInvites                      []AccountInvitation                      `json:"sentInvites"`
		UserServiceSettingConfigurations []ServiceSettingConfiguration            `json:"userServiceSettingConfigurations"`
		UserAuditLogEntries              []AuditLogEntry                          `json:"userAuditLogEntries"`
		Accounts                         []Account                                `json:"accounts"`
	}
)
