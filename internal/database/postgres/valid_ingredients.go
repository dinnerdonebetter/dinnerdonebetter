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

	"github.com/Masterminds/squirrel"
)

const (
	validIngredientsTable = "valid_ingredients"
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
	}
)

// scanValidIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient struct.
func (q *Querier) scanValidIngredient(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidIngredient{}

	targetVars := []any{
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
		&x.Slug,
		&x.ContainsAlcohol,
		&x.ShoppingSuggestions,
		&x.IsStarch,
		&x.IsProtein,
		&x.IsGrain,
		&x.IsFruit,
		&x.IsSalt,
		&x.IsFat,
		&x.IsAcid,
		&x.IsHeat,
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

// scanValidIngredients takes some database rows and turns them into a slice of valid ingredients.
func (q *Querier) scanValidIngredients(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredients []*types.ValidIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

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
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validIngredients, filteredCount, totalCount, nil
}

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

	result, err := q.generatedQuerier.CheckValidIngredientExistence(ctx, q.db, validIngredientID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient existence check")
	}

	return result, nil
}

//go:embed queries/valid_ingredients/get_one.sql
var getValidIngredientQuery string

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

	args := []any{
		validIngredientID,
	}

	row := q.getOneRow(ctx, q.db, "valid ingredient", getValidIngredientQuery, args)

	validIngredient, _, _, err := q.scanValidIngredient(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient")
	}

	return validIngredient, nil
}

//go:embed queries/valid_ingredients/get_random.sql
var getRandomValidIngredientQuery string

// GetRandomValidIngredient fetches a valid ingredient from the database.
func (q *Querier) GetRandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	row := q.getOneRow(ctx, q.db, "valid ingredient", getRandomValidIngredientQuery, nil)

	validIngredient, _, _, err := q.scanValidIngredient(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid ingredient")
	}

	return validIngredient, nil
}

//go:embed queries/valid_ingredients/search.sql
var validIngredientSearchQuery string

// SearchForValidIngredients fetches a valid ingredient from the database.
func (q *Querier) SearchForValidIngredients(ctx context.Context, query string, filter *types.QueryFilter) ([]*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidIngredientIDToSpan(span, query)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	args := []any{
		query,
	}

	rows, err := q.getRows(ctx, q.db, "valid ingredients", validIngredientSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	x, _, _, err := q.scanValidIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredients")
	}

	return x, nil
}

//go:embed queries/valid_ingredients/search_by_preparation_and_ingredient_name.sql
var searchForIngredientsByPreparationAndIngredientNameQuery string

// SearchForValidIngredientsForPreparation fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, preparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, preparationID)

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachSearchQueryToSpan(span, query)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	x = &types.QueryFilteredResult[types.ValidIngredient]{
		Pagination: filter.ToPagination(),
	}

	searchForIngredientsByPreparationAndIngredientNameArgs := []any{
		preparationID,
		query,
	}

	rows, err := q.getRows(ctx, q.db, "valid ingredient preparations search by ingredient name", searchForIngredientsByPreparationAndIngredientNameQuery, searchForIngredientsByPreparationAndIngredientNameArgs)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient preparations search by ingredient name retrieval query")
	}

	if x.Data, _, _, err = q.scanValidIngredients(ctx, rows, false); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient preparations")
	}

	return x, nil
}

// SearchForValidIngredientsForIngredientState searches for valid ingredient sates.
func (q *Querier) SearchForValidIngredientsForIngredientState(ctx context.Context, ingredientStateID, query string, filter *types.QueryFilter) ([]*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ingredientStateID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, ingredientStateID)
	tracing.AttachValidIngredientStateIDToSpan(span, ingredientStateID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	args := []any{
		wrapQueryForILIKE(query),
	}

	rows, err := q.getRows(ctx, q.db, "valid ingredients search by ingredient state", validIngredientSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	x, _, _, err := q.scanValidIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredients")
	}

	return x, nil
}

// GetValidIngredients fetches a list of valid ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidIngredient]{
		Pagination: filter.ToPagination(),
	}

	query, args := q.buildListQuery(ctx, validIngredientsTable, nil, nil, nil, householdOwnershipColumn, validIngredientsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "valid ingredients", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidIngredients(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredients")
	}

	return x, nil
}

// GetValidIngredientsWithIDs fetches a list of valid ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientsWithIDs(ctx context.Context, ids []string) ([]*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	where := squirrel.Eq{"valid_ingredients.id": ids}
	query, args := q.buildListQuery(ctx, validIngredientsTable, nil, nil, where, householdOwnershipColumn, validIngredientsTableColumns, "", false, nil)

	rows, err := q.getRows(ctx, q.db, "valid ingredients", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients id list retrieval query")
	}

	ingredients, _, _, err := q.scanValidIngredients(ctx, rows, true)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredients")
	}

	return ingredients, nil
}

// GetValidIngredientIDsThatNeedSearchIndexing fetches a list of valid ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetValidIngredientsNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid ingredients list retrieval query")
	}

	return results, err
}

