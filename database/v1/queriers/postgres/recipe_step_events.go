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
	recipeStepEventsTableName                    = "recipe_step_events"
	recipeStepEventsTableEventTypeColumn         = "event_type"
	recipeStepEventsTableDoneColumn              = "done"
	recipeStepEventsTableRecipeIterationIDColumn = "recipe_iteration_id"
	recipeStepEventsTableRecipeStepIDColumn      = "recipe_step_id"
	recipeStepEventsTableOwnershipColumn         = "belongs_to_recipe_step"
)

var (
	recipeStepEventsTableColumns = []string{
		fmt.Sprintf("%s.%s", recipeStepEventsTableName, idColumn),
		fmt.Sprintf("%s.%s", recipeStepEventsTableName, recipeStepEventsTableEventTypeColumn),
		fmt.Sprintf("%s.%s", recipeStepEventsTableName, recipeStepEventsTableDoneColumn),
		fmt.Sprintf("%s.%s", recipeStepEventsTableName, recipeStepEventsTableRecipeIterationIDColumn),
		fmt.Sprintf("%s.%s", recipeStepEventsTableName, recipeStepEventsTableRecipeStepIDColumn),
		fmt.Sprintf("%s.%s", recipeStepEventsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", recipeStepEventsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepEventsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", recipeStepEventsTableName, recipeStepEventsTableOwnershipColumn),
	}
)

// scanRecipeStepEvent takes a database Scanner (i.e. *sql.Row) and scans the result into a Recipe Step Event struct
func (p *Postgres) scanRecipeStepEvent(scan database.Scanner) (*models.RecipeStepEvent, error) {
	x := &models.RecipeStepEvent{}

	targetVars := []interface{}{
		&x.ID,
		&x.EventType,
		&x.Done,
		&x.RecipeIterationID,
		&x.RecipeStepID,
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

// scanRecipeStepEvents takes a logger and some database rows and turns them into a slice of recipe step events.
func (p *Postgres) scanRecipeStepEvents(rows database.ResultIterator) ([]models.RecipeStepEvent, error) {
	var (
		list []models.RecipeStepEvent
	)

	for rows.Next() {
		x, err := p.scanRecipeStepEvent(rows)
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

// buildRecipeStepEventExistsQuery constructs a SQL query for checking if a recipe step event with a given ID belong to a a recipe step with a given ID exists
func (p *Postgres) buildRecipeStepEventExistsQuery(recipeID, recipeStepID, recipeStepEventID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", recipeStepEventsTableName, idColumn)).
		Prefix(existencePrefix).
		From(recipeStepEventsTableName).
		Join(recipeStepsOnRecipeStepEventsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, idColumn):                             recipeStepEventID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                  recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):           recipeID,
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, recipeStepEventsTableOwnershipColumn): recipeStepID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// RecipeStepEventExists queries the database to see if a given recipe step event belonging to a given user exists.
func (p *Postgres) RecipeStepEventExists(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (exists bool, err error) {
	query, args := p.buildRecipeStepEventExistsQuery(recipeID, recipeStepID, recipeStepEventID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetRecipeStepEventQuery constructs a SQL query for fetching a recipe step event with a given ID belong to a recipe step with a given ID.
func (p *Postgres) buildGetRecipeStepEventQuery(recipeID, recipeStepID, recipeStepEventID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(recipeStepEventsTableColumns...).
		From(recipeStepEventsTableName).
		Join(recipeStepsOnRecipeStepEventsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, idColumn):                             recipeStepEventID,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                  recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):           recipeID,
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, recipeStepEventsTableOwnershipColumn): recipeStepID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepEvent fetches a recipe step event from the database.
func (p *Postgres) GetRecipeStepEvent(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (*models.RecipeStepEvent, error) {
	query, args := p.buildGetRecipeStepEventQuery(recipeID, recipeStepID, recipeStepEventID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanRecipeStepEvent(row)
}

var (
	allRecipeStepEventsCountQueryBuilder sync.Once
	allRecipeStepEventsCountQuery        string
)

// buildGetAllRecipeStepEventsCountQuery returns a query that fetches the total number of recipe step events in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllRecipeStepEventsCountQuery() string {
	allRecipeStepEventsCountQueryBuilder.Do(func() {
		var err error

		allRecipeStepEventsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, recipeStepEventsTableName)).
			From(recipeStepEventsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", recipeStepEventsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allRecipeStepEventsCountQuery
}

// GetAllRecipeStepEventsCount will fetch the count of recipe step events from the database.
func (p *Postgres) GetAllRecipeStepEventsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllRecipeStepEventsCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfRecipeStepEventsQuery returns a query that fetches every recipe step event in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfRecipeStepEventsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(recipeStepEventsTableColumns...).
		From(recipeStepEventsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllRecipeStepEvents fetches every recipe step event from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllRecipeStepEvents(ctx context.Context, resultChannel chan []models.RecipeStepEvent) error {
	count, err := p.GetAllRecipeStepEventsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfRecipeStepEventsQuery(begin, end)
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

			recipeStepEvents, err := p.scanRecipeStepEvents(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- recipeStepEvents
		}(beginID, endID)
	}

	return nil
}

// buildGetRecipeStepEventsQuery builds a SQL query selecting recipe step events that adhere to a given QueryFilter and belong to a given recipe step,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetRecipeStepEventsQuery(recipeID, recipeStepID uint64, filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(recipeStepEventsTableColumns...).
		From(recipeStepEventsTableName).
		Join(recipeStepsOnRecipeStepEventsJoinClause).
		Join(recipesOnRecipeStepsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, archivedOnColumn):                     nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                  recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):           recipeID,
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, recipeStepEventsTableOwnershipColumn): recipeStepID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", recipeStepEventsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, recipeStepEventsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepEvents fetches a list of recipe step events from the database that meet a particular filter.
func (p *Postgres) GetRecipeStepEvents(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*models.RecipeStepEventList, error) {
	query, args := p.buildGetRecipeStepEventsQuery(recipeID, recipeStepID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step events")
	}

	recipeStepEvents, err := p.scanRecipeStepEvents(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.RecipeStepEventList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		RecipeStepEvents: recipeStepEvents,
	}

	return list, nil
}

// buildGetRecipeStepEventsWithIDsQuery builds a SQL query selecting recipeStepEvents that belong to a given recipe step,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetRecipeStepEventsWithIDsQuery(recipeID, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(recipeStepEventsTableColumns...).
		From(recipeStepEventsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, archivedOnColumn):                     nil,
			fmt.Sprintf("%s.%s", recipesTableName, idColumn):                                      recipeID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, idColumn):                                  recipeStepID,
			fmt.Sprintf("%s.%s", recipeStepsTableName, recipeStepsTableOwnershipColumn):           recipeID,
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, recipeStepEventsTableOwnershipColumn): recipeStepID,
		}).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(recipeStepEventsTableColumns...).
		FromSelect(subqueryBuilder, recipeStepEventsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", recipeStepEventsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetRecipeStepEventsWithIDs fetches a list of recipe step events from the database that exist within a given set of IDs.
func (p *Postgres) GetRecipeStepEventsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]models.RecipeStepEvent, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetRecipeStepEventsWithIDsQuery(recipeID, recipeStepID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for recipe step events")
	}

	recipeStepEvents, err := p.scanRecipeStepEvents(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return recipeStepEvents, nil
}

// buildCreateRecipeStepEventQuery takes a recipe step event and returns a creation query for that recipe step event and the relevant arguments.
func (p *Postgres) buildCreateRecipeStepEventQuery(input *models.RecipeStepEvent) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(recipeStepEventsTableName).
		Columns(
			recipeStepEventsTableEventTypeColumn,
			recipeStepEventsTableDoneColumn,
			recipeStepEventsTableRecipeIterationIDColumn,
			recipeStepEventsTableRecipeStepIDColumn,
			recipeStepEventsTableOwnershipColumn,
		).
		Values(
			input.EventType,
			input.Done,
			input.RecipeIterationID,
			input.RecipeStepID,
			input.BelongsToRecipeStep,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateRecipeStepEvent creates a recipe step event in the database.
func (p *Postgres) CreateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEventCreationInput) (*models.RecipeStepEvent, error) {
	x := &models.RecipeStepEvent{
		EventType:           input.EventType,
		Done:                input.Done,
		RecipeIterationID:   input.RecipeIterationID,
		RecipeStepID:        input.RecipeStepID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
	}

	query, args := p.buildCreateRecipeStepEventQuery(x)

	// create the recipe step event.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing recipe step event creation query: %w", err)
	}

	return x, nil
}

