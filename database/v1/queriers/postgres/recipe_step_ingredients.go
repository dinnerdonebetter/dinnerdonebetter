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
	recipeStepIngredientsTableName            = "recipe_step_ingredients"
	recipeStepIngredientsTableOwnershipColumn = "belongs_to_recipe_step"
)

var (
	recipeStepIngredientsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "id"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "valid_ingredient_id"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "ingredient_notes"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "quantity_type"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "quantity_value"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "quantity_notes"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "product_of_recipe_step_id"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "created_on"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableOwnershipColumn),
	}
)

// scanRecipeStepIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Ingredient struct
func (p *Postgres) scanRecipeStepIngredient(scan database.Scanner, includeCount bool) (*models.RecipeStepIngredient, uint64, error) {
	x := &models.RecipeStepIngredient{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.ValidIngredientID,
		&x.IngredientNotes,
		&x.QuantityType,
		&x.QuantityValue,
		&x.QuantityNotes,
		&x.ProductOfRecipeStepID,
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

// scanRecipeStepIngredients takes a logger and some database rows and turns them into a slice of recipe step ingredients.
func (p *Postgres) scanRecipeStepIngredients(rows database.ResultIterator) ([]models.RecipeStepIngredient, uint64, error) {
	var (
		list  []models.RecipeStepIngredient
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanRecipeStepIngredient(rows, true)
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

// buildRecipeStepIngredientExistsQuery constructs a SQL query for checking if a recipe step ingredient with a given ID belong to a a recipe step with a given ID exists
func (p *Postgres) buildRecipeStepIngredientExistsQuery(recipeID, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", recipeStepIngredientsTableName)).
		Prefix(existencePrefix).
		From(recipeStepIngredientsTableName).
		Join(recipeStepsOnRecipeStepIngredientsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeStepIngredientsTableName):                                            recipeStepIngredientID,
			fmt.Sprintf("%s.id", recipesTableName):                                                          recipeID,
			fmt.Sprintf("%s.id", recipeStepsTableName):                                                      recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                     recipeID,
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableOwnershipColumn): recipeStepID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeStepIngredientExists queries the database to see if a given recipe step ingredient belonging to a given user exists.
func (p *Postgres) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (exists bool, err error) {
	query, args := p.buildRecipeStepIngredientExistsQuery(recipeID, recipeStepID, recipeStepIngredientID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRecipeStepIngredientQuery constructs a SQL query for fetching a recipe step ingredient with a given ID belong to a recipe step with a given ID.
func (p *Postgres) buildGetRecipeStepIngredientQuery(recipeID, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipeStepIngredientsTableColumns...).
		From(recipeStepIngredientsTableName).
		Join(recipeStepsOnRecipeStepIngredientsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", recipeStepIngredientsTableName):                                            recipeStepIngredientID,
			fmt.Sprintf("%s.id", recipesTableName):                                                          recipeID,
			fmt.Sprintf("%s.id", recipeStepsTableName):                                                      recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                     recipeID,
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableOwnershipColumn): recipeStepID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepIngredient fetches a recipe step ingredient from the database.
func (p *Postgres) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*models.RecipeStepIngredient, error) {
	query, args := p.buildGetRecipeStepIngredientQuery(recipeID, recipeStepID, recipeStepIngredientID)
	row := p.db.QueryRowContext(ctx, query, args...)

	recipeStepIngredient, _, err := p.scanRecipeStepIngredient(row, false)
	return recipeStepIngredient, err
}

var (
	allRecipeStepIngredientsCountQueryBuilder sync.Once
	allRecipeStepIngredientsCountQuery        string
)

// buildGetAllRecipeStepIngredientsCountQuery returns a query that fetches the total number of recipe step ingredients in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipeStepIngredientsCountQuery() string {
	allRecipeStepIngredientsCountQueryBuilder.Do(func() {
		var err error

		allRecipeStepIngredientsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipeStepIngredientsTableName)).
			From(recipeStepIngredientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", recipeStepIngredientsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeStepIngredientsCountQuery
}

// GetAllRecipeStepIngredientsCount will fetch the count of recipe step ingredients from the database.
func (p *Postgres) GetAllRecipeStepIngredientsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeStepIngredientsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeStepIngredientsQuery builds a SQL query selecting recipe step ingredients that adhere to a given QueryFilter and belong to a given recipe step,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeStepIngredientsQuery(recipeID, recipeStepID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(recipeStepIngredientsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllRecipeStepIngredientsCountQuery()))...).
		From(recipeStepIngredientsTableName).
		Join(recipeStepsOnRecipeStepIngredientsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", recipeStepIngredientsTableName):                                   nil,
			fmt.Sprintf("%s.id", recipesTableName):                                                          recipeID,
			fmt.Sprintf("%s.id", recipeStepsTableName):                                                      recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                     recipeID,
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableOwnershipColumn): recipeStepID,
		}).
		OrderBy(fmt.Sprintf("%s.id", recipeStepIngredientsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipeStepIngredientsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepIngredients fetches a list of recipe step ingredients from the database that meet a particular filter.
func (p *Postgres) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepIngredientList, error) {
	query, args := p.buildGetRecipeStepIngredientsQuery(recipeID, recipeStepID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step ingredients")
	}

	recipeStepIngredients, count, err := p.scanRecipeStepIngredients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeStepIngredientList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeStepIngredients: recipeStepIngredients,
	}

	return list, nil
}

// buildCreateRecipeStepIngredientQuery takes a recipe step ingredient and returns a creation query for that recipe step ingredient and the relevant arguments.
func (p *Postgres) buildCreateRecipeStepIngredientQuery(input *models.RecipeStepIngredient) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeStepIngredientsTableName).
		Columns(
			"valid_ingredient_id",
			"ingredient_notes",
			"quantity_type",
			"quantity_value",
			"quantity_notes",
			"product_of_recipe_step_id",
			recipeStepIngredientsTableOwnershipColumn,
		).
		Values(
			input.ValidIngredientID,
			input.IngredientNotes,
			input.QuantityType,
			input.QuantityValue,
			input.QuantityNotes,
			input.ProductOfRecipeStepID,
			input.BelongsToRecipeStep,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database.
func (p *Postgres) CreateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredientCreationInput) (*models.RecipeStepIngredient, error) {
	x := &models.RecipeStepIngredient{
		ValidIngredientID:     input.ValidIngredientID,
		IngredientNotes:       input.IngredientNotes,
		QuantityType:          input.QuantityType,
		QuantityValue:         input.QuantityValue,
		QuantityNotes:         input.QuantityNotes,
		ProductOfRecipeStepID: input.ProductOfRecipeStepID,
		BelongsToRecipeStep:   input.BelongsToRecipeStep,
	}

	query, args := p.buildCreateRecipeStepIngredientQuery(x)

	// create the recipe step ingredient.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step ingredient creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeStepIngredientQuery takes a recipe step ingredient and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeStepIngredientQuery(input *models.RecipeStepIngredient) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepIngredientsTableName).
		Set("valid_ingredient_id", input.ValidIngredientID).
		Set("ingredient_notes", input.IngredientNotes).
		Set("quantity_type", input.QuantityType).
		Set("quantity_value", input.QuantityValue).
		Set("quantity_notes", input.QuantityNotes).
		Set("product_of_recipe_step_id", input.ProductOfRecipeStepID).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
			recipeStepIngredientsTableOwnershipColumn: input.BelongsToRecipeStep,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepIngredient updates a particular recipe step ingredient. Note that UpdateRecipeStepIngredient expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredient) error {
	query, args := p.buildUpdateRecipeStepIngredientQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveRecipeStepIngredientQuery returns a SQL query which marks a given recipe step ingredient belonging to a given recipe step as archived.
func (p *Postgres) buildArchiveRecipeStepIngredientQuery(recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepIngredientsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeStepIngredientID,
			"archived_on": nil,
			recipeStepIngredientsTableOwnershipColumn: recipeStepID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepIngredient marks a recipe step ingredient as archived in the database.
func (p *Postgres) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID uint64) error {
	query, args := p.buildArchiveRecipeStepIngredientQuery(recipeStepID, recipeStepIngredientID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
