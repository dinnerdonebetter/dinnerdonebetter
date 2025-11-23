package mealplanning

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ types.MealPlanDataManager = (*repository)(nil)

	ErrAlreadyFinalized = errors.New("meal plan already finalized")
)

// MealPlanExists fetches whether a meal plan exists from the database.
func (q *repository) MealPlanExists(ctx context.Context, mealPlanID, accountID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if accountID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	result, err := q.generatedQuerier.CheckMealPlanExistence(ctx, q.db, &generated.CheckMealPlanExistenceParams{
		MealPlanID:       mealPlanID,
		BelongsToAccount: accountID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan existence check")
	}

	return result, nil
}

// GetMealPlan fetches a meal plan from the database.
func (q *repository) getMealPlan(ctx context.Context, mealPlanID, accountID string) (*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	result, err := q.generatedQuerier.GetMealPlan(ctx, q.db, &generated.GetMealPlanParams{
		ID:               mealPlanID,
		BelongsToAccount: accountID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan retrieval")
	}

	mealPlan := &types.MealPlan{
		CreatedAt:              result.CreatedAt,
		VotingDeadline:         result.VotingDeadline,
		ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
		ID:                     result.ID,
		Status:                 string(result.Status),
		Notes:                  result.Notes,
		ElectionMethod:         string(result.ElectionMethod),
		BelongsToAccount:       result.BelongsToAccount,
		CreatedByUser:          result.CreatedByUser,
		GroceryListInitialized: result.GroceryListInitialized,
		TasksCreated:           result.TasksCreated,
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
func (q *repository) GetMealPlan(ctx context.Context, mealPlanID, accountID string) (*types.MealPlan, error) {
	return q.getMealPlan(ctx, mealPlanID, accountID)
}

// GetMealPlansForAccount fetches a list of meal plans from the database that meet a particular filter.
func (q *repository) GetMealPlansForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.MealPlan], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.MealPlan]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetMealPlansForAccount(ctx, q.db, &generated.GetMealPlansForAccountParams{
		BelongsToAccount: accountID,
		CreatedBefore:    database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:     database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:    database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:     database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:           database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived:  database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plans retrieval")
	}

	for _, result := range results {
		// Extract counts from the first result (all rows have the same counts)
		if x.TotalCount == 0 {
			x.TotalCount = uint64(result.TotalCount)
			x.FilteredCount = uint64(result.FilteredCount)
		}

		x.Data = append(x.Data, &types.MealPlan{
			CreatedAt:              result.CreatedAt,
			VotingDeadline:         result.VotingDeadline,
			ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
			ID:                     result.ID,
			Status:                 string(result.Status),
			Notes:                  result.Notes,
			ElectionMethod:         string(result.ElectionMethod),
			BelongsToAccount:       result.BelongsToAccount,
			CreatedByUser:          result.CreatedByUser,
			Events:                 nil,
			GroceryListInitialized: result.GroceryListInitialized,
			TasksCreated:           result.TasksCreated,
		})
	}

	fullMealPlans := []*types.MealPlan{}
	for _, mp := range x.Data {
		fmp, mealPlanFetchErr := q.getMealPlan(ctx, mp.ID, accountID)
		if mealPlanFetchErr != nil {
			return nil, observability.PrepareError(mealPlanFetchErr, span, "scanning meal plans")
		}

		fullMealPlans = append(fullMealPlans, fmp)
	}
	x.Data = fullMealPlans

	return x, nil
}

// CreateMealPlan creates a meal plan in the database.
func (q *repository) CreateMealPlan(ctx context.Context, input *types.MealPlanDatabaseCreationInput) (*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
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
		ID:               input.ID,
		Notes:            input.Notes,
		Status:           generated.MealPlanStatus(status),
		VotingDeadline:   input.VotingDeadline,
		BelongsToAccount: input.BelongsToAccount,
		CreatedByUser:    input.CreatedByUser,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan")
	}

	x := &types.MealPlan{
		ID:               input.ID,
		Notes:            input.Notes,
		Status:           string(status),
		VotingDeadline:   input.VotingDeadline,
		BelongsToAccount: input.BelongsToAccount,
		ElectionMethod:   input.ElectionMethod,
		CreatedAt:        q.CurrentTime(),
		CreatedByUser:    input.CreatedByUser,
	}

	logger.WithValue("quantity", len(input.Events)).Info("creating events for meal plan")
	for _, event := range input.Events {
		event.BelongsToMealPlan = x.ID
		opt, createErr := q.createMealPlanEvent(ctx, tx, event)
		if createErr != nil {
			q.RollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createErr, span, "creating meal plan event for meal plan")
		}
		x.Events = append(x.Events, opt)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	tracing.AttachToSpan(span, keys.MealPlanIDKey, x.ID)
	logger.Info("meal plan created")

	return x, nil
}

// UpdateMealPlan updates a particular meal plan.
func (q *repository) UpdateMealPlan(ctx context.Context, updated *types.MealPlan) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealPlanIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.AccountIDKey, updated.BelongsToAccount)

	if _, err := q.generatedQuerier.UpdateMealPlan(ctx, q.db, &generated.UpdateMealPlanParams{
		Notes:            updated.Notes,
		Status:           generated.MealPlanStatus(updated.Status),
		VotingDeadline:   updated.VotingDeadline,
		BelongsToAccount: updated.BelongsToAccount,
		ID:               updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan")
	}

	logger.Info("meal plan updated")

	return nil
}

