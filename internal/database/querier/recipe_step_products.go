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
	_ types.RecipeStepProductDataManager = (*SQLQuerier)(nil)
)

// scanRecipeStepProduct takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step product struct.
func (q *SQLQuerier) scanRecipeStepProduct(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepProduct, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.RecipeStepProduct{}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.Name,
		&x.QuantityType,
		&x.QuantityValue,
		&x.QuantityNotes,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipeStep,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepProducts takes some database rows and turns them into a slice of recipe step products.
func (q *SQLQuerier) scanRecipeStepProducts(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepProducts []*types.RecipeStepProduct, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepProduct(ctx, rows, includeCounts)
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

		recipeStepProducts = append(recipeStepProducts, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipeStepProducts, filteredCount, totalCount, nil
}

// RecipeStepProductExists fetches whether a recipe step product exists from the database.
func (q *SQLQuerier) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	query, args := q.sqlQueryBuilder.BuildRecipeStepProductExistsQuery(ctx, recipeID, recipeStepID, recipeStepProductID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe step product existence check")
	}

	return result, nil
}

// GetRecipeStepProduct fetches a recipe step product from the database.
func (q *SQLQuerier) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	query, args := q.sqlQueryBuilder.BuildGetRecipeStepProductQuery(ctx, recipeID, recipeStepID, recipeStepProductID)
	row := q.getOneRow(ctx, q.db, "recipeStepProduct", query, args...)

	recipeStepProduct, _, _, err := q.scanRecipeStepProduct(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step product")
	}

	return recipeStepProduct, nil
}

// GetAllRecipeStepProductsCount fetches the count of recipe step products from the database that meet a particular filter.
func (q *SQLQuerier) GetAllRecipeStepProductsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllRecipeStepProductsCountQuery(ctx), "fetching count of recipe step products")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of recipe step products")
	}

	return count, nil
}

// GetAllRecipeStepProducts fetches a list of all recipe step products in the database.
func (q *SQLQuerier) GetAllRecipeStepProducts(ctx context.Context, results chan []*types.RecipeStepProduct, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllRecipeStepProductsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of recipe step products")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfRecipeStepProductsQuery(ctx, begin, end)
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

			recipeStepProducts, _, _, scanErr := q.scanRecipeStepProducts(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- recipeStepProducts
		}(beginID, endID)
	}

	return nil
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID uint64, filter *types.QueryFilter) (x *types.RecipeStepProductList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	x = &types.RecipeStepProductList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetRecipeStepProductsQuery(ctx, recipeID, recipeStepID, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "recipeStepProducts", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe step products list retrieval query")
	}

	if x.RecipeStepProducts, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepProducts(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step products")
	}

	return x, nil
}

// GetRecipeStepProductsWithIDs fetches recipe step products from the database within a given set of IDs.
func (q *SQLQuerier) GetRecipeStepProductsWithIDs(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) ([]*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.sqlQueryBuilder.BuildGetRecipeStepProductsWithIDsQuery(ctx, recipeStepID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "recipe step products with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe step products from database")
	}

	recipeStepProducts, _, _, err := q.scanRecipeStepProducts(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step products")
	}

	return recipeStepProducts, nil
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *SQLQuerier) CreateRecipeStepProduct(ctx context.Context, input *types.RecipeStepProductCreationInput, createdByUser uint64) (*types.RecipeStepProduct, error) {
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

	query, args := q.sqlQueryBuilder.BuildCreateRecipeStepProductQuery(ctx, input)

	// create the recipe step product.
	id, err := q.performWriteQuery(ctx, tx, false, "recipe step product creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating recipe step product")
	}

	x := &types.RecipeStepProduct{
		ID:                  id,
		Name:                input.Name,
		QuantityType:        input.QuantityType,
		QuantityValue:       input.QuantityValue,
		QuantityNotes:       input.QuantityNotes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		CreatedOn:           q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeStepProductCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing recipe step product creation audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachRecipeStepProductIDToSpan(span, x.ID)
	logger.Info("recipe step product created")

	return x, nil
}

// UpdateRecipeStepProduct updates a particular recipe step product. Note that UpdateRecipeStepProduct expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepProductIDKey, updated.ID)
	tracing.AttachRecipeStepProductIDToSpan(span, updated.ID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateRecipeStepProductQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "recipe step product update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating recipe step product")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeStepProductUpdateEventEntry(changedByUser, updated.ID, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing recipe step product update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe step product updated")

	return nil
}

// ArchiveRecipeStepProduct archives a recipe step product from the database by its ID.
func (q *SQLQuerier) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeStepProductID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveRecipeStepProductQuery(ctx, recipeStepID, recipeStepProductID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "recipe step product archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating recipe step product")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeStepProductArchiveEventEntry(archivedBy, recipeStepProductID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing recipe step product archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe step product archived")

	return nil
}

// GetAuditLogEntriesForRecipeStepProduct fetches a list of audit log entries from the database that relate to a given recipe step product.
func (q *SQLQuerier) GetAuditLogEntriesForRecipeStepProduct(ctx context.Context, recipeStepProductID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeStepProductID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForRecipeStepProductQuery(ctx, recipeStepProductID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for recipe step product", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
