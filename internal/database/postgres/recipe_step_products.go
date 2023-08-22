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
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	recipeStepsOnRecipeStepProductsJoinClause = "recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id"
)

var (
	_ types.RecipeStepProductDataManager = (*Querier)(nil)

	// recipeStepProductsTableColumns are the columns for the recipe_step_products table.
	recipeStepProductsTableColumns = []string{
		"recipe_step_products.id",
		"recipe_step_products.name",
		"recipe_step_products.type",
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
		"recipe_step_products.minimum_quantity_value",
		"recipe_step_products.maximum_quantity_value",
		"recipe_step_products.quantity_notes",
		"recipe_step_products.compostable",
		"recipe_step_products.maximum_storage_duration_in_seconds",
		"recipe_step_products.minimum_storage_temperature_in_celsius",
		"recipe_step_products.maximum_storage_temperature_in_celsius",
		"recipe_step_products.storage_instructions",
		"recipe_step_products.is_liquid",
		"recipe_step_products.is_waste",
		"recipe_step_products.index",
		"recipe_step_products.contained_in_vessel_index",
		"recipe_step_products.created_at",
		"recipe_step_products.last_updated_at",
		"recipe_step_products.archived_at",
		"recipe_step_products.belongs_to_recipe_step",
	}

	getRecipeStepProductsJoins = []string{
		recipeStepsOnRecipeStepProductsJoinClause,
		recipesOnRecipeStepsJoinClause,
		validMeasurementUnitsOnRecipeStepProductsJoinClause,
	}
)

