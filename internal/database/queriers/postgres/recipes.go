package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
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
		"recipes.created_by_user",
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
		&x.CreatedByUser,
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

	logger := q.logger.Clone()

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

const getRecipeByIDQuery = `SELECT 
	recipes.id,
	recipes.name,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.created_on,
	recipes.last_updated_on,
	recipes.archived_on,
	recipes.created_by_user,
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.created_on,
	valid_preparations.last_updated_on,
	valid_preparations.archived_on,
	recipe_steps.min_estimated_time_in_seconds,
	recipe_steps.max_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.optional,
	recipe_steps.created_on,
	recipe_steps.last_updated_on,
	recipe_steps.archived_on,
	recipe_steps.belongs_to_recipe
FROM recipes
	FULL OUTER JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
	FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_on IS NULL
	AND recipes.id = $1
ORDER BY recipe_steps.index
`

const getRecipeByIDAndAuthorIDQuery = `SELECT 
	recipes.id,
	recipes.name,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.created_on,
	recipes.last_updated_on,
	recipes.archived_on,
	recipes.created_by_user,
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.created_on,
	valid_preparations.last_updated_on,
	valid_preparations.archived_on,
	recipe_steps.min_estimated_time_in_seconds,
	recipe_steps.max_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.optional,
	recipe_steps.created_on,
	recipe_steps.last_updated_on,
	recipe_steps.archived_on,
	recipe_steps.belongs_to_recipe
FROM recipes
	FULL OUTER JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
	FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_on IS NULL
	AND recipes.id = $1
	AND recipes.created_by_user = $2
ORDER BY recipe_steps.index
`

// scanRecipeAndStep takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe struct.
func (q *SQLQuerier) scanRecipeAndStep(ctx context.Context, scan database.Scanner) (x *types.Recipe, y *types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.Recipe{}
	y = &types.RecipeStep{}

	// because a recipe is allowed to have no steps, currently, we have to
	// create temporary nil values and assign them afterwards if they're not nil
	// it sucks but we are prototyping
	var (
		recipeStepID                            *string
		recipeStepIndex                         *uint32
		recipeStepPreparationID                 *string
		recipeStepPreparationName               *string
		recipeStepPreparationDescription        *string
		recipeStepPreparationIconPath           *string
		recipeStepPreparationCreatedOn          *uint64
		recipeStepPreparationLastUpdatedOn      *uint64
		recipeStepPreparationArchivedOn         *uint64
		recipeStepMinimumEstimatedTimeInSeconds *uint32
		recipeStepMaximumEstimatedTimeInSeconds *uint32
		recipeStepMinimumTemperatureInCelsius   *uint16
		recipeStepMaximumTemperatureInCelsius   *uint16
		recipeStepNotes                         *string
		recipeStepOptional                      *bool
		recipeStepCreatedOn                     *uint64
		recipeStepLastUpdatedOn                 *uint64
		recipeStepArchivedOn                    *uint64
		recipeStepBelongsToRecipe               *string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Source,
		&x.Description,
		&x.InspiredByRecipeID,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.CreatedByUser,
		&recipeStepID,
		&recipeStepIndex,
		&recipeStepPreparationID,
		&recipeStepPreparationName,
		&recipeStepPreparationDescription,
		&recipeStepPreparationIconPath,
		&recipeStepPreparationCreatedOn,
		&recipeStepPreparationLastUpdatedOn,
		&recipeStepPreparationArchivedOn,
		&recipeStepMinimumEstimatedTimeInSeconds,
		&recipeStepMaximumEstimatedTimeInSeconds,
		&recipeStepMinimumTemperatureInCelsius,
		&recipeStepMaximumTemperatureInCelsius,
		&recipeStepNotes,
		&recipeStepOptional,
		&recipeStepCreatedOn,
		&recipeStepLastUpdatedOn,
		&recipeStepArchivedOn,
		&recipeStepBelongsToRecipe,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, nil, 0, 0, observability.PrepareError(err, q.logger, span, "")
	}

	if recipeStepID != nil {
		y.ID = *recipeStepID
	}
	if recipeStepIndex != nil {
		y.Index = *recipeStepIndex
	}
	if recipeStepPreparationID != nil {
		y.Preparation.ID = *recipeStepPreparationID
	}
	if recipeStepPreparationName != nil {
		y.Preparation.Name = *recipeStepPreparationName
	}
	if recipeStepPreparationDescription != nil {
		y.Preparation.Description = *recipeStepPreparationDescription
	}
	if recipeStepPreparationIconPath != nil {
		y.Preparation.IconPath = *recipeStepPreparationIconPath
	}
	if recipeStepPreparationCreatedOn != nil {
		y.Preparation.CreatedOn = *recipeStepPreparationCreatedOn
	}
	if recipeStepPreparationLastUpdatedOn != nil {
		y.Preparation.LastUpdatedOn = recipeStepPreparationLastUpdatedOn
	}
	if recipeStepPreparationArchivedOn != nil {
		y.Preparation.ArchivedOn = recipeStepPreparationArchivedOn
	}
	if recipeStepMinimumEstimatedTimeInSeconds != nil {
		y.MinimumEstimatedTimeInSeconds = *recipeStepMinimumEstimatedTimeInSeconds
	}
	if recipeStepMaximumEstimatedTimeInSeconds != nil {
		y.MaximumEstimatedTimeInSeconds = *recipeStepMaximumEstimatedTimeInSeconds
	}
	if recipeStepMinimumTemperatureInCelsius != nil {
		y.MinimumTemperatureInCelsius = recipeStepMinimumTemperatureInCelsius
	}
	if recipeStepMaximumTemperatureInCelsius != nil {
		y.MaximumTemperatureInCelsius = recipeStepMaximumTemperatureInCelsius
	}
	if recipeStepNotes != nil {
		y.Notes = *recipeStepNotes
	}
	if recipeStepOptional != nil {
		y.Optional = *recipeStepOptional
	}
	if recipeStepCreatedOn != nil {
		y.CreatedOn = *recipeStepCreatedOn
	}
	if recipeStepLastUpdatedOn != nil {
		y.LastUpdatedOn = recipeStepLastUpdatedOn
	}
	if recipeStepArchivedOn != nil {
		y.ArchivedOn = recipeStepArchivedOn
	}
	if recipeStepBelongsToRecipe != nil {
		y.BelongsToRecipe = *recipeStepBelongsToRecipe
	}

	return x, y, filteredCount, totalCount, nil
}

