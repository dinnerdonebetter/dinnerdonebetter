package postgres

import (
	"context"
	"math/rand"
	"time"

	"resenje.org/schulze"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	mealPlansOnMealPlanOptionsJoinClause = "meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id"
	mealsOnMealPlanOptionsJoinClause     = "meals ON meal_plan_options.meal_id=meals.id"
)

var (
	_ types.MealPlanOptionDataManager = (*Querier)(nil)

	// mealPlanOptionsTableColumns are the columns for the meal_plan_options table.
	mealPlanOptionsTableColumns = []string{
		"meal_plan_options.id",
		"meal_plan_options.day",
		"meal_plan_options.assigned_cook",
		"meal_plan_options.assigned_dishwasher",
		"meal_plan_options.meal_name",
		"meal_plan_options.chosen",
		"meal_plan_options.tiebroken",
		"meal_plan_options.meal_id",
		"meal_plan_options.notes",
		"meal_plan_options.prep_steps_created",
		"meal_plan_options.created_at",
		"meal_plan_options.last_updated_at",
		"meal_plan_options.archived_at",
		"meal_plan_options.belongs_to_meal_plan",
		"meals.id",
		"meals.name",
		"meals.description",
		"meals.created_at",
		"meals.last_updated_at",
		"meals.archived_at",
		"meals.created_by_user",
	}

	getMealPlanOptionsJoins = []string{
		mealPlansOnMealPlanOptionsJoinClause,
		mealsOnMealPlanOptionsJoinClause,
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

// scanMealPlanOption takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan option struct.
func (q *Querier) scanMealPlanOption(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealPlanOption, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.MealPlanOption{
		Votes: []*types.MealPlanOptionVote{},
	}

	targetVars := []interface{}{
		&x.ID,
		&x.Day,
		&x.AssignedCook,
		&x.AssignedDishwasher,
		&x.MealName,
		&x.Chosen,
		&x.TieBroken,
		&x.Meal.ID,
		&x.Notes,
		&x.PrepStepsCreated,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
		&x.BelongsToMealPlan,
		&x.Meal.ID,
		&x.Meal.Name,
		&x.Meal.Description,
		&x.Meal.CreatedAt,
		&x.Meal.LastUpdatedAt,
		&x.Meal.ArchivedAt,
		&x.Meal.CreatedByUser,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanMealPlanOptions takes some database rows and turns them into a slice of meal plan options.
func (q *Querier) scanMealPlanOptions(ctx context.Context, rows database.ResultIterator, includeCounts bool) (mealPlanOptions []*types.MealPlanOption, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanMealPlanOption(ctx, rows, includeCounts)
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

		mealPlanOptions = append(mealPlanOptions, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return mealPlanOptions, filteredCount, totalCount, nil
}

const mealPlanOptionExistenceQuery = "SELECT EXISTS ( SELECT meal_plan_options.id FROM meal_plan_options JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id WHERE meal_plan_options.archived_at IS NULL AND meal_plan_options.belongs_to_meal_plan = $1 AND meal_plan_options.id = $2 AND meal_plans.archived_at IS NULL AND meal_plans.id = $3 )"

// MealPlanOptionExists fetches whether a meal plan option exists from the database.
func (q *Querier) MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanOptionID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	args := []interface{}{
		mealPlanID,
		mealPlanOptionID,
		mealPlanID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanOptionExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing meal plan option existence check")
	}

	return result, nil
}

const getMealPlanOptionQuery = `SELECT
	meal_plan_options.id,
	meal_plan_options.day,
	meal_plan_options.assigned_cook,
	meal_plan_options.assigned_dishwasher,
	meal_plan_options.meal_name,
	meal_plan_options.chosen,
	meal_plan_options.tiebroken,
	meal_plan_options.meal_id,
	meal_plan_options.notes,
	meal_plan_options.prep_steps_created,
	meal_plan_options.created_at,
	meal_plan_options.last_updated_at,
	meal_plan_options.archived_at,
	meal_plan_options.belongs_to_meal_plan,
	meals.id,
	meals.name,
	meals.description,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user
FROM meal_plan_options
JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id
JOIN meals ON meal_plan_options.meal_id=meals.id
WHERE meal_plan_options.archived_at IS NULL
AND meal_plan_options.belongs_to_meal_plan = $1
AND meal_plan_options.id = $2
AND meal_plans.archived_at IS NULL
AND meal_plans.id = $3
`

// GetMealPlanOption fetches a meal plan option from the database.
func (q *Querier) GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) (*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	args := []interface{}{
		mealPlanID,
		mealPlanOptionID,
		mealPlanID,
	}

	row := q.getOneRow(ctx, q.db, "mealPlanOption", getMealPlanOptionQuery, args)

	mealPlanOption, _, _, err := q.scanMealPlanOption(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning mealPlanOption")
	}

	return mealPlanOption, nil
}

// GetMealPlanOptions fetches a list of meal plan options from the database that meet a particular filter.
func (q *Querier) GetMealPlanOptions(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (x *types.MealPlanOptionList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	x = &types.MealPlanOptionList{}
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

	groupBys := []string{"meal_plan_options.id", "meals.id"}
	query, args := q.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, groupBys, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "mealPlanOptions", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan options list retrieval query")
	}

	if x.MealPlanOptions, x.FilteredCount, x.TotalCount, err = q.scanMealPlanOptions(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plan options")
	}

	return x, nil
}

const mealPlanOptionCreationQuery = "INSERT INTO meal_plan_options (id,day,assigned_cook,assigned_dishwasher,meal_name,meal_id,notes,belongs_to_meal_plan) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"

// createMealPlanOption creates a meal plan option in the database.
func (q *Querier) createMealPlanOption(ctx context.Context, db database.SQLQueryExecutor, input *types.MealPlanOptionDatabaseCreationInput) (*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanOptionIDKey, input.ID)

	// we're leaving PrepStepsCreated out on purpose since it would be false by default.
	args := []interface{}{
		input.ID,
		input.Day,
		input.AssignedCook,
		input.AssignedDishwasher,
		input.MealName,
		input.MealID,
		input.Notes,
		input.BelongsToMealPlan,
	}

	// create the meal plan option.
	if err := q.performWriteQuery(ctx, db, "meal plan option creation", mealPlanOptionCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating meal plan option")
	}

	x := &types.MealPlanOption{
		ID:                input.ID,
		Day:               input.Day,
		AssignedCook:      input.AssignedCook,
		Meal:              types.Meal{ID: input.MealID},
		MealName:          input.MealName,
		Notes:             input.Notes,
		BelongsToMealPlan: input.BelongsToMealPlan,
		CreatedAt:         q.currentTime(),
		Votes:             []*types.MealPlanOptionVote{},
	}

	tracing.AttachMealPlanOptionIDToSpan(span, x.ID)
	logger.Info("meal plan option created")

	return x, nil
}

// CreateMealPlanOption creates a meal plan option in the database.
func (q *Querier) CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionDatabaseCreationInput) (*types.MealPlanOption, error) {
	return q.createMealPlanOption(ctx, q.db, input)
}

const updateMealPlanOptionQuery = `UPDATE meal_plan_options
SET 
	day = $1,
	assigned_cook = $2,
	assigned_dishwasher = $3,
	meal_id = $4,
	meal_name = $5,
	notes = $6,
	prep_steps_created = $7,
	last_updated_at = extract(epoch FROM NOW())
WHERE archived_at IS NULL
  AND belongs_to_meal_plan = $8 
  AND id = $9
`

// UpdateMealPlanOption updates a particular meal plan option.
func (q *Querier) UpdateMealPlanOption(ctx context.Context, updated *types.MealPlanOption) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanOptionIDKey, updated.ID)
	tracing.AttachMealPlanOptionIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Day,
		updated.AssignedCook,
		updated.AssignedDishwasher,
		updated.Meal.ID,
		updated.MealName,
		updated.Notes,
		updated.PrepStepsCreated,
		updated.BelongsToMealPlan,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan option update", updateMealPlanOptionQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option")
	}

	logger.Info("meal plan option updated")

	return nil
}

