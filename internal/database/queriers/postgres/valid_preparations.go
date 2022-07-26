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
	validPreparationsOnRecipeStepsJoinClause = "valid_preparations ON recipe_steps.preparation_id=valid_preparations.id"
)

var (
	_ types.ValidPreparationDataManager = (*SQLQuerier)(nil)

	// validPreparationsTableColumns are the columns for the valid_preparations table.
	validPreparationsTableColumns = []string{
		"valid_preparations.id",
		"valid_preparations.name",
		"valid_preparations.description",
		"valid_preparations.icon_path",
		"valid_preparations.created_on",
		"valid_preparations.last_updated_on",
		"valid_preparations.archived_on",
	}
)

// scanValidPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a valid preparation struct.
func (q *SQLQuerier) scanValidPreparation(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidPreparation{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Description,
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

// scanValidPreparations takes some database rows and turns them into a slice of valid preparations.
func (q *SQLQuerier) scanValidPreparations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validPreparations []*types.ValidPreparation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidPreparation(ctx, rows, includeCounts)
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

		validPreparations = append(validPreparations, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validPreparations, filteredCount, totalCount, nil
}

const validPreparationExistenceQuery = "SELECT EXISTS ( SELECT valid_preparations.id FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.id = $1 )"

// ValidPreparationExists fetches whether a valid preparation exists from the database.
func (q *SQLQuerier) ValidPreparationExists(ctx context.Context, validPreparationID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	args := []interface{}{
		validPreparationID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validPreparationExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid preparation existence check")
	}

	return result, nil
}

const getValidPreparationBaseQuery = `SELECT
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.created_on,
	valid_preparations.last_updated_on,
	valid_preparations.archived_on
FROM valid_preparations
WHERE valid_preparations.archived_on IS NULL
`

const getValidPreparationQuery = getValidPreparationBaseQuery + `AND valid_preparations.id = $1`

// GetValidPreparation fetches a valid preparation from the database.
func (q *SQLQuerier) GetValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	args := []interface{}{
		validPreparationID,
	}

	row := q.getOneRow(ctx, q.db, "validPreparation", getValidPreparationQuery, args)

	validPreparation, _, _, err := q.scanValidPreparation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning validPreparation")
	}

	return validPreparation, nil
}

const getRandomValidPreparationQuery = getValidPreparationBaseQuery + `ORDER BY random() LIMIT 1`

// GetRandomValidPreparation fetches a valid preparation from the database.
func (q *SQLQuerier) GetRandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()
	args := []interface{}{}

	row := q.getOneRow(ctx, q.db, "validPreparation", getRandomValidPreparationQuery, args)

	validPreparation, _, _, err := q.scanValidPreparation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning validPreparation")
	}

	return validPreparation, nil
}

const validPreparationSearchQuery = "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon_path, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.name ILIKE $1 LIMIT 50"

// SearchForValidPreparations fetches a valid preparation from the database.
func (q *SQLQuerier) SearchForValidPreparations(ctx context.Context, query string) ([]*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidPreparationIDToSpan(span, query)

	args := []interface{}{
		wrapQueryForILIKE(query),
	}

	rows, err := q.performReadQuery(ctx, q.db, "valid ingredients", validPreparationSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	x, _, _, err := q.scanValidPreparations(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid ingredients")
	}

	return x, nil
}

// GetValidPreparations fetches a list of valid preparations from the database that meet a particular filter.
func (q *SQLQuerier) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (x *types.ValidPreparationList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidPreparationList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "valid_preparations", nil, nil, nil, householdOwnershipColumn, validPreparationsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "validPreparations", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid preparations list retrieval query")
	}

	if x.ValidPreparations, x.FilteredCount, x.TotalCount, err = q.scanValidPreparations(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparations")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetValidPreparationsWithIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"valid_preparations.id":          ids,
		"valid_preparations.archived_on": nil,
	}

	subqueryBuilder := q.sqlBuilder.Select(validPreparationsTableColumns...).
		From("valid_preparations").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(validPreparationsTableColumns...).
		FromSelect(subqueryBuilder, "valid_preparations").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetValidPreparationsWithIDs fetches valid preparations from the database within a given set of IDs.
func (q *SQLQuerier) GetValidPreparationsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidPreparation, error) {
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

	query, args := q.buildGetValidPreparationsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid preparations with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid preparations from database")
	}

	validPreparations, _, _, err := q.scanValidPreparations(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparations")
	}

	return validPreparations, nil
}

const validPreparationCreationQuery = "INSERT INTO valid_preparations (id,name,description,icon_path) VALUES ($1,$2,$3,$4)"

// CreateValidPreparation creates a valid preparation in the database.
func (q *SQLQuerier) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationDatabaseCreationInput) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Name,
		input.Description,
		input.IconPath,
	}

	// create the valid preparation.
	if err := q.performWriteQuery(ctx, q.db, "valid preparation creation", validPreparationCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing valid preparation creation query")
	}

	x := &types.ValidPreparation{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		IconPath:    input.IconPath,
		CreatedOn:   q.currentTime(),
	}

	tracing.AttachValidPreparationIDToSpan(span, x.ID)
	logger.Info("valid preparation created")

	return x, nil
}

const updateValidPreparationQuery = "UPDATE valid_preparations SET name = $1, description = $2, icon_path = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4"

// UpdateValidPreparation updates a particular valid preparation.
func (q *SQLQuerier) UpdateValidPreparation(ctx context.Context, updated *types.ValidPreparation) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationIDKey, updated.ID)
	tracing.AttachValidPreparationIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Name,
		updated.Description,
		updated.IconPath,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation update", updateValidPreparationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation updated")

	return nil
}

const archiveValidPreparationQuery = "UPDATE valid_preparations SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"

// ArchiveValidPreparation archives a valid preparation from the database by its ID.
func (q *SQLQuerier) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	args := []interface{}{
		validPreparationID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation archive", archiveValidPreparationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation archived")

	return nil
}
