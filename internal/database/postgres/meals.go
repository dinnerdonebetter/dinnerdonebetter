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
	_ types.MealDataManager = (*Querier)(nil)

	// mealsTableColumns are the columns for the meals table.
	mealsTableColumns = []string{
		"meals.id",
		"meals.name",
		"meals.description",
		"meals.created_at",
		"meals.last_updated_at",
		"meals.archived_at",
		"meals.created_by_user",
	}
)

// scanMeal takes a database Scanner (i.e. *sql.Row) and scans the result into a meal struct.
func (q *Querier) scanMeal(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.Meal, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.Meal{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Description,
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

// scanMeals takes some database rows and turns them into a slice of meals.
func (q *Querier) scanMeals(ctx context.Context, rows database.ResultIterator, includeCounts bool) (meals []*types.Meal, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanMeal(ctx, rows, includeCounts)
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

		meals = append(meals, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return meals, filteredCount, totalCount, nil
}

// scanMealWithRecipes takes a database Scanner (i.e. *sql.Row) and scans the result into a meal struct.
func (q *Querier) scanMealWithRecipes(ctx context.Context, rows database.ResultIterator) (x *types.Meal, recipeIDs []string, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.Meal{}

	for rows.Next() {
		var recipeID string
		targetVars := []interface{}{
			&x.ID,
			&x.Name,
			&x.Description,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
			&x.CreatedByUser,
			&recipeID,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, nil, observability.PrepareError(err, logger, span, "scanning complete meal")
		}
		recipeIDs = append(recipeIDs, recipeID)
	}

	return x, recipeIDs, nil
}

const mealExistenceQuery = "SELECT EXISTS ( SELECT meals.id FROM meals WHERE meals.archived_at IS NULL AND meals.id = $1 )"

// MealExists fetches whether a meal exists from the database.
func (q *Querier) MealExists(ctx context.Context, mealID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	args := []interface{}{
		mealID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing meal existence check")
	}

	return result, nil
}

const getMealByIDQuery = `SELECT 
	meals.id,
	meals.name,
	meals.description,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	meal_recipes.recipe_id
FROM meals
	FULL OUTER JOIN meal_recipes ON meal_recipes.meal_id=meals.id
WHERE meals.archived_at IS NULL
	AND meal_recipes.archived_at IS NULL
	AND meals.id = $1
`

// GetMeal fetches a meal from the database.
func (q *Querier) GetMeal(ctx context.Context, mealID string) (*types.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	args := []interface{}{
		mealID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "meal", getMealByIDQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal retrieval query")
	}

	m, recipeIDs, err := q.scanMealWithRecipes(ctx, rows)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal retrieval query")
	}

	if m == nil {
		return nil, sql.ErrNoRows
	}

	for _, id := range recipeIDs {
		r, getRecipeErr := q.GetRecipe(ctx, id)
		if getRecipeErr != nil {
			return nil, observability.PrepareError(getRecipeErr, logger, span, "fetching recipe for meal")
		}

		m.Recipes = append(m.Recipes, r)
	}

	return m, nil
}

// GetMeals fetches a list of meals from the database that meet a particular filter.
func (q *Querier) GetMeals(ctx context.Context, filter *types.QueryFilter) (x *types.MealList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.MealList{}
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

	query, args := q.buildListQuery(ctx, "meals", nil, nil, nil, "", mealsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "meals", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meals list retrieval query")
	}

	if x.Meals, x.FilteredCount, x.TotalCount, err = q.scanMeals(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meals")
	}

	return x, nil
}

// SearchForMeals fetches a list of recipes from the database that match a query.
func (q *Querier) SearchForMeals(ctx context.Context, mealNameQuery string, filter *types.QueryFilter) (x *types.MealList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.MealList{}
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

	where := squirrel.ILike{"name": wrapQueryForILIKE(mealNameQuery)}
	query, args := q.buildListQueryWithILike(ctx, "meals", nil, nil, where, "", mealsTableColumns, "", false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "meals", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meals search query")
	}

	if x.Meals, x.FilteredCount, x.TotalCount, err = q.scanMeals(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meals")
	}

	return x, nil
}

const mealCreationQuery = "INSERT INTO meals (id,name,description,created_by_user) VALUES ($1,$2,$3,$4)"

// CreateMeal creates a meal in the database.
func (q *Querier) createMeal(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealDatabaseCreationInput) (*types.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealIDKey, input.ID).WithValue("meal.name", input.Name)

	args := []interface{}{
		input.ID,
		input.Name,
		input.Description,
		input.CreatedByUser,
	}

	// create the meal.
	if err := q.performWriteQuery(ctx, querier, "meal creation", mealCreationQuery, args); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareError(err, logger, span, "performing meal creation query")
	}

	x := &types.Meal{
		ID:            input.ID,
		Name:          input.Name,
		Description:   input.Description,
		CreatedByUser: input.CreatedByUser,
		CreatedAt:     q.currentTime(),
	}

	for _, recipeID := range input.Recipes {
		if err := q.CreateMealRecipe(ctx, querier, x.ID, recipeID); err != nil {
			q.rollbackTransaction(ctx, querier)
			return nil, observability.PrepareError(err, logger, span, "creating meal recipe")
		}
	}

	tracing.AttachMealIDToSpan(span, x.ID)
	logger.Info("meal created")

	return x, nil
}

// CreateMeal creates a meal in the database.
func (q *Querier) CreateMeal(ctx context.Context, input *types.MealDatabaseCreationInput) (*types.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealIDKey, input.ID).WithValue("meal.name", input.Name)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	x, err := q.createMeal(ctx, tx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating meal")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	return x, nil
}

const mealRecipeCreationQuery = "INSERT INTO meal_recipes (id,meal_id,recipe_id) VALUES ($1,$2,$3)"

// CreateMealRecipe creates a meal in the database.
func (q *Querier) CreateMealRecipe(ctx context.Context, querier database.SQLQueryExecutor, mealID, recipeID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachUserIDToSpan(span, recipeID)

	args := []interface{}{
		ksuid.New().String(),
		mealID,
		recipeID,
	}

	// create the meal.
	if err := q.performWriteQuery(ctx, querier, "meal recipe creation", mealRecipeCreationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "performing meal creation query")
	}

	return nil
}

const archiveMealQuery = "UPDATE meals SET archived_at = extract(epoch FROM NOW()) WHERE archived_at IS NULL AND created_by_user = $1 AND id = $2"

// ArchiveMeal archives a meal from the database by its ID.
func (q *Querier) ArchiveMeal(ctx context.Context, mealID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	args := []interface{}{
		userID,
		mealID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal archive", archiveMealQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal")
	}

	logger.Info("meal archived")

	return nil
}
