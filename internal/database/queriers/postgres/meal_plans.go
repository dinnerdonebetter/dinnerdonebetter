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
	_ types.MealPlanDataManager = (*SQLQuerier)(nil)

	// mealPlansTableColumns are the columns for the meal_plans table.
	mealPlansTableColumns = []string{
		"meal_plans.id",
		"meal_plans.notes",
		"meal_plans.state",
		"meal_plans.starts_at",
		"meal_plans.ends_at",
		"meal_plans.created_on",
		"meal_plans.last_updated_on",
		"meal_plans.archived_on",
		"meal_plans.belongs_to_household",
	}
)

// scanMealPlan takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan struct.
func (q *SQLQuerier) scanMealPlan(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealPlan, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.MealPlan{}

	targetVars := []interface{}{
		&x.ID,
		&x.Notes,
		&x.State,
		&x.StartsAt,
		&x.EndsAt,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToHousehold,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanMealPlans takes some database rows and turns them into a slice of meal plans.
func (q *SQLQuerier) scanMealPlans(ctx context.Context, rows database.ResultIterator, includeCounts bool) (mealPlans []*types.MealPlan, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanMealPlan(ctx, rows, includeCounts)
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

		mealPlans = append(mealPlans, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return mealPlans, filteredCount, totalCount, nil
}

const mealPlanExistenceQuery = "SELECT EXISTS ( SELECT meal_plans.id FROM meal_plans WHERE meal_plans.archived_on IS NULL AND meal_plans.id = $1 )"

// MealPlanExists fetches whether a meal plan exists from the database.
func (q *SQLQuerier) MealPlanExists(ctx context.Context, mealPlanID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	args := []interface{}{
		mealPlanID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing meal plan existence check")
	}

	return result, nil
}

const getMealPlanQuery = "SELECT meal_plans.id, meal_plans.notes, meal_plans.state, meal_plans.starts_at, meal_plans.ends_at, meal_plans.created_on, meal_plans.last_updated_on, meal_plans.archived_on, meal_plans.belongs_to_household FROM meal_plans WHERE meal_plans.archived_on IS NULL AND meal_plans.id = $1"

// GetMealPlan fetches a meal plan from the database.
func (q *SQLQuerier) GetMealPlan(ctx context.Context, mealPlanID string) (*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	args := []interface{}{
		mealPlanID,
	}

	row := q.getOneRow(ctx, q.db, "mealPlan", getMealPlanQuery, args)

	mealPlan, _, _, err := q.scanMealPlan(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning mealPlan")
	}

	return mealPlan, nil
}

const getTotalMealPlansCountQuery = "SELECT COUNT(meal_plans.id) FROM meal_plans WHERE meal_plans.archived_on IS NULL"

// GetTotalMealPlanCount fetches the count of meal plans from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalMealPlanCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, getTotalMealPlansCountQuery, "fetching count of meal plans")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of meal plans")
	}

	return count, nil
}

// GetMealPlans fetches a list of meal plans from the database that meet a particular filter.
func (q *SQLQuerier) GetMealPlans(ctx context.Context, filter *types.QueryFilter) (x *types.MealPlanList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.MealPlanList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(
		ctx,
		"meal_plans",
		nil,
		nil,
		householdOwnershipColumn,
		mealPlansTableColumns,
		"",
		false,
		filter,
	)

	rows, err := q.performReadQuery(ctx, q.db, "mealPlans", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plans list retrieval query")
	}

	if x.MealPlans, x.FilteredCount, x.TotalCount, err = q.scanMealPlans(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plans")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetMealPlansWithIDsQuery(ctx context.Context, householdID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"meal_plans.id":          ids,
		"meal_plans.archived_on": nil,
	}

	if householdID != "" {
		withIDsWhere["meal_plans.belongs_to_household"] = householdID
	}

	subqueryBuilder := q.sqlBuilder.Select(mealPlansTableColumns...).
		From("meal_plans").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(mealPlansTableColumns...).
		FromSelect(subqueryBuilder, "meal_plans").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetMealPlansWithIDs fetches meal plans from the database within a given set of IDs.
func (q *SQLQuerier) GetMealPlansWithIDs(ctx context.Context, householdID string, limit uint8, ids []string) ([]*types.MealPlan, error) {
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

	query, args := q.buildGetMealPlansWithIDsQuery(ctx, householdID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "meal plans with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching meal plans from database")
	}

	mealPlans, _, _, err := q.scanMealPlans(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plans")
	}

	return mealPlans, nil
}

const mealPlanCreationQuery = "INSERT INTO meal_plans (id,notes,state,starts_at,ends_at,belongs_to_household) VALUES ($1,$2,$3,$4,$5,$6)"

// CreateMealPlan creates a meal plan in the database.
func (q *SQLQuerier) CreateMealPlan(ctx context.Context, input *types.MealPlanDatabaseCreationInput) (*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Notes,
		input.State,
		input.StartsAt,
		input.EndsAt,
		input.BelongsToHousehold,
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	// create the meal plan.
	if err := q.performWriteQuery(ctx, tx, "meal plan creation", mealPlanCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating meal plan")
	}

	x := &types.MealPlan{
		ID:                 input.ID,
		Notes:              input.Notes,
		State:              input.State,
		StartsAt:           input.StartsAt,
		EndsAt:             input.EndsAt,
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedOn:          q.currentTime(),
	}

	for _, option := range input.Options {
		option.BelongsToMealPlan = x.ID
		opt, createErr := q.createMealPlanOption(ctx, tx, option)
		if createErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createErr, logger, span, "creating meal plan option for meal plan")
		}
		x.Options = append(x.Options, opt)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachMealPlanIDToSpan(span, x.ID)
	logger.Info("meal plan created")

	return x, nil
}

const updateMealPlanQuery = "UPDATE meal_plans SET notes = $1, state = $2, starts_at = $3, ends_at = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_household = $5 AND id = $6"

// UpdateMealPlan updates a particular meal plan.
func (q *SQLQuerier) UpdateMealPlan(ctx context.Context, updated *types.MealPlan) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanIDKey, updated.ID)
	tracing.AttachMealPlanIDToSpan(span, updated.ID)
	tracing.AttachHouseholdIDToSpan(span, updated.BelongsToHousehold)

	args := []interface{}{
		updated.Notes,
		updated.State,
		updated.StartsAt,
		updated.EndsAt,
		updated.BelongsToHousehold,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan update", updateMealPlanQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan")
	}

	logger.Info("meal plan updated")

	return nil
}

const archiveMealPlanQuery = "UPDATE meal_plans SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_household = $1 AND id = $2"

// ArchiveMealPlan archives a meal plan from the database by its ID.
func (q *SQLQuerier) ArchiveMealPlan(ctx context.Context, mealPlanID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		householdID,
		mealPlanID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan archive", archiveMealPlanQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan")
	}

	logger.Info("meal plan archived")

	return nil
}
