package mealplanning

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ mealplanning.RecipeStepDataManager = (*repository)(nil)
)

// RecipeStepExists fetches whether a recipe step exists from the database.
func (q *repository) RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	result, err := q.generatedQuerier.CheckRecipeStepExistence(ctx, q.db, &generated.CheckRecipeStepExistenceParams{
		RecipeID:     recipeID,
		RecipeStepID: recipeStepID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step existence check")
	}

	return result, nil
}

// GetRecipeStep fetches a recipe step from the database.
func (q *repository) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*mealplanning.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	result, err := q.generatedQuerier.GetRecipeStep(ctx, q.db, &generated.GetRecipeStepParams{
		RecipeID:     recipeID,
		RecipeStepID: recipeStepID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step")
	}

	recipeStep := &mealplanning.RecipeStep{
		CreatedAt: result.CreatedAt,
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Max: database.Uint32PointerFromNullInt64(result.MaximumEstimatedTimeInSeconds),
			Min: database.Uint32PointerFromNullInt64(result.MinimumEstimatedTimeInSeconds),
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Max: database.Float32PointerFromNullString(result.MaximumTemperatureInCelsius),
			Min: database.Float32PointerFromNullString(result.MinimumTemperatureInCelsius),
		},
		ArchivedAt:           database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:        database.TimePointerFromNullTime(result.LastUpdatedAt),
		BelongsToRecipe:      result.BelongsToRecipe,
		ConditionExpression:  result.ConditionExpression,
		ID:                   result.ID,
		Notes:                result.Notes,
		ExplicitInstructions: result.ExplicitInstructions,
		Media:                []*mealplanning.RecipeMedia{},
		Products:             []*mealplanning.RecipeStepProduct{},
		Instruments:          []*mealplanning.RecipeStepInstrument{},
		Vessels:              []*mealplanning.RecipeStepVessel{},
		CompletionConditions: []*mealplanning.RecipeStepCompletionCondition{},
		Ingredients:          []*mealplanning.RecipeStepIngredient{},
		Preparation: mealplanning.ValidPreparation{
			CreatedAt: result.ValidPreparationCreatedAt,
			InstrumentCount: types.Uint16RangeWithOptionalMax{
				Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
				Min: uint16(result.ValidPreparationMinimumInstrumentCount),
			},
			IngredientCount: types.Uint16RangeWithOptionalMax{
				Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
				Min: uint16(result.ValidPreparationMinimumIngredientCount),
			},
			VesselCount: types.Uint16RangeWithOptionalMax{
				Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
				Min: uint16(result.ValidPreparationMinimumVesselCount),
			},
			ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
			LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
			IconPath:                    result.ValidPreparationIconPath,
			PastTense:                   result.ValidPreparationPastTense,
			ID:                          result.ValidPreparationID,
			Name:                        result.ValidPreparationName,
			Description:                 result.ValidPreparationDescription,
			Slug:                        result.ValidPreparationSlug,
			RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
			TemperatureRequired:         result.ValidPreparationTemperatureRequired,
			TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
			ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
			ConsumesVessel:              result.ValidPreparationConsumesVessel,
			OnlyForVessels:              result.ValidPreparationOnlyForVessels,
			YieldsNothing:               result.ValidPreparationYieldsNothing,
		},
		Index:                   uint32(result.Index),
		Optional:                result.Optional,
		StartTimerAutomatically: result.StartTimerAutomatically,
	}

	// Fetch related data for this recipe step
	ingredients, err := q.getRecipeStepIngredientsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step ingredients for recipe step")
	}
	for _, ingredient := range ingredients {
		if ingredient.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.Ingredients = append(recipeStep.Ingredients, ingredient)
		}
	}

	products, err := q.getRecipeStepProductsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step products for recipe step")
	}
	for _, product := range products {
		if product.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.Products = append(recipeStep.Products, product)
		}
	}

	instruments, err := q.getRecipeStepInstrumentsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step instruments for recipe step")
	}
	for _, instrument := range instruments {
		if instrument.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.Instruments = append(recipeStep.Instruments, instrument)
		}
	}

	vessels, err := q.getRecipeStepVesselsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step vessels for recipe step")
	}
	for _, vessel := range vessels {
		if vessel.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.Vessels = append(recipeStep.Vessels, vessel)
		}
	}

	completionConditions, err := q.getRecipeStepCompletionConditionsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step completion conditions for recipe step")
	}
	for _, completionCondition := range completionConditions {
		if completionCondition.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.CompletionConditions = append(recipeStep.CompletionConditions, completionCondition)
		}
	}

	recipeMedia, err := q.getRecipeMediaForRecipeStep(ctx, recipeID, recipeStep.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe media for recipe step")
	}
	recipeStep.Media = recipeMedia

	return recipeStep, nil
}

