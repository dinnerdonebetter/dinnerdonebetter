package grpc

import (
	"database/sql"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	identitymanager "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	"github.com/dinnerdonebetter/backend/internal/platform/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	platformerrors "github.com/dinnerdonebetter/backend/internal/platform/errors"

	"google.golang.org/grpc/codes"
)

// PlatformMapper maps platform-level errors to gRPC codes.
// It does not depend on any domain (mealplanning, etc.).
var PlatformMapper GRPCErrorMapper = platformMapper{}

type platformMapper struct{}

func (platformMapper) Map(err error) (code codes.Code, ok bool) {
	if err == nil {
		return codes.OK, false
	}
	switch {
	case errors.Is(err, database.ErrUserAlreadyExists):
		return codes.AlreadyExists, true
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound, true
	case errors.Is(err, authentication.ErrInvalidTOTPToken),
		errors.Is(err, authentication.ErrPasswordDoesNotMatch):
		return codes.InvalidArgument, true
	case errors.Is(err, sessions.ErrAuthenticationNotFound),
		errors.Is(err, sessions.ErrNoSessionContextDataAvailable):
		return codes.Unauthenticated, true
	case errors.Is(err, circuitbreaking.ErrCircuitBroken):
		return codes.Unavailable, true
	case errors.Is(err, platformerrors.ErrNilInputParameter),
		errors.Is(err, platformerrors.ErrEmptyInputParameter),
		errors.Is(err, platformerrors.ErrNilInputProvided),
		errors.Is(err, platformerrors.ErrInvalidIDProvided),
		errors.Is(err, platformerrors.ErrEmptyInputProvided),
		errors.Is(err, identitymanager.ErrInvalidIDProvided),
		errors.Is(err, identitymanager.ErrNilInputProvided),
		errors.Is(err, identitymanager.ErrEmptyInputProvided):
		return codes.InvalidArgument, true
	default:
		return codes.Unknown, false
	}
}
