package sqlite

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
func (s *Sqlite) buildGetRecipeStepIngredientQuery(recipeStepIngredientID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Select(recipeStepIngredientsTableColumns...).
		From(recipeStepIngredientsTableName).
		Where(squirrel.Eq{
			"id":         recipeStepIngredientID,
			"belongs_to": userID,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepIngredient fetches a recipe step ingredient from the sqlite database
func (s *Sqlite) GetRecipeStepIngredient(ctx context.Context, recipeStepIngredientID, userID uint64) (*models.RecipeStepIngredient, error) {
	query, args := s.buildGetRecipeStepIngredientQuery(recipeStepIngredientID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)
	return scanRecipeStepIngredient(row)
}

// buildGetRecipeStepIngredientCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of recipe step ingredients belonging to a given user that meet a given query
func (s *Sqlite) buildGetRecipeStepIngredientCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := s.sqlBuilder.
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
	s.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepIngredientCount will fetch the count of recipe step ingredients from the database that meet a particular filter and belong to a particular user.
func (s *Sqlite) GetRecipeStepIngredientCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := s.buildGetRecipeStepIngredientCountQuery(filter, userID)
	err = s.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allRecipeStepIngredientsCountQueryBuilder sync.Once
	allRecipeStepIngredientsCountQuery        string
)

// buildGetAllRecipeStepIngredientsCountQuery returns a query that fetches the total number of recipe step ingredients in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (s *Sqlite) buildGetAllRecipeStepIngredientsCountQuery() string {
	allRecipeStepIngredientsCountQueryBuilder.Do(func() {
		var err error
		allRecipeStepIngredientsCountQuery, _, err = s.sqlBuilder.
			Select(CountQuery).
			From(recipeStepIngredientsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		s.logQueryBuildingError(err)
	})

	return allRecipeStepIngredientsCountQuery
}

// GetAllRecipeStepIngredientsCount will fetch the count of recipe step ingredients from the database
func (s *Sqlite) GetAllRecipeStepIngredientsCount(ctx context.Context) (count uint64, err error) {
	err = s.db.QueryRowContext(ctx, s.buildGetAllRecipeStepIngredientsCountQuery()).Scan(&count)
	return count, err
}

// buildGetRecipeStepIngredientsQuery builds a SQL query selecting recipe step ingredients that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (s *Sqlite) buildGetRecipeStepIngredientsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := s.sqlBuilder.
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
	s.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepIngredients fetches a list of recipe step ingredients from the database that meet a particular filter
func (s *Sqlite) GetRecipeStepIngredients(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeStepIngredientList, error) {
	query, args := s.buildGetRecipeStepIngredientsQuery(filter, userID)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step ingredients")
	}

	list, err := scanRecipeStepIngredients(s.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := s.GetRecipeStepIngredientCount(ctx, filter, userID)
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
func (s *Sqlite) GetAllRecipeStepIngredientsForUser(ctx context.Context, userID uint64) ([]models.RecipeStepIngredient, error) {
	query, args := s.buildGetRecipeStepIngredientsQuery(nil, userID)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching recipe step ingredients for user")
	}

	list, err := scanRecipeStepIngredients(s.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateRecipeStepIngredientQuery takes a recipe step ingredient and returns a creation query for that recipe step ingredient and the relevant arguments.
func (s *Sqlite) buildCreateRecipeStepIngredientQuery(input *models.RecipeStepIngredient) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
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
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// buildRecipeStepIngredientCreationTimeQuery takes a recipe step ingredient and returns a creation query for that recipe step ingredient and the relevant arguments
func (s *Sqlite) buildRecipeStepIngredientCreationTimeQuery(recipeStepIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select("created_on").
		From(recipeStepIngredientsTableName).
		Where(squirrel.Eq{"id": recipeStepIngredientID}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database
func (s *Sqlite) CreateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredientCreationInput) (*models.RecipeStepIngredient, error) {
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

	query, args := s.buildCreateRecipeStepIngredientQuery(x)

	// create the recipe step ingredient
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step ingredient creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := s.buildRecipeStepIngredientCreationTimeQuery(x.ID)
		s.logCreationTimeRetrievalError(s.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateRecipeStepIngredientQuery takes a recipe step ingredient and returns an update SQL query, with the relevant query parameters
func (s *Sqlite) buildUpdateRecipeStepIngredientQuery(input *models.RecipeStepIngredient) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
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
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepIngredient updates a particular recipe step ingredient. Note that UpdateRecipeStepIngredient expects the provided input to have a valid ID.
func (s *Sqlite) UpdateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredient) error {
	query, args := s.buildUpdateRecipeStepIngredientQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveRecipeStepIngredientQuery returns a SQL query which marks a given recipe step ingredient belonging to a given user as archived.
func (s *Sqlite) buildArchiveRecipeStepIngredientQuery(recipeStepIngredientID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = s.sqlBuilder.
		Update(recipeStepIngredientsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          recipeStepIngredientID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepIngredient marks a recipe step ingredient as archived in the database
func (s *Sqlite) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepIngredientID, userID uint64) error {
	query, args := s.buildArchiveRecipeStepIngredientQuery(recipeStepIngredientID, userID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
