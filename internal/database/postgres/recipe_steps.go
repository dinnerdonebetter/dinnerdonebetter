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
		"valid_preparations.minimum_ingredient_count",
		"valid_preparations.maximum_ingredient_count",
		"valid_preparations.minimum_instrument_count",
		"valid_preparations.maximum_instrument_count",
		"valid_preparations.temperature_required",
		"valid_preparations.time_estimate_required",
		"valid_preparations.condition_expression_required",
		"valid_preparations.consumes_vessel",
		"valid_preparations.only_for_vessels",
		"valid_preparations.minimum_vessel_count",
		"valid_preparations.maximum_vessel_count",
		"valid_preparations.slug",
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
		"recipe_steps.condition_expression",
		"recipe_steps.optional",
		"recipe_steps.start_timer_automatically",
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
		&x.Preparation.MinimumIngredientCount,
		&x.Preparation.MaximumIngredientCount,
		&x.Preparation.MinimumInstrumentCount,
		&x.Preparation.MaximumInstrumentCount,
		&x.Preparation.TemperatureRequired,
		&x.Preparation.TimeEstimateRequired,
		&x.Preparation.ConditionExpressionRequired,
		&x.Preparation.ConsumesVessel,
		&x.Preparation.OnlyForVessels,
		&x.Preparation.MinimumVesselCount,
		&x.Preparation.MaximumVesselCount,
		&x.Preparation.Slug,
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
		&x.ConditionExpression,
		&x.Optional,
		&x.StartTimerAutomatically,
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

	result, err := q.generatedQuerier.CheckRecipeStepExistence(ctx, q.db, &generated.CheckRecipeStepExistenceParams{
		BelongsToRecipe: recipeID,
		ID:              recipeStepID,
	})
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

//go:embed queries/recipe_steps/get_one_by_recipe_id.sql
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
func (q *Querier) GetRecipeSteps(ctx context.Context, recipeID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStep], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.RecipeStep]{
		Pagination: filter.ToPagination(),
	}

	query, args := q.buildListQuery(ctx, "recipe_steps", getRecipeStepsJoins, []string{"valid_preparations.id"}, nil, householdOwnershipColumn, recipeStepsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "recipe steps", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe steps list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanRecipeSteps(ctx, rows, true); err != nil {
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
		input.ConditionExpression,
		input.Optional,
		input.StartTimerAutomatically,
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
		ConditionExpression:           input.ConditionExpression,
		Optional:                      input.Optional,
		BelongsToRecipe:               input.BelongsToRecipe,
		StartTimerAutomatically:       input.StartTimerAutomatically,
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

	for i, vesselInput := range input.Vessels {
		vesselInput.BelongsToRecipeStep = x.ID
		vessel, createErr := q.createRecipeStepVessel(ctx, db, vesselInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, span, "creating recipe step vessel #%d", i+1)
		}

		x.Vessels = append(x.Vessels, vessel)
	}

	for i, conditionInput := range input.CompletionConditions {
		conditionInput.BelongsToRecipeStep = x.ID
		condition, createErr := q.createRecipeStepCompletionCondition(ctx, db, conditionInput)
		if createErr != nil {
			return nil, observability.PrepareError(createErr, span, "creating recipe step completion condition #%d", i+1)
		}

		x.CompletionConditions = append(x.CompletionConditions, condition)
	}

	tracing.AttachRecipeStepIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStep creates a recipe step in the database.
func (q *Querier) CreateRecipeStep(ctx context.Context, input *types.RecipeStepDatabaseCreationInput) (*types.RecipeStep, error) {
	return q.createRecipeStep(ctx, q.db, input)
}

// UpdateRecipeStep updates a particular recipe step.
func (q *Querier) UpdateRecipeStep(ctx context.Context, updated *types.RecipeStep) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepIDKey, updated.ID)
	tracing.AttachRecipeStepIDToSpan(span, updated.ID)

	if err := q.generatedQuerier.UpdateRecipeStep(ctx, q.db, &generated.UpdateRecipeStepParams{
		ConditionExpression:           updated.ConditionExpression,
		PreparationID:                 updated.Preparation.ID,
		ID:                            updated.ID,
		BelongsToRecipe:               updated.BelongsToRecipe,
		Notes:                         updated.Notes,
		ExplicitInstructions:          updated.ExplicitInstructions,
		MinimumTemperatureInCelsius:   nullStringFromFloat32Pointer(updated.MinimumTemperatureInCelsius),
		MaximumTemperatureInCelsius:   nullStringFromFloat32Pointer(updated.MaximumTemperatureInCelsius),
		MaximumEstimatedTimeInSeconds: nullInt64FromUint32Pointer(updated.MaximumEstimatedTimeInSeconds),
		MinimumEstimatedTimeInSeconds: nullInt64FromUint32Pointer(updated.MinimumEstimatedTimeInSeconds),
		Index:                         int32(updated.Index),
		Optional:                      updated.Optional,
		StartTimerAutomatically:       updated.StartTimerAutomatically,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step")
	}

	logger.Info("recipe step updated")

	return nil
}

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

	if err := q.generatedQuerier.ArchiveRecipeStep(ctx, q.db, &generated.ArchiveRecipeStepParams{
		BelongsToRecipe: recipeID,
		ID:              recipeStepID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step")
	}

	logger.Info("recipe step archived")

	return nil
}
