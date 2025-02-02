// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	User struct {
		FirstName                 string `json:"firstName"`
		Username                  string `json:"username"`
		AccountStatus             string `json:"accountStatus"`
		Avatar                    string `json:"avatar"`
		Birthday                  string `json:"birthday"`
		CreatedAt                 string `json:"createdAt"`
		EmailAddress              string `json:"emailAddress"`
		EmailAddressVerifiedAt    string `json:"emailAddressVerifiedAt"`
		ArchivedAt                string `json:"archivedAt"`
		AccountStatusExplanation  string `json:"accountStatusExplanation"`
		LastAcceptedTos           string `json:"lastAcceptedTOS"`
		LastAcceptedPrivacyPolicy string `json:"lastAcceptedPrivacyPolicy"`
		LastName                  string `json:"lastName"`
		LastUpdatedAt             string `json:"lastUpdatedAt"`
		PasswordLastChangedAt     string `json:"passwordLastChangedAt"`
		ID                        string `json:"id"`
		ServiceRoles              string `json:"serviceRoles"`
		TwoFactorSecretVerifiedAt string `json:"twoFactorSecretVerifiedAt"`
		RequiresPasswordChange    bool   `json:"requiresPasswordChange"`
	}
)
