// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
HouseholdInvitation struct {
   ArchivedAt string `json:"archivedAt"`
 CreatedAt string `json:"createdAt"`
 DestinationHousehold Household `json:"destinationHousehold"`
 ExpiresAt string `json:"expiresAt"`
 FromUser User `json:"fromUser"`
 ID string `json:"id"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Note string `json:"note"`
 Status string `json:"status"`
 StatusNote string `json:"statusNote"`
 ToEmail string `json:"toEmail"`
 ToName string `json:"toName"`
 ToUser string `json:"toUser"`
 Token string `json:"token"`

}
)
