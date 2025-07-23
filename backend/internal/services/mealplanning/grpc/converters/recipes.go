package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	types2 "github.com/dinnerdonebetter/backend/internal/platform/types"
)

func ConvertGRPCRecipeCreationRequestInputToRecipeCreationRequestInput(input *mealplanningsvc.RecipeCreationRequestInput) *mealplanning.RecipeCreationRequestInput {
	var steps []*mealplanning.RecipeStepCreationRequestInput
	for _, step := range input.Steps {
		steps = append(steps, ConvertGRPCRecipeStepCreationRequestInputToRecipeStepCreationRequestInput(step))
	}

	var prepTasks []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput
	for _, task := range input.PrepTasks {
		prepTasks = append(prepTasks, ConvertGRPCRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskWithinRecipeCreationRequestInput(task))
	}

	return &mealplanning.RecipeCreationRequestInput{
		InspiredByRecipeID:  input.InspiredByRecipeID,
		Name:                input.Name,
		Source:              input.Source,
		Description:         input.Description,
		PluralPortionName:   input.PluralPortionName,
		PortionName:         input.PortionName,
		Slug:                input.Slug,
		YieldsComponentType: input.YieldsComponentType,
		EstimatedPortions: types2.Float32RangeWithOptionalMax{
			Max: input.EstimatedPortions.Max,
			Min: input.EstimatedPortions.Min,
		},
		AlsoCreateMeal:   input.AlsoCreateMeal,
		SealOfApproval:   input.SealOfApproval,
		EligibleForMeals: input.EligibleForMeals,
		PrepTasks:        prepTasks,
		Steps:            steps,
	}
}

func ConvertGRPCRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskWithinRecipeCreationRequestInput(input *mealplanningsvc.RecipePrepTaskWithinRecipeCreationRequestInput) *mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput {
	var prepTaskSteps []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput
	for _, step := range input.RecipeSteps {
		prepTaskSteps = append(prepTaskSteps, ConvertGRPCRecipePrepTaskStepWithinRecipeCreationRequestInputToRecipePrepTaskStepWithinRecipeCreationRequestInput(step))
	}

	return &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
		StorageTemperatureInCelsius: types2.OptionalFloat32Range{
			Min: input.StorageTemperatureInCelsius.Min,
			Max: input.StorageTemperatureInCelsius.Max,
		},
		TimeBufferBeforeRecipeInSeconds: types2.Uint32RangeWithOptionalMax{
			Min: input.TimeBufferBeforeRecipeInSeconds.Min,
			Max: input.TimeBufferBeforeRecipeInSeconds.Max,
		},
		StorageType:                 input.StorageType,
		ExplicitStorageInstructions: input.ExplicitStorageInstructions,
		Notes:                       input.Notes,
		Name:                        input.Name,
		Description:                 input.Description,
		BelongsToRecipe:             input.BelongsToRecipe,
		RecipeSteps:                 prepTaskSteps,
		Optional:                    input.Optional,
	}
}

func ConvertGRPCRecipePrepTaskCreationRequestInputToRecipePrepTaskCreationRequestInput(input *mealplanningsvc.RecipePrepTaskCreationRequestInput) *mealplanning.RecipePrepTaskCreationRequestInput {
	var steps []*mealplanning.RecipePrepTaskStepCreationRequestInput
	for _, step := range input.RecipeSteps {
		steps = append(steps, ConvertGRPCRecipePrepTaskStepCreationRequestInputToRecipePrepTaskStepCreationRequestInput(step))
	}

	return &mealplanning.RecipePrepTaskCreationRequestInput{
		StorageTemperatureInCelsius:     types2.OptionalFloat32Range{},
		TimeBufferBeforeRecipeInSeconds: types2.Uint32RangeWithOptionalMax{},
		StorageType:                     input.StorageType,
		ExplicitStorageInstructions:     input.ExplicitStorageInstructions,
		Notes:                           input.Notes,
		Name:                            input.Name,
		Description:                     input.Description,
		BelongsToRecipe:                 input.BelongsToRecipe,
		RecipeSteps:                     steps,
		Optional:                        input.Optional,
	}
}

