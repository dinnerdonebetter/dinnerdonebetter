package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.MealPlanDataManager = (*SQLQuerier)(nil)

	// mealPlansTableColumns are the columns for the meal_plans table.
	mealPlansTableColumns = []string{
		"meal_plans.id",
		"meal_plans.notes",
		"meal_plans.status",
		"meal_plans.voting_deadline",
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
		&x.Status,
		&x.VotingDeadline,
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

// scanMealPlanWithOptionsAndVotesRow takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan struct.
func (q *SQLQuerier) scanMealPlanWithOptionsAndVotesRow(ctx context.Context, scan database.Scanner) (mealPlan *types.MealPlan, mealPlanOption *types.MealPlanOption, mealPlanOptionVote *types.MealPlanOptionVote, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger
	mealPlan = &types.MealPlan{}
	mealPlanOption = &types.MealPlanOption{}
	nmpov := &nullableMealPlanOptionVote{}

	targetVars := []interface{}{
		&mealPlan.ID,
		&mealPlan.Notes,
		&mealPlan.Status,
		&mealPlan.VotingDeadline,
		&mealPlan.StartsAt,
		&mealPlan.EndsAt,
		&mealPlan.CreatedOn,
		&mealPlan.LastUpdatedOn,
		&mealPlan.ArchivedOn,
		&mealPlan.BelongsToHousehold,
		&mealPlanOption.ID,
		&mealPlanOption.Day,
		&mealPlanOption.MealName,
		&mealPlanOption.Chosen,
		&mealPlanOption.TieBroken,
		&mealPlanOption.Meal.ID,
		&mealPlanOption.Notes,
		&mealPlanOption.CreatedOn,
		&mealPlanOption.LastUpdatedOn,
		&mealPlanOption.ArchivedOn,
		&mealPlanOption.BelongsToMealPlan,
		&nmpov.ID,
		&nmpov.Rank,
		&nmpov.Abstain,
		&nmpov.Notes,
		&nmpov.ByUser,
		&nmpov.CreatedOn,
		&nmpov.LastUpdatedOn,
		&nmpov.ArchivedOn,
		&nmpov.BelongsToMealPlanOption,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, nil, nil, observability.PrepareError(err, logger, span, "")
	}

	if nmpov.ID != nil {
		mealPlanOptionVote = &types.MealPlanOptionVote{
			LastUpdatedOn:           nmpov.LastUpdatedOn,
			ArchivedOn:              nmpov.ArchivedOn,
			ID:                      *nmpov.ID,
			Notes:                   *nmpov.Notes,
			BelongsToMealPlanOption: *nmpov.BelongsToMealPlanOption,
			ByUser:                  *nmpov.ByUser,
			CreatedOn:               *nmpov.CreatedOn,
			Rank:                    *nmpov.Rank,
			Abstain:                 *nmpov.Abstain,
		}
	}

	return mealPlan, mealPlanOption, mealPlanOptionVote, nil
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
func (q *SQLQuerier) MealPlanExists(ctx context.Context, mealPlanID, householdID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		mealPlanID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing meal plan existence check")
	}

	return result, nil
}

const getMealPlanQuery = `SELECT
	meal_plans.id,
	meal_plans.notes,
	meal_plans.status,
	meal_plans.voting_deadline,
	meal_plans.starts_at,
	meal_plans.ends_at,
	meal_plans.created_on,
	meal_plans.last_updated_on,
	meal_plans.archived_on,
	meal_plans.belongs_to_household,
    meal_plan_options.id,
    meal_plan_options.day,
    meal_plan_options.meal_name,
    meal_plan_options.chosen,
    meal_plan_options.tiebroken,
    meal_plan_options.meal_id,
    meal_plan_options.notes,
    meal_plan_options.created_on,
    meal_plan_options.last_updated_on,
	meal_plan_options.archived_on,
    meal_plan_options.belongs_to_meal_plan,
	meal_plan_option_votes.id,
	meal_plan_option_votes.rank,
	meal_plan_option_votes.abstain,
	meal_plan_option_votes.notes,
	meal_plan_option_votes.by_user,
	meal_plan_option_votes.created_on,
	meal_plan_option_votes.last_updated_on,
	meal_plan_option_votes.archived_on,
	meal_plan_option_votes.belongs_to_meal_plan_option
FROM meal_plans 
	FULL OUTER JOIN meal_plan_options ON meal_plan_options.belongs_to_meal_plan=meal_plans.id
	FULL OUTER JOIN meal_plan_option_votes ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
WHERE meal_plans.archived_on IS NULL 
AND meal_plans.id = $1
AND meal_plans.belongs_to_household = $2
ORDER BY meal_plan_options.id
`

// GetMealPlan fetches a meal plan from the database.
func (q *SQLQuerier) GetMealPlan(ctx context.Context, mealPlanID, householdID string) (*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		mealPlanID,
		householdID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "meal plan", getMealPlanQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan with options retrieval query")
	}

	var (
		mealPlan           *types.MealPlan
		currentOptionIndex = 0
	)
	for rows.Next() {
		rowMealPlan, rowMealPlanOption, rowMealPlanOptionVote, scanErr := q.scanMealPlanWithOptionsAndVotesRow(ctx, rows)
		if scanErr != nil {
			return nil, scanErr
		}

		if mealPlan == nil {
			mealPlan = rowMealPlan
		}

		if len(mealPlan.Options) == 0 && currentOptionIndex == 0 {
			mealPlan.Options = append(mealPlan.Options, rowMealPlanOption)
		}

		if mealPlan.Options[currentOptionIndex].ID != rowMealPlanOption.ID {
			currentOptionIndex++
			mealPlan.Options = append(mealPlan.Options, rowMealPlanOption)
		}

		if rowMealPlanOptionVote != nil {
			mealPlan.Options[currentOptionIndex].Votes = append(mealPlan.Options[currentOptionIndex].Votes, rowMealPlanOptionVote)
		}
	}

	return mealPlan, nil
}

