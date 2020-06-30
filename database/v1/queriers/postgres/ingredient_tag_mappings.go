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
	ingredientTagMappingsTableName            = "ingredient_tag_mappings"
	ingredientTagMappingsTableOwnershipColumn = "belongs_to_valid_ingredient"
)

var (
	ingredientTagMappingsTableColumns = []string{
		fmt.Sprintf("%s.%s", ingredientTagMappingsTableName, "id"),
		fmt.Sprintf("%s.%s", ingredientTagMappingsTableName, "valid_ingredient_tag_id"),
		fmt.Sprintf("%s.%s", ingredientTagMappingsTableName, "created_on"),
		fmt.Sprintf("%s.%s", ingredientTagMappingsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", ingredientTagMappingsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", ingredientTagMappingsTableName, ingredientTagMappingsTableOwnershipColumn),
	}
)

// scanIngredientTagMapping takes a database Scanner (i.e. *sql.Row) and scans the result into an Ingredient Tag Mapping struct
func (p *Postgres) scanIngredientTagMapping(scan database.Scanner, includeCount bool) (*models.IngredientTagMapping, uint64, error) {
	x := &models.IngredientTagMapping{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.ValidIngredientTagID,
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

// scanIngredientTagMappings takes a logger and some database rows and turns them into a slice of ingredient tag mappings.
func (p *Postgres) scanIngredientTagMappings(rows database.ResultIterator) ([]models.IngredientTagMapping, uint64, error) {
	var (
		list  []models.IngredientTagMapping
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanIngredientTagMapping(rows, true)
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

// buildIngredientTagMappingExistsQuery constructs a SQL query for checking if an ingredient tag mapping with a given ID belong to a a valid ingredient with a given ID exists
func (p *Postgres) buildIngredientTagMappingExistsQuery(validIngredientID, ingredientTagMappingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", ingredientTagMappingsTableName)).
		Prefix(existencePrefix).
		From(ingredientTagMappingsTableName).
		Join(validIngredientsOnIngredientTagMappingsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", ingredientTagMappingsTableName):                                            ingredientTagMappingID,
			fmt.Sprintf("%s.id", validIngredientsTableName):                                                 validIngredientID,
			fmt.Sprintf("%s.%s", ingredientTagMappingsTableName, ingredientTagMappingsTableOwnershipColumn): validIngredientID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// IngredientTagMappingExists queries the database to see if a given ingredient tag mapping belonging to a given user exists.
func (p *Postgres) IngredientTagMappingExists(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (exists bool, err error) {
	query, args := p.buildIngredientTagMappingExistsQuery(validIngredientID, ingredientTagMappingID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetIngredientTagMappingQuery constructs a SQL query for fetching an ingredient tag mapping with a given ID belong to a valid ingredient with a given ID.
func (p *Postgres) buildGetIngredientTagMappingQuery(validIngredientID, ingredientTagMappingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(ingredientTagMappingsTableColumns...).
		From(ingredientTagMappingsTableName).
		Join(validIngredientsOnIngredientTagMappingsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", ingredientTagMappingsTableName):                                            ingredientTagMappingID,
			fmt.Sprintf("%s.id", validIngredientsTableName):                                                 validIngredientID,
			fmt.Sprintf("%s.%s", ingredientTagMappingsTableName, ingredientTagMappingsTableOwnershipColumn): validIngredientID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetIngredientTagMapping fetches an ingredient tag mapping from the database.
func (p *Postgres) GetIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (*models.IngredientTagMapping, error) {
	query, args := p.buildGetIngredientTagMappingQuery(validIngredientID, ingredientTagMappingID)
	row := p.db.QueryRowContext(ctx, query, args...)

	ingredientTagMapping, _, err := p.scanIngredientTagMapping(row, false)
	return ingredientTagMapping, err
}

var (
	allIngredientTagMappingsCountQueryBuilder sync.Once
	allIngredientTagMappingsCountQuery        string
)

// buildGetAllIngredientTagMappingsCountQuery returns a query that fetches the total number of ingredient tag mappings in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllIngredientTagMappingsCountQuery() string {
	allIngredientTagMappingsCountQueryBuilder.Do(func() {
		var err error

		allIngredientTagMappingsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, ingredientTagMappingsTableName)).
			From(ingredientTagMappingsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", ingredientTagMappingsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allIngredientTagMappingsCountQuery
}

// GetAllIngredientTagMappingsCount will fetch the count of ingredient tag mappings from the database.
func (p *Postgres) GetAllIngredientTagMappingsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllIngredientTagMappingsCountQuery()).Scan(&count)
	return count, err
}

// buildGetIngredientTagMappingsQuery builds a SQL query selecting ingredient tag mappings that adhere to a given QueryFilter and belong to a given valid ingredient,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetIngredientTagMappingsQuery(validIngredientID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(ingredientTagMappingsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllIngredientTagMappingsCountQuery()))...).
		From(ingredientTagMappingsTableName).
		Join(validIngredientsOnIngredientTagMappingsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", ingredientTagMappingsTableName):                                   nil,
			fmt.Sprintf("%s.id", validIngredientsTableName):                                                 validIngredientID,
			fmt.Sprintf("%s.%s", ingredientTagMappingsTableName, ingredientTagMappingsTableOwnershipColumn): validIngredientID,
		}).
		OrderBy(fmt.Sprintf("%s.id", ingredientTagMappingsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, ingredientTagMappingsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetIngredientTagMappings fetches a list of ingredient tag mappings from the database that meet a particular filter.
func (p *Postgres) GetIngredientTagMappings(ctx context.Context, validIngredientID uint64, filter *models.QueryFilter) (*models.IngredientTagMappingList, error) {
	query, args := p.buildGetIngredientTagMappingsQuery(validIngredientID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for ingredient tag mappings")
	}

	ingredientTagMappings, count, err := p.scanIngredientTagMappings(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.IngredientTagMappingList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		IngredientTagMappings: ingredientTagMappings,
	}

	return list, nil
}

// buildCreateIngredientTagMappingQuery takes an ingredient tag mapping and returns a creation query for that ingredient tag mapping and the relevant arguments.
func (p *Postgres) buildCreateIngredientTagMappingQuery(input *models.IngredientTagMapping) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(ingredientTagMappingsTableName).
		Columns(
			"valid_ingredient_tag_id",
			ingredientTagMappingsTableOwnershipColumn,
		).
		Values(
			input.ValidIngredientTagID,
			input.BelongsToValidIngredient,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateIngredientTagMapping creates an ingredient tag mapping in the database.
func (p *Postgres) CreateIngredientTagMapping(ctx context.Context, input *models.IngredientTagMappingCreationInput) (*models.IngredientTagMapping, error) {
	x := &models.IngredientTagMapping{
		ValidIngredientTagID:     input.ValidIngredientTagID,
		BelongsToValidIngredient: input.BelongsToValidIngredient,
	}

	query, args := p.buildCreateIngredientTagMappingQuery(x)

	// create the ingredient tag mapping.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing ingredient tag mapping creation query: %w", err)
	}

	return x, nil
}

// buildUpdateIngredientTagMappingQuery takes an ingredient tag mapping and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateIngredientTagMappingQuery(input *models.IngredientTagMapping) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(ingredientTagMappingsTableName).
		Set("valid_ingredient_tag_id", input.ValidIngredientTagID).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
			ingredientTagMappingsTableOwnershipColumn: input.BelongsToValidIngredient,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateIngredientTagMapping updates a particular ingredient tag mapping. Note that UpdateIngredientTagMapping expects the provided input to have a valid ID.
func (p *Postgres) UpdateIngredientTagMapping(ctx context.Context, input *models.IngredientTagMapping) error {
	query, args := p.buildUpdateIngredientTagMappingQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveIngredientTagMappingQuery returns a SQL query which marks a given ingredient tag mapping belonging to a given valid ingredient as archived.
func (p *Postgres) buildArchiveIngredientTagMappingQuery(validIngredientID, ingredientTagMappingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(ingredientTagMappingsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          ingredientTagMappingID,
			"archived_on": nil,
			ingredientTagMappingsTableOwnershipColumn: validIngredientID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveIngredientTagMapping marks an ingredient tag mapping as archived in the database.
func (p *Postgres) ArchiveIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) error {
	query, args := p.buildArchiveIngredientTagMappingQuery(validIngredientID, ingredientTagMappingID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
