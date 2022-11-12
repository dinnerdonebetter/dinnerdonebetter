package postgres

import (
	"context"
	_ "embed"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

var (
	_ types.MealPlanDataManager = (*Querier)(nil)

	// mealPlansTableColumns are the columns for the meal_plans table.
	mealPlansTableColumns = []string{
		"meal_plans.id",
		"meal_plans.notes",
		"meal_plans.status",
		"meal_plans.voting_deadline",
		"meal_plans.grocery_list_initialized",
		"meal_plans.tasks_created",
		"meal_plans.election_method",
		"meal_plans.created_at",
		"meal_plans.last_updated_at",
		"meal_plans.archived_at",
		"meal_plans.belongs_to_household",
	}
)

// scanMealPlan takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan struct.
func (q *Querier) scanMealPlan(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealPlan, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.MealPlan{}

	targetVars := []any{
		&x.ID,
		&x.Notes,
		&x.Status,
		&x.VotingDeadline,
		&x.GroceryListInitialized,
		&x.TasksCreated,
		&x.ElectionMethod,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
		&x.BelongsToHousehold,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanMealPlans takes some database rows and turns them into a slice of meal plans.
func (q *Querier) scanMealPlans(ctx context.Context, rows database.ResultIterator, includeCounts bool) (mealPlans []*types.MealPlan, filteredCount, totalCount uint64, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

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
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return mealPlans, filteredCount, totalCount, nil
}

//go:embed queries/meal_plans/exists.sql
var mealPlanExistenceQuery string

// MealPlanExists fetches whether a meal plan exists from the database.
func (q *Querier) MealPlanExists(ctx context.Context, mealPlanID, householdID string) (exists bool, err error) {
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

	args := []any{
		mealPlanID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan existence check")
	}

	return result, nil
}

//go:embed queries/meal_plans/get_one.sql
var getMealPlanQuery string

//go:embed queries/meal_plans/get_one_past_voting_deadline.sql
var getMealPlanPastVotingDeadlineQuery string

// GetMealPlan fetches a meal plan from the database.
func (q *Querier) getMealPlan(ctx context.Context, mealPlanID, householdID string, restrictToPastVotingDeadline bool) (*types.MealPlan, error) {
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

	args := []any{
		mealPlanID,
		householdID,
	}

	query := getMealPlanQuery
	if restrictToPastVotingDeadline {
		query = getMealPlanPastVotingDeadlineQuery
	}

	row := q.getOneRow(ctx, q.db, "meal plan", query, args)
	if err := row.Err(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan with options retrieval query")
	}

	mealPlan, _, _, err := q.scanMealPlan(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meal plan")
	}

	events, err := q.getMealPlanEventsForMealPlan(ctx, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "populating meal plan events")
	}
	mealPlan.Events = events

	return mealPlan, nil
}

// GetMealPlan fetches a meal plan from the database.
func (q *Querier) GetMealPlan(ctx context.Context, mealPlanID, householdID string) (*types.MealPlan, error) {
	return q.getMealPlan(ctx, mealPlanID, householdID, false)
}

// GetMealPlans fetches a list of meal plans from the database that meet a particular filter.
func (q *Querier) GetMealPlans(ctx context.Context, householdID string, filter *types.QueryFilter) (x *types.MealPlanList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.MealPlanList{}
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, householdID, false, filter)

	rows, err := q.getRows(ctx, q.db, "mealPlans", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing meal plans list retrieval query")
	}

	if x.MealPlans, x.FilteredCount, x.TotalCount, err = q.scanMealPlans(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, span, "scanning meal plans")
	}

	fullMealPlans := []*types.MealPlan{}
	for _, mp := range x.MealPlans {
		fmp, mealPlanFetchErr := q.getMealPlan(ctx, mp.ID, householdID, false)
		if mealPlanFetchErr != nil {
			return nil, observability.PrepareError(mealPlanFetchErr, span, "scanning meal plans")
		}

		fullMealPlans = append(fullMealPlans, fmp)
	}
	x.MealPlans = fullMealPlans

	return x, nil
}

//go:embed queries/meal_plans/create.sql
var mealPlanCreationQuery string