func ConvertGRPCRecipeRatingCreationRequestInputToRecipeRatingCreationRequestInput(input *mealplanningsvc.RecipeRatingCreationRequestInput) *mealplanning.RecipeRatingCreationRequestInput {
	return &mealplanning.RecipeRatingCreationRequestInput{
		RecipeID:     input.RecipeID,
		Notes:        input.Notes,
		ByUser:       input.ByUser,
		Taste:        input.Taste,
		Difficulty:   input.Difficulty,
		Cleanup:      input.Cleanup,
		Instructions: input.Instructions,
		Overall:      input.Overall,
	}
}

func ConvertGRPCRecipePrepTaskStepCreationRequestInputToRecipePrepTaskStepCreationRequestInput(input *mealplanningsvc.RecipePrepTaskStepCreationRequestInput) *mealplanning.RecipePrepTaskStepCreationRequestInput {
	return &mealplanning.RecipePrepTaskStepCreationRequestInput{
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		SatisfiesRecipeStep: input.SatisfiesRecipeStep,
	}
}

func ConvertGRPCRecipePrepTaskStepWithinRecipeCreationRequestInputToRecipePrepTaskStepWithinRecipeCreationRequestInput(input *mealplanningsvc.RecipePrepTaskStepWithinRecipeCreationRequestInput) *mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput {
	return &mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
		BelongsToRecipeStepIndex: input.BelongsToRecipeStepIndex,
		SatisfiesRecipeStep:      input.SatisfiesRecipeStep,
	}
}

func ConvertGRPCRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentCreationRequestInput(input *mealplanningsvc.RecipeStepInstrumentCreationRequestInput) *mealplanning.RecipeStepInstrumentCreationRequestInput {
	return &mealplanning.RecipeStepInstrumentCreationRequestInput{
		InstrumentID:                    input.InstrumentID,
		RecipeStepProductID:             input.RecipeStepProductID,
		ProductOfRecipeStepIndex:        input.ProductOfRecipeStepIndex,
		ProductOfRecipeStepProductIndex: input.ProductOfRecipeStepProductIndex,
		Notes:                           input.Notes,
		Name:                            input.Name,
		OptionIndex:                     uint16(input.OptionIndex),
		Optional:                        input.Optional,
		PreferenceRank:                  uint8(input.PreferenceRank),
		Quantity: types2.Uint32RangeWithOptionalMax{
			Min: input.Quantity.Min,
			Max: input.Quantity.Max,
		},
	}
}

func ConvertGRPCRecipeStepVesselCreationRequestInputToRecipeStepVesselCreationRequestInput(input *mealplanningsvc.RecipeStepVesselCreationRequestInput) *mealplanning.RecipeStepVesselCreationRequestInput {
	return &mealplanning.RecipeStepVesselCreationRequestInput{
		RecipeStepProductID:             input.RecipeStepProductID,
		ProductOfRecipeStepIndex:        input.ProductOfRecipeStepIndex,
		ProductOfRecipeStepProductIndex: input.ProductOfRecipeStepProductIndex,
		VesselID:                        input.VesselID,
		Name:                            input.Name,
		Notes:                           input.Notes,
		VesselPreposition:               input.VesselPreposition,
		UnavailableAfterStep:            input.UnavailableAfterStep,
		Quantity: types2.Uint16RangeWithOptionalMax{
			Min: uint16(input.Quantity.Min),
			Max: pointer.To(uint16(pointer.Dereference(input.Quantity.Max))),
		},
	}
}

func ConvertGRPCRecipeStepProductCreationRequestInputToRecipeStepProductCreationRequestInput(input *mealplanningsvc.RecipeStepProductCreationRequestInput) *mealplanning.RecipeStepProductCreationRequestInput {
	return &mealplanning.RecipeStepProductCreationRequestInput{
		MeasurementUnitID:      input.MeasurementUnitID,
		ContainedInVesselIndex: pointer.To(uint16(pointer.Dereference(input.ContainedInVesselIndex))),
		QuantityNotes:          input.QuantityNotes,
		Name:                   input.Name,
		StorageInstructions:    input.StorageInstructions,
		Type:                   input.Type,
		Index:                  uint16(input.Index),
		Compostable:            input.Compostable,
		IsLiquid:               input.IsLiquid,
		IsWaste:                input.IsWaste,
		StorageTemperatureInCelsius: types2.OptionalFloat32Range{
			Min: input.StorageTemperatureInCelsius.Min,
			Max: input.StorageTemperatureInCelsius.Max,
		},
		StorageDurationInSeconds: types2.OptionalUint32Range{
			Min: input.StorageDurationInSeconds.Min,
			Max: input.StorageDurationInSeconds.Max,
		},
		Quantity: types2.OptionalFloat32Range{
			Min: input.Quantity.Min,
			Max: input.Quantity.Max,
		},
	}
}

