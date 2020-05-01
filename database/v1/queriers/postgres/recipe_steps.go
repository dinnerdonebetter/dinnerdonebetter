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
	recipeStepsTableName            = "recipe_steps"
	recipeStepsTableOwnershipColumn = "belongs_to_recipe"
)

var (
	recipeStepsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeStepsTableName, "id"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "index"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "valid_preparation_id"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "prerequisite_step_id"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "min_estimated_time_in_seconds"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "max_estimated_time_in_seconds"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "yields_product_name"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "yields_quantity"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "notes"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "created_on"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn),
	}

	recipeStepsOnRecipeStepPreparationsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.id", recipeStepsTableName, recipeStepPreparationsTableName, recipeStepPreparationsTableOwnershipColumn, recipeStepsTableName)
	recipeStepsOnRecipeStepIngredientsJoinClause  = fmt.Sprintf("%s ON %s.%s=%s.id", recipeStepsTableName, recipeStepIngredientsTableName, recipeStepIngredientsTableOwnershipColumn, recipeStepsTableName)
)

// scanRecipeStep takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step struct
func (p *Postgres) scanRecipeStep(scan database.Scanner, includeCount bool) (*models.RecipeStep, uint64, error) {
	x := &models.RecipeStep{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Index,
		&x.ValidPreparationID,
		&x.PrerequisiteStepID,
		&x.MinEstimatedTimeInSeconds,
		&x.MaxEstimatedTimeInSeconds,
		&x.YieldsProductName,
		&x.YieldsQuantity,
		&x.Notes,
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

// scanRecipeSteps takes a logger and some database rows and turns them into a slice of recipe steps.
func (p *Postgres) scanRecipeSteps(rows database.ResultIterator) ([]models.RecipeStep, uint64, error) {
	var (
		list  []models.RecipeStep
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanRecipeStep(rows, true)
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

// buildRecipeStepExistsQuery constructs a SQL query for checking if a recipe step with a given ID belong to a a recipe with a given ID exists
func (p *Postgres) buildRecipeStepExistsQuery(recipeID, recipeStepID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", recipeStepsTableName)).
		Prefix(existencePrefix).
		From(recipeStepsTableName).
		Join(recipesOnRecipeStepsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeStepsTableName):                                  recipeStepID,
			fmt.Sprintf("%s.id", recipesTableName):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn): recipeID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeStepExists queries the database to see if a given recipe step belonging to a given user exists.
func (p *Postgres) RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (exists bool, err error) {
	query, args := p.buildRecipeStepExistsQuery(recipeID, recipeStepID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRecipeStepQuery constructs a SQL query for fetching a recipe step with a given ID belong to a recipe with a given ID.
func (p *Postgres) buildGetRecipeStepQuery(recipeID, recipeStepID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipeStepsTableColumns...).
		From(recipeStepsTableName).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeStepsTableName):                                  recipeStepID,
			fmt.Sprintf("%s.id", recipesTableName):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn): recipeID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStep fetches a recipe step from the database.
func (p *Postgres) GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (*models.RecipeStep, error) {
	query, args := p.buildGetRecipeStepQuery(recipeID, recipeStepID)
	row := p.db.QueryRowContext(ctx, query, args...)

	recipeStep, _, err := p.scanRecipeStep(row, false)
	return recipeStep, err
}

var (
	allRecipeStepsCountQueryBuilder sync.Once
	allRecipeStepsCountQuery        string
)

// buildGetAllRecipeStepsCountQuery returns a query that fetches the total number of recipe steps in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipeStepsCountQuery() string {
	allRecipeStepsCountQueryBuilder.Do(func() {
		var err error

		allRecipeStepsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipeStepsTableName)).
			From(recipeStepsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", recipeStepsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeStepsCountQuery
}

// GetAllRecipeStepsCount will fetch the count of recipe steps from the database.
func (p *Postgres) GetAllRecipeStepsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeStepsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeStepsQuery builds a SQL query selecting recipe steps that adhere to a given QueryFilter and belong to a given recipe,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeStepsQuery(recipeID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(recipeStepsTableColumns, fmt.Sprintf(countQuery, recipeStepsTableName))...).
		From(recipeStepsTableName).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", recipeStepsTableName):                         nil,
			fmt.Sprintf("%s.id", recipesTableName):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn): recipeID,
		}).
		GroupBy(fmt.Sprintf("%s.id", recipeStepsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipeStepsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeSteps fetches a list of recipe steps from the database that meet a particular filter.
func (p *Postgres) GetRecipeSteps(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*models.RecipeStepList, error) {
	query, args := p.buildGetRecipeStepsQuery(recipeID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe steps")
	}

	recipeSteps, count, err := p.scanRecipeSteps(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeStepList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeSteps: recipeSteps,
	}

	return list, nil
}

// buildCreateRecipeStepQuery takes a recipe step and returns a creation query for that recipe step and the relevant arguments.
func (p *Postgres) buildCreateRecipeStepQuery(input *models.RecipeStep) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeStepsTableName).
		Columns(
			"index",
			"valid_preparation_id",
			"prerequisite_step_id",
			"min_estimated_time_in_seconds",
			"max_estimated_time_in_seconds",
			"yields_product_name",
			"yields_quantity",
			"notes",
			recipeStepsTableOwnershipColumn,
		).
		Values(
			input.Index,
			input.ValidPreparationID,
			input.PrerequisiteStepID,
			input.MinEstimatedTimeInSeconds,
			input.MaxEstimatedTimeInSeconds,
			input.YieldsProductName,
			input.YieldsQuantity,
			input.Notes,
			input.BelongsToRecipe,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStep creates a recipe step in the database.
func (p *Postgres) CreateRecipeStep(ctx context.Context, input *models.RecipeStepCreationInput) (*models.RecipeStep, error) {
	x := &models.RecipeStep{
		Index:                     input.Index,
		ValidPreparationID:        input.ValidPreparationID,
		PrerequisiteStepID:        input.PrerequisiteStepID,
		MinEstimatedTimeInSeconds: input.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: input.MaxEstimatedTimeInSeconds,
		YieldsProductName:         input.YieldsProductName,
		YieldsQuantity:            input.YieldsQuantity,
		Notes:                     input.Notes,
		BelongsToRecipe:           input.BelongsToRecipe,
	}

	query, args := p.buildCreateRecipeStepQuery(x)

	// create the recipe step.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeStepQuery takes a recipe step and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeStepQuery(input *models.RecipeStep) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepsTableName).
		Set("index", input.Index).
		Set("valid_preparation_id", input.ValidPreparationID).
		Set("prerequisite_step_id", input.PrerequisiteStepID).
		Set("min_estimated_time_in_seconds", input.MinEstimatedTimeInSeconds).
		Set("max_estimated_time_in_seconds", input.MaxEstimatedTimeInSeconds).
		Set("yields_product_name", input.YieldsProductName).
		Set("yields_quantity", input.YieldsQuantity).
		Set("notes", input.Notes).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                            input.ID,
			recipeStepsTableOwnershipColumn: input.BelongsToRecipe,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStep updates a particular recipe step. Note that UpdateRecipeStep expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeStep(ctx context.Context, input *models.RecipeStep) error {
	query, args := p.buildUpdateRecipeStepQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveRecipeStepQuery returns a SQL query which marks a given recipe step belonging to a given recipe as archived.
func (p *Postgres) buildArchiveRecipeStepQuery(recipeID, recipeStepID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                            recipeStepID,
			"archived_on":                   nil,
			recipeStepsTableOwnershipColumn: recipeID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStep marks a recipe step as archived in the database.
func (p *Postgres) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) error {
	query, args := p.buildArchiveRecipeStepQuery(recipeID, recipeStepID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
