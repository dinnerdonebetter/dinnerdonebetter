package mealplanning

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/mealplanning/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

var (
	_ mealplanning.RecipePrepTaskDataManager = (*Querier)(nil)
)

// RecipePrepTaskExists checks if a recipe prep task exists.
func (q *Querier) RecipePrepTaskExists(ctx context.Context, recipeID, recipePrepTaskID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipePrepTaskID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)

	result, err := q.generatedQuerier.CheckRecipePrepTaskExistence(ctx, q.db, &generated.CheckRecipePrepTaskExistenceParams{
		RecipeID:         recipeID,
		RecipePrepTaskID: recipePrepTaskID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe prep task existence check")
	}

	logger.Info("recipe prep task existence retrieved")

	return result, nil
}

// GetRecipePrepTask fetches a recipe prep task.
func (q *Querier) GetRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (x *mealplanning.RecipePrepTask, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipePrepTaskID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)

	results, err := q.generatedQuerier.GetRecipePrepTask(ctx, q.db, recipePrepTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe prep task")
	}

	for _, result := range results {
		if x == nil {
			x = &mealplanning.RecipePrepTask{
				CreatedAt: result.CreatedAt,
				StorageTemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
				},
				TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
					Max: database.Uint32PointerFromNullInt32(result.MaximumTimeBufferBeforeRecipeInSeconds),
					Min: uint32(result.MinimumTimeBufferBeforeRecipeInSeconds),
				},
				ArchivedAt:                  database.TimePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt:               database.TimePointerFromNullTime(result.LastUpdatedAt),
				ID:                          result.ID,
				BelongsToRecipe:             result.BelongsToRecipe,
				ExplicitStorageInstructions: result.ExplicitStorageInstructions,
				Notes:                       result.Notes,
				Name:                        result.Name,
				Description:                 result.Description,
				TaskSteps:                   []*mealplanning.RecipePrepTaskStep{},
				Optional:                    result.Optional,
			}

			logger.WithValue("storage_type", result.StorageType).Info("storage type")

			if result.StorageType.Valid {
				x.StorageType = string(result.StorageType.StorageContainerType)
			}
		}

		x.TaskSteps = append(x.TaskSteps, &mealplanning.RecipePrepTaskStep{
			ID:                      result.TaskStepID,
			BelongsToRecipeStep:     result.TaskStepBelongsToRecipeStep,
			BelongsToRecipePrepTask: result.TaskStepBelongsToRecipePrepTask,
			SatisfiesRecipeStep:     result.TaskStepSatisfiesRecipeStep,
		})
	}

	if x == nil || x.ID == "" {
		return nil, sql.ErrNoRows
	}

	return x, nil
}

// createRecipePrepTask creates a recipe prep task.
func (q *Querier) createRecipePrepTask(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *mealplanning.RecipePrepTaskDatabaseCreationInput) (*mealplanning.RecipePrepTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanIDKey, input.ID)
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, input.ID)

	// create the recipe prep task.
	if err := q.generatedQuerier.CreateRecipePrepTask(ctx, querier, &generated.CreateRecipePrepTaskParams{
		ID:                                     input.ID,
		Name:                                   input.Name,
		Description:                            input.Description,
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		BelongsToRecipe:                        input.BelongsToRecipe,
		StorageType:                            generated.NullStorageContainerType{StorageContainerType: generated.StorageContainerType(input.StorageType), Valid: true},
		MinimumStorageTemperatureInCelsius:     database.NullStringFromFloat32Pointer(input.StorageTemperatureInCelsius.Min),
		MaximumStorageTemperatureInCelsius:     database.NullStringFromFloat32Pointer(input.StorageTemperatureInCelsius.Max),
		MaximumTimeBufferBeforeRecipeInSeconds: database.NullInt32FromUint32Pointer(input.TimeBufferBeforeRecipeInSeconds.Max),
		MinimumTimeBufferBeforeRecipeInSeconds: int32(input.TimeBufferBeforeRecipeInSeconds.Min),
		Optional:                               input.Optional,
	}); err != nil {
		q.RollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
	}

	x := &mealplanning.RecipePrepTask{
		CreatedAt:                   q.CurrentTime(),
		ID:                          input.ID,
		Name:                        input.Name,
		Description:                 input.Description,
		Notes:                       input.Notes,
		Optional:                    input.Optional,
		ExplicitStorageInstructions: input.ExplicitStorageInstructions,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: input.StorageTemperatureInCelsius.Max,
			Min: input.StorageTemperatureInCelsius.Min,
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Max: input.TimeBufferBeforeRecipeInSeconds.Max,
			Min: input.TimeBufferBeforeRecipeInSeconds.Min,
		},
		StorageType:     input.StorageType,
		BelongsToRecipe: input.BelongsToRecipe,
	}

	for _, recipePrepTaskStep := range input.TaskSteps {
		s, err := q.createRecipePrepTaskStep(ctx, querier, recipePrepTaskStep)
		if err != nil {
			q.RollbackTransaction(ctx, querier)
			return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
		}

		x.TaskSteps = append(x.TaskSteps, s)
	}

	logger.Info("recipe prep task created")

	return x, nil
}

