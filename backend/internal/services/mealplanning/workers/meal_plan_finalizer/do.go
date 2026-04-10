package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"

	"github.com/verygoodsoftwarenotvirus/platform/v5/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v5/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/tracing"

	"github.com/samber/do/v2"
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
