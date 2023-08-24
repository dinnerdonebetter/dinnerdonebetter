package postgres

import (
	"context"
	_ "embed"
	"errors"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.MealPlanDataManager = (*Querier)(nil)

	// ErrAlreadyFinalized is returned when a meal plan is already finalized.
	ErrAlreadyFinalized = errors.New("meal plan already finalized")

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
		"meal_plans.created_by_user",
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
		&x.CreatedByUser,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "scanning meal plan")
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

	result, err := q.generatedQuerier.CheckMealPlanExistence(ctx, q.db, &generated.CheckMealPlanExistenceParams{
		MealPlanID:  mealPlanID,
		HouseholdID: householdID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan existence check")
	}

	return result, nil
}

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

	var mealPlan *types.MealPlan
	if restrictToPastVotingDeadline {
		result, err := q.generatedQuerier.GetMealPlanPastVotingDeadline(ctx, q.db, &generated.GetMealPlanPastVotingDeadlineParams{
			MealPlanID:  mealPlanID,
			HouseholdID: householdID,
		})
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan retrieval")
		}

		mealPlan = &types.MealPlan{
			CreatedAt:              result.CreatedAt,
			VotingDeadline:         result.VotingDeadline,
			ArchivedAt:             timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:          timePointerFromNullTime(result.LastUpdatedAt),
			ID:                     result.ID,
			Status:                 string(result.Status),
			Notes:                  result.Notes,
			ElectionMethod:         string(result.ElectionMethod),
			BelongsToHousehold:     result.BelongsToHousehold,
			CreatedByUser:          result.CreatedByUser,
			GroceryListInitialized: result.GroceryListInitialized,
			TasksCreated:           result.TasksCreated,
		}
	} else {
		result, err := q.generatedQuerier.GetMealPlan(ctx, q.db, &generated.GetMealPlanParams{
			ID:                 mealPlanID,
			BelongsToHousehold: householdID,
		})
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan retrieval")
		}

		mealPlan = &types.MealPlan{
			CreatedAt:              result.CreatedAt,
			VotingDeadline:         result.VotingDeadline,
			ArchivedAt:             timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:          timePointerFromNullTime(result.LastUpdatedAt),
			ID:                     result.ID,
			Status:                 string(result.Status),
			Notes:                  result.Notes,
			ElectionMethod:         string(result.ElectionMethod),
			BelongsToHousehold:     result.BelongsToHousehold,
			CreatedByUser:          result.CreatedByUser,
			GroceryListInitialized: result.GroceryListInitialized,
			TasksCreated:           result.TasksCreated,
		}
	}

	events, err := q.getMealPlanEventsForMealPlan(ctx, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "populating meal plan events")
	}

	if events != nil {
		mealPlan.Events = events
	}

	return mealPlan, nil
}

// GetMealPlan fetches a meal plan from the database.
func (q *Querier) GetMealPlan(ctx context.Context, mealPlanID, householdID string) (*types.MealPlan, error) {
	return q.getMealPlan(ctx, mealPlanID, householdID, false)
}

// GetMealPlans fetches a list of meal plans from the database that meet a particular filter.
func (q *Querier) GetMealPlans(ctx context.Context, householdID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.MealPlan], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.MealPlan]{
		Pagination: filter.ToPagination(),
	}

	query, args := q.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, householdID, false, filter)

	rows, err := q.getRows(ctx, q.db, "meal plans", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing meal plans list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanMealPlans(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meal plans")
	}

	fullMealPlans := []*types.MealPlan{}
	for _, mp := range x.Data {
		fmp, mealPlanFetchErr := q.getMealPlan(ctx, mp.ID, householdID, false)
		if mealPlanFetchErr != nil {
			return nil, observability.PrepareError(mealPlanFetchErr, span, "scanning meal plans")
		}

		fullMealPlans = append(fullMealPlans, fmp)
	}
	x.Data = fullMealPlans

	return x, nil
}

