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
	_ types.RecipeDataManager = (*SQLQuerier)(nil)
)

// scanRecipe takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe struct.
func (q *SQLQuerier) scanRecipe(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.Recipe, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.Recipe{
		Steps: []*types.RecipeStep{},
	}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.Name,
		&x.Source,
		&x.Description,
		&x.InspiredByRecipeID,
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

// scanRecipes takes some database rows and turns them into a slice of recipes.
func (q *SQLQuerier) scanRecipes(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipes []*types.Recipe, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipe(ctx, rows, includeCounts)
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

		recipes = append(recipes, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipes, filteredCount, totalCount, nil
}

// RecipeExists fetches whether a recipe exists from the database.
func (q *SQLQuerier) RecipeExists(ctx context.Context, recipeID uint64) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	query, args := q.sqlQueryBuilder.BuildRecipeExistsQuery(ctx, recipeID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe existence check")
	}

	return result, nil
}

// GetRecipe fetches a recipe from the database.
func (q *SQLQuerier) GetRecipe(ctx context.Context, recipeID uint64) (*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	query, args := q.sqlQueryBuilder.BuildGetRecipeQuery(ctx, recipeID)
	row := q.getOneRow(ctx, q.db, "recipe", query, args...)

	recipe, _, _, err := q.scanRecipe(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe")
	}

	return recipe, nil
}

// GetAllRecipesCount fetches the count of recipes from the database that meet a particular filter.
func (q *SQLQuerier) GetAllRecipesCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllRecipesCountQuery(ctx), "fetching count of recipes")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of recipes")
	}

	return count, nil
}

// GetAllRecipes fetches a list of all recipes in the database.
func (q *SQLQuerier) GetAllRecipes(ctx context.Context, results chan []*types.Recipe, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllRecipesCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of recipes")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfRecipesQuery(ctx, begin, end)
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

			recipes, _, _, scanErr := q.scanRecipes(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- recipes
		}(beginID, endID)
	}

	return nil
}

// GetRecipes fetches a list of recipes from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipes(ctx context.Context, filter *types.QueryFilter) (x *types.RecipeList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.RecipeList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetRecipesQuery(ctx, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "recipes", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipes list retrieval query")
	}

	if x.Recipes, x.FilteredCount, x.TotalCount, err = q.scanRecipes(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipes")
	}

	return x, nil
}

// GetRecipesWithIDs fetches recipes from the database within a given set of IDs.
func (q *SQLQuerier) GetRecipesWithIDs(ctx context.Context, accountID uint64, limit uint8, ids []uint64) ([]*types.Recipe, error) {
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

	query, args := q.sqlQueryBuilder.BuildGetRecipesWithIDsQuery(ctx, accountID, limit, ids, false)

	rows, err := q.performReadQuery(ctx, q.db, "recipes with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipes from database")
	}

	recipes, _, _, err := q.scanRecipes(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipes")
	}

	return recipes, nil
}

// CreateRecipe creates a recipe in the database.
func (q *SQLQuerier) CreateRecipe(ctx context.Context, input *types.RecipeCreationInput, createdByUser uint64) (*types.Recipe, error) {
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

	query, args := q.sqlQueryBuilder.BuildCreateRecipeQuery(ctx, input)

	// create the recipe.
	id, err := q.performWriteQuery(ctx, tx, false, "recipe creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating recipe")
	}

	x := &types.Recipe{
		ID:                 id,
		Name:               input.Name,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		BelongsToAccount:   input.BelongsToAccount,
		CreatedOn:          q.currentTime(),
	}

	for _, stepInput := range input.Steps {
		stepInput.BelongsToRecipe = x.ID
		s, createErr := q.createRecipeStep(ctx, tx, logger, stepInput, createdByUser)
		if createErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step")
		}
		x.Steps = append(x.Steps, s)
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing recipe creation audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachRecipeIDToSpan(span, x.ID)
	logger.Info("recipe created")

	return x, nil
}

// UpdateRecipe updates a particular recipe. Note that UpdateRecipe expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateRecipe(ctx context.Context, updated *types.Recipe, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.RecipeIDKey, updated.ID)
	tracing.AttachRecipeIDToSpan(span, updated.ID)
	tracing.AttachAccountIDToSpan(span, updated.BelongsToAccount)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateRecipeQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "recipe update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating recipe")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeUpdateEventEntry(changedByUser, updated.ID, updated.BelongsToAccount, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing recipe update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe updated")

	return nil
}

// ArchiveRecipe archives a recipe from the database by its ID.
func (q *SQLQuerier) ArchiveRecipe(ctx context.Context, recipeID, accountID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

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

	query, args := q.sqlQueryBuilder.BuildArchiveRecipeQuery(ctx, recipeID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "recipe archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating recipe")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildRecipeArchiveEventEntry(archivedBy, accountID, recipeID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing recipe archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe archived")

	return nil
}

// GetAuditLogEntriesForRecipe fetches a list of audit log entries from the database that relate to a given recipe.
func (q *SQLQuerier) GetAuditLogEntriesForRecipe(ctx context.Context, recipeID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForRecipeQuery(ctx, recipeID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for recipe", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