// getRecipe fetches a recipe from the database.
func (q *SQLQuerier) getRecipe(ctx context.Context, recipeID, userID string) (*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	args := []interface{}{
		recipeID,
	}

	query := getRecipeByIDQuery
	if userID != "" {
		query = getRecipeByIDAndAuthorIDQuery
		args = append(args, userID)
	}

	rows, err := q.performReadQuery(ctx, q.db, "get recipe", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe")
	}

	var x *types.Recipe
	for rows.Next() {
		recipe, recipeStep, _, _, recipeScanErr := q.scanRecipeAndStep(ctx, rows)
		if recipeScanErr != nil {
			return nil, observability.PrepareError(recipeScanErr, logger, span, "scanning recipe")
		}

		if x == nil {
			x = recipe
		}

		x.Steps = append(x.Steps, recipeStep)
	}

	if x == nil {
		return nil, sql.ErrNoRows
	}

	// need to grab ingredients here and add them to steps
	ingredients, err := q.getRecipeStepIngredientsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe step ingredients for recipe")
	}

	// need to grab products here and add them to steps
	products, err := q.getRecipeStepProductsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe step products for recipe")
	}

	for i, step := range x.Steps {
		for _, ingredient := range ingredients {
			if ingredient.BelongsToRecipeStep == step.ID {
				x.Steps[i].Ingredients = append(x.Steps[i].Ingredients, ingredient)
			}
		}

		for _, product := range products {
			if product.BelongsToRecipeStep == step.ID {
				x.Steps[i].Products = append(x.Steps[i].Products, product)
			}
		}
	}

	return x, nil
}

// GetRecipe fetches a recipe from the database.
func (q *SQLQuerier) GetRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	return q.getRecipe(ctx, recipeID, "")
}

// GetRecipeByIDAndUser fetches a recipe from the database.
func (q *SQLQuerier) GetRecipeByIDAndUser(ctx context.Context, recipeID, userID string) (*types.Recipe, error) {
	return q.getRecipe(ctx, recipeID, userID)
}

const getTotalRecipesCountQuery = "SELECT COUNT(recipes.id) FROM recipes WHERE recipes.archived_on IS NULL"

// GetTotalRecipeCount fetches the count of recipes from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalRecipeCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	logger := q.logger.Clone()

	x = &types.RecipeList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "recipes", nil, nil, nil, "", recipesTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "recipes", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipes list retrieval query")
	}

	if x.Recipes, x.FilteredCount, x.TotalCount, err = q.scanRecipes(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipes")
	}

	return x, nil
}

