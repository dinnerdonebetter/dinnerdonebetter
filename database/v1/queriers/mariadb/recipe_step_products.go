package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	recipeStepProductsTableName = "recipe_step_products"
)

var (
	recipeStepProductsTableColumns = []string{
		"id",
		"name",
		"recipe_step_id",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanRecipeStepProduct takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Product struct
func scanRecipeStepProduct(scan database.Scanner) (*models.RecipeStepProduct, error) {
	x := &models.RecipeStepProduct{}

	if err := scan.Scan(
		&x.ID,
		&x.Name,
		&x.RecipeStepID,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeStepProducts takes a logger and some database rows and turns them into a slice of recipe step products
func scanRecipeStepProducts(logger logging.Logger, rows *sql.Rows) ([]models.RecipeStepProduct, error) {
	var list []models.RecipeStepProduct

	for rows.Next() {
		x, err := scanRecipeStepProduct(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildGetRecipeStepProductQuery constructs a SQL query for fetching a recipe step product with a given ID belong to a user with a given ID.
func (m *MariaDB) buildGetRecipeStepProductQuery(recipeStepProductID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(recipeStepProductsTableColumns...).
		From(recipeStepProductsTableName).
		Where(squirrel.Eq{
			"id":         recipeStepProductID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepProduct fetches a recipe step product from the mariadb database
func (m *MariaDB) GetRecipeStepProduct(ctx context.Context, recipeStepProductID, userID uint64) (*models.RecipeStepProduct, error) {
	query, args := m.buildGetRecipeStepProductQuery(recipeStepProductID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanRecipeStepProduct(row)
}

// buildGetRecipeStepProductCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of recipe step products belonging to a given user that meet a given query
func (m *MariaDB) buildGetRecipeStepProductCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(CountQuery).
		From(recipeStepProductsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepProductCount will fetch the count of recipe step products from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetRecipeStepProductCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetRecipeStepProductCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allRecipeStepProductsCountQueryBuilder sync.Once
	allRecipeStepProductsCountQuery        string
)

// buildGetAllRecipeStepProductsCountQuery returns a query that fetches the total number of recipe step products in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllRecipeStepProductsCountQuery() string {
	allRecipeStepProductsCountQueryBuilder.Do(func() {
		var err error
		allRecipeStepProductsCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(recipeStepProductsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allRecipeStepProductsCountQuery
}

// GetAllRecipeStepProductsCount will fetch the count of recipe step products from the database
func (m *MariaDB) GetAllRecipeStepProductsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllRecipeStepProductsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeStepProductsQuery builds a SQL query selecting recipe step products that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetRecipeStepProductsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(recipeStepProductsTableColumns...).
		From(recipeStepProductsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter
func (m *MariaDB) GetRecipeStepProducts(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepProductList, error) {
	query, args := m.buildGetRecipeStepProductsQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step products")
	}

	list, err := scanRecipeStepProducts(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetRecipeStepProductCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching recipe step product count: %w", err)
	}

	x := &models.RecipeStepProductList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeStepProducts: list,
	}

	return x, nil
}

// GetAllRecipeStepProductsForUser fetches every recipe step product belonging to a user
func (m *MariaDB) GetAllRecipeStepProductsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepProduct, error) {
	query, args := m.buildGetRecipeStepProductsQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching recipe step products for user")
	}

	list, err := scanRecipeStepProducts(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRecipeStepProductQuery takes a recipe step product and returns a creation query for that recipe step product and the relevant arguments.
func (m *MariaDB) buildCreateRecipeStepProductQuery(input *models.RecipeStepProduct) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(recipeStepProductsTableName).
		Columns(
			"name",
			"recipe_step_id",
			"belongs_to",
			"created_on",
		).
		Values(
			input.Name,
			input.RecipeStepID,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildRecipeStepProductCreationTimeQuery takes a recipe step product and returns a creation query for that recipe step product and the relevant arguments
func (m *MariaDB) buildRecipeStepProductCreationTimeQuery(recipeStepProductID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(recipeStepProductsTableName).
		Where(squirrel.Eq{"id": recipeStepProductID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepProduct creates a recipe step product in the database
func (m *MariaDB) CreateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProductCreationInput) (*models.RecipeStepProduct, error) {
	x := &models.RecipeStepProduct{
		Name:         input.Name,
		RecipeStepID: input.RecipeStepID,
		BelongsTo:    input.BelongsTo,
	}

	query, args := m.buildCreateRecipeStepProductQuery(x)

	// create the recipe step product
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step product creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildRecipeStepProductCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateRecipeStepProductQuery takes a recipe step product and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdateRecipeStepProductQuery(input *models.RecipeStepProduct) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(recipeStepProductsTableName).
		Set("name", input.Name).
		Set("recipe_step_id", input.RecipeStepID).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepProduct updates a particular recipe step product. Note that UpdateRecipeStepProduct expects the provided input to have a valid ID.
func (m *MariaDB) UpdateRecipeStepProduct(ctx context.Context, input *models.RecipeStepProduct) error {
	query, args := m.buildUpdateRecipeStepProductQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveRecipeStepProductQuery returns a SQL query which marks a given recipe step product belonging to a given user as archived.
func (m *MariaDB) buildArchiveRecipeStepProductQuery(recipeStepProductID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(recipeStepProductsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeStepProductID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepProduct marks a recipe step product as archived in the database
func (m *MariaDB) ArchiveRecipeStepProduct(ctx context.Context, recipeStepProductID, userID uint64) error {
	query, args := m.buildArchiveRecipeStepProductQuery(recipeStepProductID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
