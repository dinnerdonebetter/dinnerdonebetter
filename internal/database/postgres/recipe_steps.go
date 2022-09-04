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
	recipesOnRecipeStepsJoinClause = "recipes ON recipe_steps.belongs_to_recipe=recipes.id"
)

var (
	_ types.RecipeStepDataManager = (*SQLQuerier)(nil)

	// recipeStepsTableColumns are the columns for the recipe_steps table.
	recipeStepsTableColumns = []string{
		"recipe_steps.id",
		"recipe_steps.index",
		"valid_preparations.id",
		"valid_preparations.name",
		"valid_preparations.description",
		"valid_preparations.icon_path",
		"valid_preparations.yields_nothing",
		"valid_preparations.restrict_to_ingredients",
		"valid_preparations.zero_ingredients_allowable",
		"valid_preparations.past_tense",
		"valid_preparations.created_at",
		"valid_preparations.last_updated_at",
		"valid_preparations.archived_at",
		"recipe_steps.minimum_estimated_time_in_seconds",
		"recipe_steps.maximum_estimated_time_in_seconds",
		"recipe_steps.minimum_temperature_in_celsius",
		"recipe_steps.maximum_temperature_in_celsius",
		"recipe_steps.notes",
		"recipe_steps.explicit_instructions",
		"recipe_steps.optional",
		"recipe_steps.created_at",
		"recipe_steps.last_updated_at",
		"recipe_steps.archived_at",
		"recipe_steps.belongs_to_recipe",
	}

	getRecipeStepsJoins = []string{
		recipesOnRecipeStepsJoinClause,
		validPreparationsOnRecipeStepsJoinClause,
	}
)

