package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildValidPreparationInstrumentExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		expectedQuery := "SELECT EXISTS ( SELECT valid_preparation_instruments.id FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_on IS NULL AND valid_preparation_instruments.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidPreparationInstrument.ID,
		}
		actualQuery, actualArgs := q.BuildValidPreparationInstrumentExistsQuery(ctx, exampleValidPreparationInstrument.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		expectedQuery := "SELECT valid_preparation_instruments.id, valid_preparation_instruments.external_id, valid_preparation_instruments.instrument_id, valid_preparation_instruments.preparation_id, valid_preparation_instruments.notes, valid_preparation_instruments.created_on, valid_preparation_instruments.last_updated_on, valid_preparation_instruments.archived_on FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_on IS NULL AND valid_preparation_instruments.id = $1"
		expectedArgs := []interface{}{
			exampleValidPreparationInstrument.ID,
		}
		actualQuery, actualArgs := q.BuildGetValidPreparationInstrumentQuery(ctx, exampleValidPreparationInstrument.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllValidPreparationInstrumentsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_preparation_instruments.id) FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_on IS NULL"
		actualQuery := q.BuildGetAllValidPreparationInstrumentsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfValidPreparationInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT valid_preparation_instruments.id, valid_preparation_instruments.external_id, valid_preparation_instruments.instrument_id, valid_preparation_instruments.preparation_id, valid_preparation_instruments.notes, valid_preparation_instruments.created_on, valid_preparation_instruments.last_updated_on, valid_preparation_instruments.archived_on FROM valid_preparation_instruments WHERE valid_preparation_instruments.id > $1 AND valid_preparation_instruments.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfValidPreparationInstrumentsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidPreparationInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_preparation_instruments.id, valid_preparation_instruments.external_id, valid_preparation_instruments.instrument_id, valid_preparation_instruments.preparation_id, valid_preparation_instruments.notes, valid_preparation_instruments.created_on, valid_preparation_instruments.last_updated_on, valid_preparation_instruments.archived_on, (SELECT COUNT(valid_preparation_instruments.id) FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_on IS NULL) as total_count, (SELECT COUNT(valid_preparation_instruments.id) FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_on IS NULL AND valid_preparation_instruments.created_on > $1 AND valid_preparation_instruments.created_on < $2 AND valid_preparation_instruments.last_updated_on > $3 AND valid_preparation_instruments.last_updated_on < $4) as filtered_count FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_on IS NULL AND valid_preparation_instruments.created_on > $5 AND valid_preparation_instruments.created_on < $6 AND valid_preparation_instruments.last_updated_on > $7 AND valid_preparation_instruments.last_updated_on < $8 GROUP BY valid_preparation_instruments.id ORDER BY valid_preparation_instruments.id LIMIT 20 OFFSET 180"
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
		actualQuery, actualArgs := q.BuildGetValidPreparationInstrumentsQuery(ctx, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetValidPreparationInstrumentsWithIDsQuery(T *testing.T) {
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

		expectedQuery := "SELECT valid_preparation_instruments.id, valid_preparation_instruments.external_id, valid_preparation_instruments.instrument_id, valid_preparation_instruments.preparation_id, valid_preparation_instruments.notes, valid_preparation_instruments.created_on, valid_preparation_instruments.last_updated_on, valid_preparation_instruments.archived_on FROM (SELECT valid_preparation_instruments.id, valid_preparation_instruments.external_id, valid_preparation_instruments.instrument_id, valid_preparation_instruments.preparation_id, valid_preparation_instruments.notes, valid_preparation_instruments.created_on, valid_preparation_instruments.last_updated_on, valid_preparation_instruments.archived_on FROM valid_preparation_instruments JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS valid_preparation_instruments WHERE valid_preparation_instruments.archived_on IS NULL AND valid_preparation_instruments.id IN ($1,$2,$3)"
		expectedArgs := []interface{}{
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetValidPreparationInstrumentsWithIDsQuery(ctx, defaultLimit, exampleIDs)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateValidPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleValidPreparationInstrument.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO valid_preparation_instruments (external_id,instrument_id,preparation_id,notes) VALUES ($1,$2,$3,$4) RETURNING id"
		expectedArgs := []interface{}{
			exampleValidPreparationInstrument.ExternalID,
			exampleValidPreparationInstrument.InstrumentID,
			exampleValidPreparationInstrument.PreparationID,
			exampleValidPreparationInstrument.Notes,
		}
		actualQuery, actualArgs := q.BuildCreateValidPreparationInstrumentQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateValidPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		expectedQuery := "UPDATE valid_preparation_instruments SET instrument_id = $1, preparation_id = $2, notes = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4"
		expectedArgs := []interface{}{
			exampleValidPreparationInstrument.InstrumentID,
			exampleValidPreparationInstrument.PreparationID,
			exampleValidPreparationInstrument.Notes,
			exampleValidPreparationInstrument.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateValidPreparationInstrumentQuery(ctx, exampleValidPreparationInstrument)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveValidPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparationInstrumentID := fakes.BuildFakeID()

		expectedQuery := "UPDATE valid_preparation_instruments SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"
		expectedArgs := []interface{}{
			exampleValidPreparationInstrumentID,
		}
		actualQuery, actualArgs := q.BuildArchiveValidPreparationInstrumentQuery(ctx, exampleValidPreparationInstrumentID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForValidPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'valid_preparation_instrument_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleValidPreparationInstrument.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForValidPreparationInstrumentQuery(ctx, exampleValidPreparationInstrument.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
