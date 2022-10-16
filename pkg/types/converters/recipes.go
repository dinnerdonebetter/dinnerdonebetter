package converters

import "github.com/prixfixeco/api_server/pkg/types"

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
		prepTaskDatabaseCreationInput, err := types.RecipePrepTaskDatabaseCreationInputFromRecipePrepTaskWithinRecipeCreationInput(x, task)
		if err != nil {
			return nil, err
		}
		x.PrepTasks = append(x.PrepTasks, prepTaskDatabaseCreationInput)
	}

	return x, nil
}
