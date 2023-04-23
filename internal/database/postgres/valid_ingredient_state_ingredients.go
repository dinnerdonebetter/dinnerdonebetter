package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"

	"github.com/Masterminds/squirrel"
)

const (
	validIngredientsOnValidIngredientStateIngredientsJoinClause      = "valid_ingredients ON valid_ingredient_state_ingredients.valid_ingredient = valid_ingredients.id"
	validIngredientStatesOnValidIngredientStateIngredientsJoinClause = "valid_ingredient_states ON valid_ingredient_state_ingredients.valid_ingredient_state = valid_ingredient_states.id"
)

var (
	_ types.ValidIngredientStateIngredientDataManager = (*Querier)(nil)

	// validIngredientStateIngredientsTableColumns are the columns for the valid_ingredient_state_ingredients table.
	validIngredientStateIngredientsTableColumns = []string{
		"valid_ingredient_state_ingredients.id",
		"valid_ingredient_state_ingredients.notes",
		"valid_ingredient_states.id",
		"valid_ingredient_states.name",
		"valid_ingredient_states.description",
		"valid_ingredient_states.icon_path",
		"valid_ingredient_states.slug",
		"valid_ingredient_states.past_tense",
		"valid_ingredient_states.attribute_type",
		"valid_ingredient_states.created_at",
		"valid_ingredient_states.last_updated_at",
		"valid_ingredient_states.archived_at",
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
		"valid_ingredient_state_ingredients.created_at",
		"valid_ingredient_state_ingredients.last_updated_at",
		"valid_ingredient_state_ingredients.archived_at",
	}
)

// scanValidIngredientStateIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient state ingredient struct.
func (q *Querier) scanValidIngredientStateIngredient(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredientStateIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidIngredientStateIngredient{}

	targetVars := []any{
		&x.ID,
		&x.Notes,
		&x.IngredientState.ID,
		&x.IngredientState.Name,
		&x.IngredientState.Description,
		&x.IngredientState.IconPath,
		&x.IngredientState.Slug,
		&x.IngredientState.PastTense,
		&x.IngredientState.AttributeType,
		&x.IngredientState.CreatedAt,
		&x.IngredientState.LastUpdatedAt,
		&x.IngredientState.ArchivedAt,
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
		&x.Ingredient.StorageInstructions,
		&x.Ingredient.Slug,
		&x.Ingredient.ContainsAlcohol,
		&x.Ingredient.ShoppingSuggestions,
		&x.Ingredient.IsStarch,
		&x.Ingredient.IsProtein,
		&x.Ingredient.IsGrain,
		&x.Ingredient.IsFruit,
		&x.Ingredient.IsSalt,
		&x.Ingredient.IsFat,
		&x.Ingredient.IsAcid,
		&x.Ingredient.IsHeat,
		&x.Ingredient.CreatedAt,
		&x.Ingredient.LastUpdatedAt,
		&x.Ingredient.ArchivedAt,
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

// scanValidIngredientStateIngredients takes some database rows and turns them into a slice of valid ingredient state ingredients.
func (q *Querier) scanValidIngredientStateIngredients(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredientStateIngredients []*types.ValidIngredientStateIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidIngredientStateIngredient(ctx, rows, includeCounts)
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

		validIngredientStateIngredients = append(validIngredientStateIngredients, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validIngredientStateIngredients, filteredCount, totalCount, nil
}

//go:embed queries/valid_ingredient_state_ingredients/exists.sql
var validIngredientStateIngredientExistenceQuery string

// ValidIngredientStateIngredientExists fetches whether a valid ingredient state ingredient exists from the database.
func (q *Querier) ValidIngredientStateIngredientExists(ctx context.Context, validIngredientStateIngredientID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)

	args := []any{
		validIngredientStateIngredientID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validIngredientStateIngredientExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient state ingredient existence check")
	}

	return result, nil
}

//go:embed queries/valid_ingredient_state_ingredients/get_one.sql
var getValidIngredientStateIngredientQuery string

// GetValidIngredientStateIngredient fetches a valid ingredient state ingredient from the database.
func (q *Querier) GetValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)

	args := []any{
		validIngredientStateIngredientID,
	}

	row := q.getOneRow(ctx, q.db, "valid ingredient state ingredient", getValidIngredientStateIngredientQuery, args)

	validIngredientStateIngredient, _, _, err := q.scanValidIngredientStateIngredient(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient state ingredient")
	}

	return validIngredientStateIngredient, nil
}

// GetValidIngredientStateIngredients fetches a list of valid ingredient state ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStateIngredients(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientStateIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{}
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
		validIngredientsOnValidIngredientStateIngredientsJoinClause,
		validIngredientStatesOnValidIngredientStateIngredientsJoinClause,
	}

	groupBys := []string{
		"valid_ingredients.id",
		"valid_ingredient_states.id",
		"valid_ingredient_state_ingredients.id",
	}

	query, args := q.buildListQuery(ctx, "valid_ingredient_state_ingredients", joins, groupBys, nil, householdOwnershipColumn, validIngredientStateIngredientsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "valid ingredient state ingredients", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient state ingredients list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientStateIngredients(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient state ingredients")
	}

	return x, nil
}

