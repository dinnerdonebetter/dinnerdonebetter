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
	validIngredientPreparationsTableName                     = "valid_ingredient_preparations"
	validIngredientPreparationsTableNotesColumn              = "notes"
	validIngredientPreparationsTableValidPreparationIDColumn = "valid_preparation_id"
	validIngredientPreparationsTableValidIngredientIDColumn  = "valid_ingredient_id"
)

var (
	validIngredientPreparationsTableColumns = []string{
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, idColumn),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, validIngredientPreparationsTableNotesColumn),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, validIngredientPreparationsTableValidPreparationIDColumn),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, validIngredientPreparationsTableValidIngredientIDColumn),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, archivedOnColumn),
	}
)

// scanValidIngredientPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a Valid Ingredient Preparation struct
func (p *Postgres) scanValidIngredientPreparation(scan database.Scanner, includeCount bool) (*models.ValidIngredientPreparation, uint64, error) {
	x := &models.ValidIngredientPreparation{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Notes,
		&x.ValidPreparationID,
		&x.ValidIngredientID,
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

// scanValidIngredientPreparations takes a logger and some database rows and turns them into a slice of valid ingredient preparations.
func (p *Postgres) scanValidIngredientPreparations(rows database.ResultIterator) ([]models.ValidIngredientPreparation, uint64, error) {
	var (
		list  []models.ValidIngredientPreparation
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanValidIngredientPreparation(rows, true)
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

// buildValidIngredientPreparationExistsQuery constructs a SQL query for checking if a valid ingredient preparation with a given ID exists
func (p *Postgres) buildValidIngredientPreparationExistsQuery(validIngredientPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, idColumn)).
		Prefix(existencePrefix).
		From(validIngredientPreparationsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, idColumn): validIngredientPreparationID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ValidIngredientPreparationExists queries the database to see if a given valid ingredient preparation belonging to a given user exists.
func (p *Postgres) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID uint64) (exists bool, err error) {
	query, args := p.buildValidIngredientPreparationExistsQuery(validIngredientPreparationID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetValidIngredientPreparationQuery constructs a SQL query for fetching a valid ingredient preparation with a given ID.
func (p *Postgres) buildGetValidIngredientPreparationQuery(validIngredientPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(validIngredientPreparationsTableColumns...).
		From(validIngredientPreparationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, idColumn): validIngredientPreparationID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredientPreparation fetches a valid ingredient preparation from the database.
func (p *Postgres) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) (*models.ValidIngredientPreparation, error) {
	query, args := p.buildGetValidIngredientPreparationQuery(validIngredientPreparationID)
	row := p.db.QueryRowContext(ctx, query, args...)
	vip, _, err := p.scanValidIngredientPreparation(row, false)
	return vip, err
}

var (
	allValidIngredientPreparationsCountQueryBuilder sync.Once
	allValidIngredientPreparationsCountQuery        string
)

// buildGetAllValidIngredientPreparationsCountQuery returns a query that fetches the total number of valid ingredient preparations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllValidIngredientPreparationsCountQuery() string {
	allValidIngredientPreparationsCountQueryBuilder.Do(func() {
		var err error

		allValidIngredientPreparationsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, validIngredientPreparationsTableName)).
			From(validIngredientPreparationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allValidIngredientPreparationsCountQuery
}

// GetAllValidIngredientPreparationsCount will fetch the count of valid ingredient preparations from the database.
func (p *Postgres) GetAllValidIngredientPreparationsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllValidIngredientPreparationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfValidIngredientPreparationsQuery returns a query that fetches every valid ingredient preparation in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfValidIngredientPreparationsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(validIngredientPreparationsTableColumns...).
		From(validIngredientPreparationsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllValidIngredientPreparations fetches every valid ingredient preparation from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllValidIngredientPreparations(ctx context.Context, resultChannel chan []models.ValidIngredientPreparation) error {
	count, err := p.GetAllValidIngredientPreparationsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfValidIngredientPreparationsQuery(begin, end)
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

			validIngredientPreparations, _, err := p.scanValidIngredientPreparations(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- validIngredientPreparations
		}(beginID, endID)
	}

	return nil
}

// buildGetValidIngredientPreparationsQuery builds a SQL query selecting valid ingredient preparations that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetValidIngredientPreparationsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(validIngredientPreparationsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllValidIngredientPreparationsCountQuery()))...).
		From(validIngredientPreparationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, validIngredientPreparationsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredientPreparations fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (p *Postgres) GetValidIngredientPreparations(ctx context.Context, filter *models.QueryFilter) (*models.ValidIngredientPreparationList, error) {
	query, args := p.buildGetValidIngredientPreparationsQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for valid ingredient preparations")
	}

	validIngredientPreparations, count, err := p.scanValidIngredientPreparations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.ValidIngredientPreparationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		ValidIngredientPreparations: validIngredientPreparations,
	}

	return list, nil
}

// buildGetValidIngredientPreparationsWithIDsQuery builds a SQL query selecting validIngredientPreparations
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetValidIngredientPreparationsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(validIngredientPreparationsTableColumns...).
		From(validIngredientPreparationsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(validIngredientPreparationsTableColumns...).
		FromSelect(subqueryBuilder, validIngredientPreparationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredientPreparationsWithIDs fetches a list of valid ingredient preparations from the database that exist within a given set of IDs.
func (p *Postgres) GetValidIngredientPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidIngredientPreparation, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetValidIngredientPreparationsWithIDsQuery(limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for valid ingredient preparations")
	}

	validIngredientPreparations, _, err := p.scanValidIngredientPreparations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return validIngredientPreparations, nil
}

// buildCreateValidIngredientPreparationQuery takes a valid ingredient preparation and returns a creation query for that valid ingredient preparation and the relevant arguments.
func (p *Postgres) buildCreateValidIngredientPreparationQuery(input *models.ValidIngredientPreparation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(validIngredientPreparationsTableName).
		Columns(
			validIngredientPreparationsTableNotesColumn,
			validIngredientPreparationsTableValidPreparationIDColumn,
			validIngredientPreparationsTableValidIngredientIDColumn,
		).
		Values(
			input.Notes,
			input.ValidPreparationID,
			input.ValidIngredientID,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateValidIngredientPreparation creates a valid ingredient preparation in the database.
func (p *Postgres) CreateValidIngredientPreparation(ctx context.Context, input *models.ValidIngredientPreparationCreationInput) (*models.ValidIngredientPreparation, error) {
	x := &models.ValidIngredientPreparation{
		Notes:              input.Notes,
		ValidPreparationID: input.ValidPreparationID,
		ValidIngredientID:  input.ValidIngredientID,
	}

	query, args := p.buildCreateValidIngredientPreparationQuery(x)

	// create the valid ingredient preparation.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing valid ingredient preparation creation query: %w", err)
	}

	return x, nil
}

// buildUpdateValidIngredientPreparationQuery takes a valid ingredient preparation and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateValidIngredientPreparationQuery(input *models.ValidIngredientPreparation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validIngredientPreparationsTableName).
		Set(validIngredientPreparationsTableNotesColumn, input.Notes).
		Set(validIngredientPreparationsTableValidPreparationIDColumn, input.ValidPreparationID).
		Set(validIngredientPreparationsTableValidIngredientIDColumn, input.ValidIngredientID).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateValidIngredientPreparation updates a particular valid ingredient preparation. Note that UpdateValidIngredientPreparation expects the provided input to have a valid ID.
func (p *Postgres) UpdateValidIngredientPreparation(ctx context.Context, input *models.ValidIngredientPreparation) error {
	query, args := p.buildUpdateValidIngredientPreparationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveValidIngredientPreparationQuery returns a SQL query which marks a given valid ingredient preparation as archived.
func (p *Postgres) buildArchiveValidIngredientPreparationQuery(validIngredientPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validIngredientPreparationsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         validIngredientPreparationID,
			archivedOnColumn: nil,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveValidIngredientPreparation marks a valid ingredient preparation as archived in the database.
func (p *Postgres) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) error {
	query, args := p.buildArchiveValidIngredientPreparationQuery(validIngredientPreparationID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
