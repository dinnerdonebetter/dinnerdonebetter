package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.RecipeStepCompletionConditionDataManager = (*Querier)(nil)

	// recipeStepCompletionConditionsTableColumns are the columns for the recipe_step_completion_conditions table.
	recipeStepCompletionConditionsTableColumns = []string{
		"recipe_step_completion_condition_ingredients.id",
		"recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition",
		"recipe_step_completion_condition_ingredients.recipe_step_ingredient",
		"recipe_step_completion_conditions.id",
		"recipe_step_completion_conditions.belongs_to_recipe_step",
		"valid_ingredient_states.id",
		"valid_ingredient_states.name",
		"valid_ingredient_states.description",
		"valid_ingredient_states.icon_path",
		"valid_ingredient_states.slug",
		"valid_ingredient_states.past_tense",
		"valid_ingredient_states.attribute_type",
		"valid_ingredient_states.created_at",
		"valid_ingredient_states.last_updated_at",
		"valid_ingredient_states.archived_at",
		"recipe_step_completion_conditions.optional",
		"recipe_step_completion_conditions.notes",
		"recipe_step_completion_conditions.created_at",
		"recipe_step_completion_conditions.last_updated_at",
		"recipe_step_completion_conditions.archived_at",
	}
)

// scanRecipeStepCompletionConditionWithIngredients takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step completion condition struct.
func (q *Querier) scanRecipeStepCompletionConditionWithIngredients(ctx context.Context, scan database.ResultIterator, includeCounts bool) (x *types.RecipeStepCompletionCondition, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStepCompletionCondition{}

	for scan.Next() {
		y := &types.RecipeStepCompletionConditionIngredient{}

		targetVars := []any{
			&y.ID,
			&y.BelongsToRecipeStepCompletionCondition,
			&y.RecipeStepIngredient,
			&x.ID,
			&x.BelongsToRecipeStep,
			&x.IngredientState.ID,
			&x.IngredientState.Name,
			&x.IngredientState.Description,
			&x.IngredientState.IconPath,
			&x.IngredientState.Slug,
			&x.IngredientState.PastTense,
			&x.IngredientState.AttributeType,
			&x.IngredientState.CreatedAt,
			&x.IngredientState.LastUpdatedAt,
			&x.IngredientState.ArchivedAt,
			&x.Optional,
			&x.Notes,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		if includeCounts {
			targetVars = append(targetVars, &filteredCount, &totalCount)
		}

		if err = scan.Scan(targetVars...); err != nil {
			return nil, filteredCount, totalCount, observability.PrepareError(err, span, "")
		}

		x.Ingredients = append(x.Ingredients, y)
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepCompletionConditionsWithIngredients takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step completion condition struct.
func (q *Querier) scanRecipeStepCompletionConditionsWithIngredients(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepConditions []*types.RecipeStepCompletionCondition, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	conditionsAndIngredients := map[string]*types.RecipeStepCompletionCondition{}
	idOrder := []string{}

	for rows.Next() {
		x := &types.RecipeStepCompletionCondition{}
		y := &types.RecipeStepCompletionConditionIngredient{}

		targetVars := []any{
			&y.ID,
			&y.BelongsToRecipeStepCompletionCondition,
			&y.RecipeStepIngredient,
			&x.ID,
			&x.BelongsToRecipeStep,
			&x.IngredientState.ID,
			&x.IngredientState.Name,
			&x.IngredientState.Description,
			&x.IngredientState.IconPath,
			&x.IngredientState.Slug,
			&x.IngredientState.PastTense,
			&x.IngredientState.AttributeType,
			&x.IngredientState.CreatedAt,
			&x.IngredientState.LastUpdatedAt,
			&x.IngredientState.ArchivedAt,
			&x.Optional,
			&x.Notes,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		if includeCounts {
			targetVars = append(targetVars, &filteredCount, &totalCount)
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, 0, 0, observability.PrepareError(err, span, "scanning complete recipe step completion condition")
		}

		if _, ok := conditionsAndIngredients[x.ID]; ok {
			conditionsAndIngredients[x.ID].Ingredients = append(conditionsAndIngredients[x.ID].Ingredients, y)
		} else {
			idOrder = append(idOrder, x.ID)
			x.Ingredients = append(x.Ingredients, y)
			conditionsAndIngredients[x.ID] = x
		}
	}

	for _, id := range idOrder {
		recipeStepConditions = append(recipeStepConditions, conditionsAndIngredients[id])
	}

	return recipeStepConditions, filteredCount, totalCount, nil
}

//go:embed queries/recipe_step_completion_conditions/exists.sql
var recipeStepCompletionConditionExistenceQuery string

// RecipeStepCompletionConditionExists fetches whether a recipe step completion condition exists from the database.
func (q *Querier) RecipeStepCompletionConditionExists(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (exists bool, err error) {
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

	if recipeStepCompletionConditionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepCompletionConditionID)

	args := []any{
		recipeStepID,
		recipeStepCompletionConditionID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepCompletionConditionExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step completion condition existence check")
	}

	return result, nil
}

//go:embed queries/recipe_step_completion_conditions/get_one.sql
var getRecipeStepCompletionConditionQuery string

// GetRecipeStepCompletionCondition fetches a recipe step completion condition from the database.
func (q *Querier) GetRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*types.RecipeStepCompletionCondition, error) {
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

	if recipeStepCompletionConditionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepCompletionConditionID)

	args := []any{
		recipeID,
		recipeStepID,
		recipeStepCompletionConditionID,
	}

	rows, err := q.getRows(ctx, q.db, "get recipe step completion condition", getRecipeStepCompletionConditionQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for recipe step completion condition")
	}

	recipeStepCompletionCondition, _, _, err := q.scanRecipeStepCompletionConditionWithIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step completion condition")
	}

	return recipeStepCompletionCondition, nil
}

//go:embed queries/recipe_step_completion_conditions/get_all_for_recipe.sql
var getRecipeStepCompletionConditionsForRecipeQuery string

// getRecipeStepCompletionConditionsForRecipe fetches a recipe step completion condition from the database.
func (q *Querier) getRecipeStepCompletionConditionsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepCompletionCondition, error) {
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

	rows, err := q.getRows(ctx, q.db, "get recipe step completion condition", getRecipeStepCompletionConditionsForRecipeQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for recipe step completion condition")
	}

	recipeStepCompletionCondition, _, _, err := q.scanRecipeStepCompletionConditionsWithIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step completion condition")
	}

	return recipeStepCompletionCondition, nil
}

