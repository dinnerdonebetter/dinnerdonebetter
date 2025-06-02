// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	HouseholdInstrumentOwnership struct {
		ArchivedAt         string          `json:"archivedAt"`
		BelongsToHousehold string          `json:"belongsToHousehold"`
		CreatedAt          string          `json:"createdAt"`
		ID                 string          `json:"id"`
		LastUpdatedAt      string          `json:"lastUpdatedAt"`
		Notes              string          `json:"notes"`
		Instrument         ValidInstrument `json:"instrument"`
		Quantity           uint64          `json:"quantity"`
	}
)