func ConvertGRPCRecipeStepIngredientCreationRequestInputToRecipeStepIngredientCreationRequestInput(input *mealplanningsvc.RecipeStepIngredientCreationRequestInput) *mealplanning.RecipeStepIngredientCreationRequestInput {
	return &mealplanning.RecipeStepIngredientCreationRequestInput{
		IngredientID:                    input.IngredientID,
		ProductOfRecipeStepIndex:        input.ProductOfRecipeStepIndex,
		ProductOfRecipeStepProductIndex: input.ProductOfRecipeStepProductIndex,
		VesselIndex:                     pointer.To(uint16(pointer.Dereference(input.VesselIndex))),
		ProductPercentageToUse:          input.ProductPercentageToUse,
		RecipeStepProductRecipeID:       input.RecipeStepProductRecipeID,
		IngredientNotes:                 input.IngredientNotes,
		MeasurementUnitID:               input.MeasurementUnitID,
		Name:                            input.Name,
		QuantityNotes:                   input.QuantityNotes,
		OptionIndex:                     uint16(input.OptionIndex),
		Optional:                        input.Optional,
		ToTaste:                         input.ToTaste,
		Quantity: types2.Float32RangeWithOptionalMax{
			Min: input.Quantity.Min,
			Max: input.Quantity.Max,
		},
	}
}

func ConvertGRPCRecipeStepCompletionConditionCreationRequestInputToRecipeStepCompletionConditionCreationRequestInput(input *mealplanningsvc.RecipeStepCompletionConditionCreationRequestInput) *mealplanning.RecipeStepCompletionConditionCreationRequestInput {
	return &mealplanning.RecipeStepCompletionConditionCreationRequestInput{
		IngredientStateID:   input.IngredientStateID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Notes:               input.Notes,
		Ingredients:         input.Ingredients,
		Optional:            input.Optional,
	}
}

func ConvertGRPCRecipeStepCompletionConditionForExistingRecipeCreationRequestInputToRecipeStepCompletionConditionForExistingRecipeCreationRequestInput(input *mealplanningsvc.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput) *mealplanning.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput {
	ingredients := []*mealplanning.RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput{}
	for _, ingredient := range input.Ingredients {
		ingredients = append(ingredients, &mealplanning.RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput{RecipeStepIngredient: ingredient.RecipeStepIngredient})
	}

	return &mealplanning.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{
		IngredientStateID:   input.IngredientStateID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Notes:               input.Notes,
		Ingredients:         ingredients,
		Optional:            input.Optional,
	}
}