// CreateRecipePrepTask creates a recipe prep task.
func (q *Querier) CreateRecipePrepTask(ctx context.Context, input *mealplanning.RecipePrepTaskDatabaseCreationInput) (*mealplanning.RecipePrepTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, database.ErrNilInputProvided
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

// createRecipePrepTaskStep creates a recipe prep task step.
func (q *Querier) createRecipePrepTaskStep(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *mealplanning.RecipePrepTaskStepDatabaseCreationInput) (*mealplanning.RecipePrepTaskStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, input.ID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, input.BelongsToRecipePrepTask)

	// create the meal plan.
	if err := q.generatedQuerier.CreateRecipePrepTaskStep(ctx, querier, &generated.CreateRecipePrepTaskStepParams{
		ID:                      input.ID,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}); err != nil {
		q.RollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task step")
	}

	x := &mealplanning.RecipePrepTaskStep{
		ID:                      input.ID,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}

	logger.Info("recipe prep task step created")

	return x, nil
}

// getRecipePrepTasksForRecipe gets a recipe prep task.
func (q *Querier) getRecipePrepTasksForRecipe(ctx context.Context, recipeID string) ([]*mealplanning.RecipePrepTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	results, err := q.generatedQuerier.ListAllRecipePrepTasksByRecipe(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe prep tasks list retrieval query")
	}

	x := []*mealplanning.RecipePrepTask{}

	var currentRecipePrepTask *mealplanning.RecipePrepTask
	for _, result := range results {
		prepTaskStep := &mealplanning.RecipePrepTaskStep{
			ID:                      result.TaskStepID,
			BelongsToRecipeStep:     result.TaskStepBelongsToRecipeStep,
			BelongsToRecipePrepTask: result.TaskStepBelongsToRecipePrepTask,
			SatisfiesRecipeStep:     result.TaskStepSatisfiesRecipeStep,
		}

		if currentRecipePrepTask != nil && currentRecipePrepTask.ID != result.ID {
			x = append(x, currentRecipePrepTask)
			currentRecipePrepTask = nil
		}

		if currentRecipePrepTask == nil {
			currentRecipePrepTask = &mealplanning.RecipePrepTask{
				CreatedAt:     result.CreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
				StorageTemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
				},
				TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
					Max: database.Uint32PointerFromNullInt32(result.MaximumTimeBufferBeforeRecipeInSeconds),
					Min: uint32(result.MinimumTimeBufferBeforeRecipeInSeconds),
				},
				ID:                          result.ID,
				StorageType:                 string(result.StorageType.StorageContainerType),
				BelongsToRecipe:             result.BelongsToRecipe,
				ExplicitStorageInstructions: result.ExplicitStorageInstructions,
				Notes:                       result.Notes,
				Name:                        result.Name,
				Description:                 result.Description,
				TaskSteps:                   []*mealplanning.RecipePrepTaskStep{},
				Optional:                    result.Optional,
			}
		}
		currentRecipePrepTask.TaskSteps = append(currentRecipePrepTask.TaskSteps, prepTaskStep)
	}

	if currentRecipePrepTask != nil && currentRecipePrepTask.ID != "" {
		x = append(x, currentRecipePrepTask)
	}

	return x, nil
}

// GetRecipePrepTasksForRecipe gets a recipe prep task.
func (q *Querier) GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) (x []*mealplanning.RecipePrepTask, err error) {
	return q.getRecipePrepTasksForRecipe(ctx, recipeID)
}

// UpdateRecipePrepTask updates a recipe prep task.
func (q *Querier) UpdateRecipePrepTask(ctx context.Context, updated *mealplanning.RecipePrepTask) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()
	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateRecipePrepTask(ctx, q.db, &generated.UpdateRecipePrepTaskParams{
		Name:                                   updated.Name,
		Description:                            updated.Description,
		Notes:                                  updated.Notes,
		Optional:                               updated.Optional,
		ExplicitStorageInstructions:            updated.ExplicitStorageInstructions,
		MinimumTimeBufferBeforeRecipeInSeconds: int32(updated.TimeBufferBeforeRecipeInSeconds.Min),
		MaximumTimeBufferBeforeRecipeInSeconds: database.NullInt32FromUint32Pointer(updated.TimeBufferBeforeRecipeInSeconds.Max),
		StorageType:                            generated.NullStorageContainerType{StorageContainerType: generated.StorageContainerType(updated.StorageType), Valid: true},
		MinimumStorageTemperatureInCelsius:     database.NullStringFromFloat32Pointer(updated.StorageTemperatureInCelsius.Min),
		MaximumStorageTemperatureInCelsius:     database.NullStringFromFloat32Pointer(updated.StorageTemperatureInCelsius.Max),
		BelongsToRecipe:                        updated.BelongsToRecipe,
		ID:                                     updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe prep task")
	}

	logger.Info("recipe prep task updated")

	return nil
}

// ArchiveRecipePrepTask marks a recipe prep task as archived.
func (q *Querier) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipePrepTaskID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)

	if _, err := q.generatedQuerier.ArchiveRecipePrepTask(ctx, q.db, recipePrepTaskID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe prep task")
	}

	return nil
}
