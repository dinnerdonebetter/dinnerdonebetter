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
	_ mealplanning.RecipeStepIngredientDataManager = (*repository)(nil)
)

// RecipeStepIngredientExists fetches whether a recipe step ingredient exists from the database.
func (r *repository) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (exists bool, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

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

	if recipeStepIngredientID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

	result, err := r.generatedQuerier.CheckRecipeStepIngredientExistence(ctx, r.db, &generated.CheckRecipeStepIngredientExistenceParams{
		RecipeStepID:           recipeStepID,
		RecipeStepIngredientID: recipeStepIngredientID,
		RecipeID:               recipeID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step ingredient existence check")
	}

	return result, nil
}

// GetRecipeStepIngredient fetches a recipe step ingredient from the database.
func (r *repository) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*mealplanning.RecipeStepIngredient, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

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

	if recipeStepIngredientID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

	result, err := r.generatedQuerier.GetRecipeStepIngredient(ctx, r.db, &generated.GetRecipeStepIngredientParams{
		RecipeStepID:           recipeStepID,
		RecipeStepIngredientID: recipeStepIngredientID,
		RecipeID:               recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step ingredient")
	}

	recipeStepIngredient := &mealplanning.RecipeStepIngredient{
		CreatedAt:                 result.CreatedAt,
		RecipeStepProductID:       database.StringPointerFromNullString(result.RecipeStepProductID),
		ArchivedAt:                database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:             database.TimePointerFromNullTime(result.LastUpdatedAt),
		VesselIndex:               database.Uint16PointerFromNullInt32(result.VesselIndex),
		ProductPercentageToUse:    database.Float32PointerFromNullString(result.ProductPercentageToUse),
		RecipeStepProductRecipeID: database.StringPointerFromNullString(result.RecipeStepProductRecipeID),
		QuantityNotes:             result.QuantityNotes,
		ID:                        result.ID,
		BelongsToRecipeStep:       result.BelongsToRecipeStep,
		IngredientNotes:           result.IngredientNotes,
		Name:                      result.Name,
		MeasurementUnit: mealplanning.ValidMeasurementUnit{
			CreatedAt:     result.ValidMeasurementUnitCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
			Name:          result.ValidMeasurementUnitName,
			IconPath:      result.ValidMeasurementUnitIconPath,
			ID:            result.ValidMeasurementUnitID,
			Description:   result.ValidMeasurementUnitDescription,
			PluralName:    result.ValidMeasurementUnitPluralName,
			Slug:          result.ValidMeasurementUnitSlug,
			Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
			Universal:     result.ValidMeasurementUnitUniversal,
			Metric:        result.ValidMeasurementUnitMetric,
			Imperial:      result.ValidMeasurementUnitImperial,
		},
		Quantity: types.Float32RangeWithOptionalMax{
			Max: database.Float32PointerFromNullString(result.MaximumQuantityValue),
			Min: database.Float32FromString(result.MinimumQuantityValue),
		},
		OptionIndex: uint16(result.OptionIndex),
		Optional:    result.Optional,
		ToTaste:     result.ToTaste,
	}

	if result.ValidIngredientID.Valid && result.ValidIngredientID.String != "" {
		recipeStepIngredient.Ingredient = &mealplanning.ValidIngredient{
			CreatedAt:     result.ValidIngredientCreatedAt.Time,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
			},
			IconPath:               result.ValidIngredientIconPath.String,
			Warning:                result.ValidIngredientWarning.String,
			PluralName:             result.ValidIngredientPluralName.String,
			StorageInstructions:    result.ValidIngredientStorageInstructions.String,
			Name:                   result.ValidIngredientName.String,
			ID:                     result.ValidIngredientID.String,
			Description:            result.ValidIngredientDescription.String,
			Slug:                   result.ValidIngredientSlug.String,
			ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions.String,
			ContainsShellfish:      result.ValidIngredientContainsShellfish.Bool,
			IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
			ContainsPeanut:         result.ValidIngredientContainsPeanut.Bool,
			ContainsTreeNut:        result.ValidIngredientContainsTreeNut.Bool,
			ContainsEgg:            result.ValidIngredientContainsEgg.Bool,
			ContainsWheat:          result.ValidIngredientContainsWheat.Bool,
			ContainsSoy:            result.ValidIngredientContainsSoy.Bool,
			AnimalDerived:          result.ValidIngredientAnimalDerived.Bool,
			RestrictToPreparations: result.ValidIngredientRestrictToPreparations.Bool,
			ContainsSesame:         result.ValidIngredientContainsSesame.Bool,
			ContainsFish:           result.ValidIngredientContainsFish.Bool,
			ContainsGluten:         result.ValidIngredientContainsGluten.Bool,
			ContainsDairy:          result.ValidIngredientContainsDairy.Bool,
			ContainsAlcohol:        result.ValidIngredientContainsAlcohol.Bool,
			AnimalFlesh:            result.ValidIngredientAnimalFlesh.Bool,
			IsStarch:               result.ValidIngredientIsStarch.Bool,
			IsProtein:              result.ValidIngredientIsProtein.Bool,
			IsGrain:                result.ValidIngredientIsGrain.Bool,
			IsFruit:                result.ValidIngredientIsFruit.Bool,
			IsSalt:                 result.ValidIngredientIsSalt.Bool,
			IsFat:                  result.ValidIngredientIsFat.Bool,
			IsAcid:                 result.ValidIngredientIsAcid.Bool,
			IsHeat:                 result.ValidIngredientIsHeat.Bool,
		}
	}

	return recipeStepIngredient, nil
}

