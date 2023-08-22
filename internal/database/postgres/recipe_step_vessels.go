package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	recipeStepsOnRecipeStepVesselsJoinClause      = "recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id"
	recipeStepVesselsOnValidInstrumentsJoinClause = "valid_vessels ON recipe_step_vessels.valid_vessel_id=valid_vessels.id"
)

var (
	_ types.RecipeStepVesselDataManager = (*Querier)(nil)

	// recipeStepVesselsTableColumns are the columns for the recipe_step_vessels table.
	recipeStepVesselsTableColumns = []string{
		"recipe_step_vessels.id",
		"valid_vessels.id",
		"valid_vessels.name",
		"valid_vessels.plural_name",
		"valid_vessels.description",
		"valid_vessels.icon_path",
		"valid_vessels.usable_for_storage",
		"valid_vessels.slug",
		"valid_vessels.display_in_summary_lists",
		"valid_vessels.include_in_generated_instructions",
		"valid_vessels.capacity",
		"valid_measurement_units.id",
		"valid_measurement_units.name",
		"valid_measurement_units.description",
		"valid_measurement_units.volumetric",
		"valid_measurement_units.icon_path",
		"valid_measurement_units.universal",
		"valid_measurement_units.metric",
		"valid_measurement_units.imperial",
		"valid_measurement_units.slug",
		"valid_measurement_units.plural_name",
		"valid_measurement_units.created_at",
		"valid_measurement_units.last_updated_at",
		"valid_measurement_units.archived_at",
		"valid_vessels.width_in_millimeters",
		"valid_vessels.length_in_millimeters",
		"valid_vessels.height_in_millimeters",
		"valid_vessels.shape",
		"valid_vessels.created_at",
		"valid_vessels.last_updated_at",
		"valid_vessels.archived_at",
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
		"valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id",
	}
)

// scanRecipeStepVessel takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step vessel struct.
func (q *Querier) scanRecipeStepVessel(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepVessel, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.RecipeStepVessel{}
	vessel := &types.NullableValidVessel{
		CapacityUnit: &types.NullableValidMeasurementUnit{},
	}

	targetVars := []any{
		&x.ID,
		&vessel.ID,
		&vessel.Name,
		&vessel.PluralName,
		&vessel.Description,
		&vessel.IconPath,
		&vessel.UsableForStorage,
		&vessel.Slug,
		&vessel.DisplayInSummaryLists,
		&vessel.IncludeInGeneratedInstructions,
		&vessel.Capacity,
		&vessel.CapacityUnit.ID,
		&vessel.CapacityUnit.Name,
		&vessel.CapacityUnit.Description,
		&vessel.CapacityUnit.Volumetric,
		&vessel.CapacityUnit.IconPath,
		&vessel.CapacityUnit.Universal,
		&vessel.CapacityUnit.Metric,
		&vessel.CapacityUnit.Imperial,
		&vessel.CapacityUnit.Slug,
		&vessel.CapacityUnit.PluralName,
		&vessel.CapacityUnit.CreatedAt,
		&vessel.CapacityUnit.LastUpdatedAt,
		&vessel.CapacityUnit.ArchivedAt,
		&vessel.WidthInMillimeters,
		&vessel.LengthInMillimeters,
		&vessel.HeightInMillimeters,
		&vessel.Shape,
		&vessel.CreatedAt,
		&vessel.LastUpdatedAt,
		&vessel.ArchivedAt,
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

	if vessel.ID != nil {
		x.Vessel = converters.ConvertNullableValidVesselToValidVessel(vessel)
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

	result, err := q.generatedQuerier.CheckRecipeStepVesselExistence(ctx, q.db, &generated.CheckRecipeStepVesselExistenceParams{
		RecipeStepID:       recipeStepID,
		RecipeStepVesselID: recipeStepVesselID,
		RecipeID:           recipeID,
	})
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

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.RecipeStepVessel]{
		Pagination: filter.ToPagination(),
	}

	query, args := q.buildListQuery(ctx, "recipe_step_vessels", getRecipeStepVesselsJoins, []string{"valid_vessels.id", "valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepVesselsTableColumns, "", false, filter)

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

// CreateRecipeStepVessel creates a recipe step vessel in the database.
func (q *Querier) createRecipeStepVessel(ctx context.Context, querier database.SQLQueryExecutor, input *types.RecipeStepVesselDatabaseCreationInput) (*types.RecipeStepVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepVesselIDKey, input.ID)

	// create the recipe step vessel.
	if err := q.generatedQuerier.CreateRecipeStepVessel(ctx, querier, &generated.CreateRecipeStepVesselParams{
		ID:                   input.ID,
		Name:                 input.Name,
		Notes:                input.Notes,
		BelongsToRecipeStep:  input.BelongsToRecipeStep,
		VesselPredicate:      input.VesselPreposition,
		RecipeStepProductID:  nullStringFromStringPointer(input.RecipeStepProductID),
		ValidVesselID:        nullStringFromStringPointer(input.VesselID),
		MaximumQuantity:      nullInt32FromUint32Pointer(input.MaximumQuantity),
		MinimumQuantity:      int32(input.MinimumQuantity),
		UnavailableAfterStep: input.UnavailableAfterStep,
	}); err != nil {
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

	if input.VesselID != nil {
		x.Vessel = &types.ValidVessel{ID: *input.VesselID}
	}

	tracing.AttachRecipeStepVesselIDToSpan(span, x.ID)
	logger.Info("recipe step vessel created")

	return x, nil
}

// CreateRecipeStepVessel creates a recipe step vessel in the database.
func (q *Querier) CreateRecipeStepVessel(ctx context.Context, input *types.RecipeStepVesselDatabaseCreationInput) (*types.RecipeStepVessel, error) {
	return q.createRecipeStepVessel(ctx, q.db, input)
}

// UpdateRecipeStepVessel updates a particular recipe step vessel.
func (q *Querier) UpdateRecipeStepVessel(ctx context.Context, updated *types.RecipeStepVessel) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepVesselIDKey, updated.ID)
	tracing.AttachRecipeStepVesselIDToSpan(span, updated.ID)

	var vesselID *string
	if updated.Vessel != nil {
		vesselID = &updated.Vessel.ID
	}

	if err := q.generatedQuerier.UpdateRecipeStepVessel(ctx, q.db, &generated.UpdateRecipeStepVesselParams{
		Name:                 updated.Name,
		Notes:                updated.Notes,
		RecipeStepID:         updated.BelongsToRecipeStep,
		VesselPredicate:      updated.VesselPreposition,
		ID:                   updated.ID,
		RecipeStepProductID:  nullStringFromStringPointer(updated.RecipeStepProductID),
		ValidVesselID:        nullStringFromStringPointer(vesselID),
		MaximumQuantity:      nullInt32FromUint32Pointer(updated.MaximumQuantity),
		MinimumQuantity:      int32(updated.MinimumQuantity),
		UnavailableAfterStep: updated.UnavailableAfterStep,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step vessel")
	}

	logger.Info("recipe step vessel updated")

	return nil
}

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

	if err := q.generatedQuerier.ArchiveRecipeStepVessel(ctx, q.db, &generated.ArchiveRecipeStepVesselParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepVesselID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step vessel")
	}

	logger.Info("recipe step vessel archived")

	return nil
}
