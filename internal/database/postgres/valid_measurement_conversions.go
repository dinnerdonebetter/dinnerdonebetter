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

var (
	_ types.ValidMeasurementConversionDataManager = (*Querier)(nil)

	// validMeasurementConversionsTableColumns are the columns for the valid_measurement_conversions table.
	validMeasurementConversionsTableColumns = []string{
		"valid_measurement_conversions.id",
		"valid_measurement_units_1.id",
		"valid_measurement_units_1.name",
		"valid_measurement_units_1.description",
		"valid_measurement_units_1.volumetric",
		"valid_measurement_units_1.icon_path",
		"valid_measurement_units_1.universal",
		"valid_measurement_units_1.metric",
		"valid_measurement_units_1.imperial",
		"valid_measurement_units_1.slug",
		"valid_measurement_units_1.plural_name",
		"valid_measurement_units_1.created_at",
		"valid_measurement_units_1.last_updated_at",
		"valid_measurement_units_2.archived_at",
		"valid_measurement_units_2.id",
		"valid_measurement_units_2.name",
		"valid_measurement_units_2.description",
		"valid_measurement_units_2.volumetric",
		"valid_measurement_units_2.icon_path",
		"valid_measurement_units_2.universal",
		"valid_measurement_units_2.metric",
		"valid_measurement_units_2.imperial",
		"valid_measurement_units_2.slug",
		"valid_measurement_units_2.plural_name",
		"valid_measurement_units_2.created_at",
		"valid_measurement_units_2.last_updated_at",
		"valid_measurement_units_2.archived_at",
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
		"valid_ingredients.storage_instructions",
		"valid_ingredients.slug",
		"valid_ingredients.contains_alcohol",
		"valid_ingredients.shopping_suggestions",
		"valid_ingredients.is_starch",
		"valid_ingredients.is_protein",
		"valid_ingredients.is_grain",
		"valid_ingredients.is_fruit",
		"valid_ingredients.is_salt",
		"valid_ingredients.is_fat",
		"valid_ingredients.is_acid",
		"valid_ingredients.is_heat",
		"valid_ingredients.created_at",
		"valid_ingredients.last_updated_at",
		"valid_ingredients.archived_at",
		"valid_measurement_conversions.modifier",
		"valid_measurement_conversions.notes",
		"valid_measurement_conversions.created_at",
		"valid_measurement_conversions.last_updated_at",
		"valid_measurement_conversions.archived_at",
	}
)

// scanValidMeasurementConversion takes a database Scanner (i.e. *sql.Row) and scans the result into a valid measurement conversion struct.
func (q *Querier) scanValidMeasurementConversion(ctx context.Context, scan database.Scanner) (x *types.ValidMeasurementUnitConversion, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidMeasurementUnitConversion{}
	var (
		ingredient = &types.NullableValidIngredient{}
	)

	targetVars := []any{
		&x.ID,

		// first valid measurement join
		&x.From.ID,
		&x.From.Name,
		&x.From.Description,
		&x.From.Volumetric,
		&x.From.IconPath,
		&x.From.Universal,
		&x.From.Metric,
		&x.From.Imperial,
		&x.From.Slug,
		&x.From.PluralName,
		&x.From.CreatedAt,
		&x.From.LastUpdatedAt,
		&x.From.ArchivedAt,

		// second valid measurement join
		&x.To.ID,
		&x.To.Name,
		&x.To.Description,
		&x.To.Volumetric,
		&x.To.IconPath,
		&x.To.Universal,
		&x.To.Metric,
		&x.To.Imperial,
		&x.To.Slug,
		&x.To.PluralName,
		&x.To.CreatedAt,
		&x.To.LastUpdatedAt,
		&x.To.ArchivedAt,

		// valid ingredient join
		&ingredient.ID,
		&ingredient.Name,
		&ingredient.Description,
		&ingredient.Warning,
		&ingredient.ContainsEgg,
		&ingredient.ContainsDairy,
		&ingredient.ContainsPeanut,
		&ingredient.ContainsTreeNut,
		&ingredient.ContainsSoy,
		&ingredient.ContainsWheat,
		&ingredient.ContainsShellfish,
		&ingredient.ContainsSesame,
		&ingredient.ContainsFish,
		&ingredient.ContainsGluten,
		&ingredient.AnimalFlesh,
		&ingredient.IsMeasuredVolumetrically,
		&ingredient.IsLiquid,
		&ingredient.IconPath,
		&ingredient.AnimalDerived,
		&ingredient.PluralName,
		&ingredient.RestrictToPreparations,
		&ingredient.MinimumIdealStorageTemperatureInCelsius,
		&ingredient.MaximumIdealStorageTemperatureInCelsius,
		&ingredient.StorageInstructions,
		&ingredient.Slug,
		&ingredient.ContainsAlcohol,
		&ingredient.ShoppingSuggestions,
		&ingredient.IsStarch,
		&ingredient.IsProtein,
		&ingredient.IsGrain,
		&ingredient.IsFruit,
		&ingredient.IsSalt,
		&ingredient.IsFat,
		&ingredient.IsAcid,
		&ingredient.IsHeat,
		&ingredient.CreatedAt,
		&ingredient.LastUpdatedAt,
		&ingredient.ArchivedAt,

		// rest of the valid measurement conversion
		&x.Modifier,
		&x.Notes,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, observability.PrepareError(err, span, "")
	}

	if ingredient.ID != nil {
		x.OnlyForIngredient = converters.ConvertNullableValidIngredientToValidIngredient(ingredient)
	}

	return x, nil
}

