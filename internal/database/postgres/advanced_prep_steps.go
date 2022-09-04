package postgres

import (
	"context"
	_ "embed"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.AdminUserDataManager = (*SQLQuerier)(nil)

// scanAdvancedPrepStep takes a database Scanner (i.e. *sql.Row) and scans the result into a valid instrument struct.
func (q *SQLQuerier) scanAdvancedPrepStep(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.AdvancedPrepStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.AdvancedPrepStep{}

	targetVars := []interface{}{
		&x.ID,
		&x.MealPlanOption.ID,
		&x.MealPlanOption.Day,
		&x.MealPlanOption.AssignedCook,
		&x.MealPlanOption.AssignedDishwasher,
		&x.MealPlanOption.MealName,
		&x.MealPlanOption.Chosen,
		&x.MealPlanOption.TieBroken,
		&x.MealPlanOption.Meal.ID,
		&x.MealPlanOption.Notes,
		&x.MealPlanOption.PrepStepsCreated,
		&x.MealPlanOption.CreatedAt,
		&x.MealPlanOption.LastUpdatedAt,
		&x.MealPlanOption.ArchivedAt,
		&x.MealPlanOption.BelongsToMealPlan,
		&x.RecipeStep.ID,
		&x.RecipeStep.Index,
		&x.RecipeStep.Preparation.ID,
		&x.RecipeStep.Preparation.Name,
		&x.RecipeStep.Preparation.Description,
		&x.RecipeStep.Preparation.IconPath,
		&x.RecipeStep.Preparation.YieldsNothing,
		&x.RecipeStep.Preparation.RestrictToIngredients,
		&x.RecipeStep.Preparation.ZeroIngredientsAllowable,
		&x.RecipeStep.Preparation.PastTense,
		&x.RecipeStep.Preparation.CreatedAt,
		&x.RecipeStep.Preparation.LastUpdatedAt,
		&x.RecipeStep.Preparation.ArchivedAt,
		&x.RecipeStep.MinimumEstimatedTimeInSeconds,
		&x.RecipeStep.MaximumEstimatedTimeInSeconds,
		&x.RecipeStep.MinimumTemperatureInCelsius,
		&x.RecipeStep.MaximumTemperatureInCelsius,
		&x.RecipeStep.Notes,
		&x.RecipeStep.ExplicitInstructions,
		&x.RecipeStep.Optional,
		&x.RecipeStep.CreatedAt,
		&x.RecipeStep.LastUpdatedAt,
		&x.RecipeStep.ArchivedAt,
		&x.RecipeStep.BelongsToRecipe,
		&x.CannotCompleteBefore,
		&x.CannotCompleteAfter,
		&x.CreatedAt,
		&x.CompletedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanAdvancedPrepSteps takes some database rows and turns them into a slice of valid instruments.
func (q *SQLQuerier) scanAdvancedPrepSteps(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validInstruments []*types.AdvancedPrepStep, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanAdvancedPrepStep(ctx, rows, includeCounts)
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

		validInstruments = append(validInstruments, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validInstruments, filteredCount, totalCount, nil
}

//go:embed queries/advanced_prep_step_list_all_by_meal_plan_option.sql
var listAdvancedPrepStepsQuery string

// GetAdvancedPrepStepsForMealPlanOptionID lists advanced prep steps for a given meal plan option.
func (q *SQLQuerier) GetAdvancedPrepStepsForMealPlanOptionID(ctx context.Context, mealPlanOptionID string, filter *types.QueryFilter) (x *types.AdvancedPrepStepList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.AdvancedPrepStepList{}
	logger := filter.AttachToLogger(q.logger.Clone())
	tracing.AttachQueryFilterToSpan(span, filter)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	args := []interface{}{
		mealPlanOptionID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "advanced prep steps list", listAdvancedPrepStepsQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing advanced prep steps list retrieval query")
	}

	if x.AdvancedPrepSteps, x.FilteredCount, x.TotalCount, err = q.scanAdvancedPrepSteps(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning advanced prep steps")
	}

	logger.Info("advanced steps retrieved")

	return x, nil
}

//go:embed queries/advanced_prep_step_create.sql
var createAdvancedPrepStepQuery string

// CreateAdvancedPrepStep updates a user's household status.
func (q *SQLQuerier) CreateAdvancedPrepStep(ctx context.Context, input *types.AdvancedPrepStepDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	args := []interface{}{
		input.ID,
		input.MealPlanOptionID,
		input.RecipeStepID,
		input.CannotCompleteBefore,
		input.CannotCompleteAfter,
	}

	if err := q.performWriteQuery(ctx, q.db, "create advanced prep step", createAdvancedPrepStepQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "create advanced prep step")
	}

	logger.Info("advanced step created")

	return nil
}
