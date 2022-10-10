package postgres

import (
	"context"
	_ "embed"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.RecipePrepTaskDataManager = (*Querier)(nil)
)

// scanRecipePrepTaskWithSteps takes a database Scanner (i.e. *sql.Row) and scans the result into a meal struct.
func (q *Querier) scanRecipePrepTaskWithSteps(ctx context.Context, rows database.ResultIterator) (x *types.RecipePrepTask, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipePrepTask{}

	for rows.Next() {
		recipePrepTaskRecipeStep := &types.RecipePrepTaskStep{}

		targetVars := []interface{}{
			&x.ID,
			&x.Notes,
			&x.ExplicitStorageInstructions,
			&x.MinimumTimeBufferBeforeRecipeInSeconds,
			&x.MaximumTimeBufferBeforeRecipeInSeconds,
			&x.StorageType,
			&x.MinimumStorageTemperatureInCelsius,
			&x.MaximumStorageTemperatureInCelsius,
			&x.BelongsToRecipe,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
			&recipePrepTaskRecipeStep.ID,
			&recipePrepTaskRecipeStep.BelongsToRecipeStep,
			&recipePrepTaskRecipeStep.BelongsToRecipePrepTask,
			&recipePrepTaskRecipeStep.SatisfiesRecipeStep,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, observability.PrepareError(err, span, "scanning complete meal")
		}

		x.TaskSteps = append(x.TaskSteps, recipePrepTaskRecipeStep)
	}

	return x, nil
}

// scanRecipePrepTasksWithSteps takes a database Scanner (i.e. *sql.Row) and scans the result into a meal struct.
func (q *Querier) scanRecipePrepTasksWithSteps(ctx context.Context, rows database.ResultIterator) (recipePrepTasks []*types.RecipePrepTask, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	var lastRecipePrepTaskID string

	x := &types.RecipePrepTask{}

	for rows.Next() {
		recipePrepTaskRecipeStep := &types.RecipePrepTaskStep{}

		targetVars := []interface{}{
			&x.ID,
			&x.Notes,
			&x.ExplicitStorageInstructions,
			&x.MinimumTimeBufferBeforeRecipeInSeconds,
			&x.MaximumTimeBufferBeforeRecipeInSeconds,
			&x.StorageType,
			&x.MinimumStorageTemperatureInCelsius,
			&x.MaximumStorageTemperatureInCelsius,
			&x.BelongsToRecipe,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
			&recipePrepTaskRecipeStep.ID,
			&recipePrepTaskRecipeStep.BelongsToRecipeStep,
			&recipePrepTaskRecipeStep.BelongsToRecipePrepTask,
			&recipePrepTaskRecipeStep.SatisfiesRecipeStep,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, observability.PrepareError(err, span, "scanning complete meal")
		}

		x.TaskSteps = append(x.TaskSteps, recipePrepTaskRecipeStep)

		if lastRecipePrepTaskID == "" {
			lastRecipePrepTaskID = x.ID
		}

		if len(recipePrepTasks) > 0 && lastRecipePrepTaskID != recipePrepTasks[len(recipePrepTasks)-1].ID {
			recipePrepTasks = append(recipePrepTasks, x)
		}
	}

	recipePrepTasks = append(recipePrepTasks, x)

	return recipePrepTasks, nil
}

//go:embed queries/recipe_prep_tasks/exists.sql
var recipePrepTasksExistsQuery string

// RecipePrepTaskExists checks if a recipe prep task exists.
func (q *Querier) RecipePrepTaskExists(ctx context.Context, recipePrepTaskID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipePrepTaskID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	args := []interface{}{
		recipePrepTaskID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipePrepTasksExistsQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe prep task existence check")
	}

	logger.Info("recipe prep task existence retrieved")

	return result, nil
}

//go:embed queries/recipe_prep_tasks/get_one.sql
var getRecipePrepTasksQuery string

// GetRecipePrepTask fetches a recipe prep task.
func (q *Querier) GetRecipePrepTask(ctx context.Context, recipePrepTaskID string) (x *types.RecipePrepTask, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipePrepTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	args := []interface{}{
		recipePrepTaskID,
	}

	rows, err := q.getRows(ctx, q.db, "recipe prep task", getRecipePrepTasksQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe prep task rows")
	}

	x, err = q.scanRecipePrepTaskWithSteps(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe prep task")
	}

	logger.Info("recipe prep tasks retrieved")

	return x, nil
}

//go:embed queries/recipe_prep_tasks/create.sql
var createRecipePrepTaskQuery string

