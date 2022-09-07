package postgres

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.ValidIngredientDataManager = (*Querier)(nil)

	// validIngredientsTableColumns are the columns for the valid_ingredients table.
	validIngredientsTableColumns = []string{
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
	}
)

// scanValidIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient struct.
func (q *Querier) scanValidIngredient(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidIngredient{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Description,
		&x.Warning,
		&x.ContainsEgg,
		&x.ContainsDairy,
		&x.ContainsPeanut,
		&x.ContainsTreeNut,
		&x.ContainsSoy,
		&x.ContainsWheat,
		&x.ContainsShellfish,
		&x.ContainsSesame,
		&x.ContainsFish,
		&x.ContainsGluten,
		&x.AnimalFlesh,
		&x.IsMeasuredVolumetrically,
		&x.IsLiquid,
		&x.IconPath,
		&x.AnimalDerived,
		&x.PluralName,
		&x.RestrictToPreparations,
		&x.MinimumIdealStorageTemperatureInCelsius,
		&x.MaximumIdealStorageTemperatureInCelsius,
		&x.StorageInstructions,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanValidIngredients takes some database rows and turns them into a slice of valid ingredients.
func (q *Querier) scanValidIngredients(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredients []*types.ValidIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidIngredient(ctx, rows, includeCounts)
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

		validIngredients = append(validIngredients, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validIngredients, filteredCount, totalCount, nil
}

const validIngredientExistenceQuery = "SELECT EXISTS ( SELECT valid_ingredients.id FROM valid_ingredients WHERE valid_ingredients.archived_at IS NULL AND valid_ingredients.id = $1 )"

// ValidIngredientExists fetches whether a valid ingredient exists from the database.
func (q *Querier) ValidIngredientExists(ctx context.Context, validIngredientID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	args := []interface{}{
		validIngredientID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validIngredientExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid ingredient existence check")
	}

	return result, nil
}

const getValidIngredientBaseQuery = `SELECT 
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
	valid_ingredients.storage_instructions,
	valid_ingredients.created_at, 
	valid_ingredients.last_updated_at, 
	valid_ingredients.archived_at 
FROM valid_ingredients 
WHERE valid_ingredients.archived_at IS NULL
`

const getValidIngredientQuery = getValidIngredientBaseQuery + `AND valid_ingredients.id = $1`

// GetValidIngredient fetches a valid ingredient from the database.
func (q *Querier) GetValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	args := []interface{}{
		validIngredientID,
	}

	row := q.getOneRow(ctx, q.db, "valid ingredient", getValidIngredientQuery, args)

	validIngredient, _, _, err := q.scanValidIngredient(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient")
	}

	return validIngredient, nil
}

const getRandomValidIngredientQuery = getValidIngredientBaseQuery + `ORDER BY random() LIMIT 1`

// GetRandomValidIngredient fetches a valid ingredient from the database.
func (q *Querier) GetRandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()
	args := []interface{}{}

	row := q.getOneRow(ctx, q.db, "valid ingredient", getRandomValidIngredientQuery, args)

	validIngredient, _, _, err := q.scanValidIngredient(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient")
	}

	return validIngredient, nil
}

const validIngredientSearchQuery = `SELECT 
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
	valid_ingredients.storage_instructions,
    valid_ingredients.created_at,
    valid_ingredients.last_updated_at,
    valid_ingredients.archived_at 
FROM valid_ingredients
WHERE valid_ingredients.name ILIKE $1 
    AND valid_ingredients.archived_at IS NULL
LIMIT 50`

// SearchForValidIngredients fetches a valid ingredient from the database.
func (q *Querier) SearchForValidIngredients(ctx context.Context, query string) ([]*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidIngredientIDToSpan(span, query)

	args := []interface{}{
		wrapQueryForILIKE(query),
	}

	rows, err := q.performReadQuery(ctx, q.db, "valid ingredients", validIngredientSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	x, _, _, err := q.scanValidIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredients")
	}

	return x, nil
}

// SearchForValidIngredientsForPreparation fetches a valid ingredient from the database.
func (q *Querier) SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string) ([]*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidIngredientIDToSpan(span, query)

	args := []interface{}{
		wrapQueryForILIKE(query),
	}

	// TODO: find some way to restrict by preparationID

	rows, err := q.performReadQuery(ctx, q.db, "valid ingredients", validIngredientSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	x, _, _, err := q.scanValidIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredients")
	}

	return x, nil
}

// GetValidIngredients fetches a list of valid ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (x *types.ValidIngredientList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidIngredientList{}
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

	query, args := q.buildListQuery(ctx, "valid_ingredients", nil, nil, nil, householdOwnershipColumn, validIngredientsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "validIngredients", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	if x.ValidIngredients, x.FilteredCount, x.TotalCount, err = q.scanValidIngredients(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredients")
	}

	return x, nil
}

