// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	HouseholdInvitation struct {
		LastUpdatedAt        string    `json:"lastUpdatedAt"`
		ToEmail              string    `json:"toEmail"`
		Token                string    `json:"token"`
		ExpiresAt            string    `json:"expiresAt"`
		Note                 string    `json:"note"`
		ID                   string    `json:"id"`
		Status               string    `json:"status"`
		ToUser               string    `json:"toUser"`
		CreatedAt            string    `json:"createdAt"`
		StatusNote           string    `json:"statusNote"`
		ArchivedAt           string    `json:"archivedAt"`
		ToName               string    `json:"toName"`
		FromUser             User      `json:"fromUser"`
		DestinationHousehold Household `json:"destinationHousehold"`
	}
)
