// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	AccountInstrumentOwnership struct {
		ArchivedAt       string          `json:"archivedAt"`
		BelongsToAccount string          `json:"belongsToAccount"`
		CreatedAt        string          `json:"createdAt"`
		ID               string          `json:"id"`
		LastUpdatedAt    string          `json:"lastUpdatedAt"`
		Notes            string          `json:"notes"`
		Instrument       ValidInstrument `json:"instrument"`
		Quantity         uint64          `json:"quantity"`
	}
)
