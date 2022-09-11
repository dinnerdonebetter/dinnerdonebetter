package postgres

import (
	"context"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
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
		"meal_plan_events.belongs_to_meal_plan",
		"meal_plan_events.created_at",
		"meal_plan_events.last_updated_at",
		"meal_plan_events.archived_at",
	}
)

// scanMealPlanEvent takes a database Scanner (i.e. *sql.Row) and scans the result into a mealPlanEvent struct.
func (q *Querier) scanMealPlanEvent(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealPlanEvent, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.MealPlanEvent{}

	targetVars := []interface{}{
		&x.ID,
		&x.Notes,
		&x.StartsAt,
		&x.EndsAt,
		&x.BelongsToMealPlan,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanMealPlanEvents takes some database rows and turns them into a slice of meal_plan_events.
func (q *Querier) scanMealPlanEvents(ctx context.Context, rows database.ResultIterator, includeCounts bool) (mealPlanEvents []*types.MealPlanEvent, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

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
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return mealPlanEvents, filteredCount, totalCount, nil
}

const mealPlanEventExistenceQuery = "SELECT EXISTS ( SELECT meal_plan_events.id FROM meal_plan_events WHERE meal_plan_events.archived_at IS NULL AND meal_plan_events.id = $1 )"

// MealPlanEventExists fetches whether a mealPlanEvent exists from the database.
func (q *Querier) MealPlanEventExists(ctx context.Context, mealPlanEventID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanEventID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	args := []interface{}{
		mealPlanEventID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanEventExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing mealPlanEvent existence check")
	}

	return result, nil
}

const getMealPlanEventByIDQuery = `SELECT 
	meal_plan_events.id,
	meal_plan_events.notes,
	meal_plan_events.starts_at,
	meal_plan_events.ends_at,
	meal_plan_events.belongs_to_meal_plan,
	meal_plan_events.created_at,
	meal_plan_events.last_updated_at,
	meal_plan_events.archived_at
FROM meal_plan_events
WHERE meal_plan_events.archived_at IS NULL
	AND meal_plan_events.id = $1
`

// GetMealPlanEvent fetches a mealPlanEvent from the database.
func (q *Querier) GetMealPlanEvent(ctx context.Context, mealPlanEventID string) (*types.MealPlanEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	args := []interface{}{
		mealPlanEventID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "meal plan event", getMealPlanEventByIDQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan event retrieval query")
	}

	m, _, _, err := q.scanMealPlanEvent(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan event retrieval query")
	}

	return m, nil
}

const getMealPlanEventForMealPlanQuery = `SELECT 
	meal_plan_events.id,
	meal_plan_events.notes,
	meal_plan_events.starts_at,
	meal_plan_events.ends_at,
	meal_plan_events.belongs_to_meal_plan,
	meal_plan_events.created_at,
	meal_plan_events.last_updated_at,
	meal_plan_events.archived_at
FROM meal_plan_events
WHERE meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $1
`

// getMealPlanEventsForMealPlan fetches a list of mealPlanEvents from the database that meet a particular filter.
func (q *Querier) getMealPlanEventsForMealPlan(ctx context.Context, mealPlanID string) (x []*types.MealPlanEvent, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	args := []interface{}{
		mealPlanID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "meal plan events", getMealPlanEventForMealPlanQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan events list retrieval query")
	}

	x, _, _, err = q.scanMealPlanEvents(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plan events")
	}

	for _, e := range x {
		mealPlanOptions, mealPlanOptionsErr := q.getMealPlanOptionsForMealPlanEvent(ctx, mealPlanID, e.ID)
		if mealPlanOptionsErr != nil {
			return nil, observability.PrepareError(mealPlanOptionsErr, logger, span, "fetching options for meal plan events")
		}

		e.Options = mealPlanOptions
	}

	return x, nil
}

// GetMealPlanEvents fetches a list of mealPlanEvents from the database that meet a particular filter.
func (q *Querier) GetMealPlanEvents(ctx context.Context, filter *types.QueryFilter) (x *types.MealPlanEventList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.MealPlanEventList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "meal_plan_events", nil, nil, nil, "", mealPlanEventsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "meal plan events", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan events list retrieval query")
	}

	if x.MealPlanEvents, x.FilteredCount, x.TotalCount, err = q.scanMealPlanEvents(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plan events")
	}

	return x, nil
}

const mealPlanEventCreationQuery = `INSERT INTO meal_plan_events (id,notes,starts_at,ends_at,belongs_to_meal_plan) VALUES ($1,$2,$3,$4,$5)`

// CreateMealPlanEvent creates a mealPlanEvent in the database.
func (q *Querier) createMealPlanEvent(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanEventIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Notes,
		input.StartsAt,
		input.EndsAt,
		input.BelongsToMealPlan,
	}

	// create the mealPlanEvent.
	if err := q.performWriteQuery(ctx, querier, "meal plan event creation", mealPlanEventCreationQuery, args); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareError(err, logger, span, "performing mealPlanEvent creation query")
	}

	x := &types.MealPlanEvent{
		ID:                input.ID,
		Notes:             input.Notes,
		StartsAt:          input.StartsAt,
		EndsAt:            input.EndsAt,
		BelongsToMealPlan: input.BelongsToMealPlan,
		CreatedAt:         q.currentTime(),
	}

	logger.WithValue("quantity", len(input.Options)).Info("creating options for meal plan event")
	for _, option := range input.Options {
		option.BelongsToMealPlanEvent = x.ID
		opt, createErr := q.createMealPlanOption(ctx, querier, option)
		if createErr != nil {
			q.rollbackTransaction(ctx, querier)
			return nil, observability.PrepareError(createErr, logger, span, "creating meal plan option for meal plan")
		}
		x.Options = append(x.Options, opt)
	}

	tracing.AttachMealPlanEventIDToSpan(span, x.ID)
	logger.Info("meal plan event created")

	return x, nil
}

