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
	validInstrumentsOnValidPreparationInstrumentsJoinClause  = "valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id"
	validPreparationsOnValidPreparationInstrumentsJoinClause = "valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id"
)

var (
	_ types.ValidPreparationInstrumentDataManager = (*Querier)(nil)

	// fullValidPreparationInstrumentsTableColumns are the columns for the valid_preparation_instruments table.
	fullValidPreparationInstrumentsTableColumns = []string{
		"valid_preparation_instruments.id",
		"valid_preparation_instruments.notes",
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
		"valid_preparation_instruments.created_at",
		"valid_preparation_instruments.last_updated_at",
		"valid_preparation_instruments.archived_at",
	}
)

// scanValidPreparationInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a valid preparation instrument struct.
func (q *Querier) scanValidPreparationInstrument(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidPreparationInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidPreparationInstrument{}

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
		&x.Instrument.ID,
		&x.Instrument.Name,
		&x.Instrument.PluralName,
		&x.Instrument.Description,
		&x.Instrument.IconPath,
		&x.Instrument.UsableForStorage,
		&x.Instrument.DisplayInSummaryLists,
		&x.Instrument.IncludeInGeneratedInstructions,
		&x.Instrument.Slug,
		&x.Instrument.CreatedAt,
		&x.Instrument.LastUpdatedAt,
		&x.Instrument.ArchivedAt,
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

	return x, filteredCount, totalCount, nil
}

// scanValidPreparationInstruments takes some database rows and turns them into a slice of valid preparation instruments.
func (q *Querier) scanValidPreparationInstruments(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validPreparationInstruments []*types.ValidPreparationInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidPreparationInstrument(ctx, rows, includeCounts)
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

		validPreparationInstruments = append(validPreparationInstruments, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validPreparationInstruments, filteredCount, totalCount, nil
}

//go:embed queries/valid_preparation_instruments/exists.sql
var validPreparationInstrumentExistenceQuery string

// ValidPreparationInstrumentExists fetches whether a valid preparation instrument exists from the database.
func (q *Querier) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	args := []any{
		validPreparationInstrumentID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validPreparationInstrumentExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, span, "performing valid preparation instrument existence check")
	}

	return result, nil
}

//go:embed queries/valid_preparation_instruments/get_one.sql
var getValidPreparationInstrumentQuery string

// GetValidPreparationInstrument fetches a valid preparation instrument from the database.
func (q *Querier) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	args := []any{
		validPreparationInstrumentID,
	}

	row := q.getOneRow(ctx, q.db, "validPreparationInstrument", getValidPreparationInstrumentQuery, args)

	validPreparationInstrument, _, _, err := q.scanValidPreparationInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning validPreparationInstrument")
	}

	return validPreparationInstrument, nil
}

// GetValidPreparationInstruments fetches a list of valid preparation instruments from the database that meet a particular filter.
func (q *Querier) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationInstrument], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparationInstrument]{
		Pagination: filter.ToPagination(),
	}

	joins := []string{
		validInstrumentsOnValidPreparationInstrumentsJoinClause,
		validPreparationsOnValidPreparationInstrumentsJoinClause,
	}

	groupBys := []string{
		"valid_preparations.id",
		"valid_instruments.id",
		"valid_preparation_instruments.id",
	}

	query, args := q.buildListQuery(ctx, "valid_preparation_instruments", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationInstrumentsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "valid preparation instruments", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparation instruments list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationInstruments(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid preparation instruments")
	}

	return x, nil
}

func (q *Querier) buildGetValidPreparationInstrumentsRestrictedByIDsQuery(ctx context.Context, column string, limit uint8, ids []string) (query string, args []any) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query, args, err := q.sqlBuilder.Select(fullValidPreparationInstrumentsTableColumns...).
		From("valid_preparation_instruments").
		Join(validInstrumentsOnValidPreparationInstrumentsJoinClause).
		Join(validPreparationsOnValidPreparationInstrumentsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("valid_preparation_instruments.%s", column): ids,
			"valid_preparation_instruments.archived_at":             nil,
		}).
		Limit(uint64(limit)).
		ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

