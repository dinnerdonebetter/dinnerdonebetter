package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	database "github.com/prixfixeco/api_server/internal/database"
	observability "github.com/prixfixeco/api_server/internal/observability"
	keys "github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	mealPlansOnMealPlanOptionsJoinClause = "meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id"
)

var (
	_ types.MealPlanOptionDataManager = (*SQLQuerier)(nil)

	// mealPlanOptionsTableColumns are the columns for the meal_plan_options table.
	mealPlanOptionsTableColumns = []string{
		"meal_plan_options.id",
		"meal_plan_options.day",
		"meal_plan_options.meal_name",
		"meal_plan_options.chosen",
		"meal_plan_options.tiebroken",
		"meal_plan_options.recipe_id",
		"meal_plan_options.notes",
		"meal_plan_options.created_on",
		"meal_plan_options.last_updated_on",
		"meal_plan_options.archived_on",
		"meal_plan_options.belongs_to_meal_plan",
	}

	getMealPlanOptionsJoins = []string{
		mealPlansOnMealPlanOptionsJoinClause,
	}
)

// scanMealPlanOption takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan option struct.
func (q *SQLQuerier) scanMealPlanOption(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealPlanOption, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.MealPlanOption{
		Votes: []*types.MealPlanOptionVote{},
	}

	targetVars := []interface{}{
		&x.ID,
		&x.Day,
		&x.MealName,
		&x.Chosen,
		&x.TieBroken,
		&x.RecipeID,
		&x.Notes,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToMealPlan,
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

const mealPlanOptionExistenceQuery = "SELECT EXISTS ( SELECT meal_plan_options.id FROM meal_plan_options JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id WHERE meal_plan_options.archived_on IS NULL AND meal_plan_options.belongs_to_meal_plan = $1 AND meal_plan_options.id = $2 AND meal_plans.archived_on IS NULL AND meal_plans.id = $3 )"

// MealPlanOptionExists fetches whether a meal plan option exists from the database.
func (q *SQLQuerier) MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanOptionID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	args := []interface{}{
		mealPlanID,
		mealPlanOptionID,
		mealPlanID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanOptionExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing meal plan option existence check")
	}

	return result, nil
}

const getMealPlanOptionQuery = "SELECT meal_plan_options.id, meal_plan_options.day, meal_plan_options.meal_name, meal_plan_options.chosen, meal_plan_options.tiebroken, meal_plan_options.recipe_id, meal_plan_options.notes, meal_plan_options.created_on, meal_plan_options.last_updated_on, meal_plan_options.archived_on, meal_plan_options.belongs_to_meal_plan FROM meal_plan_options JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id WHERE meal_plan_options.archived_on IS NULL AND meal_plan_options.belongs_to_meal_plan = $1 AND meal_plan_options.id = $2 AND meal_plans.archived_on IS NULL AND meal_plans.id = $3"

// GetMealPlanOption fetches a meal plan option from the database.
func (q *SQLQuerier) GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) (*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	args := []interface{}{
		mealPlanID,
		mealPlanOptionID,
		mealPlanID,
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

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getTotalMealPlanOptionsCountQuery, "fetching count of meal plan options")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of meal plan options")
	}

	return count, nil
}

// GetMealPlanOptions fetches a list of meal plan options from the database that meet a particular filter.
func (q *SQLQuerier) GetMealPlanOptions(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (x *types.MealPlanOptionList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	x = &types.MealPlanOptionList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, nil, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "mealPlanOptions", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan options list retrieval query")
	}

	if x.MealPlanOptions, x.FilteredCount, x.TotalCount, err = q.scanMealPlanOptions(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plan options")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetMealPlanOptionsWithIDsQuery(ctx context.Context, mealPlanID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"meal_plan_options.id":                   ids,
		"meal_plan_options.archived_on":          nil,
		"meal_plan_options.belongs_to_meal_plan": mealPlanID,
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
func (q *SQLQuerier) GetMealPlanOptionsWithIDs(ctx context.Context, mealPlanID string, limit uint8, ids []string) ([]*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

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

	query, args := q.buildGetMealPlanOptionsWithIDsQuery(ctx, mealPlanID, limit, ids)

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

const mealPlanOptionCreationQuery = "INSERT INTO meal_plan_options (id,day,meal_name,recipe_id,notes,belongs_to_meal_plan) VALUES ($1,$2,$3,$4,$5,$6)"

// createMealPlanOption creates a meal plan option in the database.
func (q *SQLQuerier) createMealPlanOption(ctx context.Context, db database.SQLQueryExecutor, input *types.MealPlanOptionDatabaseCreationInput) (*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanOptionIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Day,
		input.MealName,
		input.RecipeID,
		input.Notes,
		input.BelongsToMealPlan,
	}

	// create the meal plan option.
	if err := q.performWriteQuery(ctx, db, "meal plan option creation", mealPlanOptionCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating meal plan option")
	}

	x := &types.MealPlanOption{
		ID:                input.ID,
		Day:               input.Day,
		RecipeID:          input.RecipeID,
		MealName:          input.MealName,
		Notes:             input.Notes,
		BelongsToMealPlan: input.BelongsToMealPlan,
		CreatedOn:         q.currentTime(),
		Votes:             []*types.MealPlanOptionVote{},
	}

	tracing.AttachMealPlanOptionIDToSpan(span, x.ID)
	logger.Info("meal plan option created")

	return x, nil
}

// CreateMealPlanOption creates a meal plan option in the database.
func (q *SQLQuerier) CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionDatabaseCreationInput) (*types.MealPlanOption, error) {
	return q.createMealPlanOption(ctx, q.db, input)
}

const updateMealPlanOptionQuery = "UPDATE meal_plan_options SET day = $1, recipe_id = $2, meal_name = $3, notes = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_meal_plan = $5 AND id = $6"

// UpdateMealPlanOption updates a particular meal plan option.
func (q *SQLQuerier) UpdateMealPlanOption(ctx context.Context, updated *types.MealPlanOption) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanOptionIDKey, updated.ID)
	tracing.AttachMealPlanOptionIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Day,
		updated.RecipeID,
		updated.MealName,
		updated.Notes,
		updated.BelongsToMealPlan,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan option update", updateMealPlanOptionQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option")
	}

	logger.Info("meal plan option updated")

	return nil
}

const archiveMealPlanOptionQuery = "UPDATE meal_plan_options SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_meal_plan = $1 AND id = $2"

// ArchiveMealPlanOption archives a meal plan option from the database by its ID.
func (q *SQLQuerier) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	args := []interface{}{
		mealPlanID,
		mealPlanOptionID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan option archive", archiveMealPlanOptionQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option")
	}

	logger.Info("meal plan option archived")

	return nil
}
