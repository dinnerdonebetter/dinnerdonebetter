package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildValidIngredientPreparationExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		expectedQuery := "SELECT EXISTS ( SELECT valid_ingredient_preparations.id FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.ID,
		}
		actualQuery, actualArgs := q.BuildValidIngredientPreparationExistsQuery(ctx, exampleValidIngredientPreparation.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.external_id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.id = $1"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.ID,
		}
		actualQuery, actualArgs := q.BuildGetValidIngredientPreparationQuery(ctx, exampleValidIngredientPreparation.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllValidIngredientPreparationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL"
		actualQuery := q.BuildGetAllValidIngredientPreparationsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfValidIngredientPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.external_id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.id > $1 AND valid_ingredient_preparations.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfValidIngredientPreparationsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidIngredientPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.external_id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on, (SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL) as total_count, (SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.created_on > $1 AND valid_ingredient_preparations.created_on < $2 AND valid_ingredient_preparations.last_updated_on > $3 AND valid_ingredient_preparations.last_updated_on < $4) as filtered_count FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.created_on > $5 AND valid_ingredient_preparations.created_on < $6 AND valid_ingredient_preparations.last_updated_on > $7 AND valid_ingredient_preparations.last_updated_on < $8 GROUP BY valid_ingredient_preparations.id ORDER BY valid_ingredient_preparations.id LIMIT 20 OFFSET 180"
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
		actualQuery, actualArgs := q.BuildGetValidIngredientPreparationsQuery(ctx, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidIngredientPreparationsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.external_id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM (SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.external_id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.id IN ($1,$2,$3)"
		expectedArgs := []interface{}{
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetValidIngredientPreparationsWithIDsQuery(ctx, defaultLimit, exampleIDs)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleValidIngredientPreparation.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO valid_ingredient_preparations (external_id,notes,valid_ingredient_id,valid_preparation_id) VALUES ($1,$2,$3,$4) RETURNING id"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.ExternalID,
			exampleValidIngredientPreparation.Notes,
			exampleValidIngredientPreparation.ValidIngredientID,
			exampleValidIngredientPreparation.ValidPreparationID,
		}
		actualQuery, actualArgs := q.BuildCreateValidIngredientPreparationQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		expectedQuery := "UPDATE valid_ingredient_preparations SET notes = $1, valid_ingredient_id = $2, valid_preparation_id = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.Notes,
			exampleValidIngredientPreparation.ValidIngredientID,
			exampleValidIngredientPreparation.ValidPreparationID,
			exampleValidIngredientPreparation.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateValidIngredientPreparationQuery(ctx, exampleValidIngredientPreparation)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredientPreparationID := fakes.BuildFakeID()

		expectedQuery := "UPDATE valid_ingredient_preparations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparationID,
		}
		actualQuery, actualArgs := q.BuildArchiveValidIngredientPreparationQuery(ctx, exampleValidIngredientPreparationID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'valid_ingredient_preparation_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForValidIngredientPreparationQuery(ctx, exampleValidIngredientPreparation.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
