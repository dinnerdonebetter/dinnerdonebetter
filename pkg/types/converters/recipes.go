package converters

import "github.com/prixfixeco/backend/pkg/types"

// ConvertRecipeToRecipeUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeToRecipeUpdateRequestInput(input *types.Recipe) *types.RecipeUpdateRequestInput {
	x := &types.RecipeUpdateRequestInput{
		Name:               &input.Name,
		Source:             &input.Source,
		Description:        &input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		CreatedByUser:      &input.CreatedByUser,
		SealOfApproval:     &input.SealOfApproval,
		YieldsPortions:     &input.YieldsPortions,
	}

	return x
}

// ConvertRecipeCreationRequestInputToRecipeDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeCreationRequestInputToRecipeDatabaseCreationInput(input *types.RecipeCreationRequestInput) (*types.RecipeDatabaseCreationInput, error) {
	x := &types.RecipeDatabaseCreationInput{
		AlsoCreateMeal:     input.AlsoCreateMeal,
		Name:               input.Name,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		SealOfApproval:     input.SealOfApproval,
		YieldsPortions:     input.YieldsPortions,
	}

	for _, step := range input.Steps {
		x.Steps = append(x.Steps, ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(step))
	}

	for _, task := range input.PrepTasks {
		prepTaskDatabaseCreationInput, err := ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput(x, task)
		if err != nil {
			return nil, err
		}
		x.PrepTasks = append(x.PrepTasks, prepTaskDatabaseCreationInput)
	}

	return x, nil
}

// ConvertRecipeToRecipeCreationRequestInputFromRecipe builds a RecipeCreationRequestInput from a recipe.
func ConvertRecipeToRecipeCreationRequestInputFromRecipe(recipe *types.Recipe) *types.RecipeCreationRequestInput {
	steps := []*types.RecipeStepCreationRequestInput{}
	for _, step := range recipe.Steps {
		steps = append(steps, ConvertRecipeStepToRecipeStepCreationRequestInput(step))
	}

	prepTasks := []*types.RecipePrepTaskWithinRecipeCreationRequestInput{}
	for _, prepTask := range recipe.PrepTasks {
		prepTasks = append(prepTasks, ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(recipe, prepTask))
	}

	return &types.RecipeCreationRequestInput{
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		CreatedByUser:      recipe.CreatedByUser,
		SealOfApproval:     recipe.SealOfApproval,
		YieldsPortions:     recipe.YieldsPortions,
		Steps:              steps,
		PrepTasks:          prepTasks,
	}
}

// ConvertRecipeToRecipeDatabaseCreationInput builds a RecipeDatabaseCreationInput from a recipe.
func ConvertRecipeToRecipeDatabaseCreationInput(recipe *types.Recipe) *types.RecipeDatabaseCreationInput {
	steps := []*types.RecipeStepDatabaseCreationInput{}
	for _, step := range recipe.Steps {
		steps = append(steps, ConvertRecipeStepToRecipeStepDatabaseCreationInput(step))
	}

	prepTasks := []*types.RecipePrepTaskDatabaseCreationInput{}
	for _, prepTask := range recipe.PrepTasks {
		prepTasks = append(prepTasks, ConvertRecipePrepTaskToRecipePrepTaskDatabaseCreationInput(prepTask))
	}

	return &types.RecipeDatabaseCreationInput{
		ID:                 recipe.ID,
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		CreatedByUser:      recipe.CreatedByUser,
		SealOfApproval:     recipe.SealOfApproval,
		YieldsPortions:     recipe.YieldsPortions,
		Steps:              steps,
		PrepTasks:          prepTasks,
	}
}
