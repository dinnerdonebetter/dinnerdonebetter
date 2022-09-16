package postgres

import (
	"context"
	_ "embed"

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
	_ types.RecipeStepInstrumentDataManager = (*Querier)(nil)

	// recipeStepInstrumentsTableColumns are the columns for the recipe_step_instruments table.
	recipeStepInstrumentsTableColumns = []string{
		"recipe_step_instruments.id",
		"valid_instruments.id",
		"valid_instruments.name",
		"valid_instruments.plural_name",
		"valid_instruments.description",
		"valid_instruments.icon_path",
		"valid_instruments.usable_for_storage",
		"valid_instruments.created_at",
		"valid_instruments.last_updated_at",
		"valid_instruments.archived_at",
		"recipe_step_instruments.recipe_step_product_id",
		"recipe_step_instruments.name",
		"recipe_step_instruments.product_of_recipe_step",
		"recipe_step_instruments.notes",
		"recipe_step_instruments.preference_rank",
		"recipe_step_instruments.optional",
		"recipe_step_instruments.minimum_quantity",
		"recipe_step_instruments.maximum_quantity",
		"recipe_step_instruments.created_at",
		"recipe_step_instruments.last_updated_at",
		"recipe_step_instruments.archived_at",
		"recipe_step_instruments.belongs_to_recipe_step",
	}

	getRecipeStepInstrumentsJoins = []string{
		recipeStepsOnRecipeStepInstrumentsJoinClause,
		recipeStepInstrumentsOnValidInstrumentsJoinClause,
		recipesOnRecipeStepsJoinClause,
	}
)

// scanRecipeStepInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step instrument struct.
func (q *Querier) scanRecipeStepInstrument(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStepInstrument{}
	instrument := types.NullableValidInstrument{}

	targetVars := []interface{}{
		&x.ID,
		&instrument.ID,
		&instrument.Name,
		&instrument.PluralName,
		&instrument.Description,
		&instrument.IconPath,
		&instrument.UsableForStorage,
		&instrument.CreatedAt,
		&instrument.LastUpdatedAt,
		&instrument.ArchivedAt,
		&x.RecipeStepProductID,
		&x.Name,
		&x.ProductOfRecipeStep,
		&x.Notes,
		&x.PreferenceRank,
		&x.Optional,
		&x.MinimumQuantity,
		&x.MaximumQuantity,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
		&x.BelongsToRecipeStep,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	if instrument.ID != nil {
		x.Instrument = instrument.ToValidInstrument()
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepInstruments takes some database rows and turns them into a slice of recipe step instruments.
func (q *Querier) scanRecipeStepInstruments(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepInstruments []*types.RecipeStepInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

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
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return recipeStepInstruments, filteredCount, totalCount, nil
}

//go:embed queries/recipe_step_instruments/recipe_step_instruments_exists.sql
var recipeStepInstrumentExistenceQuery string

// RecipeStepInstrumentExists fetches whether a recipe step instrument exists from the database.
func (q *Querier) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (exists bool, err error) {
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
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step instrument existence check")
	}

	return result, nil
}

//go:embed queries/recipe_step_instruments/recipe_step_instruments_get_one.sql
var getRecipeStepInstrumentQuery string

// GetRecipeStepInstrument fetches a recipe step instrument from the database.
func (q *Querier) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error) {
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
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeStepInstrument")
	}

	return recipeStepInstrument, nil
}

// GetRecipeStepInstruments fetches a list of recipe step instruments from the database that meet a particular filter.
func (q *Querier) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.RecipeStepInstrumentList, err error) {
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
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "recipe_step_instruments", getRecipeStepInstrumentsJoins, []string{"valid_instruments.id", "recipe_step_instruments.id"}, nil, householdOwnershipColumn, recipeStepInstrumentsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "recipeStepInstruments", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step instruments list retrieval query")
	}

	if x.RecipeStepInstruments, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepInstruments(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step instruments")
	}

	return x, nil
}

//go:embed queries/recipe_step_instruments/recipe_step_instruments_get_for_recipe.sql
var getRecipeStepInstrumentsForRecipeQuery string

// getRecipeStepInstrumentsForRecipe fetches a list of recipe step instruments from the database that meet a particular filter.
func (q *Querier) getRecipeStepInstrumentsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepInstrument, error) {
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
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step instruments list retrieval query")
	}

	recipeStepInstruments, _, _, err := q.scanRecipeStepInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step instruments")
	}

	return recipeStepInstruments, nil
}

//go:embed queries/recipe_step_instruments/recipe_step_instruments_create.sql
var recipeStepInstrumentCreationQuery string

// CreateRecipeStepInstrument creates a recipe step instrument in the database.
func (q *Querier) createRecipeStepInstrument(ctx context.Context, querier database.SQLQueryExecutor, input *types.RecipeStepInstrumentDatabaseCreationInput) (*types.RecipeStepInstrument, error) {
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
		input.Optional,
		input.MinimumQuantity,
		input.MaximumQuantity,
		input.BelongsToRecipeStep,
	}

	// create the recipe step instrument.
	if err := q.performWriteQuery(ctx, querier, "recipe step instrument creation", recipeStepInstrumentCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe step instrument creation query")
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
		Optional:            input.Optional,
		MinimumQuantity:     input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
		CreatedAt:           q.currentTime(),
	}

	if input.InstrumentID != nil {
		x.Instrument = &types.ValidInstrument{ID: *input.InstrumentID}
	}

	tracing.AttachRecipeStepInstrumentIDToSpan(span, x.ID)
	logger.Info("recipe step instrument created")

	return x, nil
}

// CreateRecipeStepInstrument creates a recipe step instrument in the database.
func (q *Querier) CreateRecipeStepInstrument(ctx context.Context, input *types.RecipeStepInstrumentDatabaseCreationInput) (*types.RecipeStepInstrument, error) {
	return q.createRecipeStepInstrument(ctx, q.db, input)
}

//go:embed queries/recipe_step_instruments/recipe_step_instruments_update.sql
var updateRecipeStepInstrumentQuery string

// UpdateRecipeStepInstrument updates a particular recipe step instrument.
func (q *Querier) UpdateRecipeStepInstrument(ctx context.Context, updated *types.RecipeStepInstrument) error {
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
		updated.Optional,
		updated.MinimumQuantity,
		updated.MaximumQuantity,
		updated.BelongsToRecipeStep,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step instrument update", updateRecipeStepInstrumentQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step instrument")
	}

	logger.Info("recipe step instrument updated")

	return nil
}

//go:embed queries/recipe_step_instruments/recipe_step_instruments_archive.sql
var archiveRecipeStepInstrumentQuery string

// ArchiveRecipeStepInstrument archives a recipe step instrument from the database by its ID.
func (q *Querier) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error {
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
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step instrument")
	}

	logger.Info("recipe step instrument archived")

	return nil
}
