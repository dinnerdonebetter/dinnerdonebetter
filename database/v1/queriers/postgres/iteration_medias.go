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
	iterationMediasTableName                    = "iteration_medias"
	iterationMediasTablePathColumn              = "path"
	iterationMediasTableMimetypeColumn          = "mimetype"
	iterationMediasTableRecipeIterationIDColumn = "recipe_iteration_id"
	iterationMediasTableRecipeStepIDColumn      = "recipe_step_id"
	iterationMediasTableOwnershipColumn         = "belongs_to_recipe_iteration"
)

var (
	iterationMediasTableColumns = []string{
		fmt.Sprintf("%s.%s", iterationMediasTableName, idColumn),
		fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTablePathColumn),
		fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableMimetypeColumn),
		fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableRecipeIterationIDColumn),
		fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableRecipeStepIDColumn),
		fmt.Sprintf("%s.%s", iterationMediasTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", iterationMediasTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", iterationMediasTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableOwnershipColumn),
	}
)

// scanIterationMedia takes a database Scanner (i.e. *sql.Row) and scans the result into an Iteration Media struct
func (p *Postgres) scanIterationMedia(scan database.Scanner) (*models.IterationMedia, error) {
	x := &models.IterationMedia{}

	targetVars := []interface{}{
		&x.ID,
		&x.Path,
		&x.Mimetype,
		&x.RecipeIterationID,
		&x.RecipeStepID,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipeIteration,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanIterationMedias takes a logger and some database rows and turns them into a slice of iteration medias.
func (p *Postgres) scanIterationMedias(rows database.ResultIterator) ([]models.IterationMedia, error) {
	var (
		list []models.IterationMedia
	)

	for rows.Next() {
		x, err := p.scanIterationMedia(rows)
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

// buildIterationMediaExistsQuery constructs a SQL query for checking if an iteration media with a given ID belong to a a recipe iteration with a given ID exists
func (p *Postgres) buildIterationMediaExistsQuery(recipeID, recipeIterationID, iterationMediaID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", iterationMediasTableName, idColumn)).
		Prefix(existencePrefix).
		From(iterationMediasTableName).
		Join(recipeIterationsOnIterationMediasJoinClause).
		Join(recipesOnRecipeIterationsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", iterationMediasTableName, idColumn):                              iterationMediaID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn):                             recipeIterationID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
			fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableOwnershipColumn):   recipeIterationID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// IterationMediaExists queries the database to see if a given iteration media belonging to a given user exists.
func (p *Postgres) IterationMediaExists(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (exists bool, err error) {
	query, args := p.buildIterationMediaExistsQuery(recipeID, recipeIterationID, iterationMediaID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetIterationMediaQuery constructs a SQL query for fetching an iteration media with a given ID belong to a recipe iteration with a given ID.
func (p *Postgres) buildGetIterationMediaQuery(recipeID, recipeIterationID, iterationMediaID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(iterationMediasTableColumns...).
		From(iterationMediasTableName).
		Join(recipeIterationsOnIterationMediasJoinClause).
		Join(recipesOnRecipeIterationsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", iterationMediasTableName, idColumn):                              iterationMediaID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn):                             recipeIterationID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
			fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableOwnershipColumn):   recipeIterationID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetIterationMedia fetches an iteration media from the database.
func (p *Postgres) GetIterationMedia(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (*models.IterationMedia, error) {
	query, args := p.buildGetIterationMediaQuery(recipeID, recipeIterationID, iterationMediaID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanIterationMedia(row)
}

var (
	allIterationMediasCountQueryBuilder sync.Once
	allIterationMediasCountQuery        string
)

// buildGetAllIterationMediasCountQuery returns a query that fetches the total number of iteration medias in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllIterationMediasCountQuery() string {
	allIterationMediasCountQueryBuilder.Do(func() {
		var err error

		allIterationMediasCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, iterationMediasTableName)).
			From(iterationMediasTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", iterationMediasTableName, archivedOnColumn): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allIterationMediasCountQuery
}

// GetAllIterationMediasCount will fetch the count of iteration medias from the database.
func (p *Postgres) GetAllIterationMediasCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllIterationMediasCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfIterationMediasQuery returns a query that fetches every iteration media in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfIterationMediasQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(iterationMediasTableColumns...).
		From(iterationMediasTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", iterationMediasTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", iterationMediasTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllIterationMedias fetches every iteration media from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllIterationMedias(ctx context.Context, resultChannel chan []models.IterationMedia) error {
	count, err := p.GetAllIterationMediasCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfIterationMediasQuery(begin, end)
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

			iterationMedias, err := p.scanIterationMedias(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- iterationMedias
		}(beginID, endID)
	}

	return nil
}

// buildGetIterationMediasQuery builds a SQL query selecting iteration medias that adhere to a given QueryFilter and belong to a given recipe iteration,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetIterationMediasQuery(recipeID, recipeIterationID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(iterationMediasTableColumns...).
		From(iterationMediasTableName).
		Join(recipeIterationsOnIterationMediasJoinClause).
		Join(recipesOnRecipeIterationsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", iterationMediasTableName, archivedOnColumn):                      nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn):                             recipeIterationID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
			fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableOwnershipColumn):   recipeIterationID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", iterationMediasTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, iterationMediasTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetIterationMedias fetches a list of iteration medias from the database that meet a particular filter.
func (p *Postgres) GetIterationMedias(ctx context.Context, recipeID, recipeIterationID uint64, filter *models.QueryFilter) (*models.IterationMediaList, error) {
	query, args := p.buildGetIterationMediasQuery(recipeID, recipeIterationID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for iteration medias")
	}

	iterationMedias, err := p.scanIterationMedias(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.IterationMediaList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		IterationMedias: iterationMedias,
	}

	return list, nil
}

// buildGetIterationMediasWithIDsQuery builds a SQL query selecting iterationMedias that belong to a given recipe iteration,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetIterationMediasWithIDsQuery(recipeID, recipeIterationID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(iterationMediasTableColumns...).
		From(iterationMediasTableName).
		Join(recipeIterationsOnIterationMediasJoinClause).
		Join(recipesOnRecipeIterationsJoinClause).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", iterationMediasTableName, archivedOnColumn):                      nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, idColumn):                             recipeIterationID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
			fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableOwnershipColumn):   recipeIterationID,
		}).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(iterationMediasTableColumns...).
		FromSelect(subqueryBuilder, iterationMediasTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", iterationMediasTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetIterationMediasWithIDs fetches a list of iteration medias from the database that exist within a given set of IDs.
func (p *Postgres) GetIterationMediasWithIDs(ctx context.Context, recipeID, recipeIterationID uint64, limit uint8, ids []uint64) ([]models.IterationMedia, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetIterationMediasWithIDsQuery(recipeID, recipeIterationID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for iteration medias")
	}

	iterationMedias, err := p.scanIterationMedias(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return iterationMedias, nil
}

// buildCreateIterationMediaQuery takes an iteration media and returns a creation query for that iteration media and the relevant arguments.
func (p *Postgres) buildCreateIterationMediaQuery(input *models.IterationMedia) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(iterationMediasTableName).
		Columns(
			iterationMediasTablePathColumn,
			iterationMediasTableMimetypeColumn,
			iterationMediasTableRecipeIterationIDColumn,
			iterationMediasTableRecipeStepIDColumn,
			iterationMediasTableOwnershipColumn,
		).
		Values(
			input.Path,
			input.Mimetype,
			input.RecipeIterationID,
			input.RecipeStepID,
			input.BelongsToRecipeIteration,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateIterationMedia creates an iteration media in the database.
func (p *Postgres) CreateIterationMedia(ctx context.Context, input *models.IterationMediaCreationInput) (*models.IterationMedia, error) {
	x := &models.IterationMedia{
		Path:                     input.Path,
		Mimetype:                 input.Mimetype,
		RecipeIterationID:        input.RecipeIterationID,
		RecipeStepID:             input.RecipeStepID,
		BelongsToRecipeIteration: input.BelongsToRecipeIteration,
	}

	query, args := p.buildCreateIterationMediaQuery(x)

	// create the iteration media.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing iteration media creation query: %w", err)
	}

	return x, nil
}

// buildUpdateIterationMediaQuery takes an iteration media and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateIterationMediaQuery(input *models.IterationMedia) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(iterationMediasTableName).
		Set(iterationMediasTablePathColumn, input.Path).
		Set(iterationMediasTableMimetypeColumn, input.Mimetype).
		Set(iterationMediasTableRecipeIterationIDColumn, input.RecipeIterationID).
		Set(iterationMediasTableRecipeStepIDColumn, input.RecipeStepID).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                            input.ID,
			iterationMediasTableOwnershipColumn: input.BelongsToRecipeIteration,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateIterationMedia updates a particular iteration media. Note that UpdateIterationMedia expects the provided input to have a valid ID.
func (p *Postgres) UpdateIterationMedia(ctx context.Context, input *models.IterationMedia) error {
	query, args := p.buildUpdateIterationMediaQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveIterationMediaQuery returns a SQL query which marks a given iteration media belonging to a given recipe iteration as archived.
func (p *Postgres) buildArchiveIterationMediaQuery(recipeIterationID, iterationMediaID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(iterationMediasTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                            iterationMediaID,
			archivedOnColumn:                    nil,
			iterationMediasTableOwnershipColumn: recipeIterationID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveIterationMedia marks an iteration media as archived in the database.
func (p *Postgres) ArchiveIterationMedia(ctx context.Context, recipeIterationID, iterationMediaID uint64) error {
	query, args := p.buildArchiveIterationMediaQuery(recipeIterationID, iterationMediaID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
