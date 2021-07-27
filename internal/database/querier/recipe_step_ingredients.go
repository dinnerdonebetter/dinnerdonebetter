package querier

import (
	"context"
	"database/sql"
	"errors"

	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	database "gitlab.com/prixfixe/prixfixe/internal/database"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var (
	_ types.RecipeStepIngredientDataManager = (*SQLQuerier)(nil)
)

// scanRecipeStepIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step ingredient struct.
func (q *SQLQuerier) scanRecipeStepIngredient(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.RecipeStepIngredient{}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.IngredientID,
		&x.Name,
		&x.QuantityType,
		&x.QuantityValue,
		&x.QuantityNotes,
		&x.ProductOfRecipeStep,
		&x.IngredientNotes,
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

// scanRecipeStepIngredients takes some database rows and turns them into a slice of recipe step ingredients.
func (q *SQLQuerier) scanRecipeStepIngredients(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepIngredients []*types.RecipeStepIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepIngredient(ctx, rows, includeCounts)
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

		recipeStepIngredients = append(recipeStepIngredients, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipeStepIngredients, filteredCount, totalCount, nil
}

// RecipeStepIngredientExists fetches whether a recipe step ingredient exists from the database.
func (q *SQLQuerier) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (exists bool, err error) {
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

	if recipeStepIngredientID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	query, args := q.sqlQueryBuilder.BuildRecipeStepIngredientExistsQuery(ctx, recipeID, recipeStepID, recipeStepIngredientID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe step ingredient existence check")
	}

	return result, nil
}

// GetRecipeStepIngredient fetches a recipe step ingredient from the database.
func (q *SQLQuerier) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*types.RecipeStepIngredient, error) {
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

	if recipeStepIngredientID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	query, args := q.sqlQueryBuilder.BuildGetRecipeStepIngredientQuery(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	row := q.getOneRow(ctx, q.db, "recipeStepIngredient", query, args...)

	recipeStepIngredient, _, _, err := q.scanRecipeStepIngredient(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step ingredient")
	}

	return recipeStepIngredient, nil
}

// GetAllRecipeStepIngredientsCount fetches the count of recipe step ingredients from the database that meet a particular filter.
func (q *SQLQuerier) GetAllRecipeStepIngredientsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllRecipeStepIngredientsCountQuery(ctx), "fetching count of recipe step ingredients")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of recipe step ingredients")
	}

	return count, nil
}

// GetAllRecipeStepIngredients fetches a list of all recipe step ingredients in the database.
func (q *SQLQuerier) GetAllRecipeStepIngredients(ctx context.Context, results chan []*types.RecipeStepIngredient, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllRecipeStepIngredientsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of recipe step ingredients")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfRecipeStepIngredientsQuery(ctx, begin, end)
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

			recipeStepIngredients, _, _, scanErr := q.scanRecipeStepIngredients(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- recipeStepIngredients
		}(beginID, endID)
	}

	return nil
}

// GetRecipeStepIngredients fetches a list of recipe step ingredients from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *types.QueryFilter) (x *types.RecipeStepIngredientList, err error) {
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

	x = &types.RecipeStepIngredientList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetRecipeStepIngredientsQuery(ctx, recipeID, recipeStepID, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "recipeStepIngredients", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	if x.RecipeStepIngredients, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepIngredients(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step ingredients")
	}

	return x, nil
}

// GetRecipeStepIngredientsWithIDs fetches recipe step ingredients from the database within a given set of IDs.
func (q *SQLQuerier) GetRecipeStepIngredientsWithIDs(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) ([]*types.RecipeStepIngredient, error) {
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

	query, args := q.sqlQueryBuilder.BuildGetRecipeStepIngredientsWithIDsQuery(ctx, recipeStepID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "recipe step ingredients with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe step ingredients from database")
	}

	recipeStepIngredients, _, _, err := q.scanRecipeStepIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step ingredients")
	}

	return recipeStepIngredients, nil
}

func (q *SQLQuerier) createRecipeStepIngredient(ctx context.Context, tx *sql.Tx, logger logging.Logger, input *types.RecipeStepIngredientCreationInput, createdByUser uint64) (*types.RecipeStepIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query, args := q.sqlQueryBuilder.BuildCreateRecipeStepIngredientQuery(ctx, input)

	// create the recipe step ingredient.
	id, err := q.performWriteQuery(ctx, tx, false, "recipe step ingredient creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating recipe step ingredient")
	}

	x := &types.RecipeStepIngredient{
		ID:                  id,
		IngredientID:        input.IngredientID,
		Name:                input.Name,
		QuantityType:        input.QuantityType,
		QuantityValue:       input.QuantityValue,
		QuantityNotes:       input.QuantityNotes,
		ProductOfRecipeStep: input.ProductOfRecipeStep,
		IngredientNotes:     input.IngredientNotes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		CreatedOn:           q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeStepIngredientCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing recipe step ingredient creation audit log entry")
	}

	return x, nil
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database.
func (q *SQLQuerier) CreateRecipeStepIngredient(ctx context.Context, input *types.RecipeStepIngredientCreationInput, createdByUser uint64) (*types.RecipeStepIngredient, error) {
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

	x, createErr := q.createRecipeStepIngredient(ctx, tx, logger, input, createdByUser)
	if createErr != nil {
		return nil, observability.PrepareError(createErr, logger, span, "creating recipe step ingredient")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachRecipeStepIngredientIDToSpan(span, x.ID)
	logger.Info("recipe step ingredient created")

	return x, nil
}

// UpdateRecipeStepIngredient updates a particular recipe step ingredient. Note that UpdateRecipeStepIngredient expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateRecipeStepIngredient(ctx context.Context, updated *types.RecipeStepIngredient, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepIngredientIDKey, updated.ID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, updated.ID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateRecipeStepIngredientQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "recipe step ingredient update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating recipe step ingredient")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeStepIngredientUpdateEventEntry(changedByUser, updated.ID, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing recipe step ingredient update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe step ingredient updated")

	return nil
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient from the database by its ID.
func (q *SQLQuerier) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeStepIngredientID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveRecipeStepIngredientQuery(ctx, recipeStepID, recipeStepIngredientID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "recipe step ingredient archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating recipe step ingredient")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeStepIngredientArchiveEventEntry(archivedBy, recipeStepIngredientID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing recipe step ingredient archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe step ingredient archived")

	return nil
}

// GetAuditLogEntriesForRecipeStepIngredient fetches a list of audit log entries from the database that relate to a given recipe step ingredient.
func (q *SQLQuerier) GetAuditLogEntriesForRecipeStepIngredient(ctx context.Context, recipeStepIngredientID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeStepIngredientID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForRecipeStepIngredientQuery(ctx, recipeStepIngredientID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for recipe step ingredient", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
