// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	HouseholdInvitation struct {
		LastUpdatedAt        string    `json:"lastUpdatedAt"`
		Status               string    `json:"status"`
		Token                string    `json:"token"`
		ExpiresAt            string    `json:"expiresAt"`
		ToUser               string    `json:"toUser"`
		ID                   string    `json:"id"`
		CreatedAt            string    `json:"createdAt"`
		ArchivedAt           string    `json:"archivedAt"`
		ToEmail              string    `json:"toEmail"`
		StatusNote           string    `json:"statusNote"`
		Note                 string    `json:"note"`
		ToName               string    `json:"toName"`
		FromUser             User      `json:"fromUser"`
		DestinationHousehold Household `json:"destinationHousehold"`
	}
)
