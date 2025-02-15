package recipemanagement

const (
	// RecipeIDURIParamKey is a standard string that we'll use to refer to recipe IDs with.
	RecipeIDURIParamKey = "recipeID"
	// RecipeStepCompletionConditionIDURIParamKey is a standard string that we'll use to refer to recipe step completion condition IDs with.
	RecipeStepCompletionConditionIDURIParamKey = "recipeStepCompletionConditionID"
	// RecipeStepIDURIParamKey is a standard string that we'll use to refer to recipe step IDs with.
	RecipeStepIDURIParamKey = "recipeStepID"
	// RecipeStepIngredientIDURIParamKey is a standard string that we'll use to refer to recipe step ingredient IDs with.
	RecipeStepIngredientIDURIParamKey = "recipeStepIngredientID"
	// RecipeStepInstrumentIDURIParamKey is a standard string that we'll use to refer to recipe step instrument IDs with.
	RecipeStepInstrumentIDURIParamKey = "recipeStepInstrumentID"
	// RecipeStepProductIDURIParamKey is a standard string that we'll use to refer to recipe step product IDs with.
	RecipeStepProductIDURIParamKey = "recipeStepProductID"
	// RecipeStepVesselIDURIParamKey is a standard string that we'll use to refer to recipe step vessel IDs with.
	RecipeStepVesselIDURIParamKey = "recipeStepVesselID"
	// RecipePrepTaskIDURIParamKey is a standard string that we'll use to refer to a recipe prep task's ID.
	RecipePrepTaskIDURIParamKey = "recipePrepTaskID"
	// RecipeRatingIDURIParamKey is a standard string that we'll use to refer to a recipe rating's ID.
	RecipeRatingIDURIParamKey = "recipeRatingID"
)

var (
	allURIKeys = []string{
		RecipeIDURIParamKey,
		RecipeStepCompletionConditionIDURIParamKey,
		RecipeStepIDURIParamKey,
		RecipeStepIngredientIDURIParamKey,
		RecipeStepInstrumentIDURIParamKey,
		RecipeStepProductIDURIParamKey,
		RecipeStepVesselIDURIParamKey,
		RecipePrepTaskIDURIParamKey,
		RecipeRatingIDURIParamKey,
	}
)