// SearchForRecipes fetches a list of recipes from the database that match a query.
func (q *SQLQuerier) SearchForRecipes(ctx context.Context, recipeNameQuery string, filter *types.QueryFilter) (x *types.RecipeList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.RecipeList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	where := squirrel.ILike{"name": wrapQueryForILIKE(recipeNameQuery)}
	query, args := q.buildListQueryWithILike(ctx, "recipes", nil, nil, where, "", recipesTableColumns, "", false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "recipes", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipes search query")
	}

	if x.Recipes, x.FilteredCount, x.TotalCount, err = q.scanRecipes(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipes")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetRecipesWithIDsQuery(ctx context.Context, userID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"recipes.id":          ids,
		"recipes.archived_on": nil,
	}

	if userID != "" {
		withIDsWhere["recipes.created_by_user"] = userID
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
func (q *SQLQuerier) GetRecipesWithIDs(ctx context.Context, userID string, limit uint8, ids []string) ([]*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	query, args := q.buildGetRecipesWithIDsQuery(ctx, userID, limit, ids)

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

const recipeCreationQuery = "INSERT INTO recipes (id,name,source,description,inspired_by_recipe_id,created_by_user) VALUES ($1,$2,$3,$4,$5,$6)"

// CreateRecipe creates a recipe in the database.
func (q *SQLQuerier) CreateRecipe(ctx context.Context, input *types.RecipeDatabaseCreationInput) (*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeIDKey, input.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	args := []interface{}{
		input.ID,
		input.Name,
		input.Source,
		input.Description,
		input.InspiredByRecipeID,
		input.CreatedByUser,
	}

	// create the recipe.
	if err = q.performWriteQuery(ctx, q.db, "recipe creation", recipeCreationQuery, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "performing recipe creation query")
	}

	x := &types.Recipe{
		ID:                 input.ID,
		Name:               input.Name,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		CreatedByUser:      input.CreatedByUser,
		CreatedOn:          q.currentTime(),
	}

	for i, stepInput := range input.Steps {
		stepInput.Index = uint32(i)
		stepInput.BelongsToRecipe = x.ID

		// we need to go through all the prior steps and see
		// if the names of a product matches any ingredients
		// used in this step and not used in prior steps
		findCreatedRecipeStepProducts(input, i)

		s, createErr := q.createRecipeStep(ctx, tx, stepInput)
		if createErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step")
		}

		x.Steps = append(x.Steps, s)
	}

	if input.AlsoCreateMeal {
		_, mealCreateErr := q.createMeal(ctx, tx, &types.MealDatabaseCreationInput{
			ID:            ksuid.New().String(),
			Name:          x.Name,
			Description:   x.Description,
			CreatedByUser: x.CreatedByUser,
			Recipes:       []string{x.ID},
		})

		if mealCreateErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(mealCreateErr, logger, span, "creating meal from recipe")
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachRecipeIDToSpan(span, x.ID)
	logger.Info("recipe created")

	return x, nil
}

func findCreatedRecipeStepProducts(recipe *types.RecipeDatabaseCreationInput, stepIndex int) {
	step := recipe.Steps[stepIndex]

	priorSteps := []*types.RecipeStepDatabaseCreationInput{}
	for i, s := range recipe.Steps {
		if i < stepIndex {
			priorSteps = append(priorSteps, s)
		} else {
			break
		}
	}

	// created products is everything available to the step at the provided stepIndex.
	createdProducts := map[string]*types.RecipeStepProductDatabaseCreationInput{}
	for _, s := range priorSteps {
		for _, product := range s.Products {
			createdProducts[product.Name] = product
		}

		for _, ingredient := range s.Ingredients {
			if ingredient.ProductOfRecipeStep && ingredient.IngredientID == nil {
				delete(createdProducts, ingredient.Name)
			}
		}
	}

	for _, ingredient := range step.Ingredients {
		createdProduct, availableAsACreatedProduct := createdProducts[ingredient.Name]

		if ingredient.ProductOfRecipeStep && availableAsACreatedProduct {
			ingredient.RecipeStepProductID = &createdProduct.ID
		}
	}
}

const updateRecipeQuery = "UPDATE recipes SET name = $1, source = $2, description = $3, inspired_by_recipe_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND created_by_user = $5 AND id = $6"

// UpdateRecipe updates a particular recipe.
func (q *SQLQuerier) UpdateRecipe(ctx context.Context, updated *types.Recipe) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeIDKey, updated.ID)
	tracing.AttachRecipeIDToSpan(span, updated.ID)
	tracing.AttachUserIDToSpan(span, updated.CreatedByUser)

	args := []interface{}{
		updated.Name,
		updated.Source,
		updated.Description,
		updated.InspiredByRecipeID,
		updated.CreatedByUser,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe update", updateRecipeQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe")
	}

	logger.Info("recipe updated")

	return nil
}

const archiveRecipeQuery = "UPDATE recipes SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND created_by_user = $1 AND id = $2"

// ArchiveRecipe archives a recipe from the database by its ID.
func (q *SQLQuerier) ArchiveRecipe(ctx context.Context, recipeID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	args := []interface{}{
		userID,
		recipeID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe archive", archiveRecipeQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe")
	}

	logger.Info("recipe archived")

	return nil
}
