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
	recipeStepProductsTableName               = "recipe_step_products"
	recipeStepProductsTableNameColumn         = "name"
	recipeStepProductsTableRecipeStepIDColumn = "recipe_step_id"
	recipeStepProductsTableOwnershipColumn    = "belongs_to_recipe_step"
)

var (
	recipeStepProductsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeStepProductsTableName, idColumn),
		fmt.Sprintf("%s.%s", recipeStepProductsTableName, recipeStepProductsTableNameColumn),
		fmt.Sprintf("%s.%s", recipeStepProductsTableName, recipeStepProductsTableRecipeStepIDColumn),
		fmt.Sprintf("%s.%s", recipeStepProductsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", recipeStepProductsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepProductsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepProductsTableName, recipeStepProductsTableOwnershipColumn),
	}
)

// scanRecipeStepProduct takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Product struct
func (p *Postgres) scanRecipeStepProduct(scan database.Scanner) (*models.RecipeStepProduct, error) {
	x := &models.RecipeStepProduct{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.RecipeStepID,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipeStep,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeStepProducts takes a logger and some database rows and turns them into a slice of recipe step products.
func (p *Postgres) scanRecipeStepProducts(rows database.ResultIterator) ([]models.RecipeStepProduct, error) {
	var (
		list []models.RecipeStepProduct
	)

	for rows.Next() {
		x, err := p.scanRecipeStepProduct(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		p.logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildRecipeStepProductExistsQuery constructs a SQL query for checking if a recipe step product with a given ID belong to a a recipe step with a given ID exists
func (p *Postgres) buildRecipeStepProductExistsQuery(recipeID, recipeStepID, recipeStepProductID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", recipeStepProductsTableName, idColumn)).
		Prefix(existencePrefix).
		From(recipeStepProductsTableName).
		Join(recipeStepsOnRecipeStepProductsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, idColumn):                               recipeStepProductID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                          recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                      recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):               recipeID,
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, recipeStepProductsTableOwnershipColumn): recipeStepID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeStepProductExists queries the database to see if a given recipe step product belonging to a given user exists.
func (p *Postgres) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (exists bool, err error) {
	query, args := p.buildRecipeStepProductExistsQuery(recipeID, recipeStepID, recipeStepProductID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRecipeStepProductQuery constructs a SQL query for fetching a recipe step product with a given ID belong to a recipe step with a given ID.
func (p *Postgres) buildGetRecipeStepProductQuery(recipeID, recipeStepID, recipeStepProductID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipeStepProductsTableColumns...).
		From(recipeStepProductsTableName).
		Join(recipeStepsOnRecipeStepProductsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, idColumn):                               recipeStepProductID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                          recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                      recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):               recipeID,
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, recipeStepProductsTableOwnershipColumn): recipeStepID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepProduct fetches a recipe step product from the database.
func (p *Postgres) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*models.RecipeStepProduct, error) {
	query, args := p.buildGetRecipeStepProductQuery(recipeID, recipeStepID, recipeStepProductID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanRecipeStepProduct(row)
}

var (
	allRecipeStepProductsCountQueryBuilder sync.Once
	allRecipeStepProductsCountQuery        string
)

// buildGetAllRecipeStepProductsCountQuery returns a query that fetches the total number of recipe step products in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipeStepProductsCountQuery() string {
	allRecipeStepProductsCountQueryBuilder.Do(func() {
		var err error

		allRecipeStepProductsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipeStepProductsTableName)).
			From(recipeStepProductsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", recipeStepProductsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeStepProductsCountQuery
}

// GetAllRecipeStepProductsCount will fetch the count of recipe step products from the database.
func (p *Postgres) GetAllRecipeStepProductsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeStepProductsCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfRecipeStepProductsQuery returns a query that fetches every recipe step product in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfRecipeStepProductsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(recipeStepProductsTableColumns...).
		From(recipeStepProductsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllRecipeStepProducts fetches every recipe step product from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllRecipeStepProducts(ctx context.Context, resultChannel chan []models.RecipeStepProduct) error {
	count, err := p.GetAllRecipeStepProductsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfRecipeStepProductsQuery(begin, end)
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

			recipeStepProducts, err := p.scanRecipeStepProducts(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- recipeStepProducts
		}(beginID, endID)
	}

	return nil
}

// buildGetRecipeStepProductsQuery builds a SQL query selecting recipe step products that adhere to a given QueryFilter and belong to a given recipe step,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeStepProductsQuery(recipeID, recipeStepID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(recipeStepProductsTableColumns...).
		From(recipeStepProductsTableName).
		Join(recipeStepsOnRecipeStepProductsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, archivedOnColumn):                       nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                          recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                      recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):               recipeID,
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, recipeStepProductsTableOwnershipColumn): recipeStepID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", recipeStepProductsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipeStepProductsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter.
func (p *Postgres) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepProductList, error) {
	query, args := p.buildGetRecipeStepProductsQuery(recipeID, recipeStepID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step products")
	}

	recipeStepProducts, err := p.scanRecipeStepProducts(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeStepProductList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		RecipeStepProducts: recipeStepProducts,
	}

	return list, nil
}

// buildGetRecipeStepProductsWithIDsQuery builds a SQL query selecting recipeStepProducts that belong to a given recipe step,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetRecipeStepProductsWithIDsQuery(recipeID, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(recipeStepProductsTableColumns...).
		From(recipeStepProductsTableName).
		Join(recipeStepsOnRecipeStepProductsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, archivedOnColumn):                       nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                          recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                      recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):               recipeID,
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, recipeStepProductsTableOwnershipColumn): recipeStepID,
		}).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(recipeStepProductsTableColumns...).
		FromSelect(subqueryBuilder, recipeStepProductsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepProductsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepProductsWithIDs fetches a list of recipe step products from the database that exist within a given set of IDs.
func (p *Postgres) GetRecipeStepProductsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepProduct, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetRecipeStepProductsWithIDsQuery(recipeID, recipeStepID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step products")
	}

	recipeStepProducts, err := p.scanRecipeStepProducts(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return recipeStepProducts, nil
}

// buildCreateRecipeStepProductQuery takes a recipe step product and returns a creation query for that recipe step product and the relevant arguments.
func (p *Postgres) buildCreateRecipeStepProductQuery(input *models.RecipeStepProduct) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeStepProductsTableName).
		Columns(
			recipeStepProductsTableNameColumn,
			recipeStepProductsTableRecipeStepIDColumn,
			recipeStepProductsTableOwnershipColumn,
		).
		Values(
			input.Name,
			input.RecipeStepID,
			input.BelongsToRecipeStep,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (p *Postgres) CreateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProductCreationInput) (*models.RecipeStepProduct, error) {
	x := &models.RecipeStepProduct{
		Name:                input.Name,
		RecipeStepID:        input.RecipeStepID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
	}

	query, args := p.buildCreateRecipeStepProductQuery(x)

	// create the recipe step product.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step product creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeStepProductQuery takes a recipe step product and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeStepProductQuery(input *models.RecipeStepProduct) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepProductsTableName).
		Set(recipeStepProductsTableNameColumn, input.Name).
		Set(recipeStepProductsTableRecipeStepIDColumn, input.RecipeStepID).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                               input.ID,
			recipeStepProductsTableOwnershipColumn: input.BelongsToRecipeStep,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepProduct updates a particular recipe step product. Note that UpdateRecipeStepProduct expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProduct) error {
	query, args := p.buildUpdateRecipeStepProductQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveRecipeStepProductQuery returns a SQL query which marks a given recipe step product belonging to a given recipe step as archived.
func (p *Postgres) buildArchiveRecipeStepProductQuery(recipeStepID, recipeStepProductID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepProductsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                               recipeStepProductID,
			archivedOnColumn:                       nil,
			recipeStepProductsTableOwnershipColumn: recipeStepID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepProduct marks a recipe step product as archived in the database.
func (p *Postgres) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID uint64) error {
	query, args := p.buildArchiveRecipeStepProductQuery(recipeStepID, recipeStepProductID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
