package types

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// SessionContextDataKey is the non-string type we use for referencing SessionContextData structs.
	SessionContextDataKey ContextKey = "session_context_data"
	// UserIDContextKey is the non-string type we use for referencing SessionContextData structs.
	UserIDContextKey ContextKey = "user_id"
	// HouseholdIDContextKey is the non-string type we use for referencing SessionContextData structs.
	HouseholdIDContextKey ContextKey = "household_id"
	// UserRegistrationInputContextKey is the non-string type we use for referencing SessionContextData structs.
	UserRegistrationInputContextKey ContextKey = "user_registration_input"

	// TwoFactorSecretVerifiedServiceEventType indicates a user's two factor secret was verified.
	/* #nosec G101 */
	TwoFactorSecretVerifiedServiceEventType ServiceEventType = "two_factor_secret_verified"
	// TwoFactorDeactivatedServiceEventType indicates a user's two factor secret was changed and verified_at timestamp was reset.
	/* #nosec G101 */
	TwoFactorDeactivatedServiceEventType ServiceEventType = "two_factor_deactivated"
	// TwoFactorSecretChangedServiceEventType indicates a user's two factor secret was changed and verified_at timestamp was reset.
	/* #nosec G101 */
	TwoFactorSecretChangedServiceEventType ServiceEventType = "two_factor_secret_changed"
	// PasswordResetTokenCreatedEventType indicates a user created a password reset token.
	PasswordResetTokenCreatedEventType ServiceEventType = "password_reset_token_created"
	// PasswordResetTokenRedeemedEventType indicates a user created a password reset token.
	PasswordResetTokenRedeemedEventType ServiceEventType = "password_reset_token_redeemed"
	// PasswordChangedEventType indicates a user changed their password.
	PasswordChangedEventType ServiceEventType = "password_changed"
	// EmailAddressChangedEventType indicates a user changed their email address.
	EmailAddressChangedEventType ServiceEventType = "email_address_changed"
	// UsernameChangedEventType indicates a user changed their username.
	UsernameChangedEventType ServiceEventType = "username_changed"
	// UserDetailsChangedEventType indicates a user changed their information.
	UserDetailsChangedEventType ServiceEventType = "user_details_changed"
	// UsernameReminderRequestedEventType indicates a user requested a username reminder.
	UsernameReminderRequestedEventType ServiceEventType = "username_reminder_requested"
	// UserLoggedInServiceEventType indicates a user has logged in.
	UserLoggedInServiceEventType ServiceEventType = "user_logged_in"
	// UserLoggedOutServiceEventType indicates a user has logged in.
	UserLoggedOutServiceEventType ServiceEventType = "user_logged_out"
	// UserChangedActiveHouseholdServiceEventType indicates a user has logged in.
	UserChangedActiveHouseholdServiceEventType ServiceEventType = "changed_active_household"
	// UserEmailAddressVerifiedEventType indicates a user created a password reset token.
	UserEmailAddressVerifiedEventType ServiceEventType = "user_email_address_verified"
	// UserEmailAddressVerificationEmailRequestedEventType indicates a user created a password reset token.
	UserEmailAddressVerificationEmailRequestedEventType ServiceEventType = "user_email_address_verification_email_requested"
)

func init() {
	gob.Register(&SessionContextData{})
}

type (
	// SessionContextData represents what we encode in our passwords cookies.
	SessionContextData struct {
		_ struct{} `json:"-"`

		HouseholdPermissions map[string]authorization.HouseholdRolePermissionsChecker `json:"-"`
		Requester            RequesterInfo                                            `json:"-"`
		ActiveHouseholdID    string                                                   `json:"-"`
	}

	// RequesterInfo contains data relevant to the user making a request.
	RequesterInfo struct {
		_ struct{} `json:"-"`

		ServicePermissions       authorization.ServiceRolePermissionChecker `json:"-"`
		AccountStatus            string                                     `json:"-"`
		AccountStatusExplanation string                                     `json:"-"`
		UserID                   string                                     `json:"-"`
		EmailAddress             string                                     `json:"-"`
		Username                 string                                     `json:"-"`
	}

	// UserStatusResponse is what we encode when the frontend wants to check auth status.
	UserStatusResponse struct {
		_ struct{} `json:"-"`

		UserID                   string `json:"userID"`
		AccountStatus            string `json:"accountStatus,omitempty"`
		AccountStatusExplanation string `json:"accountStatusExplanation"`
		ActiveHousehold          string `json:"activeHousehold,omitempty"`
		UserIsAuthenticated      bool   `json:"isAuthenticated"`
	}

	// JWTResponse is used to return a JWT to a user.
	JWTResponse struct {
		_ struct{} `json:"-"`

		UserID      string `json:"userID"`
		HouseholdID string `json:"householdID"`
		Token       string `json:"token"`
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

// HouseholdRolePermissionsChecker returns the relevant HouseholdRolePermissionsChecker.
func (x *SessionContextData) HouseholdRolePermissionsChecker() authorization.HouseholdRolePermissionsChecker {
	return x.HouseholdPermissions[x.ActiveHouseholdID]
}

// ServiceRolePermissionChecker returns the relevant ServiceRolePermissionChecker.
func (x *SessionContextData) ServiceRolePermissionChecker() authorization.ServiceRolePermissionChecker {
	return x.Requester.ServicePermissions
}

// AttachToLogger provides a consistent way to attach a SessionContextData object to a logger.
func (x *SessionContextData) AttachToLogger(logger logging.Logger) logging.Logger {
	if x != nil {
		logger = logger.WithValue(keys.RequesterIDKey, x.Requester.UserID).
			WithValue(keys.ActiveHouseholdIDKey, x.ActiveHouseholdID)
	}

	return logger
}
