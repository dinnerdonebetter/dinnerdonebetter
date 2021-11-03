package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"

	database "github.com/prixfixeco/api_server/internal/database"
	observability "github.com/prixfixeco/api_server/internal/observability"
	keys "github.com/prixfixeco/api_server/internal/observability/keys"
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

// scanFullRecipe takes a database Scanner (i.e. *sql.Row) and scans the result into a full recipe struct.
func (q *SQLQuerier) scanFullRecipe(ctx context.Context, scan database.Scanner) (*types.FullRecipe, *types.FullRecipeStep, *types.FullRecipeStepIngredient, error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	recipe := &types.FullRecipe{Steps: []*types.FullRecipeStep{}}
	recipeStep := &types.FullRecipeStep{Ingredients: []*types.FullRecipeStepIngredient{}}
	recipeStepIngredient := &types.FullRecipeStepIngredient{}

	targetVars := []interface{}{
		&recipe.ID,
		&recipe.Name,
		&recipe.Source,
		&recipe.Description,
		&recipe.InspiredByRecipeID,
		&recipe.CreatedOn,
		&recipe.LastUpdatedOn,
		&recipe.ArchivedOn,
		&recipe.CreatedByUser,
		&recipeStep.ID,
		&recipeStep.Index,
		&recipeStep.Preparation.ID,
		&recipeStep.Preparation.Name,
		&recipeStep.Preparation.Description,
		&recipeStep.Preparation.IconPath,
		&recipeStep.Preparation.CreatedOn,
		&recipeStep.Preparation.LastUpdatedOn,
		&recipeStep.Preparation.ArchivedOn,
		&recipeStep.PrerequisiteStep,
		&recipeStep.MinEstimatedTimeInSeconds,
		&recipeStep.MaxEstimatedTimeInSeconds,
		&recipeStep.TemperatureInCelsius,
		&recipeStep.Notes,
		&recipeStep.Why,
		&recipeStep.CreatedOn,
		&recipeStep.LastUpdatedOn,
		&recipeStep.ArchivedOn,
		&recipeStep.BelongsToRecipe,
		&recipeStepIngredient.ID,
		&recipeStepIngredient.Ingredient.ID,
		&recipeStepIngredient.Ingredient.Name,
		&recipeStepIngredient.Ingredient.Variant,
		&recipeStepIngredient.Ingredient.Description,
		&recipeStepIngredient.Ingredient.Warning,
		&recipeStepIngredient.Ingredient.ContainsEgg,
		&recipeStepIngredient.Ingredient.ContainsDairy,
		&recipeStepIngredient.Ingredient.ContainsPeanut,
		&recipeStepIngredient.Ingredient.ContainsTreeNut,
		&recipeStepIngredient.Ingredient.ContainsSoy,
		&recipeStepIngredient.Ingredient.ContainsWheat,
		&recipeStepIngredient.Ingredient.ContainsShellfish,
		&recipeStepIngredient.Ingredient.ContainsSesame,
		&recipeStepIngredient.Ingredient.ContainsFish,
		&recipeStepIngredient.Ingredient.ContainsGluten,
		&recipeStepIngredient.Ingredient.AnimalFlesh,
		&recipeStepIngredient.Ingredient.AnimalDerived,
		&recipeStepIngredient.Ingredient.Volumetric,
		&recipeStepIngredient.Ingredient.IconPath,
		&recipeStepIngredient.Ingredient.CreatedOn,
		&recipeStepIngredient.Ingredient.LastUpdatedOn,
		&recipeStepIngredient.Ingredient.ArchivedOn,
		&recipeStepIngredient.QuantityType,
		&recipeStepIngredient.QuantityValue,
		&recipeStepIngredient.QuantityNotes,
		&recipeStepIngredient.ProductOfRecipeStep,
		&recipeStepIngredient.IngredientNotes,
		&recipeStepIngredient.CreatedOn,
		&recipeStepIngredient.LastUpdatedOn,
		&recipeStepIngredient.ArchivedOn,
		&recipeStepIngredient.BelongsToRecipeStep,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, nil, nil, observability.PrepareError(err, q.logger, span, "")
	}

	return recipe, recipeStep, recipeStepIngredient, nil
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

const getRecipeQuery = "SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.created_by_user FROM recipes WHERE recipes.archived_on IS NULL AND recipes.id = $1"

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

var fullRecipeColumns = []string{
	"recipes.id",
	"recipes.name",
	"recipes.source",
	"recipes.description",
	"recipes.inspired_by_recipe_id",
	"recipes.created_on",
	"recipes.last_updated_on",
	"recipes.archived_on",
	"recipes.created_by_user",
	"recipe_steps.id",
	"recipe_steps.index",
	"valid_preparations.id",
	"valid_preparations.name",
	"valid_preparations.description",
	"valid_preparations.icon_path",
	"valid_preparations.created_on",
	"valid_preparations.last_updated_on",
	"valid_preparations.archived_on",
	"recipe_steps.prerequisite_step",
	"recipe_steps.min_estimated_time_in_seconds",
	"recipe_steps.max_estimated_time_in_seconds",
	"recipe_steps.temperature_in_celsius",
	"recipe_steps.notes",
	"recipe_steps.why",
	"recipe_steps.created_on",
	"recipe_steps.last_updated_on",
	"recipe_steps.archived_on",
	"recipe_steps.belongs_to_recipe",
	"recipe_step_ingredients.id",
	"valid_ingredients.id",
	"valid_ingredients.name",
	"valid_ingredients.variant",
	"valid_ingredients.description",
	"valid_ingredients.warning",
	"valid_ingredients.contains_egg",
	"valid_ingredients.contains_dairy",
	"valid_ingredients.contains_peanut",
	"valid_ingredients.contains_tree_nut",
	"valid_ingredients.contains_soy",
	"valid_ingredients.contains_wheat",
	"valid_ingredients.contains_shellfish",
	"valid_ingredients.contains_sesame",
	"valid_ingredients.contains_fish",
	"valid_ingredients.contains_gluten",
	"valid_ingredients.animal_flesh",
	"valid_ingredients.animal_derived",
	"valid_ingredients.volumetric",
	"valid_ingredients.icon_path",
	"valid_ingredients.created_on",
	"valid_ingredients.last_updated_on",
	"valid_ingredients.archived_on",
	"recipe_step_ingredients.quantity_type",
	"recipe_step_ingredients.quantity_value",
	"recipe_step_ingredients.quantity_notes",
	"recipe_step_ingredients.product_of_recipe_step",
	"recipe_step_ingredients.ingredient_notes",
	"recipe_step_ingredients.created_on",
	"recipe_step_ingredients.last_updated_on",
	"recipe_step_ingredients.archived_on",
	"recipe_step_ingredients.belongs_to_recipe_step",
}

const getFullRecipeQuery = `SELECT 
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
	recipe_steps.prerequisite_step,
	recipe_steps.min_estimated_time_in_seconds,
	recipe_steps.max_estimated_time_in_seconds,
	recipe_steps.temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.why,
	recipe_steps.created_on,
	recipe_steps.last_updated_on,
	recipe_steps.archived_on,
	recipe_steps.belongs_to_recipe,
	recipe_step_ingredients.id,
	valid_ingredients.id,
	valid_ingredients.name,
	valid_ingredients.variant,
	valid_ingredients.description,
	valid_ingredients.warning,
	valid_ingredients.contains_egg,
	valid_ingredients.contains_dairy,
	valid_ingredients.contains_peanut,
	valid_ingredients.contains_tree_nut,
	valid_ingredients.contains_soy,
	valid_ingredients.contains_wheat,
	valid_ingredients.contains_shellfish,
	valid_ingredients.contains_sesame,
	valid_ingredients.contains_fish,
	valid_ingredients.contains_gluten,
	valid_ingredients.animal_flesh,
	valid_ingredients.animal_derived,
	valid_ingredients.volumetric,
	valid_ingredients.icon_path,
	valid_ingredients.created_on,
	valid_ingredients.last_updated_on,
	valid_ingredients.archived_on,
	recipe_step_ingredients.quantity_type,
	recipe_step_ingredients.quantity_value,
	recipe_step_ingredients.quantity_notes,
	recipe_step_ingredients.product_of_recipe_step,
	recipe_step_ingredients.ingredient_notes,
	recipe_step_ingredients.created_on,
	recipe_step_ingredients.last_updated_on,
	recipe_step_ingredients.archived_on,
	recipe_step_ingredients.belongs_to_recipe_step
FROM recipe_step_ingredients
	FULL OUTER JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id
	FULL OUTER JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	FULL OUTER JOIN valid_ingredients ON recipe_step_ingredients.ingredient_id=valid_ingredients.id
	FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_step_ingredients.archived_on IS NULL
	AND recipe_steps.archived_on IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipes.archived_on IS NULL
	AND recipes.id = $2
`

// GetFullRecipe fetches a recipe from the database.
func (q *SQLQuerier) GetFullRecipe(ctx context.Context, recipeID string) (*types.FullRecipe, error) {
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
		recipeID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "recipe", getFullRecipeQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing full recipe retrieval query")
	}

	var (
		recipe           *types.FullRecipe
		currentStepIndex = 0
	)

	for rows.Next() {
		rowRecipe, rowRecipeStep, rowRecipeStepIngredient, scanErr := q.scanFullRecipe(ctx, rows)
		if scanErr != nil {
			return nil, scanErr
		}

		if recipe == nil {
			recipe = rowRecipe
		}

		if len(recipe.Steps) == 0 && currentStepIndex == 0 {
			recipe.Steps = append(recipe.Steps, rowRecipeStep)
		}

		if recipe.Steps[currentStepIndex].ID != rowRecipeStep.ID {
			currentStepIndex++
			recipe.Steps = append(recipe.Steps, rowRecipeStep)
		}

		recipe.Steps[currentStepIndex].Ingredients = append(recipe.Steps[currentStepIndex].Ingredients, rowRecipeStepIngredient)
	}

	if recipe == nil {
		return nil, sql.ErrNoRows
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
		"",
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

	for _, stepInput := range input.Steps {
		stepInput.BelongsToRecipe = x.ID
		s, createErr := q.createRecipeStep(ctx, tx, stepInput)
		if createErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step")
		}
		x.Steps = append(x.Steps, s)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachRecipeIDToSpan(span, x.ID)
	logger.Info("recipe created")

	return x, nil
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

	logger := q.logger

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
