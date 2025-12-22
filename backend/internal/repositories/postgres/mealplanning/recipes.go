package mealplanning

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/recipevalidator"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ mealplanning.RecipeDataManager = (*repository)(nil)
)

// RecipeExists fetches whether a recipe exists from the database.
func (q *repository) RecipeExists(ctx context.Context, recipeID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return false, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	result, err := q.generatedQuerier.CheckRecipeExistence(ctx, q.db, recipeID)
	if err != nil {
		return false, observability.PrepareError(err, span, "performing recipe existence check")
	}

	return result, nil
}

// getRecipe fetches a recipe from the database.
func (q *repository) getRecipe(ctx context.Context, recipeID string) (*mealplanning.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	var x *mealplanning.Recipe
	results, err := q.generatedQuerier.GetRecipeByID(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe")
	}

	for _, result := range results {
		if x == nil {
			x = &mealplanning.Recipe{
				CreatedAt:           result.CreatedAt,
				InspiredByRecipeID:  database.StringPointerFromNullString(result.InspiredByRecipeID),
				LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
				ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
				PluralPortionName:   result.PluralPortionName,
				Description:         result.Description,
				Name:                result.Name,
				PortionName:         result.PortionName,
				ID:                  result.ID,
				CreatedByUser:       result.CreatedByUser,
				Source:              result.Source,
				Slug:                result.Slug,
				YieldsComponentType: string(result.YieldsComponentType),
				EstimatedPortions: types.Float32RangeWithOptionalMax{
					Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
					Min: database.Float32FromString(result.MinEstimatedPortions),
				},
				Status:           string(result.Status),
				EligibleForMeals: result.EligibleForMeals,
			}
		}

		// Only add the step if it actually exists (not NULL from LEFT JOIN)
		if result.RecipeStepID.Valid {
			x.Steps = append(x.Steps, &mealplanning.RecipeStep{
				CreatedAt: result.RecipeStepCreatedAt.Time,
				EstimatedTimeInSeconds: types.OptionalUint32Range{
					Max: database.Uint32PointerFromNullInt64(result.RecipeStepMaximumEstimatedTimeInSeconds),
					Min: database.Uint32PointerFromNullInt64(result.RecipeStepMinimumEstimatedTimeInSeconds),
				},
				TemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.RecipeStepMaximumTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.RecipeStepMinimumTemperatureInCelsius),
				},
				ArchivedAt:           database.TimePointerFromNullTime(result.RecipeStepArchivedAt),
				LastUpdatedAt:        database.TimePointerFromNullTime(result.RecipeStepLastUpdatedAt),
				BelongsToRecipe:      result.RecipeStepBelongsToRecipe.String,
				ConditionExpression:  result.RecipeStepConditionExpression.String,
				ID:                   result.RecipeStepID.String,
				Notes:                result.RecipeStepNotes.String,
				ExplicitInstructions: result.RecipeStepExplicitInstructions.String,
				Preparation: mealplanning.ValidPreparation{
					CreatedAt: result.RecipeStepPreparationCreatedAt.Time,
					InstrumentCount: types.Uint16RangeWithOptionalMax{
						Max: database.Uint16PointerFromNullInt32(result.RecipeStepPreparationMaximumInstrumentCount),
						Min: uint16(result.RecipeStepPreparationMinimumInstrumentCount.Int32),
					},
					IngredientCount: types.Uint16RangeWithOptionalMax{
						Max: database.Uint16PointerFromNullInt32(result.RecipeStepPreparationMaximumIngredientCount),
						Min: uint16(result.RecipeStepPreparationMinimumIngredientCount.Int32),
					},
					VesselCount: types.Uint16RangeWithOptionalMax{
						Max: database.Uint16PointerFromNullInt32(result.RecipeStepPreparationMaximumVesselCount),
						Min: uint16(result.RecipeStepPreparationMinimumVesselCount.Int32),
					},
					ArchivedAt:                  database.TimePointerFromNullTime(result.RecipeStepPreparationArchivedAt),
					LastUpdatedAt:               database.TimePointerFromNullTime(result.RecipeStepPreparationLastUpdatedAt),
					IconPath:                    result.RecipeStepPreparationIconPath.String,
					PastTense:                   result.RecipeStepPreparationPastTense.String,
					ID:                          result.RecipeStepPreparationID.String,
					Name:                        result.RecipeStepPreparationName.String,
					Description:                 result.RecipeStepPreparationDescription.String,
					Slug:                        result.RecipeStepPreparationSlug.String,
					RestrictToIngredients:       result.RecipeStepPreparationRestrictToIngredients.Bool,
					TemperatureRequired:         result.RecipeStepPreparationTemperatureRequired.Bool,
					TimeEstimateRequired:        result.RecipeStepPreparationTimeEstimateRequired.Bool,
					ConditionExpressionRequired: result.RecipeStepPreparationConditionExpressionRequired.Bool,
					ConsumesVessel:              result.RecipeStepPreparationConsumesVessel.Bool,
					OnlyForVessels:              result.RecipeStepPreparationOnlyForVessels.Bool,
					YieldsNothing:               result.RecipeStepPreparationYieldsNothing.Bool,
				},
				Index:                   uint32(result.RecipeStepIndex.Int32),
				Optional:                result.RecipeStepOptional.Bool,
				StartTimerAutomatically: result.RecipeStepStartTimerAutomatically.Bool,
				Media:                   []*mealplanning.RecipeMedia{},
			})
		}
	}

	if x == nil {
		return nil, sql.ErrNoRows
	}

	prepTasks, err := q.getRecipePrepTasksForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step prep tasks for recipe")
	}
	if prepTasks != nil {
		x.PrepTasks = prepTasks
	}

	recipeMedia, err := q.getRecipeMediaForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step media for recipe")
	}
	if recipeMedia != nil {
		x.Media = recipeMedia
	}

	ingredients, err := q.getRecipeStepIngredientsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step ingredients for recipe")
	}

	products, err := q.getRecipeStepProductsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step products for recipe")
	}

	instruments, err := q.getRecipeStepInstrumentsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step instruments for recipe")
	}

	vessels, err := q.getRecipeStepVesselsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step vessels for recipe")
	}

	completionConditions, err := q.getRecipeStepCompletionConditionsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step completion conditions for recipe")
	}

	for i, step := range x.Steps {
		for _, ingredient := range ingredients {
			if ingredient.BelongsToRecipeStep == step.ID {
				x.Steps[i].Ingredients = append(x.Steps[i].Ingredients, ingredient)
			}
		}

		for _, product := range products {
			if product.BelongsToRecipeStep == step.ID {
				x.Steps[i].Products = append(x.Steps[i].Products, product)
			}
		}

		for _, instrument := range instruments {
			if instrument.BelongsToRecipeStep == step.ID {
				x.Steps[i].Instruments = append(x.Steps[i].Instruments, instrument)
			}
		}

		for _, vessel := range vessels {
			if vessel.BelongsToRecipeStep == step.ID {
				x.Steps[i].Vessels = append(x.Steps[i].Vessels, vessel)
			}
		}

		for _, completionCondition := range completionConditions {
			if completionCondition.BelongsToRecipeStep == step.ID {
				x.Steps[i].CompletionConditions = append(x.Steps[i].CompletionConditions, completionCondition)
			}
		}

		recipeMedia, err = q.getRecipeMediaForRecipeStep(ctx, recipeID, step.ID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "fetching recipe media for recipe step")
		}
		x.Steps[i].Media = recipeMedia
	}

	return x, nil
}

