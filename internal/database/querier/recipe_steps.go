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
	_ types.RecipeStepDataManager = (*SQLQuerier)(nil)
)

// scanRecipeStep takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step struct.
func (q *SQLQuerier) scanRecipeStep(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.RecipeStep{
		Ingredients: []*types.RecipeStepIngredient{},
	}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.Index,
		&x.PreparationID,
		&x.PrerequisiteStep,
		&x.MinEstimatedTimeInSeconds,
		&x.MaxEstimatedTimeInSeconds,
		&x.TemperatureInCelsius,
		&x.Notes,
		&x.Why,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipe,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeSteps takes some database rows and turns them into a slice of recipe steps.
func (q *SQLQuerier) scanRecipeSteps(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeSteps []*types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	recipeSteps = []*types.RecipeStep{}

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStep(ctx, rows, includeCounts)
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

		recipeSteps = append(recipeSteps, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipeSteps, filteredCount, totalCount, nil
}

// RecipeStepExists fetches whether a recipe step exists from the database.
func (q *SQLQuerier) RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (exists bool, err error) {
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

	query, args := q.sqlQueryBuilder.BuildRecipeStepExistsQuery(ctx, recipeID, recipeStepID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe step existence check")
	}

	return result, nil
}

// GetRecipeStep fetches a recipe step from the database.
func (q *SQLQuerier) GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (*types.RecipeStep, error) {
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

	query, args := q.sqlQueryBuilder.BuildGetRecipeStepQuery(ctx, recipeID, recipeStepID)
	row := q.getOneRow(ctx, q.db, "recipeStep", query, args...)

	recipeStep, _, _, err := q.scanRecipeStep(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step")
	}

	return recipeStep, nil
}

// GetAllRecipeStepsCount fetches the count of recipe steps from the database that meet a particular filter.
func (q *SQLQuerier) GetAllRecipeStepsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllRecipeStepsCountQuery(ctx), "fetching count of recipe steps")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of recipe steps")
	}

	return count, nil
}

// GetAllRecipeSteps fetches a list of all recipe steps in the database.
func (q *SQLQuerier) GetAllRecipeSteps(ctx context.Context, results chan []*types.RecipeStep, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllRecipeStepsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of recipe steps")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfRecipeStepsQuery(ctx, begin, end)
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

			recipeSteps, _, _, scanErr := q.scanRecipeSteps(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- recipeSteps
		}(beginID, endID)
	}

	return nil
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipeSteps(ctx context.Context, recipeID uint64, filter *types.QueryFilter) (x *types.RecipeStepList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	x = &types.RecipeStepList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetRecipeStepsQuery(ctx, recipeID, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "recipeSteps", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe steps list retrieval query")
	}

	if x.RecipeSteps, x.FilteredCount, x.TotalCount, err = q.scanRecipeSteps(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe steps")
	}

	return x, nil
}

// GetRecipeStepsWithIDs fetches recipe steps from the database within a given set of IDs.
func (q *SQLQuerier) GetRecipeStepsWithIDs(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) ([]*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.sqlQueryBuilder.BuildGetRecipeStepsWithIDsQuery(ctx, recipeID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "recipe steps with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe steps from database")
	}

	recipeSteps, _, _, err := q.scanRecipeSteps(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe steps")
	}

	return recipeSteps, nil
}

func (q *SQLQuerier) createRecipeStep(ctx context.Context, tx *sql.Tx, logger logging.Logger, input *types.RecipeStepCreationInput, createdByUser uint64) (*types.RecipeStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query, args := q.sqlQueryBuilder.BuildCreateRecipeStepQuery(ctx, input)

	// create the recipe step.
	id, err := q.performWriteQuery(ctx, tx, false, "recipe step creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating recipe step")
	}

	x := &types.RecipeStep{
		ID:                        id,
		Index:                     input.Index,
		PreparationID:             input.PreparationID,
		PrerequisiteStep:          input.PrerequisiteStep,
		MinEstimatedTimeInSeconds: input.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: input.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      input.TemperatureInCelsius,
		Notes:                     input.Notes,
		Why:                       input.Why,
		BelongsToRecipe:           input.BelongsToRecipe,
		CreatedOn:                 q.currentTime(),
	}

	for _, ingredientInput := range input.Ingredients {
		ingredientInput.BelongsToRecipeStep = x.ID
		ingredient, createErr := q.createRecipeStepIngredient(ctx, tx, logger, ingredientInput, createdByUser)
		if createErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step ingredient")
		}

		x.Ingredients = append(x.Ingredients, ingredient)
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeStepCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing recipe step creation audit log entry")
	}

	return x, nil
}

// CreateRecipeStep creates a recipe step in the database.
func (q *SQLQuerier) CreateRecipeStep(ctx context.Context, input *types.RecipeStepCreationInput, createdByUser uint64) (*types.RecipeStep, error) {
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

	x, createErr := q.createRecipeStep(ctx, tx, logger, input, createdByUser)
	if createErr != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(createErr, logger, span, "creating recipe step")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachRecipeStepIDToSpan(span, x.ID)
	logger.Info("recipe step created")

	return x, nil
}

// UpdateRecipeStep updates a particular recipe step. Note that UpdateRecipeStep expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateRecipeStep(ctx context.Context, updated *types.RecipeStep, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepIDKey, updated.ID)
	tracing.AttachRecipeStepIDToSpan(span, updated.ID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateRecipeStepQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "recipe step update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating recipe step")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeStepUpdateEventEntry(changedByUser, updated.ID, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing recipe step update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe step updated")

	return nil
}

// ArchiveRecipeStep archives a recipe step from the database by its ID.
func (q *SQLQuerier) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeStepID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveRecipeStepQuery(ctx, recipeID, recipeStepID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "recipe step archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating recipe step")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeStepArchiveEventEntry(archivedBy, recipeStepID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing recipe step archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe step archived")

	return nil
}

// GetAuditLogEntriesForRecipeStep fetches a list of audit log entries from the database that relate to a given recipe step.
func (q *SQLQuerier) GetAuditLogEntriesForRecipeStep(ctx context.Context, recipeStepID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeStepID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForRecipeStepQuery(ctx, recipeStepID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for recipe step", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
