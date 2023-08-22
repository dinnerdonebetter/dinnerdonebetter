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

var (
	_ types.RecipeMediaDataManager = (*Querier)(nil)

	// recipeMediaTableColumns are the columns for the recipe_media table.
	recipeMediaTableColumns = []string{
		"recipe_media.id",
		"recipe_media.belongs_to_recipe",
		"recipe_media.belongs_to_recipe_step",
		"recipe_media.mime_type",
		"recipe_media.internal_path",
		"recipe_media.external_path",
		"recipe_media.index",
		"recipe_media.created_at",
		"recipe_media.last_updated_at",
		"recipe_media.archived_at",
	}
)

// scanPieceOfRecipeMedia takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe media struct.
func (q *Querier) scanPieceOfRecipeMedia(ctx context.Context, scan database.Scanner) (x *types.RecipeMedia, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeMedia{}

	targetVars := []any{
		&x.ID,
		&x.BelongsToRecipe,
		&x.BelongsToRecipeStep,
		&x.MimeType,
		&x.InternalPath,
		&x.ExternalPath,
		&x.Index,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, observability.PrepareError(err, span, "")
	}

	return x, nil
}

// scanRecipeMedia takes some database rows and turns them into a slice of recipe media.
func (q *Querier) scanRecipeMedia(ctx context.Context, rows database.ResultIterator) (recipeMedias []*types.RecipeMedia, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, scanErr := q.scanPieceOfRecipeMedia(ctx, rows)
		if scanErr != nil {
			return nil, scanErr
		}

		recipeMedias = append(recipeMedias, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "handling rows")
	}

	return recipeMedias, nil
}

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

//go:embed queries/recipe_media/for_recipe.sql
var recipeMediaForRecipeQuery string

// getRecipeMediaForRecipe fetches a list of recipe media from the database that meet a particular filter.
func (q *Querier) getRecipeMediaForRecipe(ctx context.Context, recipeID string) (x []*types.RecipeMedia, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	recipeMediaForRecipeArgs := []any{
		recipeID,
	}

	rows, err := q.getRows(ctx, q.db, "recipe media for recipe", recipeMediaForRecipeQuery, recipeMediaForRecipeArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe media list retrieval query")
	}

	if x, err = q.scanRecipeMedia(ctx, rows); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe media")
	}

	return x, nil
}

//go:embed queries/recipe_media/for_recipe_step.sql
var recipeMediaForRecipeStepQuery string

// getRecipeMediaForRecipeStep fetches a list of recipe media from the database that meet a particular filter.
func (q *Querier) getRecipeMediaForRecipeStep(ctx context.Context, recipeID, recipeStepID string) (x []*types.RecipeMedia, err error) {
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

	recipeMediaForRecipeStepArgs := []any{
		recipeID,
		recipeStepID,
	}

	rows, err := q.getRows(ctx, q.db, "recipe media for recipe step", recipeMediaForRecipeStepQuery, recipeMediaForRecipeStepArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe media list retrieval query")
	}

	if x, err = q.scanRecipeMedia(ctx, rows); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe media")
	}

	return x, nil
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

	if err := q.generatedQuerier.UpdateRecipeMedia(ctx, q.db, &generated.UpdateRecipeMediaParams{
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

	if err := q.generatedQuerier.ArchiveRecipeMedia(ctx, q.db, recipeMediaID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe media")
	}

	logger.Info("recipe media archived")

	return nil
}
