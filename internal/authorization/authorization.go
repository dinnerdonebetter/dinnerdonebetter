package authorization

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func hasPermission(p Permission, roles ...string) bool {
	if len(roles) == 0 {
		return false
	}

	for _, r := range roles {
		if !globalAuthorizer.IsGranted(r, p, nil) {
			return false
		}
	}

	return true
}

// CanCreateValidInstruments returns whether a user can create valid instruments or not.
func CanCreateValidInstruments(roles ...string) bool {
	return hasPermission(CreateValidInstrumentsPermission, roles...)
}

// CanSeeValidInstruments returns whether a user can view valid instruments or not.
func CanSeeValidInstruments(roles ...string) bool {
	return hasPermission(ReadValidInstrumentsPermission, roles...)
}

// CanSearchValidInstruments returns whether a user can search valid instruments or not.
func CanSearchValidInstruments(roles ...string) bool {
	return hasPermission(SearchValidInstrumentsPermission, roles...)
}

// CanUpdateValidInstruments returns whether a user can update valid instruments or not.
func CanUpdateValidInstruments(roles ...string) bool {
	return hasPermission(UpdateValidInstrumentsPermission, roles...)
}

// CanDeleteValidInstruments returns whether a user can delete valid instruments or not.
func CanDeleteValidInstruments(roles ...string) bool {
	return hasPermission(ArchiveValidInstrumentsPermission, roles...)
}

// CanCreateValidIngredients returns whether a user can create valid ingredients or not.
func CanCreateValidIngredients(roles ...string) bool {
	return hasPermission(CreateValidIngredientsPermission, roles...)
}

// CanSeeValidIngredients returns whether a user can view valid ingredients or not.
func CanSeeValidIngredients(roles ...string) bool {
	return hasPermission(ReadValidIngredientsPermission, roles...)
}

// CanSearchValidIngredients returns whether a user can search valid ingredients or not.
func CanSearchValidIngredients(roles ...string) bool {
	return hasPermission(SearchValidIngredientsPermission, roles...)
}

// CanUpdateValidIngredients returns whether a user can update valid ingredients or not.
func CanUpdateValidIngredients(roles ...string) bool {
	return hasPermission(UpdateValidIngredientsPermission, roles...)
}

// CanDeleteValidIngredients returns whether a user can delete valid ingredients or not.
func CanDeleteValidIngredients(roles ...string) bool {
	return hasPermission(ArchiveValidIngredientsPermission, roles...)
}

// CanCreateValidPreparations returns whether a user can create valid preparations or not.
func CanCreateValidPreparations(roles ...string) bool {
	return hasPermission(CreateValidPreparationsPermission, roles...)
}

// CanSeeValidPreparations returns whether a user can view valid preparations or not.
func CanSeeValidPreparations(roles ...string) bool {
	return hasPermission(ReadValidPreparationsPermission, roles...)
}

// CanSearchValidPreparations returns whether a user can search valid preparations or not.
func CanSearchValidPreparations(roles ...string) bool {
	return hasPermission(SearchValidPreparationsPermission, roles...)
}

// CanUpdateValidPreparations returns whether a user can update valid preparations or not.
func CanUpdateValidPreparations(roles ...string) bool {
	return hasPermission(UpdateValidPreparationsPermission, roles...)
}

// CanDeleteValidPreparations returns whether a user can delete valid preparations or not.
func CanDeleteValidPreparations(roles ...string) bool {
	return hasPermission(ArchiveValidPreparationsPermission, roles...)
}

// CanCreateValidIngredientPreparations returns whether a user can create valid ingredient preparations or not.
func CanCreateValidIngredientPreparations(roles ...string) bool {
	return hasPermission(CreateValidIngredientPreparationsPermission, roles...)
}

// CanSeeValidIngredientPreparations returns whether a user can view valid ingredient preparations or not.
func CanSeeValidIngredientPreparations(roles ...string) bool {
	return hasPermission(ReadValidIngredientPreparationsPermission, roles...)
}

