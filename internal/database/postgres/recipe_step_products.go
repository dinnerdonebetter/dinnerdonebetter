package postgres

import (
	"context"
	_ "embed"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
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

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Type,
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

	if x.MeasurementUnit != nil && x.MeasurementUnit.ID == "" {
		x.MeasurementUnit = nil
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

//go:embed queries/recipe_step_products/exists.sql
var recipeStepProductExistenceQuery string

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

	args := []any{
		recipeStepID,
		recipeStepProductID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepProductExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step product existence check")
	}

	return result, nil
}

//go:embed queries/recipe_step_products/get_one.sql
var getRecipeStepProductQuery string

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

	args := []any{
		recipeStepID,
		recipeStepProductID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "recipeStepProduct", getRecipeStepProductQuery, args)

	recipeStepProduct, _, _, err := q.scanRecipeStepProduct(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeStepProduct")
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

	x = &types.QueryFilteredResult[types.RecipeStepProduct]{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
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

//go:embed queries/recipe_step_products/create.sql
var recipeStepProductCreationQuery string

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *Querier) createRecipeStepProduct(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	args := []any{
		input.ID,
		input.Name,
		input.Type,
		input.MeasurementUnitID,
		input.MinimumQuantity,
		input.MaximumQuantity,
		input.QuantityNotes,
		input.Compostable,
		input.MaximumStorageDurationInSeconds,
		input.MinimumStorageTemperatureInCelsius,
		input.MaximumStorageTemperatureInCelsius,
		input.StorageInstructions,
		input.BelongsToRecipeStep,
		input.IsLiquid,
		input.IsWaste,
		input.Index,
		input.ContainedInVesselIndex,
	}

	// create the recipe step product.
	if err := q.performWriteQuery(ctx, db, "recipe step product creation", recipeStepProductCreationQuery, args); err != nil {
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

//go:embed queries/recipe_step_products/update.sql
var updateRecipeStepProductQuery string

// UpdateRecipeStepProduct updates a particular recipe step product.
func (q *Querier) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepProductIDKey, updated.ID)
	tracing.AttachRecipeStepProductIDToSpan(span, updated.ID)

	args := []any{
		updated.Name,
		updated.Type,
		updated.MeasurementUnit.ID,
		updated.MinimumQuantity,
		updated.MaximumQuantity,
		updated.QuantityNotes,
		updated.Compostable,
		updated.MaximumStorageDurationInSeconds,
		updated.MinimumStorageTemperatureInCelsius,
		updated.MaximumStorageTemperatureInCelsius,
		updated.StorageInstructions,
		updated.IsLiquid,
		updated.IsWaste,
		updated.Index,
		updated.ContainedInVesselIndex,
		updated.BelongsToRecipeStep,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step product update", updateRecipeStepProductQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product updated")

	return nil
}

//go:embed queries/recipe_step_products/archive.sql
var archiveRecipeStepProductQuery string

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

	args := []any{
		recipeStepID,
		recipeStepProductID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step product archive", archiveRecipeStepProductQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product archived")

	return nil
}
