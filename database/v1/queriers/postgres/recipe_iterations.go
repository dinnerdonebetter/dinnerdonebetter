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
	recipeIterationsTableName                      = "recipe_iterations"
	recipeIterationsTableRecipeIDColumn            = "recipe_id"
	recipeIterationsTableEndDifficultyRatingColumn = "end_difficulty_rating"
	recipeIterationsTableEndComplexityRatingColumn = "end_complexity_rating"
	recipeIterationsTableEndTasteRatingColumn      = "end_taste_rating"
	recipeIterationsTableEndOverallRatingColumn    = "end_overall_rating"
	recipeIterationsTableOwnershipColumn           = "belongs_to_recipe"
)

var (
	recipeIterationsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableRecipeIDColumn),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableEndDifficultyRatingColumn),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableEndComplexityRatingColumn),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableEndTasteRatingColumn),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableEndOverallRatingColumn),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn),
	}

	recipeIterationsOnIterationMediasJoinClause = fmt.Sprintf("%s ON %s.%s=%s.%s", recipeIterationsTableName, iterationMediasTableName, iterationMediasTableOwnershipColumn, recipeIterationsTableName, idColumn)
)

// scanRecipeIteration takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Iteration struct
func (p *Postgres) scanRecipeIteration(scan database.Scanner) (*models.RecipeIteration, error) {
	x := &models.RecipeIteration{}

	targetVars := []interface{}{
		&x.ID,
		&x.RecipeID,
		&x.EndDifficultyRating,
		&x.EndComplexityRating,
		&x.EndTasteRating,
		&x.EndOverallRating,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipe,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanRecipeIterations takes a logger and some database rows and turns them into a slice of recipe iterations.
func (p *Postgres) scanRecipeIterations(rows database.ResultIterator) ([]models.RecipeIteration, error) {
	var (
		list []models.RecipeIteration
	)

	for rows.Next() {
		x, err := p.scanRecipeIteration(rows)
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

// buildRecipeIterationExistsQuery constructs a SQL query for checking if a recipe iteration with a given ID belong to a a recipe with a given ID exists
func (p *Postgres) buildRecipeIterationExistsQuery(recipeID, recipeIterationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn)).
		Prefix(existencePrefix).
		From(recipeIterationsTableName).
		Join(recipesOnRecipeIterationsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn):                             recipeIterationID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
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
			fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn):                             recipeIterationID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
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
	return p.scanRecipeIteration(row)
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
				fmt.Sprintf("%s.%s", recipeIterationsTableName, archivedOnColumn): nil,
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

// buildGetBatchOfRecipeIterationsQuery returns a query that fetches every recipe iteration in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfRecipeIterationsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(recipeIterationsTableColumns...).
		From(recipeIterationsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllRecipeIterations fetches every recipe iteration from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllRecipeIterations(ctx context.Context, resultChannel chan []models.RecipeIteration) error {
	count, err := p.GetAllRecipeIterationsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfRecipeIterationsQuery(begin, end)
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

			recipeIterations, err := p.scanRecipeIterations(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- recipeIterations
		}(beginID, endID)
	}

	return nil
}

// buildGetRecipeIterationsQuery builds a SQL query selecting recipe iterations that adhere to a given QueryFilter and belong to a given recipe,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeIterationsQuery(recipeID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(recipeIterationsTableColumns...).
		From(recipeIterationsTableName).
		Join(recipesOnRecipeIterationsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeIterationsTableName, archivedOnColumn):                     nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn))

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

	recipeIterations, err := p.scanRecipeIterations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeIterationList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		RecipeIterations: recipeIterations,
	}

	return list, nil
}

// buildGetRecipeIterationsWithIDsQuery builds a SQL query selecting recipeIterations that belong to a given recipe,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetRecipeIterationsWithIDsQuery(recipeID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(recipeIterationsTableColumns...).
		From(recipeIterationsTableName).
		Join(recipesOnRecipeIterationsJoinClause).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeIterationsTableName, archivedOnColumn):                     nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
		}).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(recipeIterationsTableColumns...).
		FromSelect(subqueryBuilder, recipeIterationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeIterationsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeIterationsWithIDs fetches a list of recipe iterations from the database that exist within a given set of IDs.
func (p *Postgres) GetRecipeIterationsWithIDs(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) ([]models.RecipeIteration, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetRecipeIterationsWithIDsQuery(recipeID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe iterations")
	}

	recipeIterations, err := p.scanRecipeIterations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return recipeIterations, nil
}

// buildCreateRecipeIterationQuery takes a recipe iteration and returns a creation query for that recipe iteration and the relevant arguments.
func (p *Postgres) buildCreateRecipeIterationQuery(input *models.RecipeIteration) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeIterationsTableName).
		Columns(
			recipeIterationsTableRecipeIDColumn,
			recipeIterationsTableEndDifficultyRatingColumn,
			recipeIterationsTableEndComplexityRatingColumn,
			recipeIterationsTableEndTasteRatingColumn,
			recipeIterationsTableEndOverallRatingColumn,
			recipeIterationsTableOwnershipColumn,
		).
		Values(
			input.RecipeID,
			input.EndDifficultyRating,
			input.EndComplexityRating,
			input.EndTasteRating,
			input.EndOverallRating,
			input.BelongsToRecipe,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeIteration creates a recipe iteration in the database.
func (p *Postgres) CreateRecipeIteration(ctx context.Context, input *models.RecipeIterationCreationInput) (*models.RecipeIteration, error) {
	x := &models.RecipeIteration{
		RecipeID:            input.RecipeID,
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
		Set(recipeIterationsTableRecipeIDColumn, input.RecipeID).
		Set(recipeIterationsTableEndDifficultyRatingColumn, input.EndDifficultyRating).
		Set(recipeIterationsTableEndComplexityRatingColumn, input.EndComplexityRating).
		Set(recipeIterationsTableEndTasteRatingColumn, input.EndTasteRating).
		Set(recipeIterationsTableEndOverallRatingColumn, input.EndOverallRating).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             input.ID,
			recipeIterationsTableOwnershipColumn: input.BelongsToRecipe,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeIteration updates a particular recipe iteration. Note that UpdateRecipeIteration expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeIteration(ctx context.Context, input *models.RecipeIteration) error {
	query, args := p.buildUpdateRecipeIterationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveRecipeIterationQuery returns a SQL query which marks a given recipe iteration belonging to a given recipe as archived.
func (p *Postgres) buildArchiveRecipeIterationQuery(recipeID, recipeIterationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeIterationsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             recipeIterationID,
			archivedOnColumn:                     nil,
			recipeIterationsTableOwnershipColumn: recipeID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
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
