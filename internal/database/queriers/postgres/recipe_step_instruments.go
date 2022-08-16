package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	recipeStepsOnRecipeStepInstrumentsJoinClause      = "recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id"
	recipeStepInstrumentsOnValidInstrumentsJoinClause = "valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id"
)

var (
	_ types.RecipeStepInstrumentDataManager = (*SQLQuerier)(nil)

	// recipeStepInstrumentsTableColumns are the columns for the recipe_step_instruments table.
	recipeStepInstrumentsTableColumns = []string{
		"recipe_step_instruments.id",
		"valid_instruments.id",
		"valid_instruments.name",
		"valid_instruments.variant",
		"valid_instruments.description",
		"valid_instruments.icon_path",
		"valid_instruments.created_on",
		"valid_instruments.last_updated_on",
		"valid_instruments.archived_on",
		"recipe_step_instruments.recipe_step_product_id",
		"recipe_step_instruments.name",
		"recipe_step_instruments.product_of_recipe_step",
		"recipe_step_instruments.notes",
		"recipe_step_instruments.preference_rank",
		"recipe_step_instruments.created_on",
		"recipe_step_instruments.last_updated_on",
		"recipe_step_instruments.archived_on",
		"recipe_step_instruments.belongs_to_recipe_step",
	}

	getRecipeStepInstrumentsJoins = []string{
		recipeStepsOnRecipeStepInstrumentsJoinClause,
		recipeStepInstrumentsOnValidInstrumentsJoinClause,
		recipesOnRecipeStepsJoinClause,
	}
)

// scanRecipeStepInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step instrument struct.
func (q *SQLQuerier) scanRecipeStepInstrument(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.RecipeStepInstrument{
		Instrument: &types.ValidInstrument{},
	}

	targetVars := []interface{}{
		&x.ID,
		&x.Instrument.ID,
		&x.Instrument.Name,
		&x.Instrument.Variant,
		&x.Instrument.Description,
		&x.Instrument.IconPath,
		&x.Instrument.CreatedOn,
		&x.Instrument.LastUpdatedOn,
		&x.Instrument.ArchivedOn,
		&x.RecipeStepProductID,
		&x.Name,
		&x.ProductOfRecipeStep,
		&x.Notes,
		&x.PreferenceRank,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipeStep,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepInstruments takes some database rows and turns them into a slice of recipe step instruments.
func (q *SQLQuerier) scanRecipeStepInstruments(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepInstruments []*types.RecipeStepInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepInstrument(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		recipeStepInstruments = append(recipeStepInstruments, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipeStepInstruments, filteredCount, totalCount, nil
}

const recipeStepInstrumentExistenceQuery = "SELECT EXISTS ( SELECT recipe_step_instruments.id FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.archived_on IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5 )"

// RecipeStepInstrumentExists fetches whether a recipe step instrument exists from the database.
func (q *SQLQuerier) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	args := []interface{}{
		recipeStepID,
		recipeStepInstrumentID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepInstrumentExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe step instrument existence check")
	}

	return result, nil
}

const getRecipeStepInstrumentQuery = `SELECT
	recipe_step_instruments.id,
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.variant,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.created_on,
	valid_instruments.last_updated_on,
	valid_instruments.archived_on,
	recipe_step_instruments.recipe_step_product_id,
	recipe_step_instruments.name,
	recipe_step_instruments.product_of_recipe_step,
	recipe_step_instruments.notes,
	recipe_step_instruments.preference_rank,
	recipe_step_instruments.created_on,
	recipe_step_instruments.last_updated_on,
	recipe_step_instruments.archived_on,
	recipe_step_instruments.belongs_to_recipe_step
FROM recipe_step_instruments
LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_instruments.archived_on IS NULL
  AND recipe_step_instruments.belongs_to_recipe_step = $1
  AND recipe_step_instruments.id = $2
  AND recipe_steps.archived_on IS NULL
  AND recipe_steps.belongs_to_recipe = $3
  AND recipe_steps.id = $4
  AND recipes.archived_on IS NULL
  AND recipes.id = $5`

// GetRecipeStepInstrument fetches a recipe step instrument from the database.
func (q *SQLQuerier) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	args := []interface{}{
		recipeStepID,
		recipeStepInstrumentID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "recipe step instrument", getRecipeStepInstrumentQuery, args)

	recipeStepInstrument, _, _, err := q.scanRecipeStepInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipeStepInstrument")
	}

	return recipeStepInstrument, nil
}

const getTotalRecipeStepInstrumentsCountQuery = "SELECT COUNT(recipe_step_instruments.id) FROM recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL"

// GetTotalRecipeStepInstrumentCount fetches the count of recipe step instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalRecipeStepInstrumentCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getTotalRecipeStepInstrumentsCountQuery, "fetching count of recipe step instruments")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of recipe step instruments")
	}

	return count, nil
}

// GetRecipeStepInstruments fetches a list of recipe step instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.RecipeStepInstrumentList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	x = &types.RecipeStepInstrumentList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "recipe_step_instruments", getRecipeStepInstrumentsJoins, []string{"valid_instruments.id", "recipe_step_instruments.id"}, nil, householdOwnershipColumn, recipeStepInstrumentsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "recipeStepInstruments", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe step instruments list retrieval query")
	}

	if x.RecipeStepInstruments, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepInstruments(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step instruments")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetRecipeStepInstrumentsWithIDsQuery(ctx context.Context, recipeStepID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"recipe_step_instruments.id":                     ids,
		"recipe_step_instruments.archived_on":            nil,
		"recipe_step_instruments.belongs_to_recipe_step": recipeStepID,
	}

	subqueryBuilder := q.sqlBuilder.Select(recipeStepInstrumentsTableColumns...).
		From("recipe_step_instruments").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(recipeStepInstrumentsTableColumns...).
		FromSelect(subqueryBuilder, "recipe_step_instruments").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetRecipeStepInstrumentsWithIDs fetches recipe step instruments from the database within a given set of IDs.