// CanUpdateValidIngredientPreparations returns whether a user can update valid ingredient preparations or not.
func CanUpdateValidIngredientPreparations(roles ...string) bool {
	return hasPermission(UpdateValidIngredientPreparationsPermission, roles...)
}

// CanDeleteValidIngredientPreparations returns whether a user can delete valid ingredient preparations or not.
func CanDeleteValidIngredientPreparations(roles ...string) bool {
	return hasPermission(ArchiveValidIngredientPreparationsPermission, roles...)
}

// CanCreateRecipes returns whether a user can create recipes or not.
func CanCreateRecipes(roles ...string) bool {
	return hasPermission(CreateRecipesPermission, roles...)
}

// CanSeeRecipes returns whether a user can view recipes or not.
func CanSeeRecipes(roles ...string) bool {
	return hasPermission(ReadRecipesPermission, roles...)
}

// CanUpdateRecipes returns whether a user can update recipes or not.
func CanUpdateRecipes(roles ...string) bool {
	return hasPermission(UpdateRecipesPermission, roles...)
}

// CanDeleteRecipes returns whether a user can delete recipes or not.
func CanDeleteRecipes(roles ...string) bool {
	return hasPermission(ArchiveRecipesPermission, roles...)
}

// CanCreateRecipeSteps returns whether a user can create recipe steps or not.
func CanCreateRecipeSteps(roles ...string) bool {
	return hasPermission(CreateRecipeStepsPermission, roles...)
}

// CanSeeRecipeSteps returns whether a user can view recipe steps or not.
func CanSeeRecipeSteps(roles ...string) bool {
	return hasPermission(ReadRecipeStepsPermission, roles...)
}

// CanUpdateRecipeSteps returns whether a user can update recipe steps or not.
func CanUpdateRecipeSteps(roles ...string) bool {
	return hasPermission(UpdateRecipeStepsPermission, roles...)
}

// CanDeleteRecipeSteps returns whether a user can delete recipe steps or not.
func CanDeleteRecipeSteps(roles ...string) bool {
	return hasPermission(ArchiveRecipeStepsPermission, roles...)
}

// CanCreateRecipeStepInstruments returns whether a user can create recipe step instruments or not.
func CanCreateRecipeStepInstruments(roles ...string) bool {
	return hasPermission(CreateRecipeStepInstrumentsPermission, roles...)
}

// CanSeeRecipeStepInstruments returns whether a user can view recipe step instruments or not.
func CanSeeRecipeStepInstruments(roles ...string) bool {
	return hasPermission(ReadRecipeStepInstrumentsPermission, roles...)
}

// CanUpdateRecipeStepInstruments returns whether a user can update recipe step instruments or not.
func CanUpdateRecipeStepInstruments(roles ...string) bool {
	return hasPermission(UpdateRecipeStepInstrumentsPermission, roles...)
}

// CanDeleteRecipeStepInstruments returns whether a user can delete recipe step instruments or not.
func CanDeleteRecipeStepInstruments(roles ...string) bool {
	return hasPermission(ArchiveRecipeStepInstrumentsPermission, roles...)
}

// CanCreateRecipeStepIngredients returns whether a user can create recipe step ingredients or not.
func CanCreateRecipeStepIngredients(roles ...string) bool {
	return hasPermission(CreateRecipeStepIngredientsPermission, roles...)
}

// CanSeeRecipeStepIngredients returns whether a user can view recipe step ingredients or not.
func CanSeeRecipeStepIngredients(roles ...string) bool {
	return hasPermission(ReadRecipeStepIngredientsPermission, roles...)
}

// CanUpdateRecipeStepIngredients returns whether a user can update recipe step ingredients or not.
func CanUpdateRecipeStepIngredients(roles ...string) bool {
	return hasPermission(UpdateRecipeStepIngredientsPermission, roles...)
}

