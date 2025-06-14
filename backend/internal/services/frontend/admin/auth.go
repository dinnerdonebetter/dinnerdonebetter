package admin

type (
	contextKey string
)

const (
	apiClientContextKey contextKey = "api_client"
)

type userSessionDetails struct {
	Token     string `json:"token"`
	UserID    string `json:"userID"`
	AccountID string `json:"accountID"`
}