func (q *SQLQuerier) GetRecipeStepInstrumentsWithIDs(ctx context.Context, recipeStepID string, limit uint8, ids []string) ([]*types.RecipeStepInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if ids == nil {
		return nil, ErrNilInputProvided
	}

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.buildGetRecipeStepInstrumentsWithIDsQuery(ctx, recipeStepID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "recipe step instruments with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe step instruments from database")
	}

	recipeStepInstruments, _, _, err := q.scanRecipeStepInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step instruments")
	}

	return recipeStepInstruments, nil
}

const getRecipeStepInstrumentsForRecipeQuery = `SELECT
	recipe_step_instruments.id,
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.variant,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.created_on,
	valid_instruments.last_updated_on,
	valid_instruments.archived_on,
	recipe_step_instruments.recipe_step_product_id,
	recipe_step_instruments.name,
	recipe_step_instruments.product_of_recipe_step,
	recipe_step_instruments.notes,
	recipe_step_instruments.preference_rank,
	recipe_step_instruments.created_on,
	recipe_step_instruments.last_updated_on,
	recipe_step_instruments.archived_on,
	recipe_step_instruments.belongs_to_recipe_step
FROM recipe_step_instruments
LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_instruments.archived_on IS NULL
AND recipe_steps.archived_on IS NULL
AND recipe_steps.belongs_to_recipe = $1
AND recipes.archived_on IS NULL 
AND recipes.id = $2
`

// getRecipeStepInstrumentsForRecipe fetches a list of recipe step instruments from the database that meet a particular filter.
func (q *SQLQuerier) getRecipeStepInstrumentsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	args := []interface{}{
		recipeID,
		recipeID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "recipe step instruments", getRecipeStepInstrumentsForRecipeQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe step instruments list retrieval query")
	}

	recipeStepInstruments, _, _, err := q.scanRecipeStepInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step instruments")
	}

	return recipeStepInstruments, nil
}

const recipeStepInstrumentCreationQuery = `INSERT INTO recipe_step_instruments (id,instrument_id,recipe_step_product_id,name,product_of_recipe_step,notes,preference_rank,belongs_to_recipe_step) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

// CreateRecipeStepInstrument creates a recipe step instrument in the database.
func (q *SQLQuerier) createRecipeStepInstrument(ctx context.Context, querier database.SQLQueryExecutor, input *types.RecipeStepInstrumentDatabaseCreationInput) (*types.RecipeStepInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepInstrumentIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.InstrumentID,
		input.RecipeStepProductID,
		input.Name,
		input.ProductOfRecipeStep,
		input.Notes,
		input.PreferenceRank,
		input.BelongsToRecipeStep,
	}

	// create the recipe step instrument.
	if err := q.performWriteQuery(ctx, querier, "recipe step instrument creation", recipeStepInstrumentCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing recipe step instrument creation query")
	}

	x := &types.RecipeStepInstrument{
		ID:                  input.ID,
		Instrument:          nil,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                input.Name,
		ProductOfRecipeStep: input.ProductOfRecipeStep,
		Notes:               input.Notes,
		PreferenceRank:      input.PreferenceRank,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		CreatedOn:           q.currentTime(),
	}

	if input.InstrumentID != nil {
		x.Instrument = &types.ValidInstrument{ID: *input.InstrumentID}
	}

	tracing.AttachRecipeStepInstrumentIDToSpan(span, x.ID)
	logger.Info("recipe step instrument created")

	return x, nil
}

// CreateRecipeStepInstrument creates a recipe step instrument in the database.
func (q *SQLQuerier) CreateRecipeStepInstrument(ctx context.Context, input *types.RecipeStepInstrumentDatabaseCreationInput) (*types.RecipeStepInstrument, error) {
	return q.createRecipeStepInstrument(ctx, q.db, input)
}

const updateRecipeStepInstrumentQuery = `UPDATE recipe_step_instruments SET
   instrument_id = $1,
   recipe_step_product_id = $2,
   name = $3,
   product_of_recipe_step = $4,
   notes = $5,
   preference_rank = $6,
   last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
  AND belongs_to_recipe_step = $7
  AND id = $8`

// UpdateRecipeStepInstrument updates a particular recipe step instrument.
func (q *SQLQuerier) UpdateRecipeStepInstrument(ctx context.Context, updated *types.RecipeStepInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepInstrumentIDKey, updated.ID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, updated.ID)

	var instrumentID *string
	if updated.Instrument != nil {
		instrumentID = &updated.Instrument.ID
	}

	args := []interface{}{
		instrumentID,
		updated.RecipeStepProductID,
		updated.Name,
		updated.ProductOfRecipeStep,
		updated.Notes,
		updated.PreferenceRank,
		updated.BelongsToRecipeStep,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step instrument update", updateRecipeStepInstrumentQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step instrument")
	}

	logger.Info("recipe step instrument updated")

	return nil
}

const archiveRecipeStepInstrumentQuery = "UPDATE recipe_step_instruments SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2"

// ArchiveRecipeStepInstrument archives a recipe step instrument from the database by its ID.
func (q *SQLQuerier) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)

	args := []interface{}{
		recipeStepID,
		recipeStepInstrumentID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step instrument archive", archiveRecipeStepInstrumentQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step instrument")
	}

	logger.Info("recipe step instrument archived")

	return nil
}
