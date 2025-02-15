// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	Webhook struct {
		ArchivedAt         string                `json:"archivedAt"`
		BelongsToHousehold string                `json:"belongsToHousehold"`
		ContentType        string                `json:"contentType"`
		CreatedAt          string                `json:"createdAt"`
		ID                 string                `json:"id"`
		LastUpdatedAt      string                `json:"lastUpdatedAt"`
		Method             string                `json:"method"`
		Name               string                `json:"name"`
		URL                string                `json:"url"`
		Events             []WebhookTriggerEvent `json:"events"`
	}
)
