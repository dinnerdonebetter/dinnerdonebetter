package mariadb

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
	recipeStepEventsTableName = "recipe_step_events"
)

var (
	recipeStepEventsTableColumns = []string{
		"id",
		"event_type",
		"done",
		"recipe_iteration_id",
		"recipe_step_id",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanRecipeStepEvent takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Event struct
func scanRecipeStepEvent(scan database.Scanner) (*models.RecipeStepEvent, error) {
	x := &models.RecipeStepEvent{}

	if err := scan.Scan(
		&x.ID,
		&x.EventType,
		&x.Done,
		&x.RecipeIterationID,
		&x.RecipeStepID,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeStepEvents takes a logger and some database rows and turns them into a slice of recipe step events
func scanRecipeStepEvents(logger logging.Logger, rows *sql.Rows) ([]models.RecipeStepEvent, error) {
	var list []models.RecipeStepEvent

	for rows.Next() {
		x, err := scanRecipeStepEvent(rows)
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

// buildGetRecipeStepEventQuery constructs a SQL query for fetching a recipe step event with a given ID belong to a user with a given ID.
func (m *MariaDB) buildGetRecipeStepEventQuery(recipeStepEventID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(recipeStepEventsTableColumns...).
		From(recipeStepEventsTableName).
		Where(squirrel.Eq{
			"id":         recipeStepEventID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepEvent fetches a recipe step event from the mariadb database
func (m *MariaDB) GetRecipeStepEvent(ctx context.Context, recipeStepEventID, userID uint64) (*models.RecipeStepEvent, error) {
	query, args := m.buildGetRecipeStepEventQuery(recipeStepEventID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanRecipeStepEvent(row)
}

// buildGetRecipeStepEventCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of recipe step events belonging to a given user that meet a given query
func (m *MariaDB) buildGetRecipeStepEventCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(CountQuery).
		From(recipeStepEventsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepEventCount will fetch the count of recipe step events from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetRecipeStepEventCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetRecipeStepEventCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allRecipeStepEventsCountQueryBuilder sync.Once
	allRecipeStepEventsCountQuery        string
)

// buildGetAllRecipeStepEventsCountQuery returns a query that fetches the total number of recipe step events in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllRecipeStepEventsCountQuery() string {
	allRecipeStepEventsCountQueryBuilder.Do(func() {
		var err error
		allRecipeStepEventsCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(recipeStepEventsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allRecipeStepEventsCountQuery
}

// GetAllRecipeStepEventsCount will fetch the count of recipe step events from the database
func (m *MariaDB) GetAllRecipeStepEventsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllRecipeStepEventsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeStepEventsQuery builds a SQL query selecting recipe step events that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetRecipeStepEventsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(recipeStepEventsTableColumns...).
		From(recipeStepEventsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepEvents fetches a list of recipe step events from the database that meet a particular filter
func (m *MariaDB) GetRecipeStepEvents(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepEventList, error) {
	query, args := m.buildGetRecipeStepEventsQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step events")
	}

	list, err := scanRecipeStepEvents(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetRecipeStepEventCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching recipe step event count: %w", err)
	}

	x := &models.RecipeStepEventList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeStepEvents: list,
	}

	return x, nil
}

// GetAllRecipeStepEventsForUser fetches every recipe step event belonging to a user
func (m *MariaDB) GetAllRecipeStepEventsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepEvent, error) {
	query, args := m.buildGetRecipeStepEventsQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching recipe step events for user")
	}

	list, err := scanRecipeStepEvents(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRecipeStepEventQuery takes a recipe step event and returns a creation query for that recipe step event and the relevant arguments.
func (m *MariaDB) buildCreateRecipeStepEventQuery(input *models.RecipeStepEvent) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(recipeStepEventsTableName).
		Columns(
			"event_type",
			"done",
			"recipe_iteration_id",
			"recipe_step_id",
			"belongs_to",
			"created_on",
		).
		Values(
			input.EventType,
			input.Done,
			input.RecipeIterationID,
			input.RecipeStepID,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildRecipeStepEventCreationTimeQuery takes a recipe step event and returns a creation query for that recipe step event and the relevant arguments
func (m *MariaDB) buildRecipeStepEventCreationTimeQuery(recipeStepEventID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(recipeStepEventsTableName).
		Where(squirrel.Eq{"id": recipeStepEventID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepEvent creates a recipe step event in the database
func (m *MariaDB) CreateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEventCreationInput) (*models.RecipeStepEvent, error) {
	x := &models.RecipeStepEvent{
		EventType:         input.EventType,
		Done:              input.Done,
		RecipeIterationID: input.RecipeIterationID,
		RecipeStepID:      input.RecipeStepID,
		BelongsTo:         input.BelongsTo,
	}

	query, args := m.buildCreateRecipeStepEventQuery(x)

	// create the recipe step event
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step event creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildRecipeStepEventCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateRecipeStepEventQuery takes a recipe step event and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdateRecipeStepEventQuery(input *models.RecipeStepEvent) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(recipeStepEventsTableName).
		Set("event_type", input.EventType).
		Set("done", input.Done).
		Set("recipe_iteration_id", input.RecipeIterationID).
		Set("recipe_step_id", input.RecipeStepID).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepEvent updates a particular recipe step event. Note that UpdateRecipeStepEvent expects the provided input to have a valid ID.
func (m *MariaDB) UpdateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEvent) error {
	query, args := m.buildUpdateRecipeStepEventQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveRecipeStepEventQuery returns a SQL query which marks a given recipe step event belonging to a given user as archived.
func (m *MariaDB) buildArchiveRecipeStepEventQuery(recipeStepEventID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(recipeStepEventsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeStepEventID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepEvent marks a recipe step event as archived in the database
func (m *MariaDB) ArchiveRecipeStepEvent(ctx context.Context, recipeStepEventID, userID uint64) error {
	query, args := m.buildArchiveRecipeStepEventQuery(recipeStepEventID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
