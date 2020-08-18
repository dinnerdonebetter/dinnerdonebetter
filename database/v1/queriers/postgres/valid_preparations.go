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
	validPreparationsTableName                             = "valid_preparations"
	validPreparationsTableNameColumn                       = "name"
	validPreparationsTableDescriptionColumn                = "description"
	validPreparationsTableIconColumn                       = "icon"
	validPreparationsTableApplicableToAllIngredientsColumn = "applicable_to_all_ingredients"
)

var (
	validPreparationsTableColumns = []string{
		fmt.Sprintf("%s.%s", validPreparationsTableName, idColumn),
		fmt.Sprintf("%s.%s", validPreparationsTableName, validPreparationsTableNameColumn),
		fmt.Sprintf("%s.%s", validPreparationsTableName, validPreparationsTableDescriptionColumn),
		fmt.Sprintf("%s.%s", validPreparationsTableName, validPreparationsTableIconColumn),
		fmt.Sprintf("%s.%s", validPreparationsTableName, validPreparationsTableApplicableToAllIngredientsColumn),
		fmt.Sprintf("%s.%s", validPreparationsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", validPreparationsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", validPreparationsTableName, archivedOnColumn),
	}
)

// scanValidPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a Valid Preparation struct
func (p *Postgres) scanValidPreparation(scan database.Scanner, includeCount bool) (*models.ValidPreparation, uint64, error) {
	x := &models.ValidPreparation{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Description,
		&x.Icon,
		&x.ApplicableToAllIngredients,
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

// scanValidPreparations takes a logger and some database rows and turns them into a slice of valid preparations.
func (p *Postgres) scanValidPreparations(rows database.ResultIterator) ([]models.ValidPreparation, uint64, error) {
	var (
		list  []models.ValidPreparation
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanValidPreparation(rows, true)
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

// buildValidPreparationExistsQuery constructs a SQL query for checking if a valid preparation with a given ID exists
func (p *Postgres) buildValidPreparationExistsQuery(validPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", validPreparationsTableName, idColumn)).
		Prefix(existencePrefix).
		From(validPreparationsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validPreparationsTableName, idColumn): validPreparationID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ValidPreparationExists queries the database to see if a given valid preparation belonging to a given user exists.
func (p *Postgres) ValidPreparationExists(ctx context.Context, validPreparationID uint64) (exists bool, err error) {
	query, args := p.buildValidPreparationExistsQuery(validPreparationID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetValidPreparationQuery constructs a SQL query for fetching a valid preparation with a given ID.
func (p *Postgres) buildGetValidPreparationQuery(validPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(validPreparationsTableColumns...).
		From(validPreparationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validPreparationsTableName, idColumn): validPreparationID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetValidPreparation fetches a valid preparation from the database.
func (p *Postgres) GetValidPreparation(ctx context.Context, validPreparationID uint64) (*models.ValidPreparation, error) {
	query, args := p.buildGetValidPreparationQuery(validPreparationID)
	row := p.db.QueryRowContext(ctx, query, args...)
	vp, _, err := p.scanValidPreparation(row, false)
	return vp, err
}

var (
	allValidPreparationsCountQueryBuilder sync.Once
	allValidPreparationsCountQuery        string
)

// buildGetAllValidPreparationsCountQuery returns a query that fetches the total number of valid preparations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllValidPreparationsCountQuery() string {
	allValidPreparationsCountQueryBuilder.Do(func() {
		var err error

		allValidPreparationsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, validPreparationsTableName)).
			From(validPreparationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", validPreparationsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allValidPreparationsCountQuery
}

// GetAllValidPreparationsCount will fetch the count of valid preparations from the database.
func (p *Postgres) GetAllValidPreparationsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllValidPreparationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfValidPreparationsQuery returns a query that fetches every valid preparation in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfValidPreparationsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(validPreparationsTableColumns...).
		From(validPreparationsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", validPreparationsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", validPreparationsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllValidPreparations fetches every valid preparation from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllValidPreparations(ctx context.Context, resultChannel chan []models.ValidPreparation) error {
	count, err := p.GetAllValidPreparationsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfValidPreparationsQuery(begin, end)
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

			validPreparations, _, err := p.scanValidPreparations(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- validPreparations
		}(beginID, endID)
	}

	return nil
}

// buildGetValidPreparationsQuery builds a SQL query selecting valid preparations that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetValidPreparationsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(validPreparationsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllValidPreparationsCountQuery()))...).
		From(validPreparationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validPreparationsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", validPreparationsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, validPreparationsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidPreparations fetches a list of valid preparations from the database that meet a particular filter.
func (p *Postgres) GetValidPreparations(ctx context.Context, filter *models.QueryFilter) (*models.ValidPreparationList, error) {
	query, args := p.buildGetValidPreparationsQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for valid preparations")
	}

	validPreparations, count, err := p.scanValidPreparations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.ValidPreparationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		ValidPreparations: validPreparations,
	}

	return list, nil
}

// buildGetValidPreparationsWithIDsQuery builds a SQL query selecting validPreparations
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetValidPreparationsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(validPreparationsTableColumns...).
		From(validPreparationsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(validPreparationsTableColumns...).
		FromSelect(subqueryBuilder, validPreparationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validPreparationsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidPreparationsWithIDs fetches a list of valid preparations from the database that exist within a given set of IDs.
func (p *Postgres) GetValidPreparationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidPreparation, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetValidPreparationsWithIDsQuery(limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for valid preparations")
	}

	validPreparations, _, err := p.scanValidPreparations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return validPreparations, nil
}

// buildCreateValidPreparationQuery takes a valid preparation and returns a creation query for that valid preparation and the relevant arguments.
func (p *Postgres) buildCreateValidPreparationQuery(input *models.ValidPreparation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(validPreparationsTableName).
		Columns(
			validPreparationsTableNameColumn,
			validPreparationsTableDescriptionColumn,
			validPreparationsTableIconColumn,
			validPreparationsTableApplicableToAllIngredientsColumn,
		).
		Values(
			input.Name,
			input.Description,
			input.Icon,
			input.ApplicableToAllIngredients,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateValidPreparation creates a valid preparation in the database.
func (p *Postgres) CreateValidPreparation(ctx context.Context, input *models.ValidPreparationCreationInput) (*models.ValidPreparation, error) {
	x := &models.ValidPreparation{
		Name:                       input.Name,
		Description:                input.Description,
		Icon:                       input.Icon,
		ApplicableToAllIngredients: input.ApplicableToAllIngredients,
	}

	query, args := p.buildCreateValidPreparationQuery(x)

	// create the valid preparation.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing valid preparation creation query: %w", err)
	}

	return x, nil
}

// buildUpdateValidPreparationQuery takes a valid preparation and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateValidPreparationQuery(input *models.ValidPreparation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validPreparationsTableName).
		Set(validPreparationsTableNameColumn, input.Name).
		Set(validPreparationsTableDescriptionColumn, input.Description).
		Set(validPreparationsTableIconColumn, input.Icon).
		Set(validPreparationsTableApplicableToAllIngredientsColumn, input.ApplicableToAllIngredients).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateValidPreparation updates a particular valid preparation. Note that UpdateValidPreparation expects the provided input to have a valid ID.
func (p *Postgres) UpdateValidPreparation(ctx context.Context, input *models.ValidPreparation) error {
	query, args := p.buildUpdateValidPreparationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveValidPreparationQuery returns a SQL query which marks a given valid preparation as archived.
func (p *Postgres) buildArchiveValidPreparationQuery(validPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validPreparationsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         validPreparationID,
			archivedOnColumn: nil,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveValidPreparation marks a valid preparation as archived in the database.
func (p *Postgres) ArchiveValidPreparation(ctx context.Context, validPreparationID uint64) error {
	query, args := p.buildArchiveValidPreparationQuery(validPreparationID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
