package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildRecipeStepIngredientExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		expectedQuery := "SELECT EXISTS ( SELECT recipe_step_ingredients.id FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5 )"
		expectedArgs := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}
		actualQuery, actualArgs := q.BuildRecipeStepIngredientExistsQuery(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		expectedQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.external_id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.name, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe_step, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5"
		expectedArgs := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}
		actualQuery, actualArgs := q.BuildGetRecipeStepIngredientQuery(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllRecipeStepIngredientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_step_ingredients.id) FROM recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL"
		actualQuery := q.BuildGetAllRecipeStepIngredientsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfRecipeStepIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.external_id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.name, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe_step, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients WHERE recipe_step_ingredients.id > $1 AND recipe_step_ingredients.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfRecipeStepIngredientsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeStepIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.external_id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.name, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe_step, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step, (SELECT COUNT(recipe_step_ingredients.id) FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.archived_on IS NULL AND recipes.id = $4) as total_count, (SELECT COUNT(recipe_step_ingredients.id) FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $5 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $6 AND recipe_steps.id = $7 AND recipes.archived_on IS NULL AND recipes.id = $8 AND recipe_step_ingredients.created_on > $9 AND recipe_step_ingredients.created_on < $10 AND recipe_step_ingredients.last_updated_on > $11 AND recipe_step_ingredients.last_updated_on < $12) as filtered_count FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $13 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $14 AND recipe_steps.id = $15 AND recipes.archived_on IS NULL AND recipes.id = $16 AND recipe_step_ingredients.created_on > $17 AND recipe_step_ingredients.created_on < $18 AND recipe_step_ingredients.last_updated_on > $19 AND recipe_step_ingredients.last_updated_on < $20 GROUP BY recipe_step_ingredients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipeStepID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			exampleRecipeStepID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.BuildGetRecipeStepIngredientsQuery(ctx, exampleRecipeID, exampleRecipeStepID, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeStepIngredientsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.external_id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.name, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe_step, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM (SELECT recipe_step_ingredients.id, recipe_step_ingredients.external_id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.name, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe_step, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id IN ($2,$3,$4)"
		expectedArgs := []interface{}{
			exampleRecipeStepID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetRecipeStepIngredientsWithIDsQuery(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleRecipeStepIngredient.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO recipe_step_ingredients (external_id,ingredient_id,name,quantity_type,quantity_value,quantity_notes,product_of_recipe_step,ingredient_notes,belongs_to_recipe_step) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id"
		expectedArgs := []interface{}{
			exampleRecipeStepIngredient.ExternalID,
			exampleRecipeStepIngredient.IngredientID,
			exampleRecipeStepIngredient.Name,
			exampleRecipeStepIngredient.QuantityType,
			exampleRecipeStepIngredient.QuantityValue,
			exampleRecipeStepIngredient.QuantityNotes,
			exampleRecipeStepIngredient.ProductOfRecipeStep,
			exampleRecipeStepIngredient.IngredientNotes,
			exampleRecipeStepIngredient.BelongsToRecipeStep,
		}
		actualQuery, actualArgs := q.BuildCreateRecipeStepIngredientQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		expectedQuery := "UPDATE recipe_step_ingredients SET ingredient_id = $1, name = $2, quantity_type = $3, quantity_value = $4, quantity_notes = $5, product_of_recipe_step = $6, ingredient_notes = $7, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $8 AND id = $9"
		expectedArgs := []interface{}{
			exampleRecipeStepIngredient.IngredientID,
			exampleRecipeStepIngredient.Name,
			exampleRecipeStepIngredient.QuantityType,
			exampleRecipeStepIngredient.QuantityValue,
			exampleRecipeStepIngredient.QuantityNotes,
			exampleRecipeStepIngredient.ProductOfRecipeStep,
			exampleRecipeStepIngredient.IngredientNotes,
			exampleRecipeStepIngredient.BelongsToRecipeStep,
			exampleRecipeStepIngredient.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateRecipeStepIngredientQuery(ctx, exampleRecipeStepIngredient)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientID := fakes.BuildFakeID()

		expectedQuery := "UPDATE recipe_step_ingredients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2"
		expectedArgs := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredientID,
		}
		actualQuery, actualArgs := q.BuildArchiveRecipeStepIngredientQuery(ctx, exampleRecipeStepID, exampleRecipeStepIngredientID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'recipe_step_ingredient_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleRecipeStepIngredient.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForRecipeStepIngredientQuery(ctx, exampleRecipeStepIngredient.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
