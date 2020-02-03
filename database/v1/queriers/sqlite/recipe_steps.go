package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	recipeStepsTableName = "recipe_steps"
)

var (
	recipeStepsTableColumns = []string{
		"id",
		"index",
		"preparation_id",
		"prerequisite_step",
		"min_estimated_time_in_seconds",
		"max_estimated_time_in_seconds",
		"temperature_in_celsius",
		"notes",
		"recipe_id",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanRecipeStep takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step struct
func scanRecipeStep(scan database.Scanner) (*models.RecipeStep, error) {
	x := &models.RecipeStep{}

	if err := scan.Scan(
		&x.ID,
		&x.Index,
		&x.PreparationID,
		&x.PrerequisiteStep,
		&x.MinEstimatedTimeInSeconds,
		&x.MaxEstimatedTimeInSeconds,
		&x.TemperatureInCelsius,
		&x.Notes,
		&x.RecipeID,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeSteps takes a logger and some database rows and turns them into a slice of recipe steps
func scanRecipeSteps(logger logging.Logger, rows *sql.Rows) ([]models.RecipeStep, error) {
	var list []models.RecipeStep

	for rows.Next() {
		x, err := scanRecipeStep(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildGetRecipeStepQuery constructs a SQL query for fetching a recipe step with a given ID belong to a user with a given ID.
func (s *Sqlite) buildGetRecipeStepQuery(recipeStepID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Select(recipeStepsTableColumns...).
		From(recipeStepsTableName).
		Where(squirrel.Eq{
			"id":         recipeStepID,
			"belongs_to": userID,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStep fetches a recipe step from the sqlite database
func (s *Sqlite) GetRecipeStep(ctx context.Context, recipeStepID, userID uint64) (*models.RecipeStep, error) {
	query, args := s.buildGetRecipeStepQuery(recipeStepID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)
	return scanRecipeStep(row)
}

// buildGetRecipeStepCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of recipe steps belonging to a given user that meet a given query
func (s *Sqlite) buildGetRecipeStepCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := s.sqlBuilder.
		Select(CountQuery).
		From(recipeStepsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepCount will fetch the count of recipe steps from the database that meet a particular filter and belong to a particular user.
func (s *Sqlite) GetRecipeStepCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := s.buildGetRecipeStepCountQuery(filter, userID)
	err = s.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allRecipeStepsCountQueryBuilder sync.Once
	allRecipeStepsCountQuery        string
)

// buildGetAllRecipeStepsCountQuery returns a query that fetches the total number of recipe steps in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (s *Sqlite) buildGetAllRecipeStepsCountQuery() string {
	allRecipeStepsCountQueryBuilder.Do(func() {
		var err error
		allRecipeStepsCountQuery, _, err = s.sqlBuilder.
			Select(CountQuery).
			From(recipeStepsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		s.logQueryBuildingError(err)
	})

	return allRecipeStepsCountQuery
}

// GetAllRecipeStepsCount will fetch the count of recipe steps from the database
func (s *Sqlite) GetAllRecipeStepsCount(ctx context.Context) (count uint64, err error) {
	err = s.db.QueryRowContext(ctx, s.buildGetAllRecipeStepsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeStepsQuery builds a SQL query selecting recipe steps that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (s *Sqlite) buildGetRecipeStepsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := s.sqlBuilder.
		Select(recipeStepsTableColumns...).
		From(recipeStepsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter
func (s *Sqlite) GetRecipeSteps(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepList, error) {
	query, args := s.buildGetRecipeStepsQuery(filter, userID)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe steps")
	}

	list, err := scanRecipeSteps(s.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := s.GetRecipeStepCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching recipe step count: %w", err)
	}

	x := &models.RecipeStepList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeSteps: list,
	}

	return x, nil
}

// GetAllRecipeStepsForUser fetches every recipe step belonging to a user
func (s *Sqlite) GetAllRecipeStepsForUser(ctx context.Context, userID uint64) ([]models.RecipeStep, error) {
	query, args := s.buildGetRecipeStepsQuery(nil, userID)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching recipe steps for user")
	}

	list, err := scanRecipeSteps(s.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRecipeStepQuery takes a recipe step and returns a creation query for that recipe step and the relevant arguments.
func (s *Sqlite) buildCreateRecipeStepQuery(input *models.RecipeStep) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Insert(recipeStepsTableName).
		Columns(
			"index",
			"preparation_id",
			"prerequisite_step",
			"min_estimated_time_in_seconds",
			"max_estimated_time_in_seconds",
			"temperature_in_celsius",
			"notes",
			"recipe_id",
			"belongs_to",
		).
		Values(
			input.Index,
			input.PreparationID,
			input.PrerequisiteStep,
			input.MinEstimatedTimeInSeconds,
			input.MaxEstimatedTimeInSeconds,
			input.TemperatureInCelsius,
			input.Notes,
			input.RecipeID,
			input.BelongsTo,
		).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// buildRecipeStepCreationTimeQuery takes a recipe step and returns a creation query for that recipe step and the relevant arguments
func (s *Sqlite) buildRecipeStepCreationTimeQuery(recipeStepID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select("created_on").
		From(recipeStepsTableName).
		Where(squirrel.Eq{"id": recipeStepID}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStep creates a recipe step in the database
func (s *Sqlite) CreateRecipeStep(ctx context.Context, input *models.RecipeStepCreationInput) (*models.RecipeStep, error) {
	x := &models.RecipeStep{
		Index:                     input.Index,
		PreparationID:             input.PreparationID,
		PrerequisiteStep:          input.PrerequisiteStep,
		MinEstimatedTimeInSeconds: input.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: input.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      input.TemperatureInCelsius,
		Notes:                     input.Notes,
		RecipeID:                  input.RecipeID,
		BelongsTo:                 input.BelongsTo,
	}

	query, args := s.buildCreateRecipeStepQuery(x)

	// create the recipe step
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := s.buildRecipeStepCreationTimeQuery(x.ID)
		s.logCreationTimeRetrievalError(s.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateRecipeStepQuery takes a recipe step and returns an update SQL query, with the relevant query parameters
func (s *Sqlite) buildUpdateRecipeStepQuery(input *models.RecipeStep) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Update(recipeStepsTableName).
		Set("index", input.Index).
		Set("preparation_id", input.PreparationID).
		Set("prerequisite_step", input.PrerequisiteStep).
		Set("min_estimated_time_in_seconds", input.MinEstimatedTimeInSeconds).
		Set("max_estimated_time_in_seconds", input.MaxEstimatedTimeInSeconds).
		Set("temperature_in_celsius", input.TemperatureInCelsius).
		Set("notes", input.Notes).
		Set("recipe_id", input.RecipeID).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStep updates a particular recipe step. Note that UpdateRecipeStep expects the provided input to have a valid ID.
func (s *Sqlite) UpdateRecipeStep(ctx context.Context, input *models.RecipeStep) error {
	query, args := s.buildUpdateRecipeStepQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveRecipeStepQuery returns a SQL query which marks a given recipe step belonging to a given user as archived.
func (s *Sqlite) buildArchiveRecipeStepQuery(recipeStepID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Update(recipeStepsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeStepID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStep marks a recipe step as archived in the database
func (s *Sqlite) ArchiveRecipeStep(ctx context.Context, recipeStepID, userID uint64) error {
	query, args := s.buildArchiveRecipeStepQuery(recipeStepID, userID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
