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
	_ types.ValidIngredientPreparationDataManager = (*SQLQuerier)(nil)
)

// scanValidIngredientPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient preparation struct.
func (q *SQLQuerier) scanValidIngredientPreparation(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredientPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidIngredientPreparation{}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.Notes,
		&x.ValidIngredientID,
		&x.ValidPreparationID,
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

// scanValidIngredientPreparations takes some database rows and turns them into a slice of valid ingredient preparations.
func (q *SQLQuerier) scanValidIngredientPreparations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredientPreparations []*types.ValidIngredientPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidIngredientPreparation(ctx, rows, includeCounts)
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

		validIngredientPreparations = append(validIngredientPreparations, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validIngredientPreparations, filteredCount, totalCount, nil
}

// ValidIngredientPreparationExists fetches whether a valid ingredient preparation exists from the database.
func (q *SQLQuerier) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID uint64) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validIngredientPreparationID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	query, args := q.sqlQueryBuilder.BuildValidIngredientPreparationExistsQuery(ctx, validIngredientPreparationID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid ingredient preparation existence check")
	}

	return result, nil
}

// GetValidIngredientPreparation fetches a valid ingredient preparation from the database.
func (q *SQLQuerier) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) (*types.ValidIngredientPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validIngredientPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	query, args := q.sqlQueryBuilder.BuildGetValidIngredientPreparationQuery(ctx, validIngredientPreparationID)
	row := q.getOneRow(ctx, q.db, "validIngredientPreparation", query, args...)

	validIngredientPreparation, _, _, err := q.scanValidIngredientPreparation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient preparation")
	}

	return validIngredientPreparation, nil
}

// GetAllValidIngredientPreparationsCount fetches the count of valid ingredient preparations from the database that meet a particular filter.
func (q *SQLQuerier) GetAllValidIngredientPreparationsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllValidIngredientPreparationsCountQuery(ctx), "fetching count of valid ingredient preparations")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid ingredient preparations")
	}

	return count, nil
}

// GetAllValidIngredientPreparations fetches a list of all valid ingredient preparations in the database.
func (q *SQLQuerier) GetAllValidIngredientPreparations(ctx context.Context, results chan []*types.ValidIngredientPreparation, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllValidIngredientPreparationsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of valid ingredient preparations")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfValidIngredientPreparationsQuery(ctx, begin, end)
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

			validIngredientPreparations, _, _, scanErr := q.scanValidIngredientPreparations(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- validIngredientPreparations
		}(beginID, endID)
	}

	return nil
}

// GetValidIngredientPreparations fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *SQLQuerier) GetValidIngredientPreparations(ctx context.Context, filter *types.QueryFilter) (x *types.ValidIngredientPreparationList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.ValidIngredientPreparationList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetValidIngredientPreparationsQuery(ctx, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "validIngredientPreparations", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredient preparations list retrieval query")
	}

	if x.ValidIngredientPreparations, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientPreparations(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient preparations")
	}

	return x, nil
}

// GetValidIngredientPreparationsWithIDs fetches valid ingredient preparations from the database within a given set of IDs.
func (q *SQLQuerier) GetValidIngredientPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*types.ValidIngredientPreparation, error) {
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

	query, args := q.sqlQueryBuilder.BuildGetValidIngredientPreparationsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid ingredient preparations with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid ingredient preparations from database")
	}

	validIngredientPreparations, _, _, err := q.scanValidIngredientPreparations(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient preparations")
	}

	return validIngredientPreparations, nil
}

// CreateValidIngredientPreparation creates a valid ingredient preparation in the database.
func (q *SQLQuerier) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationInput, createdByUser uint64) (*types.ValidIngredientPreparation, error) {
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

	query, args := q.sqlQueryBuilder.BuildCreateValidIngredientPreparationQuery(ctx, input)

	// create the valid ingredient preparation.
	id, err := q.performWriteQuery(ctx, tx, false, "valid ingredient preparation creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating valid ingredient preparation")
	}

	x := &types.ValidIngredientPreparation{
		ID:                 id,
		Notes:              input.Notes,
		ValidIngredientID:  input.ValidIngredientID,
		ValidPreparationID: input.ValidPreparationID,
		CreatedOn:          q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidIngredientPreparationCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing valid ingredient preparation creation audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachValidIngredientPreparationIDToSpan(span, x.ID)
	logger.Info("valid ingredient preparation created")

	return x, nil
}

// UpdateValidIngredientPreparation updates a particular valid ingredient preparation. Note that UpdateValidIngredientPreparation expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateValidIngredientPreparation(ctx context.Context, updated *types.ValidIngredientPreparation, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientPreparationIDKey, updated.ID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, updated.ID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateValidIngredientPreparationQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid ingredient preparation update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid ingredient preparation")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidIngredientPreparationUpdateEventEntry(changedByUser, updated.ID, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid ingredient preparation update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid ingredient preparation updated")

	return nil
}

// ArchiveValidIngredientPreparation archives a valid ingredient preparation from the database by its ID.
func (q *SQLQuerier) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validIngredientPreparationID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveValidIngredientPreparationQuery(ctx, validIngredientPreparationID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid ingredient preparation archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid ingredient preparation")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidIngredientPreparationArchiveEventEntry(archivedBy, validIngredientPreparationID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid ingredient preparation archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid ingredient preparation archived")

	return nil
}

// GetAuditLogEntriesForValidIngredientPreparation fetches a list of audit log entries from the database that relate to a given valid ingredient preparation.
func (q *SQLQuerier) GetAuditLogEntriesForValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validIngredientPreparationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForValidIngredientPreparationQuery(ctx, validIngredientPreparationID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for valid ingredient preparation", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
