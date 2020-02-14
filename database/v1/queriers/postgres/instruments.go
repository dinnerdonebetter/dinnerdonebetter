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
	instrumentsTableName = "instruments"
)

var (
	instrumentsTableColumns = []string{
		"id",
		"name",
		"variant",
		"description",
		"icon",
		"created_on",
		"updated_on",
		"archived_on",
	}
)

// scanInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into an Instrument struct
func scanInstrument(scan database.Scanner) (*models.Instrument, error) {
	x := &models.Instrument{}

	if err := scan.Scan(
		&x.ID,
		&x.Name,
		&x.Variant,
		&x.Description,
		&x.Icon,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanInstruments takes a logger and some database rows and turns them into a slice of instruments
func scanInstruments(logger logging.Logger, rows *sql.Rows) ([]models.Instrument, error) {
	var list []models.Instrument

	for rows.Next() {
		x, err := scanInstrument(rows)
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

// buildGetInstrumentQuery constructs a SQL query for fetching an instrument with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetInstrumentQuery(instrumentID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(instrumentsTableColumns...).
		From(instrumentsTableName).
		Where(squirrel.Eq{
			"id": instrumentID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetInstrument fetches an instrument from the postgres database
func (p *Postgres) GetInstrument(ctx context.Context, instrumentID uint64) (*models.Instrument, error) {
	query, args := p.buildGetInstrumentQuery(instrumentID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanInstrument(row)
}

// buildGetInstrumentCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of instruments belonging to a given user that meet a given query
func (p *Postgres) buildGetInstrumentCountQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(instrumentsTableName).
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

// GetInstrumentCount will fetch the count of instruments from the database that meet a particular filter and belong to a particular user.
func (p *Postgres) GetInstrumentCount(ctx context.Context, filter *models.QueryFilter) (count uint64, err error) {
	query, args := p.buildGetInstrumentCountQuery(filter)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allInstrumentsCountQueryBuilder sync.Once
	allInstrumentsCountQuery        string
)

// buildGetAllInstrumentsCountQuery returns a query that fetches the total number of instruments in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllInstrumentsCountQuery() string {
	allInstrumentsCountQueryBuilder.Do(func() {
		var err error
		allInstrumentsCountQuery, _, err = p.sqlBuilder.
			Select(CountQuery).
			From(instrumentsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allInstrumentsCountQuery
}

// GetAllInstrumentsCount will fetch the count of instruments from the database
func (p *Postgres) GetAllInstrumentsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllInstrumentsCountQuery()).Scan(&count)
	return count, err
}

// buildGetInstrumentsQuery builds a SQL query selecting instruments that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetInstrumentsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(instrumentsTableColumns...).
		From(instrumentsTableName).
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

// GetInstruments fetches a list of instruments from the database that meet a particular filter
func (p *Postgres) GetInstruments(ctx context.Context, filter *models.QueryFilter) (*models.InstrumentList, error) {
	query, args := p.buildGetInstrumentsQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for instruments")
	}

	list, err := scanInstruments(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetInstrumentCount(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("fetching instrument count: %w", err)
	}

	x := &models.InstrumentList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Instruments: list,
	}

	return x, nil
}

// buildCreateInstrumentQuery takes an instrument and returns a creation query for that instrument and the relevant arguments.
func (p *Postgres) buildCreateInstrumentQuery(input *models.Instrument) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(instrumentsTableName).
		Columns(
			"name",
			"variant",
			"description",
			"icon",
		).
		Values(
			input.Name,
			input.Variant,
			input.Description,
			input.Icon,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateInstrument creates an instrument in the database
func (p *Postgres) CreateInstrument(ctx context.Context, input *models.InstrumentCreationInput) (*models.Instrument, error) {
	x := &models.Instrument{
		Name:        input.Name,
		Variant:     input.Variant,
		Description: input.Description,
		Icon:        input.Icon,
	}

	query, args := p.buildCreateInstrumentQuery(x)

	// create the instrument
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing instrument creation query: %w", err)
	}

	return x, nil
}

// buildUpdateInstrumentQuery takes an instrument and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdateInstrumentQuery(input *models.Instrument) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(instrumentsTableName).
		Set("name", input.Name).
		Set("variant", input.Variant).
		Set("description", input.Description).
		Set("icon", input.Icon).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateInstrument updates a particular instrument. Note that UpdateInstrument expects the provided input to have a valid ID.
func (p *Postgres) UpdateInstrument(ctx context.Context, input *models.Instrument) error {
	query, args := p.buildUpdateInstrumentQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveInstrumentQuery returns a SQL query which marks a given instrument belonging to a given user as archived.
func (p *Postgres) buildArchiveInstrumentQuery(instrumentID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(instrumentsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          instrumentID,
			"archived_on": nil,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveInstrument marks an instrument as archived in the database
func (p *Postgres) ArchiveInstrument(ctx context.Context, instrumentID uint64) error {
	query, args := p.buildArchiveInstrumentQuery(instrumentID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
