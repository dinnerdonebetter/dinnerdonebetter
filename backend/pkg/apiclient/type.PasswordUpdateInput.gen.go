// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	PasswordUpdateInput struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
		TotpToken       string `json:"totpToken"`
	}
)
