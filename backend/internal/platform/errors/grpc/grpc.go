package grpc

import (
	"database/sql"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	mp "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	platformerrors "github.com/dinnerdonebetter/backend/internal/platform/errors"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"google.golang.org/grpc/codes"
)

// PrepareAndLogGRPCStatus derives the gRPC code via MapToGRPC, then logs, traces, and returns
// a status error. Use defaultCode as the fallback for unknown errors.
func PrepareAndLogGRPCStatus(err error, logger logging.Logger, span tracing.Span, defaultCode codes.Code, descriptionFmt string, descriptionArgs ...any) error {
	code := MapToGRPC(err, defaultCode)
	return observability.PrepareAndLogGRPCStatus(err, logger, span, code, descriptionFmt, descriptionArgs...)
}

// MapToGRPC returns the appropriate gRPC code for known sentinel errors.
// Use std errors.Is for matching. Returns defaultCode if no match.
func MapToGRPC(err error, defaultCode codes.Code) codes.Code {
	if err == nil {
		return codes.OK
	}
	switch {
	case errors.Is(err, mp.ErrDuplicateMeal),
		errors.Is(err, mp.ErrDuplicateMealInList),
		errors.Is(err, mp.ErrDuplicateMealPlanOption),
		errors.Is(err, database.ErrUserAlreadyExists):
		return codes.AlreadyExists
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case errors.Is(err, authentication.ErrInvalidTOTPToken),
		errors.Is(err, authentication.ErrPasswordDoesNotMatch):
		return codes.InvalidArgument
	case errors.Is(err, circuitbreaking.ErrCircuitBroken):
		return codes.Unavailable
	case errors.Is(err, platformerrors.ErrNilInputParameter),
		errors.Is(err, platformerrors.ErrEmptyInputParameter),
		errors.Is(err, platformerrors.ErrNilInputProvided),
		errors.Is(err, platformerrors.ErrInvalidIDProvided),
		errors.Is(err, platformerrors.ErrEmptyInputProvided):
		return codes.InvalidArgument
	}
	return defaultCode
}
