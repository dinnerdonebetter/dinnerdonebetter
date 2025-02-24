package types

type (
	RecipeManagementDataManager interface {
		RecipeDataManager
		RecipeMediaDataManager
		RecipeStepCompletionConditionDataManager
		RecipeStepIngredientDataManager
		RecipeStepInstrumentDataManager
		RecipeStepProductDataManager
		RecipeStepDataManager
		RecipeStepVesselDataManager
		RecipeRatingDataManager
		RecipePrepTaskDataManager
	}

	RecipeManagementDataService interface {
		RecipeDataService
		RecipeStepCompletionConditionDataService
		RecipeStepIngredientDataService
		RecipeStepInstrumentDataService
		RecipeStepProductDataService
		RecipeStepDataService
		RecipeStepVesselDataService
		RecipeRatingDataService
		RecipePrepTaskDataService
	}
)
