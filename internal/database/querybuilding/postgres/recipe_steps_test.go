package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildRecipeStepExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		expectedQuery := "SELECT EXISTS ( SELECT recipe_steps.id FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.archived_on IS NULL AND recipes.id = $3 )"
		expectedArgs := []interface{}{
			exampleRecipeID,
			exampleRecipeStep.ID,
			exampleRecipeID,
		}
		actualQuery, actualArgs := q.BuildRecipeStepExistsQuery(ctx, exampleRecipeID, exampleRecipeStep.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		expectedQuery := "SELECT recipe_steps.id, recipe_steps.external_id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.why, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.archived_on IS NULL AND recipes.id = $3"
		expectedArgs := []interface{}{
			exampleRecipeID,
			exampleRecipeStep.ID,
			exampleRecipeID,
		}
		actualQuery, actualArgs := q.BuildGetRecipeStepQuery(ctx, exampleRecipeID, exampleRecipeStep.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllRecipeStepsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_steps.id) FROM recipe_steps WHERE recipe_steps.archived_on IS NULL"
		actualQuery := q.BuildGetAllRecipeStepsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfRecipeStepsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipe_steps.id, recipe_steps.external_id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.why, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps WHERE recipe_steps.id > $1 AND recipe_steps.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfRecipeStepsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeStepsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_steps.id, recipe_steps.external_id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.why, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe, (SELECT COUNT(recipe_steps.id) FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.archived_on IS NULL AND recipes.id = $2) as total_count, (SELECT COUNT(recipe_steps.id) FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipes.archived_on IS NULL AND recipes.id = $4 AND recipe_steps.created_on > $5 AND recipe_steps.created_on < $6 AND recipe_steps.last_updated_on > $7 AND recipe_steps.last_updated_on < $8) as filtered_count FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $9 AND recipes.archived_on IS NULL AND recipes.id = $10 AND recipe_steps.created_on > $11 AND recipe_steps.created_on < $12 AND recipe_steps.last_updated_on > $13 AND recipe_steps.last_updated_on < $14 GROUP BY recipe_steps.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipeID,
			exampleRecipeID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			exampleRecipeID,
			exampleRecipeID,
			exampleRecipeID,
			exampleRecipeID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.BuildGetRecipeStepsQuery(ctx, exampleRecipeID, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeStepsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT recipe_steps.id, recipe_steps.external_id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.why, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM (SELECT recipe_steps.id, recipe_steps.external_id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.why, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS recipe_steps WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id IN ($2,$3,$4)"
		expectedArgs := []interface{}{
			exampleRecipeID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetRecipeStepsWithIDsQuery(ctx, exampleRecipeID, defaultLimit, exampleIDs)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleRecipeStep.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO recipe_steps (external_id,index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,why,recipe_id,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id"
		expectedArgs := []interface{}{
			exampleRecipeStep.ExternalID,
			exampleRecipeStep.Index,
			exampleRecipeStep.PreparationID,
			exampleRecipeStep.PrerequisiteStep,
			exampleRecipeStep.MinEstimatedTimeInSeconds,
			exampleRecipeStep.MaxEstimatedTimeInSeconds,
			exampleRecipeStep.TemperatureInCelsius,
			exampleRecipeStep.Notes,
			exampleRecipeStep.Why,
			exampleRecipeStep.RecipeID,
			exampleRecipeStep.BelongsToRecipe,
		}
		actualQuery, actualArgs := q.BuildCreateRecipeStepQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		expectedQuery := "UPDATE recipe_steps SET index = $1, preparation_id = $2, prerequisite_step = $3, min_estimated_time_in_seconds = $4, max_estimated_time_in_seconds = $5, temperature_in_celsius = $6, notes = $7, why = $8, recipe_id = $9, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $10 AND id = $11"
		expectedArgs := []interface{}{
			exampleRecipeStep.Index,
			exampleRecipeStep.PreparationID,
			exampleRecipeStep.PrerequisiteStep,
			exampleRecipeStep.MinEstimatedTimeInSeconds,
			exampleRecipeStep.MaxEstimatedTimeInSeconds,
			exampleRecipeStep.TemperatureInCelsius,
			exampleRecipeStep.Notes,
			exampleRecipeStep.Why,
			exampleRecipeStep.RecipeID,
			exampleRecipeStep.BelongsToRecipe,
			exampleRecipeStep.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateRecipeStepQuery(ctx, exampleRecipeStep)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		expectedQuery := "UPDATE recipe_steps SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2"
		expectedArgs := []interface{}{
			exampleRecipeID,
			exampleRecipeStepID,
		}
		actualQuery, actualArgs := q.BuildArchiveRecipeStepQuery(ctx, exampleRecipeID, exampleRecipeStepID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'recipe_step_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForRecipeStepQuery(ctx, exampleRecipeStep.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