// getRecipeStepIngredientsForRecipe fetches a list of recipe step ingredients from the database that meet a particular filter.
func (r *repository) getRecipeStepIngredientsForRecipe(ctx context.Context, recipeID string) ([]*mealplanning.RecipeStepIngredient, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	results, err := r.generatedQuerier.GetAllRecipeStepIngredientsForRecipe(ctx, r.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	recipeStepIngredients := []*mealplanning.RecipeStepIngredient{}
	for _, result := range results {
		recipeStepIngredient := &mealplanning.RecipeStepIngredient{
			CreatedAt:                 result.CreatedAt,
			RecipeStepProductID:       database.StringPointerFromNullString(result.RecipeStepProductID),
			ArchivedAt:                database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:             database.TimePointerFromNullTime(result.LastUpdatedAt),
			VesselIndex:               database.Uint16PointerFromNullInt32(result.VesselIndex),
			ProductPercentageToUse:    database.Float32PointerFromNullString(result.ProductPercentageToUse),
			RecipeStepProductRecipeID: database.StringPointerFromNullString(result.RecipeStepProductRecipeID),
			QuantityNotes:             result.QuantityNotes,
			ID:                        result.ID,
			BelongsToRecipeStep:       result.BelongsToRecipeStep,
			IngredientNotes:           result.IngredientNotes,
			Name:                      result.Name,
			MeasurementUnit: mealplanning.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName,
				IconPath:      result.ValidMeasurementUnitIconPath,
				ID:            result.ValidMeasurementUnitID,
				Description:   result.ValidMeasurementUnitDescription,
				PluralName:    result.ValidMeasurementUnitPluralName,
				Slug:          result.ValidMeasurementUnitSlug,
				Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal,
				Metric:        result.ValidMeasurementUnitMetric,
				Imperial:      result.ValidMeasurementUnitImperial,
			},
			Quantity: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.MaximumQuantityValue),
				Min: database.Float32FromString(result.MinimumQuantityValue),
			},
			OptionIndex: uint16(result.OptionIndex),
			Optional:    result.Optional,
			ToTaste:     result.ToTaste,
		}

		if result.ValidIngredientID.Valid {
			recipeStepIngredient.Ingredient = &mealplanning.ValidIngredient{
				CreatedAt:     result.ValidIngredientCreatedAt.Time,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				StorageTemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				},
				IconPath:               result.ValidIngredientIconPath.String,
				Warning:                result.ValidIngredientWarning.String,
				PluralName:             result.ValidIngredientPluralName.String,
				StorageInstructions:    result.ValidIngredientStorageInstructions.String,
				Name:                   result.ValidIngredientName.String,
				ID:                     result.ValidIngredientID.String,
				Description:            result.ValidIngredientDescription.String,
				Slug:                   result.ValidIngredientSlug.String,
				ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions.String,
				ContainsShellfish:      result.ValidIngredientContainsShellfish.Bool,
				IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:         result.ValidIngredientContainsPeanut.Bool,
				ContainsTreeNut:        result.ValidIngredientContainsTreeNut.Bool,
				ContainsEgg:            result.ValidIngredientContainsEgg.Bool,
				ContainsWheat:          result.ValidIngredientContainsWheat.Bool,
				ContainsSoy:            result.ValidIngredientContainsSoy.Bool,
				AnimalDerived:          result.ValidIngredientAnimalDerived.Bool,
				RestrictToPreparations: result.ValidIngredientRestrictToPreparations.Bool,
				ContainsSesame:         result.ValidIngredientContainsSesame.Bool,
				ContainsFish:           result.ValidIngredientContainsFish.Bool,
				ContainsGluten:         result.ValidIngredientContainsGluten.Bool,
				ContainsDairy:          result.ValidIngredientContainsDairy.Bool,
				ContainsAlcohol:        result.ValidIngredientContainsAlcohol.Bool,
				AnimalFlesh:            result.ValidIngredientAnimalFlesh.Bool,
				IsStarch:               result.ValidIngredientIsStarch.Bool,
				IsProtein:              result.ValidIngredientIsProtein.Bool,
				IsGrain:                result.ValidIngredientIsGrain.Bool,
				IsFruit:                result.ValidIngredientIsFruit.Bool,
				IsSalt:                 result.ValidIngredientIsSalt.Bool,
				IsFat:                  result.ValidIngredientIsFat.Bool,
				IsAcid:                 result.ValidIngredientIsAcid.Bool,
				IsHeat:                 result.ValidIngredientIsHeat.Bool,
			}
		}

		recipeStepIngredients = append(recipeStepIngredients, recipeStepIngredient)
	}

	return recipeStepIngredients, nil
}