// GetRecipe fetches a recipe from the database.
func (q *repository) GetRecipe(ctx context.Context, recipeID string) (*mealplanning.Recipe, error) {
	return q.getRecipe(ctx, recipeID)
}

// GetRecipes fetches a list of recipes from the database that meet a particular filter.
func (q *repository) GetRecipes(ctx context.Context, status string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.Recipe], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if status == "" {
		status = mealplanning.RecipeStatusApproved
	}

	results, err := q.generatedQuerier.GetRecipes(ctx, q.db, &generated.GetRecipesParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:          database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:     database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
		Status:          generated.NullRecipeStatus{RecipeStatus: generated.RecipeStatus(status), Valid: true},
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipes list retrieval query")
	}

	var (
		data          []*mealplanning.Recipe
		filteredCount uint64
		totalCount    uint64
	)
	for _, result := range results {
		data = append(data, &mealplanning.Recipe{
			CreatedAt:           result.CreatedAt,
			InspiredByRecipeID:  database.StringPointerFromNullString(result.InspiredByRecipeID),
			LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
			PluralPortionName:   result.PluralPortionName,
			Description:         result.Description,
			Name:                result.Name,
			PortionName:         result.PortionName,
			ID:                  result.ID,
			CreatedByUser:       result.CreatedByUser,
			Source:              result.Source,
			Slug:                result.Slug,
			YieldsComponentType: string(result.YieldsComponentType),
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
				Min: database.Float32FromString(result.MinEstimatedPortions),
			},
			Status:           string(result.Status),
			EligibleForMeals: result.EligibleForMeals,
		})
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(r *mealplanning.Recipe) string { return r.ID },
		filter,
	)

	return x, nil
}