func ConvertGRPCRecipeStepCreationRequestInputToRecipeStepCreationRequestInput(input *mealplanningsvc.RecipeStepCreationRequestInput) *mealplanning.RecipeStepCreationRequestInput {
	var recipeStepInstrumentCreationRequestInputs []*mealplanning.RecipeStepInstrumentCreationRequestInput
	for _, instrument := range input.Instruments {
		recipeStepInstrumentCreationRequestInputs = append(recipeStepInstrumentCreationRequestInputs, ConvertGRPCRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentCreationRequestInput(instrument))
	}

	var recipeStepVesselCreationRequestInputs []*mealplanning.RecipeStepVesselCreationRequestInput
	for _, vessel := range input.Vessels {
		recipeStepVesselCreationRequestInputs = append(recipeStepVesselCreationRequestInputs, ConvertGRPCRecipeStepVesselCreationRequestInputToRecipeStepVesselCreationRequestInput(vessel))
	}

	var recipeStepProductCreationRequestInputs []*mealplanning.RecipeStepProductCreationRequestInput
	for _, product := range input.Products {
		recipeStepProductCreationRequestInputs = append(recipeStepProductCreationRequestInputs, ConvertGRPCRecipeStepProductCreationRequestInputToRecipeStepProductCreationRequestInput(product))
	}

	var recipeStepIngredientCreationRequestInputs []*mealplanning.RecipeStepIngredientCreationRequestInput
	for _, ingredient := range input.Ingredients {
		recipeStepIngredientCreationRequestInputs = append(recipeStepIngredientCreationRequestInputs, ConvertGRPCRecipeStepIngredientCreationRequestInputToRecipeStepIngredientCreationRequestInput(ingredient))
	}

	var recipeStepCompletionConditionCreationRequestInputs []*mealplanning.RecipeStepCompletionConditionCreationRequestInput
	for _, completionCondition := range input.CompletionConditions {
		recipeStepCompletionConditionCreationRequestInputs = append(recipeStepCompletionConditionCreationRequestInputs, ConvertGRPCRecipeStepCompletionConditionCreationRequestInputToRecipeStepCompletionConditionCreationRequestInput(completionCondition))
	}

	return &mealplanning.RecipeStepCreationRequestInput{
		EstimatedTimeInSeconds: types2.OptionalUint32Range{
			Min: input.EstimatedTimeInSeconds.Min,
			Max: input.EstimatedTimeInSeconds.Max,
		},
		TemperatureInCelsius: types2.OptionalFloat32Range{
			Min: input.TemperatureInCelsius.Min,
			Max: input.TemperatureInCelsius.Max,
		},
		PreparationID:           input.PreparationID,
		Notes:                   input.Notes,
		ConditionExpression:     input.ConditionExpression,
		ExplicitInstructions:    input.ExplicitInstructions,
		Index:                   input.Index,
		Optional:                input.Optional,
		StartTimerAutomatically: input.StartTimerAutomatically,
		Instruments:             recipeStepInstrumentCreationRequestInputs,
		Vessels:                 recipeStepVesselCreationRequestInputs,
		Products:                recipeStepProductCreationRequestInputs,
		Ingredients:             recipeStepIngredientCreationRequestInputs,
		CompletionConditions:    recipeStepCompletionConditionCreationRequestInputs,
	}
}

func ConvertRecipeToGRPCRecipe(input *mealplanning.Recipe) *mealplanningsvc.Recipe {
	var recipeSteps []*mealplanningsvc.RecipeStep
	for _, step := range input.Steps {
		recipeSteps = append(recipeSteps, ConvertRecipeStepToGRPCRecipeStep(step))
	}

	var recipeMedia []*mealplanningsvc.RecipeMedia
	for _, media := range input.Media {
		recipeMedia = append(recipeMedia, ConvertRecipeMediaToGRPCRecipeMedia(media))
	}

	var recipePrepTasks []*mealplanningsvc.RecipePrepTask
	for _, rps := range input.PrepTasks {
		recipePrepTasks = append(recipePrepTasks, ConvertRecipePrepTaskToGRPCRecipePrepTask(rps))
	}

	var recipeSupportingRecipes []*mealplanningsvc.Recipe
	for _, r := range input.SupportingRecipes {
		recipeSupportingRecipes = append(recipeSupportingRecipes, ConvertRecipeToGRPCRecipe(r))
	}

	var recipeInspiredByRecipeID string
	if input.InspiredByRecipeID != nil {
		recipeInspiredByRecipeID = *input.InspiredByRecipeID
	}

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
		Steps:               recipeSteps,
		Media:               recipeMedia,
		PrepTasks:           recipePrepTasks,
		SupportingRecipes:   recipeSupportingRecipes,
		InspiredByRecipeID:  recipeInspiredByRecipeID,
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
	}

	for _, media := range input.Media {
		step.Media = append(step.Media, ConvertRecipeMediaToGRPCRecipeMedia(media))
	}

	for _, instrument := range input.Instruments {
		step.Instruments = append(step.Instruments, ConvertRecipeStepInstrumentToGRPCRecipeStepInstrument(instrument))
	}

	for _, vessel := range input.Vessels {
		step.Vessels = append(step.Vessels, ConvertRecipeStepVesselToGRPCRecipeStepVessel(vessel))
	}

	for _, completionCondition := range input.CompletionConditions {
		step.CompletionConditions = append(step.CompletionConditions, ConvertRecipeStepCompletionConditionToGRPCRecipeStepCompletionCondition(completionCondition))
	}

	for _, ingredient := range input.Ingredients {
		step.Ingredients = append(step.Ingredients, ConvertRecipeStepIngredientToGRPCRecipeStepIngredient(ingredient))
	}

	for _, product := range input.Products {
		step.Products = append(step.Products, ConvertRecipeStepProductToGRPCRecipeStepProduct(product))
	}

	return step
}

