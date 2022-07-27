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
	_ types.ValidPreparationInstrumentDataManager = (*SQLQuerier)(nil)

	// validPreparationInstrumentsTableColumns are the columns for the valid_preparation_instrument table.
	validPreparationInstrumentsTableColumns = []string{
		"valid_preparation_instruments.id",
		"valid_preparation_instruments.notes",
		"valid_preparation_instruments.valid_preparation_id",
		"valid_preparation_instruments.valid_instrument_id",
		"valid_preparation_instruments.created_on",
		"valid_preparation_instruments.last_updated_on",
		"valid_preparation_instruments.archived_on",
	}
)

// scanValidPreparationInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a valid preparation instrument struct.
func (q *SQLQuerier) scanValidPreparationInstrument(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidPreparationInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidPreparationInstrument{}

	targetVars := []interface{}{
		&x.ID,
		&x.Notes,
		&x.ValidPreparationID,
		&x.ValidInstrumentID,
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

// scanValidPreparationInstruments takes some database rows and turns them into a slice of valid preparation instruments.
func (q *SQLQuerier) scanValidPreparationInstruments(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validPreparationInstruments []*types.ValidPreparationInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidPreparationInstrument(ctx, rows, includeCounts)
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

		validPreparationInstruments = append(validPreparationInstruments, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validPreparationInstruments, filteredCount, totalCount, nil
}

const validPreparationInstrumentExistenceQuery = "SELECT EXISTS ( SELECT valid_preparation_instrument.id FROM valid_preparation_instrument WHERE valid_preparation_instrument.archived_on IS NULL AND valid_preparation_instrument.id = $1 )"

// ValidPreparationInstrumentExists fetches whether a valid preparation instrument exists from the database.
func (q *SQLQuerier) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	args := []interface{}{
		validPreparationInstrumentID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validPreparationInstrumentExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid preparation instrument existence check")
	}

	return result, nil
}

const getValidPreparationInstrumentQuery = "SELECT valid_preparation_instrument.id, valid_preparation_instrument.notes, valid_preparation_instrument.valid_preparation_id, valid_preparation_instrument.valid_instrument_id, valid_preparation_instrument.created_on, valid_preparation_instrument.last_updated_on, valid_preparation_instrument.archived_on FROM valid_preparation_instrument WHERE valid_preparation_instrument.archived_on IS NULL AND valid_preparation_instrument.id = $1"

// GetValidPreparationInstrument fetches a valid preparation instrument from the database.
func (q *SQLQuerier) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	args := []interface{}{
		validPreparationInstrumentID,
	}

	row := q.getOneRow(ctx, q.db, "validPreparationInstrument", getValidPreparationInstrumentQuery, args)

	validPreparationInstrument, _, _, err := q.scanValidPreparationInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning validPreparationInstrument")
	}

	return validPreparationInstrument, nil
}

// GetValidPreparationInstruments fetches a list of valid preparation instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.ValidPreparationInstrumentList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidPreparationInstrumentList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "valid_preparation_instrument", nil, nil, nil, householdOwnershipColumn, validPreparationInstrumentsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "validPreparationInstruments", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid preparation instruments list retrieval query")
	}

	if x.ValidPreparationInstruments, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationInstruments(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparation instruments")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetValidPreparationInstrumentsWithIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"valid_preparation_instrument.id":          ids,
		"valid_preparation_instrument.archived_on": nil,
	}

	subqueryBuilder := q.sqlBuilder.Select(validPreparationInstrumentsTableColumns...).
		From("valid_preparation_instrument").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(validPreparationInstrumentsTableColumns...).
		FromSelect(subqueryBuilder, "valid_preparation_instrument").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetValidPreparationInstrumentsWithIDs fetches valid preparation instruments from the database within a given set of IDs.
func (q *SQLQuerier) GetValidPreparationInstrumentsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidPreparationInstrument, error) {
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

	query, args := q.buildGetValidPreparationInstrumentsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid preparation instruments with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid preparation instruments from database")
	}

	validPreparationInstruments, _, _, err := q.scanValidPreparationInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid preparation instruments")
	}

	return validPreparationInstruments, nil
}

const validPreparationInstrumentCreationQuery = "INSERT INTO valid_preparation_instrument (id,notes,valid_preparation_id,valid_instrument_id) VALUES ($1,$2,$3,$4)"

// CreateValidPreparationInstrument creates a valid preparation instrument in the database.
func (q *SQLQuerier) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentDatabaseCreationInput) (*types.ValidPreparationInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Notes,
		input.ValidPreparationID,
		input.ValidInstrumentID,
	}

	// create the valid preparation instrument.
	if err := q.performWriteQuery(ctx, q.db, "valid preparation instrument creation", validPreparationInstrumentCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing valid preparation instrument creation query")
	}

	x := &types.ValidPreparationInstrument{
		ID:                 input.ID,
		Notes:              input.Notes,
		ValidPreparationID: input.ValidPreparationID,
		ValidInstrumentID:  input.ValidInstrumentID,
		CreatedOn:          q.currentTime(),
	}

	tracing.AttachValidPreparationInstrumentIDToSpan(span, x.ID)
	logger.Info("valid preparation instrument created")

	return x, nil
}

const updateValidPreparationInstrumentQuery = "UPDATE valid_preparation_instrument SET notes = $1, valid_preparation_id = $2, valid_instrument_id = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4"

// UpdateValidPreparationInstrument updates a particular valid preparation instrument.
func (q *SQLQuerier) UpdateValidPreparationInstrument(ctx context.Context, updated *types.ValidPreparationInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, updated.ID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Notes,
		updated.ValidPreparationID,
		updated.ValidInstrumentID,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation instrument update", updateValidPreparationInstrumentQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid preparation instrument")
	}

	logger.Info("valid preparation instrument updated")

	return nil
}

const archiveValidPreparationInstrumentQuery = "UPDATE valid_preparation_instrument SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"

// ArchiveValidPreparationInstrument archives a valid preparation instrument from the database by its ID.
func (q *SQLQuerier) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	args := []interface{}{
		validPreparationInstrumentID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation instrument archive", archiveValidPreparationInstrumentQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid preparation instrument")
	}

	logger.Info("valid preparation instrument archived")

	return nil
}
