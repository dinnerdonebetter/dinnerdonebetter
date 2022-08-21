package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
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
		"valid_instruments.plural_name",
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
		&x.PluralName,
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

	result, err := q.generatedQuerier.ValidInstrumentExists(ctx, validInstrumentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, observability.PrepareError(err, logger, span, "performing valid instrument existence check")
	}

	return result, nil
}

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

	result, err := q.generatedQuerier.GetValidInstrument(ctx, validInstrumentID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid instrument")
	}

	instrument := &types.ValidInstrument{
		Description: result.Description,
		IconPath:    result.IconPath,
		ID:          result.ID,
		Name:        result.Name,
		PluralName:  result.PluralName,
		CreatedOn:   uint64(result.CreatedOn),
	}

	if result.LastUpdatedOn.Valid {
		t := uint64(result.LastUpdatedOn.Int64)
		instrument.LastUpdatedOn = &t
	}

	if result.ArchivedOn.Valid {
		t := uint64(result.ArchivedOn.Int64)
		instrument.ArchivedOn = &t
	}

	return instrument, nil
}

// GetRandomValidInstrument fetches a valid instrument from the database.
func (q *SQLQuerier) GetRandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	result, err := q.generatedQuerier.GetRandomValidInstrument(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid instrument")
	}

	instrument := &types.ValidInstrument{
		Description: result.Description,
		IconPath:    result.IconPath,
		ID:          result.ID,
		Name:        result.Name,
		PluralName:  result.PluralName,
		CreatedOn:   uint64(result.CreatedOn),
	}

	if result.LastUpdatedOn.Valid {
		t := uint64(result.LastUpdatedOn.Int64)
		instrument.LastUpdatedOn = &t
	}

	if result.ArchivedOn.Valid {
		t := uint64(result.ArchivedOn.Int64)
		instrument.ArchivedOn = &t
	}

	return instrument, nil
}

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

	results, err := q.generatedQuerier.SearchForValidInstruments(ctx, wrapQueryForILIKE(query))
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredients search query")
	}

	validInstruments := []*types.ValidInstrument{}

	for _, result := range results {
		instrument := &types.ValidInstrument{
			Description: result.Description,
			IconPath:    result.IconPath,
			ID:          result.ID,
			Name:        result.Name,
			PluralName:  result.PluralName,
			CreatedOn:   uint64(result.CreatedOn),
		}

		if result.LastUpdatedOn.Valid {
			t := uint64(result.LastUpdatedOn.Int64)
			instrument.LastUpdatedOn = &t
		}

		if result.ArchivedOn.Valid {
			t := uint64(result.ArchivedOn.Int64)
			instrument.ArchivedOn = &t
		}

		validInstruments = append(validInstruments, instrument)
	}

	return validInstruments, nil
}

const getTotalValidInstrumentsCountQuery = "SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL"

// GetTotalValidInstrumentCount fetches the count of valid instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalValidInstrumentCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getTotalValidInstrumentsCountQuery, "fetching count of valid instruments")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid instruments")
	}

	return count, nil
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

// CreateValidInstrument creates a valid instrument in the database.
func (q *SQLQuerier) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentDatabaseCreationInput) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, input.ID)

	args := &generated.CreateValidInstrumentParams{
		ID:          input.ID,
		Name:        input.Name,
		PluralName:  input.PluralName,
		Description: input.Description,
		IconPath:    input.IconPath,
	}

	if err := q.generatedQuerier.CreateValidInstrument(ctx, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating valid instrument")
	}

	x := &types.ValidInstrument{
		Description: input.Description,
		IconPath:    input.IconPath,
		ID:          input.ID,
		Name:        input.Name,
		PluralName:  input.PluralName,
		CreatedOn:   q.currentTime(),
	}

	tracing.AttachValidInstrumentIDToSpan(span, x.ID)
	logger.Info("valid instrument created")

	return x, nil
}

// UpdateValidInstrument updates a particular valid instrument.
func (q *SQLQuerier) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, updated.ID)
	tracing.AttachValidInstrumentIDToSpan(span, updated.ID)

	args := &generated.UpdateValidInstrumentParams{
		ID:          updated.ID,
		Name:        updated.Name,
		PluralName:  updated.PluralName,
		Description: updated.Description,
		IconPath:    updated.IconPath,
	}

	if err := q.generatedQuerier.UpdateValidInstrument(ctx, args); err != nil {
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

	if err := q.generatedQuerier.ArchiveValidInstrument(ctx, validInstrumentID); err != nil {
		return observability.PrepareError(err, logger, span, "archiving valid instrument")
	}

	logger.Info("valid instrument archived")

	return nil
}
