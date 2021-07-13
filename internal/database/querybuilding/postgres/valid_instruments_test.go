package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildValidInstrumentExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		expectedQuery := "SELECT EXISTS ( SELECT valid_instruments.id FROM valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidInstrument.ID,
		}
		actualQuery, actualArgs := q.BuildValidInstrumentExistsQuery(ctx, exampleValidInstrument.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		expectedQuery := "SELECT valid_instruments.id, valid_instruments.external_id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon_path, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.id = $1"
		expectedArgs := []interface{}{
			exampleValidInstrument.ID,
		}
		actualQuery, actualArgs := q.BuildGetValidInstrumentQuery(ctx, exampleValidInstrument.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllValidInstrumentsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL"
		actualQuery := q.BuildGetAllValidInstrumentsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfValidInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT valid_instruments.id, valid_instruments.external_id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon_path, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments WHERE valid_instruments.id > $1 AND valid_instruments.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfValidInstrumentsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_instruments.id, valid_instruments.external_id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon_path, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on, (SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL) as total_count, (SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.created_on > $1 AND valid_instruments.created_on < $2 AND valid_instruments.last_updated_on > $3 AND valid_instruments.last_updated_on < $4) as filtered_count FROM valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.created_on > $5 AND valid_instruments.created_on < $6 AND valid_instruments.last_updated_on > $7 AND valid_instruments.last_updated_on < $8 GROUP BY valid_instruments.id LIMIT 20 OFFSET 180"
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
		actualQuery, actualArgs := q.BuildGetValidInstrumentsQuery(ctx, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidInstrumentsWithIDsQuery(T *testing.T) {
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

		expectedQuery := "SELECT valid_instruments.id, valid_instruments.external_id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon_path, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM (SELECT valid_instruments.id, valid_instruments.external_id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon_path, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.id IN ($1,$2,$3)"
		expectedArgs := []interface{}{
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetValidInstrumentsWithIDsQuery(ctx, defaultLimit, exampleIDs)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateValidInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		exampleInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleValidInstrument.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO valid_instruments (external_id,name,variant,description,icon_path) VALUES ($1,$2,$3,$4,$5) RETURNING id"
		expectedArgs := []interface{}{
			exampleValidInstrument.ExternalID,
			exampleValidInstrument.Name,
			exampleValidInstrument.Variant,
			exampleValidInstrument.Description,
			exampleValidInstrument.IconPath,
		}
		actualQuery, actualArgs := q.BuildCreateValidInstrumentQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateValidInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		expectedQuery := "UPDATE valid_instruments SET name = $1, variant = $2, description = $3, icon_path = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $5"
		expectedArgs := []interface{}{
			exampleValidInstrument.Name,
			exampleValidInstrument.Variant,
			exampleValidInstrument.Description,
			exampleValidInstrument.IconPath,
			exampleValidInstrument.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateValidInstrumentQuery(ctx, exampleValidInstrument)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveValidInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidInstrumentID := fakes.BuildFakeID()

		expectedQuery := "UPDATE valid_instruments SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"
		expectedArgs := []interface{}{
			exampleValidInstrumentID,
		}
		actualQuery, actualArgs := q.BuildArchiveValidInstrumentQuery(ctx, exampleValidInstrumentID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForValidInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'valid_instrument_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleValidInstrument.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForValidInstrumentQuery(ctx, exampleValidInstrument.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
