package postgres

import (
	"context"
	_ "embed"
	"time"

	"github.com/lib/pq"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.AdvancedPrepStepDataManager = (*Querier)(nil)
)

// scanAdvancedPrepStep takes a database Scanner (i.e. *sql.Row) and scans the result into a valid instrument struct.
func (q *Querier) scanAdvancedPrepStep(ctx context.Context, scan database.Scanner) (x *types.AdvancedPrepStep, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.AdvancedPrepStep{}

	targetVars := []interface{}{
		&x.ID,
		&x.MealPlanOption.ID,
		&x.MealPlanOption.AssignedCook,
		&x.MealPlanOption.AssignedDishwasher,
		&x.MealPlanOption.Chosen,
		&x.MealPlanOption.TieBroken,
		&x.MealPlanOption.Meal.ID,
		&x.MealPlanOption.Notes,
		&x.MealPlanOption.PrepStepsCreated,
		&x.MealPlanOption.CreatedAt,
		&x.MealPlanOption.LastUpdatedAt,
		&x.MealPlanOption.ArchivedAt,
		&x.MealPlanOption.BelongsToMealPlanEvent,
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
		&x.Status,
		&x.StatusExplanation,
		&x.CreationExplanation,
		&x.CannotCompleteBefore,
		&x.CannotCompleteAfter,
		&x.CreatedAt,
		&x.SettledAt,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, observability.PrepareError(err, span, "")
	}

	return x, nil
}

// scanAdvancedPrepSteps takes some database rows and turns them into a slice of advanced prep steps.
func (q *Querier) scanAdvancedPrepSteps(ctx context.Context, rows database.ResultIterator) (validInstruments []*types.AdvancedPrepStep, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, scanErr := q.scanAdvancedPrepStep(ctx, rows)
		if scanErr != nil {
			return nil, scanErr
		}

		validInstruments = append(validInstruments, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "handling rows")
	}

	return validInstruments, nil
}

//go:embed queries/advanced_prep_steps/exists.sql
var advancedPrepStepsExistsQuery string

// AdvancedPrepStepExists checks if an advanced prep step exists.
func (q *Querier) AdvancedPrepStepExists(ctx context.Context, mealPlanID, advancedPrepStepID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if advancedPrepStepID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AdvancedPrepStepIDKey, advancedPrepStepID)
	tracing.AttachAdvancedPrepStepIDToSpan(span, advancedPrepStepID)

	args := []interface{}{
		mealPlanID,
		advancedPrepStepID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, advancedPrepStepsExistsQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing advanced step existence check")
	}

	logger.Info("advanced step existence retrieved")

	return result, nil
}

//go:embed queries/advanced_prep_steps/get_one.sql
var getAdvancedPrepStepsQuery string

// GetAdvancedPrepStep fetches an advanced prep step.
func (q *Querier) GetAdvancedPrepStep(ctx context.Context, advancedPrepStepID string) (x *types.AdvancedPrepStep, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if advancedPrepStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AdvancedPrepStepIDKey, advancedPrepStepID)
	tracing.AttachAdvancedPrepStepIDToSpan(span, advancedPrepStepID)

	args := []interface{}{
		advancedPrepStepID,
	}

	rows := q.getOneRow(ctx, q.db, "advanced prep step", getAdvancedPrepStepsQuery, args)
	if x, err = q.scanAdvancedPrepStep(ctx, rows); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning advanced prep step")
	}

	logger.Info("advanced steps retrieved")

	return x, nil
}

//go:embed queries/advanced_prep_steps/list_all_by_meal_plan.sql
var listAdvancedPrepStepsForMealPlanQuery string

