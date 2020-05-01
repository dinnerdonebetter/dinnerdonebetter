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
	iterationMediasTableName            = "iteration_medias"
	iterationMediasTableOwnershipColumn = "belongs_to_recipe_iteration"
)

var (
	iterationMediasTableColumns = []string{
		fmt.Sprintf("%s.%s", iterationMediasTableName, "id"),
		fmt.Sprintf("%s.%s", iterationMediasTableName, "source"),
		fmt.Sprintf("%s.%s", iterationMediasTableName, "mimetype"),
		fmt.Sprintf("%s.%s", iterationMediasTableName, "created_on"),
		fmt.Sprintf("%s.%s", iterationMediasTableName, "updated_on"),
		fmt.Sprintf("%s.%s", iterationMediasTableName, "archived_on"),
		fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableOwnershipColumn),
	}
)

// scanIterationMedia takes a database Scanner (i.e. *sql.Row) and scans the result into an Iteration Media struct
func (p *Postgres) scanIterationMedia(scan database.Scanner, includeCount bool) (*models.IterationMedia, uint64, error) {
	x := &models.IterationMedia{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Source,
		&x.Mimetype,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipeIteration,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}

// scanIterationMedias takes a logger and some database rows and turns them into a slice of iteration medias.
func (p *Postgres) scanIterationMedias(rows database.ResultIterator) ([]models.IterationMedia, uint64, error) {
	var (
		list  []models.IterationMedia
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanIterationMedia(rows, true)
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

// buildIterationMediaExistsQuery constructs a SQL query for checking if an iteration media with a given ID belong to a a recipe iteration with a given ID exists
func (p *Postgres) buildIterationMediaExistsQuery(recipeID, recipeIterationID, iterationMediaID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", iterationMediasTableName)).
		Prefix(existencePrefix).
		From(iterationMediasTableName).
		Join(recipeIterationsOnIterationMediasJoinClause).
		Join(recipesOnRecipeIterationsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", iterationMediasTableName):                                        iterationMediaID,
			fmt.Sprintf("%s.id", recipesTableName):                                                recipeID,
			fmt.Sprintf("%s.id", recipeIterationsTableName):                                       recipeIterationID,
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
			fmt.Sprintf("%s.id", iterationMediasTableName):                                        iterationMediaID,
			fmt.Sprintf("%s.id", recipesTableName):                                                recipeID,
			fmt.Sprintf("%s.id", recipeIterationsTableName):                                       recipeIterationID,
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

	iterationMedia, _, err := p.scanIterationMedia(row, false)
	return iterationMedia, err
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
				fmt.Sprintf("%s.archived_on", iterationMediasTableName): nil,
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

// buildGetIterationMediasQuery builds a SQL query selecting iteration medias that adhere to a given QueryFilter and belong to a given recipe iteration,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetIterationMediasQuery(recipeID, recipeIterationID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(iterationMediasTableColumns, fmt.Sprintf(countQuery, iterationMediasTableName))...).
		From(iterationMediasTableName).
		Join(recipeIterationsOnIterationMediasJoinClause).
		Join(recipesOnRecipeIterationsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", iterationMediasTableName):                               nil,
			fmt.Sprintf("%s.id", recipesTableName):                                                recipeID,
			fmt.Sprintf("%s.id", recipeIterationsTableName):                                       recipeIterationID,
			fmt.Sprintf("%s.%s", recipeIterationsTableName, recipeIterationsTableOwnershipColumn): recipeID,
			fmt.Sprintf("%s.%s", iterationMediasTableName, iterationMediasTableOwnershipColumn):   recipeIterationID,
		}).
		GroupBy(fmt.Sprintf("%s.id", iterationMediasTableName))

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

	iterationMedias, count, err := p.scanIterationMedias(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.IterationMediaList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		IterationMedias: iterationMedias,
	}

	return list, nil
}

// buildCreateIterationMediaQuery takes an iteration media and returns a creation query for that iteration media and the relevant arguments.
func (p *Postgres) buildCreateIterationMediaQuery(input *models.IterationMedia) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(iterationMediasTableName).
		Columns(
			"source",
			"mimetype",
			iterationMediasTableOwnershipColumn,
		).
		Values(
			input.Source,
			input.Mimetype,
			input.BelongsToRecipeIteration,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateIterationMedia creates an iteration media in the database.
func (p *Postgres) CreateIterationMedia(ctx context.Context, input *models.IterationMediaCreationInput) (*models.IterationMedia, error) {
	x := &models.IterationMedia{
		Source:                   input.Source,
		Mimetype:                 input.Mimetype,
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
		Set("source", input.Source).
		Set("mimetype", input.Mimetype).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                                input.ID,
			iterationMediasTableOwnershipColumn: input.BelongsToRecipeIteration,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateIterationMedia updates a particular iteration media. Note that UpdateIterationMedia expects the provided input to have a valid ID.
func (p *Postgres) UpdateIterationMedia(ctx context.Context, input *models.IterationMedia) error {
	query, args := p.buildUpdateIterationMediaQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveIterationMediaQuery returns a SQL query which marks a given iteration media belonging to a given recipe iteration as archived.
func (p *Postgres) buildArchiveIterationMediaQuery(recipeIterationID, iterationMediaID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(iterationMediasTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                                iterationMediaID,
			"archived_on":                       nil,
			iterationMediasTableOwnershipColumn: recipeIterationID,
		}).
		Suffix("RETURNING archived_on").
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
