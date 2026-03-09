package errors

import (
	"errors"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	identitymanager "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	"github.com/dinnerdonebetter/backend/internal/platform/errors/http"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func init() {
	http.RegisterHTTPErrorMapper(authSessionIdentityHTTPMapper{})
}

type authSessionIdentityHTTPMapper struct{}

func (authSessionIdentityHTTPMapper) Map(err error) (code types.ErrorCode, msg string, ok bool) {
	if err == nil {
		return "", "", false
	}
	switch {
	case errors.Is(err, authentication.ErrInvalidTOTPToken),
		errors.Is(err, authentication.ErrPasswordDoesNotMatch):
		return types.ErrValidatingRequestInput, "invalid credentials", true
	case errors.Is(err, sessions.ErrAuthenticationNotFound),
		errors.Is(err, sessions.ErrNoSessionContextDataAvailable):
		return types.ErrFetchingSessionContextData, "session not found", true
	case errors.Is(err, identitymanager.ErrInvalidIDProvided),
		errors.Is(err, identitymanager.ErrNilInputProvided),
		errors.Is(err, identitymanager.ErrEmptyInputProvided):
		return types.ErrValidatingRequestInput, "invalid input", true
	default:
		return "", "", false
	}
}
