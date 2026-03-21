package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

// RegisterMealPlanFinalizer registers the meal plan finalizer with the injector.
func RegisterMealPlanFinalizer(i do.Injector) {
	do.Provide[*Worker](i, func(i do.Injector) (*Worker, error) {
		return NewMealPlanFinalizer(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[mealplanning.Repository](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
		)
	})
}
