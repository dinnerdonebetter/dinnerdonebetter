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
	validMeasurementUnitsOnValidIngredientMeasurementUnitsJoinClause = "valid_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id"
	validIngredientsOnValidIngredientMeasurementUnitsJoinClause      = "valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id"
)

var (
	_ types.ValidIngredientMeasurementUnitDataManager = (*SQLQuerier)(nil)

	// validIngredientMeasurementUnitsTableColumns are the columns for the valid_ingredient_measurement_units table.
	validIngredientMeasurementUnitsTableColumns = []string{
		"valid_ingredient_measurement_units.id",
		"valid_ingredient_measurement_units.notes",
		"valid_measurement_units.id",
		"valid_measurement_units.name",
		"valid_measurement_units.description",
		"valid_measurement_units.volumetric",
		"valid_measurement_units.icon_path",
		"valid_measurement_units.created_on",
		"valid_measurement_units.last_updated_on",
		"valid_measurement_units.archived_on",
		"valid_ingredients.id",
		"valid_ingredients.name",
		"valid_ingredients.description",
		"valid_ingredients.warning",
		"valid_ingredients.contains_egg",
		"valid_ingredients.contains_dairy",
		"valid_ingredients.contains_peanut",
		"valid_ingredients.contains_tree_nut",
		"valid_ingredients.contains_soy",
		"valid_ingredients.contains_wheat",
		"valid_ingredients.contains_shellfish",
		"valid_ingredients.contains_sesame",
		"valid_ingredients.contains_fish",
		"valid_ingredients.contains_gluten",
		"valid_ingredients.animal_flesh",
		"valid_ingredients.volumetric",
		"valid_ingredients.is_liquid",
		"valid_ingredients.icon_path",
		"valid_ingredients.animal_derived",
		"valid_ingredients.plural_name",
		"valid_ingredients.restrict_to_preparations",
		"valid_ingredients.minimum_ideal_storage_temperature_in_celsius",
		"valid_ingredients.maximum_ideal_storage_temperature_in_celsius",
		"valid_ingredients.created_on",
		"valid_ingredients.last_updated_on",
		"valid_ingredients.archived_on",
		"valid_ingredient_measurement_units.created_on",
		"valid_ingredient_measurement_units.last_updated_on",
		"valid_ingredient_measurement_units.archived_on",
	}
)

