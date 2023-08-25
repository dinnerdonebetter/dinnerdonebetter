package postgres

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.MealDataManager = (*Querier)(nil)
)

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

	result, err := q.generatedQuerier.CheckMealExistence(ctx, q.db, mealID)
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
		r, getRecipeErr := q.getRecipe(ctx, mealComponent.RecipeID)
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

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.Meal]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetMeals(ctx, q.db, &generated.GetMealsParams{
		CreatedAfter:  sql.NullTime{},
		CreatedBefore: sql.NullTime{},
		UpdatedAfter:  sql.NullTime{},
		UpdatedBefore: sql.NullTime{},
		QueryOffset:   sql.NullInt32{},
		QueryLimit:    sql.NullInt32{},
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meals list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.Meal{
			CreatedAt:                result.CreatedAt,
			ArchivedAt:               timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:            timePointerFromNullTime(result.LastUpdatedAt),
			MaximumEstimatedPortions: float32PointerFromNullString(result.MaxEstimatedPortions),
			ID:                       result.ID,
			Description:              result.Description,
			CreatedByUser:            result.CreatedByUser,
			Name:                     result.Name,
			Components:               nil,
			MinimumEstimatedPortions: float32FromString(result.MinEstimatedPortions),
			EligibleForMealPlans:     result.EligibleForMealPlans,
		})

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetMealsWithIDs fetches a list of meals from the database that have IDs within a given set.
func (q *Querier) GetMealsWithIDs(ctx context.Context, ids []string) ([]*types.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	meals := []*types.Meal{}
	for _, id := range ids {
		r, err := q.GetMeal(ctx, id)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "getting meal")
		}

		meals = append(meals, r)
	}

	return meals, nil
}

// GetMealIDsThatNeedSearchIndexing fetches a list of meal IDs from the database that meet a particular filter.
func (q *Querier) GetMealIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetMealsNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing meals list retrieval query")
	}

	return results, nil
}

// SearchForMeals fetches a list of recipes from the database that match a query.
func (q *Querier) SearchForMeals(ctx context.Context, mealNameQuery string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Meal], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.Meal]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.SearchForMeals(ctx, q.db, &generated.SearchForMealsParams{
		Query:         nullStringFromString(mealNameQuery),
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meals search query")
	}

	var meal *types.Meal
	for _, result := range results {
		if meal != nil && meal.ID != result.ID {
			x.Data = append(x.Data, meal)
			meal = nil
		}

		if meal == nil {
			meal = &types.Meal{
				CreatedAt:                result.CreatedAt,
				ArchivedAt:               timePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt:            timePointerFromNullTime(result.LastUpdatedAt),
				MaximumEstimatedPortions: float32PointerFromNullString(result.MaxEstimatedPortions),
				ID:                       result.ID,
				Description:              result.Description,
				CreatedByUser:            result.CreatedByUser,
				Name:                     result.Name,
				Components:               []*types.MealComponent{},
				MinimumEstimatedPortions: float32FromString(result.MinEstimatedPortions),
				EligibleForMealPlans:     result.EligibleForMealPlans,
			}
		}

		meal.Components = append(meal.Components, &types.MealComponent{
			ComponentType: string(result.ComponentMealComponentType),
			Recipe: types.Recipe{
				ID: result.ComponentRecipeID,
			},
			RecipeScale: float32FromString(result.ComponentRecipeScale),
		})

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateMeal creates a meal in the database.
func (q *Querier) createMeal(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealDatabaseCreationInput) (*types.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealIDKey, input.ID).WithValue("meal.name", input.Name)

	// create the meal.
	if err := q.generatedQuerier.CreateMeal(ctx, querier, &generated.CreateMealParams{
		ID:                   input.ID,
		Name:                 input.Name,
		Description:          input.Description,
		MinEstimatedPortions: stringFromFloat32(input.MinimumEstimatedPortions),
		CreatedByUser:        input.CreatedByUser,
		MaxEstimatedPortions: nullStringFromFloat32Pointer(input.MaximumEstimatedPortions),
		EligibleForMealPlans: input.EligibleForMealPlans,
	}); err != nil {
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

// CreateMealComponent creates a meal component in the database.
func (q *Querier) CreateMealComponent(ctx context.Context, querier database.SQLQueryExecutor, mealID string, input *types.MealComponentDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return ErrNilInputProvided
	}

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	// create the meal.
	if err := q.generatedQuerier.CreateMealComponent(ctx, querier, &generated.CreateMealComponentParams{
		ID:                identifiers.New(),
		MealID:            mealID,
		RecipeID:          input.RecipeID,
		MealComponentType: generated.ComponentType(input.ComponentType),
		RecipeScale:       stringFromFloat32(input.RecipeScale),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "performing meal creation query")
	}

	return nil
}

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

	if err := q.generatedQuerier.UpdateMealLastIndexedAt(ctx, q.db, mealID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking meal as indexed")
	}

	logger.Info("meal marked as indexed")

	return nil
}

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

	if err := q.generatedQuerier.ArchiveMeal(ctx, q.db, &generated.ArchiveMealParams{
		CreatedByUser: userID,
		ID:            mealID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal")
	}

	logger.Info("meal archived")

	return nil
}
