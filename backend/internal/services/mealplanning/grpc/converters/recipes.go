package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
)

func ConvertRecipeToGRPCRecipe(input *mealplanning.Recipe) *mealplanningsvc.Recipe {
	recipe := &mealplanningsvc.Recipe{
		EstimatedPortions: &types.Float32RangeWithOptionalMax{
			Max: input.EstimatedPortions.Max,
			Min: input.EstimatedPortions.Min,
		},
		CreatedAt:           grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:          grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		ID:                  input.ID,
		YieldsComponentType: input.YieldsComponentType,
		Description:         input.Description,
		Name:                input.Name,
		PortionName:         input.PortionName,
		CreatedByUser:       input.CreatedByUser,
		Source:              input.Source,
		Slug:                input.Slug,
		PluralPortionName:   input.PluralPortionName,
		SealOfApproval:      input.SealOfApproval,
		EligibleForMeals:    input.EligibleForMeals,
	}

	for _, step := range input.Steps {
		recipe.Steps = append(recipe.Steps, ConvertRecipeStepToGRPCRecipeStep(step))
	}

	for _, media := range input.Media {
		recipe.Media = append(recipe.Media, ConvertRecipeMediaToGRPCRecipeMedia(media))
	}

	for _, rps := range input.PrepTasks {
		recipe.PrepTasks = append(recipe.PrepTasks, ConvertRecipePrepTaskToGRPCRecipePrepTask(rps))
	}

	for _, r := range input.SupportingRecipes {
		recipe.SupportingRecipes = append(recipe.SupportingRecipes, ConvertRecipeToGRPCRecipe(r))
	}

	if input.InspiredByRecipeID != nil {
		recipe.InspiredByRecipeID = *input.InspiredByRecipeID
	}

	return recipe
}

func ConvertRecipeStepToGRPCRecipeStep(input *mealplanning.RecipeStep) *mealplanningsvc.RecipeStep {
	step := &mealplanningsvc.RecipeStep{
		CreatedAt: grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		EstimatedTimeInSeconds: &types.OptionalUint32Range{
			Max: input.EstimatedTimeInSeconds.Max,
			Min: input.EstimatedTimeInSeconds.Min,
		},
		TemperatureInCelsius: &types.OptionalFloat32Range{
			Max: input.TemperatureInCelsius.Max,
			Min: input.TemperatureInCelsius.Min,
		},
		ArchivedAt:              grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt:           grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ExplicitInstructions:    input.ExplicitInstructions,
		BelongsToRecipe:         input.BelongsToRecipe,
		Notes:                   input.Notes,
		ConditionExpression:     input.ConditionExpression,
		ID:                      input.ID,
		Index:                   input.Index,
		Optional:                input.Optional,
		StartTimerAutomatically: input.StartTimerAutomatically,
		Preparation:             ConvertValidPreparationToGRPCValidPreparation(&input.Preparation),
		//Instruments:             input.Instruments,
		//Vessels:                 input.Vessels,
		//CompletionConditions:    input.CompletionConditions,
	}

	for _, media := range input.Media {
		step.Media = append(step.Media, ConvertRecipeMediaToGRPCRecipeMedia(media))
	}

	//for _, instrument := range input.Instruments {
	//	step.Instruments = append(step.Instruments, ConvertRecipeMediaToGRPCRecipeMedia(instrument))
	//}

	//for _, vessel := range input.Vessels {
	//	step.Vessels = append(step.Vessels, ConvertRecipeMediaToGRPCRecipeMedia(vessel))
	//}

	//for _, completionCondition := range input.CompletionConditions {
	//	step.CompletionConditions = append(step.CompletionConditions, ConvertRecipeMediaToGRPCRecipeMedia(completionCondition))
	//}

	for _, ingredient := range input.Ingredients {
		step.Ingredients = append(step.Ingredients, ConvertRecipeStepIngredientToGRPCRecipeStepIngredient(ingredient))
	}

	for _, product := range input.Products {
		step.Products = append(step.Products, ConvertRecipeStepProductToGRPCRecipeStepProduct(product))
	}

	return step
}

