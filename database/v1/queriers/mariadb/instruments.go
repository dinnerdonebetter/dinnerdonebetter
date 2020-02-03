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
		"belongs_to",
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
		&x.BelongsTo,
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
func (m *MariaDB) buildGetInstrumentQuery(instrumentID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(instrumentsTableColumns...).
		From(instrumentsTableName).
		Where(squirrel.Eq{
			"id":         instrumentID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetInstrument fetches an instrument from the mariadb database
func (m *MariaDB) GetInstrument(ctx context.Context, instrumentID, userID uint64) (*models.Instrument, error) {
	query, args := m.buildGetInstrumentQuery(instrumentID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanInstrument(row)
}

// buildGetInstrumentCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of instruments belonging to a given user that meet a given query
func (m *MariaDB) buildGetInstrumentCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(CountQuery).
		From(instrumentsTableName).
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

// GetInstrumentCount will fetch the count of instruments from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetInstrumentCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetInstrumentCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allInstrumentsCountQueryBuilder sync.Once
	allInstrumentsCountQuery        string
)

// buildGetAllInstrumentsCountQuery returns a query that fetches the total number of instruments in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllInstrumentsCountQuery() string {
	allInstrumentsCountQueryBuilder.Do(func() {
		var err error
		allInstrumentsCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(instrumentsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allInstrumentsCountQuery
}

// GetAllInstrumentsCount will fetch the count of instruments from the database
func (m *MariaDB) GetAllInstrumentsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllInstrumentsCountQuery()).Scan(&count)
	return count, err
}

// buildGetInstrumentsQuery builds a SQL query selecting instruments that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetInstrumentsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(instrumentsTableColumns...).
		From(instrumentsTableName).
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

// GetInstruments fetches a list of instruments from the database that meet a particular filter
func (m *MariaDB) GetInstruments(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.InstrumentList, error) {
	query, args := m.buildGetInstrumentsQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for instruments")
	}

	list, err := scanInstruments(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetInstrumentCount(ctx, filter, userID)
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

// GetAllInstrumentsForUser fetches every instrument belonging to a user
func (m *MariaDB) GetAllInstrumentsForUser(ctx context.Context, userID uint64) ([]models.Instrument, error) {
	query, args := m.buildGetInstrumentsQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching instruments for user")
	}

	list, err := scanInstruments(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateInstrumentQuery takes an instrument and returns a creation query for that instrument and the relevant arguments.
func (m *MariaDB) buildCreateInstrumentQuery(input *models.Instrument) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(instrumentsTableName).
		Columns(
			"name",
			"variant",
			"description",
			"icon",
			"belongs_to",
			"created_on",
		).
		Values(
			input.Name,
			input.Variant,
			input.Description,
			input.Icon,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildInstrumentCreationTimeQuery takes an instrument and returns a creation query for that instrument and the relevant arguments
func (m *MariaDB) buildInstrumentCreationTimeQuery(instrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(instrumentsTableName).
		Where(squirrel.Eq{"id": instrumentID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateInstrument creates an instrument in the database
func (m *MariaDB) CreateInstrument(ctx context.Context, input *models.InstrumentCreationInput) (*models.Instrument, error) {
	x := &models.Instrument{
		Name:        input.Name,
		Variant:     input.Variant,
		Description: input.Description,
		Icon:        input.Icon,
		BelongsTo:   input.BelongsTo,
	}

	query, args := m.buildCreateInstrumentQuery(x)

	// create the instrument
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing instrument creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildInstrumentCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateInstrumentQuery takes an instrument and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdateInstrumentQuery(input *models.Instrument) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(instrumentsTableName).
		Set("name", input.Name).
		Set("variant", input.Variant).
		Set("description", input.Description).
		Set("icon", input.Icon).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateInstrument updates a particular instrument. Note that UpdateInstrument expects the provided input to have a valid ID.
func (m *MariaDB) UpdateInstrument(ctx context.Context, input *models.Instrument) error {
	query, args := m.buildUpdateInstrumentQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveInstrumentQuery returns a SQL query which marks a given instrument belonging to a given user as archived.
func (m *MariaDB) buildArchiveInstrumentQuery(instrumentID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(instrumentsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          instrumentID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveInstrument marks an instrument as archived in the database
func (m *MariaDB) ArchiveInstrument(ctx context.Context, instrumentID, userID uint64) error {
	query, args := m.buildArchiveInstrumentQuery(instrumentID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
