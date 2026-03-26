package recipeanalysis

import (
	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

// RegisterRecipeAnalyzer registers the recipe analyzer with the injector.
func RegisterRecipeAnalyzer(i do.Injector) {
	do.Provide[RecipeAnalyzer](i, func(i do.Injector) (RecipeAnalyzer, error) {
		return NewRecipeAnalyzer(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
		), nil
	})
}
