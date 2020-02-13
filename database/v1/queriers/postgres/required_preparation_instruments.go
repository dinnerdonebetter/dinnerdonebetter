package postgres

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
func (p *Postgres) buildGetRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(requiredPreparationInstrumentsTableColumns...).
		From(requiredPreparationInstrumentsTableName).
		Where(squirrel.Eq{
			"id":         requiredPreparationInstrumentID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRequiredPreparationInstrument fetches a required preparation instrument from the postgres database
func (p *Postgres) GetRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID, userID uint64) (*models.RequiredPreparationInstrument, error) {
	query, args := p.buildGetRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanRequiredPreparationInstrument(row)
}

// buildGetRequiredPreparationInstrumentCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of required preparation instruments belonging to a given user that meet a given query
func (p *Postgres) buildGetRequiredPreparationInstrumentCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(requiredPreparationInstrumentsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRequiredPreparationInstrumentCount will fetch the count of required preparation instruments from the database that meet a particular filter and belong to a particular user.
func (p *Postgres) GetRequiredPreparationInstrumentCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := p.buildGetRequiredPreparationInstrumentCountQuery(filter, userID)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allRequiredPreparationInstrumentsCountQueryBuilder sync.Once
	allRequiredPreparationInstrumentsCountQuery        string
)

// buildGetAllRequiredPreparationInstrumentsCountQuery returns a query that fetches the total number of required preparation instruments in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRequiredPreparationInstrumentsCountQuery() string {
	allRequiredPreparationInstrumentsCountQueryBuilder.Do(func() {
		var err error
		allRequiredPreparationInstrumentsCountQuery, _, err = p.sqlBuilder.
			Select(CountQuery).
			From(requiredPreparationInstrumentsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRequiredPreparationInstrumentsCountQuery
}

// GetAllRequiredPreparationInstrumentsCount will fetch the count of required preparation instruments from the database
func (p *Postgres) GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRequiredPreparationInstrumentsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRequiredPreparationInstrumentsQuery builds a SQL query selecting required preparation instruments that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRequiredPreparationInstrumentsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(requiredPreparationInstrumentsTableColumns...).
		From(requiredPreparationInstrumentsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRequiredPreparationInstruments fetches a list of required preparation instruments from the database that meet a particular filter
func (p *Postgres) GetRequiredPreparationInstruments(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RequiredPreparationInstrumentList, error) {
	query, args := p.buildGetRequiredPreparationInstrumentsQuery(filter, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for required preparation instruments")
	}

	list, err := scanRequiredPreparationInstruments(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetRequiredPreparationInstrumentCount(ctx, filter, userID)
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
func (p *Postgres) GetAllRequiredPreparationInstrumentsForUser(ctx context.Context, userID uint64) ([]models.RequiredPreparationInstrument, error) {
	query, args := p.buildGetRequiredPreparationInstrumentsQuery(nil, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching required preparation instruments for user")
	}

	list, err := scanRequiredPreparationInstruments(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRequiredPreparationInstrumentQuery takes a required preparation instrument and returns a creation query for that required preparation instrument and the relevant arguments.
func (p *Postgres) buildCreateRequiredPreparationInstrumentQuery(input *models.RequiredPreparationInstrument) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(requiredPreparationInstrumentsTableName).
		Columns(
			"instrument_id",
			"preparation_id",
			"notes",
		).
		Values(
			input.InstrumentID,
			input.PreparationID,
			input.Notes,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRequiredPreparationInstrument creates a required preparation instrument in the database
func (p *Postgres) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (*models.RequiredPreparationInstrument, error) {
	x := &models.RequiredPreparationInstrument{
		InstrumentID:  input.InstrumentID,
		PreparationID: input.PreparationID,
		Notes:         input.Notes,
	}

	query, args := p.buildCreateRequiredPreparationInstrumentQuery(x)

	// create the required preparation instrument
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing required preparation instrument creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRequiredPreparationInstrumentQuery takes a required preparation instrument and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdateRequiredPreparationInstrumentQuery(input *models.RequiredPreparationInstrument) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(requiredPreparationInstrumentsTableName).
		Set("instrument_id", input.InstrumentID).
		Set("preparation_id", input.PreparationID).
		Set("notes", input.Notes).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRequiredPreparationInstrument updates a particular required preparation instrument. Note that UpdateRequiredPreparationInstrument expects the provided input to have a valid ID.
func (p *Postgres) UpdateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrument) error {
	query, args := p.buildUpdateRequiredPreparationInstrumentQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveRequiredPreparationInstrumentQuery returns a SQL query which marks a given required preparation instrument belonging to a given user as archived.
func (p *Postgres) buildArchiveRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(requiredPreparationInstrumentsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          requiredPreparationInstrumentID,
			"archived_on": nil,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRequiredPreparationInstrument marks a required preparation instrument as archived in the database
func (p *Postgres) ArchiveRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID, userID uint64) error {
	query, args := p.buildArchiveRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