// CreateValidIngredient creates a valid ingredient in the database.
func (q *Querier) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientDatabaseCreationInput) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachValidIngredientIDToSpan(span, input.ID)
	logger := q.logger.WithValue(keys.ValidIngredientIDKey, input.ID)

	// create the valid ingredient.
	if err := q.generatedQuerier.CreateValidIngredient(ctx, q.db, &generated.CreateValidIngredientParams{
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
		Volumetric:                              input.IsMeasuredVolumetrically,
		IsLiquid:                                nullBoolFromBool(input.IsLiquid),
		IconPath:                                input.IconPath,
		AnimalDerived:                           input.AnimalDerived,
		PluralName:                              input.PluralName,
		RestrictToPreparations:                  input.RestrictToPreparations,
		MaximumIdealStorageTemperatureInCelsius: nullFloat64FromFloat32Pointer(input.MaximumIdealStorageTemperatureInCelsius),
		MinimumIdealStorageTemperatureInCelsius: nullFloat64FromFloat32Pointer(input.MinimumIdealStorageTemperatureInCelsius),
		StorageInstructions:                     input.StorageInstructions,
		Slug:                                    input.Slug,
		ContainsAlcohol:                         input.ContainsAlcohol,
		ShoppingSuggestions:                     input.ShoppingSuggestions,
		IsStarch:                                input.IsStarch,
		IsProtein:                               input.IsProtein,
		IsGrain:                                 input.IsGrain,
		IsFruit:                                 input.IsFruit,
		IsSalt:                                  input.IsSalt,
		IsFat:                                   input.IsFat,
		IsAcid:                                  input.IsAcid,
		IsHeat:                                  input.IsHeat,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient creation query")
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
		IsStarch:                                input.IsStarch,
		IsProtein:                               input.IsProtein,
		IsGrain:                                 input.IsGrain,
		IsFruit:                                 input.IsFruit,
		IsSalt:                                  input.IsSalt,
		IsFat:                                   input.IsFat,
		IsAcid:                                  input.IsAcid,
		IsHeat:                                  input.IsHeat,
		RestrictToPreparations:                  input.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: input.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: input.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     input.StorageInstructions,
		Slug:                                    input.Slug,
		ContainsAlcohol:                         input.ContainsAlcohol,
		ShoppingSuggestions:                     input.ShoppingSuggestions,
		CreatedAt:                               q.currentTime(),
	}

	tracing.AttachValidIngredientIDToSpan(span, x.ID)
	logger.Info("valid ingredient created")

	return x, nil
}

// UpdateValidIngredient updates a particular valid ingredient.
func (q *Querier) UpdateValidIngredient(ctx context.Context, updated *types.ValidIngredient) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidIngredientIDKey, updated.ID)
	tracing.AttachValidIngredientIDToSpan(span, updated.ID)

	if err := q.generatedQuerier.UpdateValidIngredient(ctx, q.db, &generated.UpdateValidIngredientParams{
		Description:                             updated.Description,
		Warning:                                 updated.Warning,
		ID:                                      updated.ID,
		ShoppingSuggestions:                     updated.ShoppingSuggestions,
		Slug:                                    updated.Slug,
		StorageInstructions:                     updated.StorageInstructions,
		Name:                                    updated.Name,
		PluralName:                              updated.PluralName,
		IconPath:                                updated.IconPath,
		MaximumIdealStorageTemperatureInCelsius: nullFloat64FromFloat32Pointer(updated.MaximumIdealStorageTemperatureInCelsius),
		MinimumIdealStorageTemperatureInCelsius: nullFloat64FromFloat32Pointer(updated.MinimumIdealStorageTemperatureInCelsius),
		IsLiquid:                                nullBoolFromBool(updated.IsLiquid),
		ContainsWheat:                           updated.ContainsWheat,
		ContainsPeanut:                          updated.ContainsPeanut,
		Volumetric:                              updated.IsMeasuredVolumetrically,
		ContainsGluten:                          updated.ContainsGluten,
		ContainsFish:                            updated.ContainsFish,
		AnimalDerived:                           updated.AnimalDerived,
		ContainsSesame:                          updated.ContainsSesame,
		RestrictToPreparations:                  updated.RestrictToPreparations,
		ContainsShellfish:                       updated.ContainsShellfish,
		ContainsSoy:                             updated.ContainsSoy,
		ContainsTreeNut:                         updated.ContainsTreeNut,
		AnimalFlesh:                             updated.AnimalFlesh,
		ContainsAlcohol:                         updated.ContainsAlcohol,
		ContainsDairy:                           updated.ContainsDairy,
		IsStarch:                                updated.IsStarch,
		IsProtein:                               updated.IsProtein,
		IsGrain:                                 updated.IsGrain,
		IsFruit:                                 updated.IsFruit,
		IsSalt:                                  updated.IsSalt,
		IsFat:                                   updated.IsFat,
		IsAcid:                                  updated.IsAcid,
		IsHeat:                                  updated.IsHeat,
		ContainsEgg:                             updated.ContainsEgg,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient")
	}

	logger.Info("valid ingredient updated")

	return nil
}

// MarkValidIngredientAsIndexed updates a particular valid ingredient's last_indexed_at value.
func (q *Querier) MarkValidIngredientAsIndexed(ctx context.Context, validIngredientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	if err := q.generatedQuerier.UpdateValidIngredientLastIndexedAt(ctx, q.db, validIngredientID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid ingredient as indexed")
	}

	logger.Info("valid ingredient marked as indexed")

	return nil
}

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

	if err := q.generatedQuerier.ArchiveValidIngredient(ctx, q.db, validIngredientID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient")
	}

	logger.Info("valid ingredient archived")

	return nil
}
