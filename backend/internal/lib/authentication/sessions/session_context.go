package sessions

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
	gob.Register(&ContextData{})
}

var (
	// ErrNoSessionContextDataAvailable indicates no ContextData was attached to the request.
	ErrNoSessionContextDataAvailable = errors.New("no ContextData attached to session context data")
)

// FetchContextFromRequest fetches a ContextData from a request.
func FetchContextFromRequest(req *http.Request) (*ContextData, error) {
	if sessionCtxData, ok := req.Context().Value(SessionContextDataKey).(*ContextData); ok && sessionCtxData != nil {
		return sessionCtxData, nil
	}

	return nil, ErrNoSessionContextDataAvailable
}

const SessionContextDataKey routing.ContextKey = "session_context_data"

// ContextData represents what we encode in our passwords cookies.
type ContextData struct {
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
func (x *ContextData) GetUserID() string {
	return x.Requester.UserID
}

// GetServicePermissions is a simple getter.
func (x *ContextData) GetServicePermissions() authorization.ServiceRolePermissionChecker {
	return x.Requester.ServicePermissions
}

// GetActiveHouseholdID is a simple getter.
func (x *ContextData) GetActiveHouseholdID() string {
	return x.ActiveHouseholdID
}

// HouseholdRolePermissionsChecker returns the relevant HouseholdRolePermissionsChecker.
func (x *ContextData) HouseholdRolePermissionsChecker() authorization.HouseholdRolePermissionsChecker {
	return x.HouseholdPermissions[x.ActiveHouseholdID]
}

// ServiceRolePermissionChecker returns the relevant ServiceRolePermissionChecker.
func (x *ContextData) ServiceRolePermissionChecker() authorization.ServiceRolePermissionChecker {
	return x.Requester.ServicePermissions
}

// AttachToLogger provides a consistent way to attach a ContextData object to a logger.
func (x *ContextData) AttachToLogger(logger logging.Logger) logging.Logger {
	if x != nil {
		logger = logger.WithValue(keys.RequesterIDKey, x.Requester.UserID).
			WithValue(keys.ActiveHouseholdIDKey, x.ActiveHouseholdID)
	}

	return logger
}
