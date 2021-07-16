package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildRecipeStepProductExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		expectedQuery := "SELECT EXISTS ( SELECT recipe_step_products.id FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5 )"
		expectedArgs := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}
		actualQuery, actualArgs := q.BuildRecipeStepProductExistsQuery(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		expectedQuery := "SELECT recipe_step_products.id, recipe_step_products.external_id, recipe_step_products.name, recipe_step_products.quantity_type, recipe_step_products.quantity_value, recipe_step_products.quantity_notes, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5"
		expectedArgs := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}
		actualQuery, actualArgs := q.BuildGetRecipeStepProductQuery(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllRecipeStepProductsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_step_products.id) FROM recipe_step_products WHERE recipe_step_products.archived_on IS NULL"
		actualQuery := q.BuildGetAllRecipeStepProductsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfRecipeStepProductsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipe_step_products.id, recipe_step_products.external_id, recipe_step_products.name, recipe_step_products.quantity_type, recipe_step_products.quantity_value, recipe_step_products.quantity_notes, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products WHERE recipe_step_products.id > $1 AND recipe_step_products.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfRecipeStepProductsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeStepProductsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_step_products.id, recipe_step_products.external_id, recipe_step_products.name, recipe_step_products.quantity_type, recipe_step_products.quantity_value, recipe_step_products.quantity_notes, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step, (SELECT COUNT(recipe_step_products.id) FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.archived_on IS NULL AND recipes.id = $4) as total_count, (SELECT COUNT(recipe_step_products.id) FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $5 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $6 AND recipe_steps.id = $7 AND recipes.archived_on IS NULL AND recipes.id = $8 AND recipe_step_products.created_on > $9 AND recipe_step_products.created_on < $10 AND recipe_step_products.last_updated_on > $11 AND recipe_step_products.last_updated_on < $12) as filtered_count FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $13 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $14 AND recipe_steps.id = $15 AND recipes.archived_on IS NULL AND recipes.id = $16 AND recipe_step_products.created_on > $17 AND recipe_step_products.created_on < $18 AND recipe_step_products.last_updated_on > $19 AND recipe_step_products.last_updated_on < $20 GROUP BY recipe_step_products.id ORDER BY recipe_step_products.id LIMIT 20 OFFSET 180"
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
		actualQuery, actualArgs := q.BuildGetRecipeStepProductsQuery(ctx, exampleRecipeID, exampleRecipeStepID, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeStepProductsWithIDsQuery(T *testing.T) {
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

		expectedQuery := "SELECT recipe_step_products.id, recipe_step_products.external_id, recipe_step_products.name, recipe_step_products.quantity_type, recipe_step_products.quantity_value, recipe_step_products.quantity_notes, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM (SELECT recipe_step_products.id, recipe_step_products.external_id, recipe_step_products.name, recipe_step_products.quantity_type, recipe_step_products.quantity_value, recipe_step_products.quantity_notes, recipe_step_products.recipe_step_id, recipe_step_products.created_on, recipe_step_products.last_updated_on, recipe_step_products.archived_on, recipe_step_products.belongs_to_recipe_step FROM recipe_step_products JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS recipe_step_products WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id IN ($2,$3,$4)"
		expectedArgs := []interface{}{
			exampleRecipeStepID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetRecipeStepProductsWithIDsQuery(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleRecipeStepProduct.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO recipe_step_products (external_id,name,quantity_type,quantity_value,quantity_notes,recipe_step_id,belongs_to_recipe_step) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id"
		expectedArgs := []interface{}{
			exampleRecipeStepProduct.ExternalID,
			exampleRecipeStepProduct.Name,
			exampleRecipeStepProduct.QuantityType,
			exampleRecipeStepProduct.QuantityValue,
			exampleRecipeStepProduct.QuantityNotes,
			exampleRecipeStepProduct.RecipeStepID,
			exampleRecipeStepProduct.BelongsToRecipeStep,
		}
		actualQuery, actualArgs := q.BuildCreateRecipeStepProductQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		expectedQuery := "UPDATE recipe_step_products SET name = $1, quantity_type = $2, quantity_value = $3, quantity_notes = $4, recipe_step_id = $5, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $6 AND id = $7"
		expectedArgs := []interface{}{
			exampleRecipeStepProduct.Name,
			exampleRecipeStepProduct.QuantityType,
			exampleRecipeStepProduct.QuantityValue,
			exampleRecipeStepProduct.QuantityNotes,
			exampleRecipeStepProduct.RecipeStepID,
			exampleRecipeStepProduct.BelongsToRecipeStep,
			exampleRecipeStepProduct.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateRecipeStepProductQuery(ctx, exampleRecipeStepProduct)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProductID := fakes.BuildFakeID()

		expectedQuery := "UPDATE recipe_step_products SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2"
		expectedArgs := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepProductID,
		}
		actualQuery, actualArgs := q.BuildArchiveRecipeStepProductQuery(ctx, exampleRecipeStepID, exampleRecipeStepProductID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForRecipeStepProductQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'recipe_step_product_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleRecipeStepProduct.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForRecipeStepProductQuery(ctx, exampleRecipeStepProduct.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