// getRecipeStepByID fetches a recipe step from the database.
func (q *repository) getRecipeStepByID(ctx context.Context, querier database.SQLQueryExecutor, recipeStepID string) (*mealplanning.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	result, err := q.generatedQuerier.GetRecipeStepByRecipeID(ctx, querier, recipeStepID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step")
	}

	recipeStep := &mealplanning.RecipeStep{
		CreatedAt: result.CreatedAt,
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Max: database.Uint32PointerFromNullInt64(result.MaximumEstimatedTimeInSeconds),
			Min: database.Uint32PointerFromNullInt64(result.MinimumEstimatedTimeInSeconds),
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Max: database.Float32PointerFromNullString(result.MaximumTemperatureInCelsius),
			Min: database.Float32PointerFromNullString(result.MinimumTemperatureInCelsius),
		},
		ArchivedAt:           database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:        database.TimePointerFromNullTime(result.LastUpdatedAt),
		BelongsToRecipe:      result.BelongsToRecipe,
		ConditionExpression:  result.ConditionExpression,
		ID:                   result.ID,
		Notes:                result.Notes,
		ExplicitInstructions: result.ExplicitInstructions,
		Media:                []*mealplanning.RecipeMedia{},
		Products:             []*mealplanning.RecipeStepProduct{},
		Instruments:          []*mealplanning.RecipeStepInstrument{},
		Vessels:              []*mealplanning.RecipeStepVessel{},
		CompletionConditions: []*mealplanning.RecipeStepCompletionCondition{},
		Ingredients:          []*mealplanning.RecipeStepIngredient{},
		Preparation: mealplanning.ValidPreparation{
			CreatedAt: result.ValidPreparationCreatedAt,
			InstrumentCount: types.Uint16RangeWithOptionalMax{
				Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
				Min: uint16(result.ValidPreparationMinimumInstrumentCount),
			},
			IngredientCount: types.Uint16RangeWithOptionalMax{
				Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
				Min: uint16(result.ValidPreparationMinimumIngredientCount),
			},
			VesselCount: types.Uint16RangeWithOptionalMax{
				Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
				Min: uint16(result.ValidPreparationMinimumVesselCount),
			},
			ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
			LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
			IconPath:                    result.ValidPreparationIconPath,
			PastTense:                   result.ValidPreparationPastTense,
			ID:                          result.ValidPreparationID,
			Name:                        result.ValidPreparationName,
			Description:                 result.ValidPreparationDescription,
			Slug:                        result.ValidPreparationSlug,
			RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
			TemperatureRequired:         result.ValidPreparationTemperatureRequired,
			TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
			ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
			ConsumesVessel:              result.ValidPreparationConsumesVessel,
			OnlyForVessels:              result.ValidPreparationOnlyForVessels,
			YieldsNothing:               result.ValidPreparationYieldsNothing,
		},
		Index:                   uint32(result.Index),
		Optional:                result.Optional,
		StartTimerAutomatically: result.StartTimerAutomatically,
	}

	// Fetch related data for this recipe step
	ingredients, err := q.getRecipeStepIngredientsForRecipe(ctx, result.BelongsToRecipe)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step ingredients for recipe step")
	}
	for _, ingredient := range ingredients {
		if ingredient.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.Ingredients = append(recipeStep.Ingredients, ingredient)
		}
	}

	products, err := q.getRecipeStepProductsForRecipe(ctx, result.BelongsToRecipe)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step products for recipe step")
	}
	for _, product := range products {
		if product.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.Products = append(recipeStep.Products, product)
		}
	}

	instruments, err := q.getRecipeStepInstrumentsForRecipe(ctx, result.BelongsToRecipe)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step instruments for recipe step")
	}
	for _, instrument := range instruments {
		if instrument.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.Instruments = append(recipeStep.Instruments, instrument)
		}
	}

	vessels, err := q.getRecipeStepVesselsForRecipe(ctx, result.BelongsToRecipe)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step vessels for recipe step")
	}
	for _, vessel := range vessels {
		if vessel.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.Vessels = append(recipeStep.Vessels, vessel)
		}
	}

	completionConditions, err := q.getRecipeStepCompletionConditionsForRecipe(ctx, result.BelongsToRecipe)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step completion conditions for recipe step")
	}
	for _, completionCondition := range completionConditions {
		if completionCondition.BelongsToRecipeStep == recipeStep.ID {
			recipeStep.CompletionConditions = append(recipeStep.CompletionConditions, completionCondition)
		}
	}

	recipeMedia, err := q.getRecipeMediaForRecipeStep(ctx, result.BelongsToRecipe, recipeStep.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe media for recipe step")
	}
	recipeStep.Media = recipeMedia

	return recipeStep, nil
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter.
func (q *repository) GetRecipeSteps(ctx context.Context, recipeID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.RecipeStep], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[mealplanning.RecipeStep]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeSteps(ctx, q.db, &generated.GetRecipeStepsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
		RecipeID:        recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe steps")
	}

	for _, result := range results {
		recipeStep := &mealplanning.RecipeStep{
			CreatedAt: result.CreatedAt,
			EstimatedTimeInSeconds: types.OptionalUint32Range{
				Max: database.Uint32PointerFromNullInt64(result.MaximumEstimatedTimeInSeconds),
				Min: database.Uint32PointerFromNullInt64(result.MinimumEstimatedTimeInSeconds),
			},
			TemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.MaximumTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.MinimumTemperatureInCelsius),
			},
			ArchivedAt:           database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:        database.TimePointerFromNullTime(result.LastUpdatedAt),
			BelongsToRecipe:      result.BelongsToRecipe,
			ConditionExpression:  result.ConditionExpression,
			ID:                   result.ID,
			Notes:                result.Notes,
			ExplicitInstructions: result.ExplicitInstructions,
			Media:                []*mealplanning.RecipeMedia{},
			Products:             []*mealplanning.RecipeStepProduct{},
			Instruments:          []*mealplanning.RecipeStepInstrument{},
			Vessels:              []*mealplanning.RecipeStepVessel{},
			CompletionConditions: []*mealplanning.RecipeStepCompletionCondition{},
			Ingredients:          []*mealplanning.RecipeStepIngredient{},
			Preparation: mealplanning.ValidPreparation{
				CreatedAt: result.ValidPreparationCreatedAt,
				InstrumentCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
					Min: uint16(result.ValidPreparationMinimumInstrumentCount),
				},
				IngredientCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
					Min: uint16(result.ValidPreparationMinimumInstrumentCount),
				},
				VesselCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
					Min: uint16(result.ValidPreparationMinimumVesselCount),
				},
				ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
				LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
				TemperatureRequired:         result.ValidPreparationTemperatureRequired,
				TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
				ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
				ConsumesVessel:              result.ValidPreparationConsumesVessel,
				OnlyForVessels:              result.ValidPreparationOnlyForVessels,
				YieldsNothing:               result.ValidPreparationYieldsNothing,
			},
			Index:                   uint32(result.Index),
			Optional:                result.Optional,
			StartTimerAutomatically: result.StartTimerAutomatically,
		}

		x.Data = append(x.Data, recipeStep)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	// Fetch all related data for all recipe steps
	ingredients, err := q.getRecipeStepIngredientsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step ingredients for recipe steps")
	}

	products, err := q.getRecipeStepProductsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step products for recipe steps")
	}

	instruments, err := q.getRecipeStepInstrumentsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step instruments for recipe steps")
	}

	vessels, err := q.getRecipeStepVesselsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step vessels for recipe steps")
	}

	completionConditions, err := q.getRecipeStepCompletionConditionsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step completion conditions for recipe steps")
	}

	// Populate each recipe step with its related data
	for i, step := range x.Data {
		for _, ingredient := range ingredients {
			if ingredient.BelongsToRecipeStep == step.ID {
				x.Data[i].Ingredients = append(x.Data[i].Ingredients, ingredient)
			}
		}

		for _, product := range products {
			if product.BelongsToRecipeStep == step.ID {
				x.Data[i].Products = append(x.Data[i].Products, product)
			}
		}

		for _, instrument := range instruments {
			if instrument.BelongsToRecipeStep == step.ID {
				x.Data[i].Instruments = append(x.Data[i].Instruments, instrument)
			}
		}

		for _, vessel := range vessels {
			if vessel.BelongsToRecipeStep == step.ID {
				x.Data[i].Vessels = append(x.Data[i].Vessels, vessel)
			}
		}

		for _, completionCondition := range completionConditions {
			if completionCondition.BelongsToRecipeStep == step.ID {
				x.Data[i].CompletionConditions = append(x.Data[i].CompletionConditions, completionCondition)
			}
		}

		recipeMedia, mediaErr := q.getRecipeMediaForRecipeStep(ctx, recipeID, step.ID)
		if mediaErr != nil {
			return nil, observability.PrepareAndLogError(mediaErr, logger, span, "fetching recipe media for recipe step")
		}
		x.Data[i].Media = recipeMedia
	}

	return x, nil
}

