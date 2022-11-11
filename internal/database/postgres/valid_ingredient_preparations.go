package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	validIngredientsOnValidIngredientPreparationsJoinClause  = "valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id"
	validPreparationsOnValidIngredientPreparationsJoinClause = "valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id"
)

var (
	_ types.ValidIngredientPreparationDataManager = (*Querier)(nil)

	// validIngredientPreparationsTableColumns are the columns for the valid_ingredient_preparations table.
	validIngredientPreparationsTableColumns = []string{
		"valid_ingredient_preparations.id",
		"valid_ingredient_preparations.notes",
		"valid_preparations.id",
		"valid_preparations.name",
		"valid_preparations.description",
		"valid_preparations.icon_path",
		"valid_preparations.yields_nothing",
		"valid_preparations.restrict_to_ingredients",
		"valid_preparations.zero_ingredients_allowable",
		"valid_preparations.past_tense",
		"valid_preparations.created_at",
		"valid_preparations.last_updated_at",
		"valid_preparations.archived_at",
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
		"valid_ingredients.created_at",
		"valid_ingredients.last_updated_at",
		"valid_ingredients.archived_at",
		"valid_ingredient_preparations.created_at",
		"valid_ingredient_preparations.last_updated_at",
		"valid_ingredient_preparations.archived_at",
	}
)

// scanValidIngredientPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient preparation struct.
func (q *Querier) scanValidIngredientPreparation(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredientPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidIngredientPreparation{}

	targetVars := []any{
		&x.ID,
		&x.Notes,
		&x.Preparation.ID,
		&x.Preparation.Name,
		&x.Preparation.Description,
		&x.Preparation.IconPath,
		&x.Preparation.YieldsNothing,
		&x.Preparation.RestrictToIngredients,
		&x.Preparation.ZeroIngredientsAllowable,
		&x.Preparation.PastTense,
		&x.Preparation.CreatedAt,
		&x.Preparation.LastUpdatedAt,
		&x.Preparation.ArchivedAt,
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

// scanValidIngredientPreparations takes some database rows and turns them into a slice of valid ingredient preparations.
func (q *Querier) scanValidIngredientPreparations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredientPreparations []*types.ValidIngredientPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidIngredientPreparation(ctx, rows, includeCounts)
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

		validIngredientPreparations = append(validIngredientPreparations, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validIngredientPreparations, filteredCount, totalCount, nil
}

//go:embed queries/valid_ingredient_preparations/exists.sql
var validIngredientPreparationExistenceQuery string

// ValidIngredientPreparationExists fetches whether a valid ingredient preparation exists from the database.
func (q *Querier) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientPreparationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	args := []any{
		validIngredientPreparationID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validIngredientPreparationExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient preparation existence check")
	}

	return result, nil
}

//go:embed queries/valid_ingredient_preparations/get_one.sql
var getValidIngredientPreparationQuery string

// GetValidIngredientPreparation fetches a valid ingredient preparation from the database.
func (q *Querier) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	args := []any{
		validIngredientPreparationID,
	}

	row := q.getOneRow(ctx, q.db, "validIngredientPreparation", getValidIngredientPreparationQuery, args)

	validIngredientPreparation, _, _, err := q.scanValidIngredientPreparation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning validIngredientPreparation")
	}

	return validIngredientPreparation, nil
}

// GetValidIngredientPreparations fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) GetValidIngredientPreparations(ctx context.Context, filter *types.QueryFilter) (x *types.ValidIngredientPreparationList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidIngredientPreparationList{}
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
		validIngredientsOnValidIngredientPreparationsJoinClause,
		validPreparationsOnValidIngredientPreparationsJoinClause,
	}

	groupBys := []string{
		"valid_ingredients.id",
		"valid_preparations.id",
		"valid_ingredient_preparations.id",
	}

	query, args := q.buildListQuery(ctx, "valid_ingredient_preparations", joins, groupBys, nil, householdOwnershipColumn, validIngredientPreparationsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "validIngredientPreparations", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient preparations list retrieval query")
	}

	if x.ValidIngredientPreparations, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientPreparations(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient preparations")
	}

	return x, nil
}