// GetAdvancedPrepStepsForMealPlan fetches a list of advanced prep steps.
func (q *Querier) GetAdvancedPrepStepsForMealPlan(ctx context.Context, mealPlanID string) (x []*types.AdvancedPrepStep, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = []*types.AdvancedPrepStep{}
	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	args := []interface{}{
		mealPlanID,
	}

	rows, getRowsErr := q.performReadQuery(ctx, q.db, "advanced prep steps list", listAdvancedPrepStepsForMealPlanQuery, args)
	if getRowsErr != nil {
		return nil, observability.PrepareAndLogError(getRowsErr, logger, span, "executing advanced prep steps list retrieval query")
	}

	x, scanErr := q.scanAdvancedPrepSteps(ctx, rows)
	if scanErr != nil {
		return nil, observability.PrepareAndLogError(scanErr, logger, span, "scanning advanced prep steps")
	}

	logger.Info("advanced steps retrieved")

	return x, nil
}

//go:embed queries/advanced_prep_steps/create.sql
var createAdvancedPrepStepQuery string

//go:embed queries/meal_plan_options/mark_as_steps_created.sql
var markMealPlanOptionAsHavingStepsCreatedQuery string

// CreateAdvancedPrepStepsForMealPlanOption creates advanced prep steps.
func (q *Querier) CreateAdvancedPrepStepsForMealPlanOption(ctx context.Context, mealPlanOptionID string, inputs []*types.AdvancedPrepStepDatabaseCreationInput) ([]*types.AdvancedPrepStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	outputs := []*types.AdvancedPrepStep{}
	for _, input := range inputs {
		createAdvancedPrepStepArgs := []interface{}{
			input.ID,
			mealPlanOptionID,
			input.RecipeStepID,
			input.Status,
			input.StatusExplanation,
			input.CreationExplanation,
			pq.FormatTimestamp(input.CannotCompleteBefore.Truncate(time.Second)),
			pq.FormatTimestamp(input.CannotCompleteAfter.Truncate(time.Second)),
		}

		if err = q.performWriteQuery(ctx, tx, "create advanced prep step", createAdvancedPrepStepQuery, createAdvancedPrepStepArgs); err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(err, logger, span, "create advanced prep step")
		}

		outputs = append(outputs, &types.AdvancedPrepStep{
			ID:                   input.ID,
			CannotCompleteBefore: input.CannotCompleteBefore.Truncate(time.Second),
			CannotCompleteAfter:  input.CannotCompleteAfter.Truncate(time.Second),
			MealPlanOption:       types.MealPlanOption{ID: mealPlanOptionID},
			RecipeStep:           types.RecipeStep{ID: input.RecipeStepID},
			CreatedAt:            q.currentTime(),
			Status:               input.Status,
			StatusExplanation:    input.StatusExplanation,
			CreationExplanation:  input.CreationExplanation,
			SettledAt:            input.CompletedAt,
		})
	}

	// mark prep steps as created for step
	markMealPlanOptionAsHavingStepsCreatedArgs := []interface{}{
		mealPlanOptionID,
	}

	if err = q.performWriteQuery(ctx, tx, "create advanced prep step", markMealPlanOptionAsHavingStepsCreatedQuery, markMealPlanOptionAsHavingStepsCreatedArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "create advanced prep step")
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, observability.PrepareAndLogError(commitErr, logger, span, "committing transaction")
	}

	logger.Info("advanced steps created")

	return outputs, nil
}

//go:embed queries/advanced_prep_steps/change_status.sql
var changeAdvancedPrepStepStatusQuery string

// ChangeAdvancedPrepStepStatus changes an advanced prep step's status.
func (q *Querier) ChangeAdvancedPrepStepStatus(ctx context.Context, input *types.AdvancedPrepStepStatusChangeRequestInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return ErrNilInputProvided
	}
	tracing.AttachAdvancedPrepStepIDToSpan(span, input.ID)
	logger = logger.WithValue(keys.AdvancedPrepStepIDKey, input.ID)

	var settledAt *time.Time
	if input.Status == types.AdvancedPrepStepStatusFinished {
		t := q.timeFunc()
		settledAt = &t
	}

	changeAdvancedPrepStepStatusArgs := []interface{}{
		input.ID,
		input.Status,
		input.StatusExplanation,
		settledAt,
	}

	if err := q.performWriteQuery(ctx, q.db, "prep step status change", changeAdvancedPrepStepStatusQuery, changeAdvancedPrepStepStatusArgs); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "changing prep step status")
	}

	logger.Info("prep step status changed")

	return nil
}
