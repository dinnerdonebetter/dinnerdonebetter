package mealplanning

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"

	"resenje.org/schulze"
)

var (
	_ mealplanning.MealPlanOptionDataManager = (*repository)(nil)
)

// MealPlanOptionExists fetches whether a meal plan option exists from the database.
func (r *repository) MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (exists bool, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if mealPlanID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if mealPlanOptionID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	result, err := r.generatedQuerier.CheckMealPlanOptionExistence(ctx, r.db, &generated.CheckMealPlanOptionExistenceParams{
		MealPlanEventID:  database.NullStringFromString(mealPlanEventID),
		MealPlanOptionID: mealPlanOptionID,
		MealPlanID:       mealPlanID,
	})
	if err != nil {
		logger.Error("performing meal plan option existence check", err)
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan option existence check")
	}

	return result, nil
}

// GetMealPlanOption fetches a meal plan option from the database.
func (r *repository) GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*mealplanning.MealPlanOption, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if mealPlanOptionID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	result, err := r.generatedQuerier.GetMealPlanOption(ctx, r.db, &generated.GetMealPlanOptionParams{
		MealPlanID:       mealPlanID,
		MealPlanEventID:  database.NullStringFromString(mealPlanEventID),
		MealPlanOptionID: mealPlanOptionID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan option query")
	}

	mealPlanOption := &mealplanning.MealPlanOption{
		CreatedAt:              result.CreatedAt,
		LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
		AssignedCook:           database.StringPointerFromNullString(result.AssignedCook),
		ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
		AssignedDishwasher:     database.StringPointerFromNullString(result.AssignedDishwasher),
		Notes:                  result.Notes,
		BelongsToMealPlanEvent: database.StringFromNullString(result.BelongsToMealPlanEvent),
		ID:                     result.ID,
		Meal: mealplanning.Meal{
			CreatedAt:     result.MealCreatedAt,
			ArchivedAt:    database.TimePointerFromNullTime(result.MealArchivedAt),
			LastUpdatedAt: database.TimePointerFromNullTime(result.MealLastUpdatedAt),
			ID:            result.MealID,
			Description:   result.MealDescription,
			CreatedByUser: result.MealCreatedByUser,
			Name:          result.MealName,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: database.Float32FromString(result.MealMinEstimatedPortions),
				Max: database.Float32PointerFromNullString(result.MealMaxEstimatedPortions),
			},
			EligibleForMealPlans: result.MealEligibleForMealPlans,
		},
		MealScale: database.Float32FromString(result.MealScale),
		Chosen:    result.Chosen,
		TieBroken: result.Tiebroken,
	}

	return mealPlanOption, nil
}

// getMealPlanOptionByID fetches a meal plan option from the database by its ID.
func (r *repository) getMealPlanOptionByID(ctx context.Context, mealPlanOptionID string) (*mealplanning.MealPlanOption, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if mealPlanOptionID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	result, err := r.generatedQuerier.GetMealPlanOptionByID(ctx, r.db, mealPlanOptionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan option query")
	}

	mealPlanOption := &mealplanning.MealPlanOption{
		CreatedAt:              result.CreatedAt,
		LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
		AssignedCook:           database.StringPointerFromNullString(result.AssignedCook),
		ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
		AssignedDishwasher:     database.StringPointerFromNullString(result.AssignedDishwasher),
		Notes:                  result.Notes,
		BelongsToMealPlanEvent: database.StringFromNullString(result.BelongsToMealPlanEvent),
		ID:                     result.ID,
		Votes:                  nil,
		Meal: mealplanning.Meal{
			CreatedAt:     result.MealCreatedAt,
			ArchivedAt:    database.TimePointerFromNullTime(result.MealArchivedAt),
			LastUpdatedAt: database.TimePointerFromNullTime(result.MealLastUpdatedAt),
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: database.Float32FromString(result.MealMinEstimatedPortions),
				Max: database.Float32PointerFromNullString(result.MealMaxEstimatedPortions),
			},
			ID:                   result.MealID,
			Description:          result.MealDescription,
			CreatedByUser:        result.MealCreatedByUser,
			Name:                 result.MealName,
			Components:           []*mealplanning.MealComponent{},
			EligibleForMealPlans: result.MealEligibleForMealPlans,
		},
		MealScale: database.Float32FromString(result.MealScale),
		Chosen:    result.Chosen,
		TieBroken: result.Tiebroken,
	}

	return mealPlanOption, nil
}

