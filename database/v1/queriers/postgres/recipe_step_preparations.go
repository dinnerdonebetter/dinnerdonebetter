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
	recipeStepPreparationsTableName            = "recipe_step_preparations"
	recipeStepPreparationsTableOwnershipColumn = "belongs_to_recipe_step"
)

var (
	recipeStepPreparationsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, "id"),
		fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, "valid_preparation_id"),
		fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, "notes"),
		fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, "created_on"),
		fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, recipeStepPreparationsTableOwnershipColumn),
	}
)

// scanRecipeStepPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Preparation struct
func (p *Postgres) scanRecipeStepPreparation(scan database.Scanner, includeCount bool) (*models.RecipeStepPreparation, uint64, error) {
	x := &models.RecipeStepPreparation{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.ValidPreparationID,
		&x.Notes,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipeStep,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}

// scanRecipeStepPreparations takes a logger and some database rows and turns them into a slice of recipe step preparations.
func (p *Postgres) scanRecipeStepPreparations(rows database.ResultIterator) ([]models.RecipeStepPreparation, uint64, error) {
	var (
		list  []models.RecipeStepPreparation
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanRecipeStepPreparation(rows, true)
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

// buildRecipeStepPreparationExistsQuery constructs a SQL query for checking if a recipe step preparation with a given ID belong to a a recipe step with a given ID exists
func (p *Postgres) buildRecipeStepPreparationExistsQuery(recipeID, recipeStepID, recipeStepPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", recipeStepPreparationsTableName)).
		Prefix(existencePrefix).
		From(recipeStepPreparationsTableName).
		Join(recipeStepsOnRecipeStepPreparationsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeStepPreparationsTableName):                                             recipeStepPreparationID,
			fmt.Sprintf("%s.id", recipesTableName):                                                            recipeID,
			fmt.Sprintf("%s.id", recipeStepsTableName):                                                        recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                       recipeID,
			fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, recipeStepPreparationsTableOwnershipColumn): recipeStepID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeStepPreparationExists queries the database to see if a given recipe step preparation belonging to a given user exists.
func (p *Postgres) RecipeStepPreparationExists(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (exists bool, err error) {
	query, args := p.buildRecipeStepPreparationExistsQuery(recipeID, recipeStepID, recipeStepPreparationID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRecipeStepPreparationQuery constructs a SQL query for fetching a recipe step preparation with a given ID belong to a recipe step with a given ID.
func (p *Postgres) buildGetRecipeStepPreparationQuery(recipeID, recipeStepID, recipeStepPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipeStepPreparationsTableColumns...).
		From(recipeStepPreparationsTableName).
		Join(recipeStepsOnRecipeStepPreparationsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeStepPreparationsTableName):                                             recipeStepPreparationID,
			fmt.Sprintf("%s.id", recipesTableName):                                                            recipeID,
			fmt.Sprintf("%s.id", recipeStepsTableName):                                                        recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                       recipeID,
			fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, recipeStepPreparationsTableOwnershipColumn): recipeStepID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepPreparation fetches a recipe step preparation from the database.
func (p *Postgres) GetRecipeStepPreparation(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (*models.RecipeStepPreparation, error) {
	query, args := p.buildGetRecipeStepPreparationQuery(recipeID, recipeStepID, recipeStepPreparationID)
	row := p.db.QueryRowContext(ctx, query, args...)

	recipeStepPreparation, _, err := p.scanRecipeStepPreparation(row, false)
	return recipeStepPreparation, err
}

var (
	allRecipeStepPreparationsCountQueryBuilder sync.Once
	allRecipeStepPreparationsCountQuery        string
)

// buildGetAllRecipeStepPreparationsCountQuery returns a query that fetches the total number of recipe step preparations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipeStepPreparationsCountQuery() string {
	allRecipeStepPreparationsCountQueryBuilder.Do(func() {
		var err error

		allRecipeStepPreparationsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipeStepPreparationsTableName)).
			From(recipeStepPreparationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", recipeStepPreparationsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeStepPreparationsCountQuery
}

// GetAllRecipeStepPreparationsCount will fetch the count of recipe step preparations from the database.
func (p *Postgres) GetAllRecipeStepPreparationsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeStepPreparationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeStepPreparationsQuery builds a SQL query selecting recipe step preparations that adhere to a given QueryFilter and belong to a given recipe step,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeStepPreparationsQuery(recipeID, recipeStepID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(recipeStepPreparationsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllRecipeStepPreparationsCountQuery()))...).
		From(recipeStepPreparationsTableName).
		Join(recipeStepsOnRecipeStepPreparationsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", recipeStepPreparationsTableName):                                    nil,
			fmt.Sprintf("%s.id", recipesTableName):                                                            recipeID,
			fmt.Sprintf("%s.id", recipeStepsTableName):                                                        recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                       recipeID,
			fmt.Sprintf("%s.%s", recipeStepPreparationsTableName, recipeStepPreparationsTableOwnershipColumn): recipeStepID,
		}).
		OrderBy(fmt.Sprintf("%s.id", recipeStepPreparationsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipeStepPreparationsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepPreparations fetches a list of recipe step preparations from the database that meet a particular filter.
func (p *Postgres) GetRecipeStepPreparations(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepPreparationList, error) {
	query, args := p.buildGetRecipeStepPreparationsQuery(recipeID, recipeStepID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step preparations")
	}

	recipeStepPreparations, count, err := p.scanRecipeStepPreparations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeStepPreparationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeStepPreparations: recipeStepPreparations,
	}

	return list, nil
}

// buildCreateRecipeStepPreparationQuery takes a recipe step preparation and returns a creation query for that recipe step preparation and the relevant arguments.
func (p *Postgres) buildCreateRecipeStepPreparationQuery(input *models.RecipeStepPreparation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeStepPreparationsTableName).
		Columns(
			"valid_preparation_id",
			"notes",
			recipeStepPreparationsTableOwnershipColumn,
		).
		Values(
			input.ValidPreparationID,
			input.Notes,
			input.BelongsToRecipeStep,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepPreparation creates a recipe step preparation in the database.
func (p *Postgres) CreateRecipeStepPreparation(ctx context.Context, input *models.RecipeStepPreparationCreationInput) (*models.RecipeStepPreparation, error) {
	x := &models.RecipeStepPreparation{
		ValidPreparationID:  input.ValidPreparationID,
		Notes:               input.Notes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
	}

	query, args := p.buildCreateRecipeStepPreparationQuery(x)

	// create the recipe step preparation.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step preparation creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeStepPreparationQuery takes a recipe step preparation and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeStepPreparationQuery(input *models.RecipeStepPreparation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepPreparationsTableName).
		Set("valid_preparation_id", input.ValidPreparationID).
		Set("notes", input.Notes).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
			recipeStepPreparationsTableOwnershipColumn: input.BelongsToRecipeStep,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepPreparation updates a particular recipe step preparation. Note that UpdateRecipeStepPreparation expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeStepPreparation(ctx context.Context, input *models.RecipeStepPreparation) error {
	query, args := p.buildUpdateRecipeStepPreparationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveRecipeStepPreparationQuery returns a SQL query which marks a given recipe step preparation belonging to a given recipe step as archived.
func (p *Postgres) buildArchiveRecipeStepPreparationQuery(recipeStepID, recipeStepPreparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepPreparationsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeStepPreparationID,
			"archived_on": nil,
			recipeStepPreparationsTableOwnershipColumn: recipeStepID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepPreparation marks a recipe step preparation as archived in the database.
func (p *Postgres) ArchiveRecipeStepPreparation(ctx context.Context, recipeStepID, recipeStepPreparationID uint64) error {
	query, args := p.buildArchiveRecipeStepPreparationQuery(recipeStepID, recipeStepPreparationID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
