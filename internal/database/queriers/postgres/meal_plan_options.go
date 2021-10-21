package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var (
	_ types.MealPlanOptionDataManager = (*SQLQuerier)(nil)

	// mealPlanOptionsTableColumns are the columns for the meal_plan_options table.
	mealPlanOptionsTableColumns = []string{
		"meal_plan_options.id",
		"meal_plan_options.meal_plan_id",
		"meal_plan_options.day_of_week",
		"meal_plan_options.recipe_id",
		"meal_plan_options.notes",
		"meal_plan_options.created_on",
		"meal_plan_options.last_updated_on",
		"meal_plan_options.archived_on",
		"meal_plan_options.belongs_to_account",
	}
)

// scanMealPlanOption takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan option struct.
func (q *SQLQuerier) scanMealPlanOption(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealPlanOption, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.MealPlanOption{}

	targetVars := []interface{}{
		&x.ID,
		&x.MealPlanID,
		&x.DayOfWeek,
		&x.RecipeID,
		&x.Notes,
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

// scanMealPlanOptions takes some database rows and turns them into a slice of meal plan options.
func (q *SQLQuerier) scanMealPlanOptions(ctx context.Context, rows database.ResultIterator, includeCounts bool) (mealPlanOptions []*types.MealPlanOption, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanMealPlanOption(ctx, rows, includeCounts)
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

		mealPlanOptions = append(mealPlanOptions, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return mealPlanOptions, filteredCount, totalCount, nil
}

const mealPlanOptionExistenceQuery = "SELECT EXISTS ( SELECT meal_plan_options.id FROM meal_plan_options WHERE meal_plan_options.archived_on IS NULL AND meal_plan_options.id = $1 )"

// MealPlanOptionExists fetches whether a meal plan option exists from the database.
func (q *SQLQuerier) MealPlanOptionExists(ctx context.Context, mealPlanOptionID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if mealPlanOptionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	args := []interface{}{
		mealPlanOptionID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanOptionExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing meal plan option existence check")
	}

	return result, nil
}

const getMealPlanOptionQuery = "SELECT meal_plan_options.id, meal_plan_options.meal_plan_id, meal_plan_options.day_of_week, meal_plan_options.recipe_id, meal_plan_options.notes, meal_plan_options.created_on, meal_plan_options.last_updated_on, meal_plan_options.archived_on, meal_plan_options.belongs_to_account FROM meal_plan_options WHERE meal_plan_options.archived_on IS NULL AND meal_plan_options.id = $1"

// GetMealPlanOption fetches a meal plan option from the database.
func (q *SQLQuerier) GetMealPlanOption(ctx context.Context, mealPlanOptionID string) (*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	args := []interface{}{
		mealPlanOptionID,
	}

	row := q.getOneRow(ctx, q.db, "mealPlanOption", getMealPlanOptionQuery, args)

	mealPlanOption, _, _, err := q.scanMealPlanOption(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning mealPlanOption")
	}

	return mealPlanOption, nil
}

const getTotalMealPlanOptionsCountQuery = "SELECT COUNT(meal_plan_options.id) FROM meal_plan_options WHERE meal_plan_options.archived_on IS NULL"

// GetTotalMealPlanOptionCount fetches the count of meal plan options from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalMealPlanOptionCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, getTotalMealPlanOptionsCountQuery, "fetching count of meal plan options")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of meal plan options")
	}

	return count, nil
}

// GetMealPlanOptions fetches a list of meal plan options from the database that meet a particular filter.
func (q *SQLQuerier) GetMealPlanOptions(ctx context.Context, filter *types.QueryFilter) (x *types.MealPlanOptionList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.MealPlanOptionList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(
		ctx,
		"meal_plan_options",
		nil,
		nil,
		accountOwnershipColumn,
		mealPlanOptionsTableColumns,
		"",
		false,
		filter,
	)

	rows, err := q.performReadQuery(ctx, q.db, "mealPlanOptions", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan options list retrieval query")
	}

	if x.MealPlanOptions, x.FilteredCount, x.TotalCount, err = q.scanMealPlanOptions(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plan options")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetMealPlanOptionsWithIDsQuery(ctx context.Context, accountID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"meal_plan_options.id":          ids,
		"meal_plan_options.archived_on": nil,
	}

	if accountID != "" {
		withIDsWhere["meal_plan_options.belongs_to_account"] = accountID
	}

	subqueryBuilder := q.sqlBuilder.Select(mealPlanOptionsTableColumns...).
		From("meal_plan_options").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(mealPlanOptionsTableColumns...).
		FromSelect(subqueryBuilder, "meal_plan_options").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetMealPlanOptionsWithIDs fetches meal plan options from the database within a given set of IDs.
func (q *SQLQuerier) GetMealPlanOptionsWithIDs(ctx context.Context, accountID string, limit uint8, ids []string) ([]*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if ids == nil {
		return nil, ErrNilInputProvided
	}

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.buildGetMealPlanOptionsWithIDsQuery(ctx, accountID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "meal plan options with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching meal plan options from database")
	}

	mealPlanOptions, _, _, err := q.scanMealPlanOptions(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plan options")
	}

	return mealPlanOptions, nil
}

const mealPlanOptionCreationQuery = "INSERT INTO meal_plan_options (id,meal_plan_id,day_of_week,recipe_id,notes,belongs_to_account) VALUES ($1,$2,$3,$4,$5,$6)"

// CreateMealPlanOption creates a meal plan option in the database.
func (q *SQLQuerier) CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionDatabaseCreationInput) (*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanOptionIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.MealPlanID,
		input.DayOfWeek,
		input.RecipeID,
		input.Notes,
		input.BelongsToAccount,
	}

	// create the meal plan option.
	if err := q.performWriteQuery(ctx, q.db, "meal plan option creation", mealPlanOptionCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating meal plan option")
	}

	x := &types.MealPlanOption{
		ID:               input.ID,
		MealPlanID:       input.MealPlanID,
		DayOfWeek:        input.DayOfWeek,
		RecipeID:         input.RecipeID,
		Notes:            input.Notes,
		BelongsToAccount: input.BelongsToAccount,
		CreatedOn:        q.currentTime(),
	}

	tracing.AttachMealPlanOptionIDToSpan(span, x.ID)
	logger.Info("meal plan option created")

	return x, nil
}