// CreateMealPlan creates a meal plan in the database.
func (q *Querier) CreateMealPlan(ctx context.Context, input *types.MealPlanDatabaseCreationInput) (*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanIDKey, input.ID)

	args := []any{
		input.ID,
		input.Notes,
		types.AwaitingVotesMealPlanStatus,
		input.VotingDeadline,
		input.BelongsToHousehold,
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	// create the meal plan.
	if err = q.performWriteQuery(ctx, tx, "meal plan creation", mealPlanCreationQuery, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan")
	}

	status := types.FinalizedMealPlanStatus
	for _, event := range input.Events {
		if len(event.Options) > 1 {
			status = types.AwaitingVotesMealPlanStatus
		}
	}

	x := &types.MealPlan{
		ID:                 input.ID,
		Notes:              input.Notes,
		Status:             status,
		VotingDeadline:     input.VotingDeadline,
		BelongsToHousehold: input.BelongsToHousehold,
		ElectionMethod:     input.ElectionMethod,
		CreatedAt:          q.currentTime(),
	}

	logger.WithValue("quantity", len(input.Events)).Info("creating events for meal plan")
	for _, event := range input.Events {
		event.BelongsToMealPlan = x.ID
		opt, createErr := q.createMealPlanEvent(ctx, tx, event)
		if createErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createErr, span, "creating meal plan event for meal plan")
		}
		x.Events = append(x.Events, opt)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	tracing.AttachMealPlanIDToSpan(span, x.ID)
	logger.Info("meal plan created")

	return x, nil
}

//go:embed queries/meal_plans/update.sql
var updateMealPlanQuery string

// UpdateMealPlan updates a particular meal plan.
func (q *Querier) UpdateMealPlan(ctx context.Context, updated *types.MealPlan) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanIDKey, updated.ID)
	tracing.AttachMealPlanIDToSpan(span, updated.ID)
	tracing.AttachHouseholdIDToSpan(span, updated.BelongsToHousehold)

	args := []any{
		updated.Notes,
		updated.Status,
		updated.VotingDeadline,
		updated.BelongsToHousehold,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan update", updateMealPlanQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan")
	}

	logger.Info("meal plan updated")

	return nil
}

//go:embed queries/meal_plans/archive.sql
var archiveMealPlanQuery string

// ArchiveMealPlan archives a meal plan from the database by its ID.
func (q *Querier) ArchiveMealPlan(ctx context.Context, mealPlanID, householdID string) error {
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

	args := []any{
		householdID,
		mealPlanID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan archive", archiveMealPlanQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan")
	}

	logger.Info("meal plan archived")

	return nil
}

//go:embed queries/meal_plans/finalize.sql
var finalizeMealPlanQuery string

// AttemptToFinalizeMealPlan finalizes a meal plan if all of its options have a selection.
func (q *Querier) AttemptToFinalizeMealPlan(ctx context.Context, mealPlanID, householdID string) (finalized bool, err error) {
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

	logger.Info("attempting to finalize meal plan")

	household, err := q.GetHouseholdByID(ctx, householdID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching household")
	}

	// fetch meal plan
	mealPlan, err := q.getMealPlan(ctx, mealPlanID, householdID, false)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching meal plan")
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	allOptionsChosen := true
	for _, event := range mealPlan.Events {
		if len(event.Options) == 0 {
			continue
		}

		availableVotes := map[string]bool{}
		for _, member := range household.Members {
			availableVotes[member.BelongsToUser.ID] = false
		}

		alreadyChosen := false
		for _, opt := range event.Options {
			if opt.Chosen {
				alreadyChosen = true
				break
			}
			for _, vote := range opt.Votes {
				if _, ok := availableVotes[vote.ByUser]; ok {
					availableVotes[vote.ByUser] = true
				}
			}
		}

		if alreadyChosen {
			continue
		}

		for _, vote := range availableVotes {
			if !vote {
				allOptionsChosen = false
				continue
			}
		}

		// if we get here, then the tally is ready to be calculated for this set of options

		winner, tiebroken, chosen := q.decideOptionWinner(ctx, event.Options)
		if chosen {
			args := []any{
				event.ID,
				winner,
				tiebroken,
			}

			logger = logger.WithValue("winner", winner).WithValue("tiebroken", tiebroken)

			if err = q.performWriteQuery(ctx, tx, "meal plan option finalization", finalizeMealPlanOptionQuery, args); err != nil {
				q.rollbackTransaction(ctx, tx)
				return false, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan option")
			}

			logger.Debug("finalized meal plan option")
		}
	}

	if allOptionsChosen {
		args := []any{
			types.FinalizedMealPlanStatus,
			mealPlanID,
		}

		if err = q.performWriteQuery(ctx, tx, "meal plan finalization", finalizeMealPlanQuery, args); err != nil {
			q.rollbackTransaction(ctx, tx)
			return false, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan")
		}

		finalized = true
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return false, observability.PrepareAndLogError(commitErr, logger, span, "committing transaction")
	}

	return finalized, nil
}

