package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/Masterminds/squirrel"
)

const (
	validPreparationsTable = "valid_preparations"

	validPreparationsOnRecipeStepsJoinClause = "valid_preparations ON recipe_steps.preparation_id=valid_preparations.id"
)

var (
	_ types.ValidPreparationDataManager = (*Querier)(nil)

	// validPreparationsTableColumns are the columns for the valid_preparations table.
	validPreparationsTableColumns = []string{
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
	}
)

// scanValidPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a valid preparation struct.
func (q *Querier) scanValidPreparation(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidPreparation{}

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Description,
		&x.IconPath,
		&x.YieldsNothing,
		&x.RestrictToIngredients,
		&x.MinimumIngredientCount,
		&x.MaximumIngredientCount,
		&x.MinimumInstrumentCount,
		&x.MaximumInstrumentCount,
		&x.TemperatureRequired,
		&x.TimeEstimateRequired,
		&x.ConditionExpressionRequired,
		&x.ConsumesVessel,
		&x.OnlyForVessels,
		&x.MinimumVesselCount,
		&x.MaximumVesselCount,
		&x.Slug,
		&x.PastTense,
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

// scanValidPreparations takes some database rows and turns them into a slice of valid preparations.
func (q *Querier) scanValidPreparations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validPreparations []*types.ValidPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidPreparation(ctx, rows, includeCounts)
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

		validPreparations = append(validPreparations, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validPreparations, filteredCount, totalCount, nil
}

//go:embed queries/valid_preparations/exists.sql
var validPreparationExistenceQuery string

// ValidPreparationExists fetches whether a valid preparation exists from the database.
func (q *Querier) ValidPreparationExists(ctx context.Context, validPreparationID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	args := []any{
		validPreparationID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validPreparationExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid preparation existence check")
	}

	return result, nil
}

//go:embed queries/valid_preparations/get_one.sql
var getValidPreparationQuery string

// GetValidPreparation fetches a valid preparation from the database.
func (q *Querier) GetValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	args := []any{
		validPreparationID,
	}

	row := q.getOneRow(ctx, q.db, "validPreparation", getValidPreparationQuery, args)

	validPreparation, _, _, err := q.scanValidPreparation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning validPreparation")
	}

	return validPreparation, nil
}

//go:embed queries/valid_preparations/get_random.sql
var getRandomValidPreparationQuery string

// GetRandomValidPreparation fetches a valid preparation from the database.
func (q *Querier) GetRandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	args := []any{}

	row := q.getOneRow(ctx, q.db, "validPreparation", getRandomValidPreparationQuery, args)

	validPreparation, _, _, err := q.scanValidPreparation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning validPreparation")
	}

	return validPreparation, nil
}

//go:embed queries/valid_preparations/search.sql
var validPreparationSearchQuery string

// SearchForValidPreparations fetches a valid preparation from the database.
func (q *Querier) SearchForValidPreparations(ctx context.Context, query string) ([]*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidPreparationIDToSpan(span, query)

	args := []any{
		wrapQueryForILIKE(query),
	}

	rows, err := q.getRows(ctx, q.db, "valid preparations", validPreparationSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparations list retrieval query")
	}

	x, _, _, err := q.scanValidPreparations(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid preparations")
	}

	return x, nil
}

// GetValidPreparations fetches a list of valid preparations from the database that meet a particular filter.
func (q *Querier) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparation], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.ValidPreparation]{}
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

	query, args := q.buildListQuery(ctx, validPreparationsTable, nil, nil, nil, householdOwnershipColumn, validPreparationsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "valid preparations", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparations list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidPreparations(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid preparations")
	}

	return x, nil
}

// GetValidPreparationsWithIDs fetches a list of valid preparations from the database that meet a particular filter.
func (q *Querier) GetValidPreparationsWithIDs(ctx context.Context, ids []string) ([]*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	where := squirrel.Eq{"valid_preparations.id": ids}
	query, args := q.buildListQuery(ctx, validPreparationsTable, nil, nil, where, householdOwnershipColumn, validPreparationsTableColumns, "", false, nil)

	rows, err := q.getRows(ctx, q.db, "valid preparations", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparations id list retrieval query")
	}

	preparations, _, _, err := q.scanValidPreparations(ctx, rows, true)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid preparations")
	}

	return preparations, nil
}

//go:embed queries/valid_preparations/get_needing_indexing.sql
var validPreparationsNeedingIndexingQuery string