// scanValidIngredientMeasurementUnit takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient measurement unit struct.
func (q *SQLQuerier) scanValidIngredientMeasurementUnit(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredientMeasurementUnit, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidIngredientMeasurementUnit{}

	targetVars := []interface{}{
		&x.ID,
		&x.Notes,
		&x.MeasurementUnit.ID,
		&x.MeasurementUnit.Name,
		&x.MeasurementUnit.Description,
		&x.MeasurementUnit.Volumetric,
		&x.MeasurementUnit.IconPath,
		&x.MeasurementUnit.CreatedOn,
		&x.MeasurementUnit.LastUpdatedOn,
		&x.MeasurementUnit.ArchivedOn,
		&x.Ingredient.ID,
		&x.Ingredient.Name,
		&x.Ingredient.Description,
		&x.Ingredient.Warning,
		&x.Ingredient.ContainsEgg,
		&x.Ingredient.ContainsDairy,
		&x.Ingredient.ContainsPeanut,
		&x.Ingredient.ContainsTreeNut,
		&x.Ingredient.ContainsSoy,
		&x.Ingredient.ContainsWheat,
		&x.Ingredient.ContainsShellfish,
		&x.Ingredient.ContainsSesame,
		&x.Ingredient.ContainsFish,
		&x.Ingredient.ContainsGluten,
		&x.Ingredient.AnimalFlesh,
		&x.Ingredient.IsMeasuredVolumetrically,
		&x.Ingredient.IsLiquid,
		&x.Ingredient.IconPath,
		&x.Ingredient.AnimalDerived,
		&x.Ingredient.PluralName,
		&x.Ingredient.RestrictToPreparations,
		&x.Ingredient.MinimumIdealStorageTemperatureInCelsius,
		&x.Ingredient.MaximumIdealStorageTemperatureInCelsius,
		&x.Ingredient.CreatedOn,
		&x.Ingredient.LastUpdatedOn,
		&x.Ingredient.ArchivedOn,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanValidIngredientMeasurementUnits takes some database rows and turns them into a slice of valid ingredient measurement units.
func (q *SQLQuerier) scanValidIngredientMeasurementUnits(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredientMeasurementUnits []*types.ValidIngredientMeasurementUnit, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidIngredientMeasurementUnit(ctx, rows, includeCounts)
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

		validIngredientMeasurementUnits = append(validIngredientMeasurementUnits, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validIngredientMeasurementUnits, filteredCount, totalCount, nil
}

const validIngredientMeasurementUnitExistenceQuery = "SELECT EXISTS ( SELECT valid_ingredient_measurement_units.id FROM valid_ingredient_measurement_units WHERE valid_ingredient_measurement_units.archived_on IS NULL AND valid_ingredient_measurement_units.id = $1 )"

// ValidIngredientMeasurementUnitExists fetches whether a valid ingredient measurement unit exists from the database.
func (q *SQLQuerier) ValidIngredientMeasurementUnitExists(ctx context.Context, validIngredientMeasurementUnitID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachValidIngredientMeasurementUnitIDToSpan(span, validIngredientMeasurementUnitID)

	args := []interface{}{
		validIngredientMeasurementUnitID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validIngredientMeasurementUnitExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid ingredient measurement unit existence check")
	}

	return result, nil
}

const getValidIngredientMeasurementUnitQuery = `
SELECT 
	valid_ingredient_measurement_units.id,
	valid_ingredient_measurement_units.notes,
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.created_on,
	valid_measurement_units.last_updated_on,
	valid_measurement_units.archived_on,
	valid_ingredients.id,
	valid_ingredients.name,
	valid_ingredients.description,
	valid_ingredients.warning,
	valid_ingredients.contains_egg,
	valid_ingredients.contains_dairy,
	valid_ingredients.contains_peanut,
	valid_ingredients.contains_tree_nut,
	valid_ingredients.contains_soy,
	valid_ingredients.contains_wheat,
	valid_ingredients.contains_shellfish,
	valid_ingredients.contains_sesame,
	valid_ingredients.contains_fish,
	valid_ingredients.contains_gluten,
	valid_ingredients.animal_flesh,
	valid_ingredients.volumetric,
	valid_ingredients.is_liquid,
	valid_ingredients.icon_path,
	valid_ingredients.animal_derived,
	valid_ingredients.plural_name,
	valid_ingredients.restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.created_on,
	valid_ingredients.last_updated_on,
	valid_ingredients.archived_on,
	valid_ingredient_measurement_units.created_on,
	valid_ingredient_measurement_units.last_updated_on,
	valid_ingredient_measurement_units.archived_on
FROM valid_ingredient_measurement_units 
JOIN valid_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE valid_ingredient_measurement_units.archived_on IS NULL 
  AND valid_ingredient_measurement_units.id = $1
`

// GetValidIngredientMeasurementUnit fetches a valid ingredient measurement unit from the database.
func (q *SQLQuerier) GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachValidIngredientMeasurementUnitIDToSpan(span, validIngredientMeasurementUnitID)

	args := []interface{}{
		validIngredientMeasurementUnitID,
	}

	row := q.getOneRow(ctx, q.db, "validIngredientMeasurementUnit", getValidIngredientMeasurementUnitQuery, args)

	validIngredientMeasurementUnit, _, _, err := q.scanValidIngredientMeasurementUnit(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning validIngredientMeasurementUnit")
	}

	return validIngredientMeasurementUnit, nil
}

const getTotalValidIngredientMeasurementUnitsCountQuery = "SELECT COUNT(valid_ingredient_measurement_units.id) FROM valid_ingredient_measurement_units WHERE valid_ingredient_measurement_units.archived_on IS NULL"

// GetTotalValidIngredientMeasurementUnitCount fetches the count of valid ingredient measurement units from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalValidIngredientMeasurementUnitCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getTotalValidIngredientMeasurementUnitsCountQuery, "fetching count of valid ingredient measurement units")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid ingredient measurement units")
	}

	return count, nil
}

func (q *SQLQuerier) buildGetValidIngredientMeasurementUnitRestrictedByIDsQuery(ctx context.Context, column string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query, args, err := q.sqlBuilder.Select(validIngredientMeasurementUnitsTableColumns...).
		From("valid_ingredient_measurement_units").
		Join(validIngredientsOnValidIngredientMeasurementUnitsJoinClause).
		Join(validMeasurementUnitsOnValidIngredientMeasurementUnitsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("valid_ingredient_measurement_units.%s", column): ids,
			"valid_ingredient_measurement_units.archived_on":             nil,
		}).
		Limit(uint64(limit)).
		ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

func (q *SQLQuerier) buildGetValidIngredientMeasurementUnitRestrictedByIngredientIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	return q.buildGetValidIngredientMeasurementUnitRestrictedByIDsQuery(ctx, "valid_ingredient_id", limit, ids)
}

// GetValidIngredientMeasurementUnitsForIngredient fetches a list of valid measurement units from the database that belong to a given ingredient ID.
func (q *SQLQuerier) GetValidIngredientMeasurementUnitsForIngredient(ctx context.Context, ingredientID string, filter *types.QueryFilter) (x *types.ValidIngredientMeasurementUnitList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ingredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, ingredientID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, ingredientID)

	x = &types.ValidIngredientMeasurementUnitList{
		Pagination: types.Pagination{
			Limit: 20,
		},
	}
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

	// the use of filter here is so weird, since we only respect the limit, but I'm trying to get this done, okay?
	query, args := q.buildGetValidIngredientMeasurementUnitRestrictedByIngredientIDsQuery(ctx, x.Limit, []string{ingredientID})

	rows, err := q.performReadQuery(ctx, q.db, "valid measurement units for ingredient", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredient measurement units list retrieval query")
	}

	if x.ValidIngredientMeasurementUnits, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientMeasurementUnits(ctx, rows, false); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient measurement units")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetValidIngredientMeasurementUnitsRestrictedByMeasurementUnitIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	return q.buildGetValidIngredientMeasurementUnitRestrictedByIDsQuery(ctx, "valid_measurement_unit_id", limit, ids)
}

