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
	_ types.RecipeStepConditionDataManager = (*Querier)(nil)

	// recipeStepConditionsTableColumns are the columns for the recipe_step_conditions table.
	recipeStepConditionsTableColumns = []string{
		"recipe_step_conditions.id",
		"recipe_step_conditions.belongs_to_recipe_step",
		"recipe_step_conditions.ingredient_state",
		"recipe_step_conditions.optional",
		"recipe_step_conditions.optional",
		"recipe_step_conditions.created_at",
		"recipe_step_conditions.last_updated_at",
		"recipe_step_conditions.archived_at",
	}

	getRecipeStepConditionsJoins = []string{
		recipesOnRecipeStepsJoinClause,
	}
)

// scanRecipeStepCondition takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step condition struct.
func (q *Querier) scanRecipeStepCondition(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepCondition, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStepCondition{}

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

// scanRecipeStepConditions takes some database rows and turns them into a slice of recipe step conditions.
func (q *Querier) scanRecipeStepConditions(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepConditions []*types.RecipeStepCondition, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepCondition(ctx, rows, includeCounts)
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

		recipeStepConditions = append(recipeStepConditions, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return recipeStepConditions, filteredCount, totalCount, nil
}

//go:embed queries/recipe_step_conditions/exists.sql
var recipeStepConditionExistenceQuery string

// RecipeStepConditionExists fetches whether a recipe step condition exists from the database.
func (q *Querier) RecipeStepConditionExists(ctx context.Context, recipeID, recipeStepID, recipeStepConditionID string) (exists bool, err error) {
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

	if recipeStepConditionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepConditionIDKey, recipeStepConditionID)
	tracing.AttachRecipeStepConditionIDToSpan(span, recipeStepConditionID)

	args := []any{
		recipeStepID,
		recipeStepConditionID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepConditionExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step condition existence check")
	}

	return result, nil
}

//go:embed queries/recipe_step_conditions/get_one.sql
var getRecipeStepConditionQuery string

// GetRecipeStepCondition fetches a recipe step condition from the database.
func (q *Querier) GetRecipeStepCondition(ctx context.Context, recipeID, recipeStepID, recipeStepConditionID string) (*types.RecipeStepCondition, error) {
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

	if recipeStepConditionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepConditionIDKey, recipeStepConditionID)
	tracing.AttachRecipeStepConditionIDToSpan(span, recipeStepConditionID)

	args := []any{
		recipeStepID,
		recipeStepConditionID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "get recipe step condition", getRecipeStepConditionQuery, args)

	recipeStepCondition, _, _, err := q.scanRecipeStepCondition(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeStepCondition")
	}

	return recipeStepCondition, nil
}

// GetRecipeStepConditions fetches a list of recipe step conditions from the database that meet a particular filter.
func (q *Querier) GetRecipeStepConditions(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepCondition], err error) {
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

	x = &types.QueryFilteredResult[types.RecipeStepCondition]{}
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

	query, args := q.buildListQuery(ctx, "recipe_step_conditions", getRecipeStepConditionsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepConditionsTableColumns, "", false, filter)
	rows, err := q.getRows(ctx, q.db, "recipe step conditions", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step conditions list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepConditions(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step conditions")
	}

	return x, nil
}

//go:embed queries/recipe_step_conditions/create.sql
var recipeStepConditionCreationQuery string

// createRecipeStepCondition creates a recipe step condition in the database.
func (q *Querier) createRecipeStepCondition(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepConditionDatabaseCreationInput) (*types.RecipeStepCondition, error) {
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

	// create the recipe step condition.
	if err := q.performWriteQuery(ctx, db, "recipe step condition creation", recipeStepConditionCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step condition creation query")
	}

	x := &types.RecipeStepCondition{
		ID:                  input.ID,
		Notes:               input.Notes,
		Optional:            input.Optional,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		IngredientState:     types.ValidIngredientState{ID: input.IngredientStateID},
		CreatedAt:           q.currentTime(),
	}

	tracing.AttachRecipeStepConditionIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStepCondition creates a recipe step condition in the database.
func (q *Querier) CreateRecipeStepCondition(ctx context.Context, input *types.RecipeStepConditionDatabaseCreationInput) (*types.RecipeStepCondition, error) {
	return q.createRecipeStepCondition(ctx, q.db, input)
}

//go:embed queries/recipe_step_conditions/update.sql
var updateRecipeStepConditionQuery string

// UpdateRecipeStepCondition updates a particular recipe step condition.
func (q *Querier) UpdateRecipeStepCondition(ctx context.Context, updated *types.RecipeStepCondition) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepConditionIDKey, updated.ID)
	tracing.AttachRecipeStepConditionIDToSpan(span, updated.ID)

	args := []any{
		updated.Optional,
		updated.Notes,
		updated.BelongsToRecipeStep,
		updated.IngredientState.ID,
		updated.BelongsToRecipeStep,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step condition update", updateRecipeStepConditionQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step condition")
	}

	logger.Info("recipe step condition updated")

	return nil
}

//go:embed queries/recipe_step_conditions/archive.sql
var archiveRecipeStepConditionQuery string

// ArchiveRecipeStepCondition archives a recipe step condition from the database by its ID.
func (q *Querier) ArchiveRecipeStepCondition(ctx context.Context, recipeStepID, recipeStepConditionID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepConditionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepConditionIDKey, recipeStepConditionID)
	tracing.AttachRecipeStepConditionIDToSpan(span, recipeStepConditionID)

	args := []any{
		recipeStepID,
		recipeStepConditionID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step condition archive", archiveRecipeStepConditionQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step condition")
	}

	logger.Info("recipe step condition archived")

	return nil
}
