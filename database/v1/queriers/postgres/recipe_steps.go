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
	recipeStepsTableName                            = "recipe_steps"
	recipeStepsTableIndexColumn                     = "index"
	recipeStepsTablePreparationIDColumn             = "preparation_id"
	recipeStepsTablePrerequisiteStepColumn          = "prerequisite_step"
	recipeStepsTableMinEstimatedTimeInSecondsColumn = "min_estimated_time_in_seconds"
	recipeStepsTableMaxEstimatedTimeInSecondsColumn = "max_estimated_time_in_seconds"
	recipeStepsTableTemperatureInCelsiusColumn      = "temperature_in_celsius"
	recipeStepsTableNotesColumn                     = "notes"
	recipeStepsTableRecipeIDColumn                  = "recipe_id"
	recipeStepsTableOwnershipColumn                 = "belongs_to_recipe"
)

var (
	recipeStepsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableIndexColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTablePreparationIDColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTablePrerequisiteStepColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableMinEstimatedTimeInSecondsColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableMaxEstimatedTimeInSecondsColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableTemperatureInCelsiusColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableNotesColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableRecipeIDColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn),
	}

	recipeStepsOnRecipeStepInstrumentsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.%s", recipeStepsTableName, recipeStepInstrumentsTableName, recipeStepInstrumentsTableOwnershipColumn, recipeStepsTableName, idColumn)
	recipeStepsOnRecipeStepIngredientsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.%s", recipeStepsTableName, recipeStepIngredientsTableName, recipeStepIngredientsTableOwnershipColumn, recipeStepsTableName, idColumn)
	recipeStepsOnRecipeStepProductsJoinClause    = fmt.Sprintf("%s ON %s.%s=%s.%s", recipeStepsTableName, recipeStepProductsTableName, recipeStepProductsTableOwnershipColumn, recipeStepsTableName, idColumn)
	recipeStepsOnRecipeStepEventsJoinClause      = fmt.Sprintf("%s ON %s.%s=%s.%s", recipeStepsTableName, recipeStepEventsTableName, recipeStepEventsTableOwnershipColumn, recipeStepsTableName, idColumn)
)

// scanRecipeStep takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step struct
func (p *Postgres) scanRecipeStep(scan database.Scanner) (*models.RecipeStep, error) {
	x := &models.RecipeStep{}

	targetVars := []interface{}{
		&x.ID,
		&x.Index,
		&x.PreparationID,
		&x.PrerequisiteStep,
		&x.MinEstimatedTimeInSeconds,
		&x.MaxEstimatedTimeInSeconds,
		&x.TemperatureInCelsius,
		&x.Notes,
		&x.RecipeID,
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

// scanRecipeSteps takes a logger and some database rows and turns them into a slice of recipe steps.
func (p *Postgres) scanRecipeSteps(rows database.ResultIterator) ([]models.RecipeStep, error) {
	var (
		list []models.RecipeStep
	)

	for rows.Next() {
		x, err := p.scanRecipeStep(rows)
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

// buildRecipeStepExistsQuery constructs a SQL query for checking if a recipe step with a given ID belong to a a recipe with a given ID exists
func (p *Postgres) buildRecipeStepExistsQuery(recipeID, recipeStepID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn)).
		Prefix(existencePrefix).
		From(recipeStepsTableName).
		Join(recipesOnRecipeStepsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                        recipeStepID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                            recipeID,
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
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                        recipeStepID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                            recipeID,
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
	return p.scanRecipeStep(row)
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
				fmt.Sprintf("%s.%s", recipeStepsTableName, archivedOnColumn): nil,
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

// buildGetBatchOfRecipeStepsQuery returns a query that fetches every recipe step in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfRecipeStepsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(recipeStepsTableColumns...).
		From(recipeStepsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllRecipeSteps fetches every recipe step from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllRecipeSteps(ctx context.Context, resultChannel chan []models.RecipeStep) error {
	count, err := p.GetAllRecipeStepsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfRecipeStepsQuery(begin, end)
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

			recipeSteps, err := p.scanRecipeSteps(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- recipeSteps
		}(beginID, endID)
	}

	return nil
}

// buildGetRecipeStepsQuery builds a SQL query selecting recipe steps that adhere to a given QueryFilter and belong to a given recipe,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeStepsQuery(recipeID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(recipeStepsTableColumns...).
		From(recipeStepsTableName).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepsTableName, archivedOnColumn):                nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                            recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn): recipeID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn))

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

	recipeSteps, err := p.scanRecipeSteps(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeStepList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		RecipeSteps: recipeSteps,
	}

	return list, nil
}