// GetValidIngredientMeasurementUnitsForMeasurementUnit fetches a list of valid measurement units from the database that belong to a given ingredient ID.
func (q *SQLQuerier) GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx context.Context, validMeasurementUnitID string, filter *types.QueryFilter) (x *types.ValidIngredientMeasurementUnitList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validMeasurementUnitID)

	x = &types.ValidIngredientMeasurementUnitList{
		Pagination: types.Pagination{
			Limit: 20,
		},
	}
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

	// the use of filter here is so weird, since we only respect the limit, but I'm trying to get this done, okay?
	query, args := q.buildGetValidIngredientMeasurementUnitsRestrictedByMeasurementUnitIDsQuery(ctx, x.Limit, []string{validMeasurementUnitID})

	rows, err := q.performReadQuery(ctx, q.db, "valid measurement units for measurement unit", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredient measurement units list retrieval query")
	}

	if x.ValidIngredientMeasurementUnits, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientMeasurementUnits(ctx, rows, false); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient measurement units")
	}

	return x, nil
}

// GetValidIngredientMeasurementUnits fetches a list of valid ingredient measurement units from the database that meet a particular filter.
func (q *SQLQuerier) GetValidIngredientMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (x *types.ValidIngredientMeasurementUnitList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidIngredientMeasurementUnitList{}
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

	joins := []string{
		validMeasurementUnitsOnValidIngredientMeasurementUnitsJoinClause,
		validIngredientsOnValidIngredientMeasurementUnitsJoinClause,
	}

	groupBys := []string{
		"valid_ingredients.id",
		"valid_measurement_units.id",
		"valid_ingredient_measurement_units.id",
	}

	query, args := q.buildListQuery(ctx, "valid_ingredient_measurement_units", joins, groupBys, nil, householdOwnershipColumn, validIngredientMeasurementUnitsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "validIngredientMeasurementUnits", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredient measurement units list retrieval query")
	}

	if x.ValidIngredientMeasurementUnits, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientMeasurementUnits(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient measurement units")
	}

	return x, nil
}

const validIngredientMeasurementUnitCreationQuery = "INSERT INTO valid_ingredient_measurement_units (id,notes,valid_measurement_unit_id,valid_ingredient_id) VALUES ($1,$2,$3,$4)"

// CreateValidIngredientMeasurementUnit creates a valid ingredient measurement unit in the database.
func (q *SQLQuerier) CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitDatabaseCreationInput) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Notes,
		input.ValidMeasurementUnitID,
		input.ValidIngredientID,
	}

	// create the valid ingredient measurement unit.
	if err := q.performWriteQuery(ctx, q.db, "valid ingredient measurement unit creation", validIngredientMeasurementUnitCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing valid ingredient measurement unit creation query")
	}

	x := &types.ValidIngredientMeasurementUnit{
		ID:              input.ID,
		Notes:           input.Notes,
		MeasurementUnit: types.ValidMeasurementUnit{ID: input.ValidMeasurementUnitID},
		Ingredient:      types.ValidIngredient{ID: input.ValidIngredientID},
		CreatedOn:       q.currentTime(),
	}

	tracing.AttachValidIngredientMeasurementUnitIDToSpan(span, x.ID)
	logger.Info("valid ingredient measurement unit created")

	return x, nil
}

const updateValidIngredientMeasurementUnitQuery = "UPDATE valid_ingredient_measurement_units SET notes = $1, valid_measurement_unit_id = $2, valid_ingredient_id = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4"

// UpdateValidIngredientMeasurementUnit updates a particular valid ingredient measurement unit.
func (q *SQLQuerier) UpdateValidIngredientMeasurementUnit(ctx context.Context, updated *types.ValidIngredientMeasurementUnit) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, updated.ID)
	tracing.AttachValidIngredientMeasurementUnitIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Notes,
		updated.MeasurementUnit.ID,
		updated.Ingredient.ID,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient measurement unit update", updateValidIngredientMeasurementUnitQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient measurement unit")
	}

	logger.Info("valid ingredient measurement unit updated")

	return nil
}

const archiveValidIngredientMeasurementUnitQuery = "UPDATE valid_ingredient_measurement_units SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"

// ArchiveValidIngredientMeasurementUnit archives a valid ingredient measurement unit from the database by its ID.
func (q *SQLQuerier) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachValidIngredientMeasurementUnitIDToSpan(span, validIngredientMeasurementUnitID)

	args := []interface{}{
		validIngredientMeasurementUnitID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient measurement unit archive", archiveValidIngredientMeasurementUnitQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient measurement unit")
	}

	logger.Info("valid ingredient measurement unit archived")

	return nil
}