// CreateMealPlanEvent creates a mealPlanEvent in the database.
func (q *Querier) CreateMealPlanEvent(ctx context.Context, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanEventIDKey, input.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	x, err := q.createMealPlanEvent(ctx, tx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating mealPlanEvent")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	return x, nil
}

const mealPlanEventRecipeCreationQuery = "INSERT INTO mealPlanEvent_recipes (id,mealPlanEvent_id,recipe_id) VALUES ($1,$2,$3)"

// CreateMealPlanEventRecipe creates a mealPlanEvent in the database.
func (q *Querier) CreateMealPlanEventRecipe(ctx context.Context, querier database.SQLQueryExecutor, mealPlanEventID, recipeID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachUserIDToSpan(span, recipeID)

	args := []interface{}{
		ksuid.New().String(),
		mealPlanEventID,
		recipeID,
	}

	// create the mealPlanEvent.
	if err := q.performWriteQuery(ctx, querier, "meal plan event recipe creation", mealPlanEventRecipeCreationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "performing mealPlanEvent creation query")
	}

	return nil
}

const updateMealPlanEventQuery = `
UPDATE meal_plan_events 
SET notes = $1,
    starts_at = $2,
    ends_at = $3,
    belongs_to_meal_plan = $4,
    last_updated_at = NOW()
WHERE archived_at IS NULL 
  AND id = $5`

// UpdateMealPlanEvent updates a particular meal plan event.
func (q *Querier) UpdateMealPlanEvent(ctx context.Context, updated *types.MealPlanEvent) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanEventIDKey, updated.ID)
	tracing.AttachMealPlanEventIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Notes,
		updated.StartsAt,
		updated.EndsAt,
		updated.BelongsToMealPlan,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan event update", updateMealPlanEventQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan event")
	}

	logger.Info("meal plan event updated")

	return nil
}

const archiveMealPlanEventQuery = "UPDATE meal_plan_events SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1"

// ArchiveMealPlanEvent archives a mealPlanEvent from the database by its ID.
func (q *Querier) ArchiveMealPlanEvent(ctx context.Context, mealPlanEventID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)

	args := []interface{}{
		mealPlanEventID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan event archive", archiveMealPlanEventQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating mealPlanEvent")
	}

	logger.Info("meal plan event archived")

	return nil
}
