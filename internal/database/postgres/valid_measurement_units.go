package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	validMeasurementUnitsOnRecipeStepIngredientsJoinClause = `valid_measurement_units ON recipe_step_ingredients.measurement_unit=valid_measurement_units.id`
	validMeasurementUnitsOnRecipeStepProductsJoinClause    = `valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id`
)

var (
	_ types.ValidMeasurementUnitDataManager = (*Querier)(nil)

	// validMeasurementUnitsTableColumns are the columns for the valid_measurement_units table.
	validMeasurementUnitsTableColumns = []string{
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
	}
)

// scanValidMeasurementUnit takes a database Scanner (i.e. *sql.Row) and scans the result into a valid measurement unit struct.
func (q *Querier) scanValidMeasurementUnit(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidMeasurementUnit, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidMeasurementUnit{}

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Description,
		&x.Volumetric,
		&x.IconPath,
		&x.Universal,
		&x.Metric,
		&x.Imperial,
		&x.Slug,
		&x.PluralName,
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

// scanValidMeasurementUnits takes some database rows and turns them into a slice of valid measurement units.
func (q *Querier) scanValidMeasurementUnits(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validMeasurementUnits []*types.ValidMeasurementUnit, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidMeasurementUnit(ctx, rows, includeCounts)
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

		validMeasurementUnits = append(validMeasurementUnits, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	if includeCounts {
		// when we include the counts, different filter values will return duplicate rows
		// so we deduplicate them here because it's easier than learning how to write SQL better.

		uniqueValidMeasurementUnits := map[string]struct{}{}
		for _, validMeasurementUnit := range validMeasurementUnits {
			uniqueValidMeasurementUnits[validMeasurementUnit.ID] = struct{}{}
		}

		filteredValidMeasurementUnits := []*types.ValidMeasurementUnit{}
		for _, validMeasurementUnit := range validMeasurementUnits {
			if _, ok := uniqueValidMeasurementUnits[validMeasurementUnit.ID]; ok {
				filteredValidMeasurementUnits = append(filteredValidMeasurementUnits, validMeasurementUnit)
			}
		}

		validMeasurementUnits = filteredValidMeasurementUnits
	}

	return validMeasurementUnits, filteredCount, totalCount, nil
}

//go:embed queries/valid_measurement_units/exists.sql
var validMeasurementUnitExistenceQuery string

// ValidMeasurementUnitExists fetches whether a valid measurement unit exists from the database.
func (q *Querier) ValidMeasurementUnitExists(ctx context.Context, validMeasurementUnitID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	args := []any{
		validMeasurementUnitID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validMeasurementUnitExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid measurement unit existence check")
	}

	return result, nil
}

//go:embed queries/valid_measurement_units/get_one.sql
var getValidMeasurementUnitQuery string

// GetValidMeasurementUnit fetches a valid measurement unit from the database.
func (q *Querier) GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	args := []any{
		validMeasurementUnitID,
	}

	row := q.getOneRow(ctx, q.db, "valid measurement unit", getValidMeasurementUnitQuery, args)

	validMeasurementUnit, _, _, err := q.scanValidMeasurementUnit(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement unit")
	}

	return validMeasurementUnit, nil
}

//go:embed queries/valid_measurement_units/get_random.sql
var getRandomValidMeasurementUnitQuery string

// GetRandomValidMeasurementUnit fetches a valid measurement unit from the database.
func (q *Querier) GetRandomValidMeasurementUnit(ctx context.Context) (*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	args := []any{}

	row := q.getOneRow(ctx, q.db, "valid measurement unit", getRandomValidMeasurementUnitQuery, args)

	validMeasurementUnit, _, _, err := q.scanValidMeasurementUnit(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid measurement unit")
	}

	return validMeasurementUnit, nil
}

//go:embed queries/valid_measurement_units/search.sql
var validMeasurementUnitSearchQuery string

// SearchForValidMeasurementUnitsByName fetches a valid measurement unit from the database.
func (q *Querier) SearchForValidMeasurementUnitsByName(ctx context.Context, query string) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidMeasurementUnitIDToSpan(span, query)

	args := []any{
		wrapQueryForILIKE(query),
	}

	rows, err := q.getRows(ctx, q.db, "valid measurement units", validMeasurementUnitSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	x, _, _, err := q.scanValidMeasurementUnits(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement units")
	}

	return x, nil
}

//go:embed queries/valid_measurement_units/search_by_ingredient_id.sql
var validMeasurementUnitSearchByIngredientIDQuery string

// ValidMeasurementUnitsForIngredientID fetches a valid measurement unit from the database.
func (q *Querier) ValidMeasurementUnitsForIngredientID(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	if validIngredientID == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	args := []any{
		filter.CreatedAfter,
		filter.CreatedBefore,
		filter.UpdatedAfter,
		filter.UpdatedBefore,
		validIngredientID,
		filter.Limit,
		filter.QueryOffset(),
	}

	rows, err := q.getRows(ctx, q.db, "valid measurement units", validMeasurementUnitSearchByIngredientIDQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	x := &types.QueryFilteredResult[types.ValidMeasurementUnit]{}
	x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidMeasurementUnits(ctx, rows, true)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement units")
	}

	if filter.Page != nil {
		x.Page = *filter.Page
	}

	if filter.Limit != nil {
		x.Limit = *filter.Limit
	}

	return x, nil
}

// GetValidMeasurementUnits fetches a list of valid measurement units from the database that meet a particular filter.
func (q *Querier) GetValidMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidMeasurementUnit], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.ValidMeasurementUnit]{}
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

	query, args := q.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "valid measurement units", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidMeasurementUnits(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement units")
	}

	return x, nil
}