// scanRecipeStep takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step struct.
func (q *SQLQuerier) scanRecipeStep(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.RecipeStep{}

	targetVars := []interface{}{
		&x.ID,
		&x.Index,
		&x.Preparation.ID,
		&x.Preparation.Name,
		&x.Preparation.Description,
		&x.Preparation.IconPath,
		&x.Preparation.YieldsNothing,
		&x.Preparation.RestrictToIngredients,
		&x.Preparation.ZeroIngredientsAllowable,
		&x.Preparation.PastTense,
		&x.Preparation.CreatedAt,
		&x.Preparation.LastUpdatedAt,
		&x.Preparation.ArchivedAt,
		&x.MinimumEstimatedTimeInSeconds,
		&x.MaximumEstimatedTimeInSeconds,
		&x.MinimumTemperatureInCelsius,
		&x.MaximumTemperatureInCelsius,
		&x.Notes,
		&x.ExplicitInstructions,
		&x.Optional,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
		&x.BelongsToRecipe,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeSteps takes some database rows and turns them into a slice of recipe steps.
func (q *SQLQuerier) scanRecipeSteps(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeSteps []*types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStep(ctx, rows, includeCounts)
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

		recipeSteps = append(recipeSteps, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipeSteps, filteredCount, totalCount, nil
}

const recipeStepExistenceQuery = "SELECT EXISTS ( SELECT recipe_steps.id FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.archived_at IS NULL AND recipes.id = $3 )"

// RecipeStepExists fetches whether a recipe step exists from the database.
func (q *SQLQuerier) RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (exists bool, err error) {
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

	args := []interface{}{
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe step existence check")
	}

	return result, nil
}

const getRecipeStepQuery = `SELECT 
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.zero_ingredients_allowable,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.optional,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe 
FROM recipe_steps
JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_at IS NULL
AND recipe_steps.belongs_to_recipe = $1
AND recipe_steps.id = $2
AND recipes.archived_at IS NULL 
AND recipes.id = $3
`

// GetRecipeStep fetches a recipe step from the database.
func (q *SQLQuerier) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
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

	args := []interface{}{
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "recipeStep", getRecipeStepQuery, args)

	recipeStep, _, _, err := q.scanRecipeStep(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipeStep")
	}

	return recipeStep, nil
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipeSteps(ctx context.Context, recipeID string, filter *types.QueryFilter) (x *types.RecipeStepList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	x = &types.RecipeStepList{}
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

	query, args := q.buildListQuery(ctx, "recipe_steps", getRecipeStepsJoins, []string{"valid_preparations.id"}, nil, householdOwnershipColumn, recipeStepsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "recipeSteps", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe steps list retrieval query")
	}

	if x.RecipeSteps, x.FilteredCount, x.TotalCount, err = q.scanRecipeSteps(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe steps")
	}

	return x, nil
}

const recipeStepCreationQuery = `INSERT INTO recipe_steps
    (id,index,preparation_id,minimum_estimated_time_in_seconds,maximum_estimated_time_in_seconds,minimum_temperature_in_celsius,maximum_temperature_in_celsius,notes,explicit_instructions,optional,belongs_to_recipe) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`

// CreateRecipeStep creates a recipe step in the database.
func (q *SQLQuerier) createRecipeStep(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Index,
		input.PreparationID,
		input.MinimumEstimatedTimeInSeconds,
		input.MaximumEstimatedTimeInSeconds,
		input.MinimumTemperatureInCelsius,
		input.MaximumTemperatureInCelsius,
		input.Notes,
		input.ExplicitInstructions,
		input.Optional,
		input.BelongsToRecipe,
	}

	// create the recipe step.
	if err := q.performWriteQuery(ctx, db, "recipe step creation", recipeStepCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing recipe step creation")
	}

	x := &types.RecipeStep{
		ID:                            input.ID,
		Index:                         input.Index,
		Preparation:                   types.ValidPreparation{ID: input.PreparationID},
		MinimumEstimatedTimeInSeconds: input.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: input.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   input.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   input.MaximumTemperatureInCelsius,
		Notes:                         input.Notes,
		ExplicitInstructions:          input.ExplicitInstructions,
		Optional:                      input.Optional,
		BelongsToRecipe:               input.BelongsToRecipe,
		CreatedAt:                     q.currentTime(),
	}

	for i, ingredientInput := range input.Ingredients {
		ingredientInput.BelongsToRecipeStep = x.ID
		ingredient, createErr := q.createRecipeStepIngredient(ctx, db, ingredientInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step ingredient #%d", i+1)
		}

		x.Ingredients = append(x.Ingredients, ingredient)
	}

	for i, productInput := range input.Products {
		productInput.BelongsToRecipeStep = x.ID
		product, createErr := q.createRecipeStepProduct(ctx, db, productInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step product #%d", i+1)
		}

		x.Products = append(x.Products, product)
	}

	for i, instrumentInput := range input.Instruments {
		instrumentInput.BelongsToRecipeStep = x.ID
		instrument, createErr := q.createRecipeStepInstrument(ctx, db, instrumentInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step instrument #%d", i+1)
		}

		x.Instruments = append(x.Instruments, instrument)
	}

	tracing.AttachRecipeStepIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStep creates a recipe step in the database.
func (q *SQLQuerier) CreateRecipeStep(ctx context.Context, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	return q.createRecipeStep(ctx, q.db, input)
}

const updateRecipeStepQuery = `UPDATE recipe_steps SET 
	index = $1,
	preparation_id = $2,
	minimum_estimated_time_in_seconds = $3,
	maximum_estimated_time_in_seconds = $4,
	minimum_temperature_in_celsius = $5,
	maximum_temperature_in_celsius = $6,
	notes = $7,
	explicit_instructions = $8, 
	optional = $9,
	last_updated_at = extract(epoch FROM NOW())
WHERE archived_at IS NULL 
  AND belongs_to_recipe = $10 
  AND id = $11
`

// UpdateRecipeStep updates a particular recipe step.
func (q *SQLQuerier) UpdateRecipeStep(ctx context.Context, updated *types.RecipeStep) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepIDKey, updated.ID)
	tracing.AttachRecipeStepIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Index,
		updated.Preparation.ID,
		updated.MinimumEstimatedTimeInSeconds,
		updated.MaximumEstimatedTimeInSeconds,
		updated.MinimumTemperatureInCelsius,
		updated.MaximumTemperatureInCelsius,
		updated.Notes,
		updated.ExplicitInstructions,
		updated.Optional,
		updated.BelongsToRecipe,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step update", updateRecipeStepQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step")
	}

	logger.Info("recipe step updated")

	return nil
}

const archiveRecipeStepQuery = "UPDATE recipe_steps SET archived_at = extract(epoch FROM NOW()) WHERE archived_at IS NULL AND belongs_to_recipe = $1 AND id = $2"

// ArchiveRecipeStep archives a recipe step from the database by its ID.
func (q *SQLQuerier) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	args := []interface{}{
		recipeID,
		recipeStepID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step archive", archiveRecipeStepQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step")
	}

	logger.Info("recipe step archived")

	return nil
}
