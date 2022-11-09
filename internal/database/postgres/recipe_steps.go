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
	recipesOnRecipeStepsJoinClause = "recipes ON recipe_steps.belongs_to_recipe=recipes.id"
)

var (
	_ types.RecipeStepDataManager = (*Querier)(nil)

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
func (q *Querier) scanRecipeStep(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStep{}

	targetVars := []any{
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
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeSteps takes some database rows and turns them into a slice of recipe steps.
func (q *Querier) scanRecipeSteps(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeSteps []*types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

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
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return recipeSteps, filteredCount, totalCount, nil
}

//go:embed queries/recipe_steps/exists.sql
var recipeStepExistenceQuery string

// RecipeStepExists fetches whether a recipe step exists from the database.
func (q *Querier) RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (exists bool, err error) {
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

	args := []any{
		recipeID,
		recipeStepID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step existence check")
	}

	return result, nil
}

//go:embed queries/recipe_steps/get_one.sql
var getRecipeStepQuery string

// GetRecipeStep fetches a recipe step from the database.
func (q *Querier) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
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

	args := []any{
		recipeID,
		recipeStepID,
	}

	row := q.getOneRow(ctx, q.db, "recipe step", getRecipeStepQuery, args)

	recipeStep, _, _, err := q.scanRecipeStep(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeStep")
	}

	return recipeStep, nil
}

//go:embed queries/recipe_steps/get_one_by_id.sql
var getRecipeStepByIDQuery string

// getRecipeStepByID fetches a recipe step from the database.
func (q *Querier) getRecipeStepByID(ctx context.Context, querier database.SQLQueryExecutor, recipeStepID string) (*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	args := []any{
		recipeStepID,
	}

	row := q.getOneRow(ctx, querier, "recipe step", getRecipeStepByIDQuery, args)

	recipeStep, _, _, err := q.scanRecipeStep(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeStep")
	}

	return recipeStep, nil
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter.
func (q *Querier) GetRecipeSteps(ctx context.Context, recipeID string, filter *types.QueryFilter) (x *types.RecipeStepList, err error) {
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

	query, args := q.buildListQuery(ctx, "recipe_steps", getRecipeStepsJoins, []string{"valid_preparations.id"}, nil, householdOwnershipColumn, recipeStepsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "recipeSteps", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe steps list retrieval query")
	}

	if x.RecipeSteps, x.FilteredCount, x.TotalCount, err = q.scanRecipeSteps(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe steps")
	}

	return x, nil
}

//go:embed queries/recipe_steps/create.sql
var recipeStepCreationQuery string

// CreateRecipeStep creates a recipe step in the database.
func (q *Querier) createRecipeStep(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	args := []any{
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
		return nil, observability.PrepareError(err, span, "performing recipe step creation")
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

	tracing.AttachRecipeStepIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStep creates a recipe step in the database.
func (q *Querier) CreateRecipeStep(ctx context.Context, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	return q.createRecipeStep(ctx, q.db, input)
}

//go:embed queries/recipe_steps/update.sql
var updateRecipeStepQuery string

// UpdateRecipeStep updates a particular recipe step.
func (q *Querier) UpdateRecipeStep(ctx context.Context, updated *types.RecipeStep) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepIDKey, updated.ID)
	tracing.AttachRecipeStepIDToSpan(span, updated.ID)

	args := []any{
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
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step")
	}

	logger.Info("recipe step updated")

	return nil
}

//go:embed queries/recipe_steps/archive.sql
var archiveRecipeStepQuery string

// ArchiveRecipeStep archives a recipe step from the database by its ID.
func (q *Querier) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
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

	args := []any{
		recipeID,
		recipeStepID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step archive", archiveRecipeStepQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step")
	}

	logger.Info("recipe step archived")

	return nil
}