// buildGetRecipeStepsWithIDsQuery builds a SQL query selecting recipeSteps that belong to a given recipe,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetRecipeStepsWithIDsQuery(recipeID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(recipeStepsTableColumns...).
		From(recipeStepsTableName).
		Join(recipesOnRecipeStepsJoinClause).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepsTableName, archivedOnColumn):                nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                            recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn): recipeID,
		}).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(recipeStepsTableColumns...).
		FromSelect(subqueryBuilder, recipeStepsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepsWithIDs fetches a list of recipe steps from the database that exist within a given set of IDs.
func (p *Postgres) GetRecipeStepsWithIDs(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) ([]models.RecipeStep, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetRecipeStepsWithIDsQuery(recipeID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe steps")
	}

	recipeSteps, err := p.scanRecipeSteps(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return recipeSteps, nil
}

// buildCreateRecipeStepQuery takes a recipe step and returns a creation query for that recipe step and the relevant arguments.
func (p *Postgres) buildCreateRecipeStepQuery(input *models.RecipeStep) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeStepsTableName).
		Columns(
			recipeStepsTableIndexColumn,
			recipeStepsTablePreparationIDColumn,
			recipeStepsTablePrerequisiteStepColumn,
			recipeStepsTableMinEstimatedTimeInSecondsColumn,
			recipeStepsTableMaxEstimatedTimeInSecondsColumn,
			recipeStepsTableTemperatureInCelsiusColumn,
			recipeStepsTableNotesColumn,
			recipeStepsTableRecipeIDColumn,
			recipeStepsTableOwnershipColumn,
		).
		Values(
			input.Index,
			input.PreparationID,
			input.PrerequisiteStep,
			input.MinEstimatedTimeInSeconds,
			input.MaxEstimatedTimeInSeconds,
			input.TemperatureInCelsius,
			input.Notes,
			input.RecipeID,
			input.BelongsToRecipe,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStep creates a recipe step in the database.
func (p *Postgres) CreateRecipeStep(ctx context.Context, input *models.RecipeStepCreationInput) (*models.RecipeStep, error) {
	x := &models.RecipeStep{
		Index:                     input.Index,
		PreparationID:             input.PreparationID,
		PrerequisiteStep:          input.PrerequisiteStep,
		MinEstimatedTimeInSeconds: input.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: input.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      input.TemperatureInCelsius,
		Notes:                     input.Notes,
		RecipeID:                  input.RecipeID,
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
		Set(recipeStepsTableIndexColumn, input.Index).
		Set(recipeStepsTablePreparationIDColumn, input.PreparationID).
		Set(recipeStepsTablePrerequisiteStepColumn, input.PrerequisiteStep).
		Set(recipeStepsTableMinEstimatedTimeInSecondsColumn, input.MinEstimatedTimeInSeconds).
		Set(recipeStepsTableMaxEstimatedTimeInSecondsColumn, input.MaxEstimatedTimeInSeconds).
		Set(recipeStepsTableTemperatureInCelsiusColumn, input.TemperatureInCelsius).
		Set(recipeStepsTableNotesColumn, input.Notes).
		Set(recipeStepsTableRecipeIDColumn, input.RecipeID).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                        input.ID,
			recipeStepsTableOwnershipColumn: input.BelongsToRecipe,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStep updates a particular recipe step. Note that UpdateRecipeStep expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeStep(ctx context.Context, input *models.RecipeStep) error {
	query, args := p.buildUpdateRecipeStepQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveRecipeStepQuery returns a SQL query which marks a given recipe step belonging to a given recipe as archived.
func (p *Postgres) buildArchiveRecipeStepQuery(recipeID, recipeStepID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                        recipeStepID,
			archivedOnColumn:                nil,
			recipeStepsTableOwnershipColumn: recipeID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
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
