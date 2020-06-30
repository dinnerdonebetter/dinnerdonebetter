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
	validPreparationsTableName = "valid_preparations"
)

var (
	validPreparationsTableColumns = []string{
		fmt.Sprintf("%s.%s", validPreparationsTableName, "id"),
		fmt.Sprintf("%s.%s", validPreparationsTableName, "name"),
		fmt.Sprintf("%s.%s", validPreparationsTableName, "description"),
		fmt.Sprintf("%s.%s", validPreparationsTableName, "icon"),
		fmt.Sprintf("%s.%s", validPreparationsTableName, "applicable_to_all_ingredients"),
		fmt.Sprintf("%s.%s", validPreparationsTableName, "created_on"),
		fmt.Sprintf("%s.%s", validPreparationsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", validPreparationsTableName, "archived_on"),
	}

	validPreparationsOnRequiredPreparationInstrumentsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.id", validPreparationsTableName, requiredPreparationInstrumentsTableName, requiredPreparationInstrumentsTableOwnershipColumn, validPreparationsTableName)
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
		&x.UpdatedOn,
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
		Select(fmt.Sprintf("%s.id", validPreparationsTableName)).
		Prefix(existencePrefix).
		From(validPreparationsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", validPreparationsTableName): validPreparationID,
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
			fmt.Sprintf("%s.id", validPreparationsTableName): validPreparationID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetValidPreparation fetches a valid preparation from the database.
func (p *Postgres) GetValidPreparation(ctx context.Context, validPreparationID uint64) (*models.ValidPreparation, error) {
	query, args := p.buildGetValidPreparationQuery(validPreparationID)
	row := p.db.QueryRowContext(ctx, query, args...)

	validPreparation, _, err := p.scanValidPreparation(row, false)
	return validPreparation, err
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
				fmt.Sprintf("%s.archived_on", validPreparationsTableName): nil,
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

// buildGetValidPreparationsQuery builds a SQL query selecting valid preparations that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetValidPreparationsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(validPreparationsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllValidPreparationsCountQuery()))...).
		From(validPreparationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", validPreparationsTableName): nil,
		}).
		OrderBy(fmt.Sprintf("%s.id", validPreparationsTableName))

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

// buildCreateValidPreparationQuery takes a valid preparation and returns a creation query for that valid preparation and the relevant arguments.
func (p *Postgres) buildCreateValidPreparationQuery(input *models.ValidPreparation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(validPreparationsTableName).
		Columns(
			"name",
			"description",
			"icon",
		).
		Values(
			input.Name,
			input.Description,
			input.Icon,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateValidPreparation creates a valid preparation in the database.
func (p *Postgres) CreateValidPreparation(ctx context.Context, input *models.ValidPreparationCreationInput) (*models.ValidPreparation, error) {
	x := &models.ValidPreparation{
		Name:        input.Name,
		Description: input.Description,
		Icon:        input.Icon,
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
		Set("name", input.Name).
		Set("description", input.Description).
		Set("icon", input.Icon).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateValidPreparation updates a particular valid preparation. Note that UpdateValidPreparation expects the provided input to have a valid ID.
func (p *Postgres) UpdateValidPreparation(ctx context.Context, input *models.ValidPreparation) error {
	query, args := p.buildUpdateValidPreparationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveValidPreparationQuery returns a SQL query which marks a given valid preparation as archived.
func (p *Postgres) buildArchiveValidPreparationQuery(validPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validPreparationsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          validPreparationID,
			"archived_on": nil,
		}).
		Suffix("RETURNING archived_on").
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
