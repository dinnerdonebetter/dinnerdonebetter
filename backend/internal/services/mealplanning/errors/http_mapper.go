package errors

import (
	"errors"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	mealplanningrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"

	"github.com/verygoodsoftwarenotvirus/platform/v4/errors/http"
	"github.com/verygoodsoftwarenotvirus/platform/v4/types"
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