const validIngredientCreationQuery = `INSERT INTO valid_ingredients
(
	id,
	name,
	description,
	warning,
	contains_egg,
	contains_dairy,
	contains_peanut,
	contains_tree_nut,
	contains_soy,
	contains_wheat,
	contains_shellfish,
	contains_sesame,
	contains_fish,
	contains_gluten,
	animal_flesh,
	volumetric,
	is_liquid,
	icon_path,
	animal_derived,
	plural_name,
	restrict_to_preparations,
	minimum_ideal_storage_temperature_in_celsius,
	maximum_ideal_storage_temperature_in_celsius,
 	storage_instructions
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)`

// CreateValidIngredient creates a valid ingredient in the database.
func (q *Querier) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientDatabaseCreationInput) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Name,
		input.Description,
		input.Warning,
		input.ContainsEgg,
		input.ContainsDairy,
		input.ContainsPeanut,
		input.ContainsTreeNut,
		input.ContainsSoy,
		input.ContainsWheat,
		input.ContainsShellfish,
		input.ContainsSesame,
		input.ContainsFish,
		input.ContainsGluten,
		input.AnimalFlesh,
		input.IsMeasuredVolumetrically,
		input.IsLiquid,
		input.IconPath,
		input.AnimalDerived,
		input.PluralName,
		input.RestrictToPreparations,
		input.MinimumIdealStorageTemperatureInCelsius,
		input.MaximumIdealStorageTemperatureInCelsius,
		input.StorageInstructions,
	}

	// create the valid ingredient.
	if err := q.performWriteQuery(ctx, q.db, "valid ingredient creation", validIngredientCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing valid ingredient creation query")
	}

	x := &types.ValidIngredient{
		ID:                                      input.ID,
		Name:                                    input.Name,
		Description:                             input.Description,
		Warning:                                 input.Warning,
		ContainsEgg:                             input.ContainsEgg,
		ContainsDairy:                           input.ContainsDairy,
		ContainsPeanut:                          input.ContainsPeanut,
		ContainsTreeNut:                         input.ContainsTreeNut,
		ContainsSoy:                             input.ContainsSoy,
		ContainsWheat:                           input.ContainsWheat,
		ContainsShellfish:                       input.ContainsShellfish,
		ContainsSesame:                          input.ContainsSesame,
		ContainsFish:                            input.ContainsFish,
		ContainsGluten:                          input.ContainsGluten,
		AnimalFlesh:                             input.AnimalFlesh,
		IsMeasuredVolumetrically:                input.IsMeasuredVolumetrically,
		IsLiquid:                                input.IsLiquid,
		IconPath:                                input.IconPath,
		AnimalDerived:                           input.AnimalDerived,
		PluralName:                              input.PluralName,
		RestrictToPreparations:                  input.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: input.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: input.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     input.StorageInstructions,
		CreatedAt:                               q.currentTime(),
	}

	tracing.AttachValidIngredientIDToSpan(span, x.ID)
	logger.Info("valid ingredient created")

	return x, nil
}

const updateValidIngredientQuery = `
UPDATE valid_ingredients SET 
	name = $1,
	description = $2,
	warning = $3,
	contains_egg = $4,
	contains_dairy = $5,
	contains_peanut = $6,
	contains_tree_nut = $7,
	contains_soy = $8,
	contains_wheat = $9,
	contains_shellfish = $10,
	contains_sesame = $11,
	contains_fish = $12,
	contains_gluten = $13,
	animal_flesh = $14,
	volumetric = $15,
	is_liquid = $16,
	icon_path = $17,
	animal_derived = $18,
	plural_name = $19,
	restrict_to_preparations = $20,
	minimum_ideal_storage_temperature_in_celsius = $21,
	maximum_ideal_storage_temperature_in_celsius = $22,
	storage_instructions = $23,
	last_updated_at = NOW() 
WHERE archived_at IS NULL AND id = $24
`

// UpdateValidIngredient updates a particular valid ingredient.
func (q *Querier) UpdateValidIngredient(ctx context.Context, updated *types.ValidIngredient) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientIDKey, updated.ID)
	tracing.AttachValidIngredientIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Name,
		updated.Description,
		updated.Warning,
		updated.ContainsEgg,
		updated.ContainsDairy,
		updated.ContainsPeanut,
		updated.ContainsTreeNut,
		updated.ContainsSoy,
		updated.ContainsWheat,
		updated.ContainsShellfish,
		updated.ContainsSesame,
		updated.ContainsFish,
		updated.ContainsGluten,
		updated.AnimalFlesh,
		updated.IsMeasuredVolumetrically,
		updated.IsLiquid,
		updated.IconPath,
		updated.AnimalDerived,
		updated.PluralName,
		updated.RestrictToPreparations,
		updated.MinimumIdealStorageTemperatureInCelsius,
		updated.MaximumIdealStorageTemperatureInCelsius,
		updated.StorageInstructions,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient update", updateValidIngredientQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient")
	}

	logger.Info("valid ingredient updated")

	return nil
}

const archiveValidIngredientQuery = `UPDATE valid_ingredients SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1`

// ArchiveValidIngredient archives a valid ingredient from the database by its ID.
func (q *Querier) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	args := []interface{}{
		validIngredientID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient archive", archiveValidIngredientQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient")
	}

	logger.Info("valid ingredient archived")

	return nil
}
