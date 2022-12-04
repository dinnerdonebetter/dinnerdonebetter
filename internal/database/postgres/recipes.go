package postgres

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/Masterminds/squirrel"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

var (
	_ types.RecipeDataManager = (*Querier)(nil)

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
func (q *Querier) scanRecipe(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.Recipe, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.Recipe{}

	targetVars := []any{
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
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipes takes some database rows and turns them into a slice of recipes.
func (q *Querier) scanRecipes(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipes []*types.Recipe, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

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
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return recipes, filteredCount, totalCount, nil
}

//go:embed queries/recipes/exists.sql
var recipeExistenceQuery string

// RecipeExists fetches whether a recipe exists from the database.
func (q *Querier) RecipeExists(ctx context.Context, recipeID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	args := []any{
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, span, "performing recipe existence check")
	}

	return result, nil
}

//go:embed queries/recipes/get_by_id.sql
var getRecipeByIDQuery string

//go:embed queries/recipes/get_by_id_and_author_id.sql
var getRecipeByIDAndAuthorIDQuery string

// scanRecipeAndStep takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe struct.
func (q *Querier) scanRecipeAndStep(ctx context.Context, scan database.Scanner) (x *types.Recipe, y *types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.Recipe{}
	y = &types.RecipeStep{}

	targetVars := []any{
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
		&y.Preparation.MinimumIngredientCount,
		&y.Preparation.MaximumIngredientCount,
		&y.Preparation.MinimumInstrumentCount,
		&y.Preparation.MaximumInstrumentCount,
		&y.Preparation.TemperatureRequired,
		&y.Preparation.TimeEstimateRequired,
		&y.Preparation.ConditionExpressionRequired,
		&y.Preparation.Slug,
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
		&y.ConditionExpression,
		&y.Optional,
		&y.CreatedAt,
		&y.LastUpdatedAt,
		&y.ArchivedAt,
		&y.BelongsToRecipe,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, nil, 0, 0, observability.PrepareError(err, span, "")
	}

	return x, y, filteredCount, totalCount, nil
}

// getRecipe fetches a recipe from the database.
func (q *Querier) getRecipe(ctx context.Context, recipeID, userID string) (*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	args := []any{
		recipeID,
	}

	query := getRecipeByIDQuery
	if userID != "" {
		query = getRecipeByIDAndAuthorIDQuery
		args = append(args, userID)
	}

	rows, err := q.getRows(ctx, q.db, "get recipe", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning recipe")
	}

	var x *types.Recipe
	for rows.Next() {
		recipe, recipeStep, _, _, recipeScanErr := q.scanRecipeAndStep(ctx, rows)
		if recipeScanErr != nil {
			return nil, observability.PrepareError(recipeScanErr, span, "scanning recipe")
		}

		if x == nil {
			x = recipe
		}

		x.Steps = append(x.Steps, recipeStep)
	}

	if x == nil {
		return nil, sql.ErrNoRows
	}

	prepTasks, err := q.getRecipePrepTasksForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step ingredients for recipe")
	}
	x.PrepTasks = prepTasks

	recipeMedia, err := q.getRecipeMediaForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step ingredients for recipe")
	}
	x.Media = recipeMedia

	ingredients, err := q.getRecipeStepIngredientsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step ingredients for recipe")
	}

	products, err := q.getRecipeStepProductsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step products for recipe")
	}

	instruments, err := q.getRecipeStepInstrumentsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step instruments for recipe")
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

		recipeMedia, err = q.getRecipeMediaForRecipeStep(ctx, recipeID, step.ID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "fetching recipe step ingredients for recipe")
		}
		x.Steps[i].Media = recipeMedia
	}

	return x, nil
}

// GetRecipe fetches a recipe from the database.
func (q *Querier) GetRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	return q.getRecipe(ctx, recipeID, "")
}

// GetRecipeByIDAndUser fetches a recipe from the database.
func (q *Querier) GetRecipeByIDAndUser(ctx context.Context, recipeID, userID string) (*types.Recipe, error) {
	return q.getRecipe(ctx, recipeID, userID)
}

// GetRecipes fetches a list of recipes from the database that meet a particular filter.
func (q *Querier) GetRecipes(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Recipe], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.QueryFilteredResult[types.Recipe]{}
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "recipes", nil, nil, nil, "", recipesTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "recipes", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing recipes list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanRecipes(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, span, "scanning recipes")
	}

	return x, nil
}

//go:embed queries/recipes/ids_for_meal.sql
var getRecipesForMealQuery string

// getRecipeIDsForMeal fetches a list of recipe IDs from the database that are associated with a given meal.
func (q *Querier) getRecipeIDsForMeal(ctx context.Context, mealID string) (x []string, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachMealIDToSpan(span, mealID)

	args := []any{
		mealID,
	}

	rows, err := q.getRows(ctx, q.db, "recipes for meal", getRecipesForMealQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing recipes list retrieval query")
	}

	if x, err = q.scanIDs(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "scanning recipes")
	}

	return x, nil
}

// SearchForRecipes fetches a list of recipes from the database that match a query.
func (q *Querier) SearchForRecipes(ctx context.Context, recipeNameQuery string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Recipe], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.QueryFilteredResult[types.Recipe]{}
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

	rows, err := q.getRows(ctx, q.db, "recipes search", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing recipes search query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanRecipes(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, span, "scanning recipes")
	}

	return x, nil
}

//go:embed queries/recipes/create.sql
var recipeCreationQuery string

// CreateRecipe creates a recipe in the database.
func (q *Querier) CreateRecipe(ctx context.Context, input *types.RecipeDatabaseCreationInput) (*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeIDKey, input.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	args := []any{
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
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe creation query")
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
			return nil, observability.PrepareError(createErr, span, "creating recipe step #%d", i+1)
		}

		x.Steps = append(x.Steps, s)
	}

	for i, prepTaskInput := range input.PrepTasks {
		pt, createPrepTaskErr := q.createRecipePrepTask(ctx, tx, prepTaskInput)
		if createPrepTaskErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createPrepTaskErr, span, "creating recipe prep task #%d", i+1)
		}

		x.PrepTasks = append(x.PrepTasks, pt)
	}

	if input.AlsoCreateMeal {
		_, mealCreateErr := q.createMeal(ctx, tx, &types.MealDatabaseCreationInput{
			ID:            identifiers.New(),
			Name:          x.Name,
			Description:   x.Description,
			CreatedByUser: x.CreatedByUser,
			Components: []*types.MealComponentDatabaseCreationInput{
				{
					RecipeID:      x.ID,
					ComponentType: types.MealComponentTypesMain,
				},
			},
		})

		if mealCreateErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(mealCreateErr, span, "creating meal from recipe")
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
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

//go:embed queries/recipes/update.sql
var updateRecipeQuery string

// UpdateRecipe updates a particular recipe.
func (q *Querier) UpdateRecipe(ctx context.Context, updated *types.Recipe) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeIDKey, updated.ID)
	tracing.AttachRecipeIDToSpan(span, updated.ID)
	tracing.AttachUserIDToSpan(span, updated.CreatedByUser)

	args := []any{
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
		return observability.PrepareAndLogError(err, logger, span, "updating recipe")
	}

	logger.Info("recipe updated")

	return nil
}

//go:embed queries/recipes/archive.sql
var archiveRecipeQuery string

// ArchiveRecipe archives a recipe from the database by its ID.
func (q *Querier) ArchiveRecipe(ctx context.Context, recipeID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	args := []any{
		userID,
		recipeID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe archive", archiveRecipeQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe")
	}

	logger.Info("recipe archived")

	return nil
}
