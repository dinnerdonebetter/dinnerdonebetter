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
	recipeStepInstrumentsTableName = "recipe_step_instruments"
)

var (
	recipeStepInstrumentsTableColumns = []string{
		"id",
		"instrument_id",
		"recipe_step_id",
		"notes",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanRecipeStepInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Instrument struct
func scanRecipeStepInstrument(scan database.Scanner) (*models.RecipeStepInstrument, error) {
	x := &models.RecipeStepInstrument{}

	if err := scan.Scan(
		&x.ID,
		&x.InstrumentID,
		&x.RecipeStepID,
		&x.Notes,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeStepInstruments takes a logger and some database rows and turns them into a slice of recipe step instruments
func scanRecipeStepInstruments(logger logging.Logger, rows *sql.Rows) ([]models.RecipeStepInstrument, error) {
	var list []models.RecipeStepInstrument

	for rows.Next() {
		x, err := scanRecipeStepInstrument(rows)
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

// buildGetRecipeStepInstrumentQuery constructs a SQL query for fetching a recipe step instrument with a given ID belong to a user with a given ID.
func (m *MariaDB) buildGetRecipeStepInstrumentQuery(recipeStepInstrumentID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(recipeStepInstrumentsTableColumns...).
		From(recipeStepInstrumentsTableName).
		Where(squirrel.Eq{
			"id":         recipeStepInstrumentID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepInstrument fetches a recipe step instrument from the mariadb database
func (m *MariaDB) GetRecipeStepInstrument(ctx context.Context, recipeStepInstrumentID, userID uint64) (*models.RecipeStepInstrument, error) {
	query, args := m.buildGetRecipeStepInstrumentQuery(recipeStepInstrumentID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanRecipeStepInstrument(row)
}

// buildGetRecipeStepInstrumentCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of recipe step instruments belonging to a given user that meet a given query
func (m *MariaDB) buildGetRecipeStepInstrumentCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(CountQuery).
		From(recipeStepInstrumentsTableName).
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

// GetRecipeStepInstrumentCount will fetch the count of recipe step instruments from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetRecipeStepInstrumentCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetRecipeStepInstrumentCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allRecipeStepInstrumentsCountQueryBuilder sync.Once
	allRecipeStepInstrumentsCountQuery        string
)

// buildGetAllRecipeStepInstrumentsCountQuery returns a query that fetches the total number of recipe step instruments in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllRecipeStepInstrumentsCountQuery() string {
	allRecipeStepInstrumentsCountQueryBuilder.Do(func() {
		var err error
		allRecipeStepInstrumentsCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(recipeStepInstrumentsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allRecipeStepInstrumentsCountQuery
}

// GetAllRecipeStepInstrumentsCount will fetch the count of recipe step instruments from the database
func (m *MariaDB) GetAllRecipeStepInstrumentsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllRecipeStepInstrumentsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeStepInstrumentsQuery builds a SQL query selecting recipe step instruments that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetRecipeStepInstrumentsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(recipeStepInstrumentsTableColumns...).
		From(recipeStepInstrumentsTableName).
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

// GetRecipeStepInstruments fetches a list of recipe step instruments from the database that meet a particular filter
func (m *MariaDB) GetRecipeStepInstruments(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepInstrumentList, error) {
	query, args := m.buildGetRecipeStepInstrumentsQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step instruments")
	}

	list, err := scanRecipeStepInstruments(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetRecipeStepInstrumentCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching recipe step instrument count: %w", err)
	}

	x := &models.RecipeStepInstrumentList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeStepInstruments: list,
	}

	return x, nil
}

// GetAllRecipeStepInstrumentsForUser fetches every recipe step instrument belonging to a user
func (m *MariaDB) GetAllRecipeStepInstrumentsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepInstrument, error) {
	query, args := m.buildGetRecipeStepInstrumentsQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching recipe step instruments for user")
	}

	list, err := scanRecipeStepInstruments(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRecipeStepInstrumentQuery takes a recipe step instrument and returns a creation query for that recipe step instrument and the relevant arguments.
func (m *MariaDB) buildCreateRecipeStepInstrumentQuery(input *models.RecipeStepInstrument) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(recipeStepInstrumentsTableName).
		Columns(
			"instrument_id",
			"recipe_step_id",
			"notes",
			"belongs_to",
			"created_on",
		).
		Values(
			input.InstrumentID,
			input.RecipeStepID,
			input.Notes,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildRecipeStepInstrumentCreationTimeQuery takes a recipe step instrument and returns a creation query for that recipe step instrument and the relevant arguments
func (m *MariaDB) buildRecipeStepInstrumentCreationTimeQuery(recipeStepInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(recipeStepInstrumentsTableName).
		Where(squirrel.Eq{"id": recipeStepInstrumentID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepInstrument creates a recipe step instrument in the database
func (m *MariaDB) CreateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrumentCreationInput) (*models.RecipeStepInstrument, error) {
	x := &models.RecipeStepInstrument{
		InstrumentID: input.InstrumentID,
		RecipeStepID: input.RecipeStepID,
		Notes:        input.Notes,
		BelongsTo:    input.BelongsTo,
	}

	query, args := m.buildCreateRecipeStepInstrumentQuery(x)

	// create the recipe step instrument
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step instrument creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildRecipeStepInstrumentCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateRecipeStepInstrumentQuery takes a recipe step instrument and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdateRecipeStepInstrumentQuery(input *models.RecipeStepInstrument) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(recipeStepInstrumentsTableName).
		Set("instrument_id", input.InstrumentID).
		Set("recipe_step_id", input.RecipeStepID).
		Set("notes", input.Notes).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepInstrument updates a particular recipe step instrument. Note that UpdateRecipeStepInstrument expects the provided input to have a valid ID.
func (m *MariaDB) UpdateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrument) error {
	query, args := m.buildUpdateRecipeStepInstrumentQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveRecipeStepInstrumentQuery returns a SQL query which marks a given recipe step instrument belonging to a given user as archived.
func (m *MariaDB) buildArchiveRecipeStepInstrumentQuery(recipeStepInstrumentID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(recipeStepInstrumentsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeStepInstrumentID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepInstrument marks a recipe step instrument as archived in the database
func (m *MariaDB) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepInstrumentID, userID uint64) error {
	query, args := m.buildArchiveRecipeStepInstrumentQuery(recipeStepInstrumentID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
