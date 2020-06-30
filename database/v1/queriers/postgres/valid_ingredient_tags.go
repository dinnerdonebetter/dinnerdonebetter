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
	validIngredientTagsTableName = "valid_ingredient_tags"
)

var (
	validIngredientTagsTableColumns = []string{
		fmt.Sprintf("%s.%s", validIngredientTagsTableName, "id"),
		fmt.Sprintf("%s.%s", validIngredientTagsTableName, "name"),
		fmt.Sprintf("%s.%s", validIngredientTagsTableName, "created_on"),
		fmt.Sprintf("%s.%s", validIngredientTagsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", validIngredientTagsTableName, "archived_on"),
	}
)

// scanValidIngredientTag takes a database Scanner (i.e. *sql.Row) and scans the result into a Valid Ingredient Tag struct
func (p *Postgres) scanValidIngredientTag(scan database.Scanner, includeCount bool) (*models.ValidIngredientTag, uint64, error) {
	x := &models.ValidIngredientTag{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
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

// scanValidIngredientTags takes a logger and some database rows and turns them into a slice of valid ingredient tags.
func (p *Postgres) scanValidIngredientTags(rows database.ResultIterator) ([]models.ValidIngredientTag, uint64, error) {
	var (
		list  []models.ValidIngredientTag
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanValidIngredientTag(rows, true)
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

// buildValidIngredientTagExistsQuery constructs a SQL query for checking if a valid ingredient tag with a given ID exists
func (p *Postgres) buildValidIngredientTagExistsQuery(validIngredientTagID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", validIngredientTagsTableName)).
		Prefix(existencePrefix).
		From(validIngredientTagsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", validIngredientTagsTableName): validIngredientTagID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ValidIngredientTagExists queries the database to see if a given valid ingredient tag belonging to a given user exists.
func (p *Postgres) ValidIngredientTagExists(ctx context.Context, validIngredientTagID uint64) (exists bool, err error) {
	query, args := p.buildValidIngredientTagExistsQuery(validIngredientTagID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetValidIngredientTagQuery constructs a SQL query for fetching a valid ingredient tag with a given ID.
func (p *Postgres) buildGetValidIngredientTagQuery(validIngredientTagID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(validIngredientTagsTableColumns...).
		From(validIngredientTagsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", validIngredientTagsTableName): validIngredientTagID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredientTag fetches a valid ingredient tag from the database.
func (p *Postgres) GetValidIngredientTag(ctx context.Context, validIngredientTagID uint64) (*models.ValidIngredientTag, error) {
	query, args := p.buildGetValidIngredientTagQuery(validIngredientTagID)
	row := p.db.QueryRowContext(ctx, query, args...)

	validIngredientTag, _, err := p.scanValidIngredientTag(row, false)
	return validIngredientTag, err
}

var (
	allValidIngredientTagsCountQueryBuilder sync.Once
	allValidIngredientTagsCountQuery        string
)

// buildGetAllValidIngredientTagsCountQuery returns a query that fetches the total number of valid ingredient tags in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllValidIngredientTagsCountQuery() string {
	allValidIngredientTagsCountQueryBuilder.Do(func() {
		var err error

		allValidIngredientTagsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, validIngredientTagsTableName)).
			From(validIngredientTagsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", validIngredientTagsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allValidIngredientTagsCountQuery
}

// GetAllValidIngredientTagsCount will fetch the count of valid ingredient tags from the database.
func (p *Postgres) GetAllValidIngredientTagsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllValidIngredientTagsCountQuery()).Scan(&count)
	return count, err
}

// buildGetValidIngredientTagsQuery builds a SQL query selecting valid ingredient tags that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetValidIngredientTagsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(validIngredientTagsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllValidIngredientTagsCountQuery()))...).
		From(validIngredientTagsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", validIngredientTagsTableName): nil,
		}).
		OrderBy(fmt.Sprintf("%s.id", validIngredientTagsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, validIngredientTagsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredientTags fetches a list of valid ingredient tags from the database that meet a particular filter.
func (p *Postgres) GetValidIngredientTags(ctx context.Context, filter *models.QueryFilter) (*models.ValidIngredientTagList, error) {
	query, args := p.buildGetValidIngredientTagsQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for valid ingredient tags")
	}

	validIngredientTags, count, err := p.scanValidIngredientTags(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.ValidIngredientTagList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		ValidIngredientTags: validIngredientTags,
	}

	return list, nil
}

// buildCreateValidIngredientTagQuery takes a valid ingredient tag and returns a creation query for that valid ingredient tag and the relevant arguments.
func (p *Postgres) buildCreateValidIngredientTagQuery(input *models.ValidIngredientTag) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(validIngredientTagsTableName).
		Columns(
			"name",
		).
		Values(
			input.Name,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateValidIngredientTag creates a valid ingredient tag in the database.
func (p *Postgres) CreateValidIngredientTag(ctx context.Context, input *models.ValidIngredientTagCreationInput) (*models.ValidIngredientTag, error) {
	x := &models.ValidIngredientTag{
		Name: input.Name,
	}

	query, args := p.buildCreateValidIngredientTagQuery(x)

	// create the valid ingredient tag.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing valid ingredient tag creation query: %w", err)
	}

	return x, nil
}

// buildUpdateValidIngredientTagQuery takes a valid ingredient tag and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateValidIngredientTagQuery(input *models.ValidIngredientTag) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validIngredientTagsTableName).
		Set("name", input.Name).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateValidIngredientTag updates a particular valid ingredient tag. Note that UpdateValidIngredientTag expects the provided input to have a valid ID.
func (p *Postgres) UpdateValidIngredientTag(ctx context.Context, input *models.ValidIngredientTag) error {
	query, args := p.buildUpdateValidIngredientTagQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveValidIngredientTagQuery returns a SQL query which marks a given valid ingredient tag as archived.
func (p *Postgres) buildArchiveValidIngredientTagQuery(validIngredientTagID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validIngredientTagsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          validIngredientTagID,
			"archived_on": nil,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveValidIngredientTag marks a valid ingredient tag as archived in the database.
func (p *Postgres) ArchiveValidIngredientTag(ctx context.Context, validIngredientTagID uint64) error {
	query, args := p.buildArchiveValidIngredientTagQuery(validIngredientTagID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
