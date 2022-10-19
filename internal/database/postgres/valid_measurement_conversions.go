package postgres

import (
	"context"
	_ "embed"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
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
func (q *Querier) scanValidMeasurementConversion(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidMeasurementConversion, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidMeasurementConversion{}
	var (
		ingredient = &types.NullableValidIngredient{}
		modifier   int64
	)

	targetVars := []interface{}{
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
		&ingredient.CreatedAt,
		&ingredient.LastUpdatedAt,
		&ingredient.ArchivedAt,

		// rest of the valid measurement conversion
		&modifier,
		&x.Notes,
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

	if ingredient.ID != nil {
		x.OnlyForIngredient = converters.ConvertNullableValidIngredientToValidIngredient(ingredient)
	}

	x.Modifier = float32(modifier) / float32(types.ValidMeasurementConversionQuantityModifier)

	return x, filteredCount, totalCount, nil
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

	args := []interface{}{
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
func (q *Querier) GetValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) (*types.ValidMeasurementConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)

	args := []interface{}{
		validMeasurementConversionID,
	}

	row := q.getOneRow(ctx, q.db, "validMeasurementConversion", getValidMeasurementConversionQuery, args)

	validMeasurementConversion, _, _, err := q.scanValidMeasurementConversion(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning validMeasurementConversion")
	}

	return validMeasurementConversion, nil
}

//go:embed queries/valid_measurement_conversions/create.sql
var validMeasurementConversionCreationQuery string

// CreateValidMeasurementConversion creates a valid measurement conversion in the database.
func (q *Querier) CreateValidMeasurementConversion(ctx context.Context, input *types.ValidMeasurementConversionDatabaseCreationInput) (*types.ValidMeasurementConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementConversionIDKey, input.ID)

	dbSafeModifier := int64(input.Modifier * types.ValidMeasurementConversionQuantityModifier)

	args := []interface{}{
		input.ID,
		input.From,
		input.To,
		input.ForIngredient,
		dbSafeModifier,
		input.Notes,
	}

	// create the valid measurement conversion.
	if err := q.performWriteQuery(ctx, q.db, "valid measurement conversion creation", validMeasurementConversionCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid measurement conversion creation query")
	}

	x := &types.ValidMeasurementConversion{
		ID:        input.ID,
		From:      types.ValidMeasurementUnit{ID: input.From},
		To:        types.ValidMeasurementUnit{ID: input.To},
		Modifier:  float32(dbSafeModifier / types.ValidMeasurementConversionQuantityModifier),
		Notes:     input.Notes,
		CreatedAt: q.currentTime(),
	}

	if input.ForIngredient != nil {
		x.OnlyForIngredient = &types.ValidIngredient{ID: *input.ForIngredient}
	}

	tracing.AttachValidMeasurementConversionIDToSpan(span, x.ID)
	logger.Info("valid measurement conversion created")

	return x, nil
}

//go:embed queries/valid_measurement_conversions/update.sql
var updateValidMeasurementConversionQuery string

// UpdateValidMeasurementConversion updates a particular valid measurement conversion.
func (q *Querier) UpdateValidMeasurementConversion(ctx context.Context, updated *types.ValidMeasurementConversion) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementConversionIDKey, updated.ID)
	tracing.AttachValidMeasurementConversionIDToSpan(span, updated.ID)

	dbSafeModifier := int64(updated.Modifier * types.ValidMeasurementConversionQuantityModifier)

	var ingredientID *string
	if updated.OnlyForIngredient != nil {
		ingredientID = &updated.OnlyForIngredient.ID
	}

	args := []interface{}{
		updated.From.ID,
		updated.To.ID,
		ingredientID,
		dbSafeModifier,
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

	args := []interface{}{
		validMeasurementConversionID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid measurement conversion archive", archiveValidMeasurementConversionQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement conversion")
	}

	logger.Info("valid measurement conversion archived")

	return nil
}
