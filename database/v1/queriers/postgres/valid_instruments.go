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
	validInstrumentsTableName              = "valid_instruments"
	validInstrumentsTableNameColumn        = "name"
	validInstrumentsTableVariantColumn     = "variant"
	validInstrumentsTableDescriptionColumn = "description"
	validInstrumentsTableIconColumn        = "icon"
)

var (
	validInstrumentsTableColumns = []string{
		fmt.Sprintf("%s.%s", validInstrumentsTableName, idColumn),
		fmt.Sprintf("%s.%s", validInstrumentsTableName, validInstrumentsTableNameColumn),
		fmt.Sprintf("%s.%s", validInstrumentsTableName, validInstrumentsTableVariantColumn),
		fmt.Sprintf("%s.%s", validInstrumentsTableName, validInstrumentsTableDescriptionColumn),
		fmt.Sprintf("%s.%s", validInstrumentsTableName, validInstrumentsTableIconColumn),
		fmt.Sprintf("%s.%s", validInstrumentsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", validInstrumentsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", validInstrumentsTableName, archivedOnColumn),
	}
)

// scanValidInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a Valid Instrument struct
func (p *Postgres) scanValidInstrument(scan database.Scanner, includeCount bool) (*models.ValidInstrument, uint64, error) {
	x := &models.ValidInstrument{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Variant,
		&x.Description,
		&x.Icon,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}

// scanValidInstruments takes a logger and some database rows and turns them into a slice of valid instruments.
func (p *Postgres) scanValidInstruments(rows database.ResultIterator, includeCount bool) ([]models.ValidInstrument, uint64, error) {
	var (
		list  []models.ValidInstrument
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanValidInstrument(rows, includeCount)
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

// buildValidInstrumentExistsQuery constructs a SQL query for checking if a valid instrument with a given ID exists
func (p *Postgres) buildValidInstrumentExistsQuery(validInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", validInstrumentsTableName, idColumn)).
		Prefix(existencePrefix).
		From(validInstrumentsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validInstrumentsTableName, idColumn): validInstrumentID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ValidInstrumentExists queries the database to see if a given valid instrument belonging to a given user exists.
func (p *Postgres) ValidInstrumentExists(ctx context.Context, validInstrumentID uint64) (exists bool, err error) {
	query, args := p.buildValidInstrumentExistsQuery(validInstrumentID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetValidInstrumentQuery constructs a SQL query for fetching a valid instrument with a given ID.
func (p *Postgres) buildGetValidInstrumentQuery(validInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(validInstrumentsTableColumns...).
		From(validInstrumentsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validInstrumentsTableName, idColumn): validInstrumentID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetValidInstrument fetches a valid instrument from the database.
func (p *Postgres) GetValidInstrument(ctx context.Context, validInstrumentID uint64) (*models.ValidInstrument, error) {
	query, args := p.buildGetValidInstrumentQuery(validInstrumentID)
	row := p.db.QueryRowContext(ctx, query, args...)
	vi, _, err := p.scanValidInstrument(row, false)
	return vi, err
}

var (
	allValidInstrumentsCountQueryBuilder sync.Once
	allValidInstrumentsCountQuery        string
)

// buildGetAllValidInstrumentsCountQuery returns a query that fetches the total number of valid instruments in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllValidInstrumentsCountQuery() string {
	allValidInstrumentsCountQueryBuilder.Do(func() {
		var err error

		allValidInstrumentsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, validInstrumentsTableName)).
			From(validInstrumentsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", validInstrumentsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allValidInstrumentsCountQuery
}

// GetAllValidInstrumentsCount will fetch the count of valid instruments from the database.
func (p *Postgres) GetAllValidInstrumentsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllValidInstrumentsCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfValidInstrumentsQuery returns a query that fetches every valid instrument in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfValidInstrumentsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(validInstrumentsTableColumns...).
		From(validInstrumentsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", validInstrumentsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", validInstrumentsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllValidInstruments fetches every valid instrument from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllValidInstruments(ctx context.Context, resultChannel chan []models.ValidInstrument) error {
	count, err := p.GetAllValidInstrumentsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfValidInstrumentsQuery(begin, end)
			logger := p.logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, err := p.db.Query(query, args...)
			if err == sql.ErrNoRows {
				return
			} else if err != nil {
				logger.Error(err, "querying for database rows")
				return
			}

			validInstruments, _, err := p.scanValidInstruments(rows, false)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- validInstruments
		}(beginID, endID)
	}

	return nil
}

// buildGetValidInstrumentsQuery builds a SQL query selecting valid instruments that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetValidInstrumentsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(validInstrumentsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllValidInstrumentsCountQuery()))...).
		From(validInstrumentsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validInstrumentsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", validInstrumentsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, validInstrumentsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidInstruments fetches a list of valid instruments from the database that meet a particular filter.
func (p *Postgres) GetValidInstruments(ctx context.Context, filter *models.QueryFilter) (*models.ValidInstrumentList, error) {
	query, args := p.buildGetValidInstrumentsQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for valid instruments")
	}

	validInstruments, count, err := p.scanValidInstruments(rows, true)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.ValidInstrumentList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		ValidInstruments: validInstruments,
	}

	return list, nil
}

// buildGetValidInstrumentsWithIDsQuery builds a SQL query selecting validInstruments
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetValidInstrumentsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(validInstrumentsTableColumns...).
		From(validInstrumentsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(validInstrumentsTableColumns...).
		FromSelect(subqueryBuilder, validInstrumentsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validInstrumentsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidInstrumentsWithIDs fetches a list of valid instruments from the database that exist within a given set of IDs.
func (p *Postgres) GetValidInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidInstrument, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetValidInstrumentsWithIDsQuery(limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for valid instruments")
	}

	validInstruments, _, err := p.scanValidInstruments(rows, false)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return validInstruments, nil
}

// buildCreateValidInstrumentQuery takes a valid instrument and returns a creation query for that valid instrument and the relevant arguments.
func (p *Postgres) buildCreateValidInstrumentQuery(input *models.ValidInstrument) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(validInstrumentsTableName).
		Columns(
			validInstrumentsTableNameColumn,
			validInstrumentsTableVariantColumn,
			validInstrumentsTableDescriptionColumn,
			validInstrumentsTableIconColumn,
		).
		Values(
			input.Name,
			input.Variant,
			input.Description,
			input.Icon,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateValidInstrument creates a valid instrument in the database.
func (p *Postgres) CreateValidInstrument(ctx context.Context, input *models.ValidInstrumentCreationInput) (*models.ValidInstrument, error) {
	x := &models.ValidInstrument{
		Name:        input.Name,
		Variant:     input.Variant,
		Description: input.Description,
		Icon:        input.Icon,
	}

	query, args := p.buildCreateValidInstrumentQuery(x)

	// create the valid instrument.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing valid instrument creation query: %w", err)
	}

	return x, nil
}

// buildUpdateValidInstrumentQuery takes a valid instrument and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateValidInstrumentQuery(input *models.ValidInstrument) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validInstrumentsTableName).
		Set(validInstrumentsTableNameColumn, input.Name).
		Set(validInstrumentsTableVariantColumn, input.Variant).
		Set(validInstrumentsTableDescriptionColumn, input.Description).
		Set(validInstrumentsTableIconColumn, input.Icon).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateValidInstrument updates a particular valid instrument. Note that UpdateValidInstrument expects the provided input to have a valid ID.
func (p *Postgres) UpdateValidInstrument(ctx context.Context, input *models.ValidInstrument) error {
	query, args := p.buildUpdateValidInstrumentQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveValidInstrumentQuery returns a SQL query which marks a given valid instrument as archived.
func (p *Postgres) buildArchiveValidInstrumentQuery(validInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validInstrumentsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         validInstrumentID,
			archivedOnColumn: nil,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveValidInstrument marks a valid instrument as archived in the database.
func (p *Postgres) ArchiveValidInstrument(ctx context.Context, validInstrumentID uint64) error {
	query, args := p.buildArchiveValidInstrumentQuery(validInstrumentID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
