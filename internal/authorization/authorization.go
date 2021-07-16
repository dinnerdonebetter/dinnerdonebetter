package authorization

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func hasPermission(p Permission, roles ...string) bool {
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

// CanCreateValidPreparationInstruments returns whether a user can create valid preparation instruments or not.
func CanCreateValidPreparationInstruments(roles ...string) bool {
	return hasPermission(CreateValidPreparationInstrumentsPermission, roles...)
}

// CanSeeValidPreparationInstruments returns whether a user can view valid preparation instruments or not.
func CanSeeValidPreparationInstruments(roles ...string) bool {
	return hasPermission(ReadValidPreparationInstrumentsPermission, roles...)
}

// CanUpdateValidPreparationInstruments returns whether a user can update valid preparation instruments or not.
func CanUpdateValidPreparationInstruments(roles ...string) bool {
	return hasPermission(UpdateValidPreparationInstrumentsPermission, roles...)
}

// CanDeleteValidPreparationInstruments returns whether a user can delete valid preparation instruments or not.
func CanDeleteValidPreparationInstruments(roles ...string) bool {
	return hasPermission(ArchiveValidPreparationInstrumentsPermission, roles...)
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

// CanCreateInvitations returns whether a user can create invitations or not.
func CanCreateInvitations(roles ...string) bool {
	return hasPermission(CreateInvitationsPermission, roles...)
}

// CanSeeInvitations returns whether a user can view invitations or not.
func CanSeeInvitations(roles ...string) bool {
	return hasPermission(ReadInvitationsPermission, roles...)
}

// CanUpdateInvitations returns whether a user can update invitations or not.
func CanUpdateInvitations(roles ...string) bool {
	return hasPermission(UpdateInvitationsPermission, roles...)
}

// CanDeleteInvitations returns whether a user can delete invitations or not.
func CanDeleteInvitations(roles ...string) bool {
	return hasPermission(ArchiveInvitationsPermission, roles...)
}

// CanCreateReports returns whether a user can create reports or not.
func CanCreateReports(roles ...string) bool {
	return hasPermission(CreateReportsPermission, roles...)
}

// CanSeeReports returns whether a user can view reports or not.
func CanSeeReports(roles ...string) bool {
	return hasPermission(ReadReportsPermission, roles...)
}

// CanUpdateReports returns whether a user can update reports or not.
func CanUpdateReports(roles ...string) bool {
	return hasPermission(UpdateReportsPermission, roles...)
}

// CanDeleteReports returns whether a user can delete reports or not.
func CanDeleteReports(roles ...string) bool {
	return hasPermission(ArchiveReportsPermission, roles...)
}