//go:embed queries/meal_plans/get_expired_and_unresolved.sql
var getExpiredAndUnresolvedMealPlansQuery string

// GetUnfinalizedMealPlansWithExpiredVotingPeriods gets unfinalized meal plans with expired voting deadlines.
func (q *Querier) GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	rows, err := q.getRows(ctx, q.db, "expired and unresolved meal plan", getExpiredAndUnresolvedMealPlansQuery, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing unfinalized meal plans with expired voting periods retrieval query")
	}

	mealPlans := []*types.MealPlan{}
	for rows.Next() {
		mp, _, _, scanErr := q.scanMealPlan(ctx, rows, false)
		if scanErr != nil {
			return nil, observability.PrepareError(scanErr, span, "scanning meal plan response")
		}
		mealPlans = append(mealPlans, mp)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "closing rows")
	}

	return mealPlans, nil
}

//go:embed queries/meal_plans/get_finalized_for_planning.sql
var getFinalizedMealPlansQuery string

// GetFinalizedMealPlanIDsForTheNextWeek gets finalized meal plans for a given duration.
func (q *Querier) GetFinalizedMealPlanIDsForTheNextWeek(ctx context.Context) ([]*types.FinalizedMealPlanDatabaseResult, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	rows, err := q.getRows(ctx, q.db, "finalized meal plans", getFinalizedMealPlansQuery, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing finalized meal plan IDs for the week retrieval query")
	}

	results := []*types.FinalizedMealPlanDatabaseResult{}

	var result *types.FinalizedMealPlanDatabaseResult
	for rows.Next() {
		r := &types.FinalizedMealPlanDatabaseResult{}
		var recipeID string
		scanErr := rows.Scan(&r.MealPlanID, &r.MealPlanOptionID, &r.MealID, &r.MealPlanEventID, &recipeID)
		if scanErr != nil {
			return nil, observability.PrepareError(scanErr, span, "scanning finalized meal plan IDs for the week")
		}

		if result == nil {
			result = r
		}

		if r.MealID != result.MealID &&
			r.MealPlanOptionID != result.MealPlanOptionID &&
			r.MealPlanEventID != result.MealPlanEventID &&
			r.MealPlanID != result.MealPlanID {
			results = append(results, result)
			result = r
		}

		result.RecipeIDs = append(result.RecipeIDs, recipeID)
	}

	if result != nil {
		results = append(results, result)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "closing rows")
	}

	return results, nil
}

//go:embed queries/meal_plans/get_finalized_without_grocery_list_init.sql
var getFinalizedMealPlansWithoutGroceryListInitializationQuery string

func (q *Querier) GetFinalizedMealPlansWithUninitializedGroceryLists(ctx context.Context) ([]*types.MealPlan, error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	rows, err := q.getRows(ctx, q.db, "expired and unresolved meal plan", getFinalizedMealPlansWithoutGroceryListInitializationQuery, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing unfinalized meal plans with expired voting periods retrieval query")
	}

	mealPlans := []*types.MealPlan{}
	for rows.Next() {
		mp, _, _, scanErr := q.scanMealPlan(ctx, rows, false)
		if scanErr != nil {
			return nil, observability.PrepareError(scanErr, span, "scanning meal plan response")
		}
		mealPlans = append(mealPlans, mp)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "closing rows")
	}

	return mealPlans, nil
}
