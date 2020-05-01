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
	recipeIterationsTableName            = "recipe_iterations"
	recipeIterationsTableOwnershipColumn = "belongs_to_recipe"
)

var (
	recipeIterationsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeIterationsTableName, "id"),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, "end_difficulty_rating"),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, "end_complexity_rating"),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, "end_taste_rating"),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, "end_overall_rating"),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, "created_on"),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn),
	}

	recipeIterationsOnIterationMediasJoinClause = fmt.Sprintf("%s ON %s.%s=%s.id", recipeIterationsTableName, iterationMediasTableName, iterationMediasTableOwnershipColumn, recipeIterationsTableName)
)

// scanRecipeIteration takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Iteration struct
func (p *Postgres) scanRecipeIteration(scan database.Scanner, includeCount bool) (*models.RecipeIteration, uint64, error) {
	x := &models.RecipeIteration{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.EndDifficultyRating,
		&x.EndComplexityRating,
		&x.EndTasteRating,
		&x.EndOverallRating,
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

// scanRecipeIterations takes a logger and some database rows and turns them into a slice of recipe iterations.
func (p *Postgres) scanRecipeIterations(rows database.ResultIterator) ([]models.RecipeIteration, uint64, error) {
	var (
		list  []models.RecipeIteration
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanRecipeIteration(rows, true)
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

// buildRecipeIterationExistsQuery constructs a SQL query for checking if a recipe iteration with a given ID belong to a a recipe with a given ID exists
func (p *Postgres) buildRecipeIterationExistsQuery(recipeID, recipeIterationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", recipeIterationsTableName)).
		Prefix(existencePrefix).
		From(recipeIterationsTableName).
		Join(recipesOnRecipeIterationsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeIterationsTableName):                                       recipeIterationID,
			fmt.Sprintf("%s.id", recipesTableName):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeIterationExists queries the database to see if a given recipe iteration belonging to a given user exists.
func (p *Postgres) RecipeIterationExists(ctx context.Context, recipeID, recipeIterationID uint64) (exists bool, err error) {
	query, args := p.buildRecipeIterationExistsQuery(recipeID, recipeIterationID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRecipeIterationQuery constructs a SQL query for fetching a recipe iteration with a given ID belong to a recipe with a given ID.
func (p *Postgres) buildGetRecipeIterationQuery(recipeID, recipeIterationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipeIterationsTableColumns...).
		From(recipeIterationsTableName).
		Join(recipesOnRecipeIterationsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeIterationsTableName):                                       recipeIterationID,
			fmt.Sprintf("%s.id", recipesTableName):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIteration fetches a recipe iteration from the database.
func (p *Postgres) GetRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) (*models.RecipeIteration, error) {
	query, args := p.buildGetRecipeIterationQuery(recipeID, recipeIterationID)
	row := p.db.QueryRowContext(ctx, query, args...)

	recipeIteration, _, err := p.scanRecipeIteration(row, false)
	return recipeIteration, err
}

var (
	allRecipeIterationsCountQueryBuilder sync.Once
	allRecipeIterationsCountQuery        string
)

// buildGetAllRecipeIterationsCountQuery returns a query that fetches the total number of recipe iterations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipeIterationsCountQuery() string {
	allRecipeIterationsCountQueryBuilder.Do(func() {
		var err error

		allRecipeIterationsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipeIterationsTableName)).
			From(recipeIterationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", recipeIterationsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeIterationsCountQuery
}

// GetAllRecipeIterationsCount will fetch the count of recipe iterations from the database.
func (p *Postgres) GetAllRecipeIterationsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeIterationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeIterationsQuery builds a SQL query selecting recipe iterations that adhere to a given QueryFilter and belong to a given recipe,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeIterationsQuery(recipeID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(recipeIterationsTableColumns, fmt.Sprintf(countQuery, recipeIterationsTableName))...).
		From(recipeIterationsTableName).
		Join(recipesOnRecipeIterationsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", recipeIterationsTableName):                              nil,
			fmt.Sprintf("%s.id", recipesTableName):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
		}).
		GroupBy(fmt.Sprintf("%s.id", recipeIterationsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipeIterationsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIterations fetches a list of recipe iterations from the database that meet a particular filter.
func (p *Postgres) GetRecipeIterations(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeIterationList, error) {
	query, args := p.buildGetRecipeIterationsQuery(recipeID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe iterations")
	}

	recipeIterations, count, err := p.scanRecipeIterations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeIterationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeIterations: recipeIterations,
	}

	return list, nil
}

// buildCreateRecipeIterationQuery takes a recipe iteration and returns a creation query for that recipe iteration and the relevant arguments.
func (p *Postgres) buildCreateRecipeIterationQuery(input *models.RecipeIteration) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeIterationsTableName).
		Columns(
			"end_difficulty_rating",
			"end_complexity_rating",
			"end_taste_rating",
			"end_overall_rating",
			recipeIterationsTableOwnershipColumn,
		).
		Values(
			input.EndDifficultyRating,
			input.EndComplexityRating,
			input.EndTasteRating,
			input.EndOverallRating,
			input.BelongsToRecipe,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeIteration creates a recipe iteration in the database.
func (p *Postgres) CreateRecipeIteration(ctx context.Context, input *models.RecipeIterationCreationInput) (*models.RecipeIteration, error) {
	x := &models.RecipeIteration{
		EndDifficultyRating: input.EndDifficultyRating,
		EndComplexityRating: input.EndComplexityRating,
		EndTasteRating:      input.EndTasteRating,
		EndOverallRating:    input.EndOverallRating,
		BelongsToRecipe:     input.BelongsToRecipe,
	}

	query, args := p.buildCreateRecipeIterationQuery(x)

	// create the recipe iteration.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe iteration creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeIterationQuery takes a recipe iteration and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeIterationQuery(input *models.RecipeIteration) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeIterationsTableName).
		Set("end_difficulty_rating", input.EndDifficultyRating).
		Set("end_complexity_rating", input.EndComplexityRating).
		Set("end_taste_rating", input.EndTasteRating).
		Set("end_overall_rating", input.EndOverallRating).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                                 input.ID,
			recipeIterationsTableOwnershipColumn: input.BelongsToRecipe,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeIteration updates a particular recipe iteration. Note that UpdateRecipeIteration expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeIteration(ctx context.Context, input *models.RecipeIteration) error {
	query, args := p.buildUpdateRecipeIterationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveRecipeIterationQuery returns a SQL query which marks a given recipe iteration belonging to a given recipe as archived.
func (p *Postgres) buildArchiveRecipeIterationQuery(recipeID, recipeIterationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeIterationsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                                 recipeIterationID,
			"archived_on":                        nil,
			recipeIterationsTableOwnershipColumn: recipeID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeIteration marks a recipe iteration as archived in the database.
func (p *Postgres) ArchiveRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) error {
	query, args := p.buildArchiveRecipeIterationQuery(recipeID, recipeIterationID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