// getMealPlanOptionsForMealPlanEvent fetches a list of meal plan options from the database that meet a particular filter.
func (r *repository) getMealPlanOptionsForMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) ([]*mealplanning.MealPlanOption, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanEventID)

	results, err := r.generatedQuerier.GetAllMealPlanOptionsForMealPlanEvent(ctx, r.db, &generated.GetAllMealPlanOptionsForMealPlanEventParams{
		MealPlanID:      mealPlanID,
		MealPlanEventID: database.NullStringFromString(mealPlanEventID),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan option query")
	}

	x := []*mealplanning.MealPlanOption{}
	for _, result := range results {
		x = append(x, &mealplanning.MealPlanOption{
			CreatedAt:              result.CreatedAt,
			LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
			AssignedCook:           database.StringPointerFromNullString(result.AssignedCook),
			ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
			AssignedDishwasher:     database.StringPointerFromNullString(result.AssignedDishwasher),
			Notes:                  result.Notes,
			BelongsToMealPlanEvent: database.StringFromNullString(result.BelongsToMealPlanEvent),
			ID:                     result.ID,
			Votes:                  nil,
			Meal: mealplanning.Meal{
				ID: result.MealID,
			},
			MealScale: database.Float32FromString(result.MealScale),
			Chosen:    result.Chosen,
			TieBroken: result.Tiebroken,
		})
	}

	for i, opt := range x {
		votes, voteFetchErr := r.GetMealPlanOptionVotesForMealPlanOption(ctx, mealPlanID, mealPlanEventID, opt.ID)
		if voteFetchErr != nil {
			return nil, observability.PrepareError(voteFetchErr, span, "fetching meal plan option votes for meal plan option")
		}
		x[i].Votes = votes

		m, mealFetchErr := r.GetMeal(ctx, opt.Meal.ID)
		if mealFetchErr != nil {
			return nil, observability.PrepareAndLogError(mealFetchErr, logger, span, "scanning meal plan options for meal plan event")
		}
		x[i].Meal = *m
	}

	logger.WithValue("quantity", len(x)).Info("fetched meal plan options for meal plan event")

	return x, nil
}