// scanValidMeasurementConversions takes some database rows and turns them into a slice of valid measurement conversions.
func (q *Querier) scanValidMeasurementConversions(ctx context.Context, rows database.ResultIterator) (validMeasurementConversions []*types.ValidMeasurementUnitConversion, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, scanErr := q.scanValidMeasurementConversion(ctx, rows)
		if scanErr != nil {
			return nil, scanErr
		}

		validMeasurementConversions = append(validMeasurementConversions, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "handling rows")
	}

	return validMeasurementConversions, nil
}

//go:embed queries/valid_measurement_conversions/exists.sql
var validMeasurementConversionExistenceQuery string

// ValidMeasurementConversionExists fetches whether a valid measurement conversion exists from the database.
func (q *Querier) ValidMeasurementConversionExists(ctx context.Context, validMeasurementConversionID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementConversionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)

	args := []any{
		validMeasurementConversionID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validMeasurementConversionExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid measurement conversion existence check")
	}

	return result, nil
}

//go:embed queries/valid_measurement_conversions/get_one.sql
var getValidMeasurementConversionQuery string

// GetValidMeasurementConversion fetches a valid measurement conversion from the database.
func (q *Querier) GetValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)

	args := []any{
		validMeasurementConversionID,
	}

	row := q.getOneRow(ctx, q.db, "valid measurement conversion", getValidMeasurementConversionQuery, args)

	validMeasurementConversion, err := q.scanValidMeasurementConversion(ctx, row)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement conversion")
	}

	return validMeasurementConversion, nil
}

//go:embed queries/valid_measurement_conversions/get_all_from_measurement_unit.sql
var getValidMeasurementConversionsFromUnitQuery string

// GetValidMeasurementConversionsFromUnit fetches a valid measurement conversions from a given measurement unit.
func (q *Querier) GetValidMeasurementConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	getValidMeasurementConversionsFromUnitArgs := []any{
		validMeasurementUnitID,
	}

	rows, err := q.getRows(ctx, q.db, "valid measurement conversion", getValidMeasurementConversionsFromUnitQuery, getValidMeasurementConversionsFromUnitArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for valid measurement conversions")
	}

	validMeasurementConversions, err := q.scanValidMeasurementConversions(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement conversions")
	}

	return validMeasurementConversions, nil
}

//go:embed queries/valid_measurement_conversions/get_all_to_measurement_unit.sql
var getValidMeasurementConversionsToUnitQuery string

// GetValidMeasurementConversionsToUnit fetches a valid measurement conversions to a given measurement unit.
func (q *Querier) GetValidMeasurementConversionsToUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	getValidMeasurementConversionsToUnitArgs := []any{
		validMeasurementUnitID,
	}

	rows, err := q.getRows(ctx, q.db, "valid measurement conversion", getValidMeasurementConversionsToUnitQuery, getValidMeasurementConversionsToUnitArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for valid measurement conversions")
	}

	validMeasurementConversions, err := q.scanValidMeasurementConversions(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement conversions")
	}

	return validMeasurementConversions, nil
}

//go:embed queries/valid_measurement_conversions/create.sql
var validMeasurementConversionCreationQuery string

// CreateValidMeasurementConversion creates a valid measurement conversion in the database.
func (q *Querier) CreateValidMeasurementConversion(ctx context.Context, input *types.ValidMeasurementConversionDatabaseCreationInput) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementConversionIDKey, input.ID)

	args := []any{
		input.ID,
		input.From,
		input.To,
		input.OnlyForIngredient,
		input.Modifier,
		input.Notes,
	}

	// create the valid measurement conversion.
	if err := q.performWriteQuery(ctx, q.db, "valid measurement conversion creation", validMeasurementConversionCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid measurement conversion creation query")
	}

	x := &types.ValidMeasurementUnitConversion{
		ID:        input.ID,
		From:      types.ValidMeasurementUnit{ID: input.From},
		To:        types.ValidMeasurementUnit{ID: input.To},
		Modifier:  input.Modifier,
		Notes:     input.Notes,
		CreatedAt: q.currentTime(),
	}

	if input.OnlyForIngredient != nil {
		x.OnlyForIngredient = &types.ValidIngredient{ID: *input.OnlyForIngredient}
	}

	tracing.AttachValidMeasurementConversionIDToSpan(span, x.ID)
	logger.Info("valid measurement conversion created")

	return x, nil
}

//go:embed queries/valid_measurement_conversions/update.sql
var updateValidMeasurementConversionQuery string

// UpdateValidMeasurementConversion updates a particular valid measurement conversion.
func (q *Querier) UpdateValidMeasurementConversion(ctx context.Context, updated *types.ValidMeasurementUnitConversion) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementConversionIDKey, updated.ID)
	tracing.AttachValidMeasurementConversionIDToSpan(span, updated.ID)

	var ingredientID *string
	if updated.OnlyForIngredient != nil {
		ingredientID = &updated.OnlyForIngredient.ID
	}

	args := []any{
		updated.From.ID,
		updated.To.ID,
		ingredientID,
		updated.Modifier,
		updated.Notes,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid measurement conversion update", updateValidMeasurementConversionQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement conversion")
	}

	logger.Info("valid measurement conversion updated")

	return nil
}

//go:embed queries/valid_measurement_conversions/archive.sql
var archiveValidMeasurementConversionQuery string

// ArchiveValidMeasurementConversion archives a valid measurement conversion from the database by its ID.
func (q *Querier) ArchiveValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementConversionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)

	args := []any{
		validMeasurementConversionID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid measurement conversion archive", archiveValidMeasurementConversionQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement conversion")
	}

	logger.Info("valid measurement conversion archived")

	return nil
}