//go:embed queries/recipe_step_completion_conditions/get_many.sql
var getRecipeStepCompletionConditionsQuery string

// GetRecipeStepCompletionConditions fetches a list of recipe step completion conditions from the database that meet a particular filter.
func (q *Querier) GetRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepCompletionCondition], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()
	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

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

	x = &types.QueryFilteredResult[types.RecipeStepCompletionCondition]{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter.Page != nil {
		x.Page = *filter.Page
	}

	if filter.Limit != nil {
		x.Limit = *filter.Limit
	}

	args := []any{
		filter.CreatedAfter,
		filter.CreatedBefore,
		filter.UpdatedAfter,
		filter.UpdatedBefore,
		filter.Limit,
		filter.QueryOffset(),
	}

	rows, err := q.getRows(ctx, q.db, "recipe step completion conditions", getRecipeStepCompletionConditionsQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step completion conditions list retrieval query")
	}

	data, filteredCount, totalCount, err := q.scanRecipeStepCompletionConditionsWithIngredients(ctx, rows, true)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step completion conditions")
	}

	x.FilteredCount = filteredCount
	x.TotalCount = totalCount
	x.Data = data

	return x, nil
}

//go:embed queries/recipe_step_completion_conditions/create.sql
var recipeStepCompletionConditionCreationQuery string

