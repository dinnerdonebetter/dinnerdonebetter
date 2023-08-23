package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.RecipeRatingDataManager = (*Querier)(nil)
)

// RecipeRatingExists fetches whether a recipe rating exists from the database.
func (q *Querier) RecipeRatingExists(ctx context.Context, recipeRatingID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeRatingID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachRecipeRatingIDToSpan(span, recipeRatingID)

	result, err := q.generatedQuerier.CheckRecipeRatingExistence(ctx, q.db, recipeRatingID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe rating existence check")
	}

	return result, nil
}

// GetRecipeRating fetches a recipe rating from the database.
func (q *Querier) GetRecipeRating(ctx context.Context, recipeRatingID string) (*types.RecipeRating, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeRatingID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachRecipeRatingIDToSpan(span, recipeRatingID)

	result, err := q.generatedQuerier.GetRecipeRating(ctx, q.db, recipeRatingID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe rating existence check")
	}

	recipeRating := &types.RecipeRating{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
		Notes:         result.Notes,
		ID:            result.ID,
		RecipeID:      result.RecipeID,
		ByUser:        result.ByUser,
		Taste:         float32FromNullString(result.Taste),
		Instructions:  float32FromNullString(result.Instructions),
		Overall:       float32FromNullString(result.Overall),
		Cleanup:       float32FromNullString(result.Cleanup),
		Difficulty:    float32FromNullString(result.Difficulty),
	}

	return recipeRating, nil
}

// GetRecipeRatings fetches a list of recipe ratings from the database that meet a particular filter.
func (q *Querier) GetRecipeRatings(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeRating], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.RecipeRating]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeRatings(ctx, q.db, &generated.GetRecipeRatingsParams{
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe ratings list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.RecipeRating{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
			Notes:         result.Notes,
			ID:            result.ID,
			RecipeID:      result.RecipeID,
			ByUser:        result.ByUser,
			Taste:         float32FromNullString(result.Taste),
			Instructions:  float32FromNullString(result.Instructions),
			Overall:       float32FromNullString(result.Overall),
			Cleanup:       float32FromNullString(result.Cleanup),
			Difficulty:    float32FromNullString(result.Difficulty),
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateRecipeRating creates a recipe rating in the database.
func (q *Querier) CreateRecipeRating(ctx context.Context, input *types.RecipeRatingDatabaseCreationInput) (*types.RecipeRating, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeRatingIDKey, input.ID)

	// create the recipe rating.
	if err := q.generatedQuerier.CreateRecipeRating(ctx, q.db, &generated.CreateRecipeRatingParams{
		ID:           input.ID,
		RecipeID:     input.RecipeID,
		Notes:        input.Notes,
		ByUser:       input.ByUser,
		Taste:        nullStringFromFloat32(input.Taste),
		Difficulty:   nullStringFromFloat32(input.Difficulty),
		Cleanup:      nullStringFromFloat32(input.Cleanup),
		Instructions: nullStringFromFloat32(input.Instructions),
		Overall:      nullStringFromFloat32(input.Overall),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe rating creation query")
	}

	x := &types.RecipeRating{
		ID:           input.ID,
		RecipeID:     input.RecipeID,
		Taste:        input.Taste,
		Difficulty:   input.Difficulty,
		Cleanup:      input.Cleanup,
		Instructions: input.Instructions,
		Overall:      input.Overall,
		Notes:        input.Notes,
		ByUser:       input.ByUser,
		CreatedAt:    q.currentTime(),
	}

	tracing.AttachRecipeRatingIDToSpan(span, x.ID)
	logger.Info("recipe rating created")

	return x, nil
}

// UpdateRecipeRating updates a particular recipe rating.
func (q *Querier) UpdateRecipeRating(ctx context.Context, updated *types.RecipeRating) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeRatingIDKey, updated.ID)
	tracing.AttachRecipeRatingIDToSpan(span, updated.ID)

	if err := q.generatedQuerier.UpdateRecipeRating(ctx, q.db, &generated.UpdateRecipeRatingParams{
		RecipeID:     updated.RecipeID,
		Taste:        nullStringFromFloat32(updated.Taste),
		Difficulty:   nullStringFromFloat32(updated.Difficulty),
		Cleanup:      nullStringFromFloat32(updated.Cleanup),
		Instructions: nullStringFromFloat32(updated.Instructions),
		Overall:      nullStringFromFloat32(updated.Overall),
		Notes:        updated.Notes,
		ID:           updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe rating")
	}

	logger.Info("recipe rating updated")

	return nil
}

// ArchiveRecipeRating archives a recipe rating from the database by its ID.
func (q *Querier) ArchiveRecipeRating(ctx context.Context, recipeRatingID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeRatingID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachRecipeRatingIDToSpan(span, recipeRatingID)

	if err := q.generatedQuerier.ArchiveRecipeRating(ctx, q.db, recipeRatingID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe rating")
	}

	logger.Info("recipe rating archived")

	return nil
}
