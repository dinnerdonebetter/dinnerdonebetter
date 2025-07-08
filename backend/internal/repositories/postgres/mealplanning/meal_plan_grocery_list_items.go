package mealplanning

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	generated2 "github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ mealplanning.MealPlanGroceryListItemDataManager = (*Querier)(nil)
)

// MealPlanGroceryListItemExists fetches whether a meal plan grocery list exists from the database.
func (q *Querier) MealPlanGroceryListItemExists(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	result, err := q.generatedQuerier.CheckMealPlanGroceryListItemExistence(ctx, q.db, &generated2.CheckMealPlanGroceryListItemExistenceParams{
		MealPlanID:                mealPlanID,
		MealPlanGroceryListItemID: mealPlanGroceryListItemID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan grocery list existence check")
	}

	return result, nil
}

func (q *Querier) fleshOutMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItem *mealplanning.MealPlanGroceryListItem) (*mealplanning.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanGroceryListItem == nil {
		return nil, database.ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItem.ID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItem.ID)

	validIngredient, err := q.recipeenumsRepository.GetValidIngredient(ctx, mealPlanGroceryListItem.Ingredient.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching grocery list item ingredient")
	}
	mealPlanGroceryListItem.Ingredient = *validIngredient

	validMeasurementUnit, err := q.recipeenumsRepository.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.MeasurementUnit.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching grocery list item measurement unit")
	}
	mealPlanGroceryListItem.MeasurementUnit = *validMeasurementUnit

	if mealPlanGroceryListItem.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnit, getPurchasedMeasurementUnitErr := q.recipeenumsRepository.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.PurchasedMeasurementUnit.ID)
		if getPurchasedMeasurementUnitErr != nil {
			return nil, observability.PrepareAndLogError(getPurchasedMeasurementUnitErr, logger, span, "fetching grocery list item purchased measurement unit")
		}
		mealPlanGroceryListItem.PurchasedMeasurementUnit = purchasedMeasurementUnit
	}

	return mealPlanGroceryListItem, nil
}