// GetRecipesCreatedByUser fetches a list of recipes from the database that meet a particular filter.
func (q *repository) GetRecipesCreatedByUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.Recipe], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if userID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	results, err := q.generatedQuerier.GetRecipesCreatedByUser(ctx, q.db, &generated.GetRecipesCreatedByUserParams{
		CreatedByUser:   userID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:          database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:     database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipes list retrieval query")
	}

	var (
		data          []*mealplanning.Recipe
		filteredCount uint64
		totalCount    uint64
	)

	for _, result := range results {
		data = append(data, &mealplanning.Recipe{
			CreatedAt:           result.CreatedAt,
			InspiredByRecipeID:  database.StringPointerFromNullString(result.InspiredByRecipeID),
			LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
			PluralPortionName:   result.PluralPortionName,
			Description:         result.Description,
			Name:                result.Name,
			PortionName:         result.PortionName,
			ID:                  result.ID,
			CreatedByUser:       result.CreatedByUser,
			Source:              result.Source,
			Slug:                result.Slug,
			YieldsComponentType: string(result.YieldsComponentType),
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
				Min: database.Float32FromString(result.MinEstimatedPortions),
			},
			Status:           string(result.Status),
			EligibleForMeals: result.EligibleForMeals,
		})
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(r *mealplanning.Recipe) string { return r.ID },
		filter,
	)

	return x, nil
}

