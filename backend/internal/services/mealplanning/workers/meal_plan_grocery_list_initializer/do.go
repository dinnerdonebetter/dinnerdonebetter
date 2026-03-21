package mealplangrocerylistinitializer

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/grocerylistpreparation"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

// RegisterMealPlanGroceryListInitializer registers the meal plan grocery list initializer with the injector.
func RegisterMealPlanGroceryListInitializer(i do.Injector) {
	do.Provide[*Worker](i, func(i do.Injector) (*Worker, error) {
		return NewMealPlanGroceryListInitializer(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[mealplanning.Repository](i),
			do.MustInvoke[grocerylistpreparation.GroceryListCreator](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
		)
	})
}
