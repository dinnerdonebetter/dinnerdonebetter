package mealplanning

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

const (
	resourceTypeMealPlanEvents = "meal_plan_events"
)

var (
	_ types.MealPlanEventDataManager = (*repository)(nil)
)

// MealPlanEventExists fetches whether a meal plan event exists from the database.
func (r *repository) MealPlanEventExists(ctx context.Context, mealPlanID, mealPlanEventID string) (exists bool, err error) {
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

	result, err := r.generatedQuerier.CheckMealPlanEventExistence(ctx, r.db, &generated.CheckMealPlanEventExistenceParams{
		ID:         mealPlanEventID,
		MealPlanID: mealPlanID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan event existence check")
	}

	return result, nil
}

// GetMealPlanEvent fetches a meal plan event from the database.
func (r *repository) GetMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error) {
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

	result, err := r.generatedQuerier.GetMealPlanEvent(ctx, r.db, &generated.GetMealPlanEventParams{
		ID:                mealPlanEventID,
		BelongsToMealPlan: mealPlanID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan event retrieval query")
	}

	mealPlanEvent := &types.MealPlanEvent{
		CreatedAt:         result.CreatedAt,
		StartsAt:          result.StartsAt,
		EndsAt:            result.EndsAt,
		ArchivedAt:        database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:     database.TimePointerFromNullTime(result.LastUpdatedAt),
		MealName:          string(result.MealName),
		Notes:             result.Notes,
		BelongsToMealPlan: result.BelongsToMealPlan,
		ID:                result.ID,
	}

	options, err := r.getMealPlanOptionsForMealPlanEvent(ctx, mealPlanID, mealPlanEventID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan options for meal plan event")
	}
	mealPlanEvent.Options = options

	return mealPlanEvent, nil
}

// getMealPlanEventsForMealPlan fetches a list of mealPlanEvents from the database that meet a particular filter.
func (r *repository) getMealPlanEventsForMealPlan(ctx context.Context, mealPlanID string) (x []*types.MealPlanEvent, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	x = []*types.MealPlanEvent{}

	results, err := r.generatedQuerier.GetAllMealPlanEventsForMealPlan(ctx, r.db, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan events list retrieval query")
	}

	for _, result := range results {
		event := &types.MealPlanEvent{
			CreatedAt:         result.CreatedAt,
			StartsAt:          result.StartsAt,
			EndsAt:            result.EndsAt,
			ArchivedAt:        database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:     database.TimePointerFromNullTime(result.LastUpdatedAt),
			MealName:          string(result.MealName),
			Notes:             result.Notes,
			BelongsToMealPlan: result.BelongsToMealPlan,
			ID:                result.ID,
		}

		mealPlanOptions, mealPlanOptionsErr := r.getMealPlanOptionsForMealPlanEvent(ctx, mealPlanID, event.ID)
		if mealPlanOptionsErr != nil {
			return nil, observability.PrepareAndLogError(mealPlanOptionsErr, logger, span, "fetching options for meal plan events")
		}

		event.Options = mealPlanOptions

		x = append(x, event)
	}

	return x, nil
}

// GetMealPlanEvents fetches a list of meal plan events from the database that meet a particular filter.
func (r *repository) GetMealPlanEvents(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.MealPlanEvent], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.MealPlanEvent]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetMealPlanEvents(ctx, r.db, &generated.GetMealPlanEventsParams{
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
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan events list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.MealPlanEvent{
			CreatedAt:         result.CreatedAt,
			StartsAt:          result.StartsAt,
			EndsAt:            result.EndsAt,
			ArchivedAt:        database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:     database.TimePointerFromNullTime(result.LastUpdatedAt),
			MealName:          string(result.MealName),
			Notes:             result.Notes,
			BelongsToMealPlan: result.BelongsToMealPlan,
			ID:                result.ID,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// MealPlanEventIsEligibleForVoting returns if a meal plan can be voted on.
func (r *repository) MealPlanEventIsEligibleForVoting(ctx context.Context, mealPlanID, mealPlanEventID string) (bool, error) {
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

	result, err := r.generatedQuerier.MealPlanEventIsEligibleForVoting(ctx, r.db, &generated.MealPlanEventIsEligibleForVotingParams{
		MealPlanID:      mealPlanID,
		MealPlanEventID: mealPlanEventID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan event existence check")
	}

	return result, nil
}

// createMealPlanEvent creates a meal plan event in the database.
func (r *repository) createMealPlanEvent(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.MealPlanEventIDKey, input.ID)

	// create the meal plan event.
	if err := r.generatedQuerier.CreateMealPlanEvent(ctx, querier, &generated.CreateMealPlanEventParams{
		ID:                input.ID,
		Notes:             input.Notes,
		StartsAt:          input.StartsAt,
		EndsAt:            input.EndsAt,
		MealName:          generated.MealName(input.MealName),
		BelongsToMealPlan: input.BelongsToMealPlan,
	}); err != nil {
		r.RollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan event creation query")
	}

	x := &types.MealPlanEvent{
		ID:                input.ID,
		Notes:             input.Notes,
		StartsAt:          input.StartsAt,
		EndsAt:            input.EndsAt,
		MealName:          input.MealName,
		BelongsToMealPlan: input.BelongsToMealPlan,
		CreatedAt:         r.CurrentTime(),
	}

	logger.WithValue("quantity", len(input.Options)).Info("creating options for meal plan event")
	for _, option := range input.Options {
		option.BelongsToMealPlanEvent = x.ID
		opt, createErr := r.createMealPlanOption(ctx, querier, option, len(input.Options) == 1)
		if createErr != nil {
			r.RollbackTransaction(ctx, querier)
			return nil, observability.PrepareError(createErr, span, "creating meal plan option for meal plan event")
		}
		x.Options = append(x.Options, opt)
	}

	if _, err := r.auditLogEntryRepo.CreateAuditLogEntry(ctx, querier, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypeMealPlanEvents,
		RelevantID:   input.ID,
		EventType:    audit.AuditLogEventTypeCreated,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, x.ID)
	logger.Info("meal plan event created")

	return x, nil
}

// CreateMealPlanEvent creates a meal plan event in the database.
func (r *repository) CreateMealPlanEvent(ctx context.Context, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "beginning transaction")
	}

	x, err := r.createMealPlanEvent(ctx, tx, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating meal plan event")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, span, "committing transaction")
	}

	return x, nil
}

// UpdateMealPlanEvent updates a particular meal plan event.
func (r *repository) UpdateMealPlanEvent(ctx context.Context, updated *types.MealPlanEvent) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.MealPlanEventIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, updated.ID)

	if _, err := r.generatedQuerier.UpdateMealPlanEvent(ctx, r.db, &generated.UpdateMealPlanEventParams{
		Notes:             updated.Notes,
		StartsAt:          updated.StartsAt,
		EndsAt:            updated.EndsAt,
		MealName:          generated.MealName(updated.MealName),
		BelongsToMealPlan: updated.BelongsToMealPlan,
		ID:                updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan event")
	}

	logger.Info("meal plan event updated")

	return nil
}

// ArchiveMealPlanEvent archives a meal plan event from the database by its ID.
func (r *repository) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
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

	rowsAffected, err := r.generatedQuerier.ArchiveMealPlanEvent(ctx, r.db, &generated.ArchiveMealPlanEventParams{
		ID:                mealPlanEventID,
		BelongsToMealPlan: mealPlanID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan event")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
