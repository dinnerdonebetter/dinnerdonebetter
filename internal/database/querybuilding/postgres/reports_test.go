package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildReportExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleReport := fakes.BuildFakeReport()

		expectedQuery := "SELECT EXISTS ( SELECT reports.id FROM reports WHERE reports.archived_on IS NULL AND reports.id = $1 )"
		expectedArgs := []interface{}{
			exampleReport.ID,
		}
		actualQuery, actualArgs := q.BuildReportExistsQuery(ctx, exampleReport.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleReport := fakes.BuildFakeReport()

		expectedQuery := "SELECT reports.id, reports.external_id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_account FROM reports WHERE reports.archived_on IS NULL AND reports.id = $1"
		expectedArgs := []interface{}{
			exampleReport.ID,
		}
		actualQuery, actualArgs := q.BuildGetReportQuery(ctx, exampleReport.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllReportsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(reports.id) FROM reports WHERE reports.archived_on IS NULL"
		actualQuery := q.BuildGetAllReportsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfReportsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT reports.id, reports.external_id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_account FROM reports WHERE reports.id > $1 AND reports.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfReportsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetReportsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT reports.id, reports.external_id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_account, (SELECT COUNT(reports.id) FROM reports WHERE reports.archived_on IS NULL) as total_count, (SELECT COUNT(reports.id) FROM reports WHERE reports.archived_on IS NULL AND reports.created_on > $1 AND reports.created_on < $2 AND reports.last_updated_on > $3 AND reports.last_updated_on < $4) as filtered_count FROM reports WHERE reports.archived_on IS NULL AND reports.created_on > $5 AND reports.created_on < $6 AND reports.last_updated_on > $7 AND reports.last_updated_on < $8 GROUP BY reports.id LIMIT 20 OFFSET 180"
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
		actualQuery, actualArgs := q.BuildGetReportsQuery(ctx, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetReportsWithIDsQuery(T *testing.T) {
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

		expectedQuery := "SELECT reports.id, reports.external_id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_account FROM (SELECT reports.id, reports.external_id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_account FROM reports JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS reports WHERE reports.archived_on IS NULL AND reports.belongs_to_account = $1 AND reports.id IN ($2,$3,$4)"
		expectedArgs := []interface{}{
			exampleAccountID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetReportsWithIDsQuery(ctx, exampleAccountID, defaultLimit, exampleIDs, true)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleReport := fakes.BuildFakeReport()
		exampleInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleReport.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO reports (external_id,report_type,concern,belongs_to_account) VALUES ($1,$2,$3,$4) RETURNING id"
		expectedArgs := []interface{}{
			exampleReport.ExternalID,
			exampleReport.ReportType,
			exampleReport.Concern,
			exampleReport.BelongsToAccount,
		}
		actualQuery, actualArgs := q.BuildCreateReportQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleReport := fakes.BuildFakeReport()

		expectedQuery := "UPDATE reports SET report_type = $1, concern = $2, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_account = $3 AND id = $4"
		expectedArgs := []interface{}{
			exampleReport.ReportType,
			exampleReport.Concern,
			exampleReport.BelongsToAccount,
			exampleReport.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateReportQuery(ctx, exampleReport)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleReportID := fakes.BuildFakeID()

		expectedQuery := "UPDATE reports SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"
		expectedArgs := []interface{}{
			exampleReportID,
		}
		actualQuery, actualArgs := q.BuildArchiveReportQuery(ctx, exampleReportID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleReport := fakes.BuildFakeReport()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'report_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleReport.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForReportQuery(ctx, exampleReport.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
