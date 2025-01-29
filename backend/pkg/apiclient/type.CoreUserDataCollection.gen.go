// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
CoreUserDataCollection struct {
   AuditLogEntries map[string]any `json:"auditLogEntries"`
 Households []Household `json:"households"`
 ReceivedInvites []HouseholdInvitation `json:"receivedInvites"`
 SentInvites []HouseholdInvitation `json:"sentInvites"`
 ServiceSettingConfigurations map[string]any `json:"serviceSettingConfigurations"`
 UserAuditLogEntries []AuditLogEntry `json:"userAuditLogEntries"`
 UserServiceSettingConfigurations []ServiceSettingConfiguration `json:"userServiceSettingConfigurations"`
 Webhooks map[string]any `json:"webhooks"`

}
)
