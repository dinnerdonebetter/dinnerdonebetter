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
	recipesTableName = "recipes"
)

var (
	recipesTableColumns = []string{
		"id",
		"name",
		"source",
		"description",
		"inspired_by_recipe_id",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanRecipe takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe struct
func scanRecipe(scan database.Scanner) (*models.Recipe, error) {
	x := &models.Recipe{}

	if err := scan.Scan(
		&x.ID,
		&x.Name,
		&x.Source,
		&x.Description,
		&x.InspiredByRecipeID,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipes takes a logger and some database rows and turns them into a slice of recipes
func scanRecipes(logger logging.Logger, rows *sql.Rows) ([]models.Recipe, error) {
	var list []models.Recipe

	for rows.Next() {
		x, err := scanRecipe(rows)
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

// buildGetRecipeQuery constructs a SQL query for fetching a recipe with a given ID belong to a user with a given ID.
func (m *MariaDB) buildGetRecipeQuery(recipeID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(recipesTableColumns...).
		From(recipesTableName).
		Where(squirrel.Eq{
			"id":         recipeID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetRecipe fetches a recipe from the mariadb database
func (m *MariaDB) GetRecipe(ctx context.Context, recipeID, userID uint64) (*models.Recipe, error) {
	query, args := m.buildGetRecipeQuery(recipeID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanRecipe(row)
}

// buildGetRecipeCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of recipes belonging to a given user that meet a given query
func (m *MariaDB) buildGetRecipeCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(CountQuery).
		From(recipesTableName).
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

// GetRecipeCount will fetch the count of recipes from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetRecipeCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetRecipeCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allRecipesCountQueryBuilder sync.Once
	allRecipesCountQuery        string
)

// buildGetAllRecipesCountQuery returns a query that fetches the total number of recipes in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllRecipesCountQuery() string {
	allRecipesCountQueryBuilder.Do(func() {
		var err error
		allRecipesCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(recipesTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allRecipesCountQuery
}

// GetAllRecipesCount will fetch the count of recipes from the database
func (m *MariaDB) GetAllRecipesCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllRecipesCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipesQuery builds a SQL query selecting recipes that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetRecipesQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(recipesTableColumns...).
		From(recipesTableName).
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

// GetRecipes fetches a list of recipes from the database that meet a particular filter
func (m *MariaDB) GetRecipes(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeList, error) {
	query, args := m.buildGetRecipesQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipes")
	}

	list, err := scanRecipes(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetRecipeCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching recipe count: %w", err)
	}

	x := &models.RecipeList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Recipes: list,
	}

	return x, nil
}

// GetAllRecipesForUser fetches every recipe belonging to a user
func (m *MariaDB) GetAllRecipesForUser(ctx context.Context, userID uint64) ([]models.Recipe, error) {
	query, args := m.buildGetRecipesQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching recipes for user")
	}

	list, err := scanRecipes(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRecipeQuery takes a recipe and returns a creation query for that recipe and the relevant arguments.
func (m *MariaDB) buildCreateRecipeQuery(input *models.Recipe) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(recipesTableName).
		Columns(
			"name",
			"source",
			"description",
			"inspired_by_recipe_id",
			"belongs_to",
			"created_on",
		).
		Values(
			input.Name,
			input.Source,
			input.Description,
			input.InspiredByRecipeID,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildRecipeCreationTimeQuery takes a recipe and returns a creation query for that recipe and the relevant arguments
func (m *MariaDB) buildRecipeCreationTimeQuery(recipeID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(recipesTableName).
		Where(squirrel.Eq{"id": recipeID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateRecipe creates a recipe in the database
func (m *MariaDB) CreateRecipe(ctx context.Context, input *models.RecipeCreationInput) (*models.Recipe, error) {
	x := &models.Recipe{
		Name:               input.Name,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		BelongsTo:          input.BelongsTo,
	}

	query, args := m.buildCreateRecipeQuery(x)

	// create the recipe
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildRecipeCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateRecipeQuery takes a recipe and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdateRecipeQuery(input *models.Recipe) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(recipesTableName).
		Set("name", input.Name).
		Set("source", input.Source).
		Set("description", input.Description).
		Set("inspired_by_recipe_id", input.InspiredByRecipeID).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipe updates a particular recipe. Note that UpdateRecipe expects the provided input to have a valid ID.
func (m *MariaDB) UpdateRecipe(ctx context.Context, input *models.Recipe) error {
	query, args := m.buildUpdateRecipeQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveRecipeQuery returns a SQL query which marks a given recipe belonging to a given user as archived.
func (m *MariaDB) buildArchiveRecipeQuery(recipeID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(recipesTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipe marks a recipe as archived in the database
func (m *MariaDB) ArchiveRecipe(ctx context.Context, recipeID, userID uint64) error {
	query, args := m.buildArchiveRecipeQuery(recipeID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