// GetRecipesWithIDs fetches a list of recipes from the database that meet a particular filter.
func (q *repository) GetRecipesWithIDs(ctx context.Context, ids []string) ([]*mealplanning.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if len(ids) == 0 {
		return []*mealplanning.Recipe{}, nil
	}

	results, err := q.generatedQuerier.GetRecipesWithIDs(ctx, q.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipes list retrieval by ids")
	}

	recipesByID := map[string]*mealplanning.Recipe{}
	for _, result := range results {
		r, exists := recipesByID[result.ID]
		if !exists {
			r = &mealplanning.Recipe{
				CreatedAt:           result.CreatedAt,
				InspiredByRecipeID:  database.StringPointerFromNullString(result.InspiredByRecipeID),
				LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
				ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
				PluralPortionName:   result.PluralPortionName,
				Description:         result.Description,
				Name:                result.Name,
				PortionName:         result.PortionName,
				ID:                  result.ID,
				CreatedByUser:       result.CreatedByUser,
				Source:              result.Source,
				Slug:                result.Slug,
				YieldsComponentType: string(result.YieldsComponentType),
				EstimatedPortions: types.Float32RangeWithOptionalMax{
					Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
					Min: database.Float32FromString(result.MinEstimatedPortions),
				},
				Status:           string(result.Status),
				EligibleForMeals: result.EligibleForMeals,
				Steps:            []*mealplanning.RecipeStep{},
			}
			recipesByID[result.ID] = r
		}

		// optional step
		if result.RecipeStepID.Valid && result.RecipeStepID.String != "" {
			var prep mealplanning.ValidPreparation
			if result.RecipeStepPreparationID.Valid {
				ingMin := uint16(0)
				if result.RecipeStepPreparationMinimumIngredientCount.Valid {
					ingMin = uint16(result.RecipeStepPreparationMinimumIngredientCount.Int32)
				}
				instMin := uint16(0)
				if result.RecipeStepPreparationMinimumInstrumentCount.Valid {
					instMin = uint16(result.RecipeStepPreparationMinimumInstrumentCount.Int32)
				}
				vesselMin := uint16(0)
				if result.RecipeStepPreparationMinimumVesselCount.Valid {
					vesselMin = uint16(result.RecipeStepPreparationMinimumVesselCount.Int32)
				}

				prep = mealplanning.ValidPreparation{
					ID:                    result.RecipeStepPreparationID.String,
					Name:                  result.RecipeStepPreparationName.String,
					Slug:                  result.RecipeStepPreparationSlug.String,
					Description:           result.RecipeStepPreparationDescription.String,
					IconPath:              result.RecipeStepPreparationIconPath.String,
					YieldsNothing:         database.BoolFromNullBool(result.RecipeStepPreparationYieldsNothing),
					RestrictToIngredients: database.BoolFromNullBool(result.RecipeStepPreparationRestrictToIngredients),
					PastTense:             result.RecipeStepPreparationPastTense.String,
					IngredientCount: types.Uint16RangeWithOptionalMax{
						Min: ingMin,
						Max: database.Uint16PointerFromNullInt32(result.RecipeStepPreparationMaximumIngredientCount),
					},
					InstrumentCount: types.Uint16RangeWithOptionalMax{
						Min: instMin,
						Max: database.Uint16PointerFromNullInt32(result.RecipeStepPreparationMaximumInstrumentCount),
					},
					TemperatureRequired:         database.BoolFromNullBool(result.RecipeStepPreparationTemperatureRequired),
					TimeEstimateRequired:        database.BoolFromNullBool(result.RecipeStepPreparationTimeEstimateRequired),
					ConditionExpressionRequired: database.BoolFromNullBool(result.RecipeStepPreparationConditionExpressionRequired),
					ConsumesVessel:              database.BoolFromNullBool(result.RecipeStepPreparationConsumesVessel),
					OnlyForVessels:              database.BoolFromNullBool(result.RecipeStepPreparationOnlyForVessels),
					VesselCount: types.Uint16RangeWithOptionalMax{
						Min: vesselMin,
						Max: database.Uint16PointerFromNullInt32(result.RecipeStepPreparationMaximumVesselCount),
					},
					CreatedAt:     database.TimeFromNullTime(result.RecipeStepPreparationCreatedAt),
					LastUpdatedAt: database.TimePointerFromNullTime(result.RecipeStepPreparationLastUpdatedAt),
					ArchivedAt:    database.TimePointerFromNullTime(result.RecipeStepPreparationArchivedAt),
				}
			}

			stepIndex := uint32(0)
			if result.RecipeStepIndex.Valid {
				stepIndex = uint32(result.RecipeStepIndex.Int32)
			}

			r.Steps = append(r.Steps, &mealplanning.RecipeStep{
				ID:              result.RecipeStepID.String,
				BelongsToRecipe: result.RecipeStepBelongsToRecipe.String,
				Index:           stepIndex,
				EstimatedTimeInSeconds: types.OptionalUint32Range{
					Min: database.Uint32PointerFromNullInt64(result.RecipeStepMinimumEstimatedTimeInSeconds),
					Max: database.Uint32PointerFromNullInt64(result.RecipeStepMaximumEstimatedTimeInSeconds),
				},
				TemperatureInCelsius: types.OptionalFloat32Range{
					Min: database.Float32PointerFromNullString(result.RecipeStepMinimumTemperatureInCelsius),
					Max: database.Float32PointerFromNullString(result.RecipeStepMaximumTemperatureInCelsius),
				},
				Notes:                   result.RecipeStepNotes.String,
				ExplicitInstructions:    result.RecipeStepExplicitInstructions.String,
				ConditionExpression:     result.RecipeStepConditionExpression.String,
				Optional:                database.BoolFromNullBool(result.RecipeStepOptional),
				StartTimerAutomatically: database.BoolFromNullBool(result.RecipeStepStartTimerAutomatically),
				CreatedAt:               database.TimeFromNullTime(result.RecipeStepCreatedAt),
				LastUpdatedAt:           database.TimePointerFromNullTime(result.RecipeStepLastUpdatedAt),
				ArchivedAt:              database.TimePointerFromNullTime(result.RecipeStepArchivedAt),
				Preparation:             prep,
			})
		}
	}

	out := make([]*mealplanning.Recipe, 0, len(recipesByID))
	for _, r := range recipesByID {
		out = append(out, r)
	}

	return out, nil
}

