package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	mealPlansOnMealPlanEventsJoinClause       = "meal_plans on meal_plan_events.belongs_to_meal_plan = meal_plans.id"
	mealPlanEventsOnMealPlanOptionsJoinClause = "meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id"
)

var (
	_ types.MealPlanEventDataManager = (*Querier)(nil)

	// mealPlanEventsTableColumns are the columns for the mealPlanEvents table.
	mealPlanEventsTableColumns = []string{
		"meal_plan_events.id",
		"meal_plan_events.notes",
		"meal_plan_events.starts_at",
		"meal_plan_events.ends_at",
		"meal_plan_events.meal_name",
		"meal_plan_events.belongs_to_meal_plan",
		"meal_plan_events.created_at",
		"meal_plan_events.last_updated_at",
		"meal_plan_events.archived_at",
	}
)

// scanMealPlanEvent takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan event struct.
func (q *Querier) scanMealPlanEvent(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealPlanEvent, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.MealPlanEvent{}

	targetVars := []any{
		&x.ID,
		&x.Notes,
		&x.StartsAt,
		&x.EndsAt,
		&x.MealName,
		&x.BelongsToMealPlan,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanMealPlanEvents takes some database rows and turns them into a slice of meal_plan_events.
func (q *Querier) scanMealPlanEvents(ctx context.Context, rows database.ResultIterator, includeCounts bool) (mealPlanEvents []*types.MealPlanEvent, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanMealPlanEvent(ctx, rows, includeCounts)
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

		mealPlanEvents = append(mealPlanEvents, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return mealPlanEvents, filteredCount, totalCount, nil
}

// MealPlanEventExists fetches whether a meal plan event exists from the database.
func (q *Querier) MealPlanEventExists(ctx context.Context, mealPlanID, mealPlanEventID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	result, err := q.generatedQuerier.CheckMealPlanEventExistence(ctx, q.db, &generated.CheckMealPlanEventExistenceParams{
		MealPlanEventID: mealPlanEventID,
		MealPlanID:      mealPlanID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan event existence check")
	}

	return result, nil
}

//go:embed queries/meal_plan_events/get_one.sql
var getMealPlanEventByIDQuery string

// GetMealPlanEvent fetches a meal plan event from the database.
func (q *Querier) GetMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	getMealPlanEventByIDArgs := []any{
		mealPlanEventID,
		mealPlanID,
	}

	row := q.getOneRow(ctx, q.db, "meal plan event", getMealPlanEventByIDQuery, getMealPlanEventByIDArgs)
	m, _, _, err := q.scanMealPlanEvent(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan event retrieval query")
	}

	return m, nil
}

//go:embed queries/meal_plan_events/get_for_meal_plan.sql
var getMealPlanEventsForMealPlanQuery string

// getMealPlanEventsForMealPlan fetches a list of mealPlanEvents from the database that meet a particular filter.
func (q *Querier) getMealPlanEventsForMealPlan(ctx context.Context, mealPlanID string) (x []*types.MealPlanEvent, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	getMealPlanEventsForMealPlanArgs := []any{
		mealPlanID,
	}

	rows, err := q.getRows(ctx, q.db, "meal plan events", getMealPlanEventsForMealPlanQuery, getMealPlanEventsForMealPlanArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan events list retrieval query")
	}

	x, _, _, err = q.scanMealPlanEvents(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meal plan events")
	}

	for _, event := range x {
		mealPlanOptions, mealPlanOptionsErr := q.getMealPlanOptionsForMealPlanEvent(ctx, mealPlanID, event.ID)
		if mealPlanOptionsErr != nil {
			return nil, observability.PrepareAndLogError(mealPlanOptionsErr, logger, span, "fetching options for meal plan events")
		}

		event.Options = mealPlanOptions
	}

	return x, nil
}

// GetMealPlanEvents fetches a list of meal plan events from the database that meet a particular filter.
func (q *Querier) GetMealPlanEvents(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.MealPlanEvent], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.MealPlanEvent]{
		Pagination: filter.ToPagination(),
	}

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	query, args := q.buildListQuery(ctx, "meal_plan_events", nil, nil, nil, "", mealPlanEventsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "meal plan events", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan events list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanMealPlanEvents(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meal plan events")
	}

	logger.WithValue("quantity", len(x.Data)).Info("fetched meal plan events")

	return x, nil
}

// MealPlanEventIsEligibleForVoting returns if a meal plan can be voted on.
func (q *Querier) MealPlanEventIsEligibleForVoting(ctx context.Context, mealPlanID, mealPlanEventID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	result, err := q.generatedQuerier.MealPlanEventIsEligibleForVoting(ctx, q.db, &generated.MealPlanEventIsEligibleForVotingParams{
		MealPlanID:      mealPlanID,
		MealPlanEventID: mealPlanEventID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan event existence check")
	}

	return result, nil
}

// createMealPlanEvent creates a meal plan event in the database.
func (q *Querier) createMealPlanEvent(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealPlanEventIDKey, input.ID)

	// create the meal plan event.
	if err := q.generatedQuerier.CreateMealPlanEvent(ctx, querier, &generated.CreateMealPlanEventParams{
		ID:                input.ID,
		Notes:             input.Notes,
		StartsAt:          input.StartsAt,
		EndsAt:            input.EndsAt,
		MealName:          generated.MealName(input.MealName),
		BelongsToMealPlan: input.BelongsToMealPlan,
	}); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan event creation query")
	}

	x := &types.MealPlanEvent{
		ID:                input.ID,
		Notes:             input.Notes,
		StartsAt:          input.StartsAt,
		EndsAt:            input.EndsAt,
		MealName:          input.MealName,
		BelongsToMealPlan: input.BelongsToMealPlan,
		CreatedAt:         q.currentTime(),
	}

	logger.WithValue("quantity", len(input.Options)).Info("creating options for meal plan event")
	for _, option := range input.Options {
		option.BelongsToMealPlanEvent = x.ID
		opt, createErr := q.createMealPlanOption(ctx, querier, option, len(input.Options) == 1)
		if createErr != nil {
			q.rollbackTransaction(ctx, querier)
			return nil, observability.PrepareError(createErr, span, "creating meal plan option for meal plan event")
		}
		x.Options = append(x.Options, opt)
	}

	tracing.AttachMealPlanEventIDToSpan(span, x.ID)
	logger.Info("meal plan event created")

	return x, nil
}

// CreateMealPlanEvent creates a meal plan event in the database.
func (q *Querier) CreateMealPlanEvent(ctx context.Context, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "beginning transaction")
	}

	x, err := q.createMealPlanEvent(ctx, tx, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating meal plan event")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, span, "committing transaction")
	}

	return x, nil
}

// UpdateMealPlanEvent updates a particular meal plan event.
func (q *Querier) UpdateMealPlanEvent(ctx context.Context, updated *types.MealPlanEvent) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealPlanEventIDKey, updated.ID)
	tracing.AttachMealPlanEventIDToSpan(span, updated.ID)

	if err := q.generatedQuerier.UpdateMealPlanEvent(ctx, q.db, &generated.UpdateMealPlanEventParams{
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
func (q *Querier) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	if err := q.generatedQuerier.ArchiveMealPlanEvent(ctx, q.db, &generated.ArchiveMealPlanEventParams{
		ID:                mealPlanEventID,
		BelongsToMealPlan: mealPlanID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan event")
	}

	logger.Info("meal plan event archived")

	return nil
}
