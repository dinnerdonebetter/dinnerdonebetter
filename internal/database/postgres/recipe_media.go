package postgres

import (
	"context"
	_ "embed"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
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

// scanPieceOfRecipeMedia takes a database Scanner (i.e. *sql.Row) and scans the result into a valid preparation struct.
func (q *Querier) scanPieceOfRecipeMedia(ctx context.Context, scan database.Scanner) (x *types.RecipeMedia, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeMedia{}

	targetVars := []interface{}{
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

// scanRecipeMedia takes some database rows and turns them into a slice of valid preparations.
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

//go:embed queries/recipe_media/exists.sql
var recipeMediaExistenceQuery string

// RecipeMediaExists fetches whether a valid preparation exists from the database.
func (q *Querier) RecipeMediaExists(ctx context.Context, recipeMediaID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeMediaID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeMediaIDKey, recipeMediaID)
	tracing.AttachRecipeMediaIDToSpan(span, recipeMediaID)

	args := []interface{}{
		recipeMediaID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeMediaExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid preparation existence check")
	}

	return result, nil
}

//go:embed queries/recipe_media/get_one.sql
var getRecipeMediaQuery string

// GetRecipeMedia fetches a valid preparation from the database.
func (q *Querier) GetRecipeMedia(ctx context.Context, recipeMediaID string) (*types.RecipeMedia, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeMediaID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeMediaIDKey, recipeMediaID)
	tracing.AttachRecipeMediaIDToSpan(span, recipeMediaID)

	args := []interface{}{
		recipeMediaID,
	}

	row := q.getOneRow(ctx, q.db, "recipeMedia", getRecipeMediaQuery, args)

	recipeMedia, err := q.scanPieceOfRecipeMedia(ctx, row)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeMedia")
	}

	return recipeMedia, nil
}

//go:embed queries/recipe_media/for_recipe.sql
var recipeMediaForRecipeQuery string

// GetRecipeMediaForRecipe fetches a list of valid preparations from the database that meet a particular filter.
func (q *Querier) GetRecipeMediaForRecipe(ctx context.Context, recipeID string) (x []*types.RecipeMedia, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	recipeMediaForRecipeArgs := []interface{}{
		recipeID,
	}

	rows, err := q.getRows(ctx, q.db, "recipe media", recipeMediaForRecipeQuery, recipeMediaForRecipeArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparations list retrieval query")
	}

	if x, err = q.scanRecipeMedia(ctx, rows); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid preparations")
	}

	return x, nil
}

//go:embed queries/recipe_media/create.sql
var recipeMediaCreationQuery string

// CreateRecipeMedia creates a valid preparation in the database.
func (q *Querier) CreateRecipeMedia(ctx context.Context, input *types.RecipeMediaDatabaseCreationInput) (*types.RecipeMedia, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeMediaIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.BelongsToRecipe,
		input.BelongsToRecipeStep,
		input.MimeType,
		input.InternalPath,
		input.ExternalPath,
		input.Index,
	}

	// create the valid preparation.
	if err := q.performWriteQuery(ctx, q.db, "valid preparation creation", recipeMediaCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid preparation creation query")
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
	logger.Info("valid preparation created")

	return x, nil
}

//go:embed queries/recipe_media/update.sql
var updateRecipeMediaQuery string

// UpdateRecipeMedia updates a particular valid preparation.
func (q *Querier) UpdateRecipeMedia(ctx context.Context, updated *types.RecipeMedia) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeMediaIDKey, updated.ID)
	tracing.AttachRecipeMediaIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.BelongsToRecipe,
		updated.BelongsToRecipeStep,
		updated.MimeType,
		updated.InternalPath,
		updated.ExternalPath,
		updated.Index,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation update", updateRecipeMediaQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation updated")

	return nil
}

//go:embed queries/recipe_media/archive.sql
var archiveRecipeMediaQuery string

// ArchiveRecipeMedia archives a valid preparation from the database by its ID.
func (q *Querier) ArchiveRecipeMedia(ctx context.Context, recipeMediaID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeMediaID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeMediaIDKey, recipeMediaID)
	tracing.AttachRecipeMediaIDToSpan(span, recipeMediaID)

	args := []interface{}{
		recipeMediaID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation archive", archiveRecipeMediaQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation archived")

	return nil
}
