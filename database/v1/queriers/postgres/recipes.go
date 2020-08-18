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
	recipesTableName                     = "recipes"
	recipesTableNameColumn               = "name"
	recipesTableSourceColumn             = "source"
	recipesTableDescriptionColumn        = "description"
	recipesTableInspiredByRecipeIDColumn = "inspired_by_recipe_id"
	recipesUserOwnershipColumn           = "belongs_to_user"
)

var (
	recipesTableColumns = []string{
		fmt.Sprintf("%s.%s", recipesTableName, idColumn),
		fmt.Sprintf("%s.%s", recipesTableName, recipesTableNameColumn),
		fmt.Sprintf("%s.%s", recipesTableName, recipesTableSourceColumn),
		fmt.Sprintf("%s.%s", recipesTableName, recipesTableDescriptionColumn),
		fmt.Sprintf("%s.%s", recipesTableName, recipesTableInspiredByRecipeIDColumn),
		fmt.Sprintf("%s.%s", recipesTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", recipesTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", recipesTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", recipesTableName, recipesUserOwnershipColumn),
	}

	recipesOnRecipeStepsJoinClause      = fmt.Sprintf("%s ON %s.%s=%s.%s", recipesTableName, recipeStepsTableName, recipeStepsTableOwnershipColumn, recipesTableName, idColumn)
	recipesOnRecipeIterationsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.%s", recipesTableName, recipeIterationsTableName, recipeIterationsTableOwnershipColumn, recipesTableName, idColumn)
)

// scanRecipe takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe struct
func (p *Postgres) scanRecipe(scan database.Scanner) (*models.Recipe, error) {
	x := &models.Recipe{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Source,
		&x.Description,
		&x.InspiredByRecipeID,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipes takes a logger and some database rows and turns them into a slice of recipes.
func (p *Postgres) scanRecipes(rows database.ResultIterator) ([]models.Recipe, error) {
	var (
		list []models.Recipe
	)

	for rows.Next() {
		x, err := p.scanRecipe(rows)
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

//
func (p *Postgres) buildRecipeExistsQuery(recipeID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", recipesTableName, idColumn)).
		Prefix(existencePrefix).
		From(recipesTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipesTableName, idColumn): recipeID,
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
			fmt.Sprintf("%s.%s", recipesTableName, idColumn): recipeID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipe fetches a recipe from the database.
func (p *Postgres) GetRecipe(ctx context.Context, recipeID uint64) (*models.Recipe, error) {
	query, args := p.buildGetRecipeQuery(recipeID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanRecipe(row)
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
				fmt.Sprintf("%s.%s", recipesTableName, archivedOnColumn): nil,
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

// buildGetBatchOfRecipesQuery returns a query that fetches every recipe in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfRecipesQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(recipesTableColumns...).
		From(recipesTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", recipesTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", recipesTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllRecipes fetches every recipe from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllRecipes(ctx context.Context, resultChannel chan []models.Recipe) error {
	count, err := p.GetAllRecipesCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfRecipesQuery(begin, end)
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

			recipes, err := p.scanRecipes(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- recipes
		}(beginID, endID)
	}

	return nil
}

// buildGetRecipesQuery builds a SQL query selecting recipes that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipesQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(recipesTableColumns...).
		From(recipesTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipesTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", recipesTableName, idColumn))

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

	recipes, err := p.scanRecipes(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Recipes: recipes,
	}

	return list, nil
}

// buildGetRecipesWithIDsQuery builds a SQL query selecting recipes
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetRecipesWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(recipesTableColumns...).
		From(recipesTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(recipesTableColumns...).
		FromSelect(subqueryBuilder, recipesTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipesTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipesWithIDs fetches a list of recipes from the database that exist within a given set of IDs.
func (p *Postgres) GetRecipesWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.Recipe, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetRecipesWithIDsQuery(limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipes")
	}

	recipes, err := p.scanRecipes(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return recipes, nil
}

// buildCreateRecipeQuery takes a recipe and returns a creation query for that recipe and the relevant arguments.
func (p *Postgres) buildCreateRecipeQuery(input *models.Recipe) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipesTableName).
		Columns(
			recipesTableNameColumn,
			recipesTableSourceColumn,
			recipesTableDescriptionColumn,
			recipesTableInspiredByRecipeIDColumn,
			recipesUserOwnershipColumn,
		).
		Values(
			input.Name,
			input.Source,
			input.Description,
			input.InspiredByRecipeID,
			input.BelongsToUser,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
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
		Set(recipesTableNameColumn, input.Name).
		Set(recipesTableSourceColumn, input.Source).
		Set(recipesTableDescriptionColumn, input.Description).
		Set(recipesTableInspiredByRecipeIDColumn, input.InspiredByRecipeID).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                   input.ID,
			recipesUserOwnershipColumn: input.BelongsToUser,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipe updates a particular recipe. Note that UpdateRecipe expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipe(ctx context.Context, input *models.Recipe) error {
	query, args := p.buildUpdateRecipeQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveRecipeQuery returns a SQL query which marks a given recipe belonging to a given user as archived.
func (p *Postgres) buildArchiveRecipeQuery(recipeID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipesTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                   recipeID,
			archivedOnColumn:           nil,
			recipesUserOwnershipColumn: userID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
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
