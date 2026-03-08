package errors

import (
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/errors/grpc"
	mealplanningrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"

	"google.golang.org/grpc/codes"
)

func init() {
	grpc.RegisterGRPCErrorMapper(mealPlanningGRPCMapper{})
}

type mealPlanningGRPCMapper struct{}

func (mealPlanningGRPCMapper) Map(err error) (code codes.Code, ok bool) {
	if err == nil {
		return codes.Unknown, false
	}
	switch {
	case errors.Is(err, mealplanning.ErrDuplicateMeal),
		errors.Is(err, mealplanning.ErrDuplicateMealInList),
		errors.Is(err, mealplanning.ErrDuplicateMealPlanOption):
		return codes.AlreadyExists, true
	case errors.Is(err, mealplanningrepo.ErrAlreadyFinalized):
		return codes.FailedPrecondition, true
	default:
		return codes.Unknown, false
	}
}
