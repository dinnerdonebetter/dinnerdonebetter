package mealplanning

import (
	"context"
	"database/sql"
	"errors"

	identitykeys "github.com/dinnerdonebetter/backend/internal/domain/identity/keys"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mealplanningkeys "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ mealplanning.MealDataManager = (*repository)(nil)
)

// MealExists fetches whether a meal exists from the database.
func (q *repository) MealExists(ctx context.Context, mealID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(mealplanningkeys.MealIDKey, mealID)
	tracing.AttachToSpan(span, mealplanningkeys.MealIDKey, mealID)

	result, err := q.generatedQuerier.CheckMealExistence(ctx, q.readDB, mealID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal existence check")
	}

	return result, nil
}

// GetMeal fetches a meal from the database.
func (q *repository) GetMeal(ctx context.Context, mealID string) (*mealplanning.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(mealplanningkeys.MealIDKey, mealID)
	tracing.AttachToSpan(span, mealplanningkeys.MealIDKey, mealID)

	results, err := q.generatedQuerier.GetMeal(ctx, q.readDB, mealID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal retrieval query")
	}

	var meal *mealplanning.Meal
	for _, result := range results {
		if meal == nil {
			meal = &mealplanning.Meal{
				CreatedAt:     result.CreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
				ID:            result.ID,
				Description:   result.Description,
				CreatedByUser: result.CreatedByUser,
				Name:          result.Name,
				Components:    nil,
				EstimatedPortions: types.Float32RangeWithOptionalMax{
					Min: database.Float32FromString(result.MinEstimatedPortions),
					Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
				},
				EligibleForMealPlans: result.EligibleForMealPlans,
			}
		}

		recipe, recipeErr := q.getRecipe(ctx, result.ComponentRecipeID)
		if recipeErr != nil {
			return nil, observability.PrepareAndLogError(recipeErr, logger, span, "getting recipe")
		}

		meal.Components = append(meal.Components, &mealplanning.MealComponent{
			ComponentType: string(result.ComponentMealComponentType),
			Recipe:        *recipe,
			RecipeScale:   database.Float32FromString(result.ComponentRecipeScale),
		})
	}

	if meal == nil || meal.ID == "" || len(meal.Components) == 0 {
		return nil, sql.ErrNoRows
	}

	for i, mealComponent := range meal.Components {
		var r *mealplanning.Recipe
		r, err = q.getRecipe(ctx, mealComponent.Recipe.ID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "fetching recipe for meal")
		}

		meal.Components[i].Recipe = *r
	}

	return meal, nil
}

// GetMeals fetches a list of meals from the database that meet a particular filter.
func (q *repository) GetMeals(ctx context.Context, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.Meal], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	var (
		data          []*mealplanning.Meal
		filteredCount uint64
		totalCount    uint64
	)

	results, err := q.generatedQuerier.GetMeals(ctx, q.readDB, &generated.GetMealsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:          database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:     database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meals list retrieval query")
	}

	var meal *mealplanning.Meal
	for _, result := range results {
		if meal != nil && meal.ID != result.ID {
			data = append(data, meal)
			meal = nil
		}

		if meal == nil {
			meal = &mealplanning.Meal{
				CreatedAt:     result.CreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
				ID:            result.ID,
				Description:   result.Description,
				CreatedByUser: result.CreatedByUser,
				Name:          result.Name,
				Components:    []*mealplanning.MealComponent{},
				EstimatedPortions: types.Float32RangeWithOptionalMax{
					Min: database.Float32FromString(result.MinEstimatedPortions),
					Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
				},
				EligibleForMealPlans: result.EligibleForMealPlans,
			}
		}

		if result.ComponentRecipeID.Valid {
			recipe, recipeErr := q.getRecipe(ctx, result.ComponentRecipeID.String)
			if recipeErr != nil {
				if errors.Is(recipeErr, sql.ErrNoRows) {
					// Recipe missing or archived (e.g. orphaned reference from another test).
					// Skip this component so listing succeeds; avoids cross-test pollution.
					logger.WithValue(mealplanningkeys.MealIDKey, result.ID).
						WithValue(mealplanningkeys.RecipeIDKey, result.ComponentRecipeID.String).
						Info("skipping meal component with missing or archived recipe")
					continue
				}
				return nil, observability.PrepareAndLogError(recipeErr, logger, span, "getting recipe for meal component")
			}

			componentType := ""
			if result.ComponentMealComponentType.Valid {
				componentType = string(result.ComponentMealComponentType.ComponentType)
			}
			recipeScale := float32(0)
			if result.ComponentRecipeScale.Valid {
				recipeScale = database.Float32FromString(result.ComponentRecipeScale.String)
			}

			meal.Components = append(meal.Components, &mealplanning.MealComponent{
				ComponentType: componentType,
				Recipe:        *recipe,
				RecipeScale:   recipeScale,
			})
		}

		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	if meal != nil {
		data = append(data, meal)
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(m *mealplanning.Meal) string { return m.ID },
		filter,
	)

	return x, nil
}