func ConvertRecipeStepIngredientToGRPCRecipeStepIngredient(input *mealplanning.RecipeStepIngredient) *mealplanningsvc.RecipeStepIngredient {
	ingredient := &mealplanningsvc.RecipeStepIngredient{
		MeasurementUnit: ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(&input.MeasurementUnit),
		CreatedAt:       grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		Quantity: &types.Float32RangeWithOptionalMax{
			Max: input.Quantity.Max,
			Min: input.Quantity.Min,
		},
		ArchivedAt:                grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		Ingredient:                ConvertValidIngredientToGRPCValidIngredient(input.Ingredient),
		LastUpdatedAt:             grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		BelongsToRecipeStep:       input.BelongsToRecipeStep,
		IngredientNotes:           input.IngredientNotes,
		QuantityNotes:             input.QuantityNotes,
		ID:                        input.ID,
		Name:                      input.Name,
		OptionIndex:               uint32(input.OptionIndex),
		Optional:                  input.Optional,
		ToTaste:                   input.ToTaste,
		ProductPercentageToUse:    input.ProductPercentageToUse,
		RecipeStepProductRecipeID: input.RecipeStepProductRecipeID,
		RecipeStepProductID:       input.RecipeStepProductID,
	}

	if input.VesselIndex != nil {
		ingredient.VesselIndex = pointer.To(uint32(*input.VesselIndex))
	}

	return ingredient
}

func ConvertRecipeStepProductToGRPCRecipeStepProduct(input *mealplanning.RecipeStepProduct) *mealplanningsvc.RecipeStepProduct {
	rsp := &mealplanningsvc.RecipeStepProduct{
		CreatedAt: grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		StorageTemperatureInCelsius: &types.OptionalFloat32Range{
			Max: input.StorageTemperatureInCelsius.Max,
			Min: input.StorageTemperatureInCelsius.Min,
		},
		StorageDurationInSeconds: &types.OptionalUint32Range{
			Max: input.StorageDurationInSeconds.Max,
			Min: input.StorageDurationInSeconds.Min,
		},
		Quantity: &types.OptionalFloat32Range{
			Max: input.Quantity.Max,
			Min: input.Quantity.Min,
		},
		ArchivedAt:          grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		MeasurementUnit:     ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(input.MeasurementUnit),
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Name:                input.Name,
		Type:                input.Type,
		ID:                  input.ID,
		StorageInstructions: input.StorageInstructions,
		QuantityNotes:       input.QuantityNotes,
		Index:               uint32(input.Index),
		IsWaste:             input.IsWaste,
		IsLiquid:            input.IsLiquid,
		Compostable:         input.Compostable,
	}

	if input.ContainedInVesselIndex != nil {
		rsp.ContainedInVesselIndex = uint32(*input.ContainedInVesselIndex)
	}

	return rsp
}

func ConvertRecipeMediaToGRPCRecipeMedia(input *mealplanning.RecipeMedia) *mealplanningsvc.RecipeMedia {
	x := &mealplanningsvc.RecipeMedia{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		ArchivedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ID:            input.ID,
		MimeType:      input.MimeType,
		InternalPath:  input.InternalPath,
		ExternalPath:  input.ExternalPath,
		Index:         uint32(input.Index),
	}

	if input.BelongsToRecipe != nil {
		x.BelongsToRecipe = *input.BelongsToRecipe
	}

	if input.BelongsToRecipe != nil {
		x.BelongsToRecipe = *input.BelongsToRecipe
	}

	return x
}

func ConvertRecipePrepTaskToGRPCRecipePrepTask(input *mealplanning.RecipePrepTask) *mealplanningsvc.RecipePrepTask {
	recipePrepTask := &mealplanningsvc.RecipePrepTask{
		CreatedAt: grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		StorageTemperatureInCelsius: &types.OptionalFloat32Range{
			Max: input.StorageTemperatureInCelsius.Max,
			Min: input.StorageTemperatureInCelsius.Min,
		},
		TimeBufferBeforeRecipeInSeconds: &types.Uint32RangeWithOptionalMax{
			Max: input.TimeBufferBeforeRecipeInSeconds.Max,
			Min: input.TimeBufferBeforeRecipeInSeconds.Min,
		},
		ArchivedAt:                  grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt:               grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		BelongsToRecipe:             input.BelongsToRecipe,
		StorageType:                 input.StorageType,
		ID:                          input.ID,
		ExplicitStorageInstructions: input.ExplicitStorageInstructions,
		Notes:                       input.Notes,
		Name:                        input.Name,
		Description:                 input.Description,
		Optional:                    input.Optional,
	}

	for _, taskStep := range input.TaskSteps {
		recipePrepTask.TaskSteps = append(recipePrepTask.TaskSteps, ConvertRecipePrepTaskStepToGRPCRecipePrepTaskStep(taskStep))
	}

	return recipePrepTask
}

func ConvertRecipePrepTaskStepToGRPCRecipePrepTaskStep(input *mealplanning.RecipePrepTaskStep) *mealplanningsvc.RecipePrepTaskStep {
	return &mealplanningsvc.RecipePrepTaskStep{
		ID:                      input.ID,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}
}