// GetMealPlanGroceryListItem fetches a meal plan grocery list from the database.
func (q *Querier) GetMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*mealplanning.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	result, err := q.generatedQuerier.GetMealPlanGroceryListItem(ctx, q.db, &generated2.GetMealPlanGroceryListItemParams{
		MealPlanID:                mealPlanID,
		MealPlanGroceryListItemID: mealPlanGroceryListItemID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan grocery list item")
	}

	mealPlanGroceryListItem := &mealplanning.MealPlanGroceryListItem{
		CreatedAt:         result.CreatedAt,
		LastUpdatedAt:     database.TimePointerFromNullTime(result.LastUpdatedAt),
		PurchasePrice:     database.Float32PointerFromNullString(result.PurchasePrice),
		PurchasedUPC:      database.StringPointerFromNullString(result.PurchasedUpc),
		ArchivedAt:        database.TimePointerFromNullTime(result.ArchivedAt),
		QuantityPurchased: database.Float32PointerFromNullString(result.QuantityPurchased),

		BelongsToMealPlan: result.BelongsToMealPlan,
		Status:            string(result.Status),
		StatusExplanation: result.StatusExplanation,
		ID:                result.ID,
		QuantityNeeded: types.Float32RangeWithOptionalMax{
			Max: database.Float32PointerFromNullString(result.MaximumQuantityNeeded),
			Min: database.Float32FromString(result.MinimumQuantityNeeded),
		},
		MeasurementUnit: recipeenums.ValidMeasurementUnit{
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
		Ingredient: recipeenums.ValidIngredient{
			CreatedAt:     result.ValidIngredientCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
			},
			IconPath:               result.ValidIngredientIconPath,
			Warning:                result.ValidIngredientWarning,
			PluralName:             result.ValidIngredientPluralName,
			StorageInstructions:    result.ValidIngredientStorageInstructions,
			Name:                   result.ValidIngredientName,
			ID:                     result.ValidIngredientID,
			Description:            result.ValidIngredientDescription,
			Slug:                   result.ValidIngredientSlug,
			ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions,
			ContainsShellfish:      result.ValidIngredientContainsShellfish,
			IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
			ContainsPeanut:         result.ValidIngredientContainsPeanut,
			ContainsTreeNut:        result.ValidIngredientContainsTreeNut,
			ContainsEgg:            result.ValidIngredientContainsEgg,
			ContainsWheat:          result.ValidIngredientContainsWheat,
			ContainsSoy:            result.ValidIngredientContainsSoy,
			AnimalDerived:          result.ValidIngredientAnimalDerived,
			RestrictToPreparations: result.ValidIngredientRestrictToPreparations,
			ContainsSesame:         result.ValidIngredientContainsSesame,
			ContainsFish:           result.ValidIngredientContainsFish,
			ContainsGluten:         result.ValidIngredientContainsGluten,
			ContainsDairy:          result.ValidIngredientContainsDairy,
			ContainsAlcohol:        result.ValidIngredientContainsAlcohol,
			AnimalFlesh:            result.ValidIngredientAnimalFlesh,
			IsStarch:               result.ValidIngredientIsStarch,
			IsProtein:              result.ValidIngredientIsProtein,
			IsGrain:                result.ValidIngredientIsGrain,
			IsFruit:                result.ValidIngredientIsFruit,
			IsSalt:                 result.ValidIngredientIsSalt,
			IsFat:                  result.ValidIngredientIsFat,
			IsAcid:                 result.ValidIngredientIsAcid,
			IsHeat:                 result.ValidIngredientIsHeat,
		},
	}

	if result.PurchasedMeasurementUnit.Valid {
		mealPlanGroceryListItem.PurchasedMeasurementUnit = &recipeenums.ValidMeasurementUnit{
			ID: result.PurchasedMeasurementUnit.String,
		}
	}

	if mealPlanGroceryListItem.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnit, getPurchasedMeasurementUnitErr := q.recipeenumsRepository.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.PurchasedMeasurementUnit.ID)
		if getPurchasedMeasurementUnitErr != nil {
			return nil, observability.PrepareAndLogError(getPurchasedMeasurementUnitErr, logger, span, "fetching grocery list item purchased measurement unit")
		}
		mealPlanGroceryListItem.PurchasedMeasurementUnit = purchasedMeasurementUnit
	}

	return mealPlanGroceryListItem, nil
}

