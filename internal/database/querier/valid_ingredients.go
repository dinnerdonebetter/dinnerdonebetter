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
	_ types.ValidIngredientDataManager = (*SQLQuerier)(nil)
)

// scanValidIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient struct.
func (q *SQLQuerier) scanValidIngredient(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidIngredient{}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.Name,
		&x.Variant,
		&x.Description,
		&x.Warning,
		&x.ContainsEgg,
		&x.ContainsDairy,
		&x.ContainsPeanut,
		&x.ContainsTreeNut,
		&x.ContainsSoy,
		&x.ContainsWheat,
		&x.ContainsShellfish,
		&x.ContainsSesame,
		&x.ContainsFish,
		&x.ContainsGluten,
		&x.AnimalFlesh,
		&x.AnimalDerived,
		&x.Volumetric,
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

// scanValidIngredients takes some database rows and turns them into a slice of valid ingredients.
func (q *SQLQuerier) scanValidIngredients(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredients []*types.ValidIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidIngredient(ctx, rows, includeCounts)
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

		validIngredients = append(validIngredients, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validIngredients, filteredCount, totalCount, nil
}

// ValidIngredientExists fetches whether a valid ingredient exists from the database.
func (q *SQLQuerier) ValidIngredientExists(ctx context.Context, validIngredientID uint64) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validIngredientID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	query, args := q.sqlQueryBuilder.BuildValidIngredientExistsQuery(ctx, validIngredientID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid ingredient existence check")
	}

	return result, nil
}

// GetValidIngredient fetches a valid ingredient from the database.
func (q *SQLQuerier) GetValidIngredient(ctx context.Context, validIngredientID uint64) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validIngredientID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	query, args := q.sqlQueryBuilder.BuildGetValidIngredientQuery(ctx, validIngredientID)
	row := q.getOneRow(ctx, q.db, "validIngredient", query, args...)

	validIngredient, _, _, err := q.scanValidIngredient(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient")
	}

	return validIngredient, nil
}

// GetAllValidIngredientsCount fetches the count of valid ingredients from the database that meet a particular filter.
func (q *SQLQuerier) GetAllValidIngredientsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllValidIngredientsCountQuery(ctx), "fetching count of valid ingredients")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid ingredients")
	}

	return count, nil
}

// GetAllValidIngredients fetches a list of all valid ingredients in the database.
func (q *SQLQuerier) GetAllValidIngredients(ctx context.Context, results chan []*types.ValidIngredient, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllValidIngredientsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of valid ingredients")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfValidIngredientsQuery(ctx, begin, end)
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

			validIngredients, _, _, scanErr := q.scanValidIngredients(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- validIngredients
		}(beginID, endID)
	}

	return nil
}

// GetValidIngredients fetches a list of valid ingredients from the database that meet a particular filter.
func (q *SQLQuerier) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (x *types.ValidIngredientList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.ValidIngredientList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetValidIngredientsQuery(ctx, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "validIngredients", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	if x.ValidIngredients, x.FilteredCount, x.TotalCount, err = q.scanValidIngredients(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredients")
	}

	return x, nil
}

// GetValidIngredientsWithIDs fetches valid ingredients from the database within a given set of IDs.
func (q *SQLQuerier) GetValidIngredientsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]*types.ValidIngredient, error) {
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

	query, args := q.sqlQueryBuilder.BuildGetValidIngredientsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid ingredients with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid ingredients from database")
	}

	validIngredients, _, _, err := q.scanValidIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredients")
	}

	return validIngredients, nil
}

// CreateValidIngredient creates a valid ingredient in the database.
func (q *SQLQuerier) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientCreationInput, createdByUser uint64) (*types.ValidIngredient, error) {
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

	query, args := q.sqlQueryBuilder.BuildCreateValidIngredientQuery(ctx, input)

	// create the valid ingredient.
	id, err := q.performWriteQuery(ctx, tx, false, "valid ingredient creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating valid ingredient")
	}

	x := &types.ValidIngredient{
		ID:                id,
		Name:              input.Name,
		Variant:           input.Variant,
		Description:       input.Description,
		Warning:           input.Warning,
		ContainsEgg:       input.ContainsEgg,
		ContainsDairy:     input.ContainsDairy,
		ContainsPeanut:    input.ContainsPeanut,
		ContainsTreeNut:   input.ContainsTreeNut,
		ContainsSoy:       input.ContainsSoy,
		ContainsWheat:     input.ContainsWheat,
		ContainsShellfish: input.ContainsShellfish,
		ContainsSesame:    input.ContainsSesame,
		ContainsFish:      input.ContainsFish,
		ContainsGluten:    input.ContainsGluten,
		AnimalFlesh:       input.AnimalFlesh,
		AnimalDerived:     input.AnimalDerived,
		Volumetric:        input.Volumetric,
		IconPath:          input.IconPath,
		CreatedOn:         q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidIngredientCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing valid ingredient creation audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachValidIngredientIDToSpan(span, x.ID)
	logger.Info("valid ingredient created")

	return x, nil
}

// UpdateValidIngredient updates a particular valid ingredient. Note that UpdateValidIngredient expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateValidIngredient(ctx context.Context, updated *types.ValidIngredient, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientIDKey, updated.ID)
	tracing.AttachValidIngredientIDToSpan(span, updated.ID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateValidIngredientQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid ingredient update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid ingredient")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidIngredientUpdateEventEntry(changedByUser, updated.ID, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid ingredient update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid ingredient updated")

	return nil
}

// ArchiveValidIngredient archives a valid ingredient from the database by its ID.
func (q *SQLQuerier) ArchiveValidIngredient(ctx context.Context, validIngredientID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validIngredientID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveValidIngredientQuery(ctx, validIngredientID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "valid ingredient archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating valid ingredient")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildValidIngredientArchiveEventEntry(archivedBy, validIngredientID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing valid ingredient archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("valid ingredient archived")

	return nil
}

// GetAuditLogEntriesForValidIngredient fetches a list of audit log entries from the database that relate to a given valid ingredient.
func (q *SQLQuerier) GetAuditLogEntriesForValidIngredient(ctx context.Context, validIngredientID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if validIngredientID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForValidIngredientQuery(ctx, validIngredientID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for valid ingredient", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
