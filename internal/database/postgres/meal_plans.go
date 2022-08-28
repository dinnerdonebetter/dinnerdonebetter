package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.MealPlanDataManager = (*SQLQuerier)(nil)

	// mealPlansTableColumns are the columns for the meal_plans table.
	mealPlansTableColumns = []string{
		"meal_plans.id",
		"meal_plans.notes",
		"meal_plans.status",
		"meal_plans.voting_deadline",
		"meal_plans.starts_at",
		"meal_plans.ends_at",
		"meal_plans.created_at",
		"meal_plans.last_updated_at",
		"meal_plans.archived_at",
		"meal_plans.belongs_to_household",
	}
)

// scanMealPlan takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan struct.
func (q *SQLQuerier) scanMealPlan(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealPlan, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.MealPlan{}

	targetVars := []interface{}{
		&x.ID,
		&x.Notes,
		&x.Status,
		&x.VotingDeadline,
		&x.StartsAt,
		&x.EndsAt,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
		&x.BelongsToHousehold,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanMealPlans takes some database rows and turns them into a slice of meal plans.
func (q *SQLQuerier) scanMealPlans(ctx context.Context, rows database.ResultIterator, includeCounts bool) (mealPlans []*types.MealPlan, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanMealPlan(ctx, rows, includeCounts)
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

		mealPlans = append(mealPlans, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return mealPlans, filteredCount, totalCount, nil
}

// scanFullMealPlan takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan struct.
func (q *SQLQuerier) scanFullMealPlan(ctx context.Context, scan database.Scanner) (mealPlan *types.MealPlan, mealPlanOption *types.MealPlanOption, mealPlanOptionVote *types.MealPlanOptionVote, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger
	mealPlan = &types.MealPlan{}
	mealPlanOption = &types.MealPlanOption{
		Meal: types.Meal{Recipes: []*types.Recipe{{}}},
	}
	nmpov := &nullableMealPlanOptionVote{}

	targetVars := []interface{}{
		&mealPlan.ID,
		&mealPlan.Notes,
		&mealPlan.Status,
		&mealPlan.VotingDeadline,
		&mealPlan.StartsAt,
		&mealPlan.EndsAt,
		&mealPlan.CreatedAt,
		&mealPlan.LastUpdatedAt,
		&mealPlan.ArchivedAt,
		&mealPlan.BelongsToHousehold,
		&mealPlanOption.ID,
		&mealPlanOption.Day,
		&mealPlanOption.AssignedCook,
		&mealPlanOption.MealName,
		&mealPlanOption.Chosen,
		&mealPlanOption.TieBroken,
		&mealPlanOption.Meal.ID,
		&mealPlanOption.Notes,
		&mealPlanOption.CreatedAt,
		&mealPlanOption.LastUpdatedAt,
		&mealPlanOption.ArchivedAt,
		&mealPlanOption.BelongsToMealPlan,
		&nmpov.ID,
		&nmpov.Rank,
		&nmpov.Abstain,
		&nmpov.Notes,
		&nmpov.ByUser,
		&nmpov.CreatedAt,
		&nmpov.LastUpdatedAT,
		&nmpov.ArchivedAt,
		&nmpov.BelongsToMealPlanOption,
		&mealPlanOption.Meal.ID,
		&mealPlanOption.Meal.Name,
		&mealPlanOption.Meal.Description,
		&mealPlanOption.Meal.CreatedAt,
		&mealPlanOption.Meal.LastUpdatedAt,
		&mealPlanOption.Meal.ArchivedAt,
		&mealPlanOption.Meal.CreatedByUser,
		&mealPlanOption.Meal.Recipes[0].ID,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, nil, nil, observability.PrepareError(err, logger, span, "")
	}

	if nmpov.ID != nil {
		mealPlanOptionVote = &types.MealPlanOptionVote{
			LastUpdatedAt:           nmpov.LastUpdatedAT,
			ArchivedAt:              nmpov.ArchivedAt,
			ID:                      *nmpov.ID,
			Notes:                   *nmpov.Notes,
			BelongsToMealPlanOption: *nmpov.BelongsToMealPlanOption,
			ByUser:                  *nmpov.ByUser,
			CreatedAt:               *nmpov.CreatedAt,
			Rank:                    *nmpov.Rank,
			Abstain:                 *nmpov.Abstain,
		}
	}

	return mealPlan, mealPlanOption, mealPlanOptionVote, nil
}

const mealPlanExistenceQuery = "SELECT EXISTS ( SELECT meal_plans.id FROM meal_plans WHERE meal_plans.archived_at IS NULL AND meal_plans.id = $1 )"

// MealPlanExists fetches whether a meal plan exists from the database.
func (q *SQLQuerier) MealPlanExists(ctx context.Context, mealPlanID, householdID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		mealPlanID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing meal plan existence check")
	}

	return result, nil
}

const baseGetMealPlanQuery = `SELECT
	meal_plans.id,
	meal_plans.notes,
	meal_plans.status,
	meal_plans.voting_deadline,
	meal_plans.starts_at,
	meal_plans.ends_at,
	meal_plans.created_at,
	meal_plans.last_updated_at,
	meal_plans.archived_at,
	meal_plans.belongs_to_household,
    meal_plan_options.id,
    meal_plan_options.day,
    meal_plan_options.assigned_cook,
    meal_plan_options.meal_name,
    meal_plan_options.chosen,
    meal_plan_options.tiebroken,
    meal_plan_options.meal_id,
    meal_plan_options.notes,
    meal_plan_options.created_at,
    meal_plan_options.last_updated_at,
	meal_plan_options.archived_at,
    meal_plan_options.belongs_to_meal_plan,
	meal_plan_option_votes.id,
	meal_plan_option_votes.rank,
	meal_plan_option_votes.abstain,
	meal_plan_option_votes.notes,
	meal_plan_option_votes.by_user,
	meal_plan_option_votes.created_at,
	meal_plan_option_votes.last_updated_at,
	meal_plan_option_votes.archived_at,
	meal_plan_option_votes.belongs_to_meal_plan_option,
	meals.id,
	meals.name,
	meals.description,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
    meal_recipes.recipe_id
FROM meal_plans
	FULL OUTER JOIN meal_plan_options ON meal_plan_options.belongs_to_meal_plan=meal_plans.id
	FULL OUTER JOIN meal_plan_option_votes ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
	FULL OUTER JOIN meal_recipes ON meal_plan_options.meal_id=meal_recipes.meal_id
	FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
WHERE meal_plans.archived_at IS NULL 
AND meal_plans.id = $1
AND meal_plans.belongs_to_household = $2
`

const getMealPlanQuery = baseGetMealPlanQuery + `
ORDER BY meal_plan_options.id
`

const getMealPlanPastVotingDeadlineQuery = baseGetMealPlanQuery + `
AND meal_plans.status = 'awaiting_votes'
AND extract(epoch from NOW()) > meal_plans.voting_deadline
ORDER BY meal_plan_options.id
`

// GetMealPlan fetches a meal plan from the database.
func (q *SQLQuerier) getMealPlan(ctx context.Context, mealPlanID, householdID string, restrictToPastVotingDeadline bool) (*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		mealPlanID,
		householdID,
	}

	query := getMealPlanQuery
	if restrictToPastVotingDeadline {
		query = getMealPlanPastVotingDeadlineQuery
	}

	rows, err := q.performReadQuery(ctx, q.db, "meal plan", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan with options retrieval query")
	}

	var (
		// this is Go parlance for "map of string to set of strings"
		mealRecipeIDs      = map[string]map[string]struct{}{}
		allRecipeIDs       = []string{}
		mealPlan           *types.MealPlan
		currentOptionIndex = 0
	)
	for rows.Next() {
		rowMealPlan, rowMealPlanOption, rowMealPlanOptionVote, scanErr := q.scanFullMealPlan(ctx, rows)
		if scanErr != nil {
			return nil, scanErr
		}

		if len(rowMealPlanOption.Meal.Recipes) > 0 {
			allRecipeIDs = append(allRecipeIDs, rowMealPlanOption.Meal.Recipes[0].ID)
			if _, ok := mealRecipeIDs[rowMealPlanOption.Meal.ID]; ok {
				mealRecipeIDs[rowMealPlanOption.Meal.ID][rowMealPlanOption.Meal.Recipes[0].ID] = struct{}{}
			} else {
				mealRecipeIDs[rowMealPlanOption.Meal.ID] = map[string]struct{}{rowMealPlanOption.Meal.Recipes[0].ID: {}}
			}
		}

		if mealPlan == nil {
			mealPlan = rowMealPlan
		}

		if len(mealPlan.Options) == 0 && currentOptionIndex == 0 {
			mealPlan.Options = append(mealPlan.Options, rowMealPlanOption)
		}

		if mealPlan.Options[currentOptionIndex].ID != rowMealPlanOption.ID {
			currentOptionIndex++
			mealPlan.Options = append(mealPlan.Options, rowMealPlanOption)
		}

		if rowMealPlanOptionVote != nil {
			mealPlan.Options[currentOptionIndex].Votes = append(mealPlan.Options[currentOptionIndex].Votes, rowMealPlanOptionVote)
		}
	}

	if mealPlan == nil {
		return nil, sql.ErrNoRows
	}

	// I'm sure this is like `O((n^5)^2)` or something, but fuck off, it works, and I'm busy.
	// go through all the meal recipe IDs and fetch the into a map, so that we only fetch them once.
	// as of this writing they have to be unique across a meal plan, but we shouldn't bank on that.
	fetchedRecipes := map[string]*types.Recipe{}
	for _, recipeID := range allRecipeIDs {
		if _, ok := fetchedRecipes[recipeID]; !ok {
			recipe, getRecipeErr := q.getRecipe(ctx, recipeID, "")
			if getRecipeErr != nil {
				tracing.AttachRecipeIDToSpan(span, recipeID)
				return nil, observability.PrepareError(getRecipeErr, logger, span, "fetching recipe from meal plan")
			}

			fetchedRecipes[recipeID] = recipe
		}
	}

	// hydrate the options in the meal plan with the appropriate recipes.
	for i, option := range mealPlan.Options {
		mealPlan.Options[i].Meal.Recipes = []*types.Recipe{}
		for recipeID := range mealRecipeIDs[option.Meal.ID] {
			option.Meal.Recipes = append(option.Meal.Recipes, fetchedRecipes[recipeID])
		}
	}

	return mealPlan, nil
}

// GetMealPlan fetches a meal plan from the database.
func (q *SQLQuerier) GetMealPlan(ctx context.Context, mealPlanID, householdID string) (*types.MealPlan, error) {
	return q.getMealPlan(ctx, mealPlanID, householdID, false)
}

// GetMealPlans fetches a list of meal plans from the database that meet a particular filter.
func (q *SQLQuerier) GetMealPlans(ctx context.Context, householdID string, filter *types.QueryFilter) (x *types.MealPlanList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.MealPlanList{}
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

	query, args := q.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, householdID, false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "mealPlans", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plans list retrieval query")
	}

	if x.MealPlans, x.FilteredCount, x.TotalCount, err = q.scanMealPlans(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plans")
	}

	return x, nil
}

