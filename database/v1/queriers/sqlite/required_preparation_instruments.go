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
	requiredPreparationInstrumentsTableName = "required_preparation_instruments"
)

var (
	requiredPreparationInstrumentsTableColumns = []string{
		"id",
		"instrument_id",
		"preparation_id",
		"notes",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanRequiredPreparationInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a Required Preparation Instrument struct
func scanRequiredPreparationInstrument(scan database.Scanner) (*models.RequiredPreparationInstrument, error) {
	x := &models.RequiredPreparationInstrument{}

	if err := scan.Scan(
		&x.ID,
		&x.InstrumentID,
		&x.PreparationID,
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

// scanRequiredPreparationInstruments takes a logger and some database rows and turns them into a slice of required preparation instruments
func scanRequiredPreparationInstruments(logger logging.Logger, rows *sql.Rows) ([]models.RequiredPreparationInstrument, error) {
	var list []models.RequiredPreparationInstrument

	for rows.Next() {
		x, err := scanRequiredPreparationInstrument(rows)
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

// buildGetRequiredPreparationInstrumentQuery constructs a SQL query for fetching a required preparation instrument with a given ID belong to a user with a given ID.
func (s *Sqlite) buildGetRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Select(requiredPreparationInstrumentsTableColumns...).
		From(requiredPreparationInstrumentsTableName).
		Where(squirrel.Eq{
			"id":         requiredPreparationInstrumentID,
			"belongs_to": userID,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetRequiredPreparationInstrument fetches a required preparation instrument from the sqlite database
func (s *Sqlite) GetRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID, userID uint64) (*models.RequiredPreparationInstrument, error) {
	query, args := s.buildGetRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)
	return scanRequiredPreparationInstrument(row)
}

// buildGetRequiredPreparationInstrumentCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of required preparation instruments belonging to a given user that meet a given query
func (s *Sqlite) buildGetRequiredPreparationInstrumentCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := s.sqlBuilder.
		Select(CountQuery).
		From(requiredPreparationInstrumentsTableName).
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

// GetRequiredPreparationInstrumentCount will fetch the count of required preparation instruments from the database that meet a particular filter and belong to a particular user.
func (s *Sqlite) GetRequiredPreparationInstrumentCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := s.buildGetRequiredPreparationInstrumentCountQuery(filter, userID)
	err = s.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allRequiredPreparationInstrumentsCountQueryBuilder sync.Once
	allRequiredPreparationInstrumentsCountQuery        string
)

// buildGetAllRequiredPreparationInstrumentsCountQuery returns a query that fetches the total number of required preparation instruments in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (s *Sqlite) buildGetAllRequiredPreparationInstrumentsCountQuery() string {
	allRequiredPreparationInstrumentsCountQueryBuilder.Do(func() {
		var err error
		allRequiredPreparationInstrumentsCountQuery, _, err = s.sqlBuilder.
			Select(CountQuery).
			From(requiredPreparationInstrumentsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		s.logQueryBuildingError(err)
	})

	return allRequiredPreparationInstrumentsCountQuery
}

// GetAllRequiredPreparationInstrumentsCount will fetch the count of required preparation instruments from the database
func (s *Sqlite) GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (count uint64, err error) {
	err = s.db.QueryRowContext(ctx, s.buildGetAllRequiredPreparationInstrumentsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRequiredPreparationInstrumentsQuery builds a SQL query selecting required preparation instruments that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (s *Sqlite) buildGetRequiredPreparationInstrumentsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := s.sqlBuilder.
		Select(requiredPreparationInstrumentsTableColumns...).
		From(requiredPreparationInstrumentsTableName).
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

// GetRequiredPreparationInstruments fetches a list of required preparation instruments from the database that meet a particular filter
func (s *Sqlite) GetRequiredPreparationInstruments(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RequiredPreparationInstrumentList, error) {
	query, args := s.buildGetRequiredPreparationInstrumentsQuery(filter, userID)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for required preparation instruments")
	}

	list, err := scanRequiredPreparationInstruments(s.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := s.GetRequiredPreparationInstrumentCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching required preparation instrument count: %w", err)
	}

	x := &models.RequiredPreparationInstrumentList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RequiredPreparationInstruments: list,
	}

	return x, nil
}

// GetAllRequiredPreparationInstrumentsForUser fetches every required preparation instrument belonging to a user
func (s *Sqlite) GetAllRequiredPreparationInstrumentsForUser(ctx context.Context, userID uint64) ([]models.RequiredPreparationInstrument, error) {
	query, args := s.buildGetRequiredPreparationInstrumentsQuery(nil, userID)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching required preparation instruments for user")
	}

	list, err := scanRequiredPreparationInstruments(s.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRequiredPreparationInstrumentQuery takes a required preparation instrument and returns a creation query for that required preparation instrument and the relevant arguments.
func (s *Sqlite) buildCreateRequiredPreparationInstrumentQuery(input *models.RequiredPreparationInstrument) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Insert(requiredPreparationInstrumentsTableName).
		Columns(
			"instrument_id",
			"preparation_id",
			"notes",
			"belongs_to",
		).
		Values(
			input.InstrumentID,
			input.PreparationID,
			input.Notes,
			input.BelongsTo,
		).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// buildRequiredPreparationInstrumentCreationTimeQuery takes a required preparation instrument and returns a creation query for that required preparation instrument and the relevant arguments
func (s *Sqlite) buildRequiredPreparationInstrumentCreationTimeQuery(requiredPreparationInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select("created_on").
		From(requiredPreparationInstrumentsTableName).
		Where(squirrel.Eq{"id": requiredPreparationInstrumentID}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// CreateRequiredPreparationInstrument creates a required preparation instrument in the database
func (s *Sqlite) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (*models.RequiredPreparationInstrument, error) {
	x := &models.RequiredPreparationInstrument{
		InstrumentID:  input.InstrumentID,
		PreparationID: input.PreparationID,
		Notes:         input.Notes,
		BelongsTo:     input.BelongsTo,
	}

	query, args := s.buildCreateRequiredPreparationInstrumentQuery(x)

	// create the required preparation instrument
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing required preparation instrument creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := s.buildRequiredPreparationInstrumentCreationTimeQuery(x.ID)
		s.logCreationTimeRetrievalError(s.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateRequiredPreparationInstrumentQuery takes a required preparation instrument and returns an update SQL query, with the relevant query parameters
func (s *Sqlite) buildUpdateRequiredPreparationInstrumentQuery(input *models.RequiredPreparationInstrument) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Update(requiredPreparationInstrumentsTableName).
		Set("instrument_id", input.InstrumentID).
		Set("preparation_id", input.PreparationID).
		Set("notes", input.Notes).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// UpdateRequiredPreparationInstrument updates a particular required preparation instrument. Note that UpdateRequiredPreparationInstrument expects the provided input to have a valid ID.
func (s *Sqlite) UpdateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrument) error {
	query, args := s.buildUpdateRequiredPreparationInstrumentQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveRequiredPreparationInstrumentQuery returns a SQL query which marks a given required preparation instrument belonging to a given user as archived.
func (s *Sqlite) buildArchiveRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Update(requiredPreparationInstrumentsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          requiredPreparationInstrumentID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// ArchiveRequiredPreparationInstrument marks a required preparation instrument as archived in the database
func (s *Sqlite) ArchiveRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID, userID uint64) error {
	query, args := s.buildArchiveRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID, userID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
