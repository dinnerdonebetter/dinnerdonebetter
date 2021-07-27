package postgres

import (
	"context"
	"fmt"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildValidPreparationExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		expectedQuery := "SELECT EXISTS ( SELECT valid_preparations.id FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := q.BuildValidPreparationExistsQuery(ctx, exampleValidPreparation.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidPreparationIDForNameQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		expectedQuery := "SELECT valid_preparations.id FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.name = $1"
		expectedArgs := []interface{}{
			exampleValidPreparation.Name,
		}
		actualQuery, actualArgs := q.BuildGetValidPreparationIDForNameQuery(ctx, exampleValidPreparation.Name)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildSearchForValidPreparationByNameQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.external_id, valid_preparations.name, valid_preparations.description, valid_preparations.icon_path, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.name ILIKE $1 AND valid_preparations.archived_on IS NULL"
		expectedArgs := []interface{}{
			fmt.Sprintf("%s%%", exampleValidPreparation.Name),
		}
		actualQuery, actualArgs := q.BuildSearchForValidPreparationByNameQuery(ctx, exampleValidPreparation.Name)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.external_id, valid_preparations.name, valid_preparations.description, valid_preparations.icon_path, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.id = $1"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := q.BuildGetValidPreparationQuery(ctx, exampleValidPreparation.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllValidPreparationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL"
		actualQuery := q.BuildGetAllValidPreparationsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfValidPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.external_id, valid_preparations.name, valid_preparations.description, valid_preparations.icon_path, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.id > $1 AND valid_preparations.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfValidPreparationsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.external_id, valid_preparations.name, valid_preparations.description, valid_preparations.icon_path, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on, (SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL) as total_count, (SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.created_on > $1 AND valid_preparations.created_on < $2 AND valid_preparations.last_updated_on > $3 AND valid_preparations.last_updated_on < $4) as filtered_count FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.created_on > $5 AND valid_preparations.created_on < $6 AND valid_preparations.last_updated_on > $7 AND valid_preparations.last_updated_on < $8 GROUP BY valid_preparations.id ORDER BY valid_preparations.id LIMIT 20 OFFSET 180"
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
		actualQuery, actualArgs := q.BuildGetValidPreparationsQuery(ctx, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidPreparationsWithIDsQuery(T *testing.T) {
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

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.external_id, valid_preparations.name, valid_preparations.description, valid_preparations.icon_path, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM (SELECT valid_preparations.id, valid_preparations.external_id, valid_preparations.name, valid_preparations.description, valid_preparations.icon_path, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.id IN ($1,$2,$3)"
		expectedArgs := []interface{}{
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetValidPreparationsWithIDsQuery(ctx, defaultLimit, exampleIDs)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateValidPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleValidPreparation.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO valid_preparations (external_id,name,description,icon_path) VALUES ($1,$2,$3,$4) RETURNING id"
		expectedArgs := []interface{}{
			exampleValidPreparation.ExternalID,
			exampleValidPreparation.Name,
			exampleValidPreparation.Description,
			exampleValidPreparation.IconPath,
		}
		actualQuery, actualArgs := q.BuildCreateValidPreparationQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateValidPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		expectedQuery := "UPDATE valid_preparations SET name = $1, description = $2, icon_path = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4"
		expectedArgs := []interface{}{
			exampleValidPreparation.Name,
			exampleValidPreparation.Description,
			exampleValidPreparation.IconPath,
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateValidPreparationQuery(ctx, exampleValidPreparation)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveValidPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparationID := fakes.BuildFakeID()

		expectedQuery := "UPDATE valid_preparations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"
		expectedArgs := []interface{}{
			exampleValidPreparationID,
		}
		actualQuery, actualArgs := q.BuildArchiveValidPreparationQuery(ctx, exampleValidPreparationID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForValidPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'valid_preparation_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForValidPreparationQuery(ctx, exampleValidPreparation.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