const mealPlanCreationQuery = "INSERT INTO meal_plans (id,notes,status,voting_deadline,starts_at,ends_at,belongs_to_household) VALUES ($1,$2,$3,$4,$5,$6,$7)"

// CreateMealPlan creates a meal plan in the database.
func (q *SQLQuerier) CreateMealPlan(ctx context.Context, input *types.MealPlanDatabaseCreationInput) (*types.MealPlan, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Notes,
		input.Status,
		input.VotingDeadline,
		input.StartsAt,
		input.EndsAt,
		input.BelongsToHousehold,
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	// create the meal plan.
	if err = q.performWriteQuery(ctx, tx, "meal plan creation", mealPlanCreationQuery, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating meal plan")
	}

	x := &types.MealPlan{
		ID:                 input.ID,
		Notes:              input.Notes,
		Status:             input.Status,
		VotingDeadline:     input.VotingDeadline,
		StartsAt:           input.StartsAt,
		EndsAt:             input.EndsAt,
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedAt:          q.currentTime(),
	}

	for _, option := range input.Options {
		option.BelongsToMealPlan = x.ID
		opt, createErr := q.createMealPlanOption(ctx, tx, option)
		if createErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(createErr, logger, span, "creating meal plan option for meal plan")
		}
		x.Options = append(x.Options, opt)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachMealPlanIDToSpan(span, x.ID)
	logger.Info("meal plan created")

	return x, nil
}