const getTotalMealPlansCountQuery = "SELECT COUNT(meal_plans.id) FROM meal_plans WHERE meal_plans.archived_on IS NULL"

// GetTotalMealPlanCount fetches the count of meal plans from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalMealPlanCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	logger := q.logger.Clone()

	x = &types.MealPlanList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, "", false, filter)

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

	logger := q.logger.Clone()

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

const mealPlanCreationQuery = "INSERT INTO meal_plans (id,notes,status,voting_deadline,starts_at,ends_at,belongs_to_household) VALUES ($1,$2,$3,$4,$5,$6,$7)"

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
		input.Status,
		input.VotingDeadline,
		input.StartsAt,
		input.EndsAt,
		input.BelongsToHousehold,
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	// create the meal plan.
	if err = q.performWriteQuery(ctx, tx, "meal plan creation", mealPlanCreationQuery, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating meal plan")
	}

	x := &types.MealPlan{
		ID:                 input.ID,
		Notes:              input.Notes,
		Status:             input.Status,
		VotingDeadline:     input.VotingDeadline,
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

const updateMealPlanQuery = "UPDATE meal_plans SET notes = $1, status = $2, voting_deadline = $3, starts_at = $4, ends_at = $5, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_household = $6 AND id = $7"

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
		updated.Status,
		updated.VotingDeadline,
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

	logger := q.logger.Clone()

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

var allDays = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
	time.Sunday,
}

var allMealNames = []types.MealName{
	types.BreakfastMealName,
	types.SecondBreakfastMealName,
	types.BrunchMealName,
	types.LunchMealName,
	types.SupperMealName,
	types.DinnerMealName,
}

func byDayAndMeal(l []*types.MealPlanOption, day time.Weekday, meal types.MealName) []*types.MealPlanOption {
	out := []*types.MealPlanOption{}

	for _, o := range l {
		if o.Day == day && o.MealName == meal {
			out = append(out, o)
		}
	}

	return out
}

const finalizeMealPlanQuery = `
	UPDATE meal_plans SET status = $1 WHERE archived_on IS NULL AND id = $2
`

