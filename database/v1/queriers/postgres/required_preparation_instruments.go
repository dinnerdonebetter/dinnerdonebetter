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
	requiredPreparationInstrumentsTableName                = "required_preparation_instruments"
	requiredPreparationInstrumentsTableInstrumentIDColumn  = "instrument_id"
	requiredPreparationInstrumentsTablePreparationIDColumn = "preparation_id"
	requiredPreparationInstrumentsTableNotesColumn         = "notes"
)

var (
	requiredPreparationInstrumentsTableColumns = []string{
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, idColumn),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, requiredPreparationInstrumentsTableInstrumentIDColumn),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, requiredPreparationInstrumentsTablePreparationIDColumn),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, requiredPreparationInstrumentsTableNotesColumn),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, archivedOnColumn),
	}
)

// scanRequiredPreparationInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a Required Preparation Instrument struct
func (p *Postgres) scanRequiredPreparationInstrument(scan database.Scanner, includeCount bool) (*models.RequiredPreparationInstrument, uint64, error) {
	x := &models.RequiredPreparationInstrument{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.InstrumentID,
		&x.PreparationID,
		&x.Notes,
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

// buildRequiredPreparationInstrumentExistsQuery constructs a SQL query for checking if a required preparation instrument with a given ID exists
func (p *Postgres) buildRequiredPreparationInstrumentExistsQuery(requiredPreparationInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, idColumn)).
		Prefix(existencePrefix).
		From(requiredPreparationInstrumentsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, idColumn): requiredPreparationInstrumentID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RequiredPreparationInstrumentExists queries the database to see if a given required preparation instrument belonging to a given user exists.
func (p *Postgres) RequiredPreparationInstrumentExists(ctx context.Context, requiredPreparationInstrumentID uint64) (exists bool, err error) {
	query, args := p.buildRequiredPreparationInstrumentExistsQuery(requiredPreparationInstrumentID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRequiredPreparationInstrumentQuery constructs a SQL query for fetching a required preparation instrument with a given ID.
func (p *Postgres) buildGetRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(requiredPreparationInstrumentsTableColumns...).
		From(requiredPreparationInstrumentsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, idColumn): requiredPreparationInstrumentID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRequiredPreparationInstrument fetches a required preparation instrument from the database.
func (p *Postgres) GetRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) (*models.RequiredPreparationInstrument, error) {
	query, args := p.buildGetRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID)
	row := p.db.QueryRowContext(ctx, query, args...)
	rpi, _, err := p.scanRequiredPreparationInstrument(row, false)
	return rpi, err
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
				fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, archivedOnColumn): nil,
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

// buildGetBatchOfRequiredPreparationInstrumentsQuery returns a query that fetches every required preparation instrument in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfRequiredPreparationInstrumentsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(requiredPreparationInstrumentsTableColumns...).
		From(requiredPreparationInstrumentsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllRequiredPreparationInstruments fetches every required preparation instrument from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllRequiredPreparationInstruments(ctx context.Context, resultChannel chan []models.RequiredPreparationInstrument) error {
	count, err := p.GetAllRequiredPreparationInstrumentsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfRequiredPreparationInstrumentsQuery(begin, end)
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

			requiredPreparationInstruments, _, err := p.scanRequiredPreparationInstruments(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- requiredPreparationInstruments
		}(beginID, endID)
	}

	return nil
}

// buildGetRequiredPreparationInstrumentsQuery builds a SQL query selecting required preparation instruments that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRequiredPreparationInstrumentsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(requiredPreparationInstrumentsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllRequiredPreparationInstrumentsCountQuery()))...).
		From(requiredPreparationInstrumentsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, requiredPreparationInstrumentsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRequiredPreparationInstruments fetches a list of required preparation instruments from the database that meet a particular filter.
func (p *Postgres) GetRequiredPreparationInstruments(ctx context.Context, filter *models.QueryFilter) (*models.RequiredPreparationInstrumentList, error) {
	query, args := p.buildGetRequiredPreparationInstrumentsQuery(filter)

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

// buildGetRequiredPreparationInstrumentsWithIDsQuery builds a SQL query selecting requiredPreparationInstruments
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetRequiredPreparationInstrumentsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(requiredPreparationInstrumentsTableColumns...).
		From(requiredPreparationInstrumentsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(requiredPreparationInstrumentsTableColumns...).
		FromSelect(subqueryBuilder, requiredPreparationInstrumentsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", requiredPreparationInstrumentsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRequiredPreparationInstrumentsWithIDs fetches a list of required preparation instruments from the database that exist within a given set of IDs.
func (p *Postgres) GetRequiredPreparationInstrumentsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.RequiredPreparationInstrument, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetRequiredPreparationInstrumentsWithIDsQuery(limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for required preparation instruments")
	}

	requiredPreparationInstruments, _, err := p.scanRequiredPreparationInstruments(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return requiredPreparationInstruments, nil
}

// buildCreateRequiredPreparationInstrumentQuery takes a required preparation instrument and returns a creation query for that required preparation instrument and the relevant arguments.
func (p *Postgres) buildCreateRequiredPreparationInstrumentQuery(input *models.RequiredPreparationInstrument) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(requiredPreparationInstrumentsTableName).
		Columns(
			requiredPreparationInstrumentsTableInstrumentIDColumn,
			requiredPreparationInstrumentsTablePreparationIDColumn,
			requiredPreparationInstrumentsTableNotesColumn,
		).
		Values(
			input.InstrumentID,
			input.PreparationID,
			input.Notes,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRequiredPreparationInstrument creates a required preparation instrument in the database.
func (p *Postgres) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (*models.RequiredPreparationInstrument, error) {
	x := &models.RequiredPreparationInstrument{
		InstrumentID:  input.InstrumentID,
		PreparationID: input.PreparationID,
		Notes:         input.Notes,
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
		Set(requiredPreparationInstrumentsTableInstrumentIDColumn, input.InstrumentID).
		Set(requiredPreparationInstrumentsTablePreparationIDColumn, input.PreparationID).
		Set(requiredPreparationInstrumentsTableNotesColumn, input.Notes).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRequiredPreparationInstrument updates a particular required preparation instrument. Note that UpdateRequiredPreparationInstrument expects the provided input to have a valid ID.
func (p *Postgres) UpdateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrument) error {
	query, args := p.buildUpdateRequiredPreparationInstrumentQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveRequiredPreparationInstrumentQuery returns a SQL query which marks a given required preparation instrument as archived.
func (p *Postgres) buildArchiveRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(requiredPreparationInstrumentsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         requiredPreparationInstrumentID,
			archivedOnColumn: nil,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRequiredPreparationInstrument marks a required preparation instrument as archived in the database.
func (p *Postgres) ArchiveRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrumentID uint64) error {
	query, args := p.buildArchiveRequiredPreparationInstrumentQuery(requiredPreparationInstrumentID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
