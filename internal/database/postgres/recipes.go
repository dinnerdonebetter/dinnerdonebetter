package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"log"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.RecipeDataManager = (*Querier)(nil)
)

// RecipeExists fetches whether a recipe exists from the database.
func (q *Querier) RecipeExists(ctx context.Context, recipeID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	tracing.AttachRecipeIDToSpan(span, recipeID)

	result, err := q.generatedQuerier.CheckRecipeExistence(ctx, q.db, recipeID)
	if err != nil {
		return false, observability.PrepareError(err, span, "performing recipe existence check")
	}

	return result, nil
}

// scanRecipeAndStep takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe struct.
func (q *Querier) scanRecipeAndStep(ctx context.Context, scan database.Scanner) (x *types.Recipe, y *types.RecipeStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.Recipe{}
	y = &types.RecipeStep{}

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Slug,
		&x.Source,
		&x.Description,
		&x.InspiredByRecipeID,
		&x.MinimumEstimatedPortions,
		&x.MaximumEstimatedPortions,
		&x.PortionName,
		&x.PluralPortionName,
		&x.SealOfApproval,
		&x.EligibleForMeals,
		&x.YieldsComponentType,
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
		&y.Preparation.ConsumesVessel,
		&y.Preparation.OnlyForVessels,
		&y.Preparation.MinimumVesselCount,
		&y.Preparation.MaximumVesselCount,
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
		&y.StartTimerAutomatically,
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

//go:embed queries/recipes/get_by_id.sql
var getRecipeByIDQuery string

//go:embed queries/recipes/get_by_id_and_author_id.sql
var getRecipeByIDAndAuthorIDQuery string

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
		return nil, observability.PrepareError(err, span, "fetching recipe data")
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

	x.PrepTasks = []*types.RecipePrepTask{}
	x.SupportingRecipes = []*types.Recipe{}
	x.Media = []*types.RecipeMedia{}

	prepTasks, err := q.getRecipePrepTasksForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step prep tasks for recipe")
	}
	if prepTasks != nil {
		x.PrepTasks = prepTasks
	}

	recipeMedia, err := q.getRecipeMediaForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step ingredients for recipe")
	}
	if recipeMedia != nil {
		x.Media = recipeMedia
	}

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

	vessels, err := q.getRecipeStepVesselsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step vessels for recipe")
	}

	completionConditions, err := q.getRecipeStepCompletionConditionsForRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching recipe step completion conditions for recipe")
	}

	supportingRecipeInfo := map[string]bool{}
	for i, step := range x.Steps {
		for _, ingredient := range ingredients {
			if ingredient.RecipeStepProductRecipeID != nil {
				supportingRecipeInfo[*ingredient.RecipeStepProductRecipeID] = true
			}

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

		for _, vessel := range vessels {
			if vessel.BelongsToRecipeStep == step.ID {
				x.Steps[i].Vessels = append(x.Steps[i].Vessels, vessel)
			}
		}

		for _, completionCondition := range completionConditions {
			if completionCondition.BelongsToRecipeStep == step.ID {
				x.Steps[i].CompletionConditions = append(x.Steps[i].CompletionConditions, completionCondition)
			}
		}

		recipeMedia, err = q.getRecipeMediaForRecipeStep(ctx, recipeID, step.ID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "fetching recipe step ingredients for recipe")
		}
		x.Steps[i].Media = recipeMedia
	}

	var supportingRecipeIDs []string
	for supportingRecipe := range supportingRecipeInfo {
		supportingRecipeIDs = append(supportingRecipeIDs, supportingRecipe)
	}

	x.SupportingRecipes, err = q.GetRecipesWithIDs(ctx, supportingRecipeIDs)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching supporting recipes")
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

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.Recipe]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipes(ctx, q.db, &generated.GetRecipesParams{
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipes list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.Recipe{
			CreatedAt:                result.CreatedAt,
			InspiredByRecipeID:       stringPointerFromNullString(result.InspiredByRecipeID),
			LastUpdatedAt:            timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:               timePointerFromNullTime(result.ArchivedAt),
			MaximumEstimatedPortions: float32PointerFromNullString(result.MaxEstimatedPortions),
			PluralPortionName:        result.PluralPortionName,
			Description:              result.Description,
			Name:                     result.Name,
			PortionName:              result.PortionName,
			ID:                       result.ID,
			CreatedByUser:            result.CreatedByUser,
			Source:                   result.Source,
			Slug:                     result.Slug,
			YieldsComponentType:      string(result.YieldsComponentType),
			MinimumEstimatedPortions: float32FromString(result.MinEstimatedPortions),
			SealOfApproval:           result.SealOfApproval,
			EligibleForMeals:         result.EligibleForMeals,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetRecipesWithIDs fetches a list of recipes from the database that meet a particular filter.
func (q *Querier) GetRecipesWithIDs(ctx context.Context, ids []string) ([]*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	recipes := []*types.Recipe{}
	for _, id := range ids {
		r, err := q.getRecipe(ctx, id, "")
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe")
		}

		recipes = append(recipes, r)
	}

	return recipes, nil
}

// GetRecipeIDsThatNeedSearchIndexing fetches a list of recipe IDs from the database that meet a particular filter.
func (q *Querier) GetRecipeIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetRecipesNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing recipes list retrieval query")
	}

	return results, nil
}

