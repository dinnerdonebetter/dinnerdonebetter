package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	recipeStepsOnRecipeStepVesselsJoinClause      = "recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id"
	recipeStepVesselsOnValidInstrumentsJoinClause = "valid_instruments ON recipe_step_vessels.valid_instrument_id=valid_instruments.id"
)

var (
	_ types.RecipeStepVesselDataManager = (*Querier)(nil)

	// recipeStepVesselsTableColumns are the columns for the recipe_step_vessels table.
	recipeStepVesselsTableColumns = []string{
		"recipe_step_vessels.id",
		"valid_instruments.id",
		"valid_instruments.name",
		"valid_instruments.plural_name",
		"valid_instruments.description",
		"valid_instruments.icon_path",
		"valid_instruments.usable_for_storage",
		"valid_instruments.display_in_summary_lists",
		"valid_instruments.include_in_generated_instructions",
		"valid_instruments.slug",
		"valid_instruments.created_at",
		"valid_instruments.last_updated_at",
		"valid_instruments.archived_at",
		"recipe_step_vessels.name",
		"recipe_step_vessels.notes",
		"recipe_step_vessels.belongs_to_recipe_step",
		"recipe_step_vessels.recipe_step_product_id",
		"recipe_step_vessels.vessel_predicate",
		"recipe_step_vessels.minimum_quantity",
		"recipe_step_vessels.maximum_quantity",
		"recipe_step_vessels.unavailable_after_step",
		"recipe_step_vessels.created_at",
		"recipe_step_vessels.last_updated_at",
		"recipe_step_vessels.archived_at",
	}

	getRecipeStepVesselsJoins = []string{
		recipeStepsOnRecipeStepVesselsJoinClause,
		recipeStepVesselsOnValidInstrumentsJoinClause,
		recipesOnRecipeStepsJoinClause,
	}
)

// scanRecipeStepVessel takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step vessel struct.
func (q *Querier) scanRecipeStepVessel(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepVessel, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStepVessel{}
	instrument := &types.NullableValidInstrument{}

	targetVars := []any{
		&x.ID,
		&instrument.ID,
		&instrument.Name,
		&instrument.PluralName,
		&instrument.Description,
		&instrument.IconPath,
		&instrument.UsableForStorage,
		&instrument.DisplayInSummaryLists,
		&instrument.IncludeInGeneratedInstructions,
		&instrument.Slug,
		&instrument.CreatedAt,
		&instrument.LastUpdatedAt,
		&instrument.ArchivedAt,
		&x.Name,
		&x.Notes,
		&x.BelongsToRecipeStep,
		&x.RecipeStepProductID,
		&x.VesselPreposition,
		&x.MinimumQuantity,
		&x.MaximumQuantity,
		&x.UnavailableAfterStep,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	if instrument.ID != nil {
		x.Instrument = converters.ConvertNullableValidInstrumentToValidInstrument(instrument)
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepVessels takes some database rows and turns them into a slice of recipe step vessels.
func (q *Querier) scanRecipeStepVessels(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepVessels []*types.RecipeStepVessel, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepVessel(ctx, rows, includeCounts)
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

		recipeStepVessels = append(recipeStepVessels, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return recipeStepVessels, filteredCount, totalCount, nil
}

//go:embed queries/recipe_step_vessels/exists.sql
var recipeStepVesselExistenceQuery string

// RecipeStepVesselExists fetches whether a recipe step vessel exists from the database.
func (q *Querier) RecipeStepVesselExists(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (exists bool, err error) {
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

	if recipeStepVesselID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepVesselID)
	tracing.AttachRecipeStepVesselIDToSpan(span, recipeStepVesselID)

	args := []any{
		recipeStepID,
		recipeStepVesselID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepVesselExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step vessel existence check")
	}

	return result, nil
}

//go:embed queries/recipe_step_vessels/get_one.sql
var getRecipeStepVesselQuery string

// GetRecipeStepVessel fetches a recipe step vessel from the database.
func (q *Querier) GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error) {
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

	if recipeStepVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepVesselID)
	tracing.AttachRecipeStepVesselIDToSpan(span, recipeStepVesselID)

	args := []any{
		recipeStepID,
		recipeStepVesselID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "recipe step vessel", getRecipeStepVesselQuery, args)

	recipeStepVessel, _, _, err := q.scanRecipeStepVessel(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipeStepVessel")
	}

	return recipeStepVessel, nil
}

// GetRecipeStepVessels fetches a list of recipe step vessels from the database that meet a particular filter.
func (q *Querier) GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepVessel], err error) {
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

	x = &types.QueryFilteredResult[types.RecipeStepVessel]{}
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

	query, args := q.buildListQuery(ctx, "recipe_step_vessels", getRecipeStepVesselsJoins, []string{"valid_instruments.id"}, nil, householdOwnershipColumn, recipeStepVesselsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "recipe step vessels", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step vessels list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepVessels(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step vessels")
	}

	return x, nil
}

//go:embed queries/recipe_step_vessels/get_for_recipe.sql
var getRecipeStepVesselsForRecipeQuery string

// getRecipeStepVesselsForRecipe fetches a list of recipe step vessels from the database that meet a particular filter.
func (q *Querier) getRecipeStepVesselsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	args := []any{
		recipeID,
	}

	rows, err := q.getRows(ctx, q.db, "recipe step vessels for recipe", getRecipeStepVesselsForRecipeQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step vessels list retrieval query")
	}

	recipeStepVessels, _, _, err := q.scanRecipeStepVessels(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning recipe step vessels")
	}

	return recipeStepVessels, nil
}

