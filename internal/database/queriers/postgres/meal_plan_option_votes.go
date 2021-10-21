package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var (
	_ types.MealPlanOptionVoteDataManager = (*SQLQuerier)(nil)

	// mealPlanOptionVotesTableColumns are the columns for the meal_plan_option_votes table.
	mealPlanOptionVotesTableColumns = []string{
		"meal_plan_option_votes.id",
		"meal_plan_option_votes.meal_plan_option_id",
		"meal_plan_option_votes.day_of_week",
		"meal_plan_option_votes.points",
		"meal_plan_option_votes.abstain",
		"meal_plan_option_votes.notes",
		"meal_plan_option_votes.created_on",
		"meal_plan_option_votes.last_updated_on",
		"meal_plan_option_votes.archived_on",
		"meal_plan_option_votes.belongs_to_account",
	}
)

// scanMealPlanOptionVote takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan option vote struct.
func (q *SQLQuerier) scanMealPlanOptionVote(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.MealPlanOptionVote, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.MealPlanOptionVote{}

	targetVars := []interface{}{
		&x.ID,
		&x.MealPlanOptionID,
		&x.DayOfWeek,
		&x.Points,
		&x.Abstain,
		&x.Notes,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToAccount,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanMealPlanOptionVotes takes some database rows and turns them into a slice of meal plan option votes.
func (q *SQLQuerier) scanMealPlanOptionVotes(ctx context.Context, rows database.ResultIterator, includeCounts bool) (mealPlanOptionVotes []*types.MealPlanOptionVote, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanMealPlanOptionVote(ctx, rows, includeCounts)
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

		mealPlanOptionVotes = append(mealPlanOptionVotes, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return mealPlanOptionVotes, filteredCount, totalCount, nil
}

const mealPlanOptionVoteExistenceQuery = "SELECT EXISTS ( SELECT meal_plan_option_votes.id FROM meal_plan_option_votes WHERE meal_plan_option_votes.archived_on IS NULL AND meal_plan_option_votes.id = $1 )"

// MealPlanOptionVoteExists fetches whether a meal plan option vote exists from the database.
func (q *SQLQuerier) MealPlanOptionVoteExists(ctx context.Context, mealPlanOptionVoteID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if mealPlanOptionVoteID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)

	args := []interface{}{
		mealPlanOptionVoteID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanOptionVoteExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing meal plan option vote existence check")
	}

	return result, nil
}

const getMealPlanOptionVoteQuery = "SELECT meal_plan_option_votes.id, meal_plan_option_votes.meal_plan_option_id, meal_plan_option_votes.day_of_week, meal_plan_option_votes.points, meal_plan_option_votes.abstain, meal_plan_option_votes.notes, meal_plan_option_votes.created_on, meal_plan_option_votes.last_updated_on, meal_plan_option_votes.archived_on, meal_plan_option_votes.belongs_to_account FROM meal_plan_option_votes WHERE meal_plan_option_votes.archived_on IS NULL AND meal_plan_option_votes.id = $1"

// GetMealPlanOptionVote fetches a meal plan option vote from the database.
func (q *SQLQuerier) GetMealPlanOptionVote(ctx context.Context, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if mealPlanOptionVoteID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)

	args := []interface{}{
		mealPlanOptionVoteID,
	}

	row := q.getOneRow(ctx, q.db, "mealPlanOptionVote", getMealPlanOptionVoteQuery, args)

	mealPlanOptionVote, _, _, err := q.scanMealPlanOptionVote(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning mealPlanOptionVote")
	}

	return mealPlanOptionVote, nil
}

const getTotalMealPlanOptionVotesCountQuery = "SELECT COUNT(meal_plan_option_votes.id) FROM meal_plan_option_votes WHERE meal_plan_option_votes.archived_on IS NULL"

// GetTotalMealPlanOptionVoteCount fetches the count of meal plan option votes from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalMealPlanOptionVoteCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, getTotalMealPlanOptionVotesCountQuery, "fetching count of meal plan option votes")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of meal plan option votes")
	}

	return count, nil
}