// CreateMealPlan creates a meal plan in the database.
func (q *Querier) CreateMealPlan(ctx context.Context, input *types.MealPlanDatabaseCreationInput) (*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanIDKey, input.ID)

	status := types.MealPlanStatusFinalized
	for _, event := range input.Events {
		if len(event.Options) > 1 {
			status = types.MealPlanStatusAwaitingVotes
		}
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	// create the meal plan.
	if err = q.generatedQuerier.CreateMealPlan(ctx, q.db, &generated.CreateMealPlanParams{
		ID:                 input.ID,
		Notes:              input.Notes,
		Status:             generated.MealPlanStatus(status),
		VotingDeadline:     input.VotingDeadline,
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedByUser:      input.CreatedByUser,
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan")
	}

	x := &types.MealPlan{
		ID:                 input.ID,
		Notes:              input.Notes,
		Status:             string(status),
		VotingDeadline:     input.VotingDeadline,
		BelongsToHousehold: input.BelongsToHousehold,
		ElectionMethod:     input.ElectionMethod,
		CreatedAt:          q.currentTime(),
		CreatedByUser:      input.CreatedByUser,
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

	if err := q.generatedQuerier.UpdateMealPlan(ctx, q.db, &generated.UpdateMealPlanParams{
		Notes:              updated.Notes,
		Status:             generated.MealPlanStatus(updated.Status),
		VotingDeadline:     updated.VotingDeadline,
		BelongsToHousehold: updated.BelongsToHousehold,
		ID:                 updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan")
	}

	logger.Info("meal plan updated")

	return nil
}

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

	if err := q.generatedQuerier.ArchiveMealPlan(ctx, q.db, &generated.ArchiveMealPlanParams{
		BelongsToHousehold: householdID,
		ID:                 mealPlanID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan")
	}

	logger.Info("meal plan archived")

	return nil
}

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

	household, err := q.GetHousehold(ctx, householdID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching household")
	}

	// fetch meal plan
	mealPlan, err := q.getMealPlan(ctx, mealPlanID, householdID, false)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching meal plan")
	}

	votingDeadlineHasPassed := mealPlan.VotingDeadline.Before(q.currentTime())
	if strings.EqualFold(mealPlan.Status, string(types.MealPlanStatusFinalized)) {
		return false, ErrAlreadyFinalized
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	allVotesAreSubmitted := true
	for _, event := range mealPlan.Events {
		if len(event.Options) == 0 {
			continue
		}

		// we load this map with false for each member of the household
		// and then iterate through the votes and mark each voter as true
		userHasVoted := map[string]bool{}
		for _, member := range household.Members {
			userHasVoted[member.BelongsToUser.ID] = true
		}

		alreadyChosen := false
		for _, opt := range event.Options {
			if opt.Chosen {
				alreadyChosen = true
				break
			}

			for _, vote := range opt.Votes {
				if _, ok := userHasVoted[vote.ByUser]; ok {
					userHasVoted[vote.ByUser] = false
				}
			}
		}

		// if we've previously marked an event option as chosen, then we don't need to do anything else
		if alreadyChosen {
			continue
		}

		for _, hasVoted := range userHasVoted {
			if hasVoted {
				allVotesAreSubmitted = false
			}
		}

		// if we're missing votes from household members, and the deadline hasn't passed, then we can't finalize the meal plan.
		if !allVotesAreSubmitted && !votingDeadlineHasPassed {
			continue
		}

		// the ballot is ready to be tallied for this event
		winner, tiebroken, chosen := q.decideOptionWinner(ctx, event.Options)
		if chosen {
			logger = logger.WithValue("winner", winner).WithValue("tiebroken", tiebroken)

			if err = q.generatedQuerier.FinalizeMealPlanOption(ctx, q.db, &generated.FinalizeMealPlanOptionParams{
				BelongsToMealPlanEvent: nullStringFromString(event.ID),
				ID:                     winner,
				Tiebroken:              tiebroken,
			}); err != nil {
				return false, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan option")
			}

			logger.Debug("finalized meal plan option")
		}
	}

	if allVotesAreSubmitted {
		if err = q.generatedQuerier.FinalizeMealPlan(ctx, q.db, &generated.FinalizeMealPlanParams{
			Status: generated.MealPlanStatus(types.MealPlanStatusFinalized),
			ID:     mealPlanID,
		}); err != nil {
			return false, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan option")
		}

		finalized = true
	}

	if err = tx.Commit(); err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "committing transaction")
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
		return nil, observability.PrepareError(err, span, "executing unfinalized meal plans with uninitialized grocery lists query")
	}

	mealPlanDetails := map[string]string{}
	for rows.Next() {
		var mealPlanID, householdID string

		targetVars := []any{
			&mealPlanID,
			&householdID,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, observability.PrepareError(err, span, "scanning meal plan ID")
		}

		mealPlanDetails[mealPlanID] = householdID
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "closing rows")
	}

	mealPlans := []*types.MealPlan{}
	for mealPlanID, householdID := range mealPlanDetails {
		mealPlan, getMealPlanErr := q.GetMealPlan(ctx, mealPlanID, householdID)
		if getMealPlanErr != nil {
			return nil, observability.PrepareError(getMealPlanErr, span, "getting meal plan")
		}

		mealPlans = append(mealPlans, mealPlan)
	}

	return mealPlans, nil
}

// FetchMissingVotesForMealPlan determines the missing votes for a given meal plan.
func (q *Querier) FetchMissingVotesForMealPlan(ctx context.Context, mealPlanID, householdID string) ([]*types.MissingVote, error) {
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

	household, err := q.GetHousehold(ctx, householdID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching household to determine missing votes")
	}

	mealPlan, err := q.GetMealPlan(ctx, mealPlanID, householdID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan to determine missing votes")
	}

	var missingVotes []*types.MissingVote
	for _, event := range mealPlan.Events {
		for _, option := range event.Options {
			for _, membership := range household.Members {
				var voteFoundForMemberForOption bool
				for _, vote := range option.Votes {
					if vote.ByUser == membership.BelongsToUser.ID {
						voteFoundForMemberForOption = true
						break
					}
				}

				if !voteFoundForMemberForOption {
					missingVotes = append(missingVotes, &types.MissingVote{
						EventID:  event.ID,
						OptionID: option.ID,
						UserID:   membership.BelongsToUser.ID,
					})
				}
			}
		}
	}

	return missingVotes, nil
}