func (q *Querier) buildGetValidPreparationInstrumentsWithPreparationIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []any) {
	return q.buildGetValidPreparationInstrumentsRestrictedByIDsQuery(ctx, "valid_preparation_id", limit, ids)
}

// GetValidPreparationInstrumentsForPreparation fetches a list of valid preparation instruments from the database that meet a particular filter.
func (q *Querier) GetValidPreparationInstrumentsForPreparation(ctx context.Context, preparationID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationInstrument], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationInstrumentIDToSpan(span, preparationID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparationInstrument]{
		Pagination: filter.ToPagination(),
	}

	// the use of filter here is so weird, since we only respect the limit, but I'm trying to get this done, okay?
	query, args := q.buildGetValidPreparationInstrumentsWithPreparationIDsQuery(ctx, x.Limit, []string{preparationID})

	rows, err := q.getRows(ctx, q.db, "valid preparation instruments for preparation", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparation instruments list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationInstruments(ctx, rows, false); err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid preparation instruments")
	}

	return x, nil
}

func (q *Querier) buildGetValidPreparationInstrumentsWithInstrumentIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []any) {
	return q.buildGetValidPreparationInstrumentsRestrictedByIDsQuery(ctx, "valid_instrument_id", limit, ids)
}

// GetValidPreparationInstrumentsForInstrument fetches a list of valid preparation instruments from the database that meet a particular filter.
func (q *Querier) GetValidPreparationInstrumentsForInstrument(ctx context.Context, instrumentID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationInstrument], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if instrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationInstrumentIDToSpan(span, instrumentID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparationInstrument]{
		Pagination: filter.ToPagination(),
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	// the use of filter here is so weird, since we only respect the limit, but I'm trying to get this done, okay?
	query, args := q.buildGetValidPreparationInstrumentsWithInstrumentIDsQuery(ctx, x.Limit, []string{instrumentID})

	rows, err := q.getRows(ctx, q.db, "valid preparation instruments for instrument", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparation instruments list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationInstruments(ctx, rows, false); err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid preparation instruments")
	}

	return x, nil
}

//go:embed queries/valid_preparation_instruments/create.sql
var validPreparationInstrumentCreationQuery string

// CreateValidPreparationInstrument creates a valid preparation instrument in the database.
func (q *Querier) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentDatabaseCreationInput) (*types.ValidPreparationInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, input.ID)

	args := []any{
		input.ID,
		input.Notes,
		input.ValidPreparationID,
		input.ValidInstrumentID,
	}

	// create the valid preparation instrument.
	if err := q.performWriteQuery(ctx, q.db, "valid preparation instrument creation", validPreparationInstrumentCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid preparation instrument creation query")
	}

	x := &types.ValidPreparationInstrument{
		ID:          input.ID,
		Notes:       input.Notes,
		Preparation: types.ValidPreparation{ID: input.ValidPreparationID},
		Instrument:  types.ValidInstrument{ID: input.ValidInstrumentID},
		CreatedAt:   q.currentTime(),
	}

	tracing.AttachValidPreparationInstrumentIDToSpan(span, x.ID)
	logger.Info("valid preparation instrument created")

	return x, nil
}

//go:embed queries/valid_preparation_instruments/update.sql
var updateValidPreparationInstrumentQuery string

// UpdateValidPreparationInstrument updates a particular valid preparation instrument.
func (q *Querier) UpdateValidPreparationInstrument(ctx context.Context, updated *types.ValidPreparationInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, updated.ID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, updated.ID)

	args := []any{
		updated.Notes,
		updated.Preparation.ID,
		updated.Instrument.ID,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation instrument update", updateValidPreparationInstrumentQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation instrument")
	}

	logger.Info("valid preparation instrument updated")

	return nil
}

// ArchiveValidPreparationInstrument archives a valid preparation instrument from the database by its ID.
func (q *Querier) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	if err := q.generatedQuerier.ArchiveValidPreparationInstrument(ctx, q.db, validPreparationInstrumentID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation instrument")
	}

	logger.Info("valid preparation instrument archived")

	return nil
}