// GetRecipeStepIngredients fetches a list of recipe step ingredients from the database that meet a particular filter.
func (r *repository) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.RecipeStepIngredient], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

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

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[mealplanning.RecipeStepIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetRecipeStepIngredients(ctx, r.db, &generated.GetRecipeStepIngredientsParams{
		RecipeID:        recipeID,
		RecipeStepID:    recipeStepID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	for _, result := range results {
		recipeStepIngredient := &mealplanning.RecipeStepIngredient{
			CreatedAt:                 result.CreatedAt,
			RecipeStepProductID:       database.StringPointerFromNullString(result.RecipeStepProductID),
			ArchivedAt:                database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:             database.TimePointerFromNullTime(result.LastUpdatedAt),
			VesselIndex:               database.Uint16PointerFromNullInt32(result.VesselIndex),
			ProductPercentageToUse:    database.Float32PointerFromNullString(result.ProductPercentageToUse),
			RecipeStepProductRecipeID: database.StringPointerFromNullString(result.RecipeStepProductRecipeID),
			QuantityNotes:             result.QuantityNotes,
			ID:                        result.ID,
			BelongsToRecipeStep:       result.BelongsToRecipeStep,
			IngredientNotes:           result.IngredientNotes,
			Name:                      result.Name,
			MeasurementUnit: mealplanning.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName,
				IconPath:      result.ValidMeasurementUnitIconPath,
				ID:            result.ValidMeasurementUnitID,
				Description:   result.ValidMeasurementUnitDescription,
				PluralName:    result.ValidMeasurementUnitPluralName,
				Slug:          result.ValidMeasurementUnitSlug,
				Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal,
				Metric:        result.ValidMeasurementUnitMetric,
				Imperial:      result.ValidMeasurementUnitImperial,
			},
			Quantity: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.MaximumQuantityValue),
				Min: database.Float32FromString(result.MinimumQuantityValue),
			},
			OptionIndex: uint16(result.OptionIndex),
			Optional:    result.Optional,
			ToTaste:     result.ToTaste,
		}

		if result.ValidIngredientID.Valid {
			recipeStepIngredient.Ingredient = &mealplanning.ValidIngredient{
				CreatedAt:     result.ValidIngredientCreatedAt.Time,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				StorageTemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				},
				IconPath:               result.ValidIngredientIconPath.String,
				Warning:                result.ValidIngredientWarning.String,
				PluralName:             result.ValidIngredientPluralName.String,
				StorageInstructions:    result.ValidIngredientStorageInstructions.String,
				Name:                   result.ValidIngredientName.String,
				ID:                     result.ValidIngredientID.String,
				Description:            result.ValidIngredientDescription.String,
				Slug:                   result.ValidIngredientSlug.String,
				ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions.String,
				ContainsShellfish:      result.ValidIngredientContainsShellfish.Bool,
				IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:         result.ValidIngredientContainsPeanut.Bool,
				ContainsTreeNut:        result.ValidIngredientContainsTreeNut.Bool,
				ContainsEgg:            result.ValidIngredientContainsEgg.Bool,
				ContainsWheat:          result.ValidIngredientContainsWheat.Bool,
				ContainsSoy:            result.ValidIngredientContainsSoy.Bool,
				AnimalDerived:          result.ValidIngredientAnimalDerived.Bool,
				RestrictToPreparations: result.ValidIngredientRestrictToPreparations.Bool,
				ContainsSesame:         result.ValidIngredientContainsSesame.Bool,
				ContainsFish:           result.ValidIngredientContainsFish.Bool,
				ContainsGluten:         result.ValidIngredientContainsGluten.Bool,
				ContainsDairy:          result.ValidIngredientContainsDairy.Bool,
				ContainsAlcohol:        result.ValidIngredientContainsAlcohol.Bool,
				AnimalFlesh:            result.ValidIngredientAnimalFlesh.Bool,
				IsStarch:               result.ValidIngredientIsStarch.Bool,
				IsProtein:              result.ValidIngredientIsProtein.Bool,
				IsGrain:                result.ValidIngredientIsGrain.Bool,
				IsFruit:                result.ValidIngredientIsFruit.Bool,
				IsSalt:                 result.ValidIngredientIsSalt.Bool,
				IsFat:                  result.ValidIngredientIsFat.Bool,
				IsAcid:                 result.ValidIngredientIsAcid.Bool,
				IsHeat:                 result.ValidIngredientIsHeat.Bool,
			}
		}

		x.Data = append(x.Data, recipeStepIngredient)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// createRecipeStepIngredient creates a recipe step ingredient in the database.