// scanRecipeStepProduct takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step product struct.
func (q *Querier) scanRecipeStepProduct(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepProduct, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStepProduct{
		MeasurementUnit: &types.ValidMeasurementUnit{},
	}
	nmu := &types.NullableValidMeasurementUnit{}

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Type,
		&nmu.ID,
		&nmu.Name,
		&nmu.Description,
		&nmu.Volumetric,
		&nmu.IconPath,
		&nmu.Universal,
		&nmu.Metric,
		&nmu.Imperial,
		&nmu.Slug,
		&nmu.PluralName,
		&nmu.CreatedAt,
		&nmu.LastUpdatedAt,
		&nmu.ArchivedAt,
		&x.MinimumQuantity,
		&x.MaximumQuantity,
		&x.QuantityNotes,
		&x.Compostable,
		&x.MaximumStorageDurationInSeconds,
		&x.MinimumStorageTemperatureInCelsius,
		&x.MaximumStorageTemperatureInCelsius,
		&x.StorageInstructions,
		&x.IsLiquid,
		&x.IsWaste,
		&x.Index,
		&x.ContainedInVesselIndex,
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

	if nmu.ID != nil {
		x.MeasurementUnit = converters.ConvertNullableValidMeasurementUnitToValidMeasurementUnit(nmu)
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepProducts takes some database rows and turns them into a slice of recipe step products.
func (q *Querier) scanRecipeStepProducts(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepProducts []*types.RecipeStepProduct, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepProduct(ctx, rows, includeCounts)
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

		recipeStepProducts = append(recipeStepProducts, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return recipeStepProducts, filteredCount, totalCount, nil
}

// RecipeStepProductExists fetches whether a recipe step product exists from the database.
func (q *Querier) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (exists bool, err error) {
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

	if recipeStepProductID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	result, err := q.generatedQuerier.CheckRecipeStepProductExistence(ctx, q.db, &generated.CheckRecipeStepProductExistenceParams{
		RecipeStepID:        recipeStepID,
		RecipeStepProductID: recipeStepProductID,
		RecipeID:            recipeID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step product existence check")
	}

	return result, nil
}

// GetRecipeStepProduct fetches a recipe step product from the database.
func (q *Querier) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
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

	if recipeStepProductID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	result, err := q.generatedQuerier.GetRecipeStepProduct(ctx, q.db, &generated.GetRecipeStepProductParams{
		RecipeStepID:        recipeStepID,
		RecipeStepProductID: recipeStepProductID,
		RecipeID:            recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step product")
	}

	recipeStepProduct := &types.RecipeStepProduct{
		CreatedAt:                          result.CreatedAt,
		MaximumStorageTemperatureInCelsius: float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
		MaximumStorageDurationInSeconds:    uint32PointerFromNullInt32(result.MaximumStorageDurationInSeconds),
		MinimumStorageTemperatureInCelsius: float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
		ArchivedAt:                         timePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:                      timePointerFromNullTime(result.LastUpdatedAt),
		MinimumQuantity:                    float32PointerFromNullString(result.MinimumQuantityValue),
		MeasurementUnit:                    nil,
		MaximumQuantity:                    float32PointerFromNullString(result.MaximumQuantityValue),
		ContainedInVesselIndex:             uint16PointerFromNullInt32(result.ContainedInVesselIndex),
		Name:                               result.Name,
		BelongsToRecipeStep:                result.BelongsToRecipeStep,
		Type:                               string(result.Type),
		ID:                                 result.ID,
		StorageInstructions:                result.StorageInstructions,
		QuantityNotes:                      result.QuantityNotes,
		Index:                              uint16(result.Index),
		IsWaste:                            result.IsWaste,
		IsLiquid:                           result.IsLiquid,
		Compostable:                        result.Compostable,
	}

	if result.ValidMeasurementUnitID != "" {
		recipeStepProduct.MeasurementUnit = &types.ValidMeasurementUnit{
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
		}
	}

	return recipeStepProduct, nil
}

//go:embed queries/recipe_step_products/get_for_recipe.sql
var getRecipeStepProductsForRecipeQuery string

// getRecipeStepProductsForRecipe fetches a list of recipe step products from the database that meet a particular filter.
func (q *Querier) getRecipeStepProductsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepProduct, error) {
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

	rows, err := q.getRows(ctx, q.db, "recipe step products", getRecipeStepProductsForRecipeQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step products list retrieval query")
	}

	recipeStepProducts, _, _, err := q.scanRecipeStepProducts(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step products")
	}

	return recipeStepProducts, nil
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter.
func (q *Querier) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepProduct], err error) {
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

	x = &types.QueryFilteredResult[types.RecipeStepProduct]{
		Pagination: filter.ToPagination(),
	}

	query, args := q.buildListQuery(ctx, "recipe_step_products", getRecipeStepProductsJoins, []string{"valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepProductsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "recipe step products", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step products list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepProducts(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step products")
	}

	return x, nil
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *Querier) createRecipeStepProduct(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// create the recipe step product.
	if err := q.generatedQuerier.CreateRecipeStepProduct(ctx, db, &generated.CreateRecipeStepProductParams{
		QuantityNotes:                      input.QuantityNotes,
		Name:                               input.Name,
		Type:                               generated.RecipeStepProductType(input.Type),
		BelongsToRecipeStep:                input.BelongsToRecipeStep,
		ID:                                 input.ID,
		StorageInstructions:                input.StorageInstructions,
		MinimumQuantityValue:               nullStringFromFloat32Pointer(input.MinimumQuantity),
		MinimumStorageTemperatureInCelsius: nullStringFromFloat32Pointer(input.MinimumStorageTemperatureInCelsius),
		MaximumStorageTemperatureInCelsius: nullStringFromFloat32Pointer(input.MaximumStorageTemperatureInCelsius),
		MaximumQuantityValue:               nullStringFromFloat32Pointer(input.MaximumQuantity),
		MeasurementUnit:                    nullStringFromStringPointer(input.MeasurementUnitID),
		MaximumStorageDurationInSeconds:    nullInt32FromUint32Pointer(input.MaximumStorageDurationInSeconds),
		ContainedInVesselIndex:             nullInt32FromUint16Pointer(input.ContainedInVesselIndex),
		Index:                              int32(input.Index),
		Compostable:                        input.Compostable,
		IsLiquid:                           input.IsLiquid,
		IsWaste:                            input.IsWaste,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step product creation query")
	}

	x := &types.RecipeStepProduct{
		ID:                                 input.ID,
		Name:                               input.Name,
		Type:                               input.Type,
		MinimumQuantity:                    input.MinimumQuantity,
		MaximumQuantity:                    input.MaximumQuantity,
		QuantityNotes:                      input.QuantityNotes,
		BelongsToRecipeStep:                input.BelongsToRecipeStep,
		Compostable:                        input.Compostable,
		MaximumStorageDurationInSeconds:    input.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: input.MaximumStorageTemperatureInCelsius,
		StorageInstructions:                input.StorageInstructions,
		IsLiquid:                           input.IsLiquid,
		IsWaste:                            input.IsWaste,
		Index:                              input.Index,
		ContainedInVesselIndex:             input.ContainedInVesselIndex,
		CreatedAt:                          q.currentTime(),
	}

	if input.MeasurementUnitID != nil {
		x.MeasurementUnit = &types.ValidMeasurementUnit{ID: *input.MeasurementUnitID}
	}

	tracing.AttachRecipeStepProductIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *Querier) CreateRecipeStepProduct(ctx context.Context, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	return q.createRecipeStepProduct(ctx, q.db, input)
}

// UpdateRecipeStepProduct updates a particular recipe step product.
func (q *Querier) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepProductIDKey, updated.ID)
	tracing.AttachRecipeStepProductIDToSpan(span, updated.ID)

	var measurementUnitID *string
	if updated.MeasurementUnit != nil {
		measurementUnitID = &updated.MeasurementUnit.ID
	}

	if err := q.generatedQuerier.UpdateRecipeStepProduct(ctx, q.db, &generated.UpdateRecipeStepProductParams{
		Name:                               updated.Name,
		Type:                               generated.RecipeStepProductType(updated.Type),
		MeasurementUnit:                    nullStringFromStringPointer(measurementUnitID),
		MinimumQuantityValue:               nullStringFromFloat32Pointer(updated.MinimumQuantity),
		MaximumQuantityValue:               nullStringFromFloat32Pointer(updated.MaximumQuantity),
		QuantityNotes:                      updated.QuantityNotes,
		Compostable:                        updated.Compostable,
		MaximumStorageDurationInSeconds:    nullInt32FromUint32Pointer(updated.MaximumStorageDurationInSeconds),
		MinimumStorageTemperatureInCelsius: nullStringFromFloat32Pointer(updated.MinimumStorageTemperatureInCelsius),
		MaximumStorageTemperatureInCelsius: nullStringFromFloat32Pointer(updated.MaximumStorageTemperatureInCelsius),
		StorageInstructions:                updated.StorageInstructions,
		IsLiquid:                           updated.IsLiquid,
		IsWaste:                            updated.IsWaste,
		Index:                              int32(updated.Index),
		ContainedInVesselIndex:             nullInt32FromUint16Pointer(updated.ContainedInVesselIndex),
		BelongsToRecipeStep:                updated.BelongsToRecipeStep,
		ID:                                 updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product updated")

	return nil
}

// ArchiveRecipeStepProduct archives a recipe step product from the database by its ID.
func (q *Querier) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	if err := q.generatedQuerier.ArchiveRecipeStepProduct(ctx, q.db, &generated.ArchiveRecipeStepProductParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepProductID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product archived")

	return nil
}
