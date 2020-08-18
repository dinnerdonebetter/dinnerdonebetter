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
	recipeStepInstrumentsTableName               = "recipe_step_instruments"
	recipeStepInstrumentsTableInstrumentIDColumn = "instrument_id"
	recipeStepInstrumentsTableRecipeStepIDColumn = "recipe_step_id"
	recipeStepInstrumentsTableNotesColumn        = "notes"
	recipeStepInstrumentsTableOwnershipColumn    = "belongs_to_recipe_step"
)

var (
	recipeStepInstrumentsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, idColumn),
		fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, recipeStepInstrumentsTableInstrumentIDColumn),
		fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, recipeStepInstrumentsTableRecipeStepIDColumn),
		fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, recipeStepInstrumentsTableNotesColumn),
		fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, recipeStepInstrumentsTableOwnershipColumn),
	}
)

// scanRecipeStepInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Instrument struct
func (p *Postgres) scanRecipeStepInstrument(scan database.Scanner) (*models.RecipeStepInstrument, error) {
	x := &models.RecipeStepInstrument{}

	targetVars := []interface{}{
		&x.ID,
		&x.InstrumentID,
		&x.RecipeStepID,
		&x.Notes,
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

// scanRecipeStepInstruments takes a logger and some database rows and turns them into a slice of recipe step instruments.
func (p *Postgres) scanRecipeStepInstruments(rows database.ResultIterator) ([]models.RecipeStepInstrument, error) {
	var (
		list []models.RecipeStepInstrument
	)

	for rows.Next() {
		x, err := p.scanRecipeStepInstrument(rows)
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

// buildRecipeStepInstrumentExistsQuery constructs a SQL query for checking if a recipe step instrument with a given ID belong to a a recipe step with a given ID exists
func (p *Postgres) buildRecipeStepInstrumentExistsQuery(recipeID, recipeStepID, recipeStepInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, idColumn)).
		Prefix(existencePrefix).
		From(recipeStepInstrumentsTableName).
		Join(recipeStepsOnRecipeStepInstrumentsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, idColumn):                                  recipeStepInstrumentID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                            recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                     recipeID,
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, recipeStepInstrumentsTableOwnershipColumn): recipeStepID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeStepInstrumentExists queries the database to see if a given recipe step instrument belonging to a given user exists.
func (p *Postgres) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (exists bool, err error) {
	query, args := p.buildRecipeStepInstrumentExistsQuery(recipeID, recipeStepID, recipeStepInstrumentID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRecipeStepInstrumentQuery constructs a SQL query for fetching a recipe step instrument with a given ID belong to a recipe step with a given ID.
func (p *Postgres) buildGetRecipeStepInstrumentQuery(recipeID, recipeStepID, recipeStepInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipeStepInstrumentsTableColumns...).
		From(recipeStepInstrumentsTableName).
		Join(recipeStepsOnRecipeStepInstrumentsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, idColumn):                                  recipeStepInstrumentID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                            recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                     recipeID,
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, recipeStepInstrumentsTableOwnershipColumn): recipeStepID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepInstrument fetches a recipe step instrument from the database.
func (p *Postgres) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (*models.RecipeStepInstrument, error) {
	query, args := p.buildGetRecipeStepInstrumentQuery(recipeID, recipeStepID, recipeStepInstrumentID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanRecipeStepInstrument(row)
}

var (
	allRecipeStepInstrumentsCountQueryBuilder sync.Once
	allRecipeStepInstrumentsCountQuery        string
)

// buildGetAllRecipeStepInstrumentsCountQuery returns a query that fetches the total number of recipe step instruments in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipeStepInstrumentsCountQuery() string {
	allRecipeStepInstrumentsCountQueryBuilder.Do(func() {
		var err error

		allRecipeStepInstrumentsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipeStepInstrumentsTableName)).
			From(recipeStepInstrumentsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeStepInstrumentsCountQuery
}

// GetAllRecipeStepInstrumentsCount will fetch the count of recipe step instruments from the database.
func (p *Postgres) GetAllRecipeStepInstrumentsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeStepInstrumentsCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfRecipeStepInstrumentsQuery returns a query that fetches every recipe step instrument in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfRecipeStepInstrumentsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(recipeStepInstrumentsTableColumns...).
		From(recipeStepInstrumentsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllRecipeStepInstruments fetches every recipe step instrument from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllRecipeStepInstruments(ctx context.Context, resultChannel chan []models.RecipeStepInstrument) error {
	count, err := p.GetAllRecipeStepInstrumentsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfRecipeStepInstrumentsQuery(begin, end)
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

			recipeStepInstruments, err := p.scanRecipeStepInstruments(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- recipeStepInstruments
		}(beginID, endID)
	}

	return nil
}

// buildGetRecipeStepInstrumentsQuery builds a SQL query selecting recipe step instruments that adhere to a given QueryFilter and belong to a given recipe step,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeStepInstrumentsQuery(recipeID, recipeStepID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(recipeStepInstrumentsTableColumns...).
		From(recipeStepInstrumentsTableName).
		Join(recipeStepsOnRecipeStepInstrumentsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, archivedOnColumn):                          nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                            recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                     recipeID,
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, recipeStepInstrumentsTableOwnershipColumn): recipeStepID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipeStepInstrumentsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepInstruments fetches a list of recipe step instruments from the database that meet a particular filter.
func (p *Postgres) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepInstrumentList, error) {
	query, args := p.buildGetRecipeStepInstrumentsQuery(recipeID, recipeStepID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step instruments")
	}

	recipeStepInstruments, err := p.scanRecipeStepInstruments(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeStepInstrumentList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		RecipeStepInstruments: recipeStepInstruments,
	}

	return list, nil
}

// buildGetRecipeStepInstrumentsWithIDsQuery builds a SQL query selecting recipeStepInstruments that belong to a given recipe step,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetRecipeStepInstrumentsWithIDsQuery(recipeID, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(recipeStepInstrumentsTableColumns...).
		From(recipeStepInstrumentsTableName).
		Join(recipeStepsOnRecipeStepInstrumentsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, archivedOnColumn):                          nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                                recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                            recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):                     recipeID,
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, recipeStepInstrumentsTableOwnershipColumn): recipeStepID,
		}).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(recipeStepInstrumentsTableColumns...).
		FromSelect(subqueryBuilder, recipeStepInstrumentsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepInstrumentsWithIDs fetches a list of recipe step instruments from the database that exist within a given set of IDs.
func (p *Postgres) GetRecipeStepInstrumentsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepInstrument, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetRecipeStepInstrumentsWithIDsQuery(recipeID, recipeStepID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step instruments")
	}

	recipeStepInstruments, err := p.scanRecipeStepInstruments(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return recipeStepInstruments, nil
}

// buildCreateRecipeStepInstrumentQuery takes a recipe step instrument and returns a creation query for that recipe step instrument and the relevant arguments.
func (p *Postgres) buildCreateRecipeStepInstrumentQuery(input *models.RecipeStepInstrument) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeStepInstrumentsTableName).
		Columns(
			recipeStepInstrumentsTableInstrumentIDColumn,
			recipeStepInstrumentsTableRecipeStepIDColumn,
			recipeStepInstrumentsTableNotesColumn,
			recipeStepInstrumentsTableOwnershipColumn,
		).
		Values(
			input.InstrumentID,
			input.RecipeStepID,
			input.Notes,
			input.BelongsToRecipeStep,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepInstrument creates a recipe step instrument in the database.
func (p *Postgres) CreateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrumentCreationInput) (*models.RecipeStepInstrument, error) {
	x := &models.RecipeStepInstrument{
		InstrumentID:        input.InstrumentID,
		RecipeStepID:        input.RecipeStepID,
		Notes:               input.Notes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
	}

	query, args := p.buildCreateRecipeStepInstrumentQuery(x)

	// create the recipe step instrument.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step instrument creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeStepInstrumentQuery takes a recipe step instrument and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeStepInstrumentQuery(input *models.RecipeStepInstrument) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepInstrumentsTableName).
		Set(recipeStepInstrumentsTableInstrumentIDColumn, input.InstrumentID).
		Set(recipeStepInstrumentsTableRecipeStepIDColumn, input.RecipeStepID).
		Set(recipeStepInstrumentsTableNotesColumn, input.Notes).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
			recipeStepInstrumentsTableOwnershipColumn: input.BelongsToRecipeStep,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepInstrument updates a particular recipe step instrument. Note that UpdateRecipeStepInstrument expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeStepInstrument(ctx context.Context, input *models.RecipeStepInstrument) error {
	query, args := p.buildUpdateRecipeStepInstrumentQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveRecipeStepInstrumentQuery returns a SQL query which marks a given recipe step instrument belonging to a given recipe step as archived.
func (p *Postgres) buildArchiveRecipeStepInstrumentQuery(recipeStepID, recipeStepInstrumentID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepInstrumentsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         recipeStepInstrumentID,
			archivedOnColumn: nil,
			recipeStepInstrumentsTableOwnershipColumn: recipeStepID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepInstrument marks a recipe step instrument as archived in the database.
func (p *Postgres) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID uint64) error {
	query, args := p.buildArchiveRecipeStepInstrumentQuery(recipeStepID, recipeStepInstrumentID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
