package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/Masterminds/squirrel"
)

const (
	validInstrumentsOnValidPreparationVesselsJoinClause  = "valid_vessels ON valid_preparation_vessels.valid_vessel_id = valid_vessels.id"
	validPreparationsOnValidPreparationVesselsJoinClause = "valid_preparations ON valid_preparation_vessels.valid_preparation_id = valid_preparations.id"
	validMeasurementUnitsOnValidVesselsJoinClause        = "valid_measurement_units ON valid_vessels.capacity_unit = valid_measurement_units.id"
)

var (
	_ types.ValidPreparationVesselDataManager = (*Querier)(nil)

	// fullValidPreparationVesselsTableColumns are the columns for the valid_preparation_vessels table.
	fullValidPreparationVesselsTableColumns = []string{
		"valid_preparation_vessels.id",
		"valid_preparation_vessels.notes",
		"valid_preparations.id",
		"valid_preparations.name",
		"valid_preparations.description",
		"valid_preparations.icon_path",
		"valid_preparations.yields_nothing",
		"valid_preparations.restrict_to_ingredients",
		"valid_preparations.minimum_ingredient_count",
		"valid_preparations.maximum_ingredient_count",
		"valid_preparations.minimum_instrument_count",
		"valid_preparations.maximum_instrument_count",
		"valid_preparations.temperature_required",
		"valid_preparations.time_estimate_required",
		"valid_preparations.condition_expression_required",
		"valid_preparations.consumes_vessel",
		"valid_preparations.only_for_vessels",
		"valid_preparations.minimum_vessel_count",
		"valid_preparations.maximum_vessel_count",
		"valid_preparations.slug",
		"valid_preparations.past_tense",
		"valid_preparations.created_at",
		"valid_preparations.last_updated_at",
		"valid_preparations.archived_at",
		// BEGIN valid_vessel embed
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
		// END valid_vessel embed
		"valid_preparation_vessels.created_at",
		"valid_preparation_vessels.last_updated_at",
		"valid_preparation_vessels.archived_at",
	}
)

