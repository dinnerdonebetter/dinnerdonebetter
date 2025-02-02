// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	HouseholdInstrumentOwnershipCreationRequestInput struct {
		BelongsToHousehold string `json:"belongsToHousehold"`
		Notes              string `json:"notes"`
		ValidInstrumentID  string `json:"validInstrumentID"`
		Quantity           uint64 `json:"quantity"`
	}
)
