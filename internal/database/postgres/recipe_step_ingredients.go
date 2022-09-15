package postgres

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	recipeStepsOnRecipeStepIngredientsJoinClause = "recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id"
)

var (
	_ types.RecipeStepIngredientDataManager = (*Querier)(nil)

	// recipeStepIngredientsTableColumns are the columns for the recipe_step_ingredients table.
	recipeStepIngredientsTableColumns = []string{
		"recipe_step_ingredients.id",
		"recipe_step_ingredients.name",
		"recipe_step_ingredients.optional",
		"recipe_step_ingredients.ingredient_id",
		"valid_measurement_units.id",
		"valid_measurement_units.name",
		"valid_measurement_units.description",
		"valid_measurement_units.volumetric",
		"valid_measurement_units.icon_path",
		"valid_measurement_units.universal",
		"valid_measurement_units.metric",
		"valid_measurement_units.imperial",
		"valid_measurement_units.plural_name",
		"valid_measurement_units.created_at",
		"valid_measurement_units.last_updated_at",
		"valid_measurement_units.archived_at",
		"recipe_step_ingredients.minimum_quantity_value",
		"recipe_step_ingredients.maximum_quantity_value",
		"recipe_step_ingredients.quantity_notes",
		"recipe_step_ingredients.product_of_recipe_step",
		"recipe_step_ingredients.recipe_step_product_id",
		"recipe_step_ingredients.ingredient_notes",
		"recipe_step_ingredients.created_at",
		"recipe_step_ingredients.last_updated_at",
		"recipe_step_ingredients.archived_at",
		"recipe_step_ingredients.belongs_to_recipe_step",
	}

	getRecipeStepIngredientsJoins = []string{
		recipeStepsOnRecipeStepIngredientsJoinClause,
		recipesOnRecipeStepsJoinClause,
		validMeasurementUnitsOnRecipeStepIngredientsJoinClause,
	}
)

// scanRecipeStepIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step ingredient struct.
func (q *Querier) scanRecipeStepIngredient(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStepIngredient{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Optional,
		&x.IngredientID,
		&x.MeasurementUnit.ID,
		&x.MeasurementUnit.Name,
		&x.MeasurementUnit.Description,
		&x.MeasurementUnit.Volumetric,
		&x.MeasurementUnit.IconPath,
		&x.MeasurementUnit.Universal,
		&x.MeasurementUnit.Metric,
		&x.MeasurementUnit.Imperial,
		&x.MeasurementUnit.PluralName,
		&x.MeasurementUnit.CreatedAt,
		&x.MeasurementUnit.LastUpdatedAt,
		&x.MeasurementUnit.ArchivedAt,
		&x.MinimumQuantity,
		&x.MaximumQuantity,
		&x.QuantityNotes,
		&x.ProductOfRecipeStep,
		&x.RecipeStepProductID,
		&x.IngredientNotes,
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

const recipeStepIngredientExistenceQuery = "SELECT EXISTS ( SELECT recipe_step_ingredients.id FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_at IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id = $2 AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_at IS NULL AND recipes.id = $5 )"

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

	args := []interface{}{
		recipeStepID,
		recipeStepIngredientID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepIngredientExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step ingredient existence check")
	}

	return result, nil
}

const getRecipeStepIngredientQuery = `SELECT
	recipe_step_ingredients.id,
	recipe_step_ingredients.name,
	recipe_step_ingredients.optional,
	recipe_step_ingredients.ingredient_id,
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.plural_name,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at,
	recipe_step_ingredients.minimum_quantity_value,
	recipe_step_ingredients.maximum_quantity_value,
	recipe_step_ingredients.quantity_notes,
	recipe_step_ingredients.product_of_recipe_step,
	recipe_step_ingredients.recipe_step_product_id,
	recipe_step_ingredients.ingredient_notes,
	recipe_step_ingredients.created_at,
	recipe_step_ingredients.last_updated_at,
	recipe_step_ingredients.archived_at,
	recipe_step_ingredients.belongs_to_recipe_step
FROM recipe_step_ingredients
JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id
JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
JOIN valid_ingredients ON recipe_step_ingredients.ingredient_id=valid_ingredients.id
JOIN valid_measurement_units ON recipe_step_ingredients.measurement_unit=valid_measurement_units.id
WHERE recipe_step_ingredients.archived_at IS NULL
AND recipe_step_ingredients.belongs_to_recipe_step = $1
AND recipe_step_ingredients.id = $2
AND recipe_steps.archived_at IS NULL
AND recipe_steps.belongs_to_recipe = $3
AND recipe_steps.id = $4
AND recipes.archived_at IS NULL
AND recipes.id = $5
`

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

	args := []interface{}{
		recipeStepID,
		recipeStepIngredientID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "get recipe step ingredient", getRecipeStepIngredientQuery, args)

	recipeStepIngredient, _, _, err := q.scanRecipeStepIngredient(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeStepIngredient")
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

	query, args := q.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, nil, false)
	rows, err := q.performReadQuery(ctx, q.db, "recipe step ingredients", query, args)
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
func (q *Querier) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.RecipeStepIngredientList, err error) {
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

	x = &types.RecipeStepIngredientList{}
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

	query, args := q.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, filter, true)
	rows, err := q.performReadQuery(ctx, q.db, "recipeStepIngredients", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	if x.RecipeStepIngredients, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepIngredients(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step ingredients")
	}

	return x, nil
}

