package postgres

import (
	"context"
	_ "embed"
	"errors"
	"strings"

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
)

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

	results, err := q.generatedQuerier.GetMealPlans(ctx, q.db, &generated.GetMealPlansParams{
		HouseholdID:   householdID,
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plans retrieval")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.MealPlan{
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
			Events:                 nil,
			GroceryListInitialized: result.GroceryListInitialized,
			TasksCreated:           result.TasksCreated,
		})
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

// GetUnfinalizedMealPlansWithExpiredVotingPeriods gets unfinalized meal plans with expired voting deadlines.
func (q *Querier) GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetExpiredAndUnresolvedMealPlans(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing unfinalized meal plans with expired voting periods retrieval query")
	}

	mealPlans := []*types.MealPlan{}
	for _, result := range results {
		mealPlans = append(mealPlans, &types.MealPlan{
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
			Events:                 nil,
			GroceryListInitialized: result.GroceryListInitialized,
			TasksCreated:           result.TasksCreated,
		})
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

func (q *Querier) GetFinalizedMealPlansWithUninitializedGroceryLists(ctx context.Context) ([]*types.MealPlan, error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetFinalizedMealPlansWithoutGroceryListInit(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing finalized meal plans without grocery list initialization query")
	}

	mealPlanDetails := map[string]string{}
	for _, result := range results {
		mealPlanDetails[result.ID] = result.BelongsToHousehold
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
