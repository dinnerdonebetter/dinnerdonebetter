package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/Masterminds/squirrel"
)

const (
	requiredPreparationInstrumentsTableName            = "required_preparation_instruments"
	requiredPreparationInstrumentsTableOwnershipColumn = "belongs_to_valid_preparation"
)

var (
	requiredPreparationInstrumentsTableColumns = []string{
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, "id"),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, "valid_instrument_id"),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, "notes"),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, "created_on"),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, requiredPreparationInstrumentsTableOwnershipColumn),
	}
)

// scanRequiredPreparationInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a Required Preparation Instrument struct
func (p *Postgres) scanRequiredPreparationInstrument(scan database.Scanner, includeCount bool) (*models.RequiredPreparationInstrument, uint64, error) {
	x := &models.RequiredPreparationInstrument{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.ValidInstrumentID,
		&x.Notes,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToValidPreparation,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}

// scanRequiredPreparationInstruments takes a logger and some database rows and turns them into a slice of required preparation instruments.
func (p *Postgres) scanRequiredPreparationInstruments(rows database.ResultIterator) ([]models.RequiredPreparationInstrument, uint64, error) {
	var (
		list  []models.RequiredPreparationInstrument
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanRequiredPreparationInstrument(rows, true)
		if err != nil {
			return nil, 0, err
		}

		if count == 0 {
			count = c
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		p.logger.Error(closeErr, "closing database rows")
	}

	return list, count, nil
}

// buildRequiredPreparationInstrumentExistsQuery constructs a SQL query for checking if a required preparation instrument with a given ID belong to a a valid preparation with a given ID exists
func (p *Postgres) buildRequiredPreparationInstrumentExistsQuery(validPreparationID, requiredPreparationInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", requiredPreparationInstrumentsTableName)).
		Prefix(existencePrefix).
		From(requiredPreparationInstrumentsTableName).
		Join(validPreparationsOnRequiredPreparationInstrumentsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", requiredPreparationInstrumentsTableName):                                                     requiredPreparationInstrumentID,
			fmt.Sprintf("%s.id", validPreparationsTableName):                                                                  validPreparationID,
			fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, requiredPreparationInstrumentsTableOwnershipColumn): validPreparationID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RequiredPreparationInstrumentExists queries the database to see if a given required preparation instrument belonging to a given user exists.
func (p *Postgres) RequiredPreparationInstrumentExists(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) (exists bool, err error) {
	query, args := p.buildRequiredPreparationInstrumentExistsQuery(validPreparationID, requiredPreparationInstrumentID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRequiredPreparationInstrumentQuery constructs a SQL query for fetching a required preparation instrument with a given ID belong to a valid preparation with a given ID.
func (p *Postgres) buildGetRequiredPreparationInstrumentQuery(validPreparationID, requiredPreparationInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(requiredPreparationInstrumentsTableColumns...).
		From(requiredPreparationInstrumentsTableName).
		Join(validPreparationsOnRequiredPreparationInstrumentsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", requiredPreparationInstrumentsTableName):                                                     requiredPreparationInstrumentID,
			fmt.Sprintf("%s.id", validPreparationsTableName):                                                                  validPreparationID,
			fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, requiredPreparationInstrumentsTableOwnershipColumn): validPreparationID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRequiredPreparationInstrument fetches a required preparation instrument from the database.
func (p *Postgres) GetRequiredPreparationInstrument(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) (*models.RequiredPreparationInstrument, error) {
	query, args := p.buildGetRequiredPreparationInstrumentQuery(validPreparationID, requiredPreparationInstrumentID)
	row := p.db.QueryRowContext(ctx, query, args...)

	requiredPreparationInstrument, _, err := p.scanRequiredPreparationInstrument(row, false)
	return requiredPreparationInstrument, err
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
			Select(fmt.Sprintf(countQuery, requiredPreparationInstrumentsTableName)).
			From(requiredPreparationInstrumentsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", requiredPreparationInstrumentsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRequiredPreparationInstrumentsCountQuery
}

// GetAllRequiredPreparationInstrumentsCount will fetch the count of required preparation instruments from the database.
func (p *Postgres) GetAllRequiredPreparationInstrumentsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRequiredPreparationInstrumentsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRequiredPreparationInstrumentsQuery builds a SQL query selecting required preparation instruments that adhere to a given QueryFilter and belong to a given valid preparation,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRequiredPreparationInstrumentsQuery(validPreparationID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(requiredPreparationInstrumentsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllRequiredPreparationInstrumentsCountQuery()))...).
		From(requiredPreparationInstrumentsTableName).
		Join(validPreparationsOnRequiredPreparationInstrumentsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", requiredPreparationInstrumentsTableName):                                            nil,
			fmt.Sprintf("%s.id", validPreparationsTableName):                                                                  validPreparationID,
			fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, requiredPreparationInstrumentsTableOwnershipColumn): validPreparationID,
		}).
		OrderBy(fmt.Sprintf("%s.id", requiredPreparationInstrumentsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, requiredPreparationInstrumentsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRequiredPreparationInstruments fetches a list of required preparation instruments from the database that meet a particular filter.
func (p *Postgres) GetRequiredPreparationInstruments(ctx context.Context, validPreparationID uint64, filter *models.QueryFilter) (*models.RequiredPreparationInstrumentList, error) {
	query, args := p.buildGetRequiredPreparationInstrumentsQuery(validPreparationID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for required preparation instruments")
	}

	requiredPreparationInstruments, count, err := p.scanRequiredPreparationInstruments(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RequiredPreparationInstrumentList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RequiredPreparationInstruments: requiredPreparationInstruments,
	}

	return list, nil
}

// buildCreateRequiredPreparationInstrumentQuery takes a required preparation instrument and returns a creation query for that required preparation instrument and the relevant arguments.
func (p *Postgres) buildCreateRequiredPreparationInstrumentQuery(input *models.RequiredPreparationInstrument) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(requiredPreparationInstrumentsTableName).
		Columns(
			"valid_instrument_id",
			"notes",
			requiredPreparationInstrumentsTableOwnershipColumn,
		).
		Values(
			input.ValidInstrumentID,
			input.Notes,
			input.BelongsToValidPreparation,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRequiredPreparationInstrument creates a required preparation instrument in the database.
func (p *Postgres) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (*models.RequiredPreparationInstrument, error) {
	x := &models.RequiredPreparationInstrument{
		ValidInstrumentID:         input.ValidInstrumentID,
		Notes:                     input.Notes,
		BelongsToValidPreparation: input.BelongsToValidPreparation,
	}

	query, args := p.buildCreateRequiredPreparationInstrumentQuery(x)

	// create the required preparation instrument.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing required preparation instrument creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRequiredPreparationInstrumentQuery takes a required preparation instrument and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRequiredPreparationInstrumentQuery(input *models.RequiredPreparationInstrument) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(requiredPreparationInstrumentsTableName).
		Set("valid_instrument_id", input.ValidInstrumentID).
		Set("notes", input.Notes).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
			requiredPreparationInstrumentsTableOwnershipColumn: input.BelongsToValidPreparation,
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

// buildArchiveRequiredPreparationInstrumentQuery returns a SQL query which marks a given required preparation instrument belonging to a given valid preparation as archived.
func (p *Postgres) buildArchiveRequiredPreparationInstrumentQuery(validPreparationID, requiredPreparationInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(requiredPreparationInstrumentsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          requiredPreparationInstrumentID,
			"archived_on": nil,
			requiredPreparationInstrumentsTableOwnershipColumn: validPreparationID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRequiredPreparationInstrument marks a required preparation instrument as archived in the database.
func (p *Postgres) ArchiveRequiredPreparationInstrument(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) error {
	query, args := p.buildArchiveRequiredPreparationInstrumentQuery(validPreparationID, requiredPreparationInstrumentID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