// GetRecipeIDsThatNeedSearchIndexing fetches a list of recipe IDs from the database that meet a particular filter.
func (q *repository) GetRecipeIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetRecipesNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing recipes list retrieval query")
	}

	return results, nil
}

// SearchForRecipes fetches a list of recipes from the database that match a query.
func (q *repository) SearchForRecipes(ctx context.Context, recipeNameQuery string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.Recipe], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.RecipeSearch(ctx, q.db, &generated.RecipeSearchParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:          database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:     database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
		Query:           recipeNameQuery,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipes search query")
	}

	var (
		data          []*mealplanning.Recipe
		filteredCount uint64
		totalCount    uint64
	)

	for _, result := range results {
		data = append(data, &mealplanning.Recipe{
			CreatedAt:           result.CreatedAt,
			InspiredByRecipeID:  database.StringPointerFromNullString(result.InspiredByRecipeID),
			LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
			PluralPortionName:   result.PluralPortionName,
			Description:         result.Description,
			Name:                result.Name,
			PortionName:         result.PortionName,
			ID:                  result.ID,
			CreatedByUser:       result.CreatedByUser,
			Source:              result.Source,
			Slug:                result.Slug,
			YieldsComponentType: string(result.YieldsComponentType),
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
				Min: database.Float32FromString(result.MinEstimatedPortions),
			},
			Status:           string(result.Status),
			EligibleForMeals: result.EligibleForMeals,
		})
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(r *mealplanning.Recipe) string { return r.ID },
		filter,
	)

	return x, nil
}

// SearchForMealEligibleRecipes fetches a list of recipes from the database that match a query.
func (q *repository) SearchForMealEligibleRecipes(ctx context.Context, recipeNameQuery string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.Recipe], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.SearchForMealEligibleRecipes(ctx, q.db, &generated.SearchForMealEligibleRecipesParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:          database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:     database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
		Query:           recipeNameQuery,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipes search query")
	}

	var (
		data          []*mealplanning.Recipe
		filteredCount uint64
		totalCount    uint64
	)

	for _, result := range results {
		data = append(data, &mealplanning.Recipe{
			CreatedAt:           result.CreatedAt,
			InspiredByRecipeID:  database.StringPointerFromNullString(result.InspiredByRecipeID),
			LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
			PluralPortionName:   result.PluralPortionName,
			Description:         result.Description,
			Name:                result.Name,
			PortionName:         result.PortionName,
			ID:                  result.ID,
			CreatedByUser:       result.CreatedByUser,
			Source:              result.Source,
			Slug:                result.Slug,
			YieldsComponentType: string(result.YieldsComponentType),
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
				Min: database.Float32FromString(result.MinEstimatedPortions),
			},
			Status:           string(result.Status),
			EligibleForMeals: result.EligibleForMeals,
		})
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(r *mealplanning.Recipe) string { return r.ID },
		filter,
	)

	return x, nil
}

