package mealplantaskcreator

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v2/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v2/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing"
)

// RegisterMealPlanTaskCreator registers the meal plan task creator with the injector.
func RegisterMealPlanTaskCreator(i do.Injector) {
	do.Provide[*Worker](i, func(i do.Injector) (*Worker, error) {
		return NewMealPlanTaskCreator(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[recipeanalysis.RecipeAnalyzer](i),
			do.MustInvoke[mealplanning.Repository](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
		)
	})
}