func ConvertRecipeStepInstrumentToGRPCRecipeStepInstrument(input *mealplanning.RecipeStepInstrument) *mealplanningsvc.RecipeStepInstrument {
	return &mealplanningsvc.RecipeStepInstrument{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		Instrument:    ConvertValidInstrumentToGRPCValidInstrument(input.Instrument),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		Quantity: &types.Uint32RangeWithOptionalMax{
			Max: input.Quantity.Max,
			Min: input.Quantity.Min,
		},
		ArchivedAt:          grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Name:                input.Name,
		Notes:               input.Notes,
		ID:                  input.ID,
		RecipeStepProductID: input.RecipeStepProductID,
		OptionIndex:         uint32(input.OptionIndex),
		PreferenceRank:      uint32(input.PreferenceRank),
		Optional:            input.Optional,
	}
}

func ConvertRecipeStepVesselToGRPCRecipeStepVessel(input *mealplanning.RecipeStepVessel) *mealplanningsvc.RecipeStepVessel {
	return &mealplanningsvc.RecipeStepVessel{
		Vessel: ConvertValidVesselToGRPCValidVessel(input.Vessel),
		Quantity: &types.Uint16RangeWithOptionalMax{
			Max: pointer.To(uint32(pointer.Dereference(input.Quantity.Max))),
			Min: uint32(input.Quantity.Min),
		},
		CreatedAt:            grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:        grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:           grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		VesselPreposition:    input.VesselPreposition,
		Notes:                input.Notes,
		RecipeStepProductID:  input.RecipeStepProductID,
		BelongsToRecipeStep:  input.BelongsToRecipeStep,
		ID:                   input.ID,
		Name:                 input.Name,
		UnavailableAfterStep: input.UnavailableAfterStep,
	}
}

func ConvertRecipeStepCompletionConditionToGRPCRecipeStepCompletionCondition(input *mealplanning.RecipeStepCompletionCondition) *mealplanningsvc.RecipeStepCompletionCondition {
	recipeStepCompletionCondition := &mealplanningsvc.RecipeStepCompletionCondition{
		CreatedAt:           grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		ArchivedAt:          grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ID:                  input.ID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Notes:               input.Notes,
		Optional:            input.Optional,
		IngredientState:     ConvertValidIngredientStateToGRPCValidIngredientState(&input.IngredientState),
	}

	for _, ingredient := range input.Ingredients {
		recipeStepCompletionCondition.Ingredients = append(recipeStepCompletionCondition.Ingredients, ConvertRecipeStepCompletionConditionIngredientToGRPCRecipeStepCompletionConditionIngredient(ingredient))
	}

	return recipeStepCompletionCondition
}

func ConvertRecipeStepCompletionConditionIngredientToGRPCRecipeStepCompletionConditionIngredient(input *mealplanning.RecipeStepCompletionConditionIngredient) *mealplanningsvc.RecipeStepCompletionConditionIngredient {
	return &mealplanningsvc.RecipeStepCompletionConditionIngredient{
		CreatedAt:                              grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		ArchivedAt:                             grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt:                          grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ID:                                     input.ID,
		BelongsToRecipeStepCompletionCondition: input.BelongsToRecipeStepCompletionCondition,
		RecipeStepIngredient:                   input.RecipeStepIngredient,
	}
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

func ConvertRecipeRatingToGRPCRecipeRating(input *mealplanning.RecipeRating) *mealplanningsvc.RecipeRating {
	return &mealplanningsvc.RecipeRating{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		RecipeID:      input.RecipeID,
		ID:            input.ID,
		Notes:         input.Notes,
		ByUser:        input.ByUser,
		Taste:         input.Taste,
		Instructions:  input.Instructions,
		Overall:       input.Overall,
		Cleanup:       input.Cleanup,
		Difficulty:    input.Difficulty,
	}
}

func ConvertGRPCRecipeUpdateRequestInputToRecipeUpdateRequestInput(input *mealplanningsvc.RecipeUpdateRequestInput) *mealplanning.RecipeUpdateRequestInput {
	return &mealplanning.RecipeUpdateRequestInput{
		Name:                input.Name,
		Slug:                input.Slug,
		Source:              input.Source,
		Description:         input.Description,
		InspiredByRecipeID:  input.InspiredByRecipeID,
		SealOfApproval:      input.SealOfApproval,
		PortionName:         input.PortionName,
		PluralPortionName:   input.PluralPortionName,
		EligibleForMeals:    input.EligibleForMeals,
		YieldsComponentType: input.YieldsComponentType,
		EstimatedPortions: types2.Float32RangeWithOptionalMaxUpdateRequestInput{
			Min: input.EstimatedPortions.Min,
			Max: input.EstimatedPortions.Max,
		},
	}
}

func ConvertGRPCRecipePrepTaskUpdateRequestInputToRecipePrepTaskUpdateRequestInput(input *mealplanningsvc.RecipePrepTaskUpdateRequestInput) *mealplanning.RecipePrepTaskUpdateRequestInput {
	var taskSteps []*mealplanning.RecipePrepTaskStepUpdateRequestInput
	for _, step := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertGRPCRecipePrepTaskStepUpdateRequestInputToRecipePrepTaskStepUpdateRequestInput(step))
	}

	return &mealplanning.RecipePrepTaskUpdateRequestInput{
		Notes:                           input.Notes,
		ExplicitStorageInstructions:     input.ExplicitStorageInstructions,
		StorageType:                     input.StorageType,
		Name:                            input.Name,
		Optional:                        input.Optional,
		Description:                     input.Description,
		BelongsToRecipe:                 input.BelongsToRecipe,
		TaskSteps:                       taskSteps,
		TimeBufferBeforeRecipeInSeconds: types2.Uint32RangeWithOptionalMaxUpdateRequestInput{},
		StorageTemperatureInCelsius:     types2.OptionalFloat32Range{},
	}
}