// validateAndPopulateRecipeInput validates bridge table IDs and populates derived fields.
// This is a no-op if no bridge table IDs are present (backward compatible).
func (q *repository) validateAndPopulateRecipeInput(ctx context.Context, input *mealplanning.RecipeDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	// Collect bridge table IDs using helper methods
	vipIDs := input.GetAllValidIngredientPreparationIDs()
	vimuIDs := input.GetAllValidIngredientMeasurementUnitIDs()
	vpiIDs := input.GetAllValidPreparationInstrumentIDs()
	vpvIDs := input.GetAllValidPreparationVesselIDs()

	// Only proceed with validation if any bridge table IDs are present
	if len(vipIDs) == 0 && len(vimuIDs) == 0 && len(vpiIDs) == 0 && len(vpvIDs) == 0 {
		return nil
	}

	// Batch fetch bridge table records
	vipMap, err := q.GetValidIngredientPreparationsByIDs(ctx, vipIDs)
	if err != nil {
		return observability.PrepareError(err, span, "fetching valid ingredient preparations")
	}

	vimuMap, err := q.GetValidIngredientMeasurementUnitsByIDs(ctx, vimuIDs)
	if err != nil {
		return observability.PrepareError(err, span, "fetching valid ingredient measurement units")
	}

	vpiMap, err := q.GetValidPreparationInstrumentsByIDs(ctx, vpiIDs)
	if err != nil {
		return observability.PrepareError(err, span, "fetching valid preparation instruments")
	}

	vpvMap, err := q.GetValidPreparationVesselsByIDs(ctx, vpvIDs)
	if err != nil {
		return observability.PrepareError(err, span, "fetching valid preparation vessels")
	}

	// Create validator and validate/populate the input
	validator := recipevalidator.NewRecipeValidator(vipMap, vimuMap, vpiMap, vpvMap)
	if err = validator.ValidateAndPopulate(input); err != nil {
		return observability.PrepareError(err, span, "validating recipe input")
	}

	return nil
}

// CreateRecipe creates a recipe in the database.
func (q *repository) CreateRecipe(ctx context.Context, input *mealplanning.RecipeDatabaseCreationInput) (*mealplanning.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeIDKey, input.ID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, input.ID)

	// Validate and populate bridge table IDs if any are present
	if err := q.validateAndPopulateRecipeInput(ctx, input); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating recipe input")
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	// create the recipe.
	if err = q.generatedQuerier.CreateRecipe(ctx, tx, &generated.CreateRecipeParams{
		MinEstimatedPortions: database.StringFromFloat32(input.EstimatedPortions.Min),
		ID:                   input.ID,
		Slug:                 input.Slug,
		Source:               input.Source,
		Description:          input.Description,
		CreatedByUser:        input.CreatedByUser,
		Name:                 input.Name,
		YieldsComponentType:  generated.ComponentType(input.YieldsComponentType),
		PortionName:          input.PortionName,
		PluralPortionName:    input.PluralPortionName,
		MaxEstimatedPortions: database.NullStringFromFloat32Pointer(input.EstimatedPortions.Max),
		InspiredByRecipeID:   database.NullStringFromStringPointer(input.InspiredByRecipeID),
		Status:               mealplanning.RecipeStatusSubmitted,
		EligibleForMeals:     input.EligibleForMeals,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe creation query")
	}

	x := &mealplanning.Recipe{
		ID:                 input.ID,
		Name:               input.Name,
		Slug:               input.Slug,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		CreatedByUser:      input.CreatedByUser,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Max: input.EstimatedPortions.Max,
			Min: input.EstimatedPortions.Min,
		},
		Status:              mealplanning.RecipeStatusSubmitted,
		EligibleForMeals:    input.EligibleForMeals,
		PortionName:         input.PortionName,
		PluralPortionName:   input.PluralPortionName,
		YieldsComponentType: input.YieldsComponentType,
		CreatedAt:           q.CurrentTime(),
		PrepTasks:           []*mealplanning.RecipePrepTask{},
		Steps:               []*mealplanning.RecipeStep{},
		Media:               []*mealplanning.RecipeMedia{},
	}

	findCreatedRecipeStepProductsForIngredients(input)
	findCreatedRecipeStepProductsForInstruments(input)
	findCreatedRecipeStepProductsForVessels(input)

	for i, stepInput := range input.Steps {
		stepInput.Index = uint32(i)
		stepInput.BelongsToRecipe = x.ID

		q.logger.Info(fmt.Sprintf("creating recipe step #%d", i+1))

		var s *mealplanning.RecipeStep
		s, err = q.createRecipeStep(ctx, tx, stepInput)
		if err != nil {
			q.RollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(err, span, "creating recipe step #%d", i+1)
		}

		x.Steps = append(x.Steps, s)
	}

	for i, prepTaskInput := range input.PrepTasks {
		var pt *mealplanning.RecipePrepTask
		pt, err = q.createRecipePrepTask(ctx, tx, prepTaskInput)
		if err != nil {
			q.RollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(err, span, "creating recipe prep task #%d", i+1)
		}

		x.PrepTasks = append(x.PrepTasks, pt)
	}

	for i, m := range input.Media {
		var rm *mealplanning.RecipeMedia
		rm, err = q.CreateRecipeMedia(ctx, &mealplanning.RecipeMediaDatabaseCreationInput{
			ID:                  m.ID,
			BelongsToRecipe:     m.BelongsToRecipe,
			BelongsToRecipeStep: m.BelongsToRecipeStep,
			MimeType:            m.MimeType,
			InternalPath:        m.InternalPath,
			ExternalPath:        m.ExternalPath,
			Index:               m.Index,
		})
		if err != nil {
			q.RollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(err, span, "creating recipe media #%d", i+1)
		}

		x.Media = append(x.Media, rm)
	}

	if input.AlsoCreateMeal {
		if _, err = q.createMeal(ctx, tx, &mealplanning.MealDatabaseCreationInput{
			ID:          identifiers.New(),
			Name:        x.Name,
			Description: x.Description,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: x.EstimatedPortions.Min,
				Max: x.EstimatedPortions.Max,
			},
			EligibleForMealPlans: x.EligibleForMeals,
			CreatedByUser:        x.CreatedByUser,
			Components: []*mealplanning.MealComponentDatabaseCreationInput{
				{
					RecipeID:      x.ID,
					RecipeScale:   1.0,
					ComponentType: mealplanning.MealComponentTypesMain,
				},
			},
		}); err != nil {
			q.RollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(err, span, "creating meal from recipe")
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe created")

	return x, nil
}