// createRecipeStepCompletionCondition creates a recipe step completion condition in the database.
func (q *Querier) createRecipeStepCompletionCondition(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepCompletionConditionDatabaseCreationInput) (*types.RecipeStepCompletionCondition, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	args := []any{
		input.ID,
		input.BelongsToRecipeStep,
		input.IngredientStateID,
		input.Optional,
		input.Notes,
	}

	// create the recipe step completion condition.
	if err := q.performWriteQuery(ctx, db, "recipe step completion condition creation", recipeStepCompletionConditionCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step completion condition creation query")
	}

	x := &types.RecipeStepCompletionCondition{
		ID:                  input.ID,
		Notes:               input.Notes,
		Optional:            input.Optional,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		IngredientState:     types.ValidIngredientState{ID: input.IngredientStateID},
		CreatedAt:           q.currentTime(),
	}

	for _, ingredient := range input.Ingredients {
		ingredient.BelongsToRecipeStepCompletionCondition = x.ID
		completionConditionIngredient, err := q.createRecipeStepCompletionConditionIngredient(ctx, db, ingredient)
		if err != nil {
			return nil, observability.PrepareError(err, span, "creating ingredient for recipe step completion condition")
		}

		x.Ingredients = append(x.Ingredients, completionConditionIngredient)
	}

	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, x.ID)

	return x, nil
}

//go:embed queries/recipe_step_completion_condition_ingredients/create.sql
var recipeStepCompletionConditionIngredientCreationQuery string

// createRecipeStepCompletionConditionIngredient creates a recipe step completion condition ingredient in the database.
func (q *Querier) createRecipeStepCompletionConditionIngredient(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepCompletionConditionIngredientDatabaseCreationInput) (*types.RecipeStepCompletionConditionIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	args := []any{
		input.ID,
		input.BelongsToRecipeStepCompletionCondition,
		input.RecipeStepIngredient,
	}

	// create the recipe step completion condition.
	if err := q.performWriteQuery(ctx, db, "recipe step completion condition ingredient creation", recipeStepCompletionConditionIngredientCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step completion condition ingredient creation query")
	}

	x := &types.RecipeStepCompletionConditionIngredient{
		ID:                                     input.ID,
		BelongsToRecipeStepCompletionCondition: input.BelongsToRecipeStepCompletionCondition,
		RecipeStepIngredient:                   input.RecipeStepIngredient,
		CreatedAt:                              q.currentTime(),
	}

	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStepCompletionCondition creates a recipe step completion condition in the database.
func (q *Querier) CreateRecipeStepCompletionCondition(ctx context.Context, input *types.RecipeStepCompletionConditionDatabaseCreationInput) (*types.RecipeStepCompletionCondition, error) {
	return q.createRecipeStepCompletionCondition(ctx, q.db, input)
}

//go:embed queries/recipe_step_completion_conditions/update.sql
var updateRecipeStepCompletionConditionQuery string

// UpdateRecipeStepCompletionCondition updates a particular recipe step completion condition.
func (q *Querier) UpdateRecipeStepCompletionCondition(ctx context.Context, updated *types.RecipeStepCompletionCondition) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepCompletionConditionIDKey, updated.ID)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, updated.ID)

	args := []any{
		updated.Optional,
		updated.Notes,
		updated.BelongsToRecipeStep,
		updated.IngredientState.ID,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step completion condition update", updateRecipeStepCompletionConditionQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step completion condition")
	}

	logger.Info("recipe step completion condition updated")

	return nil
}

//go:embed queries/recipe_step_completion_conditions/archive.sql
var archiveRecipeStepCompletionConditionQuery string

// ArchiveRecipeStepCompletionCondition archives a recipe step completion condition from the database by its ID.
func (q *Querier) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeStepID, recipeStepCompletionConditionID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepCompletionConditionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepCompletionConditionID)

	args := []any{
		recipeStepID,
		recipeStepCompletionConditionID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step completion condition archive", archiveRecipeStepCompletionConditionQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step completion condition")
	}

	logger.Info("recipe step completion condition archived")

	return nil
}