// SearchForRecipes fetches a list of recipes from the database that match a query.
func (q *Querier) SearchForRecipes(ctx context.Context, recipeNameQuery string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Recipe], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.Recipe]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.RecipeSearch(ctx, q.db, &generated.RecipeSearchParams{
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
		Query:         recipeNameQuery,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipes search query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.Recipe{
			CreatedAt:                result.CreatedAt,
			InspiredByRecipeID:       stringPointerFromNullString(result.InspiredByRecipeID),
			LastUpdatedAt:            timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:               timePointerFromNullTime(result.ArchivedAt),
			MaximumEstimatedPortions: float32PointerFromNullString(result.MaxEstimatedPortions),
			PluralPortionName:        result.PluralPortionName,
			Description:              result.Description,
			Name:                     result.Name,
			PortionName:              result.PortionName,
			ID:                       result.ID,
			CreatedByUser:            result.CreatedByUser,
			Source:                   result.Source,
			Slug:                     result.Slug,
			YieldsComponentType:      string(result.YieldsComponentType),
			MinimumEstimatedPortions: float32FromString(result.MinEstimatedPortions),
			SealOfApproval:           result.SealOfApproval,
			EligibleForMeals:         result.EligibleForMeals,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateRecipe creates a recipe in the database.
func (q *Querier) CreateRecipe(ctx context.Context, input *types.RecipeDatabaseCreationInput) (*types.Recipe, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeIDKey, input.ID)
	tracing.AttachRecipeIDToSpan(span, input.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	// create the recipe.
	if err = q.generatedQuerier.CreateRecipe(ctx, tx, &generated.CreateRecipeParams{
		MinEstimatedPortions: stringFromFloat32(input.MinimumEstimatedPortions),
		ID:                   input.ID,
		Slug:                 input.Slug,
		Source:               input.Source,
		Description:          input.Description,
		CreatedByUser:        input.CreatedByUser,
		Name:                 input.Name,
		YieldsComponentType:  generated.ComponentType(input.YieldsComponentType),
		PortionName:          input.PortionName,
		PluralPortionName:    input.PluralPortionName,
		MaxEstimatedPortions: nullStringFromFloat32Pointer(input.MaximumEstimatedPortions),
		InspiredByRecipeID:   nullStringFromStringPointer(input.InspiredByRecipeID),
		SealOfApproval:       input.SealOfApproval,
		EligibleForMeals:     input.EligibleForMeals,
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe creation query")
	}

	x := &types.Recipe{
		ID:                       input.ID,
		Name:                     input.Name,
		Slug:                     input.Slug,
		Source:                   input.Source,
		Description:              input.Description,
		InspiredByRecipeID:       input.InspiredByRecipeID,
		CreatedByUser:            input.CreatedByUser,
		MinimumEstimatedPortions: input.MinimumEstimatedPortions,
		MaximumEstimatedPortions: input.MaximumEstimatedPortions,
		SealOfApproval:           input.SealOfApproval,
		EligibleForMeals:         input.EligibleForMeals,
		PortionName:              input.PortionName,
		PluralPortionName:        input.PluralPortionName,
		YieldsComponentType:      input.YieldsComponentType,
		CreatedAt:                q.currentTime(),
		PrepTasks:                []*types.RecipePrepTask{},
		Steps:                    []*types.RecipeStep{},
		Media:                    []*types.RecipeMedia{},
		SupportingRecipes:        []*types.Recipe{},
	}

	findCreatedRecipeStepProductsForIngredients(input)
	findCreatedRecipeStepProductsForInstruments(input)
	findCreatedRecipeStepProductsForVessels(input)

	for i, stepInput := range input.Steps {
		stepInput.Index = uint32(i)
		stepInput.BelongsToRecipe = x.ID

		var s *types.RecipeStep
		s, err = q.createRecipeStep(ctx, tx, stepInput)
		if err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(err, span, "creating recipe step #%d", i+1)
		}

		x.Steps = append(x.Steps, s)
	}

	for i, prepTaskInput := range input.PrepTasks {
		var pt *types.RecipePrepTask
		pt, err = q.createRecipePrepTask(ctx, tx, prepTaskInput)
		if err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(err, span, "creating recipe prep task #%d", i+1)
		}

		x.PrepTasks = append(x.PrepTasks, pt)
	}

	if input.AlsoCreateMeal {
		if _, err = q.createMeal(ctx, tx, &types.MealDatabaseCreationInput{
			ID:                       identifiers.New(),
			Name:                     x.Name,
			Description:              x.Description,
			MinimumEstimatedPortions: x.MinimumEstimatedPortions,
			MaximumEstimatedPortions: x.MaximumEstimatedPortions,
			EligibleForMealPlans:     x.EligibleForMeals,
			CreatedByUser:            x.CreatedByUser,
			Components: []*types.MealComponentDatabaseCreationInput{
				{
					RecipeID:      x.ID,
					RecipeScale:   1.0,
					ComponentType: types.MealComponentTypesMain,
				},
			},
		}); err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(err, span, "creating meal from recipe")
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe created")

	return x, nil
}

func findCreatedRecipeStepProductsForIngredients(recipe *types.RecipeDatabaseCreationInput) {
	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.ProductOfRecipeStepIndex != nil && ingredient.ProductOfRecipeStepProductIndex != nil {
				enoughSteps := len(recipe.Steps) > int(*ingredient.ProductOfRecipeStepIndex)
				enoughRecipeStepProducts := len(recipe.Steps[int(*ingredient.ProductOfRecipeStepIndex)].Products) > int(*ingredient.ProductOfRecipeStepProductIndex)
				relevantProductIsIngredient := recipe.Steps[*ingredient.ProductOfRecipeStepIndex].Products[*ingredient.ProductOfRecipeStepProductIndex].Type == types.RecipeStepProductIngredientType
				if enoughSteps && enoughRecipeStepProducts && relevantProductIsIngredient {
					ingredient.RecipeStepProductID = &recipe.Steps[*ingredient.ProductOfRecipeStepIndex].Products[*ingredient.ProductOfRecipeStepProductIndex].ID
				}
			}
		}
	}
}

func findCreatedRecipeStepProductsForInstruments(recipe *types.RecipeDatabaseCreationInput) {
	for _, step := range recipe.Steps {
		for _, instrument := range step.Instruments {
			if instrument.ProductOfRecipeStepIndex != nil && instrument.ProductOfRecipeStepProductIndex != nil {
				enoughSteps := len(recipe.Steps) > int(*instrument.ProductOfRecipeStepIndex)
				enoughRecipeStepProducts := len(recipe.Steps[int(*instrument.ProductOfRecipeStepIndex)].Products) > int(*instrument.ProductOfRecipeStepProductIndex)
				relevantProductIsInstrument := recipe.Steps[*instrument.ProductOfRecipeStepIndex].Products[*instrument.ProductOfRecipeStepProductIndex].Type == types.RecipeStepProductInstrumentType
				if enoughSteps && enoughRecipeStepProducts && relevantProductIsInstrument {
					instrument.RecipeStepProductID = &recipe.Steps[*instrument.ProductOfRecipeStepIndex].Products[*instrument.ProductOfRecipeStepProductIndex].ID
				}
			}
		}
	}
}

func findCreatedRecipeStepProductsForVessels(recipe *types.RecipeDatabaseCreationInput) {
	for _, step := range recipe.Steps {
		for _, vessel := range step.Vessels {
			if vessel.ProductOfRecipeStepIndex != nil && vessel.ProductOfRecipeStepProductIndex != nil {
				enoughSteps := len(recipe.Steps) > int(*vessel.ProductOfRecipeStepIndex)
				enoughRecipeStepProducts := len(recipe.Steps[int(*vessel.ProductOfRecipeStepIndex)].Products) > int(*vessel.ProductOfRecipeStepProductIndex)
				relevantProductIsVessel := recipe.Steps[*vessel.ProductOfRecipeStepIndex].Products[*vessel.ProductOfRecipeStepProductIndex].Type == types.RecipeStepProductVesselType
				if enoughSteps && enoughRecipeStepProducts && relevantProductIsVessel {
					vessel.RecipeStepProductID = &recipe.Steps[*vessel.ProductOfRecipeStepIndex].Products[*vessel.ProductOfRecipeStepProductIndex].ID
				} else {
					log.Printf("for recipe step id %q, vessel ID %q, not enough steps: %t, not enough recipe step products: %t, relevant product is vessel: %t", step.ID, vessel.ID, enoughSteps, enoughRecipeStepProducts, relevantProductIsVessel)
				}
			}
		}
	}
}

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

	if err := q.generatedQuerier.UpdateRecipe(ctx, q.db, &generated.UpdateRecipeParams{
		Name:                 updated.Name,
		Slug:                 updated.Slug,
		Source:               updated.Source,
		Description:          updated.Description,
		InspiredByRecipeID:   nullStringFromStringPointer(updated.InspiredByRecipeID),
		MinEstimatedPortions: stringFromFloat32(updated.MinimumEstimatedPortions),
		MaxEstimatedPortions: nullStringFromFloat32Pointer(updated.MaximumEstimatedPortions),
		PortionName:          updated.PortionName,
		PluralPortionName:    updated.PluralPortionName,
		SealOfApproval:       updated.SealOfApproval,
		EligibleForMeals:     updated.EligibleForMeals,
		YieldsComponentType:  generated.ComponentType(updated.YieldsComponentType),
		CreatedByUser:        updated.CreatedByUser,
		ID:                   updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe")
	}

	logger.Info("recipe updated")

	return nil
}

// MarkRecipeAsIndexed updates a particular recipe's last_indexed_at value.
func (q *Querier) MarkRecipeAsIndexed(ctx context.Context, recipeID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if err := q.generatedQuerier.UpdateRecipeLastIndexedAt(ctx, q.db, recipeID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking recipe as indexed")
	}

	logger.Info("recipe marked as indexed")

	return nil
}

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

	if err := q.generatedQuerier.ArchiveRecipe(ctx, q.db, &generated.ArchiveRecipeParams{
		CreatedByUser: userID,
		ID:            recipeID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe")
	}

	logger.Info("recipe archived")

	return nil
}