func findCreatedRecipeStepProductsForIngredients(recipe *mealplanning.RecipeDatabaseCreationInput) {
	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.ProductOfRecipeStepIndex != nil && ingredient.ProductOfRecipeStepProductIndex != nil {
				enoughSteps := len(recipe.Steps) > int(*ingredient.ProductOfRecipeStepIndex)
				enoughRecipeStepProducts := len(recipe.Steps[int(*ingredient.ProductOfRecipeStepIndex)].Products) > int(*ingredient.ProductOfRecipeStepProductIndex)
				relevantProductIsIngredient := recipe.Steps[*ingredient.ProductOfRecipeStepIndex].Products[*ingredient.ProductOfRecipeStepProductIndex].Type == mealplanning.RecipeStepProductIngredientType
				if enoughSteps && enoughRecipeStepProducts && relevantProductIsIngredient {
					product := recipe.Steps[*ingredient.ProductOfRecipeStepIndex].Products[*ingredient.ProductOfRecipeStepProductIndex]
					ingredient.RecipeStepProductID = &product.ID
					// Inherit measurement unit from the product if not already set
					if ingredient.MeasurementUnitID == "" && product.MeasurementUnitID != nil {
						ingredient.MeasurementUnitID = *product.MeasurementUnitID
					}
				}
			}
		}
	}
}

func findCreatedRecipeStepProductsForInstruments(recipe *mealplanning.RecipeDatabaseCreationInput) {
	for _, step := range recipe.Steps {
		for _, instrument := range step.Instruments {
			if instrument.ProductOfRecipeStepIndex != nil && instrument.ProductOfRecipeStepProductIndex != nil {
				enoughSteps := len(recipe.Steps) > int(*instrument.ProductOfRecipeStepIndex)
				enoughRecipeStepProducts := len(recipe.Steps[int(*instrument.ProductOfRecipeStepIndex)].Products) > int(*instrument.ProductOfRecipeStepProductIndex)
				relevantProductIsInstrument := recipe.Steps[*instrument.ProductOfRecipeStepIndex].Products[*instrument.ProductOfRecipeStepProductIndex].Type == mealplanning.RecipeStepProductInstrumentType
				if enoughSteps && enoughRecipeStepProducts && relevantProductIsInstrument {
					instrument.RecipeStepProductID = &recipe.Steps[*instrument.ProductOfRecipeStepIndex].Products[*instrument.ProductOfRecipeStepProductIndex].ID
				}
			}
		}
	}
}