const archiveMealPlanOptionQuery = "UPDATE meal_plan_options SET archived_at = extract(epoch FROM NOW()) WHERE archived_at IS NULL AND belongs_to_meal_plan = $1 AND id = $2"

// ArchiveMealPlanOption archives a meal plan option from the database by its ID.
func (q *Querier) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	args := []interface{}{
		mealPlanID,
		mealPlanOptionID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan option archive", archiveMealPlanOptionQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option")
	}

	logger.Info("meal plan option archived")

	return nil
}

func (q *Querier) determineWinner(winners []schulze.Score) string {
	var (
		highestScore int
		scoreWinners []string
	)

	for _, winner := range winners {
		if winner.Wins == highestScore {
			scoreWinners = append(scoreWinners, winner.Choice)
		} else if winner.Wins > highestScore {
			highestScore = winner.Wins
			scoreWinners = []string{winner.Choice}
		}
	}

	/* #nosec: G404 */
	return scoreWinners[rand.Intn(len(scoreWinners))]
}

func (q *Querier) decideOptionWinner(ctx context.Context, options []*types.MealPlanOption) (_ string, _, _ bool) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	candidateMap := map[string]struct{}{}
	votesByUser := map[string]schulze.Ballot{}

	logger := q.logger.WithValue("options.count", len(options))

	for _, option := range options {
		for _, v := range option.Votes {
			if votesByUser[v.ByUser] == nil {
				votesByUser[v.ByUser] = schulze.Ballot{}
			}

			if !v.Abstain {
				votesByUser[v.ByUser][v.BelongsToMealPlanOption] = int(v.Rank)
			}

			candidateMap[v.BelongsToMealPlanOption] = struct{}{}
		}
	}

	candidates := []string{}
	for c := range candidateMap {
		candidates = append(candidates, c)
	}

	e := schulze.NewVoting(candidates...)
	for _, vote := range votesByUser {
		if voteErr := e.Vote(vote); voteErr != nil {
			// this actually can never happen because we use uints for ranks, lol
			observability.AcknowledgeError(voteErr, logger, span, "an invalid vote was received")
		}
	}

	winners, tie := e.Compute()
	if tie {
		return q.determineWinner(winners), true, true
	}

	if len(winners) > 0 {
		return winners[0].Choice, false, true
	}

	return "", false, false
}

