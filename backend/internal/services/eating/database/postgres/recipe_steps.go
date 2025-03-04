package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

var (
	_ types.RecipeStepDataManager = (*Querier)(nil)
)

// RecipeStepExists fetches whether a recipe step exists from the database.
func (q *Querier) RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return false, ErrInvalidIDProvided
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
func (q *Querier) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
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

	recipeStep := &types.RecipeStep{
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
		Media:                nil,
		Products:             nil,
		Instruments:          nil,
		Vessels:              nil,
		CompletionConditions: nil,
		Ingredients:          nil,
		Preparation: types.ValidPreparation{
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

	return recipeStep, nil
}

// getRecipeStepByID fetches a recipe step from the database.
func (q *Querier) getRecipeStepByID(ctx context.Context, querier database.SQLQueryExecutor, recipeStepID string) (*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	result, err := q.generatedQuerier.GetRecipeStepByRecipeID(ctx, querier, recipeStepID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe step")
	}

	recipeStep := &types.RecipeStep{
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
		Media:                nil,
		Products:             nil,
		Instruments:          nil,
		Vessels:              nil,
		CompletionConditions: nil,
		Ingredients:          nil,
		Preparation: types.ValidPreparation{
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

	return recipeStep, nil
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter.
func (q *Querier) GetRecipeSteps(ctx context.Context, recipeID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.RecipeStep], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.RecipeStep]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeSteps(ctx, q.db, &generated.GetRecipeStepsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
		RecipeID:        recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe steps")
	}

	for _, result := range results {
		recipeStep := &types.RecipeStep{
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
			Media:                nil,
			Products:             nil,
			Instruments:          nil,
			Vessels:              nil,
			CompletionConditions: nil,
			Ingredients:          nil,
			Preparation: types.ValidPreparation{
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

	return x, nil
}

// CreateRecipeStep creates a recipe step in the database.
func (q *Querier) createRecipeStep(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
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

	x := &types.RecipeStep{
		ID:                      input.ID,
		Index:                   input.Index,
		Preparation:             types.ValidPreparation{ID: input.PreparationID},
		EstimatedTimeInSeconds:  input.EstimatedTimeInSeconds,
		TemperatureInCelsius:    input.TemperatureInCelsius,
		Notes:                   input.Notes,
		ExplicitInstructions:    input.ExplicitInstructions,
		ConditionExpression:     input.ConditionExpression,
		Optional:                input.Optional,
		BelongsToRecipe:         input.BelongsToRecipe,
		StartTimerAutomatically: input.StartTimerAutomatically,
		CreatedAt:               q.currentTime(),
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
func (q *Querier) CreateRecipeStep(ctx context.Context, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	return q.createRecipeStep(ctx, q.db, input)
}

// UpdateRecipeStep updates a particular recipe step.
func (q *Querier) UpdateRecipeStep(ctx context.Context, updated *types.RecipeStep) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
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

	logger.Info("recipe step updated")

	return nil
}

// ArchiveRecipeStep archives a recipe step from the database by its ID.
func (q *Querier) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if _, err := q.generatedQuerier.ArchiveRecipeStep(ctx, q.db, &generated.ArchiveRecipeStepParams{
		BelongsToRecipe: recipeID,
		ID:              recipeStepID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step")
	}

	logger.Info("recipe step archived")

	return nil
}
