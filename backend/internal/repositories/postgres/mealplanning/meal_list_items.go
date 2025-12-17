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
	_ types.MealListItemDataManager = (*repository)(nil)
)

// GetMealListItems fetches meal list items for a given list with filtering.
func (q *repository) GetMealListItems(ctx context.Context, mealListID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.MealListItem], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealListID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealListIDKey, mealListID)
	tracing.AttachToSpan(span, keys.MealListIDKey, mealListID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	var (
		data          []*types.MealListItem
		filteredCount uint64
		totalCount    uint64
	)

	results, err := q.generatedQuerier.GetMealListItems(ctx, q.db, &generated.GetMealListItemsParams{
		MealListID:      mealListID,
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:          database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:     database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal list items list retrieval query")
	}

	for _, result := range results {
		if totalCount == 0 {
			filteredCount = uint64(result.FilteredCount)
			totalCount = uint64(result.TotalCount)
		}

		data = append(data, &types.MealListItem{
			CreatedAt:         result.CreatedAt,
			LastUpdatedAt:     database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:        database.TimePointerFromNullTime(result.ArchivedAt),
			ID:                result.ID,
			MealID:            result.MealID,
			Notes:             result.Notes,
			BelongsToMealList: result.BelongsToMealList,
		})
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(mli *types.MealListItem) string { return mli.ID },
		filter,
	)

	return x, nil
}

// CreateMealListItem creates a meal list item in the database.
func (q *repository) CreateMealListItem(ctx context.Context, input *types.MealListItemDatabaseCreationInput) (*types.MealListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.MealListItemIDKey, input.ID)
	logger := q.logger.WithValue(keys.MealListItemIDKey, input.ID)

	if err := q.generatedQuerier.CreateMealListItem(ctx, q.db, &generated.CreateMealListItemParams{
		ID:                input.ID,
		MealID:            input.MealID,
		Notes:             input.Notes,
		BelongsToMealList: input.BelongsToMealList,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal list item creation query")
	}

	x := &types.MealListItem{
		ID:                input.ID,
		MealID:            input.MealID,
		Notes:             input.Notes,
		BelongsToMealList: input.BelongsToMealList,
		CreatedAt:         q.CurrentTime(),
	}

	logger.Info("meal list item created")

	return x, nil
}

// UpdateMealListItem updates a particular meal list item.
func (q *repository) UpdateMealListItem(ctx context.Context, updated *types.MealListItem) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealListItemIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.MealListItemIDKey, updated.ID)

	rowsAffected, err := q.generatedQuerier.UpdateMealListItem(ctx, q.db, &generated.UpdateMealListItemParams{
		MealID:            updated.MealID,
		Notes:             updated.Notes,
		BelongsToMealList: updated.BelongsToMealList,
		ID:                updated.ID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal list item")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Info("meal list item updated")

	return nil
}

// ArchiveMealListItem archives a meal list item from the database by its ID.
func (q *repository) ArchiveMealListItem(ctx context.Context, mealListItemID, mealListID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealListID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealListIDKey, mealListID)
	tracing.AttachToSpan(span, keys.MealListIDKey, mealListID)

	if mealListItemID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealListItemIDKey, mealListItemID)
	tracing.AttachToSpan(span, keys.MealListItemIDKey, mealListItemID)

	rowsAffected, err := q.generatedQuerier.ArchiveMealListItem(ctx, q.db, &generated.ArchiveMealListItemParams{
		BelongsToMealList: mealListID,
		ID:                mealListItemID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal list item")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	logger.Info("meal list item archived")

	return nil
}