// GetMealPlanGroceryListItemsForMealPlan fetches a list of meal plan grocery lists from the database that meet a particular filter.
func (q *Querier) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string) ([]*mealplanning.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	results, err := q.generatedQuerier.GetMealPlanGroceryListItemsForMealPlan(ctx, q.db, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan grocery list items list retrieval query")
	}

	x := []*mealplanning.MealPlanGroceryListItem{}
	for _, result := range results {
		mealPlanGroceryListItem := &mealplanning.MealPlanGroceryListItem{
			CreatedAt:         result.CreatedAt,
			LastUpdatedAt:     database.TimePointerFromNullTime(result.LastUpdatedAt),
			PurchasePrice:     database.Float32PointerFromNullString(result.PurchasePrice),
			PurchasedUPC:      database.StringPointerFromNullString(result.PurchasedUpc),
			ArchivedAt:        database.TimePointerFromNullTime(result.ArchivedAt),
			QuantityPurchased: database.Float32PointerFromNullString(result.QuantityPurchased),
			BelongsToMealPlan: result.BelongsToMealPlan,
			Status:            string(result.Status),
			StatusExplanation: result.StatusExplanation,
			ID:                result.ID,
			QuantityNeeded: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.MaximumQuantityNeeded),
				Min: database.Float32FromString(result.MinimumQuantityNeeded),
			},
			MeasurementUnit: recipeenums.ValidMeasurementUnit{
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
			Ingredient: recipeenums.ValidIngredient{
				CreatedAt:     result.ValidIngredientCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				StorageTemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				},
				IconPath:               result.ValidIngredientIconPath,
				Warning:                result.ValidIngredientWarning,
				PluralName:             result.ValidIngredientPluralName,
				StorageInstructions:    result.ValidIngredientStorageInstructions,
				Name:                   result.ValidIngredientName,
				ID:                     result.ValidIngredientID,
				Description:            result.ValidIngredientDescription,
				Slug:                   result.ValidIngredientSlug,
				ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions,
				ContainsShellfish:      result.ValidIngredientContainsShellfish,
				IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:         result.ValidIngredientContainsPeanut,
				ContainsTreeNut:        result.ValidIngredientContainsTreeNut,
				ContainsEgg:            result.ValidIngredientContainsEgg,
				ContainsWheat:          result.ValidIngredientContainsWheat,
				ContainsSoy:            result.ValidIngredientContainsSoy,
				AnimalDerived:          result.ValidIngredientAnimalDerived,
				RestrictToPreparations: result.ValidIngredientRestrictToPreparations,
				ContainsSesame:         result.ValidIngredientContainsSesame,
				ContainsFish:           result.ValidIngredientContainsFish,
				ContainsGluten:         result.ValidIngredientContainsGluten,
				ContainsDairy:          result.ValidIngredientContainsDairy,
				ContainsAlcohol:        result.ValidIngredientContainsAlcohol,
				AnimalFlesh:            result.ValidIngredientAnimalFlesh,
				IsStarch:               result.ValidIngredientIsStarch,
				IsProtein:              result.ValidIngredientIsProtein,
				IsGrain:                result.ValidIngredientIsGrain,
				IsFruit:                result.ValidIngredientIsFruit,
				IsSalt:                 result.ValidIngredientIsSalt,
				IsFat:                  result.ValidIngredientIsFat,
				IsAcid:                 result.ValidIngredientIsAcid,
				IsHeat:                 result.ValidIngredientIsHeat,
			},
		}

		if result.PurchasedMeasurementUnit.Valid {
			mealPlanGroceryListItem.PurchasedMeasurementUnit = &recipeenums.ValidMeasurementUnit{
				ID: result.PurchasedMeasurementUnit.String,
			}
		}

		if mealPlanGroceryListItem.PurchasedMeasurementUnit != nil {
			purchasedMeasurementUnit, getPurchasedMeasurementUnitErr := q.recipeenumsRepository.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.PurchasedMeasurementUnit.ID)
			if getPurchasedMeasurementUnitErr != nil {
				return nil, observability.PrepareAndLogError(getPurchasedMeasurementUnitErr, logger, span, "fetching grocery list item purchased measurement unit")
			}
			mealPlanGroceryListItem.PurchasedMeasurementUnit = purchasedMeasurementUnit
		}

		x = append(x, mealPlanGroceryListItem)
	}

	for i := range x {
		x[i], err = q.fleshOutMealPlanGroceryListItem(ctx, x[i])
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "augmenting grocery list item data")
		}
	}

	return x, nil
}

