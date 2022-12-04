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

var (
	_ types.RecipeStepCompletionConditionDataManager = (*Querier)(nil)

	// recipeStepCompletionConditionsTableColumns are the columns for the recipe_step_completion_conditions table.
	recipeStepCompletionConditionsTableColumns = []string{
		"recipe_step_completion_conditions.id",
		"recipe_step_completion_conditions.belongs_to_recipe_step",
		"recipe_step_completion_conditions.ingredient_state",
		"recipe_step_completion_conditions.optional",
		"recipe_step_completion_conditions.optional",
		"recipe_step_completion_conditions.created_at",
		"recipe_step_completion_conditions.last_updated_at",
		"recipe_step_completion_conditions.archived_at",
	}

	getRecipeStepCompletionConditionsJoins = []string{
		recipesOnRecipeStepsJoinClause,
	}
)

// scanRecipeStepCompletionCondition takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step completion condition struct.
func (q *Querier) scanRecipeStepCompletionCondition(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepCompletionCondition, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStepCompletionCondition{}

	targetVars := []any{
		&x.ID,
		&x.BelongsToRecipeStep,
		&x.IngredientState.ID,
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
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepCompletionConditions takes some database rows and turns them into a slice of recipe step completion conditions.
func (q *Querier) scanRecipeStepCompletionConditions(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepCompletionConditions []*types.RecipeStepCompletionCondition, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepCompletionCondition(ctx, rows, includeCounts)
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

		recipeStepCompletionConditions = append(recipeStepCompletionConditions, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return recipeStepCompletionConditions, filteredCount, totalCount, nil
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
		recipeStepID,
		recipeStepCompletionConditionID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "get recipe step completion condition", getRecipeStepCompletionConditionQuery, args)

	recipeStepCompletionCondition, _, _, err := q.scanRecipeStepCompletionCondition(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeStepCompletionCondition")
	}

	return recipeStepCompletionCondition, nil
}

// GetRecipeStepCompletionConditions fetches a list of recipe step completion conditions from the database that meet a particular filter.
func (q *Querier) GetRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepCompletionCondition], err error) {
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

	x = &types.QueryFilteredResult[types.RecipeStepCompletionCondition]{}
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

	query, args := q.buildListQuery(ctx, "recipe_step_completion_conditions", getRecipeStepCompletionConditionsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepCompletionConditionsTableColumns, "", false, filter)
	rows, err := q.getRows(ctx, q.db, "recipe step completion conditions", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step completion conditions list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepCompletionConditions(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step completion conditions")
	}

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
		updated.BelongsToRecipeStep,
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