// CreateRecipeStep creates a recipe step in the database.
func (q *repository) createRecipeStep(ctx context.Context, db database.SQLQueryExecutor, input *mealplanning.RecipeStepDatabaseCreationInput) (*mealplanning.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	// create the recipe step.
	if err := q.generatedQuerier.CreateRecipeStep(ctx, db, &generated.CreateRecipeStepParams{
		ID:                            input.ID,
		BelongsToRecipe:               input.BelongsToRecipe,
		PreparationID:                 input.PreparationID,
		ConditionExpression:           input.ConditionExpression,
		ExplicitInstructions:          input.ExplicitInstructions,
		Notes:                         input.Notes,
		MaximumTemperatureInCelsius:   database.NullStringFromFloat32Pointer(input.TemperatureInCelsius.Max),
		MinimumTemperatureInCelsius:   database.NullStringFromFloat32Pointer(input.TemperatureInCelsius.Min),
		MaximumEstimatedTimeInSeconds: database.NullInt64FromUint32Pointer(input.EstimatedTimeInSeconds.Max),
		MinimumEstimatedTimeInSeconds: database.NullInt64FromUint32Pointer(input.EstimatedTimeInSeconds.Min),
		Index:                         int32(input.Index),
		Optional:                      input.Optional,
		StartTimerAutomatically:       input.StartTimerAutomatically,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step creation")
	}

	// Fetch the preparation data
	preparation, err := q.GetValidPreparation(ctx, input.PreparationID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching preparation data")
	}

	x := &mealplanning.RecipeStep{
		ID:                      input.ID,
		Index:                   input.Index,
		Preparation:             *preparation,
		EstimatedTimeInSeconds:  input.EstimatedTimeInSeconds,
		TemperatureInCelsius:    input.TemperatureInCelsius,
		Notes:                   input.Notes,
		ExplicitInstructions:    input.ExplicitInstructions,
		ConditionExpression:     input.ConditionExpression,
		Optional:                input.Optional,
		BelongsToRecipe:         input.BelongsToRecipe,
		StartTimerAutomatically: input.StartTimerAutomatically,
		CreatedAt:               q.CurrentTime(),
	}
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, x.ID)

	for i, ingredientInput := range input.Ingredients {
		ingredientInput.BelongsToRecipeStep = x.ID
		ingredient, createErr := q.createRecipeStepIngredient(ctx, db, ingredientInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, span, "creating recipe step ingredient #%d", i+1)
		}

		x.Ingredients = append(x.Ingredients, ingredient)
	}

	for i, productInput := range input.Products {
		productInput.BelongsToRecipeStep = x.ID
		product, createErr := q.createRecipeStepProduct(ctx, db, productInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, span, "creating recipe step product #%d", i+1)
		}

		x.Products = append(x.Products, product)
	}

	for i, instrumentInput := range input.Instruments {
		instrumentInput.BelongsToRecipeStep = x.ID
		instrument, createErr := q.createRecipeStepInstrument(ctx, db, instrumentInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, span, "creating recipe step instrument #%d", i+1)
		}

		x.Instruments = append(x.Instruments, instrument)
	}

	for i, vesselInput := range input.Vessels {
		vesselInput.BelongsToRecipeStep = x.ID
		vessel, createErr := q.createRecipeStepVessel(ctx, db, vesselInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, span, "creating recipe step vessel #%d", i+1)
		}

		x.Vessels = append(x.Vessels, vessel)
	}

	for i, conditionInput := range input.CompletionConditions {
		conditionInput.BelongsToRecipeStep = x.ID
		condition, createErr := q.createRecipeStepCompletionCondition(ctx, db, conditionInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, span, "creating recipe step completion condition #%d", i+1)
		}

		x.CompletionConditions = append(x.CompletionConditions, condition)
	}

	return x, nil
}

