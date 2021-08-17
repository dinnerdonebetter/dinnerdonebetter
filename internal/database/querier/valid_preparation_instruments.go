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
	_ types.ValidPreparationInstrumentDataManager = (*SQLQuerier)(nil)
)

// scanValidPreparationInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a valid preparation instrument struct.
func (q *SQLQuerier) scanValidPreparationInstrument(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidPreparationInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidPreparationInstrument{}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.InstrumentID,
		&x.PreparationID,
		&x.Notes,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanValidPreparationInstruments takes some database rows and turns them into a slice of valid preparation instruments.
func (q *SQLQuerier) scanValidPreparationInstruments(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validPreparationInstruments []*types.ValidPreparationInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	validPreparationInstruments = []*types.ValidPreparationInstrument{}

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidPreparationInstrument(ctx, rows, includeCounts)
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

		validPreparationInstruments = append(validPreparationInstruments, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validPreparationInstruments, filteredCount, totalCount, nil
}

// ValidPreparationInstrumentExists fetches whether a valid preparation instrument exists from the database.
func (q *SQLQuerier) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID uint64) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validPreparationInstrumentID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	query, args := q.sqlQueryBuilder.BuildValidPreparationInstrumentExistsQuery(ctx, validPreparationInstrumentID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid preparation instrument existence check")
	}

	return result, nil
}

// GetValidPreparationInstrument fetches a valid preparation instrument from the database.
func (q *SQLQuerier) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID uint64) (*types.ValidPreparationInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validPreparationInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	query, args := q.sqlQueryBuilder.BuildGetValidPreparationInstrumentQuery(ctx, validPreparationInstrumentID)
	row := q.getOneRow(ctx, q.db, "validPreparationInstrument", query, args...)

	validPreparationInstrument, _, _, err := q.scanValidPreparationInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparation instrument")
	}

	return validPreparationInstrument, nil
}

// GetAllValidPreparationInstrumentsCount fetches the count of valid preparation instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetAllValidPreparationInstrumentsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllValidPreparationInstrumentsCountQuery(ctx), "fetching count of valid preparation instruments")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid preparation instruments")
	}

	return count, nil
}

// GetAllValidPreparationInstruments fetches a list of all valid preparation instruments in the database.
func (q *SQLQuerier) GetAllValidPreparationInstruments(ctx context.Context, results chan []*types.ValidPreparationInstrument, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllValidPreparationInstrumentsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of valid preparation instruments")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfValidPreparationInstrumentsQuery(ctx, begin, end)
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

			validPreparationInstruments, _, _, scanErr := q.scanValidPreparationInstruments(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- validPreparationInstruments
		}(beginID, endID)
	}

	return nil
}

// GetValidPreparationInstruments fetches a list of valid preparation instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.ValidPreparationInstrumentList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.ValidPreparationInstrumentList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetValidPreparationInstrumentsQuery(ctx, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "validPreparationInstruments", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid preparation instruments list retrieval query")
	}

	if x.ValidPreparationInstruments, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationInstruments(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparation instruments")
	}

	return x, nil
}

// GetValidPreparationInstrumentsWithIDs fetches valid preparation instruments from the database within a given set of IDs.
func (q *SQLQuerier) GetValidPreparationInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*types.ValidPreparationInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.sqlQueryBuilder.BuildGetValidPreparationInstrumentsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid preparation instruments with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid preparation instruments from database")
	}

	validPreparationInstruments, _, _, err := q.scanValidPreparationInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparation instruments")
	}

	return validPreparationInstruments, nil
}

// CreateValidPreparationInstrument creates a valid preparation instrument in the database.
func (q *SQLQuerier) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationInput, createdByUser uint64) (*types.ValidPreparationInstrument, error) {
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

	query, args := q.sqlQueryBuilder.BuildCreateValidPreparationInstrumentQuery(ctx, input)

	// create the valid preparation instrument.
	id, err := q.performWriteQuery(ctx, tx, false, "valid preparation instrument creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating valid preparation instrument")
	}

	x := &types.ValidPreparationInstrument{
		ID:            id,
		InstrumentID:  input.InstrumentID,
		PreparationID: input.PreparationID,
		Notes:         input.Notes,
		CreatedOn:     q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidPreparationInstrumentCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing valid preparation instrument creation audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachValidPreparationInstrumentIDToSpan(span, x.ID)
	logger.Info("valid preparation instrument created")

	return x, nil
}

// UpdateValidPreparationInstrument updates a particular valid preparation instrument. Note that UpdateValidPreparationInstrument expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateValidPreparationInstrument(ctx context.Context, updated *types.ValidPreparationInstrument, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, updated.ID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, updated.ID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateValidPreparationInstrumentQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid preparation instrument update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid preparation instrument")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidPreparationInstrumentUpdateEventEntry(changedByUser, updated.ID, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid preparation instrument update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid preparation instrument updated")

	return nil
}

// ArchiveValidPreparationInstrument archives a valid preparation instrument from the database by its ID.
func (q *SQLQuerier) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validPreparationInstrumentID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveValidPreparationInstrumentQuery(ctx, validPreparationInstrumentID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid preparation instrument archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid preparation instrument")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidPreparationInstrumentArchiveEventEntry(archivedBy, validPreparationInstrumentID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid preparation instrument archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid preparation instrument archived")

	return nil
}

// GetAuditLogEntriesForValidPreparationInstrument fetches a list of audit log entries from the database that relate to a given valid preparation instrument.
func (q *SQLQuerier) GetAuditLogEntriesForValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validPreparationInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForValidPreparationInstrumentQuery(ctx, validPreparationInstrumentID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for valid preparation instrument", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
