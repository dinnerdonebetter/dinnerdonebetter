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
	recipeStepIngredientsTableName                  = "recipe_step_ingredients"
	recipeStepIngredientsTableIngredientIDColumn    = "ingredient_id"
	recipeStepIngredientsTableQuantityTypeColumn    = "quantity_type"
	recipeStepIngredientsTableQuantityValueColumn   = "quantity_value"
	recipeStepIngredientsTableQuantityNotesColumn   = "quantity_notes"
	recipeStepIngredientsTableProductOfRecipeColumn = "product_of_recipe"
	recipeStepIngredientsTableIngredientNotesColumn = "ingredient_notes"
	recipeStepIngredientsTableOwnershipColumn       = "belongs_to_recipe_step"
)

var (
	recipeStepIngredientsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, idColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableIngredientIDColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableQuantityTypeColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableQuantityValueColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableQuantityNotesColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableProductOfRecipeColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableIngredientNotesColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableOwnershipColumn),
	}
)

// scanRecipeStepIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Ingredient struct
func (p *Postgres) scanRecipeStepIngredient(scan database.Scanner) (*models.RecipeStepIngredient, error) {
	x := &models.RecipeStepIngredient{}

	targetVars := []interface{}{
		&x.ID,
		&x.IngredientID,
		&x.QuantityType,
		&x.QuantityValue,
		&x.QuantityNotes,
		&x.ProductOfRecipe,
		&x.IngredientNotes,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipeStep,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeStepIngredients takes a logger and some database rows and turns them into a slice of recipe step ingredients.
func (p *Postgres) scanRecipeStepIngredients(rows database.ResultIterator) ([]models.RecipeStepIngredient, error) {
	var (
		list []models.RecipeStepIngredient
	)

	for rows.Next() {
		x, err := p.scanRecipeStepIngredient(rows)
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

// buildRecipeStepIngredientExistsQuery constructs a SQL query for checking if a recipe step ingredient with a given ID belong to a a recipe step with a given ID exists
func (p *Postgres) buildRecipeStepIngredientExistsQuery(recipeID, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, idColumn)).
		Prefix(existencePrefix).
		From(recipeStepIngredientsTableName).
		Join(recipeStepsOnRecipeStepIngredientsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, idColumn):                                  recipeStepIngredientID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                            recipeStepID,
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
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, idColumn):                                  recipeStepIngredientID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                            recipeStepID,
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
	return p.scanRecipeStepIngredient(row)
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
				fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, archivedOnColumn): nil,
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

// buildGetBatchOfRecipeStepIngredientsQuery returns a query that fetches every recipe step ingredient in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfRecipeStepIngredientsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(recipeStepIngredientsTableColumns...).
		From(recipeStepIngredientsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllRecipeStepIngredients fetches every recipe step ingredient from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllRecipeStepIngredients(ctx context.Context, resultChannel chan []models.RecipeStepIngredient) error {
	count, err := p.GetAllRecipeStepIngredientsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfRecipeStepIngredientsQuery(begin, end)
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

			recipeStepIngredients, err := p.scanRecipeStepIngredients(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- recipeStepIngredients
		}(beginID, endID)
	}

	return nil
}

// buildGetRecipeStepIngredientsQuery builds a SQL query selecting recipe step ingredients that adhere to a given QueryFilter and belong to a given recipe step,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeStepIngredientsQuery(recipeID, recipeStepID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(recipeStepIngredientsTableColumns...).
		From(recipeStepIngredientsTableName).
		Join(recipeStepsOnRecipeStepIngredientsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, archivedOnColumn):                          nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                            recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                     recipeID,
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableOwnershipColumn): recipeStepID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, idColumn))

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

	recipeStepIngredients, err := p.scanRecipeStepIngredients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeStepIngredientList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		RecipeStepIngredients: recipeStepIngredients,
	}

	return list, nil
}

// buildGetRecipeStepIngredientsWithIDsQuery builds a SQL query selecting recipeStepIngredients that belong to a given recipe step,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetRecipeStepIngredientsWithIDsQuery(recipeID, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(recipeStepIngredientsTableColumns...).
		From(recipeStepIngredientsTableName).
		Join(recipeStepsOnRecipeStepIngredientsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, archivedOnColumn):                          nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                            recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                     recipeID,
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, recipeStepIngredientsTableOwnershipColumn): recipeStepID,
		}).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(recipeStepIngredientsTableColumns...).
		FromSelect(subqueryBuilder, recipeStepIngredientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepIngredientsWithIDs fetches a list of recipe step ingredients from the database that exist within a given set of IDs.
func (p *Postgres) GetRecipeStepIngredientsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepIngredient, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetRecipeStepIngredientsWithIDsQuery(recipeID, recipeStepID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step ingredients")
	}

	recipeStepIngredients, err := p.scanRecipeStepIngredients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return recipeStepIngredients, nil
}

// buildCreateRecipeStepIngredientQuery takes a recipe step ingredient and returns a creation query for that recipe step ingredient and the relevant arguments.
func (p *Postgres) buildCreateRecipeStepIngredientQuery(input *models.RecipeStepIngredient) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeStepIngredientsTableName).
		Columns(
			recipeStepIngredientsTableIngredientIDColumn,
			recipeStepIngredientsTableQuantityTypeColumn,
			recipeStepIngredientsTableQuantityValueColumn,
			recipeStepIngredientsTableQuantityNotesColumn,
			recipeStepIngredientsTableProductOfRecipeColumn,
			recipeStepIngredientsTableIngredientNotesColumn,
			recipeStepIngredientsTableOwnershipColumn,
		).
		Values(
			input.IngredientID,
			input.QuantityType,
			input.QuantityValue,
			input.QuantityNotes,
			input.ProductOfRecipe,
			input.IngredientNotes,
			input.BelongsToRecipeStep,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database.
func (p *Postgres) CreateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredientCreationInput) (*models.RecipeStepIngredient, error) {
	x := &models.RecipeStepIngredient{
		IngredientID:        input.IngredientID,
		QuantityType:        input.QuantityType,
		QuantityValue:       input.QuantityValue,
		QuantityNotes:       input.QuantityNotes,
		ProductOfRecipe:     input.ProductOfRecipe,
		IngredientNotes:     input.IngredientNotes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
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
		Set(recipeStepIngredientsTableIngredientIDColumn, input.IngredientID).
		Set(recipeStepIngredientsTableQuantityTypeColumn, input.QuantityType).
		Set(recipeStepIngredientsTableQuantityValueColumn, input.QuantityValue).
		Set(recipeStepIngredientsTableQuantityNotesColumn, input.QuantityNotes).
		Set(recipeStepIngredientsTableProductOfRecipeColumn, input.ProductOfRecipe).
		Set(recipeStepIngredientsTableIngredientNotesColumn, input.IngredientNotes).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
			recipeStepIngredientsTableOwnershipColumn: input.BelongsToRecipeStep,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepIngredient updates a particular recipe step ingredient. Note that UpdateRecipeStepIngredient expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeStepIngredient(ctx context.Context, input *models.RecipeStepIngredient) error {
	query, args := p.buildUpdateRecipeStepIngredientQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveRecipeStepIngredientQuery returns a SQL query which marks a given recipe step ingredient belonging to a given recipe step as archived.
func (p *Postgres) buildArchiveRecipeStepIngredientQuery(recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepIngredientsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         recipeStepIngredientID,
			archivedOnColumn: nil,
			recipeStepIngredientsTableOwnershipColumn: recipeStepID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
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