// createRecipePrepTask creates a recipe prep task.
func (q *Querier) createRecipePrepTask(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.RecipePrepTaskDatabaseCreationInput) (*types.RecipePrepTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Notes,
		input.ExplicitStorageInstructions,
		input.MinimumTimeBufferBeforeRecipeInSeconds,
		input.MaximumTimeBufferBeforeRecipeInSeconds,
		input.StorageType,
		input.MinimumStorageTemperatureInCelsius,
		input.MaximumStorageTemperatureInCelsius,
		input.BelongsToRecipe,
	}

	// create the recipe prep task.
	if err := q.performWriteQuery(ctx, querier, "recipe prep task creation", createRecipePrepTaskQuery, args); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
	}

	x := &types.RecipePrepTask{
		CreatedAt: q.timeFunc(),
		ID:        input.ID,
	}

	for _, recipeStep := range input.TaskSteps {
		if err := q.createRecipePrepTaskStep(ctx, querier, recipeStep); err != nil {
			q.rollbackTransaction(ctx, querier)
			return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
		}
	}

	tracing.AttachMealPlanIDToSpan(span, x.ID)
	logger.Info("recipe prep task created")

	return x, nil
}

// CreateRecipePrepTask creates a recipe prep task.
func (q *Querier) CreateRecipePrepTask(ctx context.Context, input *types.RecipePrepTaskDatabaseCreationInput) (*types.RecipePrepTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, input.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	x, err := q.createRecipePrepTask(ctx, tx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe prep task created")

	return x, nil
}

//go:embed queries/recipe_prep_task_steps/create.sql
var createRecipePrepTaskStepQuery string

// createRecipePrepTaskStep creates a recipe prep task step.
func (q *Querier) createRecipePrepTaskStep(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.RecipePrepTaskStepDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.BelongsToRecipePrepTask,
		input.BelongsToRecipeStep,
		input.SatisfiesRecipeStep,
	}
	tracing.AttachRecipePrepTaskIDToSpan(span, input.BelongsToRecipePrepTask)

	// create the meal plan.
	if err := q.performWriteQuery(ctx, querier, "recipe prep task step creation", createRecipePrepTaskStepQuery, args); err != nil {
		q.rollbackTransaction(ctx, querier)
		return observability.PrepareAndLogError(err, logger, span, "creating recipe prep task step")
	}

	logger.Info("recipe prep task step created")

	return nil
}

//go:embed queries/recipe_prep_tasks/list_all_by_recipe.sql
var listRecipePrepTasksForRecipeQuery string

// GetRecipePrepTasksForRecipe gets a recipe prep task.
func (q *Querier) GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) (x []*types.RecipePrepTask, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	args := []interface{}{
		recipeID,
	}

	rows, getRowsErr := q.getRows(ctx, q.db, "recipe prep tasks list", listRecipePrepTasksForRecipeQuery, args)
	if getRowsErr != nil {
		return nil, observability.PrepareAndLogError(getRowsErr, logger, span, "executing recipe prep tasks list retrieval query")
	}

	x, scanErr := q.scanRecipePrepTasksWithSteps(ctx, rows)
	if scanErr != nil {
		return nil, observability.PrepareAndLogError(scanErr, logger, span, "scanning recipe prep tasks")
	}

	logger.Info("recipe prep tasks retrieved")

	return x, nil
}

//go:embed queries/recipe_prep_tasks/update.sql
var updateRecipePrepStepTaskQuery string

// UpdateRecipePrepTask updates a recipe prep task.
func (q *Querier) UpdateRecipePrepTask(ctx context.Context, updated *types.RecipePrepTask) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, updated.ID)

	updateRecipePrepStepTaskArgs := []interface{}{
		updated.Notes,
		updated.ExplicitStorageInstructions,
		updated.MinimumTimeBufferBeforeRecipeInSeconds,
		updated.MaximumTimeBufferBeforeRecipeInSeconds,
		updated.StorageType,
		updated.MinimumStorageTemperatureInCelsius,
		updated.MaximumStorageTemperatureInCelsius,
		updated.BelongsToRecipe,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe prep task update", updateRecipePrepStepTaskQuery, updateRecipePrepStepTaskArgs); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe prep task")
	}

	logger.Info("recipe prep task updated")

	return nil
}

//go:embed queries/recipe_prep_tasks/archive.sql
var archiveRecipePrepStepTaskQuery string

// ArchiveRecipePrepTask marks a recipe prep task as archived.
func (q *Querier) ArchiveRecipePrepTask(ctx context.Context, recipePrepTaskID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipePrepTaskID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	args := []interface{}{
		recipePrepTaskID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe prep task archive", archiveRecipePrepStepTaskQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe prep task")
	}

	return nil
}
