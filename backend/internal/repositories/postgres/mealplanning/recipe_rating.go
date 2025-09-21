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
	_ types.RecipeRatingDataManager = (*repository)(nil)
)

// RecipeRatingExists fetches whether a recipe rating exists from the database.
func (r *repository) RecipeRatingExists(ctx context.Context, recipeID, recipeRatingID string) (exists bool, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeRatingID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRatingID)

	result, err := r.generatedQuerier.CheckRecipeRatingExistence(ctx, r.db, recipeRatingID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe rating existence check")
	}

	return result, nil
}

// GetRecipeRating fetches a recipe rating from the database.
func (r *repository) GetRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*types.RecipeRating, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeRatingID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRatingID)

	result, err := r.generatedQuerier.GetRecipeRating(ctx, r.db, recipeRatingID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe rating")
	}

	recipeRating := &types.RecipeRating{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		Notes:         result.Notes,
		ID:            result.ID,
		RecipeID:      result.RecipeID,
		ByUser:        result.ByUser,
		Taste:         database.Float32FromNullString(result.Taste),
		Instructions:  database.Float32FromNullString(result.Instructions),
		Overall:       database.Float32FromNullString(result.Overall),
		Cleanup:       database.Float32FromNullString(result.Cleanup),
		Difficulty:    database.Float32FromNullString(result.Difficulty),
	}

	return recipeRating, nil
}

// GetRecipeRatingsForRecipe fetches a list of recipe ratings from the database that meet a particular filter.
func (r *repository) GetRecipeRatingsForRecipe(ctx context.Context, recipeID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.RecipeRating], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	x = &filtering.QueryFilteredResult[types.RecipeRating]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetRecipeRatingsForRecipe(ctx, r.db, &generated.GetRecipeRatingsForRecipeParams{
		RecipeID:        recipeID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe ratings list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.RecipeRating{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Notes:         result.Notes,
			ID:            result.ID,
			RecipeID:      result.RecipeID,
			ByUser:        result.ByUser,
			Taste:         database.Float32FromNullString(result.Taste),
			Instructions:  database.Float32FromNullString(result.Instructions),
			Overall:       database.Float32FromNullString(result.Overall),
			Cleanup:       database.Float32FromNullString(result.Cleanup),
			Difficulty:    database.Float32FromNullString(result.Difficulty),
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetRecipeRatingsForUser fetches a list of recipe ratings from the database that meet a particular filter.
func (r *repository) GetRecipeRatingsForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.RecipeRating], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if userID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	x = &filtering.QueryFilteredResult[types.RecipeRating]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetRecipeRatingsForUser(ctx, r.db, &generated.GetRecipeRatingsForUserParams{
		ByUser:          userID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe ratings list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.RecipeRating{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Notes:         result.Notes,
			ID:            result.ID,
			RecipeID:      result.RecipeID,
			ByUser:        result.ByUser,
			Taste:         database.Float32FromNullString(result.Taste),
			Instructions:  database.Float32FromNullString(result.Instructions),
			Overall:       database.Float32FromNullString(result.Overall),
			Cleanup:       database.Float32FromNullString(result.Cleanup),
			Difficulty:    database.Float32FromNullString(result.Difficulty),
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateRecipeRating creates a recipe rating in the database.
func (r *repository) CreateRecipeRating(ctx context.Context, input *types.RecipeRatingDatabaseCreationInput) (*types.RecipeRating, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	logger := r.logger.WithValue(keys.RecipeRatingIDKey, input.ID)

	// create the recipe rating.
	if err := r.generatedQuerier.CreateRecipeRating(ctx, r.db, &generated.CreateRecipeRatingParams{
		ID:           input.ID,
		RecipeID:     input.RecipeID,
		Notes:        input.Notes,
		ByUser:       input.ByUser,
		Taste:        database.NullStringFromFloat32(input.Taste),
		Difficulty:   database.NullStringFromFloat32(input.Difficulty),
		Cleanup:      database.NullStringFromFloat32(input.Cleanup),
		Instructions: database.NullStringFromFloat32(input.Instructions),
		Overall:      database.NullStringFromFloat32(input.Overall),
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
		CreatedAt:    r.CurrentTime(),
	}

	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, x.ID)
	logger.Info("recipe rating created")

	return x, nil
}

// UpdateRecipeRating updates a particular recipe rating.
func (r *repository) UpdateRecipeRating(ctx context.Context, updated *types.RecipeRating) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.RecipeRatingIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, updated.ID)

	if _, err := r.generatedQuerier.UpdateRecipeRating(ctx, r.db, &generated.UpdateRecipeRatingParams{
		RecipeID:     updated.RecipeID,
		Taste:        database.NullStringFromFloat32(updated.Taste),
		Difficulty:   database.NullStringFromFloat32(updated.Difficulty),
		Cleanup:      database.NullStringFromFloat32(updated.Cleanup),
		Instructions: database.NullStringFromFloat32(updated.Instructions),
		Overall:      database.NullStringFromFloat32(updated.Overall),
		Notes:        updated.Notes,
		ID:           updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe rating")
	}

	logger.Info("recipe rating updated")

	return nil
}

// ArchiveRecipeRating archives a recipe rating from the database by its ID.
func (r *repository) ArchiveRecipeRating(ctx context.Context, recipeID, recipeRatingID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeRatingID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRatingID)

	rowsAffected, err := r.generatedQuerier.ArchiveRecipeRating(ctx, r.db, recipeRatingID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe rating")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
