package admin

type (
	contextKey string
)

const (
	apiClientContextKey contextKey = "api_client"
)

type userSessionDetails struct {
	Token       string `json:"token"`
	UserID      string `json:"userID"`
	HouseholdID string `json:"householdID"`
}
