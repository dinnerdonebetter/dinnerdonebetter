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
	recipeIterationStepsTableName            = "recipe_iteration_steps"
	recipeIterationStepsTableOwnershipColumn = "belongs_to_recipe"
)

var (
	recipeIterationStepsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeIterationStepsTableName, "id"),
		fmt.Sprintf("%s.%s", recipeIterationStepsTableName, "started_on"),
		fmt.Sprintf("%s.%s", recipeIterationStepsTableName, "ended_on"),
		fmt.Sprintf("%s.%s", recipeIterationStepsTableName, "state"),
		fmt.Sprintf("%s.%s", recipeIterationStepsTableName, "created_on"),
		fmt.Sprintf("%s.%s", recipeIterationStepsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", recipeIterationStepsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", recipeIterationStepsTableName, recipeIterationStepsTableOwnershipColumn),
	}
)

// scanRecipeIterationStep takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Iteration Step struct
func (p *Postgres) scanRecipeIterationStep(scan database.Scanner, includeCount bool) (*models.RecipeIterationStep, uint64, error) {
	x := &models.RecipeIterationStep{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.StartedOn,
		&x.EndedOn,
		&x.State,
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

// scanRecipeIterationSteps takes a logger and some database rows and turns them into a slice of recipe iteration steps.
func (p *Postgres) scanRecipeIterationSteps(rows database.ResultIterator) ([]models.RecipeIterationStep, uint64, error) {
	var (
		list  []models.RecipeIterationStep
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanRecipeIterationStep(rows, true)
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

// buildRecipeIterationStepExistsQuery constructs a SQL query for checking if a recipe iteration step with a given ID belong to a a recipe with a given ID exists
func (p *Postgres) buildRecipeIterationStepExistsQuery(recipeID, recipeIterationStepID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", recipeIterationStepsTableName)).
		Prefix(existencePrefix).
		From(recipeIterationStepsTableName).
		Join(recipesOnRecipeIterationStepsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeIterationStepsTableName):                                           recipeIterationStepID,
			fmt.Sprintf("%s.id", recipesTableName):                                                        recipeID,
			fmt.Sprintf("%s.%s", recipeIterationStepsTableName, recipeIterationStepsTableOwnershipColumn): recipeID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeIterationStepExists queries the database to see if a given recipe iteration step belonging to a given user exists.
func (p *Postgres) RecipeIterationStepExists(ctx context.Context, recipeID, recipeIterationStepID uint64) (exists bool, err error) {
	query, args := p.buildRecipeIterationStepExistsQuery(recipeID, recipeIterationStepID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRecipeIterationStepQuery constructs a SQL query for fetching a recipe iteration step with a given ID belong to a recipe with a given ID.
func (p *Postgres) buildGetRecipeIterationStepQuery(recipeID, recipeIterationStepID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipeIterationStepsTableColumns...).
		From(recipeIterationStepsTableName).
		Join(recipesOnRecipeIterationStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeIterationStepsTableName):                                           recipeIterationStepID,
			fmt.Sprintf("%s.id", recipesTableName):                                                        recipeID,
			fmt.Sprintf("%s.%s", recipeIterationStepsTableName, recipeIterationStepsTableOwnershipColumn): recipeID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIterationStep fetches a recipe iteration step from the database.
func (p *Postgres) GetRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) (*models.RecipeIterationStep, error) {
	query, args := p.buildGetRecipeIterationStepQuery(recipeID, recipeIterationStepID)
	row := p.db.QueryRowContext(ctx, query, args...)

	recipeIterationStep, _, err := p.scanRecipeIterationStep(row, false)
	return recipeIterationStep, err
}

var (
	allRecipeIterationStepsCountQueryBuilder sync.Once
	allRecipeIterationStepsCountQuery        string
)

// buildGetAllRecipeIterationStepsCountQuery returns a query that fetches the total number of recipe iteration steps in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipeIterationStepsCountQuery() string {
	allRecipeIterationStepsCountQueryBuilder.Do(func() {
		var err error

		allRecipeIterationStepsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipeIterationStepsTableName)).
			From(recipeIterationStepsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", recipeIterationStepsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeIterationStepsCountQuery
}

// GetAllRecipeIterationStepsCount will fetch the count of recipe iteration steps from the database.
func (p *Postgres) GetAllRecipeIterationStepsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeIterationStepsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeIterationStepsQuery builds a SQL query selecting recipe iteration steps that adhere to a given QueryFilter and belong to a given recipe,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeIterationStepsQuery(recipeID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(recipeIterationStepsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllRecipeIterationStepsCountQuery()))...).
		From(recipeIterationStepsTableName).
		Join(recipesOnRecipeIterationStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", recipeIterationStepsTableName):                                  nil,
			fmt.Sprintf("%s.id", recipesTableName):                                                        recipeID,
			fmt.Sprintf("%s.%s", recipeIterationStepsTableName, recipeIterationStepsTableOwnershipColumn): recipeID,
		}).
		OrderBy(fmt.Sprintf("%s.id", recipeIterationStepsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipeIterationStepsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIterationSteps fetches a list of recipe iteration steps from the database that meet a particular filter.
func (p *Postgres) GetRecipeIterationSteps(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeIterationStepList, error) {
	query, args := p.buildGetRecipeIterationStepsQuery(recipeID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe iteration steps")
	}

	recipeIterationSteps, count, err := p.scanRecipeIterationSteps(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeIterationStepList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeIterationSteps: recipeIterationSteps,
	}

	return list, nil
}

// buildCreateRecipeIterationStepQuery takes a recipe iteration step and returns a creation query for that recipe iteration step and the relevant arguments.
func (p *Postgres) buildCreateRecipeIterationStepQuery(input *models.RecipeIterationStep) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeIterationStepsTableName).
		Columns(
			"started_on",
			"ended_on",
			"state",
			recipeIterationStepsTableOwnershipColumn,
		).
		Values(
			input.StartedOn,
			input.EndedOn,
			input.State,
			input.BelongsToRecipe,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeIterationStep creates a recipe iteration step in the database.
func (p *Postgres) CreateRecipeIterationStep(ctx context.Context, input *models.RecipeIterationStepCreationInput) (*models.RecipeIterationStep, error) {
	x := &models.RecipeIterationStep{
		StartedOn:       input.StartedOn,
		EndedOn:         input.EndedOn,
		State:           input.State,
		BelongsToRecipe: input.BelongsToRecipe,
	}

	query, args := p.buildCreateRecipeIterationStepQuery(x)

	// create the recipe iteration step.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe iteration step creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeIterationStepQuery takes a recipe iteration step and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeIterationStepQuery(input *models.RecipeIterationStep) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeIterationStepsTableName).
		Set("started_on", input.StartedOn).
		Set("ended_on", input.EndedOn).
		Set("state", input.State).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                                     input.ID,
			recipeIterationStepsTableOwnershipColumn: input.BelongsToRecipe,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeIterationStep updates a particular recipe iteration step. Note that UpdateRecipeIterationStep expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeIterationStep(ctx context.Context, input *models.RecipeIterationStep) error {
	query, args := p.buildUpdateRecipeIterationStepQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveRecipeIterationStepQuery returns a SQL query which marks a given recipe iteration step belonging to a given recipe as archived.
func (p *Postgres) buildArchiveRecipeIterationStepQuery(recipeID, recipeIterationStepID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeIterationStepsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                                     recipeIterationStepID,
			"archived_on":                            nil,
			recipeIterationStepsTableOwnershipColumn: recipeID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeIterationStep marks a recipe iteration step as archived in the database.
func (p *Postgres) ArchiveRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) error {
	query, args := p.buildArchiveRecipeIterationStepQuery(recipeID, recipeIterationStepID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
