package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.RecipeMediaDataManager = (*Querier)(nil)
)

// RecipeMediaExists fetches whether a recipe media exists from the database.
func (q *Querier) RecipeMediaExists(ctx context.Context, recipeMediaID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeMediaID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeMediaIDKey, recipeMediaID)
	tracing.AttachRecipeMediaIDToSpan(span, recipeMediaID)

	result, err := q.generatedQuerier.CheckRecipeMediaExistence(ctx, q.db, recipeMediaID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe media existence check")
	}

	return result, nil
}

// GetRecipeMedia fetches a recipe media from the database.
func (q *Querier) GetRecipeMedia(ctx context.Context, recipeMediaID string) (*types.RecipeMedia, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeMediaID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeMediaIDKey, recipeMediaID)
	tracing.AttachRecipeMediaIDToSpan(span, recipeMediaID)

	result, err := q.generatedQuerier.GetRecipeMedia(ctx, q.db, recipeMediaID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe media")
	}

	recipeMedia := &types.RecipeMedia{
		CreatedAt:           result.CreatedAt,
		ArchivedAt:          timePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:       timePointerFromNullTime(result.LastUpdatedAt),
		ID:                  result.ID,
		BelongsToRecipe:     stringPointerFromNullString(result.BelongsToRecipe),
		BelongsToRecipeStep: stringPointerFromNullString(result.BelongsToRecipeStep),
		MimeType:            result.MimeType,
		InternalPath:        result.InternalPath,
		ExternalPath:        result.ExternalPath,
		Index:               uint16(result.Index),
	}

	return recipeMedia, nil
}

// getRecipeMediaForRecipe fetches a list of recipe media from the database that meet a particular filter.
func (q *Querier) getRecipeMediaForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeMedia, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	results, err := q.generatedQuerier.GetRecipeMediaForRecipe(ctx, q.db, nullStringFromString(recipeID))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe media list retrieval query")
	}

	recipeMedia := make([]*types.RecipeMedia, len(results))
	for i, result := range results {
		recipeMedia[i] = &types.RecipeMedia{
			CreatedAt:           result.CreatedAt,
			ArchivedAt:          timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:       timePointerFromNullTime(result.LastUpdatedAt),
			ID:                  result.ID,
			BelongsToRecipe:     stringPointerFromNullString(result.BelongsToRecipe),
			BelongsToRecipeStep: stringPointerFromNullString(result.BelongsToRecipeStep),
			MimeType:            result.MimeType,
			InternalPath:        result.InternalPath,
			ExternalPath:        result.ExternalPath,
			Index:               uint16(result.Index),
		}
	}

	return recipeMedia, nil
}

// getRecipeMediaForRecipeStep fetches a list of recipe media from the database that meet a particular filter.
func (q *Querier) getRecipeMediaForRecipeStep(ctx context.Context, recipeID, recipeStepID string) ([]*types.RecipeMedia, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	results, err := q.generatedQuerier.GetRecipeMediaForRecipeStep(ctx, q.db, &generated.GetRecipeMediaForRecipeStepParams{
		BelongsToRecipe:     nullStringFromString(recipeID),
		BelongsToRecipeStep: nullStringFromString(recipeStepID),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe media list retrieval query")
	}

	recipeMedia := []*types.RecipeMedia{}
	for _, result := range results {
		recipeMedia = append(recipeMedia, &types.RecipeMedia{
			CreatedAt:           result.CreatedAt,
			ArchivedAt:          timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:       timePointerFromNullTime(result.LastUpdatedAt),
			ID:                  result.ID,
			BelongsToRecipe:     stringPointerFromNullString(result.BelongsToRecipe),
			BelongsToRecipeStep: stringPointerFromNullString(result.BelongsToRecipeStep),
			MimeType:            result.MimeType,
			InternalPath:        result.InternalPath,
			ExternalPath:        result.ExternalPath,
			Index:               uint16(result.Index),
		})
	}

	return recipeMedia, nil
}

// CreateRecipeMedia creates a recipe media in the database.
func (q *Querier) CreateRecipeMedia(ctx context.Context, input *types.RecipeMediaDatabaseCreationInput) (*types.RecipeMedia, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeMediaIDKey, input.ID)

	// create the recipe media.
	if err := q.generatedQuerier.CreateRecipeMedia(ctx, q.db, &generated.CreateRecipeMediaParams{
		ID:                  input.ID,
		MimeType:            input.MimeType,
		InternalPath:        input.InternalPath,
		ExternalPath:        input.ExternalPath,
		BelongsToRecipe:     nullStringFromStringPointer(input.BelongsToRecipe),
		BelongsToRecipeStep: nullStringFromStringPointer(input.BelongsToRecipeStep),
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
		CreatedAt:           q.currentTime(),
	}

	tracing.AttachRecipeMediaIDToSpan(span, x.ID)
	logger.Info("recipe media created")

	return x, nil
}

// UpdateRecipeMedia updates a particular recipe media.
func (q *Querier) UpdateRecipeMedia(ctx context.Context, updated *types.RecipeMedia) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeMediaIDKey, updated.ID)
	tracing.AttachRecipeMediaIDToSpan(span, updated.ID)

	if _, err := q.generatedQuerier.UpdateRecipeMedia(ctx, q.db, &generated.UpdateRecipeMediaParams{
		BelongsToRecipe:     nullStringFromStringPointer(updated.BelongsToRecipe),
		BelongsToRecipeStep: nullStringFromStringPointer(updated.BelongsToRecipeStep),
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
func (q *Querier) ArchiveRecipeMedia(ctx context.Context, recipeMediaID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeMediaID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeMediaIDKey, recipeMediaID)
	tracing.AttachRecipeMediaIDToSpan(span, recipeMediaID)

	if _, err := q.generatedQuerier.ArchiveRecipeMedia(ctx, q.db, recipeMediaID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe media")
	}

	logger.Info("recipe media archived")

	return nil
}
