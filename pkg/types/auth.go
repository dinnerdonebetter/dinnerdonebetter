package types

import (
	"bytes"
	"context"
	"encoding/gob"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// SessionContextDataKey is the non-string type we use for referencing SessionContextData structs.
	SessionContextDataKey ContextKey = "session_context_data"
	// UserIDContextKey is the non-string type we use for referencing SessionContextData structs.
	UserIDContextKey ContextKey = "user_id"
	// HouseholdIDContextKey is the non-string type we use for referencing SessionContextData structs.
	HouseholdIDContextKey ContextKey = "household_id"
	// UserLoginInputContextKey is the non-string type we use for referencing SessionContextData structs.
	UserLoginInputContextKey ContextKey = "user_login_input"
	// UserRegistrationInputContextKey is the non-string type we use for referencing SessionContextData structs.
	UserRegistrationInputContextKey ContextKey = "user_registration_input"
)

func init() {
	gob.Register(&SessionContextData{})
}

type (
	// UserHouseholdMembershipInfo represents key information about an household membership.
	UserHouseholdMembershipInfo struct {
		HouseholdName  string   `json:"name"`
		HouseholdRoles []string `json:"-"`
		HouseholdID    uint64   `json:"householdID"`
	}

	// SessionContextData represents what we encode in our passwords cookies.
	SessionContextData struct {
		HouseholdPermissions map[uint64]authorization.HouseholdRolePermissionsChecker `json:"-"`
		Requester            RequesterInfo                                            `json:"-"`
		ActiveHouseholdID    uint64                                                   `json:"-"`
	}

	// RequesterInfo contains data relevant to the user making a request.
	RequesterInfo struct {
		ServicePermissions    authorization.ServiceRolePermissionChecker `json:"-"`
		Reputation            householdStatus                            `json:"-"`
		ReputationExplanation string                                     `json:"-"`
		UserID                uint64                                     `json:"-"`
	}

	// UserStatusResponse is what we encode when the frontend wants to check auth status.
	UserStatusResponse struct {
		UserReputation            householdStatus `json:"householdStatus,omitempty"`
		UserReputationExplanation string          `json:"reputationExplanation"`
		ActiveHousehold           uint64          `json:"activeHousehold,omitempty"`
		UserIsServiceAdmin        bool            `json:"userIsServiceAdmin"`
		UserIsAuthenticated       bool            `json:"isAuthenticated"`
	}

	// ChangeActiveHouseholdInput represents what a User could set as input for switching households.
	ChangeActiveHouseholdInput struct {
		HouseholdID uint64 `json:"householdID"`
	}

	// PASETOCreationInput is used to create a PASETO.
	PASETOCreationInput struct {
		ClientID          string `json:"clientID"`
		HouseholdID       uint64 `json:"householdID"`
		RequestTime       int64  `json:"requestTime"`
		RequestedLifetime uint64 `json:"requestedLifetime,omitempty"`
	}

	// PASETOResponse is used to respond to a PASETO request.
	PASETOResponse struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expiresAt"`
	}

	// AuthService describes a structure capable of handling passwords and authorization requests.
	AuthService interface {
		StatusHandler(res http.ResponseWriter, req *http.Request)
		BeginSessionHandler(res http.ResponseWriter, req *http.Request)
		EndSessionHandler(res http.ResponseWriter, req *http.Request)
		CycleCookieSecretHandler(res http.ResponseWriter, req *http.Request)
		PASETOHandler(res http.ResponseWriter, req *http.Request)
		ChangeActiveHouseholdHandler(res http.ResponseWriter, req *http.Request)

		PermissionFilterMiddleware(permissions ...authorization.Permission) func(next http.Handler) http.Handler
		CookieRequirementMiddleware(next http.Handler) http.Handler
		UserAttributionMiddleware(next http.Handler) http.Handler
		AuthorizationMiddleware(next http.Handler) http.Handler
		ServiceAdminMiddleware(next http.Handler) http.Handler

		AuthenticateUser(ctx context.Context, loginData *UserLoginInput) (*User, *http.Cookie, error)
		LogoutUser(ctx context.Context, sessionCtxData *SessionContextData, req *http.Request, res http.ResponseWriter) error
	}

	// AuthAuditManager describes a structure capable of auditing auth events.
	AuthAuditManager interface {
		LogCycleCookieSecretEvent(ctx context.Context, userID uint64)
		LogSuccessfulLoginEvent(ctx context.Context, userID uint64)
		LogBannedUserLoginAttemptEvent(ctx context.Context, userID uint64)
		LogUnsuccessfulLoginBadPasswordEvent(ctx context.Context, userID uint64)
		LogUnsuccessfulLoginBad2FATokenEvent(ctx context.Context, userID uint64)
		LogLogoutEvent(ctx context.Context, userID uint64)
	}
)

var _ validation.ValidatableWithContext = (*ChangeActiveHouseholdInput)(nil)

// ValidateWithContext validates a ChangeActiveHouseholdInput.
func (x *ChangeActiveHouseholdInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.HouseholdID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*PASETOCreationInput)(nil)

// ValidateWithContext ensures our  provided UserLoginInput meets expectations.
func (i *PASETOCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.ClientID, validation.Required),
		validation.Field(&i.RequestTime, validation.Required),
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

// ToBytes returns the gob encoded session info.
func (x *SessionContextData) ToBytes() []byte {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(x); err != nil {
		panic(err)
	}

	return b.Bytes()
}

// AttachToLogger provides a consistent way to attach a SessionContextData object to a logger.
func (x *SessionContextData) AttachToLogger(logger logging.Logger) logging.Logger {
	if x != nil {
		logger = logger.WithValue(keys.RequesterIDKey, x.Requester.UserID).
			WithValue(keys.ActiveHouseholdIDKey, x.ActiveHouseholdID)

		if x.Requester.ServicePermissions != nil {
			logger = logger.WithValue(keys.ServiceRoleKey, x.Requester.ServicePermissions.IsServiceAdmin())
		}
	}

	return logger
}
