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
	_ types.RecipeRatingDataManager = (*Querier)(nil)

	// recipeRatingsTableColumns are the columns for the recipe_ratings table.
	recipeRatingsTableColumns = []string{
		"recipe_ratings.id",
		"recipe_ratings.recipe_id",
		"recipe_ratings.taste",
		"recipe_ratings.difficulty",
		"recipe_ratings.cleanup",
		"recipe_ratings.instructions",
		"recipe_ratings.overall",
		"recipe_ratings.notes",
		"recipe_ratings.by_user",
		"recipe_ratings.created_at",
		"recipe_ratings.last_updated_at",
		"recipe_ratings.archived_at",
	}
)

// scanRecipeRating takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe rating struct.
func (q *Querier) scanRecipeRating(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeRating, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeRating{}

	targetVars := []any{
		&x.ID,
		&x.RecipeID,
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

// scanRecipeRatings takes some database rows and turns them into a slice of recipe ratings.
func (q *Querier) scanRecipeRatings(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeRatings []*types.RecipeRating, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeRating(ctx, rows, includeCounts)
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

		recipeRatings = append(recipeRatings, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return recipeRatings, filteredCount, totalCount, nil
}

//go:embed queries/recipe_ratings/exists.sql
var recipeRatingExistenceQuery string

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

	args := []any{
		recipeRatingID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeRatingExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe rating existence check")
	}

	return result, nil
}

//go:embed queries/recipe_ratings/get_one.sql
var getRecipeRatingQuery string

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

	args := []any{
		recipeRatingID,
	}

	row := q.getOneRow(ctx, q.db, "recipe rating", getRecipeRatingQuery, args)

	recipeRating, _, _, err := q.scanRecipeRating(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeRating")
	}

	return recipeRating, nil
}

// GetRecipeRatings fetches a list of recipe ratings from the database that meet a particular filter.
func (q *Querier) GetRecipeRatings(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeRating], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.RecipeRating]{}
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

	query, args := q.buildListQuery(ctx, "recipe_ratings", nil, nil, nil, "", recipeRatingsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "recipe ratings", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe ratings list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanRecipeRatings(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe ratings")
	}

	return x, nil
}

//go:embed queries/recipe_ratings/create.sql
var recipeRatingCreationQuery string

// CreateRecipeRating creates a recipe rating in the database.
func (q *Querier) CreateRecipeRating(ctx context.Context, input *types.RecipeRatingDatabaseCreationInput) (*types.RecipeRating, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeRatingIDKey, input.ID)

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

	// create the recipe rating.
	if err := q.performWriteQuery(ctx, q.db, "recipe rating creation", recipeRatingCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe rating creation query")
	}

	x := &types.RecipeRating{
		ID:           input.ID,
		RecipeID:     input.MealID,
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

//go:embed queries/recipe_ratings/update.sql
var updateRecipeRatingQuery string

// UpdateRecipeRating updates a particular recipe rating.
func (q *Querier) UpdateRecipeRating(ctx context.Context, updated *types.RecipeRating) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeRatingIDKey, updated.ID)
	tracing.AttachRecipeRatingIDToSpan(span, updated.ID)

	args := []any{
		updated.RecipeID,
		updated.Taste,
		updated.Difficulty,
		updated.Cleanup,
		updated.Instructions,
		updated.Overall,
		updated.Notes,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe rating update", updateRecipeRatingQuery, args); err != nil {
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
