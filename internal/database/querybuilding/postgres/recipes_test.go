package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildRecipeExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipe := fakes.BuildFakeRecipe()

		expectedQuery := "SELECT EXISTS ( SELECT recipes.id FROM recipes WHERE recipes.archived_on IS NULL AND recipes.id = $1 )"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := q.BuildRecipeExistsQuery(ctx, exampleRecipe.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipeQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipe := fakes.BuildFakeRecipe()

		expectedQuery := "SELECT recipes.id, recipes.external_id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_account FROM recipes WHERE recipes.archived_on IS NULL AND recipes.id = $1"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := q.BuildGetRecipeQuery(ctx, exampleRecipe.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllRecipesCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipes.id) FROM recipes WHERE recipes.archived_on IS NULL"
		actualQuery := q.BuildGetAllRecipesCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfRecipesQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipes.id, recipes.external_id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_account FROM recipes WHERE recipes.id > $1 AND recipes.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfRecipesQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipesQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipes.id, recipes.external_id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_account, (SELECT COUNT(recipes.id) FROM recipes WHERE recipes.archived_on IS NULL) as total_count, (SELECT COUNT(recipes.id) FROM recipes WHERE recipes.archived_on IS NULL AND recipes.created_on > $1 AND recipes.created_on < $2 AND recipes.last_updated_on > $3 AND recipes.last_updated_on < $4) as filtered_count FROM recipes WHERE recipes.archived_on IS NULL AND recipes.created_on > $5 AND recipes.created_on < $6 AND recipes.last_updated_on > $7 AND recipes.last_updated_on < $8 GROUP BY recipes.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.BuildGetRecipesQuery(ctx, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetRecipesWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleAccountID := fakes.BuildFakeID()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT recipes.id, recipes.external_id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_account FROM (SELECT recipes.id, recipes.external_id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_account FROM recipes JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS recipes WHERE recipes.archived_on IS NULL AND recipes.belongs_to_account = $1 AND recipes.id IN ($2,$3,$4)"
		expectedArgs := []interface{}{
			exampleAccountID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetRecipesWithIDsQuery(ctx, exampleAccountID, defaultLimit, exampleIDs, true)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateRecipeQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleRecipe.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO recipes (external_id,name,source,description,inspired_by_recipe_id,belongs_to_account) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"
		expectedArgs := []interface{}{
			exampleRecipe.ExternalID,
			exampleRecipe.Name,
			exampleRecipe.Source,
			exampleRecipe.Description,
			exampleRecipe.InspiredByRecipeID,
			exampleRecipe.BelongsToAccount,
		}
		actualQuery, actualArgs := q.BuildCreateRecipeQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateRecipeQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipe := fakes.BuildFakeRecipe()

		expectedQuery := "UPDATE recipes SET name = $1, source = $2, description = $3, inspired_by_recipe_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_account = $5 AND id = $6"
		expectedArgs := []interface{}{
			exampleRecipe.Name,
			exampleRecipe.Source,
			exampleRecipe.Description,
			exampleRecipe.InspiredByRecipeID,
			exampleRecipe.BelongsToAccount,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateRecipeQuery(ctx, exampleRecipe)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveRecipeQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()

		expectedQuery := "UPDATE recipes SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"
		expectedArgs := []interface{}{
			exampleRecipeID,
		}
		actualQuery, actualArgs := q.BuildArchiveRecipeQuery(ctx, exampleRecipeID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForRecipeQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleRecipe := fakes.BuildFakeRecipe()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'recipe_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForRecipeQuery(ctx, exampleRecipe.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