// CanDeleteRecipeStepIngredients returns whether a user can delete recipe step ingredients or not.
func CanDeleteRecipeStepIngredients(roles ...string) bool {
	return hasPermission(ArchiveRecipeStepIngredientsPermission, roles...)
}

// CanCreateRecipeStepProducts returns whether a user can create recipe step products or not.
func CanCreateRecipeStepProducts(roles ...string) bool {
	return hasPermission(CreateRecipeStepProductsPermission, roles...)
}

// CanSeeRecipeStepProducts returns whether a user can view recipe step products or not.
func CanSeeRecipeStepProducts(roles ...string) bool {
	return hasPermission(ReadRecipeStepProductsPermission, roles...)
}

// CanUpdateRecipeStepProducts returns whether a user can update recipe step products or not.
func CanUpdateRecipeStepProducts(roles ...string) bool {
	return hasPermission(UpdateRecipeStepProductsPermission, roles...)
}

// CanDeleteRecipeStepProducts returns whether a user can delete recipe step products or not.
func CanDeleteRecipeStepProducts(roles ...string) bool {
	return hasPermission(ArchiveRecipeStepProductsPermission, roles...)
}

// CanCreateMealPlans returns whether a user can create meal plans or not.
func CanCreateMealPlans(roles ...string) bool {
	return hasPermission(CreateMealPlansPermission, roles...)
}

// CanSeeMealPlans returns whether a user can view meal plans or not.
func CanSeeMealPlans(roles ...string) bool {
	return hasPermission(ReadMealPlansPermission, roles...)
}

// CanUpdateMealPlans returns whether a user can update meal plans or not.
func CanUpdateMealPlans(roles ...string) bool {
	return hasPermission(UpdateMealPlansPermission, roles...)
}

// CanDeleteMealPlans returns whether a user can delete meal plans or not.
func CanDeleteMealPlans(roles ...string) bool {
	return hasPermission(ArchiveMealPlansPermission, roles...)
}

// CanCreateMealPlanOptions returns whether a user can create meal plan options or not.
func CanCreateMealPlanOptions(roles ...string) bool {
	return hasPermission(CreateMealPlanOptionsPermission, roles...)
}

// CanSeeMealPlanOptions returns whether a user can view meal plan options or not.
func CanSeeMealPlanOptions(roles ...string) bool {
	return hasPermission(ReadMealPlanOptionsPermission, roles...)
}

// CanUpdateMealPlanOptions returns whether a user can update meal plan options or not.
func CanUpdateMealPlanOptions(roles ...string) bool {
	return hasPermission(UpdateMealPlanOptionsPermission, roles...)
}

// CanDeleteMealPlanOptions returns whether a user can delete meal plan options or not.
func CanDeleteMealPlanOptions(roles ...string) bool {
	return hasPermission(ArchiveMealPlanOptionsPermission, roles...)
}

// CanCreateMealPlanOptionVotes returns whether a user can create meal plan option votes or not.
func CanCreateMealPlanOptionVotes(roles ...string) bool {
	return hasPermission(CreateMealPlanOptionVotesPermission, roles...)
}

// CanSeeMealPlanOptionVotes returns whether a user can view meal plan option votes or not.
func CanSeeMealPlanOptionVotes(roles ...string) bool {
	return hasPermission(ReadMealPlanOptionVotesPermission, roles...)
}

// CanUpdateMealPlanOptionVotes returns whether a user can update meal plan option votes or not.
func CanUpdateMealPlanOptionVotes(roles ...string) bool {
	return hasPermission(UpdateMealPlanOptionVotesPermission, roles...)
}

// CanDeleteMealPlanOptionVotes returns whether a user can delete meal plan option votes or not.
func CanDeleteMealPlanOptionVotes(roles ...string) bool {
	return hasPermission(ArchiveMealPlanOptionVotesPermission, roles...)
}