//go:embed queries/recipe_step_vessels/create.sql
var recipeStepVesselCreationQuery string

// CreateRecipeStepVessel creates a recipe step vessel in the database.
func (q *Querier) createRecipeStepVessel(ctx context.Context, querier database.SQLQueryExecutor, input *types.RecipeStepVesselDatabaseCreationInput) (*types.RecipeStepVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepVesselIDKey, input.ID)

	args := []any{
		input.ID,
		input.Name,
		input.Notes,
		input.BelongsToRecipeStep,
		input.RecipeStepProductID,
		input.InstrumentID,
		input.VesselPreposition,
		input.MinimumQuantity,
		input.MaximumQuantity,
		input.UnavailableAfterStep,
	}

	// create the recipe step vessel.
	if err := q.performWriteQuery(ctx, querier, "recipe step vessel creation", recipeStepVesselCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe step vessel creation query")
	}

	x := &types.RecipeStepVessel{
		ID:                   input.ID,
		RecipeStepProductID:  input.RecipeStepProductID,
		Name:                 input.Name,
		Notes:                input.Notes,
		BelongsToRecipeStep:  input.BelongsToRecipeStep,
		MinimumQuantity:      input.MinimumQuantity,
		MaximumQuantity:      input.MaximumQuantity,
		VesselPreposition:    input.VesselPreposition,
		UnavailableAfterStep: input.UnavailableAfterStep,
		CreatedAt:            q.currentTime(),
	}

	if input.InstrumentID != nil {
		x.Instrument = &types.ValidInstrument{ID: *input.InstrumentID}
	}

	tracing.AttachRecipeStepVesselIDToSpan(span, x.ID)
	logger.Info("recipe step vessel created")

	return x, nil
}

// CreateRecipeStepVessel creates a recipe step vessel in the database.
func (q *Querier) CreateRecipeStepVessel(ctx context.Context, input *types.RecipeStepVesselDatabaseCreationInput) (*types.RecipeStepVessel, error) {
	return q.createRecipeStepVessel(ctx, q.db, input)
}

//go:embed queries/recipe_step_vessels/update.sql
var updateRecipeStepVesselQuery string

// UpdateRecipeStepVessel updates a particular recipe step vessel.
func (q *Querier) UpdateRecipeStepVessel(ctx context.Context, updated *types.RecipeStepVessel) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepVesselIDKey, updated.ID)
	tracing.AttachRecipeStepVesselIDToSpan(span, updated.ID)

	var instrumentID *string
	if updated.Instrument != nil {
		instrumentID = &updated.Instrument.ID
	}

	args := []any{
		updated.Name,
		updated.Notes,
		updated.BelongsToRecipeStep,
		updated.RecipeStepProductID,
		instrumentID,
		updated.VesselPreposition,
		updated.MinimumQuantity,
		updated.MaximumQuantity,
		updated.UnavailableAfterStep,
		updated.BelongsToRecipeStep,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step vessel update", updateRecipeStepVesselQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step vessel")
	}

	logger.Info("recipe step vessel updated")

	return nil
}

//go:embed queries/recipe_step_vessels/archive.sql
var archiveRecipeStepVesselQuery string

// ArchiveRecipeStepVessel archives a recipe step vessel from the database by its ID.
func (q *Querier) ArchiveRecipeStepVessel(ctx context.Context, recipeStepID, recipeStepVesselID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepVesselID)
	tracing.AttachRecipeStepVesselIDToSpan(span, recipeStepVesselID)

	args := []any{
		recipeStepID,
		recipeStepVesselID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step vessel archive", archiveRecipeStepVesselQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step vessel")
	}

	logger.Info("recipe step vessel archived")

	return nil
}