func findCreatedRecipeStepProductsForVessels(recipe *mealplanning.RecipeDatabaseCreationInput) {
	for _, step := range recipe.Steps {
		for _, vessel := range step.Vessels {
			if vessel.ProductOfRecipeStepIndex != nil && vessel.ProductOfRecipeStepProductIndex != nil {
				enoughSteps := len(recipe.Steps) > int(*vessel.ProductOfRecipeStepIndex)
				enoughRecipeStepProducts := len(recipe.Steps[int(*vessel.ProductOfRecipeStepIndex)].Products) > int(*vessel.ProductOfRecipeStepProductIndex)
				relevantProductIsVessel := recipe.Steps[*vessel.ProductOfRecipeStepIndex].Products[*vessel.ProductOfRecipeStepProductIndex].Type == mealplanning.RecipeStepProductVesselType
				if enoughSteps && enoughRecipeStepProducts && relevantProductIsVessel {
					vessel.RecipeStepProductID = &recipe.Steps[*vessel.ProductOfRecipeStepIndex].Products[*vessel.ProductOfRecipeStepProductIndex].ID
				} else {
					log.Printf("for recipe step id %q, vessel ID %q, not enough steps: %t, not enough recipe step products: %t, relevant product is vessel: %t", step.ID, vessel.ID, enoughSteps, enoughRecipeStepProducts, relevantProductIsVessel)
				}
			}
		}
	}
}

// UpdateRecipe updates a particular recipe.
func (q *repository) UpdateRecipe(ctx context.Context, updated *mealplanning.Recipe) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.UserIDKey, updated.CreatedByUser)

	if _, err := q.generatedQuerier.UpdateRecipe(ctx, q.db, &generated.UpdateRecipeParams{
		Name:                 updated.Name,
		Slug:                 updated.Slug,
		Source:               updated.Source,
		Description:          updated.Description,
		InspiredByRecipeID:   database.NullStringFromStringPointer(updated.InspiredByRecipeID),
		MinEstimatedPortions: database.StringFromFloat32(updated.EstimatedPortions.Min),
		MaxEstimatedPortions: database.NullStringFromFloat32Pointer(updated.EstimatedPortions.Max),
		PortionName:          updated.PortionName,
		PluralPortionName:    updated.PluralPortionName,
		EligibleForMeals:     updated.EligibleForMeals,
		YieldsComponentType:  generated.ComponentType(updated.YieldsComponentType),
		CreatedByUser:        updated.CreatedByUser,
		ID:                   updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe")
	}

	logger.Info("recipe updated")

	return nil
}

// UpdateRecipeStatus updates a particular recipe's status exclusively.
func (q *repository) UpdateRecipeStatus(ctx context.Context, recipeID, newStatus string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithSpan(span)

	if recipeID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if _, err := q.generatedQuerier.UpdateRecipeStatus(ctx, q.db, &generated.UpdateRecipeStatusParams{
		Status: generated.RecipeStatus(newStatus),
		ID:     recipeID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe status")
	}

	return nil
}

// MarkRecipeAsIndexed updates a particular recipe's last_indexed_at value.
func (q *repository) MarkRecipeAsIndexed(ctx context.Context, recipeID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if _, err := q.generatedQuerier.UpdateRecipeLastIndexedAt(ctx, q.db, recipeID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking recipe as indexed")
	}

	logger.Info("recipe marked as indexed")

	return nil
}

// ArchiveRecipe archives a recipe from the database by its ID.
func (q *repository) ArchiveRecipe(ctx context.Context, recipeID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return database.ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if userID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	rowsAffected, err := q.generatedQuerier.ArchiveRecipe(ctx, q.db, &generated.ArchiveRecipeParams{
		CreatedByUser: userID,
		ID:            recipeID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