func ConvertGRPCRecipePrepTaskStepUpdateRequestInputToRecipePrepTaskStepUpdateRequestInput(input *mealplanningsvc.RecipePrepTaskStepUpdateRequestInput) *mealplanning.RecipePrepTaskStepUpdateRequestInput {
	return &mealplanning.RecipePrepTaskStepUpdateRequestInput{
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
	}
}

func ConvertGRPCRecipeRatingUpdateRequestInputToRecipeRatingUpdateRequestInput(input *mealplanningsvc.RecipeRatingUpdateRequestInput) *mealplanning.RecipeRatingUpdateRequestInput {
	return &mealplanning.RecipeRatingUpdateRequestInput{
		RecipeID:     input.RecipeID,
		Taste:        input.Taste,
		Difficulty:   input.Difficulty,
		Cleanup:      input.Cleanup,
		Instructions: input.Instructions,
		Overall:      input.Overall,
		Notes:        input.Notes,
		ByUser:       input.ByUser,
	}
}

func ConvertGRPCRecipeStepUpdateRequestInputToRecipeStepUpdateRequestInput(input *mealplanningsvc.RecipeStepUpdateRequestInput) *mealplanning.RecipeStepUpdateRequestInput {
	return &mealplanning.RecipeStepUpdateRequestInput{
		EstimatedTimeInSeconds: types2.OptionalUint32Range{
			Min: input.EstimatedTimeInSeconds.Min,
			Max: input.EstimatedTimeInSeconds.Max,
		},
		TemperatureInCelsius: types2.OptionalFloat32Range{
			Min: input.TemperatureInCelsius.Min,
			Max: input.TemperatureInCelsius.Max,
		},
		Notes:                   input.Notes,
		Preparation:             ConvertGRPCValidPreparationToValidPreparation(input.Preparation),
		Index:                   input.Index,
		Optional:                input.Optional,
		ExplicitInstructions:    input.ExplicitInstructions,
		ConditionExpression:     input.ConditionExpression,
		StartTimerAutomatically: input.StartTimerAutomatically,
		BelongsToRecipe:         input.BelongsToRecipe,
	}
}

func ConvertGRPCRecipeStepCompletionConditionUpdateRequestInputToRecipeStepCompletionConditionUpdateRequestInput(input *mealplanningsvc.RecipeStepCompletionConditionUpdateRequestInput) *mealplanning.RecipeStepCompletionConditionUpdateRequestInput {
	return &mealplanning.RecipeStepCompletionConditionUpdateRequestInput{
		IngredientStateID:   input.IngredientStateID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Notes:               input.Notes,
		Optional:            input.Optional,
	}
}

