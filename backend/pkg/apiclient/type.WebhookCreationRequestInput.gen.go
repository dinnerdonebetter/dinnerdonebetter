// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	WebhookCreationRequestInput struct {
		ContentType string   `json:"contentType"`
		Method      string   `json:"method"`
		Name        string   `json:"name"`
		URL         string   `json:"url"`
		Events      []string `json:"events"`
	}
)
