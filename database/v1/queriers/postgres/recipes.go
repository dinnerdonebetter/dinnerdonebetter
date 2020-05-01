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
	recipesTableName           = "recipes"
	recipesUserOwnershipColumn = "belongs_to_user"
)

var (
	recipesTableColumns = []string{
		fmt.Sprintf("%s.%s", recipesTableName, "id"),
		fmt.Sprintf("%s.%s", recipesTableName, "name"),
		fmt.Sprintf("%s.%s", recipesTableName, "source"),
		fmt.Sprintf("%s.%s", recipesTableName, "description"),
		fmt.Sprintf("%s.%s", recipesTableName, "inspired_by_recipe_id"),
		fmt.Sprintf("%s.%s", recipesTableName, "private"),
		fmt.Sprintf("%s.%s", recipesTableName, "created_on"),
		fmt.Sprintf("%s.%s", recipesTableName, "updated_on"),
		fmt.Sprintf("%s.%s", recipesTableName, "archived_on"),
		fmt.Sprintf("%s.%s", recipesTableName, recipesUserOwnershipColumn),
	}

	recipesOnRecipeTagsJoinClause           = fmt.Sprintf("%s ON %s.%s=%s.id", recipesTableName, recipeTagsTableName, recipeTagsTableOwnershipColumn, recipesTableName)
	recipesOnRecipeStepsJoinClause          = fmt.Sprintf("%s ON %s.%s=%s.id", recipesTableName, recipeStepsTableName, recipeStepsTableOwnershipColumn, recipesTableName)
	recipesOnRecipeIterationsJoinClause     = fmt.Sprintf("%s ON %s.%s=%s.id", recipesTableName, recipeIterationsTableName, recipeIterationsTableOwnershipColumn, recipesTableName)
	recipesOnRecipeIterationStepsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.id", recipesTableName, recipeIterationStepsTableName, recipeIterationStepsTableOwnershipColumn, recipesTableName)
)

// scanRecipe takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe struct
func (p *Postgres) scanRecipe(scan database.Scanner, includeCount bool) (*models.Recipe, uint64, error) {
	x := &models.Recipe{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Source,
		&x.Description,
		&x.InspiredByRecipeID,
		&x.Private,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}

// scanRecipes takes a logger and some database rows and turns them into a slice of recipes.
func (p *Postgres) scanRecipes(rows database.ResultIterator) ([]models.Recipe, uint64, error) {
	var (
		list  []models.Recipe
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanRecipe(rows, true)
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

//
func (p *Postgres) buildRecipeExistsQuery(recipeID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", recipesTableName)).
		Prefix(existencePrefix).
		From(recipesTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipesTableName): recipeID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeExists queries the database to see if a given recipe belonging to a given user exists.
func (p *Postgres) RecipeExists(ctx context.Context, recipeID uint64) (exists bool, err error) {
	query, args := p.buildRecipeExistsQuery(recipeID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

//
func (p *Postgres) buildGetRecipeQuery(recipeID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipesTableColumns...).
		From(recipesTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipesTableName): recipeID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipe fetches a recipe from the database.
func (p *Postgres) GetRecipe(ctx context.Context, recipeID uint64) (*models.Recipe, error) {
	query, args := p.buildGetRecipeQuery(recipeID)
	row := p.db.QueryRowContext(ctx, query, args...)

	recipe, _, err := p.scanRecipe(row, false)
	return recipe, err
}

var (
	allRecipesCountQueryBuilder sync.Once
	allRecipesCountQuery        string
)

// buildGetAllRecipesCountQuery returns a query that fetches the total number of recipes in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipesCountQuery() string {
	allRecipesCountQueryBuilder.Do(func() {
		var err error

		allRecipesCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipesTableName)).
			From(recipesTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", recipesTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipesCountQuery
}

// GetAllRecipesCount will fetch the count of recipes from the database.
func (p *Postgres) GetAllRecipesCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipesCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipesQuery builds a SQL query selecting recipes that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipesQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(recipesTableColumns, fmt.Sprintf(countQuery, recipesTableName))...).
		From(recipesTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", recipesTableName): nil,
		}).
		GroupBy(fmt.Sprintf("%s.id", recipesTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipesTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipes fetches a list of recipes from the database that meet a particular filter.
func (p *Postgres) GetRecipes(ctx context.Context, filter *models.QueryFilter) (*models.RecipeList, error) {
	query, args := p.buildGetRecipesQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipes")
	}

	recipes, count, err := p.scanRecipes(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Recipes: recipes,
	}

	return list, nil
}

// buildCreateRecipeQuery takes a recipe and returns a creation query for that recipe and the relevant arguments.
func (p *Postgres) buildCreateRecipeQuery(input *models.Recipe) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipesTableName).
		Columns(
			"name",
			"source",
			"description",
			"inspired_by_recipe_id",
			"private",
			recipesUserOwnershipColumn,
		).
		Values(
			input.Name,
			input.Source,
			input.Description,
			input.InspiredByRecipeID,
			input.Private,
			input.BelongsToUser,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipe creates a recipe in the database.
func (p *Postgres) CreateRecipe(ctx context.Context, input *models.RecipeCreationInput) (*models.Recipe, error) {
	x := &models.Recipe{
		Name:               input.Name,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		Private:            input.Private,
		BelongsToUser:      input.BelongsToUser,
	}

	query, args := p.buildCreateRecipeQuery(x)

	// create the recipe.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeQuery takes a recipe and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeQuery(input *models.Recipe) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipesTableName).
		Set("name", input.Name).
		Set("source", input.Source).
		Set("description", input.Description).
		Set("inspired_by_recipe_id", input.InspiredByRecipeID).
		Set("private", input.Private).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                       input.ID,
			recipesUserOwnershipColumn: input.BelongsToUser,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipe updates a particular recipe. Note that UpdateRecipe expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipe(ctx context.Context, input *models.Recipe) error {
	query, args := p.buildUpdateRecipeQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveRecipeQuery returns a SQL query which marks a given recipe belonging to a given user as archived.
func (p *Postgres) buildArchiveRecipeQuery(recipeID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipesTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                       recipeID,
			"archived_on":              nil,
			recipesUserOwnershipColumn: userID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipe marks a recipe as archived in the database.
func (p *Postgres) ArchiveRecipe(ctx context.Context, recipeID, userID uint64) error {
	query, args := p.buildArchiveRecipeQuery(recipeID, userID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
