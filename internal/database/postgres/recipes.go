package postgres

import (
	"context"
	"database/sql"

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
		"recipes.yields_portions",
		"recipes.seal_of_approval",
		"recipes.created_at",
		"recipes.last_updated_at",
		"recipes.archived_at",
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
		&x.YieldsPortions,
		&x.SealOfApproval,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
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

const recipeExistenceQuery = "SELECT EXISTS ( SELECT recipes.id FROM recipes WHERE recipes.archived_at IS NULL AND recipes.id = $1 )"

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
	recipes.yields_portions,
	recipes.seal_of_approval,
	recipes.created_at,
	recipes.last_updated_at,
	recipes.archived_at,
	recipes.created_by_user,
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.zero_ingredients_allowable,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.optional,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe
FROM recipes
	FULL OUTER JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
	FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_at IS NULL
	AND recipes.id = $1
ORDER BY recipe_steps.index
`

const getRecipeByIDAndAuthorIDQuery = `SELECT 
	recipes.id,
	recipes.name,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.yields_portions,
	recipes.seal_of_approval,
	recipes.created_at,
	recipes.last_updated_at,
	recipes.archived_at,
	recipes.created_by_user,
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.zero_ingredients_allowable,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.optional,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe
FROM recipes
	FULL OUTER JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
	FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_at IS NULL
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

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Source,
		&x.Description,
		&x.InspiredByRecipeID,
		&x.YieldsPortions,
		&x.SealOfApproval,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
		&x.CreatedByUser,
		&y.ID,
		&y.Index,
		&y.Preparation.ID,
		&y.Preparation.Name,
		&y.Preparation.Description,
		&y.Preparation.IconPath,
		&y.Preparation.YieldsNothing,
		&y.Preparation.RestrictToIngredients,
		&y.Preparation.ZeroIngredientsAllowable,
		&y.Preparation.PastTense,
		&y.Preparation.CreatedAt,
		&y.Preparation.LastUpdatedAt,
		&y.Preparation.ArchivedAt,
		&y.MinimumEstimatedTimeInSeconds,
		&y.MaximumEstimatedTimeInSeconds,
		&y.MinimumTemperatureInCelsius,
		&y.MaximumTemperatureInCelsius,
		&y.Notes,
		&y.ExplicitInstructions,
		&y.Optional,
		&y.CreatedAt,
		&y.LastUpdatedAt,
		&y.ArchivedAt,
		&y.BelongsToRecipe,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, nil, 0, 0, observability.PrepareError(err, q.logger, span, "")
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

	// need to grab instruments here and add them to steps
	instruments, err := q.getRecipeStepInstrumentsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe step instruments for recipe")
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

		for _, instrument := range instruments {
			if instrument.BelongsToRecipeStep == step.ID {
				x.Steps[i].Instruments = append(x.Steps[i].Instruments, instrument)
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

// GetRecipes fetches a list of recipes from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipes(ctx context.Context, filter *types.QueryFilter) (x *types.RecipeList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.RecipeList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
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
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
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

const recipeCreationQuery = "INSERT INTO recipes (id,name,source,description,inspired_by_recipe_id,yields_portions,seal_of_approval,created_by_user) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"

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
		input.YieldsPortions,
		input.SealOfApproval,
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
		YieldsPortions:     input.YieldsPortions,
		SealOfApproval:     input.SealOfApproval,
		CreatedAt:          q.currentTime(),
	}

	for i, stepInput := range input.Steps {
		stepInput.Index = uint32(i)
		stepInput.BelongsToRecipe = x.ID

		// we need to go through all the prior steps and see
		// if the names of a product matches any ingredients
		// used in this step and not used in prior steps
		findCreatedRecipeStepProductsForIngredients(input, i)
		findCreatedRecipeStepProductsForInstruments(input, i)

		s, createErr := q.createRecipeStep(ctx, tx, stepInput)
		if createErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createErr, logger, span, "creating recipe step #%d", i+1)
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

func findCreatedRecipeStepProductsForIngredients(recipe *types.RecipeDatabaseCreationInput, stepIndex int) {
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
			if product.Type == types.RecipeStepProductIngredientType {
				createdProducts[product.Name] = product
			}
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

func findCreatedRecipeStepProductsForInstruments(recipe *types.RecipeDatabaseCreationInput, stepIndex int) {
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
			if product.Type == types.RecipeStepProductInstrumentType {
				createdProducts[product.Name] = product
			}
		}

		for _, instrument := range s.Instruments {
			if instrument.RecipeStepProductID == nil {
				delete(createdProducts, instrument.Name)
			}
		}
	}

	for _, instrument := range step.Instruments {
		createdProduct, availableAsACreatedProduct := createdProducts[instrument.Name]

		if instrument.ProductOfRecipeStep && availableAsACreatedProduct {
			instrument.RecipeStepProductID = &createdProduct.ID
		}
	}
}

const updateRecipeQuery = `UPDATE recipes SET name = $1, source = $2, description = $3, inspired_by_recipe_id = $4, yields_portions = $5, seal_of_approval = $6, last_updated_at = extract(epoch FROM NOW()) WHERE archived_at IS NULL AND created_by_user = $7 AND id = $8`

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
		updated.YieldsPortions,
		updated.SealOfApproval,
		updated.CreatedByUser,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe update", updateRecipeQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe")
	}

	logger.Info("recipe updated")

	return nil
}

const archiveRecipeQuery = "UPDATE recipes SET archived_at = extract(epoch FROM NOW()) WHERE archived_at IS NULL AND created_by_user = $1 AND id = $2"

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