const recipeStepIngredientCreationQuery = `INSERT INTO recipe_step_ingredients (
	id,
	name,
	optional,
	ingredient_id,
	measurement_unit,
	minimum_quantity_value,
	maximum_quantity_value,
	quantity_notes,
	product_of_recipe_step,
	recipe_step_product_id,
	ingredient_notes,
	belongs_to_recipe_step
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
`

// createRecipeStepIngredient creates a recipe step ingredient in the database.
func (q *Querier) createRecipeStepIngredient(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepIngredientDatabaseCreationInput) (*types.RecipeStepIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	args := []interface{}{
		input.ID,
		input.Name,
		input.Optional,
		input.IngredientID,
		input.MeasurementUnitID,
		input.MinimumQuantity,
		input.MaximumQuantity,
		input.QuantityNotes,
		input.ProductOfRecipeStep,
		input.RecipeStepProductID,
		input.IngredientNotes,
		input.BelongsToRecipeStep,
	}

	// create the recipe step ingredient.
	if err := q.performWriteQuery(ctx, db, "recipe step ingredient creation", recipeStepIngredientCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step ingredient creation query")
	}

	x := &types.RecipeStepIngredient{
		ID:                  input.ID,
		Name:                input.Name,
		Optional:            input.Optional,
		IngredientID:        input.IngredientID,
		MeasurementUnit:     types.ValidMeasurementUnit{ID: input.MeasurementUnitID},
		MinimumQuantity:     input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
		QuantityNotes:       input.QuantityNotes,
		ProductOfRecipeStep: input.ProductOfRecipeStep,
		IngredientNotes:     input.IngredientNotes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		RecipeStepProductID: input.RecipeStepProductID,
		CreatedAt:           q.currentTime(),
	}

	tracing.AttachRecipeStepIngredientIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database.
func (q *Querier) CreateRecipeStepIngredient(ctx context.Context, input *types.RecipeStepIngredientDatabaseCreationInput) (*types.RecipeStepIngredient, error) {
	return q.createRecipeStepIngredient(ctx, q.db, input)
}

const updateRecipeStepIngredientQuery = `
UPDATE recipe_step_ingredients SET
	ingredient_id = $1,
	name = $2,
	optional = $3,
	measurement_unit = $4,
	minimum_quantity_value = $5,
	maximum_quantity_value = $6,
	quantity_notes = $7,
	product_of_recipe_step = $8,
	recipe_step_product_id = $9,
	ingredient_notes = $10,
	last_updated_at = NOW() 
WHERE archived_at IS NULL AND belongs_to_recipe_step = $11
AND id = $12
`

// UpdateRecipeStepIngredient updates a particular recipe step ingredient.
func (q *Querier) UpdateRecipeStepIngredient(ctx context.Context, updated *types.RecipeStepIngredient) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepIngredientIDKey, updated.ID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.IngredientID,
		updated.Name,
		updated.Optional,
		updated.MeasurementUnit.ID,
		updated.MinimumQuantity,
		updated.MaximumQuantity,
		updated.QuantityNotes,
		updated.ProductOfRecipeStep,
		updated.RecipeStepProductID,
		updated.IngredientNotes,
		updated.BelongsToRecipeStep,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step ingredient update", updateRecipeStepIngredientQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step ingredient")
	}

	logger.Info("recipe step ingredient updated")

	return nil
}

const archiveRecipeStepIngredientQuery = "UPDATE recipe_step_ingredients SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2"

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

	args := []interface{}{
		recipeStepID,
		recipeStepIngredientID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step ingredient archive", archiveRecipeStepIngredientQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step ingredient")
	}

	logger.Info("recipe step ingredient archived")

	return nil
}
