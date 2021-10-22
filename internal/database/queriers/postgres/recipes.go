package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var (
	_ types.RecipeDataManager = (*SQLQuerier)(nil)

	// recipesTableColumns are the columns for the recipes table.
	recipesTableColumns = []string{
		"recipes.id",
		"recipes.name",
		"recipes.source",
		"recipes.description",
		"recipes.inspired_by_recipe_id",
		"recipes.created_on",
		"recipes.last_updated_on",
		"recipes.archived_on",
		"recipes.belongs_to_household",
	}
)

// scanRecipe takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe struct.
func (q *SQLQuerier) scanRecipe(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.Recipe, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.Recipe{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Source,
		&x.Description,
		&x.InspiredByRecipeID,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToHousehold,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipes takes some database rows and turns them into a slice of recipes.
func (q *SQLQuerier) scanRecipes(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipes []*types.Recipe, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipe(ctx, rows, includeCounts)
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

		recipes = append(recipes, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipes, filteredCount, totalCount, nil
}

const recipeExistenceQuery = "SELECT EXISTS ( SELECT recipes.id FROM recipes WHERE recipes.archived_on IS NULL AND recipes.id = $1 )"

// RecipeExists fetches whether a recipe exists from the database.
func (q *SQLQuerier) RecipeExists(ctx context.Context, recipeID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	args := []interface{}{
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe existence check")
	}

	return result, nil
}

const getRecipeQuery = "SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_household FROM recipes WHERE recipes.archived_on IS NULL AND recipes.id = $1"

// GetRecipe fetches a recipe from the database.
func (q *SQLQuerier) GetRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	args := []interface{}{
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "recipe", getRecipeQuery, args)

	recipe, _, _, err := q.scanRecipe(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe")
	}

	return recipe, nil
}

const getTotalRecipesCountQuery = "SELECT COUNT(recipes.id) FROM recipes WHERE recipes.archived_on IS NULL"

// GetTotalRecipeCount fetches the count of recipes from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalRecipeCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, getTotalRecipesCountQuery, "fetching count of recipes")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of recipes")
	}

	return count, nil
}

// GetRecipes fetches a list of recipes from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipes(ctx context.Context, filter *types.QueryFilter) (x *types.RecipeList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.RecipeList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(
		ctx,
		"recipes",
		nil,
		nil,
		householdOwnershipColumn,
		recipesTableColumns,
		"",
		false,
		filter,
	)

	rows, err := q.performReadQuery(ctx, q.db, "recipes", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipes list retrieval query")
	}

	if x.Recipes, x.FilteredCount, x.TotalCount, err = q.scanRecipes(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipes")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetRecipesWithIDsQuery(ctx context.Context, householdID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"recipes.id":          ids,
		"recipes.archived_on": nil,
	}

	if householdID != "" {
		withIDsWhere["recipes.belongs_to_household"] = householdID
	}

	subqueryBuilder := q.sqlBuilder.Select(recipesTableColumns...).
		From("recipes").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(recipesTableColumns...).
		FromSelect(subqueryBuilder, "recipes").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetRecipesWithIDs fetches recipes from the database within a given set of IDs.
func (q *SQLQuerier) GetRecipesWithIDs(ctx context.Context, householdID string, limit uint8, ids []string) ([]*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if ids == nil {
		return nil, ErrNilInputProvided
	}

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.buildGetRecipesWithIDsQuery(ctx, householdID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "recipes with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipes from database")
	}

	recipes, _, _, err := q.scanRecipes(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipes")
	}

	return recipes, nil
}

const recipeCreationQuery = "INSERT INTO recipes (id,name,source,description,inspired_by_recipe_id,belongs_to_household) VALUES ($1,$2,$3,$4,$5,$6)"

// CreateRecipe creates a recipe in the database.
func (q *SQLQuerier) CreateRecipe(ctx context.Context, input *types.RecipeDatabaseCreationInput) (*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Name,
		input.Source,
		input.Description,
		input.InspiredByRecipeID,
		input.BelongsToHousehold,
	}

	// create the recipe.
	if err := q.performWriteQuery(ctx, q.db, "recipe creation", recipeCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating recipe")
	}

	x := &types.Recipe{
		ID:                 input.ID,
		Name:               input.Name,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedOn:          q.currentTime(),
	}

	tracing.AttachRecipeIDToSpan(span, x.ID)
	logger.Info("recipe created")

	return x, nil
}

const updateRecipeQuery = "UPDATE recipes SET name = $1, source = $2, description = $3, inspired_by_recipe_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_household = $5 AND id = $6"

// UpdateRecipe updates a particular recipe.
func (q *SQLQuerier) UpdateRecipe(ctx context.Context, updated *types.Recipe) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeIDKey, updated.ID)
	tracing.AttachRecipeIDToSpan(span, updated.ID)
	tracing.AttachHouseholdIDToSpan(span, updated.BelongsToHousehold)

	args := []interface{}{
		updated.Name,
		updated.Source,
		updated.Description,
		updated.InspiredByRecipeID,
		updated.BelongsToHousehold,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe update", updateRecipeQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe")
	}

	logger.Info("recipe updated")

	return nil
}

const archiveRecipeQuery = "UPDATE recipes SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_household = $1 AND id = $2"

// ArchiveRecipe archives a recipe from the database by its ID.
func (q *SQLQuerier) ArchiveRecipe(ctx context.Context, recipeID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		householdID,
		recipeID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe archive", archiveRecipeQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe")
	}

	logger.Info("recipe archived")

	return nil
}
