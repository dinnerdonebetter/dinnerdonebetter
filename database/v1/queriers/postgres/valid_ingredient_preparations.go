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
	validIngredientPreparationsTableName            = "valid_ingredient_preparations"
	validIngredientPreparationsTableOwnershipColumn = "belongs_to_valid_ingredient"
)

var (
	validIngredientPreparationsTableColumns = []string{
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, "id"),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, "notes"),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, "created_on"),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, validIngredientPreparationsTableOwnershipColumn),
	}
)

// scanValidIngredientPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a Valid Ingredient Preparation struct
func (p *Postgres) scanValidIngredientPreparation(scan database.Scanner, includeCount bool) (*models.ValidIngredientPreparation, uint64, error) {
	x := &models.ValidIngredientPreparation{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Notes,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToValidIngredient,
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

// buildValidIngredientPreparationExistsQuery constructs a SQL query for checking if a valid ingredient preparation with a given ID belong to a a valid ingredient with a given ID exists
func (p *Postgres) buildValidIngredientPreparationExistsQuery(validIngredientID, validIngredientPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", validIngredientPreparationsTableName)).
		Prefix(existencePrefix).
		From(validIngredientPreparationsTableName).
		Join(validIngredientsOnValidIngredientPreparationsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", validIngredientPreparationsTableName):                                                  validIngredientPreparationID,
			fmt.Sprintf("%s.id", validIngredientsTableName):                                                             validIngredientID,
			fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, validIngredientPreparationsTableOwnershipColumn): validIngredientID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ValidIngredientPreparationExists queries the database to see if a given valid ingredient preparation belonging to a given user exists.
func (p *Postgres) ValidIngredientPreparationExists(ctx context.Context, validIngredientID, validIngredientPreparationID uint64) (exists bool, err error) {
	query, args := p.buildValidIngredientPreparationExistsQuery(validIngredientID, validIngredientPreparationID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetValidIngredientPreparationQuery constructs a SQL query for fetching a valid ingredient preparation with a given ID belong to a valid ingredient with a given ID.
func (p *Postgres) buildGetValidIngredientPreparationQuery(validIngredientID, validIngredientPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(validIngredientPreparationsTableColumns...).
		From(validIngredientPreparationsTableName).
		Join(validIngredientsOnValidIngredientPreparationsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", validIngredientPreparationsTableName):                                                  validIngredientPreparationID,
			fmt.Sprintf("%s.id", validIngredientsTableName):                                                             validIngredientID,
			fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, validIngredientPreparationsTableOwnershipColumn): validIngredientID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredientPreparation fetches a valid ingredient preparation from the database.
func (p *Postgres) GetValidIngredientPreparation(ctx context.Context, validIngredientID, validIngredientPreparationID uint64) (*models.ValidIngredientPreparation, error) {
	query, args := p.buildGetValidIngredientPreparationQuery(validIngredientID, validIngredientPreparationID)
	row := p.db.QueryRowContext(ctx, query, args...)

	validIngredientPreparation, _, err := p.scanValidIngredientPreparation(row, false)
	return validIngredientPreparation, err
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
				fmt.Sprintf("%s.archived_on", validIngredientPreparationsTableName): nil,
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

// buildGetValidIngredientPreparationsQuery builds a SQL query selecting valid ingredient preparations that adhere to a given QueryFilter and belong to a given valid ingredient,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetValidIngredientPreparationsQuery(validIngredientID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(validIngredientPreparationsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllValidIngredientPreparationsCountQuery()))...).
		From(validIngredientPreparationsTableName).
		Join(validIngredientsOnValidIngredientPreparationsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", validIngredientPreparationsTableName):                                         nil,
			fmt.Sprintf("%s.id", validIngredientsTableName):                                                             validIngredientID,
			fmt.Sprintf("%s.%s", validIngredientPreparationsTableName, validIngredientPreparationsTableOwnershipColumn): validIngredientID,
		}).
		OrderBy(fmt.Sprintf("%s.id", validIngredientPreparationsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, validIngredientPreparationsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredientPreparations fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (p *Postgres) GetValidIngredientPreparations(ctx context.Context, validIngredientID uint64, filter *models.QueryFilter) (*models.ValidIngredientPreparationList, error) {
	query, args := p.buildGetValidIngredientPreparationsQuery(validIngredientID, filter)

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

// buildCreateValidIngredientPreparationQuery takes a valid ingredient preparation and returns a creation query for that valid ingredient preparation and the relevant arguments.
func (p *Postgres) buildCreateValidIngredientPreparationQuery(input *models.ValidIngredientPreparation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(validIngredientPreparationsTableName).
		Columns(
			"notes",
			validIngredientPreparationsTableOwnershipColumn,
		).
		Values(
			input.Notes,
			input.BelongsToValidIngredient,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateValidIngredientPreparation creates a valid ingredient preparation in the database.
func (p *Postgres) CreateValidIngredientPreparation(ctx context.Context, input *models.ValidIngredientPreparationCreationInput) (*models.ValidIngredientPreparation, error) {
	x := &models.ValidIngredientPreparation{
		Notes:                    input.Notes,
		BelongsToValidIngredient: input.BelongsToValidIngredient,
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
		Set("notes", input.Notes).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
			validIngredientPreparationsTableOwnershipColumn: input.BelongsToValidIngredient,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateValidIngredientPreparation updates a particular valid ingredient preparation. Note that UpdateValidIngredientPreparation expects the provided input to have a valid ID.
func (p *Postgres) UpdateValidIngredientPreparation(ctx context.Context, input *models.ValidIngredientPreparation) error {
	query, args := p.buildUpdateValidIngredientPreparationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveValidIngredientPreparationQuery returns a SQL query which marks a given valid ingredient preparation belonging to a given valid ingredient as archived.
func (p *Postgres) buildArchiveValidIngredientPreparationQuery(validIngredientID, validIngredientPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validIngredientPreparationsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          validIngredientPreparationID,
			"archived_on": nil,
			validIngredientPreparationsTableOwnershipColumn: validIngredientID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveValidIngredientPreparation marks a valid ingredient preparation as archived in the database.
func (p *Postgres) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientID, validIngredientPreparationID uint64) error {
	query, args := p.buildArchiveValidIngredientPreparationQuery(validIngredientID, validIngredientPreparationID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
