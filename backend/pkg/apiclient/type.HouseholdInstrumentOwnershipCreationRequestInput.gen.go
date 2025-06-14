// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	AccountInstrumentOwnershipCreationRequestInput struct {
		BelongsToAccount  string `json:"belongsToAccount"`
		Notes             string `json:"notes"`
		ValidInstrumentID string `json:"validInstrumentID"`
		Quantity          uint64 `json:"quantity"`
	}
)
