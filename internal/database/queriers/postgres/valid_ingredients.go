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

var (
	_ types.ValidIngredientDataManager = (*SQLQuerier)(nil)

	// validIngredientsTableColumns are the columns for the valid_ingredients table.
	validIngredientsTableColumns = []string{
		"valid_ingredients.id",
		"valid_ingredients.name",
		"valid_ingredients.variant",
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
		"valid_ingredients.animal_derived",
		"valid_ingredients.volumetric",
		"valid_ingredients.icon_path",
		"valid_ingredients.created_on",
		"valid_ingredients.last_updated_on",
		"valid_ingredients.archived_on",
	}
)

// scanValidIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient struct.
func (q *SQLQuerier) scanValidIngredient(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredient, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidIngredient{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Variant,
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
		&x.AnimalDerived,
		&x.Volumetric,
		&x.IconPath,
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

// scanValidIngredients takes some database rows and turns them into a slice of valid ingredients.
func (q *SQLQuerier) scanValidIngredients(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredients []*types.ValidIngredient, filteredCount, totalCount uint64, err error) {
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

const validIngredientExistenceQuery = "SELECT EXISTS ( SELECT valid_ingredients.id FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.id = $1 )"

// ValidIngredientExists fetches whether a valid ingredient exists from the database.
func (q *SQLQuerier) ValidIngredientExists(ctx context.Context, validIngredientID string) (exists bool, err error) {
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

const getValidIngredientQuery = "SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.volumetric, valid_ingredients.icon_path, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.id = $1"

// GetValidIngredient fetches a valid ingredient from the database.
func (q *SQLQuerier) GetValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
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

const validIngredientSearchQuery = "SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.variant, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.animal_derived, valid_ingredients.volumetric, valid_ingredients.icon_path, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients WHERE valid_ingredients.name ILIKE $1 AND valid_ingredients.archived_on IS NULL LIMIT 50"

// SearchForValidIngredients fetches a valid ingredient from the database.
func (q *SQLQuerier) SearchForValidIngredients(ctx context.Context, query string) ([]*types.ValidIngredient, error) {
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

const getTotalValidIngredientsCountQuery = "SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL"

// GetTotalValidIngredientCount fetches the count of valid ingredients from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalValidIngredientCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getTotalValidIngredientsCountQuery, "fetching count of valid ingredients")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid ingredients")
	}

	return count, nil
}

// GetValidIngredients fetches a list of valid ingredients from the database that meet a particular filter.
func (q *SQLQuerier) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (x *types.ValidIngredientList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidIngredientList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "valid_ingredients", nil, nil, nil, householdOwnershipColumn, validIngredientsTableColumns, "", false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "validIngredients", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	if x.ValidIngredients, x.FilteredCount, x.TotalCount, err = q.scanValidIngredients(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredients")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetValidIngredientsWithIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"valid_ingredients.id":          ids,
		"valid_ingredients.archived_on": nil,
	}

	subqueryBuilder := q.sqlBuilder.Select(validIngredientsTableColumns...).
		From("valid_ingredients").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(validIngredientsTableColumns...).
		FromSelect(subqueryBuilder, "valid_ingredients").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetValidIngredientsWithIDs fetches valid ingredients from the database within a given set of IDs.
func (q *SQLQuerier) GetValidIngredientsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ids == nil {
		return nil, ErrNilInputProvided
	}

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.buildGetValidIngredientsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid ingredients with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid ingredients from database")
	}

	validIngredients, _, _, err := q.scanValidIngredients(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredients")
	}

	return validIngredients, nil
}

const validIngredientCreationQuery = "INSERT INTO valid_ingredients (id,name,variant,description,warning,contains_egg,contains_dairy,contains_peanut,contains_tree_nut,contains_soy,contains_wheat,contains_shellfish,contains_sesame,contains_fish,contains_gluten,animal_flesh,animal_derived,volumetric,icon_path) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)"

// CreateValidIngredient creates a valid ingredient in the database.
func (q *SQLQuerier) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientDatabaseCreationInput) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Name,
		input.Variant,
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
		input.AnimalDerived,
		input.Volumetric,
		input.IconPath,
	}

	// create the valid ingredient.
	if err := q.performWriteQuery(ctx, q.db, "valid ingredient creation", validIngredientCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing valid ingredient creation query")
	}

	x := &types.ValidIngredient{
		ID:                input.ID,
		Name:              input.Name,
		Variant:           input.Variant,
		Description:       input.Description,
		Warning:           input.Warning,
		ContainsEgg:       input.ContainsEgg,
		ContainsDairy:     input.ContainsDairy,
		ContainsPeanut:    input.ContainsPeanut,
		ContainsTreeNut:   input.ContainsTreeNut,
		ContainsSoy:       input.ContainsSoy,
		ContainsWheat:     input.ContainsWheat,
		ContainsShellfish: input.ContainsShellfish,
		ContainsSesame:    input.ContainsSesame,
		ContainsFish:      input.ContainsFish,
		ContainsGluten:    input.ContainsGluten,
		AnimalFlesh:       input.AnimalFlesh,
		AnimalDerived:     input.AnimalDerived,
		Volumetric:        input.Volumetric,
		IconPath:          input.IconPath,
		CreatedOn:         q.currentTime(),
	}

	tracing.AttachValidIngredientIDToSpan(span, x.ID)
	logger.Info("valid ingredient created")

	return x, nil
}

const updateValidIngredientQuery = "UPDATE valid_ingredients SET name = $1, variant = $2, description = $3, warning = $4, contains_egg = $5, contains_dairy = $6, contains_peanut = $7, contains_tree_nut = $8, contains_soy = $9, contains_wheat = $10, contains_shellfish = $11, contains_sesame = $12, contains_fish = $13, contains_gluten = $14, animal_flesh = $15, animal_derived = $16, volumetric = $17, icon_path = $18, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $19"

// UpdateValidIngredient updates a particular valid ingredient.
func (q *SQLQuerier) UpdateValidIngredient(ctx context.Context, updated *types.ValidIngredient) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientIDKey, updated.ID)
	tracing.AttachValidIngredientIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Name,
		updated.Variant,
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
		updated.AnimalDerived,
		updated.Volumetric,
		updated.IconPath,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient update", updateValidIngredientQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient")
	}

	logger.Info("valid ingredient updated")

	return nil
}

const archiveValidIngredientQuery = "UPDATE valid_ingredients SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"

// ArchiveValidIngredient archives a valid ingredient from the database by its ID.
func (q *SQLQuerier) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
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