func (r *repository) createRecipeStepIngredient(ctx context.Context, db database.SQLQueryExecutor, input *mealplanning.RecipeStepIngredientDatabaseCreationInput) (*mealplanning.RecipeStepIngredient, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	// create the recipe step ingredient.
	if err := r.generatedQuerier.CreateRecipeStepIngredient(ctx, db, &generated.CreateRecipeStepIngredientParams{
		QuantityNotes:             input.QuantityNotes,
		Name:                      input.Name,
		BelongsToRecipeStep:       input.BelongsToRecipeStep,
		IngredientNotes:           input.IngredientNotes,
		ID:                        input.ID,
		MinimumQuantityValue:      database.StringFromFloat32(input.Quantity.Min),
		RecipeStepProductID:       database.NullStringFromStringPointer(input.RecipeStepProductID),
		MaximumQuantityValue:      database.NullStringFromFloat32Pointer(input.Quantity.Max),
		MeasurementUnit:           database.NullStringFromString(input.MeasurementUnitID),
		IngredientID:              database.NullStringFromStringPointer(input.IngredientID),
		ProductPercentageToUse:    database.NullStringFromFloat32Pointer(input.ProductPercentageToUse),
		RecipeStepProductRecipeID: database.NullStringFromStringPointer(input.RecipeStepProductRecipeID),
		VesselIndex:               database.NullInt32FromUint16Pointer(input.VesselIndex),
		OptionIndex:               int32(input.OptionIndex),
		ToTaste:                   input.ToTaste,
		Optional:                  input.Optional,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step ingredient creation query")
	}

	x := &mealplanning.RecipeStepIngredient{
		ID:              input.ID,
		Name:            input.Name,
		Optional:        input.Optional,
		MeasurementUnit: mealplanning.ValidMeasurementUnit{ID: input.MeasurementUnitID},
		Quantity: types.Float32RangeWithOptionalMax{
			Max: input.Quantity.Max,
			Min: input.Quantity.Min,
		},
		QuantityNotes:             input.QuantityNotes,
		IngredientNotes:           input.IngredientNotes,
		BelongsToRecipeStep:       input.BelongsToRecipeStep,
		RecipeStepProductID:       input.RecipeStepProductID,
		OptionIndex:               input.OptionIndex,
		ToTaste:                   input.ToTaste,
		ProductPercentageToUse:    input.ProductPercentageToUse,
		VesselIndex:               input.VesselIndex,
		RecipeStepProductRecipeID: input.RecipeStepProductRecipeID,
		CreatedAt:                 r.CurrentTime(),
	}

	if input.IngredientID != nil {
		x.Ingredient = &mealplanning.ValidIngredient{ID: *input.IngredientID}
	}

	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, x.ID)

	return x, nil
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database.
func (r *repository) CreateRecipeStepIngredient(ctx context.Context, input *mealplanning.RecipeStepIngredientDatabaseCreationInput) (*mealplanning.RecipeStepIngredient, error) {
	return r.createRecipeStepIngredient(ctx, r.db, input)
}

// UpdateRecipeStepIngredient updates a particular recipe step ingredient.
func (r *repository) UpdateRecipeStepIngredient(ctx context.Context, updated *mealplanning.RecipeStepIngredient) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.RecipeStepIngredientIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, updated.ID)

	var ingredientID *string
	if updated.Ingredient != nil {
		ingredientID = &updated.Ingredient.ID
	}

	if _, err := r.generatedQuerier.UpdateRecipeStepIngredient(ctx, r.db, &generated.UpdateRecipeStepIngredientParams{
		IngredientID:              database.NullStringFromStringPointer(ingredientID),
		Name:                      updated.Name,
		Optional:                  updated.Optional,
		MeasurementUnit:           database.NullStringFromString(updated.MeasurementUnit.ID),
		MinimumQuantityValue:      database.StringFromFloat32(updated.Quantity.Min),
		MaximumQuantityValue:      database.NullStringFromFloat32Pointer(updated.Quantity.Max),
		QuantityNotes:             updated.QuantityNotes,
		RecipeStepProductID:       database.NullStringFromStringPointer(updated.RecipeStepProductID),
		IngredientNotes:           updated.IngredientNotes,
		OptionIndex:               int32(updated.OptionIndex),
		ToTaste:                   updated.ToTaste,
		ProductPercentageToUse:    database.NullStringFromFloat32Pointer(updated.ProductPercentageToUse),
		VesselIndex:               database.NullInt32FromUint16Pointer(updated.VesselIndex),
		RecipeStepProductRecipeID: database.NullStringFromStringPointer(updated.RecipeStepProductRecipeID),
		BelongsToRecipeStep:       updated.BelongsToRecipeStep,
		ID:                        updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step ingredient")
	}

	logger.Info("recipe step ingredient updated")

	return nil
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient from the database by its ID.
func (r *repository) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeStepID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepIngredientID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

	rowsAffected, err := r.generatedQuerier.ArchiveRecipeStepIngredient(ctx, r.db, &generated.ArchiveRecipeStepIngredientParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepIngredientID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step ingredient")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