// scanValidPreparationVessel takes a database Scanner (i.e. *sql.Row) and scans the result into a valid preparation vessel struct.
func (q *Querier) scanValidPreparationVessel(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidPreparationVessel, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidPreparationVessel{
		Vessel: types.ValidVessel{
			CapacityUnit: &types.ValidMeasurementUnit{},
		},
	}

	targetVars := []any{
		&x.ID,
		&x.Notes,
		&x.Preparation.ID,
		&x.Preparation.Name,
		&x.Preparation.Description,
		&x.Preparation.IconPath,
		&x.Preparation.YieldsNothing,
		&x.Preparation.RestrictToIngredients,
		&x.Preparation.MinimumIngredientCount,
		&x.Preparation.MaximumIngredientCount,
		&x.Preparation.MinimumInstrumentCount,
		&x.Preparation.MaximumInstrumentCount,
		&x.Preparation.TemperatureRequired,
		&x.Preparation.TimeEstimateRequired,
		&x.Preparation.ConditionExpressionRequired,
		&x.Preparation.ConsumesVessel,
		&x.Preparation.OnlyForVessels,
		&x.Preparation.MinimumVesselCount,
		&x.Preparation.MaximumVesselCount,
		&x.Preparation.Slug,
		&x.Preparation.PastTense,
		&x.Preparation.CreatedAt,
		&x.Preparation.LastUpdatedAt,
		&x.Preparation.ArchivedAt,
		&x.Vessel.ID,
		&x.Vessel.Name,
		&x.Vessel.PluralName,
		&x.Vessel.Description,
		&x.Vessel.IconPath,
		&x.Vessel.UsableForStorage,
		&x.Vessel.Slug,
		&x.Vessel.DisplayInSummaryLists,
		&x.Vessel.IncludeInGeneratedInstructions,
		&x.Vessel.Capacity,
		&x.Vessel.CapacityUnit.ID,
		&x.Vessel.CapacityUnit.Name,
		&x.Vessel.CapacityUnit.Description,
		&x.Vessel.CapacityUnit.Volumetric,
		&x.Vessel.CapacityUnit.IconPath,
		&x.Vessel.CapacityUnit.Universal,
		&x.Vessel.CapacityUnit.Metric,
		&x.Vessel.CapacityUnit.Imperial,
		&x.Vessel.CapacityUnit.Slug,
		&x.Vessel.CapacityUnit.PluralName,
		&x.Vessel.CapacityUnit.CreatedAt,
		&x.Vessel.CapacityUnit.LastUpdatedAt,
		&x.Vessel.CapacityUnit.ArchivedAt,
		&x.Vessel.WidthInMillimeters,
		&x.Vessel.LengthInMillimeters,
		&x.Vessel.HeightInMillimeters,
		&x.Vessel.Shape,
		&x.Vessel.CreatedAt,
		&x.Vessel.LastUpdatedAt,
		&x.Vessel.ArchivedAt,
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

	if x.Vessel.CapacityUnit == nil || x.Vessel.CapacityUnit.ID == "" {
		x.Vessel.CapacityUnit = nil
	}

	return x, filteredCount, totalCount, nil
}

// scanValidPreparationVessels takes some database rows and turns them into a slice of valid preparation vessels.
func (q *Querier) scanValidPreparationVessels(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validPreparationVessels []*types.ValidPreparationVessel, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidPreparationVessel(ctx, rows, includeCounts)
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

		validPreparationVessels = append(validPreparationVessels, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validPreparationVessels, filteredCount, totalCount, nil
}

//go:embed queries/valid_preparation_vessels/exists.sql
var validPreparationVesselExistenceQuery string

// ValidPreparationVesselExists fetches whether a valid preparation vessel exists from the database.
func (q *Querier) ValidPreparationVesselExists(ctx context.Context, validPreparationVesselID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationVesselID == "" {
		return false, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationVesselIDToSpan(span, validPreparationVesselID)

	args := []any{
		validPreparationVesselID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validPreparationVesselExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, span, "performing valid preparation vessel existence check")
	}

	return result, nil
}

//go:embed queries/valid_preparation_vessels/get_one.sql
var getValidPreparationVesselQuery string

// GetValidPreparationVessel fetches a valid preparation vessel from the database.
func (q *Querier) GetValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationVesselIDToSpan(span, validPreparationVesselID)

	args := []any{
		validPreparationVesselID,
	}

	row := q.getOneRow(ctx, q.db, "validPreparationVessel", getValidPreparationVesselQuery, args)

	validPreparationVessel, _, _, err := q.scanValidPreparationVessel(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning validPreparationVessel")
	}

	return validPreparationVessel, nil
}

// GetValidPreparationVessels fetches a list of valid preparation vessels from the database that meet a particular filter.
func (q *Querier) GetValidPreparationVessels(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.QueryFilteredResult[types.ValidPreparationVessel]{}
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	joins := []string{
		validInstrumentsOnValidPreparationVesselsJoinClause,
		validPreparationsOnValidPreparationVesselsJoinClause,
		validMeasurementUnitsOnValidVesselsJoinClause,
	}
	groupBys := []string{
		"valid_preparations.id",
		"valid_vessels.id",
		"valid_preparation_vessels.id",
		"valid_measurement_units.id",
	}
	query, args := q.buildListQuery(ctx, "valid_preparation_vessels", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationVesselsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "valid preparation instruments", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid preparation vessels list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationVessels(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid preparation vessels")
	}

	return x, nil
}

func (q *Querier) buildGetValidPreparationVesselsRestrictedByIDsQuery(ctx context.Context, column string, limit uint8, ids []string) (query string, args []any) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query, args, err := q.sqlBuilder.Select(fullValidPreparationVesselsTableColumns...).
		From("valid_preparation_vessels").
		Join(validInstrumentsOnValidPreparationVesselsJoinClause).
		Join(validPreparationsOnValidPreparationVesselsJoinClause).
		LeftJoin(validMeasurementUnitsOnValidVesselsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("valid_preparation_vessels.%s", column): ids,
			"valid_preparation_vessels.archived_at":             nil,
		}).
		Limit(uint64(limit)).
		ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

func (q *Querier) buildGetValidPreparationVesselsWithPreparationIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []any) {
	return q.buildGetValidPreparationVesselsRestrictedByIDsQuery(ctx, "valid_preparation_id", limit, ids)
}

