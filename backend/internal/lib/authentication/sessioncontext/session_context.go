package sessioncontext

import (
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
)

func init() {
	gob.Register(&SessionContextData{})
}

var (
	// ErrNoSessionContextDataAvailable indicates no SessionContextData was attached to the request.
	ErrNoSessionContextDataAvailable = errors.New("no SessionContextData attached to session context data")
)

// FetchContextFromRequest fetches a SessionContextData from a request.
func FetchContextFromRequest(req *http.Request) (*SessionContextData, error) {
	if sessionCtxData, ok := req.Context().Value(SessionContextDataKey).(*SessionContextData); ok && sessionCtxData != nil {
		return sessionCtxData, nil
	}

	return nil, ErrNoSessionContextDataAvailable
}

const SessionContextDataKey routing.ContextKey = "session_context_data"

// SessionContextData represents what we encode in our passwords cookies.
type SessionContextData struct {
	_ struct{} `json:"-"`

	HouseholdPermissions map[string]authorization.HouseholdRolePermissionsChecker `json:"-"`
	Requester            RequesterInfo                                            `json:"-"`
	ActiveHouseholdID    string                                                   `json:"-"`
}

// RequesterInfo contains data relevant to the user making a request.
type RequesterInfo struct {
	_ struct{} `json:"-"`

	ServicePermissions       authorization.ServiceRolePermissionChecker `json:"-"`
	AccountStatus            string                                     `json:"-"`
	AccountStatusExplanation string                                     `json:"-"`
	UserID                   string                                     `json:"-"`
	EmailAddress             string                                     `json:"-"`
	Username                 string                                     `json:"-"`
}

// GetUserID is a simple getter.
func (x *SessionContextData) GetUserID() string {
	return x.Requester.UserID
}

// GetServicePermissions is a simple getter.
func (x *SessionContextData) GetServicePermissions() authorization.ServiceRolePermissionChecker {
	return x.Requester.ServicePermissions
}

// GetActiveHouseholdID is a simple getter.
func (x *SessionContextData) GetActiveHouseholdID() string {
	return x.ActiveHouseholdID
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