// CreateRecipeStep creates a recipe step in the database.
func (q *repository) CreateRecipeStep(ctx context.Context, input *mealplanning.RecipeStepDatabaseCreationInput) (*mealplanning.RecipeStep, error) {
	return q.createRecipeStep(ctx, q.db, input)
}

// UpdateRecipeStep updates a particular recipe step.
func (q *repository) UpdateRecipeStep(ctx context.Context, updated *mealplanning.RecipeStep) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateRecipeStep(ctx, q.db, &generated.UpdateRecipeStepParams{
		ConditionExpression:           updated.ConditionExpression,
		PreparationID:                 updated.Preparation.ID,
		ID:                            updated.ID,
		BelongsToRecipe:               updated.BelongsToRecipe,
		Notes:                         updated.Notes,
		ExplicitInstructions:          updated.ExplicitInstructions,
		MaximumTemperatureInCelsius:   database.NullStringFromFloat32Pointer(updated.TemperatureInCelsius.Max),
		MinimumTemperatureInCelsius:   database.NullStringFromFloat32Pointer(updated.TemperatureInCelsius.Min),
		MaximumEstimatedTimeInSeconds: database.NullInt64FromUint32Pointer(updated.EstimatedTimeInSeconds.Max),
		MinimumEstimatedTimeInSeconds: database.NullInt64FromUint32Pointer(updated.EstimatedTimeInSeconds.Min),
		Index:                         int32(updated.Index),
		Optional:                      updated.Optional,
		StartTimerAutomatically:       updated.StartTimerAutomatically,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step")
	}

	return nil
}

// ArchiveRecipeStep archives a recipe step from the database by its ID.
func (q *repository) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	rowsAffected, err := q.generatedQuerier.ArchiveRecipeStep(ctx, q.db, &generated.ArchiveRecipeStepParams{
		BelongsToRecipe: recipeID,
		ID:              recipeStepID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
