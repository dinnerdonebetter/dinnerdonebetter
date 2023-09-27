package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.RecipeStepIngredientDataManager = (*Querier)(nil)
)

// RecipeStepIngredientExists fetches whether a recipe step ingredient exists from the database.
func (q *Querier) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepIngredientID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	result, err := q.generatedQuerier.CheckRecipeStepIngredientExistence(ctx, q.db, &generated.CheckRecipeStepIngredientExistenceParams{
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
func (q *Querier) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	result, err := q.generatedQuerier.GetRecipeStepIngredient(ctx, q.db, &generated.GetRecipeStepIngredientParams{
		RecipeStepID:           recipeStepID,
		RecipeStepIngredientID: recipeStepIngredientID,
		RecipeID:               recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step ingredient")
	}

	recipeStepIngredient := &types.RecipeStepIngredient{
		CreatedAt:                 result.CreatedAt,
		RecipeStepProductID:       stringPointerFromNullString(result.RecipeStepProductID),
		ArchivedAt:                timePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:             timePointerFromNullTime(result.LastUpdatedAt),
		MaximumQuantity:           float32PointerFromNullString(result.MaximumQuantityValue),
		VesselIndex:               uint16PointerFromNullInt32(result.VesselIndex),
		ProductPercentageToUse:    float32PointerFromNullString(result.ProductPercentageToUse),
		RecipeStepProductRecipeID: stringPointerFromNullString(result.RecipeStepProductRecipeID),
		QuantityNotes:             result.QuantityNotes,
		ID:                        result.ID,
		BelongsToRecipeStep:       result.BelongsToRecipeStep,
		IngredientNotes:           result.IngredientNotes,
		Name:                      result.Name,
		MeasurementUnit: types.ValidMeasurementUnit{
			CreatedAt:     result.ValidMeasurementUnitCreatedAt,
			LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
			Name:          result.ValidMeasurementUnitName,
			IconPath:      result.ValidMeasurementUnitIconPath,
			ID:            result.ValidMeasurementUnitID,
			Description:   result.ValidMeasurementUnitDescription,
			PluralName:    result.ValidMeasurementUnitPluralName,
			Slug:          result.ValidMeasurementUnitSlug,
			Volumetric:    boolFromNullBool(result.ValidMeasurementUnitVolumetric),
			Universal:     result.ValidMeasurementUnitUniversal,
			Metric:        result.ValidMeasurementUnitMetric,
			Imperial:      result.ValidMeasurementUnitImperial,
		},
		MinimumQuantity: float32FromString(result.MinimumQuantityValue),
		OptionIndex:     uint16(result.OptionIndex),
		Optional:        result.Optional,
		ToTaste:         result.ToTaste,
	}

	if result.ValidIngredientID.Valid && result.ValidIngredientID.String != "" {
		recipeStepIngredient.Ingredient = &types.ValidIngredient{
			CreatedAt:                               result.ValidIngredientCreatedAt.Time,
			LastUpdatedAt:                           timePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
			ArchivedAt:                              timePointerFromNullTime(result.ValidIngredientArchivedAt),
			MaximumIdealStorageTemperatureInCelsius: float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
			MinimumIdealStorageTemperatureInCelsius: float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
			IconPath:                                result.ValidIngredientIconPath.String,
			Warning:                                 result.ValidIngredientWarning.String,
			PluralName:                              result.ValidIngredientPluralName.String,
			StorageInstructions:                     result.ValidIngredientStorageInstructions.String,
			Name:                                    result.ValidIngredientName.String,
			ID:                                      result.ValidIngredientID.String,
			Description:                             result.ValidIngredientDescription.String,
			Slug:                                    result.ValidIngredientSlug.String,
			ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions.String,
			ContainsShellfish:                       result.ValidIngredientContainsShellfish.Bool,
			IsMeasuredVolumetrically:                result.ValidIngredientVolumetric.Bool,
			IsLiquid:                                boolFromNullBool(result.ValidIngredientIsLiquid),
			ContainsPeanut:                          result.ValidIngredientContainsPeanut.Bool,
			ContainsTreeNut:                         result.ValidIngredientContainsTreeNut.Bool,
			ContainsEgg:                             result.ValidIngredientContainsEgg.Bool,
			ContainsWheat:                           result.ValidIngredientContainsWheat.Bool,
			ContainsSoy:                             result.ValidIngredientContainsSoy.Bool,
			AnimalDerived:                           result.ValidIngredientAnimalDerived.Bool,
			RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations.Bool,
			ContainsSesame:                          result.ValidIngredientContainsSesame.Bool,
			ContainsFish:                            result.ValidIngredientContainsFish.Bool,
			ContainsGluten:                          result.ValidIngredientContainsGluten.Bool,
			ContainsDairy:                           result.ValidIngredientContainsDairy.Bool,
			ContainsAlcohol:                         result.ValidIngredientContainsAlcohol.Bool,
			AnimalFlesh:                             result.ValidIngredientAnimalFlesh.Bool,
			IsStarch:                                result.ValidIngredientIsStarch.Bool,
			IsProtein:                               result.ValidIngredientIsProtein.Bool,
			IsGrain:                                 result.ValidIngredientIsGrain.Bool,
			IsFruit:                                 result.ValidIngredientIsFruit.Bool,
			IsSalt:                                  result.ValidIngredientIsSalt.Bool,
			IsFat:                                   result.ValidIngredientIsFat.Bool,
			IsAcid:                                  result.ValidIngredientIsAcid.Bool,
			IsHeat:                                  result.ValidIngredientIsHeat.Bool,
		}
	}

	return recipeStepIngredient, nil
}

// getRecipeStepIngredientsForRecipe fetches a list of recipe step ingredients from the database that meet a particular filter.
func (q *Querier) getRecipeStepIngredientsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	results, err := q.generatedQuerier.GetAllRecipeStepIngredientsForRecipe(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	recipeStepIngredients := []*types.RecipeStepIngredient{}
	for _, result := range results {
		recipeStepIngredient := &types.RecipeStepIngredient{
			CreatedAt:                 result.CreatedAt,
			RecipeStepProductID:       stringPointerFromNullString(result.RecipeStepProductID),
			ArchivedAt:                timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:             timePointerFromNullTime(result.LastUpdatedAt),
			MaximumQuantity:           float32PointerFromNullString(result.MaximumQuantityValue),
			VesselIndex:               uint16PointerFromNullInt32(result.VesselIndex),
			ProductPercentageToUse:    float32PointerFromNullString(result.ProductPercentageToUse),
			RecipeStepProductRecipeID: stringPointerFromNullString(result.RecipeStepProductRecipeID),
			QuantityNotes:             result.QuantityNotes,
			ID:                        result.ID,
			BelongsToRecipeStep:       result.BelongsToRecipeStep,
			IngredientNotes:           result.IngredientNotes,
			Name:                      result.Name,
			MeasurementUnit: types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt,
				LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName,
				IconPath:      result.ValidMeasurementUnitIconPath,
				ID:            result.ValidMeasurementUnitID,
				Description:   result.ValidMeasurementUnitDescription,
				PluralName:    result.ValidMeasurementUnitPluralName,
				Slug:          result.ValidMeasurementUnitSlug,
				Volumetric:    boolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal,
				Metric:        result.ValidMeasurementUnitMetric,
				Imperial:      result.ValidMeasurementUnitImperial,
			},
			MinimumQuantity: float32FromString(result.MinimumQuantityValue),
			OptionIndex:     uint16(result.OptionIndex),
			Optional:        result.Optional,
			ToTaste:         result.ToTaste,
		}

		if result.ValidIngredientID.Valid {
			recipeStepIngredient.Ingredient = &types.ValidIngredient{
				CreatedAt:                               result.ValidIngredientCreatedAt.Time,
				LastUpdatedAt:                           timePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:                              timePointerFromNullTime(result.ValidIngredientArchivedAt),
				MaximumIdealStorageTemperatureInCelsius: float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				MinimumIdealStorageTemperatureInCelsius: float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				IconPath:                                result.ValidIngredientIconPath.String,
				Warning:                                 result.ValidIngredientWarning.String,
				PluralName:                              result.ValidIngredientPluralName.String,
				StorageInstructions:                     result.ValidIngredientStorageInstructions.String,
				Name:                                    result.ValidIngredientName.String,
				ID:                                      result.ValidIngredientID.String,
				Description:                             result.ValidIngredientDescription.String,
				Slug:                                    result.ValidIngredientSlug.String,
				ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions.String,
				ContainsShellfish:                       result.ValidIngredientContainsShellfish.Bool,
				IsMeasuredVolumetrically:                result.ValidIngredientVolumetric.Bool,
				IsLiquid:                                boolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:                          result.ValidIngredientContainsPeanut.Bool,
				ContainsTreeNut:                         result.ValidIngredientContainsTreeNut.Bool,
				ContainsEgg:                             result.ValidIngredientContainsEgg.Bool,
				ContainsWheat:                           result.ValidIngredientContainsWheat.Bool,
				ContainsSoy:                             result.ValidIngredientContainsSoy.Bool,
				AnimalDerived:                           result.ValidIngredientAnimalDerived.Bool,
				RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations.Bool,
				ContainsSesame:                          result.ValidIngredientContainsSesame.Bool,
				ContainsFish:                            result.ValidIngredientContainsFish.Bool,
				ContainsGluten:                          result.ValidIngredientContainsGluten.Bool,
				ContainsDairy:                           result.ValidIngredientContainsDairy.Bool,
				ContainsAlcohol:                         result.ValidIngredientContainsAlcohol.Bool,
				AnimalFlesh:                             result.ValidIngredientAnimalFlesh.Bool,
				IsStarch:                                result.ValidIngredientIsStarch.Bool,
				IsProtein:                               result.ValidIngredientIsProtein.Bool,
				IsGrain:                                 result.ValidIngredientIsGrain.Bool,
				IsFruit:                                 result.ValidIngredientIsFruit.Bool,
				IsSalt:                                  result.ValidIngredientIsSalt.Bool,
				IsFat:                                   result.ValidIngredientIsFat.Bool,
				IsAcid:                                  result.ValidIngredientIsAcid.Bool,
				IsHeat:                                  result.ValidIngredientIsHeat.Bool,
			}
		}

		recipeStepIngredients = append(recipeStepIngredients, recipeStepIngredient)
	}

	return recipeStepIngredients, nil
}

// GetRecipeStepIngredients fetches a list of recipe step ingredients from the database that meet a particular filter.
func (q *Querier) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.RecipeStepIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeStepIngredients(ctx, q.db, &generated.GetRecipeStepIngredientsParams{
		RecipeID:      recipeID,
		RecipeStepID:  recipeStepID,
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	for _, result := range results {
		recipeStepIngredient := &types.RecipeStepIngredient{
			CreatedAt:                 result.CreatedAt,
			RecipeStepProductID:       stringPointerFromNullString(result.RecipeStepProductID),
			ArchivedAt:                timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:             timePointerFromNullTime(result.LastUpdatedAt),
			MaximumQuantity:           float32PointerFromNullString(result.MaximumQuantityValue),
			VesselIndex:               uint16PointerFromNullInt32(result.VesselIndex),
			ProductPercentageToUse:    float32PointerFromNullString(result.ProductPercentageToUse),
			RecipeStepProductRecipeID: stringPointerFromNullString(result.RecipeStepProductRecipeID),
			QuantityNotes:             result.QuantityNotes,
			ID:                        result.ID,
			BelongsToRecipeStep:       result.BelongsToRecipeStep,
			IngredientNotes:           result.IngredientNotes,
			Name:                      result.Name,
			MeasurementUnit: types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt,
				LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName,
				IconPath:      result.ValidMeasurementUnitIconPath,
				ID:            result.ValidMeasurementUnitID,
				Description:   result.ValidMeasurementUnitDescription,
				PluralName:    result.ValidMeasurementUnitPluralName,
				Slug:          result.ValidMeasurementUnitSlug,
				Volumetric:    boolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal,
				Metric:        result.ValidMeasurementUnitMetric,
				Imperial:      result.ValidMeasurementUnitImperial,
			},
			MinimumQuantity: float32FromString(result.MinimumQuantityValue),
			OptionIndex:     uint16(result.OptionIndex),
			Optional:        result.Optional,
			ToTaste:         result.ToTaste,
		}

		if result.ValidIngredientID.Valid {
			recipeStepIngredient.Ingredient = &types.ValidIngredient{
				CreatedAt:                               result.ValidIngredientCreatedAt.Time,
				LastUpdatedAt:                           timePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:                              timePointerFromNullTime(result.ValidIngredientArchivedAt),
				MaximumIdealStorageTemperatureInCelsius: float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				MinimumIdealStorageTemperatureInCelsius: float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				IconPath:                                result.ValidIngredientIconPath.String,
				Warning:                                 result.ValidIngredientWarning.String,
				PluralName:                              result.ValidIngredientPluralName.String,
				StorageInstructions:                     result.ValidIngredientStorageInstructions.String,
				Name:                                    result.ValidIngredientName.String,
				ID:                                      result.ValidIngredientID.String,
				Description:                             result.ValidIngredientDescription.String,
				Slug:                                    result.ValidIngredientSlug.String,
				ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions.String,
				ContainsShellfish:                       result.ValidIngredientContainsShellfish.Bool,
				IsMeasuredVolumetrically:                result.ValidIngredientVolumetric.Bool,
				IsLiquid:                                boolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:                          result.ValidIngredientContainsPeanut.Bool,
				ContainsTreeNut:                         result.ValidIngredientContainsTreeNut.Bool,
				ContainsEgg:                             result.ValidIngredientContainsEgg.Bool,
				ContainsWheat:                           result.ValidIngredientContainsWheat.Bool,
				ContainsSoy:                             result.ValidIngredientContainsSoy.Bool,
				AnimalDerived:                           result.ValidIngredientAnimalDerived.Bool,
				RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations.Bool,
				ContainsSesame:                          result.ValidIngredientContainsSesame.Bool,
				ContainsFish:                            result.ValidIngredientContainsFish.Bool,
				ContainsGluten:                          result.ValidIngredientContainsGluten.Bool,
				ContainsDairy:                           result.ValidIngredientContainsDairy.Bool,
				ContainsAlcohol:                         result.ValidIngredientContainsAlcohol.Bool,
				AnimalFlesh:                             result.ValidIngredientAnimalFlesh.Bool,
				IsStarch:                                result.ValidIngredientIsStarch.Bool,
				IsProtein:                               result.ValidIngredientIsProtein.Bool,
				IsGrain:                                 result.ValidIngredientIsGrain.Bool,
				IsFruit:                                 result.ValidIngredientIsFruit.Bool,
				IsSalt:                                  result.ValidIngredientIsSalt.Bool,
				IsFat:                                   result.ValidIngredientIsFat.Bool,
				IsAcid:                                  result.ValidIngredientIsAcid.Bool,
				IsHeat:                                  result.ValidIngredientIsHeat.Bool,
			}
		}

		x.Data = append(x.Data, recipeStepIngredient)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// createRecipeStepIngredient creates a recipe step ingredient in the database.
func (q *Querier) createRecipeStepIngredient(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepIngredientDatabaseCreationInput) (*types.RecipeStepIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// create the recipe step ingredient.
	if err := q.generatedQuerier.CreateRecipeStepIngredient(ctx, db, &generated.CreateRecipeStepIngredientParams{
		QuantityNotes:             input.QuantityNotes,
		Name:                      input.Name,
		BelongsToRecipeStep:       input.BelongsToRecipeStep,
		IngredientNotes:           input.IngredientNotes,
		ID:                        input.ID,
		MinimumQuantityValue:      stringFromFloat32(input.MinimumQuantity),
		RecipeStepProductID:       nullStringFromStringPointer(input.RecipeStepProductID),
		MaximumQuantityValue:      nullStringFromFloat32Pointer(input.MaximumQuantity),
		MeasurementUnit:           nullStringFromString(input.MeasurementUnitID),
		IngredientID:              nullStringFromStringPointer(input.IngredientID),
		ProductPercentageToUse:    nullStringFromFloat32Pointer(input.ProductPercentageToUse),
		RecipeStepProductRecipeID: nullStringFromStringPointer(input.RecipeStepProductRecipeID),
		VesselIndex:               nullInt32FromUint16Pointer(input.VesselIndex),
		OptionIndex:               int32(input.OptionIndex),
		ToTaste:                   input.ToTaste,
		Optional:                  input.Optional,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step ingredient creation query")
	}

	x := &types.RecipeStepIngredient{
		ID:                        input.ID,
		Name:                      input.Name,
		Optional:                  input.Optional,
		MeasurementUnit:           types.ValidMeasurementUnit{ID: input.MeasurementUnitID},
		MinimumQuantity:           input.MinimumQuantity,
		MaximumQuantity:           input.MaximumQuantity,
		QuantityNotes:             input.QuantityNotes,
		IngredientNotes:           input.IngredientNotes,
		BelongsToRecipeStep:       input.BelongsToRecipeStep,
		RecipeStepProductID:       input.RecipeStepProductID,
		OptionIndex:               input.OptionIndex,
		ToTaste:                   input.ToTaste,
		ProductPercentageToUse:    input.ProductPercentageToUse,
		VesselIndex:               input.VesselIndex,
		RecipeStepProductRecipeID: input.RecipeStepProductRecipeID,
		CreatedAt:                 q.currentTime(),
	}

	if input.IngredientID != nil {
		x.Ingredient = &types.ValidIngredient{ID: *input.IngredientID}
	}

	tracing.AttachRecipeStepIngredientIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database.
func (q *Querier) CreateRecipeStepIngredient(ctx context.Context, input *types.RecipeStepIngredientDatabaseCreationInput) (*types.RecipeStepIngredient, error) {
	return q.createRecipeStepIngredient(ctx, q.db, input)
}

// UpdateRecipeStepIngredient updates a particular recipe step ingredient.
func (q *Querier) UpdateRecipeStepIngredient(ctx context.Context, updated *types.RecipeStepIngredient) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepIngredientIDKey, updated.ID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, updated.ID)

	var ingredientID *string
	if updated.Ingredient != nil {
		ingredientID = &updated.Ingredient.ID
	}

	if _, err := q.generatedQuerier.UpdateRecipeStepIngredient(ctx, q.db, &generated.UpdateRecipeStepIngredientParams{
		IngredientID:              nullStringFromStringPointer(ingredientID),
		Name:                      updated.Name,
		Optional:                  updated.Optional,
		MeasurementUnit:           nullStringFromString(updated.MeasurementUnit.ID),
		MinimumQuantityValue:      stringFromFloat32(updated.MinimumQuantity),
		MaximumQuantityValue:      nullStringFromFloat32Pointer(updated.MaximumQuantity),
		QuantityNotes:             updated.QuantityNotes,
		RecipeStepProductID:       nullStringFromStringPointer(updated.RecipeStepProductID),
		IngredientNotes:           updated.IngredientNotes,
		OptionIndex:               int32(updated.OptionIndex),
		ToTaste:                   updated.ToTaste,
		ProductPercentageToUse:    nullStringFromFloat32Pointer(updated.ProductPercentageToUse),
		VesselIndex:               nullInt32FromUint16Pointer(updated.VesselIndex),
		RecipeStepProductRecipeID: nullStringFromStringPointer(updated.RecipeStepProductRecipeID),
		BelongsToRecipeStep:       updated.BelongsToRecipeStep,
		ID:                        updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step ingredient")
	}

	logger.Info("recipe step ingredient updated")

	return nil
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient from the database by its ID.
func (q *Querier) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	if _, err := q.generatedQuerier.ArchiveRecipeStepIngredient(ctx, q.db, &generated.ArchiveRecipeStepIngredientParams{
		RecipeStepID: recipeStepID,
		ID:           recipeStepIngredientID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step ingredient")
	}

	logger.Info("recipe step ingredient archived")

	return nil
}