// GetValidPreparationVesselsForPreparation fetches a list of valid preparation vessels from the database that meet a particular filter.
func (q *Querier) GetValidPreparationVesselsForPreparation(ctx context.Context, preparationID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationVesselIDToSpan(span, preparationID)

	x = &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Pagination: types.Pagination{
			Limit: 20,
		},
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	// the use of filter here is so weird, since we only respect the limit, but I'm trying to get this done, okay?
	query, args := q.buildGetValidPreparationVesselsWithPreparationIDsQuery(ctx, x.Limit, []string{preparationID})

	rows, err := q.getRows(ctx, q.db, "valid preparation instruments for preparation", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid preparation vessels list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationVessels(ctx, rows, false); err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid preparation vessels")
	}

	return x, nil
}

func (q *Querier) buildGetValidPreparationVesselsWithVesselIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []any) {
	return q.buildGetValidPreparationVesselsRestrictedByIDsQuery(ctx, "valid_vessel_id", limit, ids)
}

// GetValidPreparationVesselsForVessel fetches a list of valid preparation vessels from the database that meet a particular filter.
func (q *Querier) GetValidPreparationVesselsForVessel(ctx context.Context, vesselID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if vesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidVesselIDToSpan(span, vesselID)

	x = &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Pagination: types.Pagination{
			Limit: 20,
		},
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	// the use of filter here is so weird, since we only respect the limit, but I'm trying to get this done, okay?
	query, args := q.buildGetValidPreparationVesselsWithVesselIDsQuery(ctx, x.Limit, []string{vesselID})

	q.logger.WithValue("query", query).Info("querying")

	rows, err := q.getRows(ctx, q.db, "valid preparation instruments for instrument", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid preparation vessels list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationVessels(ctx, rows, false); err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid preparation vessels")
	}

	return x, nil
}

//go:embed queries/valid_preparation_vessels/create.sql
var validPreparationVesselCreationQuery string

// CreateValidPreparationVessel creates a valid preparation vessel in the database.
func (q *Querier) CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselDatabaseCreationInput) (*types.ValidPreparationVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationVesselIDKey, input.ID)

	args := []any{
		input.ID,
		input.Notes,
		input.ValidPreparationID,
		input.ValidVesselID,
	}

	// create the valid preparation vessel.
	if err := q.performWriteQuery(ctx, q.db, "valid preparation vessel creation", validPreparationVesselCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid preparation vessel creation query")
	}

	x := &types.ValidPreparationVessel{
		ID:          input.ID,
		Notes:       input.Notes,
		Preparation: types.ValidPreparation{ID: input.ValidPreparationID},
		Vessel:      types.ValidVessel{ID: input.ValidVesselID},
		CreatedAt:   q.currentTime(),
	}

	tracing.AttachValidPreparationVesselIDToSpan(span, x.ID)
	logger.Info("valid preparation vessel created")

	return x, nil
}

//go:embed queries/valid_preparation_vessels/update.sql
var updateValidPreparationVesselQuery string

// UpdateValidPreparationVessel updates a particular valid preparation vessel.
func (q *Querier) UpdateValidPreparationVessel(ctx context.Context, updated *types.ValidPreparationVessel) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationVesselIDKey, updated.ID)
	tracing.AttachValidPreparationVesselIDToSpan(span, updated.ID)

	args := []any{
		updated.Notes,
		updated.Preparation.ID,
		updated.Vessel.ID,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation vessel update", updateValidPreparationVesselQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation vessel")
	}

	logger.Info("valid preparation vessel updated")

	return nil
}

// ArchiveValidPreparationVessel archives a valid preparation vessel from the database by its ID.
func (q *Querier) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachValidPreparationVesselIDToSpan(span, validPreparationVesselID)

	if err := q.generatedQuerier.ArchiveValidPreparationVessel(ctx, q.db, validPreparationVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation vessel")
	}

	logger.Info("valid preparation vessel archived")

	return nil
}
