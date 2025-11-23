package mealplanning

import (
	"context"
	"database/sql"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ types.MealPlanOptionVoteDataManager = (*repository)(nil)
)

// MealPlanOptionVoteExists fetches whether a meal plan option vote exists from the database.
func (q *repository) MealPlanOptionVoteExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	if mealPlanOptionVoteID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	result, err := q.generatedQuerier.CheckMealPlanOptionVoteExistence(ctx, q.db, &generated.CheckMealPlanOptionVoteExistenceParams{
		MealPlanOptionID:     mealPlanOptionID,
		MealPlanOptionVoteID: mealPlanOptionVoteID,
		MealPlanEventID:      database.NullStringFromString(mealPlanEventID),
		MealPlanID:           mealPlanID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan option vote existence check")
	}

	return result, nil
}

// GetMealPlanOptionVote fetches a meal plan option vote from the database.
func (q *repository) GetMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanOptionID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	if mealPlanOptionVoteID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	result, err := q.generatedQuerier.GetMealPlanOptionVote(ctx, q.db, &generated.GetMealPlanOptionVoteParams{
		MealPlanOptionID:     mealPlanOptionID,
		MealPlanOptionVoteID: mealPlanOptionVoteID,
		MealPlanID:           mealPlanID,
		MealPlanEventID:      database.NullStringFromString(mealPlanEventID),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting meal plan option vote")
	}

	mealPlanOptionVote := &types.MealPlanOptionVote{
		CreatedAt:               result.CreatedAt,
		ArchivedAt:              database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:           database.TimePointerFromNullTime(result.LastUpdatedAt),
		ID:                      result.ID,
		Notes:                   result.Notes,
		BelongsToMealPlanOption: result.BelongsToMealPlanOption,
		ByUser:                  result.ByUser,
		Rank:                    uint8(result.Rank),
		Abstain:                 result.Abstain,
	}

	return mealPlanOptionVote, nil
}

// GetMealPlanOptionVotesForMealPlanOption fetches a list of meal plan option votes from the database that meet a particular filter.
func (q *repository) GetMealPlanOptionVotesForMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (x []*types.MealPlanOptionVote, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	results, err := q.generatedQuerier.GetMealPlanOptionVotesForMealPlanOption(ctx, q.db, &generated.GetMealPlanOptionVotesForMealPlanOptionParams{
		MealPlanID:       mealPlanID,
		MealPlanOptionID: mealPlanOptionID,
		MealPlanEventID:  database.NullStringFromString(mealPlanEventID),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan option votes for meal plan option")
	}

	x = make([]*types.MealPlanOptionVote, len(results))
	for i, result := range results {
		x[i] = &types.MealPlanOptionVote{
			CreatedAt:               result.CreatedAt,
			ArchivedAt:              database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:           database.TimePointerFromNullTime(result.LastUpdatedAt),
			ID:                      result.ID,
			Notes:                   result.Notes,
			BelongsToMealPlanOption: result.BelongsToMealPlanOption,
			ByUser:                  result.ByUser,
			Rank:                    uint8(result.Rank),
			Abstain:                 result.Abstain,
		}
	}

	return x, nil
}

// GetMealPlanOptionVotes fetches a list of meal plan option votes from the database that meet a particular filter.
func (q *repository) GetMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.MealPlanOptionVote], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	var (
		data          []*types.MealPlanOptionVote
		filteredCount uint64
		totalCount    uint64
	)

	results, err := q.generatedQuerier.GetMealPlanOptionVotes(ctx, q.db, &generated.GetMealPlanOptionVotesParams{
		MealPlanID:       mealPlanID,
		MealPlanOptionID: mealPlanOptionID,
		MealPlanEventID:  database.NullStringFromString(mealPlanEventID),
		CreatedBefore:    database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:     database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:    database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:     database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:           database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived:  database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan option votes list retrieval query")
	}

	for _, result := range results {
		if totalCount == 0 {
			filteredCount = uint64(result.FilteredCount)
			totalCount = uint64(result.TotalCount)
		}
		data = append(data, &types.MealPlanOptionVote{
			CreatedAt:               result.CreatedAt,
			ArchivedAt:              database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:           database.TimePointerFromNullTime(result.LastUpdatedAt),
			ID:                      result.ID,
			Notes:                   result.Notes,
			BelongsToMealPlanOption: result.BelongsToMealPlanOption,
			ByUser:                  result.ByUser,
			Rank:                    uint8(result.Rank),
			Abstain:                 result.Abstain,
		})
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(mpov *types.MealPlanOptionVote) string { return mpov.ID },
		filter,
	)

	return x, nil
}

// CreateMealPlanOptionVote creates a meal plan option vote in the database.
func (q *repository) CreateMealPlanOptionVote(ctx context.Context, input *types.MealPlanOptionVotesDatabaseCreationInput) ([]*types.MealPlanOptionVote, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	logger := q.logger.WithValue("vote_count", len(input.Votes))

	// begin transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	votes := []*types.MealPlanOptionVote{}
	for _, vote := range input.Votes {
		l := logger.WithValue(keys.MealPlanOptionIDKey, vote.BelongsToMealPlanOption).
			WithValue(keys.MealPlanOptionVoteIDKey, vote.ID)

		// create the meal plan option vote.
		if err = q.generatedQuerier.CreateMealPlanOptionVote(ctx, tx, &generated.CreateMealPlanOptionVoteParams{
			ID:                      vote.ID,
			Notes:                   vote.Notes,
			ByUser:                  vote.ByUser,
			BelongsToMealPlanOption: vote.BelongsToMealPlanOption,
			Rank:                    int32(vote.Rank),
			Abstain:                 vote.Abstain,
		}); err != nil {
			q.RollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(err, l, span, "creating meal plan option vote")
		}

		x := &types.MealPlanOptionVote{
			ID:                      vote.ID,
			Rank:                    vote.Rank,
			Abstain:                 vote.Abstain,
			Notes:                   vote.Notes,
			ByUser:                  vote.ByUser,
			BelongsToMealPlanOption: vote.BelongsToMealPlanOption,
			CreatedAt:               q.CurrentTime(),
		}

		tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, x.ID)
		l.Info("meal plan option vote created")

		votes = append(votes, x)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return votes, nil
}

// UpdateMealPlanOptionVote updates a particular meal plan option vote.
func (q *repository) UpdateMealPlanOptionVote(ctx context.Context, updated *types.MealPlanOptionVote) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealPlanOptionVoteIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateMealPlanOptionVote(ctx, q.db, &generated.UpdateMealPlanOptionVoteParams{
		Notes:                   updated.Notes,
		ByUser:                  updated.ByUser,
		BelongsToMealPlanOption: updated.BelongsToMealPlanOption,
		ID:                      updated.ID,
		Rank:                    int32(updated.Rank),
		Abstain:                 updated.Abstain,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan option vote")
	}

	logger.Info("meal plan option vote updated")

	return nil
}

// ArchiveMealPlanOptionVote archives a meal plan option vote from the database by its ID.
func (q *repository) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	if mealPlanOptionVoteID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	rowsAffected, err := q.generatedQuerier.ArchiveMealPlanOptionVote(ctx, q.db, &generated.ArchiveMealPlanOptionVoteParams{
		BelongsToMealPlanOption: mealPlanOptionID,
		ID:                      mealPlanOptionVoteID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan option vote")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
