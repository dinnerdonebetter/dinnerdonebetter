package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ConvertRecipeToRecipeUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeToRecipeUpdateRequestInput(input *mealplanning.Recipe) *mealplanning.RecipeUpdateRequestInput {
	x := &mealplanning.RecipeUpdateRequestInput{
		Name:                &input.Name,
		Slug:                &input.Slug,
		Source:              &input.Source,
		Description:         &input.Description,
		InspiredByRecipeID:  input.InspiredByRecipeID,
		YieldsComponentType: &input.YieldsComponentType,
		SealOfApproval:      &input.SealOfApproval,
		EstimatedPortions: types.Float32RangeWithOptionalMaxUpdateRequestInput{
			Max: input.EstimatedPortions.Max,
			Min: &input.EstimatedPortions.Min,
		},
		PortionName:       &input.PortionName,
		PluralPortionName: &input.PluralPortionName,
		EligibleForMeals:  &input.EligibleForMeals,
	}

	return x
}

// ConvertRecipeCreationRequestInputToRecipeDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeCreationRequestInputToRecipeDatabaseCreationInput(input *mealplanning.RecipeCreationRequestInput) (*mealplanning.RecipeDatabaseCreationInput, error) {
	x := &mealplanning.RecipeDatabaseCreationInput{
		ID:                 identifiers.New(),
		AlsoCreateMeal:     input.AlsoCreateMeal,
		Name:               input.Name,
		Slug:               input.Slug,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		SealOfApproval:     input.SealOfApproval,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Max: input.EstimatedPortions.Max,
			Min: input.EstimatedPortions.Min,
		},
		PortionName:         input.PortionName,
		PluralPortionName:   input.PluralPortionName,
		EligibleForMeals:    input.EligibleForMeals,
		YieldsComponentType: input.YieldsComponentType,
	}

	for _, step := range input.Steps {
		s := ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(step)
		s.BelongsToRecipe = x.ID
		x.Steps = append(x.Steps, s)
	}

	for _, task := range input.PrepTasks {
		prepTaskDatabaseCreationInput, err := ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput(x, task)
		if err != nil {
			return nil, err
		}
		prepTaskDatabaseCreationInput.BelongsToRecipe = x.ID
		x.PrepTasks = append(x.PrepTasks, prepTaskDatabaseCreationInput)
	}

	for _, m := range input.Media {
		x.Media = append(x.Media, ConvertRecipeMediaCreationRequestInputToRecipeMediaDatabaseCreationInput(m))
	}

	return x, nil
}

// ConvertRecipeToRecipeCreationRequestInput builds a RecipeCreationRequestInput from a recipe.
func ConvertRecipeToRecipeCreationRequestInput(input *mealplanning.Recipe) *mealplanning.RecipeCreationRequestInput {
	steps := []*mealplanning.RecipeStepCreationRequestInput{}
	for _, step := range input.Steps {
		steps = append(steps, ConvertRecipeStepToRecipeStepCreationRequestInput(step))
	}

	prepTasks := []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{}
	for _, prepTask := range input.PrepTasks {
		prepTasks = append(prepTasks, ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(input, prepTask))
	}

	media := []*mealplanning.RecipeMediaCreationRequestInput{}
	for _, m := range input.Media {
		media = append(media, ConvertRecipeMediaToRecipeMediaCreationRequestInput(m))
	}

	return &mealplanning.RecipeCreationRequestInput{
		Name:               input.Name,
		Slug:               input.Slug,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		SealOfApproval:     input.SealOfApproval,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Max: input.EstimatedPortions.Max,
			Min: input.EstimatedPortions.Min,
		},
		PortionName:         input.PortionName,
		PluralPortionName:   input.PluralPortionName,
		Steps:               steps,
		PrepTasks:           prepTasks,
		EligibleForMeals:    input.EligibleForMeals,
		YieldsComponentType: input.YieldsComponentType,
	}
}

// ConvertRecipeToRecipeDatabaseCreationInput builds a RecipeDatabaseCreationInput from a recipe.
func ConvertRecipeToRecipeDatabaseCreationInput(input *mealplanning.Recipe) *mealplanning.RecipeDatabaseCreationInput {
	steps := []*mealplanning.RecipeStepDatabaseCreationInput{}
	for _, step := range input.Steps {
		steps = append(steps, ConvertRecipeStepToRecipeStepDatabaseCreationInput(step))
	}

	prepTasks := []*mealplanning.RecipePrepTaskDatabaseCreationInput{}
	for _, prepTask := range input.PrepTasks {
		prepTasks = append(prepTasks, ConvertRecipePrepTaskToRecipePrepTaskDatabaseCreationInput(prepTask))
	}

	return &mealplanning.RecipeDatabaseCreationInput{
		ID:                 input.ID,
		Name:               input.Name,
		Slug:               input.Slug,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		CreatedByUser:      input.CreatedByUser,
		SealOfApproval:     input.SealOfApproval,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Max: input.EstimatedPortions.Max,
			Min: input.EstimatedPortions.Min,
		},
		PortionName:         input.PortionName,
		PluralPortionName:   input.PluralPortionName,
		Steps:               steps,
		PrepTasks:           prepTasks,
		EligibleForMeals:    input.EligibleForMeals,
		YieldsComponentType: input.YieldsComponentType,
	}
}
