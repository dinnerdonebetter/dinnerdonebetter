package postgres

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
	recipeIterationsTableName = "recipe_iterations"
)

var (
	recipeIterationsTableColumns = []string{
		"id",
		"recipe_id",
		"end_difficulty_rating",
		"end_complexity_rating",
		"end_taste_rating",
		"end_overall_rating",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanRecipeIteration takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Iteration struct
func scanRecipeIteration(scan database.Scanner) (*models.RecipeIteration, error) {
	x := &models.RecipeIteration{}

	if err := scan.Scan(
		&x.ID,
		&x.RecipeID,
		&x.EndDifficultyRating,
		&x.EndComplexityRating,
		&x.EndTasteRating,
		&x.EndOverallRating,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeIterations takes a logger and some database rows and turns them into a slice of recipe iterations
func scanRecipeIterations(logger logging.Logger, rows *sql.Rows) ([]models.RecipeIteration, error) {
	var list []models.RecipeIteration

	for rows.Next() {
		x, err := scanRecipeIteration(rows)
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

// buildGetRecipeIterationQuery constructs a SQL query for fetching a recipe iteration with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetRecipeIterationQuery(recipeIterationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(recipeIterationsTableColumns...).
		From(recipeIterationsTableName).
		Where(squirrel.Eq{
			"id":         recipeIterationID,
			"belongs_to": userID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIteration fetches a recipe iteration from the postgres database
func (p *Postgres) GetRecipeIteration(ctx context.Context, recipeIterationID, userID uint64) (*models.RecipeIteration, error) {
	query, args := p.buildGetRecipeIterationQuery(recipeIterationID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanRecipeIteration(row)
}

// buildGetRecipeIterationCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of recipe iterations belonging to a given user that meet a given query
func (p *Postgres) buildGetRecipeIterationCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(recipeIterationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIterationCount will fetch the count of recipe iterations from the database that meet a particular filter and belong to a particular user.
func (p *Postgres) GetRecipeIterationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := p.buildGetRecipeIterationCountQuery(filter, userID)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
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
			Select(CountQuery).
			From(recipeIterationsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeIterationsCountQuery
}

// GetAllRecipeIterationsCount will fetch the count of recipe iterations from the database
func (p *Postgres) GetAllRecipeIterationsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeIterationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeIterationsQuery builds a SQL query selecting recipe iterations that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeIterationsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(recipeIterationsTableColumns...).
		From(recipeIterationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIterations fetches a list of recipe iterations from the database that meet a particular filter
func (p *Postgres) GetRecipeIterations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeIterationList, error) {
	query, args := p.buildGetRecipeIterationsQuery(filter, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe iterations")
	}

	list, err := scanRecipeIterations(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetRecipeIterationCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching recipe iteration count: %w", err)
	}

	x := &models.RecipeIterationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeIterations: list,
	}

	return x, nil
}

// GetAllRecipeIterationsForUser fetches every recipe iteration belonging to a user
func (p *Postgres) GetAllRecipeIterationsForUser(ctx context.Context, userID uint64) ([]models.RecipeIteration, error) {
	query, args := p.buildGetRecipeIterationsQuery(nil, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching recipe iterations for user")
	}

	list, err := scanRecipeIterations(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRecipeIterationQuery takes a recipe iteration and returns a creation query for that recipe iteration and the relevant arguments.
func (p *Postgres) buildCreateRecipeIterationQuery(input *models.RecipeIteration) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(recipeIterationsTableName).
		Columns(
			"recipe_id",
			"end_difficulty_rating",
			"end_complexity_rating",
			"end_taste_rating",
			"end_overall_rating",
			"belongs_to",
		).
		Values(
			input.RecipeID,
			input.EndDifficultyRating,
			input.EndComplexityRating,
			input.EndTasteRating,
			input.EndOverallRating,
			input.BelongsTo,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeIteration creates a recipe iteration in the database
func (p *Postgres) CreateRecipeIteration(ctx context.Context, input *models.RecipeIterationCreationInput) (*models.RecipeIteration, error) {
	x := &models.RecipeIteration{
		RecipeID:            input.RecipeID,
		EndDifficultyRating: input.EndDifficultyRating,
		EndComplexityRating: input.EndComplexityRating,
		EndTasteRating:      input.EndTasteRating,
		EndOverallRating:    input.EndOverallRating,
		BelongsTo:           input.BelongsTo,
	}

	query, args := p.buildCreateRecipeIterationQuery(x)

	// create the recipe iteration
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe iteration creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeIterationQuery takes a recipe iteration and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdateRecipeIterationQuery(input *models.RecipeIteration) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(recipeIterationsTableName).
		Set("recipe_id", input.RecipeID).
		Set("end_difficulty_rating", input.EndDifficultyRating).
		Set("end_complexity_rating", input.EndComplexityRating).
		Set("end_taste_rating", input.EndTasteRating).
		Set("end_overall_rating", input.EndOverallRating).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
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

// buildArchiveRecipeIterationQuery returns a SQL query which marks a given recipe iteration belonging to a given user as archived.
func (p *Postgres) buildArchiveRecipeIterationQuery(recipeIterationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(recipeIterationsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeIterationID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeIteration marks a recipe iteration as archived in the database
func (p *Postgres) ArchiveRecipeIteration(ctx context.Context, recipeIterationID, userID uint64) error {
	query, args := p.buildArchiveRecipeIterationQuery(recipeIterationID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