const updateMealPlanOptionQuery = "UPDATE meal_plan_options SET meal_plan_id = $1, day_of_week = $2, recipe_id = $3, notes = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_account = $5 AND id = $6"

// UpdateMealPlanOption updates a particular meal plan option.
func (q *SQLQuerier) UpdateMealPlanOption(ctx context.Context, updated *types.MealPlanOption) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanOptionIDKey, updated.ID)
	tracing.AttachMealPlanOptionIDToSpan(span, updated.ID)
	tracing.AttachAccountIDToSpan(span, updated.BelongsToAccount)

	args := []interface{}{
		updated.MealPlanID,
		updated.DayOfWeek,
		updated.RecipeID,
		updated.Notes,
		updated.BelongsToAccount,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan option update", updateMealPlanOptionQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option")
	}

	logger.Info("meal plan option updated")

	return nil
}

const archiveMealPlanOptionQuery = "UPDATE meal_plan_options SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_account = $1 AND id = $2"

// ArchiveMealPlanOption archives a meal plan option from the database by its ID.
func (q *SQLQuerier) ArchiveMealPlanOption(ctx context.Context, mealPlanOptionID, accountID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if mealPlanOptionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	if accountID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachAccountIDToSpan(span, accountID)

	args := []interface{}{
		accountID,
		mealPlanOptionID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan option archive", archiveMealPlanOptionQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option")
	}

	logger.Info("meal plan option archived")

	return nil
}