// buildUpdateRecipeStepEventQuery takes a recipe step event and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateRecipeStepEventQuery(input *models.RecipeStepEvent) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepEventsTableName).
		Set(recipeStepEventsTableEventTypeColumn, input.EventType).
		Set(recipeStepEventsTableDoneColumn, input.Done).
		Set(recipeStepEventsTableRecipeIterationIDColumn, input.RecipeIterationID).
		Set(recipeStepEventsTableRecipeStepIDColumn, input.RecipeStepID).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             input.ID,
			recipeStepEventsTableOwnershipColumn: input.BelongsToRecipeStep,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateRecipeStepEvent updates a particular recipe step event. Note that UpdateRecipeStepEvent expects the provided input to have a valid ID.
func (p *Postgres) UpdateRecipeStepEvent(ctx context.Context, input *models.RecipeStepEvent) error {
	query, args := p.buildUpdateRecipeStepEventQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveRecipeStepEventQuery returns a SQL query which marks a given recipe step event belonging to a given recipe step as archived.
func (p *Postgres) buildArchiveRecipeStepEventQuery(recipeStepID, recipeStepEventID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(recipeStepEventsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             recipeStepEventID,
			archivedOnColumn:                     nil,
			recipeStepEventsTableOwnershipColumn: recipeStepID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveRecipeStepEvent marks a recipe step event as archived in the database.
func (p *Postgres) ArchiveRecipeStepEvent(ctx context.Context, recipeStepID, recipeStepEventID uint64) error {
	query, args := p.buildArchiveRecipeStepEventQuery(recipeStepID, recipeStepEventID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
