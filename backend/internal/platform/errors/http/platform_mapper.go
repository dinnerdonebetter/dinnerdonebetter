package http

import (
	"database/sql"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	identitymanager "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	"github.com/dinnerdonebetter/backend/internal/platform/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	platformerrors "github.com/dinnerdonebetter/backend/internal/platform/errors"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// PlatformMapper maps platform-level errors to HTTP error codes and messages.
// It does not depend on any domain (mealplanning, etc.).
var PlatformMapper HTTPErrorMapper = platformMapper{}

type platformMapper struct{}

func (platformMapper) Map(err error) (code types.ErrorCode, msg string, ok bool) {
	if err == nil {
		return types.ErrNothingSpecific, "", false
	}
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return types.ErrDataNotFound, "data not found", true
	case errors.Is(err, database.ErrUserAlreadyExists):
		return types.ErrValidatingRequestInput, "user already exists", true
	case errors.Is(err, authentication.ErrInvalidTOTPToken),
		errors.Is(err, authentication.ErrPasswordDoesNotMatch):
		return types.ErrValidatingRequestInput, "invalid credentials", true
	case errors.Is(err, sessions.ErrAuthenticationNotFound),
		errors.Is(err, sessions.ErrNoSessionContextDataAvailable):
		return types.ErrFetchingSessionContextData, "session not found", true
	case errors.Is(err, circuitbreaking.ErrCircuitBroken):
		return types.ErrCircuitBroken, "service temporarily unavailable", true
	case errors.Is(err, platformerrors.ErrNilInputParameter),
		errors.Is(err, platformerrors.ErrEmptyInputParameter),
		errors.Is(err, platformerrors.ErrNilInputProvided),
		errors.Is(err, platformerrors.ErrInvalidIDProvided),
		errors.Is(err, platformerrors.ErrEmptyInputProvided),
		errors.Is(err, identitymanager.ErrInvalidIDProvided),
		errors.Is(err, identitymanager.ErrNilInputProvided),
		errors.Is(err, identitymanager.ErrEmptyInputProvided):
		return types.ErrValidatingRequestInput, "invalid input", true
	default:
		return "", "", false
	}
}