func (q *Querier) buildGetValidIngredientStateIngredientsRestrictedByIDsQuery(ctx context.Context, column string, limit uint8, ids []string) (query string, args []any) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query, args, err := q.sqlBuilder.Select(validIngredientStateIngredientsTableColumns...).
		From("valid_ingredient_state_ingredients").
		Join(validIngredientsOnValidIngredientStateIngredientsJoinClause).
		Join(validIngredientStatesOnValidIngredientStateIngredientsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("valid_ingredient_state_ingredients.%s", column): ids,
			"valid_ingredient_state_ingredients.archived_at":             nil,
		}).
		Limit(uint64(limit)).
		ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

func (q *Querier) buildGetValidIngredientStateIngredientsWithPreparationIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []any) {
	return q.buildGetValidIngredientStateIngredientsRestrictedByIDsQuery(ctx, "valid_ingredient_state", limit, ids)
}

// GetValidIngredientStateIngredientsForIngredientState fetches a list of valid ingredient state ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStateIngredientsForIngredientState(ctx context.Context, ingredientStateID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientStateIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ingredientStateID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, ingredientStateID)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, ingredientStateID)

	x = &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{}
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
	query, args := q.buildGetValidIngredientStateIngredientsWithPreparationIDsQuery(ctx, *filter.Limit, []string{ingredientStateID})

	rows, err := q.getRows(ctx, q.db, "valid ingredient state ingredients for preparation", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient state ingredients list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientStateIngredients(ctx, rows, false); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient state ingredients")
	}

	return x, nil
}

func (q *Querier) buildGetValidIngredientStateIngredientsWithIngredientIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []any) {
	return q.buildGetValidIngredientStateIngredientsRestrictedByIDsQuery(ctx, "valid_ingredient", limit, ids)
}

// GetValidIngredientStateIngredientsForIngredient fetches a list of valid ingredient state ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStateIngredientsForIngredient(ctx context.Context, ingredientID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientStateIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ingredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, ingredientID)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, ingredientID)

	x = &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{
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
	query, args := q.buildGetValidIngredientStateIngredientsWithIngredientIDsQuery(ctx, x.Limit, []string{ingredientID})

	rows, err := q.getRows(ctx, q.db, "valid ingredient state ingredients for ingredient", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient state ingredients list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientStateIngredients(ctx, rows, false); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient state ingredients")
	}

	return x, nil
}

//go:embed queries/valid_ingredient_state_ingredients/create.sql
var validIngredientStateIngredientCreationQuery string

// CreateValidIngredientStateIngredient creates a valid ingredient state ingredient in the database.
func (q *Querier) CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientDatabaseCreationInput) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientStateIngredientIDKey, input.ID)

	args := []any{
		input.ID,
		input.Notes,
		input.ValidIngredientStateID,
		input.ValidIngredientID,
	}

	// create the valid ingredient state ingredient.
	if err := q.performWriteQuery(ctx, q.db, "valid ingredient state ingredient creation", validIngredientStateIngredientCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient state ingredient creation query")
	}

	x := &types.ValidIngredientStateIngredient{
		ID:              input.ID,
		Notes:           input.Notes,
		IngredientState: types.ValidIngredientState{ID: input.ValidIngredientStateID},
		Ingredient:      types.ValidIngredient{ID: input.ValidIngredientID},
		CreatedAt:       q.currentTime(),
	}

	tracing.AttachValidIngredientStateIngredientIDToSpan(span, x.ID)
	logger.Info("valid ingredient state ingredient created")

	return x, nil
}

//go:embed queries/valid_ingredient_state_ingredients/update.sql
var updateValidIngredientStateIngredientQuery string

// UpdateValidIngredientStateIngredient updates a particular valid ingredient state ingredient.
func (q *Querier) UpdateValidIngredientStateIngredient(ctx context.Context, updated *types.ValidIngredientStateIngredient) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientStateIngredientIDKey, updated.ID)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, updated.ID)

	args := []any{
		updated.Notes,
		updated.IngredientState.ID,
		updated.Ingredient.ID,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient state ingredient update", updateValidIngredientStateIngredientQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state ingredient")
	}

	logger.Info("valid ingredient state ingredient updated")

	return nil
}

//go:embed queries/valid_ingredient_state_ingredients/archive.sql
var archiveValidIngredientStateIngredientQuery string

// ArchiveValidIngredientStateIngredient archives a valid ingredient state ingredient from the database by its ID.
func (q *Querier) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)

	args := []any{
		validIngredientStateIngredientID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient state ingredient archive", archiveValidIngredientStateIngredientQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state ingredient")
	}

	logger.Info("valid ingredient state ingredient archived")

	return nil
}