// ArchiveMealPlan archives a meal plan from the database by its ID.
func (q *repository) ArchiveMealPlan(ctx context.Context, mealPlanID, accountID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if accountID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	rowsAffected, err := q.generatedQuerier.ArchiveMealPlan(ctx, q.db, &generated.ArchiveMealPlanParams{
		BelongsToAccount: accountID,
		ID:               mealPlanID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// AttemptToFinalizeMealPlan finalizes a meal plan if all of its options have a selection.
func (q *repository) AttemptToFinalizeMealPlan(ctx context.Context, mealPlanID, accountID string) (finalized bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if accountID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	logger.Info("attempting to finalize meal plan")

	account, err := q.identityRepo.GetAccount(ctx, accountID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching account")
	}

	// fetch meal plan
	mealPlan, err := q.getMealPlan(ctx, mealPlanID, accountID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching meal plan")
	}

	votingDeadlineHasPassed := mealPlan.VotingDeadline.Before(q.CurrentTime())
	if strings.EqualFold(mealPlan.Status, string(types.MealPlanStatusFinalized)) {
		return false, ErrAlreadyFinalized
	}

	usersWhoHaveNotVoted := []string{}
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	allVotesAreSubmitted := true
	for _, event := range mealPlan.Events {
		if len(event.Options) == 0 {
			continue
		}

		// we load this map with false for each member of the account
		// and then iterate through the votes and mark each voter as true
		userHasVoted := map[string]bool{}
		for _, member := range account.Members {
			userHasVoted[member.BelongsToUser.ID] = false
		}

		alreadyChosen := false
		for _, opt := range event.Options {
			if opt.Chosen {
				alreadyChosen = true
				break
			}

			for _, vote := range opt.Votes {
				userHasVoted[vote.ByUser] = true
			}
		}

		// if we've previously marked an event option as chosen, then we don't need to do anything else
		if alreadyChosen {
			continue
		}

		for userID, hasVoted := range userHasVoted {
			if !hasVoted {
				allVotesAreSubmitted = false
				usersWhoHaveNotVoted = append(usersWhoHaveNotVoted, userID)
			}
		}

		// if we're missing votes from account members, and the deadline hasn't passed, then we can't finalize the meal plan.
		if !allVotesAreSubmitted && !votingDeadlineHasPassed {
			logger.WithValue("users_without_votes", usersWhoHaveNotVoted).Info("not all votes are submitted, and the voting deadline hasn't passed yet")
			continue
		}

		// the ballot is ready to be tallied for this event
		winner, tiebroken, chosen := q.decideOptionWinner(ctx, event.Options)
		if chosen {
			logger = logger.WithValue("winner", winner).WithValue("tiebroken", tiebroken)

			if err = q.generatedQuerier.FinalizeMealPlanOption(ctx, q.db, &generated.FinalizeMealPlanOptionParams{
				MealPlanEventID: database.NullStringFromString(event.ID),
				ID:              winner,
				Tiebroken:       tiebroken,
			}); err != nil {
				return false, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan option")
			}

			logger.Info("finalized meal plan option")
		} else {
			logger.Info("no winner chosen")
		}
	}

	if allVotesAreSubmitted || (!allVotesAreSubmitted && votingDeadlineHasPassed) {
		logger.Info("finalizing meal plan")

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

	logger.WithValue("finalized", finalized).
		WithValue("usersWhoHaveNotVoted", usersWhoHaveNotVoted).
		WithValue("allVotesAreSubmitted", allVotesAreSubmitted).
		WithValue("votingDeadlineHasPassed", votingDeadlineHasPassed).
		Info("done attempting to finalize meal plan")

	return finalized, nil
}

// GetUnfinalizedMealPlansWithExpiredVotingPeriods gets unfinalized meal plans with expired voting deadlines.
func (q *repository) GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlan, error) {
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
			ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
			ID:                     result.ID,
			Status:                 string(result.Status),
			Notes:                  result.Notes,
			ElectionMethod:         string(result.ElectionMethod),
			BelongsToAccount:       result.BelongsToAccount,
			CreatedByUser:          result.CreatedByUser,
			Events:                 nil,
			GroceryListInitialized: result.GroceryListInitialized,
			TasksCreated:           result.TasksCreated,
		})
	}

	return mealPlans, nil
}

