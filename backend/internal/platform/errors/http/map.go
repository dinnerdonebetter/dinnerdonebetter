package http

import (
	"database/sql"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	identitymanager "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	mp "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	platformerrors "github.com/dinnerdonebetter/backend/internal/platform/errors"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	mealplanningrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
)

// ToAPIError maps known sentinel errors to types.ErrorCode and a safe user-facing message.
// Returns (code, message). Use types.ErrTalkingToDatabase and "an error occurred" as fallback for unknown errors.
func ToAPIError(err error) (code types.ErrorCode, msg string) {
	if err == nil {
		return types.ErrNothingSpecific, ""
	}
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return types.ErrDataNotFound, "data not found"
	case errors.Is(err, database.ErrUserAlreadyExists):
		return types.ErrValidatingRequestInput, "user already exists"
	case errors.Is(err, authentication.ErrInvalidTOTPToken),
		errors.Is(err, authentication.ErrPasswordDoesNotMatch):
		return types.ErrValidatingRequestInput, "invalid credentials"
	case errors.Is(err, sessions.ErrAuthenticationNotFound),
		errors.Is(err, sessions.ErrNoSessionContextDataAvailable):
		return types.ErrFetchingSessionContextData, "session not found"
	case errors.Is(err, circuitbreaking.ErrCircuitBroken):
		return types.ErrCircuitBroken, "service temporarily unavailable"
	case errors.Is(err, platformerrors.ErrNilInputParameter),
		errors.Is(err, platformerrors.ErrEmptyInputParameter),
		errors.Is(err, platformerrors.ErrNilInputProvided),
		errors.Is(err, platformerrors.ErrInvalidIDProvided),
		errors.Is(err, platformerrors.ErrEmptyInputProvided),
		errors.Is(err, identitymanager.ErrInvalidIDProvided),
		errors.Is(err, identitymanager.ErrNilInputProvided),
		errors.Is(err, identitymanager.ErrEmptyInputProvided):
		return types.ErrValidatingRequestInput, "invalid input"
	case errors.Is(err, mp.ErrDuplicateMeal),
		errors.Is(err, mp.ErrDuplicateMealInList),
		errors.Is(err, mp.ErrDuplicateMealPlanOption):
		return types.ErrValidatingRequestInput, "duplicate entry"
	case errors.Is(err, mealplanningrepo.ErrAlreadyFinalized):
		return types.ErrValidatingRequestInput, "already finalized"
	default:
		code = types.ErrTalkingToDatabase
		msg = "an error occurred"
		return code, msg
	}
}