// GetMealsCreatedByUser fetches a list of meals from the database that meet a particular filter.
func (q *repository) GetMealsCreatedByUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.Meal], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	var (
		data          []*mealplanning.Meal
		filteredCount uint64
		totalCount    uint64
	)

	results, err := q.generatedQuerier.GetMealsCreatedByUser(ctx, q.readDB, &generated.GetMealsCreatedByUserParams{
		CreatedByUser:   userID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:          database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:     database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meals list retrieval query")
	}

	var meal *mealplanning.Meal
	for _, result := range results {
		if meal != nil && meal.ID != result.ID {
			data = append(data, meal)
			meal = nil
		}

		if meal == nil {
			meal = &mealplanning.Meal{
				CreatedAt:     result.CreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
				ID:            result.ID,
				Description:   result.Description,
				CreatedByUser: result.CreatedByUser,
				Name:          result.Name,
				Components:    []*mealplanning.MealComponent{},
				EstimatedPortions: types.Float32RangeWithOptionalMax{
					Min: database.Float32FromString(result.MinEstimatedPortions),
					Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
				},
				EligibleForMealPlans: result.EligibleForMealPlans,
			}
		}

		if result.ComponentRecipeID.Valid {
			recipe, recipeErr := q.getRecipe(ctx, result.ComponentRecipeID.String)
			if recipeErr != nil {
				if errors.Is(recipeErr, sql.ErrNoRows) {
					logger.WithValue(mealplanningkeys.MealIDKey, result.ID).
						WithValue(mealplanningkeys.RecipeIDKey, result.ComponentRecipeID.String).
						Info("skipping meal component with missing or archived recipe")
					continue
				}
				return nil, observability.PrepareAndLogError(recipeErr, logger, span, "getting recipe for meal component")
			}

			componentType := ""
			if result.ComponentMealComponentType.Valid {
				componentType = string(result.ComponentMealComponentType.ComponentType)
			}
			recipeScale := float32(0)
			if result.ComponentRecipeScale.Valid {
				recipeScale = database.Float32FromString(result.ComponentRecipeScale.String)
			}

			meal.Components = append(meal.Components, &mealplanning.MealComponent{
				ComponentType: componentType,
				Recipe:        *recipe,
				RecipeScale:   recipeScale,
			})
		}

		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	if meal != nil {
		data = append(data, meal)
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(m *mealplanning.Meal) string { return m.ID },
		filter,
	)

	return x, nil
}

// GetMealsWithIDs fetches a list of meals from the database that have IDs within a given set.
func (q *repository) GetMealsWithIDs(ctx context.Context, ids []string) ([]*mealplanning.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if len(ids) == 0 {
		return []*mealplanning.Meal{}, nil
	}

	results, err := q.generatedQuerier.GetMealsWithIDs(ctx, q.readDB, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meals list retrieval by ids")
	}

	mealsByID := map[string]*mealplanning.Meal{}
	for _, result := range results {
		m, exists := mealsByID[result.ID]
		if !exists {
			m = &mealplanning.Meal{
				CreatedAt:     result.CreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
				ID:            result.ID,
				Description:   result.Description,
				CreatedByUser: result.CreatedByUser,
				Name:          result.Name,
				Components:    nil,
				EstimatedPortions: types.Float32RangeWithOptionalMax{
					Min: database.Float32FromString(result.MinEstimatedPortions),
					Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
				},
				EligibleForMealPlans: result.EligibleForMealPlans,
			}
			mealsByID[result.ID] = m
		}

		m.Components = append(m.Components, &mealplanning.MealComponent{
			ComponentType: string(result.ComponentMealComponentType),
			Recipe: mealplanning.Recipe{
				ID: result.ComponentRecipeID,
			},
			RecipeScale: database.Float32FromString(result.ComponentRecipeScale),
		})
	}

	meals := make([]*mealplanning.Meal, 0, len(mealsByID))
	for _, m := range mealsByID {
		meals = append(meals, m)
	}

	return meals, nil
}

// GetMealIDsThatNeedSearchIndexing fetches a list of meal IDs from the database that meet a particular filter.
func (q *repository) GetMealIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetMealsNeedingIndexing(ctx, q.readDB)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing meals list retrieval query")
	}

	return results, nil
}

