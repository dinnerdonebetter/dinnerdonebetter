// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	HouseholdUserMembershipWithUser struct {
		ArchivedAt         string `json:"archivedAt"`
		BelongsToHousehold string `json:"belongsToHousehold"`
		CreatedAt          string `json:"createdAt"`
		HouseholdRole      string `json:"householdRole"`
		ID                 string `json:"id"`
		LastUpdatedAt      string `json:"lastUpdatedAt"`
		BelongsToUser      User   `json:"belongsToUser"`
		DefaultHousehold   bool   `json:"defaultHousehold"`
	}
)
