package errors

import (
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/errors/http"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	mealplanningrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
)

func init() {
	http.RegisterHTTPErrorMapper(mealPlanningHTTPMapper{})
}

type mealPlanningHTTPMapper struct{}

func (mealPlanningHTTPMapper) Map(err error) (code types.ErrorCode, msg string, ok bool) {
	if err == nil {
		return "", "", false
	}
	switch {
	case errors.Is(err, mealplanning.ErrDuplicateMeal),
		errors.Is(err, mealplanning.ErrDuplicateMealInList),
		errors.Is(err, mealplanning.ErrDuplicateMealPlanOption):
		return types.ErrValidatingRequestInput, "duplicate entry", true
	case errors.Is(err, mealplanningrepo.ErrAlreadyFinalized):
		return types.ErrValidatingRequestInput, "already finalized", true
	default:
		return "", "", false
	}
}