func ConvertGRPCRecipeStepIngredientUpdateRequestInputToRecipeStepIngredientUpdateRequestInput(input *mealplanningsvc.RecipeStepIngredientUpdateRequestInput) *mealplanning.RecipeStepIngredientUpdateRequestInput {
	return &mealplanning.RecipeStepIngredientUpdateRequestInput{
		IngredientID:              input.IngredientID,
		RecipeStepProductID:       input.RecipeStepProductID,
		Name:                      input.Name,
		Optional:                  input.Optional,
		MeasurementUnitID:         input.MeasurementUnitID,
		QuantityNotes:             input.QuantityNotes,
		IngredientNotes:           input.IngredientNotes,
		BelongsToRecipeStep:       input.BelongsToRecipeStep,
		OptionIndex:               grpcconverters.ConvertUint32PointerToUint16Pointer(input.OptionIndex),
		VesselIndex:               grpcconverters.ConvertUint32PointerToUint16Pointer(input.VesselIndex),
		ToTaste:                   input.ToTaste,
		ProductPercentageToUse:    input.ProductPercentageToUse,
		RecipeStepProductRecipeID: input.RecipeStepProductRecipeID,
		Quantity:                  types2.Float32RangeWithOptionalMaxUpdateRequestInput{},
	}
}

func ConvertGRPCRecipeStepInstrumentUpdateRequestInputToRecipeStepInstrumentUpdateRequestInput(input *mealplanningsvc.RecipeStepInstrumentUpdateRequestInput) *mealplanning.RecipeStepInstrumentUpdateRequestInput {
	return &mealplanning.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID:        input.InstrumentID,
		RecipeStepProductID: input.RecipeStepProductID,
		Notes:               input.Notes,
		PreferenceRank:      pointer.To(uint8(pointer.Dereference(input.PreferenceRank))),
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Name:                input.Name,
		Optional:            input.Optional,
		OptionIndex:         grpcconverters.ConvertUint32PointerToUint16Pointer(input.OptionIndex),
		Quantity: types2.Uint32RangeWithOptionalMaxUpdateRequestInput{
			Min: input.Quantity.Min,
			Max: input.Quantity.Max,
		},
	}
}

func ConvertGRPCRecipeStepProductUpdateRequestInputToRecipeStepProductUpdateRequestInput(input *mealplanningsvc.RecipeStepProductUpdateRequestInput) *mealplanning.RecipeStepProductUpdateRequestInput {
	return &mealplanning.RecipeStepProductUpdateRequestInput{
		Name:                input.Name,
		Type:                input.Type,
		MeasurementUnitID:   input.MeasurementUnitID,
		QuantityNotes:       input.QuantityNotes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		StorageTemperatureInCelsius: types2.OptionalFloat32Range{
			Max: input.StorageTemperatureInCelsius.Max,
			Min: input.StorageTemperatureInCelsius.Min,
		},
		StorageDurationInSeconds: types2.OptionalUint32Range{
			Max: input.StorageDurationInSeconds.Max,
			Min: input.StorageDurationInSeconds.Min,
		},
		Quantity: types2.OptionalFloat32Range{
			Max: input.Quantity.Max,
			Min: input.Quantity.Min,
		},
		Compostable:            input.Compostable,
		StorageInstructions:    input.StorageInstructions,
		IsLiquid:               input.IsLiquid,
		IsWaste:                input.IsWaste,
		Index:                  grpcconverters.ConvertUint32PointerToUint16Pointer(input.Index),
		ContainedInVesselIndex: grpcconverters.ConvertUint32PointerToUint16Pointer(input.ContainedInVesselIndex),
	}
}

func ConvertGRPCRecipeStepVesselUpdateRequestInputToRecipeStepVesselUpdateRequestInput(input *mealplanningsvc.RecipeStepVesselUpdateRequestInput) *mealplanning.RecipeStepVesselUpdateRequestInput {
	return &mealplanning.RecipeStepVesselUpdateRequestInput{
		RecipeStepProductID:  input.RecipeStepProductID,
		Name:                 input.Name,
		Notes:                input.Notes,
		BelongsToRecipeStep:  input.BelongsToRecipeStep,
		VesselID:             input.VesselID,
		VesselPreposition:    input.VesselPreposition,
		UnavailableAfterStep: input.UnavailableAfterStep,
		Quantity: types2.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: pointer.To(uint16(pointer.Dereference(input.Quantity.Min))),
			Max: grpcconverters.ConvertUint32PointerToUint16Pointer(input.Quantity.Max),
		},
	}
}