const finalizeMealPlanOptionQuery = `UPDATE meal_plan_options SET chosen = (belongs_to_meal_plan = $1 AND id = $2), tiebroken = $3 WHERE archived_at IS NULL AND belongs_to_meal_plan = $1 AND id = $2`

// FinalizeMealPlanOption archives a meal plan option vote from the database by its ID.
func (q *Querier) FinalizeMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID, householdID string) (changed bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanOptionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// fetch meal plan
	mealPlan, err := q.GetMealPlan(ctx, mealPlanID, householdID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "fetching meal plan")
	}

	// fetch meal plan option
	mealPlanOption, err := q.GetMealPlanOption(ctx, mealPlan.ID, mealPlanOptionID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "fetching meal plan option")
	}

	// fetch household data
	household, err := q.GetHouseholdByID(ctx, mealPlan.BelongsToHousehold)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "fetching household")
	}

	relevantOptions := byDayAndMeal(mealPlan.Options, mealPlanOption.Day, mealPlanOption.MealName)

	// go through all the votes for this meal plan option and determine if they're all there
	for _, member := range household.Members {
		for _, option := range relevantOptions {
			memberVoteFound := false
			for _, vote := range option.Votes {
				if vote.ByUser == member.BelongsToUser.ID {
					memberVoteFound = true
					break
				}
			}

			if !memberVoteFound {
				return false, nil
			}
		}
	}

	winner, tiebroken, chosen := q.decideOptionWinner(ctx, relevantOptions)
	if chosen {
		args := []interface{}{
			mealPlanID,
			winner,
			tiebroken,
		}

		if err = q.performWriteQuery(ctx, q.db, "meal plan option finalization", finalizeMealPlanOptionQuery, args); err != nil {
			return false, observability.PrepareError(err, logger, span, "finalizing meal plan option")
		}

		logger.Debug("finalized meal plan option")
	}

	return chosen, nil
}
