package types

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// SessionContextDataKey is the non-string type we use for referencing SessionContextData structs.

	// UserRegistrationInputContextKey is the non-string type we use for referencing SessionContextData structs.
	UserRegistrationInputContextKey routing.ContextKey = "user_registration_input"

	// TwoFactorSecretVerifiedServiceEventType indicates a user's two factor secret was verified.
	/* #nosec G101 */
	TwoFactorSecretVerifiedServiceEventType = "two_factor_secret_verified"
	// TwoFactorDeactivatedServiceEventType indicates a user's two factor secret was changed and verified_at timestamp was reset.
	/* #nosec G101 */
	TwoFactorDeactivatedServiceEventType = "two_factor_deactivated"
	// TwoFactorSecretChangedServiceEventType indicates a user's two factor secret was changed and verified_at timestamp was reset.
	/* #nosec G101 */
	TwoFactorSecretChangedServiceEventType = "two_factor_secret_changed"
	// PasswordResetTokenCreatedEventType indicates a user created a password reset token.
	PasswordResetTokenCreatedEventType = "password_reset_token_created"
	// PasswordResetTokenRedeemedEventType indicates a user created a password reset token.
	PasswordResetTokenRedeemedEventType = "password_reset_token_redeemed"
	// PasswordChangedEventType indicates a user changed their password.
	PasswordChangedEventType = "password_changed"
	// EmailAddressChangedEventType indicates a user changed their email address.
	EmailAddressChangedEventType = "email_address_changed"
	// UsernameChangedEventType indicates a user changed their username.
	UsernameChangedEventType = "username_changed"
	// UserDetailsChangedEventType indicates a user changed their information.
	UserDetailsChangedEventType = "user_details_changed"
	// UsernameReminderRequestedEventType indicates a user requested a username reminder.
	UsernameReminderRequestedEventType = "username_reminder_requested"
	// UserLoggedInServiceEventType indicates a user has logged in.
	UserLoggedInServiceEventType = "user_logged_in"
	// UserRefreshedTokenServiceEventType indicates a user has refreshed a token.
	UserRefreshedTokenServiceEventType = "user_refreshed_token"
	// UserLoggedOutServiceEventType indicates a user has logged in.
	UserLoggedOutServiceEventType = "user_logged_out"
	// UserChangedActiveHouseholdServiceEventType indicates a user has logged in.
	UserChangedActiveHouseholdServiceEventType = "changed_active_household"
	// UserEmailAddressVerifiedEventType indicates a user created a password reset token.
	UserEmailAddressVerifiedEventType = "user_email_address_verified"
	// UserEmailAddressVerificationEmailRequestedEventType indicates a user created a password reset token.
	UserEmailAddressVerificationEmailRequestedEventType = "user_email_address_verification_email_requested"
)

type (
	// UserStatusResponse is what we encode when the frontend wants to check auth status.
	UserStatusResponse struct {
		_ struct{} `json:"-"`

		UserID                   string `json:"userID"`
		AccountStatus            string `json:"accountStatus,omitempty"`
		AccountStatusExplanation string `json:"accountStatusExplanation"`
		ActiveHousehold          string `json:"activeHousehold,omitempty"`
		UserIsAuthenticated      bool   `json:"isAuthenticated"`
	}

	// UserPermissionsRequestInput is what we decode when the frontend wants to check permission status.
	UserPermissionsRequestInput struct {
		_ struct{} `json:"-"`

		Permissions []string `json:"permissions"`
	}

	// UserPermissionsResponse is what we encode when the frontend wants to check permission status.
	UserPermissionsResponse struct {
		_ struct{} `json:"-"`

		Permissions map[string]bool `json:"permissions"`
	}

	// ChangeActiveHouseholdInput represents what a User could set as input for switching households.
	ChangeActiveHouseholdInput struct {
		_ struct{} `json:"-"`

		HouseholdID string `json:"householdID"`
	}

	// AuthDataService describes a structure capable of handling passwords and authorization requests.
	AuthDataService interface {
		StatusHandler(http.ResponseWriter, *http.Request)
		BuildLoginHandler(adminOnly bool) func(http.ResponseWriter, *http.Request)

		SSOLoginHandler(http.ResponseWriter, *http.Request)
		SSOLoginCallbackHandler(http.ResponseWriter, *http.Request)

		PermissionFilterMiddleware(permissions ...authorization.Permission) func(next http.Handler) http.Handler
		UserAttributionMiddleware(next http.Handler) http.Handler
		AuthorizationMiddleware(next http.Handler) http.Handler
		ServiceAdminMiddleware(next http.Handler) http.Handler

		OAuth2Service
	}
)

var _ validation.ValidatableWithContext = (*ChangeActiveHouseholdInput)(nil)

// ValidateWithContext validates a ChangeActiveHouseholdInput.
func (x *ChangeActiveHouseholdInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.HouseholdID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*UserPermissionsRequestInput)(nil)

// ValidateWithContext validates a UserPermissionsRequestInput.
func (x *UserPermissionsRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Permissions, validation.Required),
	)
}