//go:embed queries/valid_measurement_units/create.sql
var validMeasurementUnitCreationQuery string

// CreateValidMeasurementUnit creates a valid measurement unit in the database.
func (q *Querier) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitDatabaseCreationInput) (*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementUnitIDKey, input.ID)

	args := []any{
		input.ID,
		input.Name,
		input.Description,
		input.Volumetric,
		input.IconPath,
		input.Universal,
		input.Metric,
		input.Imperial,
		input.PluralName,
		input.Slug,
	}

	// create the valid measurement unit.
	if err := q.performWriteQuery(ctx, q.db, "valid measurement unit creation", validMeasurementUnitCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid measurement unit creation query")
	}

	x := &types.ValidMeasurementUnit{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Volumetric:  input.Volumetric,
		IconPath:    input.IconPath,
		Universal:   input.Universal,
		Metric:      input.Metric,
		Imperial:    input.Imperial,
		Slug:        input.Slug,
		PluralName:  input.PluralName,
		CreatedAt:   q.currentTime(),
	}

	tracing.AttachValidMeasurementUnitIDToSpan(span, x.ID)
	logger.Info("valid measurement unit created")

	return x, nil
}

//go:embed queries/valid_measurement_units/update.sql
var updateValidMeasurementUnitQuery string

// UpdateValidMeasurementUnit updates a particular valid measurement unit.
func (q *Querier) UpdateValidMeasurementUnit(ctx context.Context, updated *types.ValidMeasurementUnit) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementUnitIDKey, updated.ID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, updated.ID)

	args := []any{
		updated.Name,
		updated.Description,
		updated.Volumetric,
		updated.IconPath,
		updated.Universal,
		updated.Metric,
		updated.Imperial,
		updated.Slug,
		updated.PluralName,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid measurement unit update", updateValidMeasurementUnitQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement unit")
	}

	logger.Info("valid measurement unit updated")

	return nil
}

//go:embed queries/valid_measurement_units/update_last_indexed_at.sql
var updateValidMeasurementUnitLastIndexedAtQuery string

// MarkValidMeasurementUnitAsIndexed updates a particular valid measurement unit's last_indexed_at value.
func (q *Querier) MarkValidMeasurementUnitAsIndexed(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	args := []any{
		validMeasurementUnitID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid measurement unit last_indexed_at", updateValidMeasurementUnitLastIndexedAtQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid measurement unit as indexed")
	}

	logger.Info("valid measurement unit marked as indexed")

	return nil
}

//go:embed queries/valid_measurement_units/archive.sql
var archiveValidMeasurementUnitQuery string

// ArchiveValidMeasurementUnit archives a valid measurement unit from the database by its ID.
func (q *Querier) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	args := []any{
		validMeasurementUnitID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid measurement unit archive", archiveValidMeasurementUnitQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement unit")
	}

	logger.Info("valid measurement unit archived")

	return nil
}
