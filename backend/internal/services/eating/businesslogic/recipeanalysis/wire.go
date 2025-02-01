package recipeanalysis

import (
	"github.com/google/wire"
)

// ProvidersRecipeAnalysis are our collection of what we provide to other services.
var ProvidersRecipeAnalysis = wire.NewSet(
	NewRecipeAnalyzer,
)
