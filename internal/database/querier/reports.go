package querier

import (
	"context"
	"database/sql"
	"errors"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	database "gitlab.com/prixfixe/prixfixe/internal/database"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var (
	_ types.ReportDataManager = (*SQLQuerier)(nil)
)

// scanReport takes a database Scanner (i.e. *sql.Row) and scans the result into a report struct.
func (q *SQLQuerier) scanReport(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.Report, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.Report{}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.ReportType,
		&x.Concern,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToAccount,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanReports takes some database rows and turns them into a slice of reports.
func (q *SQLQuerier) scanReports(ctx context.Context, rows database.ResultIterator, includeCounts bool) (reports []*types.Report, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	reports = []*types.Report{}

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanReport(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		reports = append(reports, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return reports, filteredCount, totalCount, nil
}

// ReportExists fetches whether a report exists from the database.
func (q *SQLQuerier) ReportExists(ctx context.Context, reportID uint64) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if reportID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	query, args := q.sqlQueryBuilder.BuildReportExistsQuery(ctx, reportID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing report existence check")
	}

	return result, nil
}

// GetReport fetches a report from the database.
func (q *SQLQuerier) GetReport(ctx context.Context, reportID uint64) (*types.Report, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if reportID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	query, args := q.sqlQueryBuilder.BuildGetReportQuery(ctx, reportID)
	row := q.getOneRow(ctx, q.db, "report", query, args...)

	report, _, _, err := q.scanReport(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning report")
	}

	return report, nil
}

// GetAllReportsCount fetches the count of reports from the database that meet a particular filter.
func (q *SQLQuerier) GetAllReportsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllReportsCountQuery(ctx), "fetching count of reports")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of reports")
	}

	return count, nil
}

// GetAllReports fetches a list of all reports in the database.
func (q *SQLQuerier) GetAllReports(ctx context.Context, results chan []*types.Report, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllReportsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of reports")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfReportsQuery(ctx, begin, end)
			logger = logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, queryErr := q.db.Query(query, args...)
			if errors.Is(queryErr, sql.ErrNoRows) {
				return
			} else if queryErr != nil {
				logger.Error(queryErr, "querying for database rows")
				return
			}

			reports, _, _, scanErr := q.scanReports(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- reports
		}(beginID, endID)
	}

	return nil
}

// GetReports fetches a list of reports from the database that meet a particular filter.
func (q *SQLQuerier) GetReports(ctx context.Context, filter *types.QueryFilter) (x *types.ReportList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.ReportList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetReportsQuery(ctx, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "reports", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing reports list retrieval query")
	}

	if x.Reports, x.FilteredCount, x.TotalCount, err = q.scanReports(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning reports")
	}

	return x, nil
}

// GetReportsWithIDs fetches reports from the database within a given set of IDs.
func (q *SQLQuerier) GetReportsWithIDs(ctx context.Context, accountID uint64, limit uint8, ids []uint64) ([]*types.Report, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if accountID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachAccountIDToSpan(span, accountID)

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.sqlQueryBuilder.BuildGetReportsWithIDsQuery(ctx, accountID, limit, ids, false)

	rows, err := q.performReadQuery(ctx, q.db, "reports with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching reports from database")
	}

	reports, _, _, err := q.scanReports(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning reports")
	}

	return reports, nil
}

// CreateReport creates a report in the database.
func (q *SQLQuerier) CreateReport(ctx context.Context, input *types.ReportCreationInput, createdByUser uint64) (*types.Report, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if createdByUser == 0 {
		return nil, ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.RequesterIDKey, createdByUser)
	tracing.AttachRequestingUserIDToSpan(span, createdByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildCreateReportQuery(ctx, input)

	// create the report.
	id, err := q.performWriteQuery(ctx, tx, false, "report creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating report")
	}

	x := &types.Report{
		ID:               id,
		ReportType:       input.ReportType,
		Concern:          input.Concern,
		BelongsToAccount: input.BelongsToAccount,
		CreatedOn:        q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildReportCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing report creation audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachReportIDToSpan(span, x.ID)
	logger.Info("report created")

	return x, nil
}

// UpdateReport updates a particular report. Note that UpdateReport expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateReport(ctx context.Context, updated *types.Report, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.ReportIDKey, updated.ID)
	tracing.AttachReportIDToSpan(span, updated.ID)
	tracing.AttachAccountIDToSpan(span, updated.BelongsToAccount)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateReportQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "report update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating report")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildReportUpdateEventEntry(changedByUser, updated.ID, updated.BelongsToAccount, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing report update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("report updated")

	return nil
}

// ArchiveReport archives a report from the database by its ID.
func (q *SQLQuerier) ArchiveReport(ctx context.Context, reportID, accountID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if reportID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	if accountID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachAccountIDToSpan(span, accountID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveReportQuery(ctx, reportID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "report archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating report")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildReportArchiveEventEntry(archivedBy, accountID, reportID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing report archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("report archived")

	return nil
}

// GetAuditLogEntriesForReport fetches a list of audit log entries from the database that relate to a given report.
func (q *SQLQuerier) GetAuditLogEntriesForReport(ctx context.Context, reportID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if reportID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForReportQuery(ctx, reportID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for report", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
