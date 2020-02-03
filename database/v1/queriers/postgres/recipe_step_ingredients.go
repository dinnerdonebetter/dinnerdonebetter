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
	recipeStepIngredientsTableName = "recipe_step_ingredients"
)

var (
	recipeStepIngredientsTableColumns = []string{
		"id",
		"ingredient_id",
		"quantity_type",
		"quantity_value",
		"quantity_notes",
		"product_of_recipe",
		"ingredient_notes",
		"recipe_step_id",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanRecipeStepIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Ingredient struct
func scanRecipeStepIngredient(scan database.Scanner) (*models.RecipeStepIngredient, error) {
	x := &models.RecipeStepIngredient{}

	if err := scan.Scan(
		&x.ID,
		&x.IngredientID,
		&x.QuantityType,
		&x.QuantityValue,
		&x.QuantityNotes,
		&x.ProductOfRecipe,
		&x.IngredientNotes,
		&x.RecipeStepID,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeStepIngredients takes a logger and some database rows and turns them into a slice of recipe step ingredients
func scanRecipeStepIngredients(logger logging.Logger, rows *sql.Rows) ([]models.RecipeStepIngredient, error) {
	var list []models.RecipeStepIngredient

	for rows.Next() {
		x, err := scanRecipeStepIngredient(rows)
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

// buildGetRecipeStepIngredientQuery constructs a SQL query for fetching a recipe step ingredient with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetRecipeStepIngredientQuery(recipeStepIngredientID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(recipeStepIngredientsTableColumns...).
		From(recipeStepIngredientsTableName).
		Where(squirrel.Eq{
			"id":         recipeStepIngredientID,
			"belongs_to": userID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepIngredient fetches a recipe step ingredient from the postgres database
func (p *Postgres) GetRecipeStepIngredient(ctx context.Context, recipeStepIngredientID, userID uint64) (*models.RecipeStepIngredient, error) {
	query, args := p.buildGetRecipeStepIngredientQuery(recipeStepIngredientID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanRecipeStepIngredient(row)
}

// buildGetRecipeStepIngredientCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of recipe step ingredients belonging to a given user that meet a given query
func (p *Postgres) buildGetRecipeStepIngredientCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(recipeStepIngredientsTableName).
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

// GetRecipeStepIngredientCount will fetch the count of recipe step ingredients from the database that meet a particular filter and belong to a particular user.
func (p *Postgres) GetRecipeStepIngredientCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := p.buildGetRecipeStepIngredientCountQuery(filter, userID)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
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
			Select(CountQuery).
			From(recipeStepIngredientsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeStepIngredientsCountQuery
}

// GetAllRecipeStepIngredientsCount will fetch the count of recipe step ingredients from the database
func (p *Postgres) GetAllRecipeStepIngredientsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeStepIngredientsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeStepIngredientsQuery builds a SQL query selecting recipe step ingredients that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeStepIngredientsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(recipeStepIngredientsTableColumns...).
		From(recipeStepIngredientsTableName).
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

// GetRecipeStepIngredients fetches a list of recipe step ingredients from the database that meet a particular filter
func (p *Postgres) GetRecipeStepIngredients(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepIngredientList, error) {
	query, args := p.buildGetRecipeStepIngredientsQuery(filter, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step ingredients")
	}

	list, err := scanRecipeStepIngredients(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetRecipeStepIngredientCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching recipe step ingredient count: %w", err)
	}

	x := &models.RecipeStepIngredientList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		RecipeStepIngredients: list,
	}

	return x, nil
}

// GetAllRecipeStepIngredientsForUser fetches every recipe step ingredient belonging to a user
func (p *Postgres) GetAllRecipeStepIngredientsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepIngredient, error) {
	query, args := p.buildGetRecipeStepIngredientsQuery(nil, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching recipe step ingredients for user")
	}

	list, err := scanRecipeStepIngredients(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRecipeStepIngredientQuery takes a recipe step ingredient and returns a creation query for that recipe step ingredient and the relevant arguments.
func (p *Postgres) buildCreateRecipeStepIngredientQuery(input *models.RecipeStepIngredient) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(recipeStepIngredientsTableName).
		Columns(
			"ingredient_id",
			"quantity_type",
			"quantity_value",
			"quantity_notes",
			"product_of_recipe",
			"ingredient_notes",
			"recipe_step_id",
			"belongs_to",
		).
		Values(
			input.IngredientID,
			input.QuantityType,
			input.QuantityValue,
			input.QuantityNotes,
			input.ProductOfRecipe,
			input.IngredientNotes,
			input.RecipeStepID,
			input.BelongsTo,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database
func (p *Postgres) CreateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredientCreationInput) (*models.RecipeStepIngredient, error) {
	x := &models.RecipeStepIngredient{
		IngredientID:    input.IngredientID,
		QuantityType:    input.QuantityType,
		QuantityValue:   input.QuantityValue,
		QuantityNotes:   input.QuantityNotes,
		ProductOfRecipe: input.ProductOfRecipe,
		IngredientNotes: input.IngredientNotes,
		RecipeStepID:    input.RecipeStepID,
		BelongsTo:       input.BelongsTo,
	}

	query, args := p.buildCreateRecipeStepIngredientQuery(x)

	// create the recipe step ingredient
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step ingredient creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeStepIngredientQuery takes a recipe step ingredient and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdateRecipeStepIngredientQuery(input *models.RecipeStepIngredient) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(recipeStepIngredientsTableName).
		Set("ingredient_id", input.IngredientID).
		Set("quantity_type", input.QuantityType).
		Set("quantity_value", input.QuantityValue).
		Set("quantity_notes", input.QuantityNotes).
		Set("product_of_recipe", input.ProductOfRecipe).
		Set("ingredient_notes", input.IngredientNotes).
		Set("recipe_step_id", input.RecipeStepID).
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

// UpdateRecipeStepIngredient updates a particular recipe step ingredient. Note that UpdateRecipeStepIngredient expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredient) error {
	query, args := p.buildUpdateRecipeStepIngredientQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveRecipeStepIngredientQuery returns a SQL query which marks a given recipe step ingredient belonging to a given user as archived.
func (p *Postgres) buildArchiveRecipeStepIngredientQuery(recipeStepIngredientID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(recipeStepIngredientsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeStepIngredientID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepIngredient marks a recipe step ingredient as archived in the database
func (p *Postgres) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepIngredientID, userID uint64) error {
	query, args := p.buildArchiveRecipeStepIngredientQuery(recipeStepIngredientID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