// AttemptToFinalizeCompleteMealPlan finalizes a meal plan if all of its options have a selection.
func (q *SQLQuerier) AttemptToFinalizeCompleteMealPlan(ctx context.Context, mealPlanID, householdID string) (changed bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// fetch meal plan
	mealPlan, err := q.GetMealPlan(ctx, mealPlanID, householdID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "fetching meal plan")
	}

	for _, day := range allDays {
		for _, mealName := range allMealNames {
			winnerChosen := false
			options := byDayAndMeal(mealPlan.Options, day, mealName)
			if len(options) > 0 {
				for _, opt := range options {
					if opt.Chosen {
						winnerChosen = true
					}
				}

				if !winnerChosen {
					return false, nil
				}
			}
		}
	}

	args := []interface{}{
		types.FinalizedMealPlanStatus,
		mealPlanID,
	}

	if err = q.performWriteQuery(ctx, q.db, "meal plan option finalization", finalizeMealPlanQuery, args); err != nil {
		return false, observability.PrepareError(err, logger, span, "finalizing meal plan option")
	}

	logger.Debug("finalized meal plan")

	return true, nil
}

// FinalizeMealPlanWithExpiredVotingPeriod finalizes a meal plan if all of its options have a selection.
func (q *SQLQuerier) FinalizeMealPlanWithExpiredVotingPeriod(ctx context.Context, mealPlanID, householdID string) (changed bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// fetch meal plan
	mealPlan, err := q.GetMealPlan(ctx, mealPlanID, householdID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "fetching meal plan")
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	for _, day := range allDays {
		for _, mealName := range allMealNames {
			options := byDayAndMeal(mealPlan.Options, day, mealName)
			if len(options) > 0 {
				winner, tiebroken := q.decideOptionWinner(options)

				args := []interface{}{
					mealPlanID,
					winner,
					tiebroken,
				}

				if err = q.performWriteQuery(ctx, tx, "meal plan option finalization", finalizeMealPlanOptionQuery, args); err != nil {
					q.rollbackTransaction(ctx, tx)
					return false, observability.PrepareError(err, logger, span, "finalizing meal plan option")
				}

				logger.Debug("finalized meal plan option")
			}
		}
	}

	args := []interface{}{
		types.FinalizedMealPlanStatus,
		mealPlanID,
	}

	if err = q.performWriteQuery(ctx, tx, "meal plan option finalization", finalizeMealPlanQuery, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return false, observability.PrepareError(err, logger, span, "finalizing meal plan option")
	}

	if err = tx.Commit(); err != nil {
		return false, observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Debug("finalized meal plan")

	return true, nil
}

const getExpiredAndUnresolvedMealPlanIDsQuery = `
SELECT
	meal_plans.id,
	meal_plans.notes,
	meal_plans.status,
	meal_plans.voting_deadline,
	meal_plans.starts_at,
	meal_plans.ends_at,
	meal_plans.created_on,
	meal_plans.last_updated_on,
	meal_plans.archived_on,
	meal_plans.belongs_to_household
FROM meal_plans
WHERE meal_plans.archived_on IS NULL 
	AND meal_plans.status = 'awaiting_votes'
	AND to_timestamp(voting_deadline)::date < now()
GROUP BY meal_plans.id
ORDER BY meal_plans.id
`

// GetUnfinalizedMealPlansWithExpiredVotingPeriods gets unfinalized meal plans with expired voting deadlines.
func (q *SQLQuerier) GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlan, error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	rows, err := q.performReadQuery(ctx, q.db, "meal plan", getExpiredAndUnresolvedMealPlanIDsQuery, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan with options retrieval query")
	}

	mealPlans := []*types.MealPlan{}
	for rows.Next() {
		mp, _, _, scanErr := q.scanMealPlan(ctx, rows, false)
		if scanErr != nil {
			return nil, observability.PrepareError(scanErr, logger, span, "scanning meal plan response")
		}
		mealPlans = append(mealPlans, mp)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, logger, span, "closing rows")
	}

	return mealPlans, nil
}
