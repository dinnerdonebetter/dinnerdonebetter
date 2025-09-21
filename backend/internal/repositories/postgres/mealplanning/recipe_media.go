package mealplanning

import (
	"context"
	"database/sql"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ types.RecipeMediaDataManager = (*repository)(nil)
)

// RecipeMediaExists fetches whether a recipe media exists from the database.
func (r *repository) RecipeMediaExists(ctx context.Context, recipeMediaID string) (exists bool, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeMediaID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeMediaIDKey, recipeMediaID)
	tracing.AttachToSpan(span, keys.RecipeMediaIDKey, recipeMediaID)

	result, err := r.generatedQuerier.CheckRecipeMediaExistence(ctx, r.db, recipeMediaID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe media existence check")
	}

	return result, nil
}

// GetRecipeMedia fetches a recipe media from the database.
func (r *repository) GetRecipeMedia(ctx context.Context, recipeMediaID string) (*types.RecipeMedia, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeMediaID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeMediaIDKey, recipeMediaID)
	tracing.AttachToSpan(span, keys.RecipeMediaIDKey, recipeMediaID)

	result, err := r.generatedQuerier.GetRecipeMedia(ctx, r.db, recipeMediaID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe media")
	}

	recipeMedia := &types.RecipeMedia{
		CreatedAt:           result.CreatedAt,
		ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
		ID:                  result.ID,
		BelongsToRecipe:     database.StringPointerFromNullString(result.BelongsToRecipe),
		BelongsToRecipeStep: database.StringPointerFromNullString(result.BelongsToRecipeStep),
		MimeType:            result.MimeType,
		InternalPath:        result.InternalPath,
		ExternalPath:        result.ExternalPath,
		Index:               uint16(result.Index),
	}

	return recipeMedia, nil
}

// getRecipeMediaForRecipe fetches a list of recipe media from the database that meet a particular filter.
func (r *repository) getRecipeMediaForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeMedia, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	results, err := r.generatedQuerier.GetRecipeMediaForRecipe(ctx, r.db, database.NullStringFromString(recipeID))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe media list retrieval query")
	}

	recipeMedia := make([]*types.RecipeMedia, len(results))
	for i, result := range results {
		recipeMedia[i] = &types.RecipeMedia{
			CreatedAt:           result.CreatedAt,
			ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
			ID:                  result.ID,
			BelongsToRecipe:     database.StringPointerFromNullString(result.BelongsToRecipe),
			BelongsToRecipeStep: database.StringPointerFromNullString(result.BelongsToRecipeStep),
			MimeType:            result.MimeType,
			InternalPath:        result.InternalPath,
			ExternalPath:        result.ExternalPath,
			Index:               uint16(result.Index),
		}
	}

	return recipeMedia, nil
}

// getRecipeMediaForRecipeStep fetches a list of recipe media from the database that meet a particular filter.
func (r *repository) getRecipeMediaForRecipeStep(ctx context.Context, recipeID, recipeStepID string) ([]*types.RecipeMedia, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	results, err := r.generatedQuerier.GetRecipeMediaForRecipeStep(ctx, r.db, &generated.GetRecipeMediaForRecipeStepParams{
		RecipeID:     database.NullStringFromString(recipeID),
		RecipeStepID: database.NullStringFromString(recipeStepID),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe media list retrieval query")
	}

	recipeMedia := []*types.RecipeMedia{}
	for _, result := range results {
		recipeMedia = append(recipeMedia, &types.RecipeMedia{
			CreatedAt:           result.CreatedAt,
			ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
			ID:                  result.ID,
			BelongsToRecipe:     database.StringPointerFromNullString(result.BelongsToRecipe),
			BelongsToRecipeStep: database.StringPointerFromNullString(result.BelongsToRecipeStep),
			MimeType:            result.MimeType,
			InternalPath:        result.InternalPath,
			ExternalPath:        result.ExternalPath,
			Index:               uint16(result.Index),
		})
	}

	return recipeMedia, nil
}

// CreateRecipeMedia creates a recipe media in the database.
func (r *repository) CreateRecipeMedia(ctx context.Context, input *types.RecipeMediaDatabaseCreationInput) (*types.RecipeMedia, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	logger := r.logger.WithValue(keys.RecipeMediaIDKey, input.ID)

	// create the recipe media.
	if err := r.generatedQuerier.CreateRecipeMedia(ctx, r.db, &generated.CreateRecipeMediaParams{
		ID:                  input.ID,
		MimeType:            input.MimeType,
		InternalPath:        input.InternalPath,
		ExternalPath:        input.ExternalPath,
		BelongsToRecipe:     database.NullStringFromStringPointer(input.BelongsToRecipe),
		BelongsToRecipeStep: database.NullStringFromStringPointer(input.BelongsToRecipeStep),
		Index:               int32(input.Index),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe media creation query")
	}

	x := &types.RecipeMedia{
		ID:                  input.ID,
		BelongsToRecipe:     input.BelongsToRecipe,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		MimeType:            input.MimeType,
		InternalPath:        input.InternalPath,
		ExternalPath:        input.ExternalPath,
		Index:               input.Index,
		CreatedAt:           r.CurrentTime(),
	}

	tracing.AttachToSpan(span, keys.RecipeMediaIDKey, x.ID)
	logger.Info("recipe media created")

	return x, nil
}

// UpdateRecipeMedia updates a particular recipe media.
func (r *repository) UpdateRecipeMedia(ctx context.Context, updated *types.RecipeMedia) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.RecipeMediaIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeMediaIDKey, updated.ID)

	if _, err := r.generatedQuerier.UpdateRecipeMedia(ctx, r.db, &generated.UpdateRecipeMediaParams{
		BelongsToRecipe:     database.NullStringFromStringPointer(updated.BelongsToRecipe),
		BelongsToRecipeStep: database.NullStringFromStringPointer(updated.BelongsToRecipeStep),
		MimeType:            updated.MimeType,
		InternalPath:        updated.InternalPath,
		ExternalPath:        updated.ExternalPath,
		Index:               int32(updated.Index),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe media")
	}

	logger.Info("recipe media updated")

	return nil
}

// ArchiveRecipeMedia archives a recipe media from the database by its ID.
func (r *repository) ArchiveRecipeMedia(ctx context.Context, recipeMediaID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if recipeMediaID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeMediaIDKey, recipeMediaID)
	tracing.AttachToSpan(span, keys.RecipeMediaIDKey, recipeMediaID)

	rowsAffected, err := r.generatedQuerier.ArchiveRecipeMedia(ctx, r.db, recipeMediaID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe media")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