// GetMealPlanOptionVotes fetches a list of meal plan option votes from the database that meet a particular filter.
func (q *SQLQuerier) GetMealPlanOptionVotes(ctx context.Context, filter *types.QueryFilter) (x *types.MealPlanOptionVoteList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.MealPlanOptionVoteList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(
		ctx,
		"meal_plan_option_votes",
		nil,
		nil,
		accountOwnershipColumn,
		mealPlanOptionVotesTableColumns,
		"",
		false,
		filter,
	)

	rows, err := q.performReadQuery(ctx, q.db, "mealPlanOptionVotes", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing meal plan option votes list retrieval query")
	}

	if x.MealPlanOptionVotes, x.FilteredCount, x.TotalCount, err = q.scanMealPlanOptionVotes(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plan option votes")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetMealPlanOptionVotesWithIDsQuery(ctx context.Context, accountID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"meal_plan_option_votes.id":          ids,
		"meal_plan_option_votes.archived_on": nil,
	}

	if accountID != "" {
		withIDsWhere["meal_plan_option_votes.belongs_to_account"] = accountID
	}

	subqueryBuilder := q.sqlBuilder.Select(mealPlanOptionVotesTableColumns...).
		From("meal_plan_option_votes").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(mealPlanOptionVotesTableColumns...).
		FromSelect(subqueryBuilder, "meal_plan_option_votes").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetMealPlanOptionVotesWithIDs fetches meal plan option votes from the database within a given set of IDs.
func (q *SQLQuerier) GetMealPlanOptionVotesWithIDs(ctx context.Context, accountID string, limit uint8, ids []string) ([]*types.MealPlanOptionVote, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if ids == nil {
		return nil, ErrNilInputProvided
	}

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.buildGetMealPlanOptionVotesWithIDsQuery(ctx, accountID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "meal plan option votes with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching meal plan option votes from database")
	}

	mealPlanOptionVotes, _, _, err := q.scanMealPlanOptionVotes(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning meal plan option votes")
	}

	return mealPlanOptionVotes, nil
}

const mealPlanOptionVoteCreationQuery = "INSERT INTO meal_plan_option_votes (id,meal_plan_option_id,day_of_week,points,abstain,notes,belongs_to_account) VALUES ($1,$2,$3,$4,$5,$6,$7)"

// CreateMealPlanOptionVote creates a meal plan option vote in the database.
func (q *SQLQuerier) CreateMealPlanOptionVote(ctx context.Context, input *types.MealPlanOptionVoteDatabaseCreationInput) (*types.MealPlanOptionVote, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanOptionVoteIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.MealPlanOptionID,
		input.DayOfWeek,
		input.Points,
		input.Abstain,
		input.Notes,
		input.BelongsToAccount,
	}

	// create the meal plan option vote.
	if err := q.performWriteQuery(ctx, q.db, "meal plan option vote creation", mealPlanOptionVoteCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating meal plan option vote")
	}

	x := &types.MealPlanOptionVote{
		ID:               input.ID,
		MealPlanOptionID: input.MealPlanOptionID,
		DayOfWeek:        input.DayOfWeek,
		Points:           input.Points,
		Abstain:          input.Abstain,
		Notes:            input.Notes,
		BelongsToAccount: input.BelongsToAccount,
		CreatedOn:        q.currentTime(),
	}

	tracing.AttachMealPlanOptionVoteIDToSpan(span, x.ID)
	logger.Info("meal plan option vote created")

	return x, nil
}

const updateMealPlanOptionVoteQuery = "UPDATE meal_plan_option_votes SET meal_plan_option_id = $1, day_of_week = $2, points = $3, abstain = $4, notes = $5, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_account = $6 AND id = $7"

// UpdateMealPlanOptionVote updates a particular meal plan option vote.
func (q *SQLQuerier) UpdateMealPlanOptionVote(ctx context.Context, updated *types.MealPlanOptionVote) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanOptionVoteIDKey, updated.ID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, updated.ID)
	tracing.AttachAccountIDToSpan(span, updated.BelongsToAccount)

	args := []interface{}{
		updated.MealPlanOptionID,
		updated.DayOfWeek,
		updated.Points,
		updated.Abstain,
		updated.Notes,
		updated.BelongsToAccount,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan option vote update", updateMealPlanOptionVoteQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option vote")
	}

	logger.Info("meal plan option vote updated")

	return nil
}

const archiveMealPlanOptionVoteQuery = "UPDATE meal_plan_option_votes SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_account = $1 AND id = $2"

// ArchiveMealPlanOptionVote archives a meal plan option vote from the database by its ID.
func (q *SQLQuerier) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanOptionVoteID, accountID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if mealPlanOptionVoteID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)

	if accountID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachAccountIDToSpan(span, accountID)

	args := []interface{}{
		accountID,
		mealPlanOptionVoteID,
	}

	if err := q.performWriteQuery(ctx, q.db, "meal plan option vote archive", archiveMealPlanOptionVoteQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating meal plan option vote")
	}

	logger.Info("meal plan option vote archived")

	return nil
}