// GetFinalizedMealPlanIDsForTheNextWeek gets finalized meal plans for a given duration.
func (q *repository) GetFinalizedMealPlanIDsForTheNextWeek(ctx context.Context) ([]*types.FinalizedMealPlanDatabaseResult, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetFinalizedMealPlansForPlanning(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing finalized meal plan IDs for the week retrieval query")
	}

	output := []*types.FinalizedMealPlanDatabaseResult{}
	var databaseResult *types.FinalizedMealPlanDatabaseResult
	for _, result := range results {
		r := &types.FinalizedMealPlanDatabaseResult{
			MealPlanID:       result.MealPlanID,
			MealPlanEventID:  result.MealPlanEventID,
			MealPlanOptionID: result.MealPlanOptionID,
			MealID:           result.MealID,
			RecipeIDs:        nil,
		}

		if databaseResult == nil {
			databaseResult = r
		}

		if r.MealID != databaseResult.MealID &&
			r.MealPlanOptionID != databaseResult.MealPlanOptionID &&
			r.MealPlanEventID != databaseResult.MealPlanEventID &&
			r.MealPlanID != databaseResult.MealPlanID {
			output = append(output, databaseResult)
			databaseResult = r
		}

		databaseResult.RecipeIDs = append(databaseResult.RecipeIDs, result.RecipeID)
	}

	if databaseResult != nil {
		output = append(output, databaseResult)
	}

	return output, nil
}

func (q *repository) GetFinalizedMealPlansWithUninitializedGroceryLists(ctx context.Context) ([]*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetFinalizedMealPlansWithoutGroceryListInit(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing finalized meal plans without grocery list initialization query")
	}

	mealPlanDetails := map[string]string{}
	for _, result := range results {
		mealPlanDetails[result.ID] = result.BelongsToAccount
	}

	mealPlans := []*types.MealPlan{}
	for mealPlanID, accountID := range mealPlanDetails {
		mealPlan, getMealPlanErr := q.GetMealPlan(ctx, mealPlanID, accountID)
		if getMealPlanErr != nil {
			return nil, observability.PrepareError(getMealPlanErr, span, "getting meal plan")
		}

		mealPlans = append(mealPlans, mealPlan)
	}

	return mealPlans, nil
}

// FetchMissingVotesForMealPlan determines the missing votes for a given meal plan.
func (q *repository) FetchMissingVotesForMealPlan(ctx context.Context, mealPlanID, accountID string) ([]*types.MissingVote, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	account, err := q.identityRepo.GetAccount(ctx, accountID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching account to determine missing votes")
	}

	mealPlan, err := q.GetMealPlan(ctx, mealPlanID, accountID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan to determine missing votes")
	}

	var missingVotes []*types.MissingVote
	for _, event := range mealPlan.Events {
		for _, option := range event.Options {
			for _, membership := range account.Members {
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
