package authentication

import (
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	// ErrNoSessionContextDataAvailable indicates no SessionContextData was attached to the request.
	ErrNoSessionContextDataAvailable = errors.New("no SessionContextData attached to session context data")
)

// FetchContextFromRequest fetches a SessionContextData from a request.
func FetchContextFromRequest(req *http.Request) (*types.SessionContextData, error) {
	if sessionCtxData, ok := req.Context().Value(types.SessionContextDataKey).(*types.SessionContextData); ok && sessionCtxData != nil {
		return sessionCtxData, nil
	}

	return nil, ErrNoSessionContextDataAvailable
}
