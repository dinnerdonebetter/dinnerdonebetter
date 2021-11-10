package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	database "github.com/prixfixeco/api_server/internal/database"
	observability "github.com/prixfixeco/api_server/internal/observability"
	keys "github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.ValidIngredientPreparationDataManager = (*SQLQuerier)(nil)

	// validIngredientPreparationsTableColumns are the columns for the valid_ingredient_preparations table.
	validIngredientPreparationsTableColumns = []string{
		"valid_ingredient_preparations.id",
		"valid_ingredient_preparations.notes",
		"valid_ingredient_preparations.valid_preparation_id",
		"valid_ingredient_preparations.valid_ingredient_id",
		"valid_ingredient_preparations.created_on",
		"valid_ingredient_preparations.last_updated_on",
		"valid_ingredient_preparations.archived_on",
	}
)

// scanValidIngredientPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient preparation struct.
func (q *SQLQuerier) scanValidIngredientPreparation(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredientPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidIngredientPreparation{}

	targetVars := []interface{}{
		&x.ID,
		&x.Notes,
		&x.ValidPreparationID,
		&x.ValidIngredientID,
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

// scanValidIngredientPreparations takes some database rows and turns them into a slice of valid ingredient preparations.
func (q *SQLQuerier) scanValidIngredientPreparations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredientPreparations []*types.ValidIngredientPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

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
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validIngredientPreparations, filteredCount, totalCount, nil
}

const validIngredientPreparationExistenceQuery = "SELECT EXISTS ( SELECT valid_ingredient_preparations.id FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.id = $1 )"

// ValidIngredientPreparationExists fetches whether a valid ingredient preparation exists from the database.
func (q *SQLQuerier) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientPreparationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	args := []interface{}{
		validIngredientPreparationID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validIngredientPreparationExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid ingredient preparation existence check")
	}

	return result, nil
}

const getValidIngredientPreparationQuery = "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.id = $1"

// GetValidIngredientPreparation fetches a valid ingredient preparation from the database.
func (q *SQLQuerier) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	args := []interface{}{
		validIngredientPreparationID,
	}

	row := q.getOneRow(ctx, q.db, "validIngredientPreparation", getValidIngredientPreparationQuery, args)

	validIngredientPreparation, _, _, err := q.scanValidIngredientPreparation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning validIngredientPreparation")
	}

	return validIngredientPreparation, nil
}

const getTotalValidIngredientPreparationsCountQuery = "SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL"

// GetTotalValidIngredientPreparationCount fetches the count of valid ingredient preparations from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalValidIngredientPreparationCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getTotalValidIngredientPreparationsCountQuery, "fetching count of valid ingredient preparations")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid ingredient preparations")
	}

	return count, nil
}

// GetValidIngredientPreparations fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *SQLQuerier) GetValidIngredientPreparations(ctx context.Context, filter *types.QueryFilter) (x *types.ValidIngredientPreparationList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidIngredientPreparationList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "valid_ingredient_preparations", nil, nil, nil, householdOwnershipColumn, validIngredientPreparationsTableColumns, "", false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "validIngredientPreparations", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredient preparations list retrieval query")
	}

	if x.ValidIngredientPreparations, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientPreparations(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient preparations")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetValidIngredientPreparationsWithIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"valid_ingredient_preparations.id":          ids,
		"valid_ingredient_preparations.archived_on": nil,
	}

	subqueryBuilder := q.sqlBuilder.Select(validIngredientPreparationsTableColumns...).
		From("valid_ingredient_preparations").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(validIngredientPreparationsTableColumns...).
		FromSelect(subqueryBuilder, "valid_ingredient_preparations").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetValidIngredientPreparationsWithIDs fetches valid ingredient preparations from the database within a given set of IDs.
func (q *SQLQuerier) GetValidIngredientPreparationsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidIngredientPreparation, error) {
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

	query, args := q.buildGetValidIngredientPreparationsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid ingredient preparations with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid ingredient preparations from database")
	}

	validIngredientPreparations, _, _, err := q.scanValidIngredientPreparations(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredient preparations")
	}

	return validIngredientPreparations, nil
}

const validIngredientPreparationCreationQuery = "INSERT INTO valid_ingredient_preparations (id,notes,valid_preparation_id,valid_ingredient_id) VALUES ($1,$2,$3,$4)"

// CreateValidIngredientPreparation creates a valid ingredient preparation in the database.
func (q *SQLQuerier) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationDatabaseCreationInput) (*types.ValidIngredientPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientPreparationIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Notes,
		input.ValidPreparationID,
		input.ValidIngredientID,
	}

	// create the valid ingredient preparation.
	if err := q.performWriteQuery(ctx, q.db, "valid ingredient preparation creation", validIngredientPreparationCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing valid ingredient preparation creation query")
	}

	x := &types.ValidIngredientPreparation{
		ID:                 input.ID,
		Notes:              input.Notes,
		ValidPreparationID: input.ValidPreparationID,
		ValidIngredientID:  input.ValidIngredientID,
		CreatedOn:          q.currentTime(),
	}

	tracing.AttachValidIngredientPreparationIDToSpan(span, x.ID)
	logger.Info("valid ingredient preparation created")

	return x, nil
}

const updateValidIngredientPreparationQuery = "UPDATE valid_ingredient_preparations SET notes = $1, valid_preparation_id = $2, valid_ingredient_id = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4"

// UpdateValidIngredientPreparation updates a particular valid ingredient preparation.
func (q *SQLQuerier) UpdateValidIngredientPreparation(ctx context.Context, updated *types.ValidIngredientPreparation) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientPreparationIDKey, updated.ID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Notes,
		updated.ValidPreparationID,
		updated.ValidIngredientID,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient preparation update", updateValidIngredientPreparationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient preparation")
	}

	logger.Info("valid ingredient preparation updated")

	return nil
}

const archiveValidIngredientPreparationQuery = "UPDATE valid_ingredient_preparations SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"

// ArchiveValidIngredientPreparation archives a valid ingredient preparation from the database by its ID.
func (q *SQLQuerier) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	args := []interface{}{
		validIngredientPreparationID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient preparation archive", archiveValidIngredientPreparationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid ingredient preparation")
	}

	logger.Info("valid ingredient preparation archived")

	return nil
}
