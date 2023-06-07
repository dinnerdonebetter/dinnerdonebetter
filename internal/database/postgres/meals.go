package postgres

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/Masterminds/squirrel"
)

var (
	_ types.MealDataManager = (*Querier)(nil)

	// mealsTableColumns are the columns for the meals table.
	mealsTableColumns = []string{
		"meals.id",
		"meals.name",
		"meals.description",
		"meals.min_estimated_portions",
		"meals.max_estimated_portions",
		"meals.eligible_for_meal_plans",
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

	x = &types.Meal{}

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Description,
		&x.MinimumEstimatedPortions,
		&x.MaximumEstimatedPortions,
		&x.EligibleForMealPlans,
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

// scanMeals takes some database rows and turns them into a slice of meals.
func (q *Querier) scanMeals(ctx context.Context, rows database.ResultIterator, includeCounts bool) (meals []*types.Meal, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

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
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return meals, filteredCount, totalCount, nil
}

// scanMealWithRecipes takes a database Scanner (i.e. *sql.Row) and scans the result into a meal struct.
func (q *Querier) scanMealWithRecipes(ctx context.Context, rows database.ResultIterator) (x *types.Meal, mealComponents []*types.MealComponentDatabaseCreationInput, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.Meal{}

	for rows.Next() {
		var (
			recipeID, componentType string
			recipeScale             float32
		)
		targetVars := []any{
			&x.ID,
			&x.Name,
			&x.Description,
			&x.MinimumEstimatedPortions,
			&x.MaximumEstimatedPortions,
			&x.EligibleForMealPlans,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
			&x.CreatedByUser,
			&recipeID,
			&recipeScale,
			&componentType,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, nil, observability.PrepareError(err, span, "scanning complete meal")
		}
		mealComponents = append(mealComponents, &types.MealComponentDatabaseCreationInput{ComponentType: componentType, RecipeScale: recipeScale, RecipeID: recipeID})
	}

	if err = rows.Err(); err != nil {
		return nil, nil, observability.PrepareError(err, span, "querying for complete meal")
	}

	return x, mealComponents, nil
}

//go:embed queries/meals/exists.sql
var mealExistenceQuery string

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

	args := []any{
		mealID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal existence check")
	}

	return result, nil
}

//go:embed queries/meals/get_one.sql
var getMealByIDQuery string

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

	getMealByIDArgs := []any{
		mealID,
	}

	rows, err := q.getRows(ctx, q.db, "meal", getMealByIDQuery, getMealByIDArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal retrieval query")
	}

	m, mealComponents, err := q.scanMealWithRecipes(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meal retrieval query")
	}

	if m == nil || m.ID == "" || len(mealComponents) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, mealComponent := range mealComponents {
		r, getRecipeErr := q.getRecipe(ctx, mealComponent.RecipeID, "")
		if getRecipeErr != nil {
			return nil, observability.PrepareError(getRecipeErr, span, "fetching recipe for meal")
		}

		m.Components = append(m.Components, &types.MealComponent{
			ComponentType: mealComponent.ComponentType,
			RecipeScale:   mealComponent.RecipeScale,
			Recipe:        *r,
		})
	}

	return m, nil
}

// GetMeals fetches a list of meals from the database that meet a particular filter.
func (q *Querier) GetMeals(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Meal], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.Meal]{}
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

	query, args := q.buildListQuery(ctx, "meals", nil, nil, nil, "", mealsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "meals", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meals list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanMeals(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meals")
	}

	return x, nil
}

// SearchForMeals fetches a list of recipes from the database that match a query.
func (q *Querier) SearchForMeals(ctx context.Context, mealNameQuery string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Meal], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.Meal]{}
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

	rows, err := q.getRows(ctx, q.db, "meals search", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meals search query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanMeals(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meals")
	}

	return x, nil
}

//go:embed queries/meals/create.sql
var mealCreationQuery string

// CreateMeal creates a meal in the database.
func (q *Querier) createMeal(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealDatabaseCreationInput) (*types.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealIDKey, input.ID).WithValue("meal.name", input.Name)

	args := []any{
		input.ID,
		input.Name,
		input.Description,
		input.MinimumEstimatedPortions,
		input.MaximumEstimatedPortions,
		input.EligibleForMealPlans,
		input.CreatedByUser,
	}

	// create the meal.
	if err := q.performWriteQuery(ctx, querier, "meal creation", mealCreationQuery, args); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal creation query")
	}

	x := &types.Meal{
		ID:                       input.ID,
		Name:                     input.Name,
		Description:              input.Description,
		MinimumEstimatedPortions: input.MinimumEstimatedPortions,
		MaximumEstimatedPortions: input.MaximumEstimatedPortions,
		EligibleForMealPlans:     input.EligibleForMealPlans,
		CreatedByUser:            input.CreatedByUser,
		CreatedAt:                q.currentTime(),
	}

	for _, recipeID := range input.Components {
		if err := q.CreateMealComponent(ctx, querier, x.ID, recipeID); err != nil {
			q.rollbackTransaction(ctx, querier)
			return nil, observability.PrepareAndLogError(err, logger, span, "creating meal recipe")
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

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "beginning transaction")
	}

	x, err := q.createMeal(ctx, tx, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating meal")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, span, "committing transaction")
	}

	return x, nil
}

//go:embed queries/meal_components/create.sql
var mealRecipeCreationQuery string

// CreateMealComponent creates a meal component in the database.
func (q *Querier) CreateMealComponent(ctx context.Context, querier database.SQLQueryExecutor, mealID string, input *types.MealComponentDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	if input == nil {
		return ErrNilInputProvided
	}

	args := []any{
		identifiers.New(),
		mealID,
		input.RecipeID,
		input.ComponentType,
		input.RecipeScale,
	}

	// create the meal.
	if err := q.performWriteQuery(ctx, querier, "meal recipe creation", mealRecipeCreationQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "performing meal creation query")
	}

	return nil
}

//go:embed queries/meals/update_last_indexed_at.sql
var updateMealLastIndexedAtQuery string

// MarkMealAsIndexed updates a particular meal's last_indexed_at value.
func (q *Querier) MarkMealAsIndexed(ctx context.Context, mealID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	args := []any{
		mealID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal last_indexed_at", updateMealLastIndexedAtQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking meal as indexed")
	}

	logger.Info("meal marked as indexed")

	return nil
}

//go:embed queries/meals/archive.sql
var archiveMealQuery string

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

	args := []any{
		userID,
		mealID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal archive", archiveMealQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal")
	}

	logger.Info("meal archived")

	return nil
}