func (q *Querier) buildGetValidIngredientPreparationsRestrictedByIDsQuery(ctx context.Context, column string, limit uint8, ids []string) (query string, args []any) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query, args, err := q.sqlBuilder.Select(validIngredientPreparationsTableColumns...).
		From("valid_ingredient_preparations").
		Join(validIngredientsOnValidIngredientPreparationsJoinClause).
		Join(validPreparationsOnValidIngredientPreparationsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("valid_ingredient_preparations.%s", column): ids,
			"valid_ingredient_preparations.archived_at":             nil,
		}).
		Limit(uint64(limit)).
		ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

func (q *Querier) buildGetValidIngredientPreparationsWithPreparationIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []any) {
	return q.buildGetValidIngredientPreparationsRestrictedByIDsQuery(ctx, "valid_preparation_id", limit, ids)
}

// GetValidIngredientPreparationsForPreparation fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) GetValidIngredientPreparationsForPreparation(ctx context.Context, preparationID string, filter *types.QueryFilter) (x *types.ValidIngredientPreparationList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, preparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, preparationID)

	x = &types.ValidIngredientPreparationList{}
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
	query, args := q.buildGetValidIngredientPreparationsWithPreparationIDsQuery(ctx, *filter.Limit, []string{preparationID})

	rows, err := q.getRows(ctx, q.db, "valid preparation instruments for preparation", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient preparations list retrieval query")
	}

	if x.ValidIngredientPreparations, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientPreparations(ctx, rows, false); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient preparations")
	}

	return x, nil
}

func (q *Querier) buildGetValidIngredientPreparationsWithIngredientIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []any) {
	return q.buildGetValidIngredientPreparationsRestrictedByIDsQuery(ctx, "valid_ingredient_id", limit, ids)
}

// GetValidIngredientPreparationsForIngredient fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) GetValidIngredientPreparationsForIngredient(ctx context.Context, ingredientID string, filter *types.QueryFilter) (x *types.ValidIngredientPreparationList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ingredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, ingredientID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, ingredientID)

	x = &types.ValidIngredientPreparationList{
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
	query, args := q.buildGetValidIngredientPreparationsWithIngredientIDsQuery(ctx, x.Limit, []string{ingredientID})

	rows, err := q.getRows(ctx, q.db, "valid preparation ingredients for ingredient", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient preparations list retrieval query")
	}

	if x.ValidIngredientPreparations, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientPreparations(ctx, rows, false); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient preparations")
	}

	return x, nil
}

//go:embed queries/valid_ingredient_preparations/create.sql
var validIngredientPreparationCreationQuery string

// CreateValidIngredientPreparation creates a valid ingredient preparation in the database.
func (q *Querier) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationDatabaseCreationInput) (*types.ValidIngredientPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientPreparationIDKey, input.ID)

	args := []any{
		input.ID,
		input.Notes,
		input.ValidPreparationID,
		input.ValidIngredientID,
	}

	// create the valid ingredient preparation.
	if err := q.performWriteQuery(ctx, q.db, "valid ingredient preparation creation", validIngredientPreparationCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient preparation creation query")
	}

	x := &types.ValidIngredientPreparation{
		ID:          input.ID,
		Notes:       input.Notes,
		Preparation: types.ValidPreparation{ID: input.ValidPreparationID},
		Ingredient:  types.ValidIngredient{ID: input.ValidIngredientID},
		CreatedAt:   q.currentTime(),
	}

	tracing.AttachValidIngredientPreparationIDToSpan(span, x.ID)
	logger.Info("valid ingredient preparation created")

	return x, nil
}

//go:embed queries/valid_ingredient_preparations/update.sql
var updateValidIngredientPreparationQuery string

// UpdateValidIngredientPreparation updates a particular valid ingredient preparation.
func (q *Querier) UpdateValidIngredientPreparation(ctx context.Context, updated *types.ValidIngredientPreparation) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientPreparationIDKey, updated.ID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, updated.ID)

	args := []any{
		updated.Notes,
		updated.Preparation.ID,
		updated.Ingredient.ID,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient preparation update", updateValidIngredientPreparationQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation")
	}

	logger.Info("valid ingredient preparation updated")

	return nil
}

//go:embed queries/valid_ingredient_preparations/archive.sql
var archiveValidIngredientPreparationQuery string

// ArchiveValidIngredientPreparation archives a valid ingredient preparation from the database by its ID.
func (q *Querier) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	args := []any{
		validIngredientPreparationID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient preparation archive", archiveValidIngredientPreparationQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation")
	}

	logger.Info("valid ingredient preparation archived")

	return nil
}