// SearchForMeals fetches a list of recipes from the database that match a query.
func (q *repository) SearchForMeals(ctx context.Context, mealNameQuery string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.Meal], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	var (
		data          []*mealplanning.Meal
		filteredCount uint64
		totalCount    uint64
	)

	results, err := q.generatedQuerier.SearchForMeals(ctx, q.readDB, &generated.SearchForMealsParams{
		Query:           mealNameQuery,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:          database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:     database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meals list retrieval query")
	}

	var meal *mealplanning.Meal
	for _, result := range results {
		if meal != nil && meal.ID != result.ID {
			data = append(data, meal)
			meal = nil
		}

		if meal == nil {
			meal = &mealplanning.Meal{
				CreatedAt:     result.CreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
				ID:            result.ID,
				Description:   result.Description,
				CreatedByUser: result.CreatedByUser,
				Name:          result.Name,
				Components:    []*mealplanning.MealComponent{},
				EstimatedPortions: types.Float32RangeWithOptionalMax{
					Min: database.Float32FromString(result.MinEstimatedPortions),
					Max: database.Float32PointerFromNullString(result.MaxEstimatedPortions),
				},
				EligibleForMealPlans: result.EligibleForMealPlans,
			}
		}

		recipe, recipeErr := q.getRecipe(ctx, result.ComponentRecipeID)
		if recipeErr != nil {
			if errors.Is(recipeErr, sql.ErrNoRows) {
				logger.WithValue(mealplanningkeys.MealIDKey, result.ID).
					WithValue(mealplanningkeys.RecipeIDKey, result.ComponentRecipeID).
					Info("skipping meal component with missing or archived recipe")
				continue
			}
			return nil, observability.PrepareAndLogError(recipeErr, logger, span, "getting recipe for meal component")
		}

		meal.Components = append(meal.Components, &mealplanning.MealComponent{
			ComponentType: string(result.ComponentMealComponentType),
			Recipe:        *recipe,
			RecipeScale:   database.Float32FromString(result.ComponentRecipeScale),
		})

		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	if meal != nil {
		data = append(data, meal)
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(m *mealplanning.Meal) string { return m.ID },
		filter,
	)

	return x, nil
}

// CreateMeal creates a meal in the database.
func (q *repository) createMeal(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *mealplanning.MealDatabaseCreationInput) (*mealplanning.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	logger := q.logger.WithValue(mealplanningkeys.MealIDKey, input.ID).WithValue("meal.name", input.Name)

	// create the meal.
	if err := q.generatedQuerier.CreateMeal(ctx, querier, &generated.CreateMealParams{
		ID:                   input.ID,
		Name:                 input.Name,
		Description:          input.Description,
		MinEstimatedPortions: database.StringFromFloat32(input.EstimatedPortions.Min),
		CreatedByUser:        input.CreatedByUser,
		MaxEstimatedPortions: database.NullStringFromFloat32Pointer(input.EstimatedPortions.Max),
		EligibleForMealPlans: input.EligibleForMealPlans,
	}); err != nil {
		q.RollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal creation query")
	}

	x := &mealplanning.Meal{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: input.EstimatedPortions.Min,
			Max: input.EstimatedPortions.Max,
		},
		EligibleForMealPlans: input.EligibleForMealPlans,
		CreatedByUser:        input.CreatedByUser,
		CreatedAt:            q.CurrentTime(),
	}

	for _, recipeID := range input.Components {
		if err := q.CreateMealComponent(ctx, querier, x.ID, recipeID); err != nil {
			q.RollbackTransaction(ctx, querier)
			return nil, observability.PrepareAndLogError(err, logger, span, "creating meal recipe")
		}
	}

	tracing.AttachToSpan(span, mealplanningkeys.MealIDKey, x.ID)
	logger.Info("meal created")

	return x, nil
}

// CreateMeal creates a meal in the database.
func (q *repository) CreateMeal(ctx context.Context, input *mealplanning.MealDatabaseCreationInput) (*mealplanning.Meal, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	tx, err := q.writeDB.BeginTx(ctx, nil)
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
func (q *repository) CreateMealComponent(ctx context.Context, querier database.SQLQueryExecutor, mealID string, input *mealplanning.MealComponentDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return database.ErrNilInputProvided
	}

	if mealID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(mealplanningkeys.MealIDKey, mealID)
	tracing.AttachToSpan(span, mealplanningkeys.MealIDKey, mealID)

	// create the meal.
	if err := q.generatedQuerier.CreateMealComponent(ctx, querier, &generated.CreateMealComponentParams{
		ID:                identifiers.New(),
		MealID:            mealID,
		RecipeID:          input.RecipeID,
		MealComponentType: generated.ComponentType(input.ComponentType),
		RecipeScale:       database.StringFromFloat32(input.RecipeScale),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "creating meal component")
	}

	return nil
}

// MarkMealAsIndexed updates a particular meal's last_indexed_at value.
func (q *repository) MarkMealAsIndexed(ctx context.Context, mealID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(mealplanningkeys.MealIDKey, mealID)
	tracing.AttachToSpan(span, mealplanningkeys.MealIDKey, mealID)

	if _, err := q.generatedQuerier.UpdateMealLastIndexedAt(ctx, q.writeDB, mealID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking meal as indexed")
	}

	logger.Info("meal marked as indexed")

	return nil
}

// ArchiveMeal archives a meal from the database by its ID.
func (q *repository) ArchiveMeal(ctx context.Context, mealID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(mealplanningkeys.MealIDKey, mealID)
	tracing.AttachToSpan(span, mealplanningkeys.MealIDKey, mealID)

	if userID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(identitykeys.UserIDKey, userID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, userID)

	rowsAffected, err := q.generatedQuerier.ArchiveMeal(ctx, q.writeDB, &generated.ArchiveMealParams{
		CreatedByUser: userID,
		ID:            mealID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
