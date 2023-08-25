package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	recipeStepsOnRecipeStepIngredientsJoinClause      = "recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id"
	validIngredientsOnRecipeStepIngredientsJoinClause = "valid_ingredients ON recipe_step_ingredients.ingredient_id=valid_ingredients.id"
)

var (
	_ types.RecipeStepIngredientDataManager = (*Querier)(nil)

	// recipeStepIngredientsTableColumns are the columns for the recipe_step_ingredients table.
	recipeStepIngredientsTableColumns = []string{
		"recipe_step_ingredients.id",
		"recipe_step_ingredients.name",
		"recipe_step_ingredients.optional",
		"valid_ingredients.id",
		"valid_ingredients.name",
		"valid_ingredients.description",
		"valid_ingredients.warning",
		"valid_ingredients.contains_egg",
		"valid_ingredients.contains_dairy",
		"valid_ingredients.contains_peanut",
		"valid_ingredients.contains_tree_nut",
		"valid_ingredients.contains_soy",
		"valid_ingredients.contains_wheat",
		"valid_ingredients.contains_shellfish",
		"valid_ingredients.contains_sesame",
		"valid_ingredients.contains_fish",
		"valid_ingredients.contains_gluten",
		"valid_ingredients.animal_flesh",
		"valid_ingredients.volumetric",
		"valid_ingredients.is_liquid",
		"valid_ingredients.icon_path",
		"valid_ingredients.animal_derived",
		"valid_ingredients.plural_name",
		"valid_ingredients.restrict_to_preparations",
		"valid_ingredients.minimum_ideal_storage_temperature_in_celsius",
		"valid_ingredients.maximum_ideal_storage_temperature_in_celsius",
		"valid_ingredients.storage_instructions",
		"valid_ingredients.slug",
		"valid_ingredients.contains_alcohol",
		"valid_ingredients.shopping_suggestions",
		"valid_ingredients.is_starch",
		"valid_ingredients.is_protein",
		"valid_ingredients.is_grain",
		"valid_ingredients.is_fruit",
		"valid_ingredients.is_salt",
		"valid_ingredients.is_fat",
		"valid_ingredients.is_acid",
		"valid_ingredients.is_heat",
		"valid_ingredients.created_at",
		"valid_ingredients.last_updated_at",
		"valid_ingredients.archived_at",
		"valid_measurement_units.id",
		"valid_measurement_units.name",
		"valid_measurement_units.description",
		"valid_measurement_units.volumetric",
		"valid_measurement_units.icon_path",
		"valid_measurement_units.universal",
		"valid_measurement_units.metric",
		"valid_measurement_units.imperial",
		"valid_measurement_units.slug",
		"valid_measurement_units.plural_name",
		"valid_measurement_units.created_at",
		"valid_measurement_units.last_updated_at",
		"valid_measurement_units.archived_at",
		"recipe_step_ingredients.minimum_quantity_value",
		"recipe_step_ingredients.maximum_quantity_value",
		"recipe_step_ingredients.quantity_notes",
		"recipe_step_ingredients.recipe_step_product_id",
		"recipe_step_ingredients.ingredient_notes",
		"recipe_step_ingredients.option_index",
		"recipe_step_ingredients.to_taste",
		"recipe_step_ingredients.product_percentage_to_use",
		"recipe_step_ingredients.vessel_index",
		"recipe_step_ingredients.recipe_step_product_recipe_id",
		"recipe_step_ingredients.created_at",
		"recipe_step_ingredients.last_updated_at",
		"recipe_step_ingredients.archived_at",
		"recipe_step_ingredients.belongs_to_recipe_step",
	}

	getRecipeStepIngredientsJoins = []string{
		recipeStepsOnRecipeStepIngredientsJoinClause,
		recipesOnRecipeStepsJoinClause,
		validIngredientsOnRecipeStepIngredientsJoinClause,
		validMeasurementUnitsOnRecipeStepIngredientsJoinClause,
	}
)

// scanRecipeStepIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step ingredient struct.
func (q *Querier) scanRecipeStepIngredient(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStepIngredient{}

	ingredient := &types.NullableValidIngredient{}

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Optional,
		&ingredient.ID,
		&ingredient.Name,
		&ingredient.Description,
		&ingredient.Warning,
		&ingredient.ContainsEgg,
		&ingredient.ContainsDairy,
		&ingredient.ContainsPeanut,
		&ingredient.ContainsTreeNut,
		&ingredient.ContainsSoy,
		&ingredient.ContainsWheat,
		&ingredient.ContainsShellfish,
		&ingredient.ContainsSesame,
		&ingredient.ContainsFish,
		&ingredient.ContainsGluten,
		&ingredient.AnimalFlesh,
		&ingredient.IsMeasuredVolumetrically,
		&ingredient.IsLiquid,
		&ingredient.IconPath,
		&ingredient.AnimalDerived,
		&ingredient.PluralName,
		&ingredient.RestrictToPreparations,
		&ingredient.MinimumIdealStorageTemperatureInCelsius,
		&ingredient.MaximumIdealStorageTemperatureInCelsius,
		&ingredient.StorageInstructions,
		&ingredient.Slug,
		&ingredient.ContainsAlcohol,
		&ingredient.ShoppingSuggestions,
		&ingredient.IsStarch,
		&ingredient.IsProtein,
		&ingredient.IsGrain,
		&ingredient.IsFruit,
		&ingredient.IsSalt,
		&ingredient.IsFat,
		&ingredient.IsAcid,
		&ingredient.IsHeat,
		&ingredient.CreatedAt,
		&ingredient.LastUpdatedAt,
		&ingredient.ArchivedAt,
		&x.MeasurementUnit.ID,
		&x.MeasurementUnit.Name,
		&x.MeasurementUnit.Description,
		&x.MeasurementUnit.Volumetric,
		&x.MeasurementUnit.IconPath,
		&x.MeasurementUnit.Universal,
		&x.MeasurementUnit.Metric,
		&x.MeasurementUnit.Imperial,
		&x.MeasurementUnit.Slug,
		&x.MeasurementUnit.PluralName,
		&x.MeasurementUnit.CreatedAt,
		&x.MeasurementUnit.LastUpdatedAt,
		&x.MeasurementUnit.ArchivedAt,
		&x.MinimumQuantity,
		&x.MaximumQuantity,
		&x.QuantityNotes,
		&x.RecipeStepProductID,
		&x.IngredientNotes,
		&x.OptionIndex,
		&x.ToTaste,
		&x.ProductPercentageToUse,
		&x.VesselIndex,
		&x.RecipeStepProductRecipeID,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
		&x.BelongsToRecipeStep,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	if ingredient.ID != nil {
		x.Ingredient = &types.ValidIngredient{
			CreatedAt:                               *ingredient.CreatedAt,
			LastUpdatedAt:                           ingredient.LastUpdatedAt,
			ArchivedAt:                              ingredient.ArchivedAt,
			ID:                                      *ingredient.ID,
			Warning:                                 *ingredient.Warning,
			Description:                             *ingredient.Description,
			IconPath:                                *ingredient.IconPath,
			PluralName:                              *ingredient.PluralName,
			StorageInstructions:                     *ingredient.StorageInstructions,
			Slug:                                    *ingredient.Slug,
			ContainsAlcohol:                         *ingredient.ContainsAlcohol,
			ShoppingSuggestions:                     *ingredient.ShoppingSuggestions,
			IsStarch:                                *ingredient.IsStarch,
			IsProtein:                               *ingredient.IsProtein,
			IsGrain:                                 *ingredient.IsGrain,
			IsFruit:                                 *ingredient.IsFruit,
			IsSalt:                                  *ingredient.IsSalt,
			IsFat:                                   *ingredient.IsFat,
			IsAcid:                                  *ingredient.IsAcid,
			IsHeat:                                  *ingredient.IsHeat,
			Name:                                    *ingredient.Name,
			MaximumIdealStorageTemperatureInCelsius: ingredient.MaximumIdealStorageTemperatureInCelsius,
			MinimumIdealStorageTemperatureInCelsius: ingredient.MinimumIdealStorageTemperatureInCelsius,
			ContainsShellfish:                       *ingredient.ContainsShellfish,
			ContainsDairy:                           *ingredient.ContainsDairy,
			AnimalFlesh:                             *ingredient.AnimalFlesh,
			IsMeasuredVolumetrically:                *ingredient.IsMeasuredVolumetrically,
			IsLiquid:                                *ingredient.IsLiquid,
			ContainsPeanut:                          *ingredient.ContainsPeanut,
			ContainsTreeNut:                         *ingredient.ContainsTreeNut,
			ContainsEgg:                             *ingredient.ContainsEgg,
			ContainsWheat:                           *ingredient.ContainsWheat,
			ContainsSoy:                             *ingredient.ContainsSoy,
			AnimalDerived:                           *ingredient.AnimalDerived,
			RestrictToPreparations:                  *ingredient.RestrictToPreparations,
			ContainsSesame:                          *ingredient.ContainsSesame,
			ContainsFish:                            *ingredient.ContainsFish,
			ContainsGluten:                          *ingredient.ContainsGluten,
		}
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepIngredients takes some database rows and turns them into a slice of recipe step ingredients.
func (q *Querier) scanRecipeStepIngredients(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepIngredients []*types.RecipeStepIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepIngredient(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		recipeStepIngredients = append(recipeStepIngredients, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return recipeStepIngredients, filteredCount, totalCount, nil
}

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

	if result.ValidIngredientID != "" {
		recipeStepIngredient.Ingredient = &types.ValidIngredient{
			CreatedAt:                               result.ValidIngredientCreatedAt,
			LastUpdatedAt:                           timePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
			ArchivedAt:                              timePointerFromNullTime(result.ValidIngredientArchivedAt),
			MaximumIdealStorageTemperatureInCelsius: float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
			MinimumIdealStorageTemperatureInCelsius: float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
			IconPath:                                result.ValidIngredientIconPath,
			Warning:                                 result.ValidIngredientWarning,
			PluralName:                              result.ValidIngredientPluralName,
			StorageInstructions:                     result.ValidIngredientStorageInstructions,
			Name:                                    result.ValidIngredientName,
			ID:                                      result.ValidIngredientID,
			Description:                             result.ValidIngredientDescription,
			Slug:                                    result.ValidIngredientSlug,
			ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions,
			ContainsShellfish:                       result.ValidIngredientContainsShellfish,
			IsMeasuredVolumetrically:                result.ValidIngredientVolumetric,
			IsLiquid:                                boolFromNullBool(result.ValidIngredientIsLiquid),
			ContainsPeanut:                          result.ValidIngredientContainsPeanut,
			ContainsTreeNut:                         result.ValidIngredientContainsTreeNut,
			ContainsEgg:                             result.ValidIngredientContainsEgg,
			ContainsWheat:                           result.ValidIngredientContainsWheat,
			ContainsSoy:                             result.ValidIngredientContainsSoy,
			AnimalDerived:                           result.ValidIngredientAnimalDerived,
			RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations,
			ContainsSesame:                          result.ValidIngredientContainsSesame,
			ContainsFish:                            result.ValidIngredientContainsFish,
			ContainsGluten:                          result.ValidIngredientContainsGluten,
			ContainsDairy:                           result.ValidIngredientContainsDairy,
			ContainsAlcohol:                         result.ValidIngredientContainsAlcohol,
			AnimalFlesh:                             result.ValidIngredientAnimalFlesh,
			IsStarch:                                result.ValidIngredientIsStarch,
			IsProtein:                               result.ValidIngredientIsProtein,
			IsGrain:                                 result.ValidIngredientIsGrain,
			IsFruit:                                 result.ValidIngredientIsFruit,
			IsSalt:                                  result.ValidIngredientIsSalt,
			IsFat:                                   result.ValidIngredientIsFat,
			IsAcid:                                  result.ValidIngredientIsAcid,
			IsHeat:                                  result.ValidIngredientIsHeat,
		}
	}

	return recipeStepIngredient, nil
}

//go:embed queries/recipe_step_ingredients/get_for_recipe.sql
var getRecipeStepIngredientsForRecipeQuery string

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

	args := []any{
		recipeID,
	}

	rows, err := q.getRows(ctx, q.db, "recipe step ingredients for recipe", getRecipeStepIngredientsForRecipeQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	recipeStepIngredients, _, _, err := q.scanRecipeStepIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step ingredients")
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

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	query, args := q.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, recipeStepIngredientsTableColumns, filter)
	rows, err := q.getRows(ctx, q.db, "recipe step ingredients", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepIngredients(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step ingredients")
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
		ingredientID = &updated.MeasurementUnit.ID
	}

	if err := q.generatedQuerier.UpdateRecipeStepIngredient(ctx, q.db, &generated.UpdateRecipeStepIngredientParams{
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

	if err := q.generatedQuerier.ArchiveRecipeStepIngredient(ctx, q.db, &generated.ArchiveRecipeStepIngredientParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepIngredientID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step ingredient")
	}

	logger.Info("recipe step ingredient archived")

	return nil
}