// GetMealPlanOptions fetches a list of meal plan options from the database that meet a particular filter.
func (r *repository) GetMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.MealPlanOption], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[mealplanning.MealPlanOption]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetMealPlanOptions(ctx, r.db, &generated.GetMealPlanOptionsParams{
		MealPlanEventID: database.NullStringFromString(mealPlanEventID),
		MealPlanID:      mealPlanID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan option query")
	}

	for _, result := range results {
		mealPlanOption := &mealplanning.MealPlanOption{
			CreatedAt:              result.CreatedAt,
			LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
			AssignedCook:           database.StringPointerFromNullString(result.AssignedCook),
			ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
			AssignedDishwasher:     database.StringPointerFromNullString(result.AssignedDishwasher),
			Notes:                  result.Notes,
			BelongsToMealPlanEvent: database.StringFromNullString(result.BelongsToMealPlanEvent),
			ID:                     result.ID,
			Votes:                  nil,
			Meal: mealplanning.Meal{
				ID: result.MealID,
			},
			MealScale: database.Float32FromString(result.MealScale),
			Chosen:    result.Chosen,
			TieBroken: result.Tiebroken,
		}
		x.Data = append(x.Data, mealPlanOption)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// createMealPlanOption creates a meal plan option in the database.
func (r *repository) createMealPlanOption(ctx context.Context, db database.SQLQueryExecutor, input *mealplanning.MealPlanOptionDatabaseCreationInput, markAsChosen bool) (*mealplanning.MealPlanOption, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, input.ID)
	logger := r.logger.WithValue(keys.MealPlanOptionIDKey, input.ID)

	// create the meal plan option.
	if err := r.generatedQuerier.CreateMealPlanOption(ctx, db, &generated.CreateMealPlanOptionParams{
		ID:                     input.ID,
		AssignedCook:           database.NullStringFromStringPointer(input.AssignedCook),
		AssignedDishwasher:     database.NullStringFromStringPointer(input.AssignedDishwasher),
		MealID:                 input.MealID,
		Notes:                  input.Notes,
		MealScale:              database.StringFromFloat32(input.MealScale),
		BelongsToMealPlanEvent: database.NullStringFromString(input.BelongsToMealPlanEvent),
		Chosen:                 markAsChosen,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan option")
	}

	x := &mealplanning.MealPlanOption{
		ID:                     input.ID,
		AssignedCook:           input.AssignedCook,
		Meal:                   mealplanning.Meal{ID: input.MealID},
		Notes:                  input.Notes,
		BelongsToMealPlanEvent: input.BelongsToMealPlanEvent,
		CreatedAt:              r.CurrentTime(),
		MealScale:              input.MealScale,
		Votes:                  []*mealplanning.MealPlanOptionVote{},
	}

	logger.Info("meal plan option created")

	return x, nil
}

// CreateMealPlanOption creates a meal plan option in the database.
func (r *repository) CreateMealPlanOption(ctx context.Context, input *mealplanning.MealPlanOptionDatabaseCreationInput) (*mealplanning.MealPlanOption, error) {
	return r.createMealPlanOption(ctx, r.db, input, false)
}

// UpdateMealPlanOption updates a particular meal plan option.
func (r *repository) UpdateMealPlanOption(ctx context.Context, updated *mealplanning.MealPlanOption) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.MealPlanOptionIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, updated.ID)

	if _, err := r.generatedQuerier.UpdateMealPlanOption(ctx, r.db, &generated.UpdateMealPlanOptionParams{
		MealID:             updated.Meal.ID,
		Notes:              updated.Notes,
		MealScale:          database.StringFromFloat32(updated.MealScale),
		MealPlanOptionID:   updated.ID,
		AssignedCook:       database.NullStringFromStringPointer(updated.AssignedCook),
		AssignedDishwasher: database.NullStringFromStringPointer(updated.AssignedDishwasher),
		MealPlanEventID:    database.NullStringFromString(updated.BelongsToMealPlanEvent),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan option")
	}

	logger.Info("meal plan option updated")

	return nil
}

// ArchiveMealPlanOption archives a meal plan option from the database by its ID.
func (r *repository) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if mealPlanID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if mealPlanOptionID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	rowsAffected, err := r.generatedQuerier.ArchiveMealPlanOption(ctx, r.db, &generated.ArchiveMealPlanOptionParams{
		ID:                     mealPlanOptionID,
		BelongsToMealPlanEvent: sql.NullString{String: mealPlanEventID, Valid: true},
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan option")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repository) determineWinner(winners []schulze.Result[string]) string {
	var (
		highestScore int
		scoreWinners []string
	)

	for _, winner := range winners {
		if winner.Wins == highestScore {
			scoreWinners = append(scoreWinners, winner.Choice)
		} else if winner.Wins > highestScore {
			highestScore = winner.Wins
			scoreWinners = []string{winner.Choice}
		}
	}

	/* #nosec: G404 */
	return scoreWinners[rand.N(len(scoreWinners))]
}

func (r *repository) decideOptionWinner(ctx context.Context, options []*mealplanning.MealPlanOption) (_ string, _, _ bool) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	candidateMap := map[string]struct{}{}
	votesByUser := map[string]schulze.Ballot[string]{}

	logger := r.logger.WithValue("options.count", len(options))

	for _, option := range options {
		for _, v := range option.Votes {
			if votesByUser[v.ByUser] == nil {
				votesByUser[v.ByUser] = schulze.Ballot[string]{}
			}

			if !v.Abstain {
				votesByUser[v.ByUser][v.BelongsToMealPlanOption] = int(v.Rank)
			}

			candidateMap[v.BelongsToMealPlanOption] = struct{}{}
		}
	}

	candidates := []string{}
	for c := range candidateMap {
		candidates = append(candidates, c)
	}

	e := schulze.NewVoting(candidates)
	for _, vote := range votesByUser {
		if _, err := e.Vote(vote); err != nil {
			// this actually can never happen because we use uints for ranks, lol
			observability.AcknowledgeError(err, logger, span, "an invalid vote was received")
		}
	}

	winners, _, tie := e.Compute()
	if tie {
		return r.determineWinner(winners), true, true
	}

	if len(winners) > 0 {
		return winners[0].Choice, false, true
	}

	return "", false, false
}

// FinalizeMealPlanOption archives a meal plan option vote from the database by its ID.
func (r *repository) FinalizeMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, accountID string) (changed bool, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if mealPlanID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if mealPlanOptionID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	if accountID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	mealPlan, err := r.GetMealPlan(ctx, mealPlanID, accountID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching meal plan")
	}

	var (
		mealPlanEvent  *mealplanning.MealPlanEvent
		mealPlanOption *mealplanning.MealPlanOption
	)
	for _, event := range mealPlan.Events {
		if event.ID == mealPlanEventID {
			mealPlanEvent = event
			for _, option := range event.Options {
				if option.ID == mealPlanOptionID {
					mealPlanOption = option
					break
				}
			}
		}
	}

	if mealPlanEvent == nil {
		return false, fmt.Errorf("meal plan event %s for meal plan %s not found", mealPlanEventID, mealPlanID)
	}

	// fetch account data
	account, err := r.identityRepo.GetAccount(ctx, mealPlan.BelongsToAccount)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching account")
	}

	// go through all the votes for this meal plan option and determine if they're all there
	for _, member := range account.Members {
		memberVoteFound := false
		for _, vote := range mealPlanOption.Votes {
			if vote.ByUser == member.BelongsToUser.ID {
				memberVoteFound = true
				break
			}
		}

		if !memberVoteFound {
			return false, nil
		}
	}

	winner, tiebroken, chosen := r.decideOptionWinner(ctx, mealPlanEvent.Options)
	if chosen {
		if err = r.generatedQuerier.FinalizeMealPlanOption(ctx, r.db, &generated.FinalizeMealPlanOptionParams{
			MealPlanEventID: database.NullStringFromString(mealPlanEventID),
			ID:              winner,
			Tiebroken:       tiebroken,
		}); err != nil {
			return false, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan option")
		}

		logger.Debug("finalized meal plan option")
	}

	return chosen, nil
}
