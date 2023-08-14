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
	_ types.ValidMeasurementUnitConversionDataManager = (*Querier)(nil)

	// validMeasurementUnitConversionsTableColumns are the columns for the valid_measurement_conversions table.
	validMeasurementUnitConversionsTableColumns = []string{
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

// scanValidMeasurementUnitConversion takes a database Scanner (i.e. *sql.Row) and scans the result into a valid measurement conversion struct.
func (q *Querier) scanValidMeasurementUnitConversion(ctx context.Context, scan database.Scanner) (x *types.ValidMeasurementUnitConversion, err error) {
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

// scanValidMeasurementUnitConversions takes some database rows and turns them into a slice of valid measurement conversions.
func (q *Querier) scanValidMeasurementUnitConversions(ctx context.Context, rows database.ResultIterator) (validMeasurementUnitConversions []*types.ValidMeasurementUnitConversion, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, scanErr := q.scanValidMeasurementUnitConversion(ctx, rows)
		if scanErr != nil {
			return nil, scanErr
		}

		validMeasurementUnitConversions = append(validMeasurementUnitConversions, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "handling rows")
	}

	return validMeasurementUnitConversions, nil
}

//go:embed queries/valid_measurement_conversions/exists.sql
var validMeasurementUnitConversionExistenceQuery string

// ValidMeasurementUnitConversionExists fetches whether a valid measurement conversion exists from the database.
func (q *Querier) ValidMeasurementUnitConversionExists(ctx context.Context, validMeasurementUnitConversionID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitConversionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	args := []any{
		validMeasurementUnitConversionID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validMeasurementUnitConversionExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid measurement conversion existence check")
	}

	return result, nil
}

//go:embed queries/valid_measurement_conversions/get_one.sql
var getValidMeasurementUnitConversionQuery string

// GetValidMeasurementUnitConversion fetches a valid measurement conversion from the database.
func (q *Querier) GetValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	args := []any{
		validMeasurementUnitConversionID,
	}

	row := q.getOneRow(ctx, q.db, "valid measurement conversion", getValidMeasurementUnitConversionQuery, args)

	validMeasurementUnitConversion, err := q.scanValidMeasurementUnitConversion(ctx, row)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement conversion")
	}

	return validMeasurementUnitConversion, nil
}

//go:embed queries/valid_measurement_conversions/get_all_from_measurement_unit.sql
var getValidMeasurementUnitConversionsFromUnitQuery string

// GetValidMeasurementUnitConversionsFromUnit fetches a valid measurement conversions from a given measurement unit.
func (q *Querier) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	getValidMeasurementUnitConversionsFromUnitArgs := []any{
		validMeasurementUnitID,
	}

	rows, err := q.getRows(ctx, q.db, "valid measurement conversion", getValidMeasurementUnitConversionsFromUnitQuery, getValidMeasurementUnitConversionsFromUnitArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for valid measurement conversions")
	}

	validMeasurementUnitConversions, err := q.scanValidMeasurementUnitConversions(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement conversions")
	}

	return validMeasurementUnitConversions, nil
}

//go:embed queries/valid_measurement_conversions/get_all_to_measurement_unit.sql
var getValidMeasurementUnitConversionsToUnitQuery string

// GetValidMeasurementUnitConversionsToUnit fetches a valid measurement conversions to a given measurement unit.
func (q *Querier) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	getValidMeasurementUnitConversionsToUnitArgs := []any{
		validMeasurementUnitID,
	}

	rows, err := q.getRows(ctx, q.db, "valid measurement conversion", getValidMeasurementUnitConversionsToUnitQuery, getValidMeasurementUnitConversionsToUnitArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for valid measurement conversions")
	}

	validMeasurementUnitConversions, err := q.scanValidMeasurementUnitConversions(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement conversions")
	}

	return validMeasurementUnitConversions, nil
}

//go:embed queries/valid_measurement_conversions/create.sql
var validMeasurementUnitConversionCreationQuery string

// CreateValidMeasurementUnitConversion creates a valid measurement conversion in the database.
func (q *Querier) CreateValidMeasurementUnitConversion(ctx context.Context, input *types.ValidMeasurementUnitConversionDatabaseCreationInput) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, input.ID)

	args := []any{
		input.ID,
		input.From,
		input.To,
		input.OnlyForIngredient,
		input.Modifier,
		input.Notes,
	}

	// create the valid measurement conversion.
	if err := q.performWriteQuery(ctx, q.db, "valid measurement conversion creation", validMeasurementUnitConversionCreationQuery, args); err != nil {
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

	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, x.ID)
	logger.Info("valid measurement conversion created")

	return x, nil
}

//go:embed queries/valid_measurement_conversions/update.sql
var updateValidMeasurementUnitConversionQuery string

// UpdateValidMeasurementUnitConversion updates a particular valid measurement conversion.
func (q *Querier) UpdateValidMeasurementUnitConversion(ctx context.Context, updated *types.ValidMeasurementUnitConversion) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, updated.ID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, updated.ID)

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

	if err := q.performWriteQuery(ctx, q.db, "valid measurement conversion update", updateValidMeasurementUnitConversionQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement conversion")
	}

	logger.Info("valid measurement conversion updated")

	return nil
}

// ArchiveValidMeasurementUnitConversion archives a valid measurement conversion from the database by its ID.
func (q *Querier) ArchiveValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitConversionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	if err := q.generatedQuerier.ArchiveValidMeasurementUnitConversion(ctx, q.db, validMeasurementUnitConversionID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement conversion")
	}

	logger.Info("valid measurement conversion archived")

	return nil
}
