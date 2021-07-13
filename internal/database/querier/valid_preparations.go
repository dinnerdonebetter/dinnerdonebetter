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
	_ types.ValidPreparationDataManager = (*SQLQuerier)(nil)
)

// scanValidPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a valid preparation struct.
func (q *SQLQuerier) scanValidPreparation(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidPreparation{}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.Name,
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

// scanValidPreparations takes some database rows and turns them into a slice of valid preparations.
func (q *SQLQuerier) scanValidPreparations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validPreparations []*types.ValidPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidPreparation(ctx, rows, includeCounts)
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

		validPreparations = append(validPreparations, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validPreparations, filteredCount, totalCount, nil
}

// ValidPreparationExists fetches whether a valid preparation exists from the database.
func (q *SQLQuerier) ValidPreparationExists(ctx context.Context, validPreparationID uint64) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validPreparationID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	query, args := q.sqlQueryBuilder.BuildValidPreparationExistsQuery(ctx, validPreparationID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid preparation existence check")
	}

	return result, nil
}

// GetValidPreparation fetches a valid preparation from the database.
func (q *SQLQuerier) GetValidPreparation(ctx context.Context, validPreparationID uint64) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	query, args := q.sqlQueryBuilder.BuildGetValidPreparationQuery(ctx, validPreparationID)
	row := q.getOneRow(ctx, q.db, "validPreparation", query, args...)

	validPreparation, _, _, err := q.scanValidPreparation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparation")
	}

	return validPreparation, nil
}

// GetAllValidPreparationsCount fetches the count of valid preparations from the database that meet a particular filter.
func (q *SQLQuerier) GetAllValidPreparationsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllValidPreparationsCountQuery(ctx), "fetching count of valid preparations")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid preparations")
	}

	return count, nil
}

// GetAllValidPreparations fetches a list of all valid preparations in the database.
func (q *SQLQuerier) GetAllValidPreparations(ctx context.Context, results chan []*types.ValidPreparation, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllValidPreparationsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of valid preparations")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfValidPreparationsQuery(ctx, begin, end)
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

			validPreparations, _, _, scanErr := q.scanValidPreparations(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- validPreparations
		}(beginID, endID)
	}

	return nil
}

// GetValidPreparations fetches a list of valid preparations from the database that meet a particular filter.
func (q *SQLQuerier) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (x *types.ValidPreparationList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.ValidPreparationList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetValidPreparationsQuery(ctx, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "validPreparations", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid preparations list retrieval query")
	}

	if x.ValidPreparations, x.FilteredCount, x.TotalCount, err = q.scanValidPreparations(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparations")
	}

	return x, nil
}

// GetValidPreparationsWithIDs fetches valid preparations from the database within a given set of IDs.
func (q *SQLQuerier) GetValidPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*types.ValidPreparation, error) {
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

	query, args := q.sqlQueryBuilder.BuildGetValidPreparationsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid preparations with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid preparations from database")
	}

	validPreparations, _, _, err := q.scanValidPreparations(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparations")
	}

	return validPreparations, nil
}

// CreateValidPreparation creates a valid preparation in the database.
func (q *SQLQuerier) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationCreationInput, createdByUser uint64) (*types.ValidPreparation, error) {
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

	query, args := q.sqlQueryBuilder.BuildCreateValidPreparationQuery(ctx, input)

	// create the valid preparation.
	id, err := q.performWriteQuery(ctx, tx, false, "valid preparation creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating valid preparation")
	}

	x := &types.ValidPreparation{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		IconPath:    input.IconPath,
		CreatedOn:   q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidPreparationCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing valid preparation creation audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachValidPreparationIDToSpan(span, x.ID)
	logger.Info("valid preparation created")

	return x, nil
}

// UpdateValidPreparation updates a particular valid preparation. Note that UpdateValidPreparation expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateValidPreparation(ctx context.Context, updated *types.ValidPreparation, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationIDKey, updated.ID)
	tracing.AttachValidPreparationIDToSpan(span, updated.ID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateValidPreparationQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid preparation update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid preparation")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidPreparationUpdateEventEntry(changedByUser, updated.ID, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid preparation update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid preparation updated")

	return nil
}

// ArchiveValidPreparation archives a valid preparation from the database by its ID.
func (q *SQLQuerier) ArchiveValidPreparation(ctx context.Context, validPreparationID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validPreparationID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveValidPreparationQuery(ctx, validPreparationID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid preparation archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid preparation")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidPreparationArchiveEventEntry(archivedBy, validPreparationID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid preparation archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid preparation archived")

	return nil
}

// GetAuditLogEntriesForValidPreparation fetches a list of audit log entries from the database that relate to a given valid preparation.
func (q *SQLQuerier) GetAuditLogEntriesForValidPreparation(ctx context.Context, validPreparationID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForValidPreparationQuery(ctx, validPreparationID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for valid preparation", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
