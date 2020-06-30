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
	recipeTagsTableName            = "recipe_tags"
	recipeTagsTableOwnershipColumn = "belongs_to_recipe"
)

var (
	recipeTagsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeTagsTableName, "id"),
		fmt.Sprintf("%s.%s", recipeTagsTableName, "name"),
		fmt.Sprintf("%s.%s", recipeTagsTableName, "created_on"),
		fmt.Sprintf("%s.%s", recipeTagsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", recipeTagsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", recipeTagsTableName, recipeTagsTableOwnershipColumn),
	}
)

// scanRecipeTag takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Tag struct
func (p *Postgres) scanRecipeTag(scan database.Scanner, includeCount bool) (*models.RecipeTag, uint64, error) {
	x := &models.RecipeTag{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipe,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}

// scanRecipeTags takes a logger and some database rows and turns them into a slice of recipe tags.
func (p *Postgres) scanRecipeTags(rows database.ResultIterator) ([]models.RecipeTag, uint64, error) {
	var (
		list  []models.RecipeTag
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanRecipeTag(rows, true)
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

// buildRecipeTagExistsQuery constructs a SQL query for checking if a recipe tag with a given ID belong to a a recipe with a given ID exists
func (p *Postgres) buildRecipeTagExistsQuery(recipeID, recipeTagID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", recipeTagsTableName)).
		Prefix(existencePrefix).
		From(recipeTagsTableName).
		Join(recipesOnRecipeTagsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeTagsTableName):                                 recipeTagID,
			fmt.Sprintf("%s.id", recipesTableName):                                    recipeID,
			fmt.Sprintf("%s.%s", recipeTagsTableName, recipeTagsTableOwnershipColumn): recipeID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeTagExists queries the database to see if a given recipe tag belonging to a given user exists.
func (p *Postgres) RecipeTagExists(ctx context.Context, recipeID, recipeTagID uint64) (exists bool, err error) {
	query, args := p.buildRecipeTagExistsQuery(recipeID, recipeTagID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRecipeTagQuery constructs a SQL query for fetching a recipe tag with a given ID belong to a recipe with a given ID.
func (p *Postgres) buildGetRecipeTagQuery(recipeID, recipeTagID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipeTagsTableColumns...).
		From(recipeTagsTableName).
		Join(recipesOnRecipeTagsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeTagsTableName):                                 recipeTagID,
			fmt.Sprintf("%s.id", recipesTableName):                                    recipeID,
			fmt.Sprintf("%s.%s", recipeTagsTableName, recipeTagsTableOwnershipColumn): recipeID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeTag fetches a recipe tag from the database.
func (p *Postgres) GetRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) (*models.RecipeTag, error) {
	query, args := p.buildGetRecipeTagQuery(recipeID, recipeTagID)
	row := p.db.QueryRowContext(ctx, query, args...)

	recipeTag, _, err := p.scanRecipeTag(row, false)
	return recipeTag, err
}

var (
	allRecipeTagsCountQueryBuilder sync.Once
	allRecipeTagsCountQuery        string
)

// buildGetAllRecipeTagsCountQuery returns a query that fetches the total number of recipe tags in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipeTagsCountQuery() string {
	allRecipeTagsCountQueryBuilder.Do(func() {
		var err error

		allRecipeTagsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipeTagsTableName)).
			From(recipeTagsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", recipeTagsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeTagsCountQuery
}

// GetAllRecipeTagsCount will fetch the count of recipe tags from the database.
func (p *Postgres) GetAllRecipeTagsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeTagsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeTagsQuery builds a SQL query selecting recipe tags that adhere to a given QueryFilter and belong to a given recipe,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeTagsQuery(recipeID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(recipeTagsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllRecipeTagsCountQuery()))...).
		From(recipeTagsTableName).
		Join(recipesOnRecipeTagsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", recipeTagsTableName):                        nil,
			fmt.Sprintf("%s.id", recipesTableName):                                    recipeID,
			fmt.Sprintf("%s.%s", recipeTagsTableName, recipeTagsTableOwnershipColumn): recipeID,
		}).
		OrderBy(fmt.Sprintf("%s.id", recipeTagsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipeTagsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeTags fetches a list of recipe tags from the database that meet a particular filter.
func (p *Postgres) GetRecipeTags(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeTagList, error) {
	query, args := p.buildGetRecipeTagsQuery(recipeID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe tags")
	}

	recipeTags, count, err := p.scanRecipeTags(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeTagList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeTags: recipeTags,
	}

	return list, nil
}

// buildCreateRecipeTagQuery takes a recipe tag and returns a creation query for that recipe tag and the relevant arguments.
func (p *Postgres) buildCreateRecipeTagQuery(input *models.RecipeTag) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeTagsTableName).
		Columns(
			"name",
			recipeTagsTableOwnershipColumn,
		).
		Values(
			input.Name,
			input.BelongsToRecipe,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeTag creates a recipe tag in the database.
func (p *Postgres) CreateRecipeTag(ctx context.Context, input *models.RecipeTagCreationInput) (*models.RecipeTag, error) {
	x := &models.RecipeTag{
		Name:            input.Name,
		BelongsToRecipe: input.BelongsToRecipe,
	}

	query, args := p.buildCreateRecipeTagQuery(x)

	// create the recipe tag.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe tag creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeTagQuery takes a recipe tag and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeTagQuery(input *models.RecipeTag) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeTagsTableName).
		Set("name", input.Name).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                           input.ID,
			recipeTagsTableOwnershipColumn: input.BelongsToRecipe,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeTag updates a particular recipe tag. Note that UpdateRecipeTag expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeTag(ctx context.Context, input *models.RecipeTag) error {
	query, args := p.buildUpdateRecipeTagQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveRecipeTagQuery returns a SQL query which marks a given recipe tag belonging to a given recipe as archived.
func (p *Postgres) buildArchiveRecipeTagQuery(recipeID, recipeTagID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeTagsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                           recipeTagID,
			"archived_on":                  nil,
			recipeTagsTableOwnershipColumn: recipeID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeTag marks a recipe tag as archived in the database.
func (p *Postgres) ArchiveRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) error {
	query, args := p.buildArchiveRecipeTagQuery(recipeID, recipeTagID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
