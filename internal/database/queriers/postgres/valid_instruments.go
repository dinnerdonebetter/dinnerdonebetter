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
	_ types.ValidInstrumentDataManager = (*SQLQuerier)(nil)

	// validInstrumentsTableColumns are the columns for the valid_instruments table.
	validInstrumentsTableColumns = []string{
		"valid_instruments.id",
		"valid_instruments.name",
		"valid_instruments.variant",
		"valid_instruments.description",
		"valid_instruments.icon_path",
		"valid_instruments.created_on",
		"valid_instruments.last_updated_on",
		"valid_instruments.archived_on",
	}
)

// scanValidInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a valid instrument struct.
func (q *SQLQuerier) scanValidInstrument(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidInstrument{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Variant,
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

// scanValidInstruments takes some database rows and turns them into a slice of valid instruments.
func (q *SQLQuerier) scanValidInstruments(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validInstruments []*types.ValidInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidInstrument(ctx, rows, includeCounts)
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

		validInstruments = append(validInstruments, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validInstruments, filteredCount, totalCount, nil
}

const validInstrumentExistenceQuery = "SELECT EXISTS ( SELECT valid_instruments.id FROM valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.id = $1 )"

// ValidInstrumentExists fetches whether a valid instrument exists from the database.
func (q *SQLQuerier) ValidInstrumentExists(ctx context.Context, validInstrumentID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	args := []interface{}{
		validInstrumentID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validInstrumentExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid instrument existence check")
	}

	return result, nil
}

const getValidInstrumentBaseQuery = `SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.variant,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.created_on,
	valid_instruments.last_updated_on,
	valid_instruments.archived_on
FROM valid_instruments
WHERE valid_instruments.archived_on IS NULL
`

const getValidInstrumentQuery = getValidInstrumentBaseQuery + `AND valid_instruments.id = $1`

// GetValidInstrument fetches a valid instrument from the database.
func (q *SQLQuerier) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	args := []interface{}{
		validInstrumentID,
	}

	row := q.getOneRow(ctx, q.db, "validInstrument", getValidInstrumentQuery, args)

	validInstrument, _, _, err := q.scanValidInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning validInstrument")
	}

	return validInstrument, nil
}

const getRandomValidInstrumentQuery = getValidInstrumentBaseQuery + `ORDER BY random() LIMIT 1`

// GetRandomValidInstrument fetches a valid instrument from the database.
func (q *SQLQuerier) GetRandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()
	args := []interface{}{}

	row := q.getOneRow(ctx, q.db, "validInstrument", getRandomValidInstrumentQuery, args)

	validInstrument, _, _, err := q.scanValidInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning validInstrument")
	}

	return validInstrument, nil
}

const validInstrumentSearchQuery = `SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon_path, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.name ILIKE $1 LIMIT 50`

// SearchForValidInstruments fetches a valid instrument from the database.
func (q *SQLQuerier) SearchForValidInstruments(ctx context.Context, query string) ([]*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidInstrumentIDToSpan(span, query)

	args := []interface{}{
		wrapQueryForILIKE(query),
	}

	rows, err := q.performReadQuery(ctx, q.db, "valid ingredients", validInstrumentSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	validInstruments, _, _, err := q.scanValidInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning validInstrument")
	}

	return validInstruments, nil
}

// GetValidInstruments fetches a list of valid instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.ValidInstrumentList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidInstrumentList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "valid_instruments", nil, nil, nil, householdOwnershipColumn, validInstrumentsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "validInstruments", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid instruments list retrieval query")
	}

	if x.ValidInstruments, x.FilteredCount, x.TotalCount, err = q.scanValidInstruments(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid instruments")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetValidInstrumentsWithIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"valid_instruments.id":          ids,
		"valid_instruments.archived_on": nil,
	}

	subqueryBuilder := q.sqlBuilder.Select(validInstrumentsTableColumns...).
		From("valid_instruments").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(validInstrumentsTableColumns...).
		FromSelect(subqueryBuilder, "valid_instruments").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetValidInstrumentsWithIDs fetches valid instruments from the database within a given set of IDs.
func (q *SQLQuerier) GetValidInstrumentsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidInstrument, error) {
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

	query, args := q.buildGetValidInstrumentsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid instruments with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid instruments from database")
	}

	validInstruments, _, _, err := q.scanValidInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid instruments")
	}

	return validInstruments, nil
}

const validInstrumentCreationQuery = "INSERT INTO valid_instruments (id,name,variant,description,icon_path) VALUES ($1,$2,$3,$4,$5)"

// CreateValidInstrument creates a valid instrument in the database.
func (q *SQLQuerier) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentDatabaseCreationInput) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Name,
		input.Variant,
		input.Description,
		input.IconPath,
	}

	// create the valid instrument.
	if err := q.performWriteQuery(ctx, q.db, "valid instrument creation", validInstrumentCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing valid instrument creation query")
	}

	x := &types.ValidInstrument{
		ID:          input.ID,
		Name:        input.Name,
		Variant:     input.Variant,
		Description: input.Description,
		IconPath:    input.IconPath,
		CreatedOn:   q.currentTime(),
	}

	tracing.AttachValidInstrumentIDToSpan(span, x.ID)
	logger.Info("valid instrument created")

	return x, nil
}

const updateValidInstrumentQuery = "UPDATE valid_instruments SET name = $1, variant = $2, description = $3, icon_path = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $5"

// UpdateValidInstrument updates a particular valid instrument.
func (q *SQLQuerier) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, updated.ID)
	tracing.AttachValidInstrumentIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Name,
		updated.Variant,
		updated.Description,
		updated.IconPath,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid instrument update", updateValidInstrumentQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid instrument")
	}

	logger.Info("valid instrument updated")

	return nil
}

const archiveValidInstrumentQuery = "UPDATE valid_instruments SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"

// ArchiveValidInstrument archives a valid instrument from the database by its ID.
func (q *SQLQuerier) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	args := []interface{}{
		validInstrumentID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid instrument archive", archiveValidInstrumentQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid instrument")
	}

	logger.Info("valid instrument archived")

	return nil
}
