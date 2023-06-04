package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.MealRatingDataManager = (*Querier)(nil)

	// mealRatingsTableColumns are the columns for the meal_ratings table.
	mealRatingsTableColumns = []string{
		"meal_ratings.id",
		"meal_ratings.meal_id",
		"meal_ratings.taste",
		"meal_ratings.difficulty",
		"meal_ratings.cleanup",
		"meal_ratings.instructions",
		"meal_ratings.overall",
		"meal_ratings.notes",
		"meal_ratings.by_user",
		"meal_ratings.created_at",
		"meal_ratings.last_updated_at",
		"meal_ratings.archived_at",
	}
)

// scanMealRating takes a database Scanner (i.e. *sql.Row) and scans the result into a meal rating struct.
func (q *Querier) scanMealRating(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealRating, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.MealRating{}

	targetVars := []any{
		&x.ID,
		&x.MealID,
		&x.Taste,
		&x.Difficulty,
		&x.Cleanup,
		&x.Instructions,
		&x.Overall,
		&x.Notes,
		&x.ByUser,
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

// scanMealRatings takes some database rows and turns them into a slice of meal ratings.
func (q *Querier) scanMealRatings(ctx context.Context, rows database.ResultIterator, includeCounts bool) (mealRatings []*types.MealRating, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanMealRating(ctx, rows, includeCounts)
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

		mealRatings = append(mealRatings, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return mealRatings, filteredCount, totalCount, nil
}

//go:embed queries/meal_ratings/exists.sql
var mealRatingExistenceQuery string

// MealRatingExists fetches whether a meal rating exists from the database.
func (q *Querier) MealRatingExists(ctx context.Context, mealRatingID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealRatingID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealRatingIDKey, mealRatingID)
	tracing.AttachMealRatingIDToSpan(span, mealRatingID)

	args := []any{
		mealRatingID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealRatingExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal rating existence check")
	}

	return result, nil
}

//go:embed queries/meal_ratings/get_one.sql
var getMealRatingQuery string

// GetMealRating fetches a meal rating from the database.
func (q *Querier) GetMealRating(ctx context.Context, mealRatingID string) (*types.MealRating, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealRatingID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealRatingIDKey, mealRatingID)
	tracing.AttachMealRatingIDToSpan(span, mealRatingID)

	args := []any{
		mealRatingID,
	}

	row := q.getOneRow(ctx, q.db, "mealRating", getMealRatingQuery, args)

	mealRating, _, _, err := q.scanMealRating(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning mealRating")
	}

	return mealRating, nil
}

//go:embed queries/meal_ratings/get_many.sql
var getMealRatingsQuery string

// GetMealRatings fetches a list of meal ratings from the database that meet a particular filter.
func (q *Querier) GetMealRatings(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.MealRating], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.MealRating]{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	} else {
		filter = types.DefaultQueryFilter()
	}

	args := []any{
		filter.CreatedAfter,
		filter.CreatedBefore,
		filter.UpdatedAfter,
		filter.UpdatedBefore,
		filter.QueryOffset(),
	}

	rows, err := q.getRows(ctx, q.db, "meal ratings", getMealRatingsQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal ratings list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanMealRatings(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meal ratings")
	}

	return x, nil
}

//go:embed queries/meal_ratings/create.sql
var mealRatingCreationQuery string

// CreateMealRating creates a meal rating in the database.
func (q *Querier) CreateMealRating(ctx context.Context, input *types.MealRatingDatabaseCreationInput) (*types.MealRating, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealRatingIDKey, input.ID)

	args := []any{
		input.ID,
		input.MealID,
		input.Taste,
		input.Difficulty,
		input.Cleanup,
		input.Instructions,
		input.Overall,
		input.Notes,
		input.ByUser,
	}

	// create the meal rating.
	if err := q.performWriteQuery(ctx, q.db, "meal rating creation", mealRatingCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal rating creation query")
	}

	x := &types.MealRating{
		ID:           input.ID,
		MealID:       input.MealID,
		Taste:        input.Taste,
		Difficulty:   input.Difficulty,
		Cleanup:      input.Cleanup,
		Instructions: input.Instructions,
		Overall:      input.Overall,
		Notes:        input.Notes,
		ByUser:       input.ByUser,
		CreatedAt:    q.currentTime(),
	}

	tracing.AttachMealRatingIDToSpan(span, x.ID)
	logger.Info("meal rating created")

	return x, nil
}

//go:embed queries/meal_ratings/update.sql
var updateMealRatingQuery string

// UpdateMealRating updates a particular meal rating.
func (q *Querier) UpdateMealRating(ctx context.Context, updated *types.MealRating) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealRatingIDKey, updated.ID)
	tracing.AttachMealRatingIDToSpan(span, updated.ID)

	args := []any{
		updated.Notes,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal rating update", updateMealRatingQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal rating")
	}

	logger.Info("meal rating updated")

	return nil
}

//go:embed queries/meal_ratings/archive.sql
var archiveMealRatingQuery string

// ArchiveMealRating archives a meal rating from the database by its ID.
func (q *Querier) ArchiveMealRating(ctx context.Context, mealRatingID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealRatingID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealRatingIDKey, mealRatingID)
	tracing.AttachMealRatingIDToSpan(span, mealRatingID)

	args := []any{
		mealRatingID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal rating archive", archiveMealRatingQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal rating")
	}

	logger.Info("meal rating archived")

	return nil
}