// GetValidPreparationIDsThatNeedSearchIndexing fetches a list of valid preparations from the database that meet a particular filter.
func (q *Querier) GetValidPreparationIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	rows, err := q.getRows(ctx, q.db, "valid preparations needing indexing", validPreparationsNeedingIndexingQuery, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid preparations list retrieval query")
	}

	return q.scanIDs(ctx, rows)
}

//go:embed queries/valid_preparations/create.sql
var validPreparationCreationQuery string

// CreateValidPreparation creates a valid preparation in the database.
func (q *Querier) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationDatabaseCreationInput) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationIDKey, input.ID)

	args := []any{
		input.ID,
		input.Name,
		input.Description,
		input.IconPath,
		input.YieldsNothing,
		input.RestrictToIngredients,
		input.MinimumIngredientCount,
		input.MaximumIngredientCount,
		input.MinimumInstrumentCount,
		input.MaximumInstrumentCount,
		input.TemperatureRequired,
		input.TimeEstimateRequired,
		input.ConditionExpressionRequired,
		input.ConsumesVessel,
		input.OnlyForVessels,
		input.MinimumVesselCount,
		input.MaximumVesselCount,
		input.PastTense,
		input.Slug,
	}

	// create the valid preparation.
	if err := q.performWriteQuery(ctx, q.db, "valid preparation creation", validPreparationCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid preparation creation query")
	}

	x := &types.ValidPreparation{
		ID:                          input.ID,
		Name:                        input.Name,
		Description:                 input.Description,
		IconPath:                    input.IconPath,
		YieldsNothing:               input.YieldsNothing,
		RestrictToIngredients:       input.RestrictToIngredients,
		Slug:                        input.Slug,
		PastTense:                   input.PastTense,
		MinimumIngredientCount:      input.MinimumIngredientCount,
		MaximumIngredientCount:      input.MaximumIngredientCount,
		MinimumInstrumentCount:      input.MinimumInstrumentCount,
		MaximumInstrumentCount:      input.MaximumInstrumentCount,
		TemperatureRequired:         input.TemperatureRequired,
		TimeEstimateRequired:        input.TimeEstimateRequired,
		ConditionExpressionRequired: input.ConditionExpressionRequired,
		ConsumesVessel:              input.ConsumesVessel,
		OnlyForVessels:              input.OnlyForVessels,
		MinimumVesselCount:          input.MinimumVesselCount,
		MaximumVesselCount:          input.MaximumVesselCount,
		CreatedAt:                   q.currentTime(),
	}

	tracing.AttachValidPreparationIDToSpan(span, x.ID)
	logger.Info("valid preparation created")

	return x, nil
}

//go:embed queries/valid_preparations/update.sql
var updateValidPreparationQuery string

// UpdateValidPreparation updates a particular valid preparation.
func (q *Querier) UpdateValidPreparation(ctx context.Context, updated *types.ValidPreparation) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationIDKey, updated.ID)
	tracing.AttachValidPreparationIDToSpan(span, updated.ID)

	args := []any{
		updated.Name,
		updated.Description,
		updated.IconPath,
		updated.YieldsNothing,
		updated.RestrictToIngredients,
		updated.MinimumIngredientCount,
		updated.MaximumIngredientCount,
		updated.MinimumInstrumentCount,
		updated.MaximumInstrumentCount,
		updated.TemperatureRequired,
		updated.TimeEstimateRequired,
		updated.ConditionExpressionRequired,
		updated.ConsumesVessel,
		updated.OnlyForVessels,
		updated.MinimumVesselCount,
		updated.MaximumVesselCount,
		updated.Slug,
		updated.PastTense,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation update", updateValidPreparationQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation updated")

	return nil
}

//go:embed queries/valid_preparations/update_last_indexed_at.sql
var updateValidPreparationLastIndexedAtQuery string

// MarkValidPreparationAsIndexed updates a particular valid preparation's last_indexed_at value.
func (q *Querier) MarkValidPreparationAsIndexed(ctx context.Context, validPreparationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	args := []any{
		validPreparationID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation last_indexed_at", updateValidPreparationLastIndexedAtQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid preparation as indexed")
	}

	logger.Info("valid preparation marked as indexed")

	return nil
}

//go:embed queries/valid_preparations/archive.sql
var archiveValidPreparationQuery string

// ArchiveValidPreparation archives a valid preparation from the database by its ID.
func (q *Querier) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	args := []any{
		validPreparationID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation archive", archiveValidPreparationQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation archived")

	return nil
}