// createMealPlanGroceryListItem creates a meal plan grocery list in the database.
func (q *Querier) createMealPlanGroceryListItem(ctx context.Context, querier database.SQLQueryExecutor, input *mealplanning.MealPlanGroceryListItemDatabaseCreationInput) (*mealplanning.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanGroceryListItemIDKey, input.ID)

	// create the meal plan grocery list.
	if err := q.generatedQuerier.CreateMealPlanGroceryListItem(ctx, querier, &generated2.CreateMealPlanGroceryListItemParams{
		ID:                       input.ID,
		BelongsToMealPlan:        input.BelongsToMealPlan,
		ValidIngredient:          input.ValidIngredientID,
		ValidMeasurementUnit:     input.ValidMeasurementUnitID,
		MinimumQuantityNeeded:    database.StringFromFloat32(input.QuantityNeeded.Min),
		StatusExplanation:        input.StatusExplanation,
		Status:                   generated2.GroceryListItemStatus(input.Status),
		MaximumQuantityNeeded:    database.NullStringFromFloat32Pointer(input.QuantityNeeded.Max),
		QuantityPurchased:        database.NullStringFromFloat32Pointer(input.QuantityPurchased),
		PurchasedMeasurementUnit: database.NullStringFromStringPointer(input.PurchasedMeasurementUnitID),
		PurchasedUpc:             database.NullStringFromStringPointer(input.PurchasedUPC),
		PurchasePrice:            database.NullStringFromFloat32Pointer(input.PurchasePrice),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan grocery list creation query")
	}

	x := &mealplanning.MealPlanGroceryListItem{
		ID:                input.ID,
		BelongsToMealPlan: input.BelongsToMealPlan,
		Ingredient:        recipeenums.ValidIngredient{ID: input.ValidIngredientID},
		MeasurementUnit:   recipeenums.ValidMeasurementUnit{ID: input.ValidMeasurementUnitID},
		QuantityNeeded: types.Float32RangeWithOptionalMax{
			Max: input.QuantityNeeded.Max,
			Min: input.QuantityNeeded.Min,
		},
		QuantityPurchased: input.QuantityPurchased,
		PurchasedUPC:      input.PurchasedUPC,
		PurchasePrice:     input.PurchasePrice,
		StatusExplanation: input.StatusExplanation,
		Status:            input.Status,
		CreatedAt:         q.CurrentTime(),
	}

	if input.PurchasedMeasurementUnitID != nil {
		x.PurchasedMeasurementUnit = &recipeenums.ValidMeasurementUnit{ID: *input.PurchasedMeasurementUnitID}
	}

	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, x.ID)
	logger.Info("meal plan grocery list created")

	return x, nil
}

// CreateMealPlanGroceryListItem creates a meal plan grocery list in the database.
func (q *Querier) CreateMealPlanGroceryListItem(ctx context.Context, input *mealplanning.MealPlanGroceryListItemDatabaseCreationInput) (*mealplanning.MealPlanGroceryListItem, error) {
	return q.createMealPlanGroceryListItem(ctx, q.db, input)
}

// UpdateMealPlanGroceryListItem updates a particular meal plan grocery list.
func (q *Querier) UpdateMealPlanGroceryListItem(ctx context.Context, updated *mealplanning.MealPlanGroceryListItem) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealPlanGroceryListItemIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, updated.ID)

	var purchasedMeasurementUnitID *string
	if updated.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnitID = &updated.PurchasedMeasurementUnit.ID
	}

	if _, err := q.generatedQuerier.UpdateMealPlanGroceryListItem(ctx, q.db, &generated2.UpdateMealPlanGroceryListItemParams{
		BelongsToMealPlan:        updated.BelongsToMealPlan,
		ValidIngredient:          updated.Ingredient.ID,
		ValidMeasurementUnit:     updated.MeasurementUnit.ID,
		MinimumQuantityNeeded:    database.StringFromFloat32(updated.QuantityNeeded.Min),
		StatusExplanation:        updated.StatusExplanation,
		Status:                   generated2.GroceryListItemStatus(updated.Status),
		ID:                       updated.ID,
		MaximumQuantityNeeded:    database.NullStringFromFloat32Pointer(updated.QuantityNeeded.Max),
		QuantityPurchased:        database.NullStringFromFloat32Pointer(updated.QuantityPurchased),
		PurchasedMeasurementUnit: database.NullStringFromStringPointer(purchasedMeasurementUnitID),
		PurchasedUpc:             database.NullStringFromStringPointer(updated.PurchasedUPC),
		PurchasePrice:            database.NullStringFromFloat32Pointer(updated.PurchasePrice),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan grocery list")
	}

	logger.Info("meal plan grocery list updated")

	return nil
}

// ArchiveMealPlanGroceryListItem archives a meal plan grocery list from the database by its ID.
func (q *Querier) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItemID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanGroceryListItemID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	if _, err := q.generatedQuerier.ArchiveMealPlanGroceryListItem(ctx, q.db, mealPlanGroceryListItemID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan grocery list")
	}

	logger.Info("meal plan grocery list archived")

	return nil
}