const updateMealPlanQuery = "UPDATE meal_plans SET notes = $1, status = $2, voting_deadline = $3, starts_at = $4, ends_at = $5, last_updated_at = extract(epoch FROM NOW()) WHERE archived_at IS NULL AND belongs_to_household = $6 AND id = $7"

// UpdateMealPlan updates a particular meal plan.
func (q *SQLQuerier) UpdateMealPlan(ctx context.Context, updated *types.MealPlan) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanIDKey, updated.ID)
	tracing.AttachMealPlanIDToSpan(span, updated.ID)
	tracing.AttachHouseholdIDToSpan(span, updated.BelongsToHousehold)

	args := []interface{}{
		updated.Notes,
		updated.Status,
		updated.VotingDeadline,
		updated.StartsAt,
		updated.EndsAt,
		updated.BelongsToHousehold,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan update", updateMealPlanQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan")
	}

	logger.Info("meal plan updated")

	return nil
}

const archiveMealPlanQuery = "UPDATE meal_plans SET archived_at = extract(epoch FROM NOW()) WHERE archived_at IS NULL AND belongs_to_household = $1 AND id = $2"

// ArchiveMealPlan archives a meal plan from the database by its ID.
func (q *SQLQuerier) ArchiveMealPlan(ctx context.Context, mealPlanID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		householdID,
		mealPlanID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan archive", archiveMealPlanQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan")
	}

	logger.Info("meal plan archived")

	return nil
}

