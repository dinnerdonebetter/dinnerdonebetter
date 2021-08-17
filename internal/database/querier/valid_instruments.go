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
	_ types.ValidInstrumentDataManager = (*SQLQuerier)(nil)
)

// scanValidInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a valid instrument struct.
func (q *SQLQuerier) scanValidInstrument(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidInstrument{}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.Name,
		&x.Variant,
		&x.Description,
		&x.IconPath,
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

// scanValidInstruments takes some database rows and turns them into a slice of valid instruments.
func (q *SQLQuerier) scanValidInstruments(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validInstruments []*types.ValidInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	validInstruments = []*types.ValidInstrument{}

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidInstrument(ctx, rows, includeCounts)
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

		validInstruments = append(validInstruments, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validInstruments, filteredCount, totalCount, nil
}

// ValidInstrumentExists fetches whether a valid instrument exists from the database.
func (q *SQLQuerier) ValidInstrumentExists(ctx context.Context, validInstrumentID uint64) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validInstrumentID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	query, args := q.sqlQueryBuilder.BuildValidInstrumentExistsQuery(ctx, validInstrumentID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid instrument existence check")
	}

	return result, nil
}

// GetValidInstrument fetches a valid instrument from the database.
func (q *SQLQuerier) GetValidInstrument(ctx context.Context, validInstrumentID uint64) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	query, args := q.sqlQueryBuilder.BuildGetValidInstrumentQuery(ctx, validInstrumentID)
	row := q.getOneRow(ctx, q.db, "validInstrument", query, args...)

	validInstrument, _, _, err := q.scanValidInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid instrument")
	}

	return validInstrument, nil
}

// GetAllValidInstrumentsCount fetches the count of valid instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetAllValidInstrumentsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllValidInstrumentsCountQuery(ctx), "fetching count of valid instruments")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid instruments")
	}

	return count, nil
}

// GetAllValidInstruments fetches a list of all valid instruments in the database.
func (q *SQLQuerier) GetAllValidInstruments(ctx context.Context, results chan []*types.ValidInstrument, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllValidInstrumentsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of valid instruments")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfValidInstrumentsQuery(ctx, begin, end)
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

			validInstruments, _, _, scanErr := q.scanValidInstruments(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- validInstruments
		}(beginID, endID)
	}

	return nil
}

// GetValidInstruments fetches a list of valid instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.ValidInstrumentList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.ValidInstrumentList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetValidInstrumentsQuery(ctx, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "validInstruments", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid instruments list retrieval query")
	}

	if x.ValidInstruments, x.FilteredCount, x.TotalCount, err = q.scanValidInstruments(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid instruments")
	}

	return x, nil
}

// GetValidInstrumentsWithIDs fetches valid instruments from the database within a given set of IDs.
func (q *SQLQuerier) GetValidInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*types.ValidInstrument, error) {
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

	query, args := q.sqlQueryBuilder.BuildGetValidInstrumentsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid instruments with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid instruments from database")
	}

	validInstruments, _, _, err := q.scanValidInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid instruments")
	}

	return validInstruments, nil
}

// CreateValidInstrument creates a valid instrument in the database.
func (q *SQLQuerier) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentCreationInput, createdByUser uint64) (*types.ValidInstrument, error) {
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

	query, args := q.sqlQueryBuilder.BuildCreateValidInstrumentQuery(ctx, input)

	// create the valid instrument.
	id, err := q.performWriteQuery(ctx, tx, false, "valid instrument creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating valid instrument")
	}

	x := &types.ValidInstrument{
		ID:          id,
		Name:        input.Name,
		Variant:     input.Variant,
		Description: input.Description,
		IconPath:    input.IconPath,
		CreatedOn:   q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidInstrumentCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing valid instrument creation audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachValidInstrumentIDToSpan(span, x.ID)
	logger.Info("valid instrument created")

	return x, nil
}

// UpdateValidInstrument updates a particular valid instrument. Note that UpdateValidInstrument expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, updated.ID)
	tracing.AttachValidInstrumentIDToSpan(span, updated.ID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateValidInstrumentQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid instrument update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid instrument")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidInstrumentUpdateEventEntry(changedByUser, updated.ID, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid instrument update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid instrument updated")

	return nil
}

// ArchiveValidInstrument archives a valid instrument from the database by its ID.
func (q *SQLQuerier) ArchiveValidInstrument(ctx context.Context, validInstrumentID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validInstrumentID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveValidInstrumentQuery(ctx, validInstrumentID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid instrument archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid instrument")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidInstrumentArchiveEventEntry(archivedBy, validInstrumentID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid instrument archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid instrument archived")

	return nil
}

// GetAuditLogEntriesForValidInstrument fetches a list of audit log entries from the database that relate to a given valid instrument.
func (q *SQLQuerier) GetAuditLogEntriesForValidInstrument(ctx context.Context, validInstrumentID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validInstrumentID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForValidInstrumentQuery(ctx, validInstrumentID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for valid instrument", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