var allDays = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
	time.Sunday,
}

var allMealNames = []types.MealName{
	types.BreakfastMealName,
	types.SecondBreakfastMealName,
	types.BrunchMealName,
	types.LunchMealName,
	types.SupperMealName,
	types.DinnerMealName,
}

func byDayAndMeal(l []*types.MealPlanOption, day time.Weekday, meal types.MealName) []*types.MealPlanOption {
	out := []*types.MealPlanOption{}

	for _, o := range l {
		if o.Day == day && o.MealName == meal {
			out = append(out, o)
		}
	}

	return out
}

const finalizeMealPlanQuery = `
	UPDATE meal_plans SET status = $1 WHERE archived_at IS NULL AND id = $2
`

// AttemptToFinalizeMealPlan finalizes a meal plan if all of its options have a selection.
func (q *SQLQuerier) AttemptToFinalizeMealPlan(ctx context.Context, mealPlanID, householdID string) (finalized bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// fetch household
	household, err := q.GetHouseholdByID(ctx, householdID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "fetching household")
	}

	// fetch meal plan
	mealPlan, err := q.getMealPlan(ctx, mealPlanID, householdID, false)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "fetching meal plan")
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	allOptionsChosen := true
	for _, day := range allDays {
		for _, mealName := range allMealNames {
			options := byDayAndMeal(mealPlan.Options, day, mealName)

			if len(options) > 0 {
				availableVotes := map[string]bool{}
				for _, member := range household.Members {
					availableVotes[member.BelongsToUser.ID] = false
				}

				alreadyChosen := false
				for _, opt := range options {
					if opt.Chosen {
						alreadyChosen = true
					}
					for _, vote := range opt.Votes {
						if _, ok := availableVotes[vote.ByUser]; ok {
							availableVotes[vote.ByUser] = true
						}
					}
				}

				if alreadyChosen {
					continue
				}

				for _, vote := range availableVotes {
					if !vote {
						allOptionsChosen = false
						continue
					}
				}

				// if we get here, then the tally is ready to be calculated for this set of options

				winner, tiebroken, chosen := q.decideOptionWinner(ctx, options)
				if chosen {
					args := []interface{}{
						mealPlanID,
						winner,
						tiebroken,
					}

					logger = logger.WithValue("winner", winner).WithValue("tiebroken", tiebroken)

					if err = q.performWriteQuery(ctx, tx, "meal plan option finalization", finalizeMealPlanOptionQuery, args); err != nil {
						q.rollbackTransaction(ctx, tx)
						return false, observability.PrepareError(err, logger, span, "finalizing meal plan option")
					}

					logger.Debug("finalized meal plan option")
				}
			}
		}
	}

	if allOptionsChosen {
		args := []interface{}{
			types.FinalizedMealPlanStatus,
			mealPlanID,
		}

		if err = q.performWriteQuery(ctx, tx, "meal plan finalization", finalizeMealPlanQuery, args); err != nil {
			q.rollbackTransaction(ctx, tx)
			return false, observability.PrepareError(err, logger, span, "finalizing meal plan")
		}

		finalized = true
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return false, observability.PrepareError(commitErr, logger, span, "committing transaction")
	}

	return finalized, nil
}

const getExpiredAndUnresolvedMealPlanIDsQuery = `
SELECT
	meal_plans.id,
	meal_plans.notes,
	meal_plans.status,
	meal_plans.voting_deadline,
	meal_plans.starts_at,
	meal_plans.ends_at,
	meal_plans.created_at,
	meal_plans.last_updated_at,
	meal_plans.archived_at,
	meal_plans.belongs_to_household
FROM meal_plans
WHERE meal_plans.archived_at IS NULL 
	AND meal_plans.status = 'awaiting_votes'
	AND to_timestamp(voting_deadline)::date < now()
GROUP BY meal_plans.id
ORDER BY meal_plans.id
`

// GetUnfinalizedMealPlansWithExpiredVotingPeriods gets unfinalized meal plans with expired voting deadlines.
func (q *SQLQuerier) GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlan, error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	rows, err := q.performReadQuery(ctx, q.db, "meal plan", getExpiredAndUnresolvedMealPlanIDsQuery, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan with options retrieval query")
	}

	mealPlans := []*types.MealPlan{}
	for rows.Next() {
		mp, _, _, scanErr := q.scanMealPlan(ctx, rows, false)
		if scanErr != nil {
			return nil, observability.PrepareError(scanErr, logger, span, "scanning meal plan response")
		}
		mealPlans = append(mealPlans, mp)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, logger, span, "closing rows")
	}

	return mealPlans, nil
}
