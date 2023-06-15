package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/Masterminds/squirrel"
)

const (
	validInstrumentsTable = "valid_instruments"
)

var (
	_ types.ValidInstrumentDataManager = (*Querier)(nil)

	// validInstrumentsTableColumns are the columns for the valid_instruments table.
	validInstrumentsTableColumns = []string{
		"valid_instruments.id",
		"valid_instruments.name",
		"valid_instruments.plural_name",
		"valid_instruments.description",
		"valid_instruments.icon_path",
		"valid_instruments.usable_for_storage",
		"valid_instruments.display_in_summary_lists",
		"valid_instruments.include_in_generated_instructions",
		"valid_instruments.is_vessel",
		"valid_instruments.is_exclusively_vessel",
		"valid_instruments.slug",
		"valid_instruments.created_at",
		"valid_instruments.last_updated_at",
		"valid_instruments.archived_at",
	}
)

// scanValidInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a valid instrument struct.
func (q *Querier) scanValidInstrument(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidInstrument{}

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.PluralName,
		&x.Description,
		&x.IconPath,
		&x.UsableForStorage,
		&x.DisplayInSummaryLists,
		&x.IncludeInGeneratedInstructions,
		&x.IsVessel,
		&x.IsExclusivelyVessel,
		&x.Slug,
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

// scanValidInstruments takes some database rows and turns them into a slice of valid instruments.
func (q *Querier) scanValidInstruments(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validInstruments []*types.ValidInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

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
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validInstruments, filteredCount, totalCount, nil
}

//go:embed queries/valid_instruments/exists.sql
var validInstrumentExistenceQuery string

// ValidInstrumentExists fetches whether a valid instrument exists from the database.
func (q *Querier) ValidInstrumentExists(ctx context.Context, validInstrumentID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	args := []any{
		validInstrumentID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validInstrumentExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid instrument existence check")
	}

	return result, nil
}

//go:embed queries/valid_instruments/get_one.sql
var getValidInstrumentQuery string

// GetValidInstrument fetches a valid instrument from the database.
func (q *Querier) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	args := []any{
		validInstrumentID,
	}

	row := q.getOneRow(ctx, q.db, "validInstrument", getValidInstrumentQuery, args)

	validInstrument, _, _, err := q.scanValidInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning validInstrument")
	}

	return validInstrument, nil
}

//go:embed queries/valid_instruments/get_random.sql
var getRandomValidInstrumentQuery string

// GetRandomValidInstrument fetches a valid instrument from the database.
func (q *Querier) GetRandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	args := []any{}

	row := q.getOneRow(ctx, q.db, "validInstrument", getRandomValidInstrumentQuery, args)

	validInstrument, _, _, err := q.scanValidInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning validInstrument")
	}

	return validInstrument, nil
}

//go:embed queries/valid_instruments/search.sql
var validInstrumentSearchQuery string

// SearchForValidInstruments fetches a valid instrument from the database.
func (q *Querier) SearchForValidInstruments(ctx context.Context, query string) ([]*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidInstrumentIDToSpan(span, query)

	args := []any{
		wrapQueryForILIKE(query),
	}

	rows, err := q.getRows(ctx, q.db, "valid ingredients", validInstrumentSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	validInstruments, _, _, err := q.scanValidInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning validInstrument")
	}

	return validInstruments, nil
}

// SearchForValidInstrumentsForPreparation fetches a valid instrument from the database.
func (q *Querier) SearchForValidInstrumentsForPreparation(ctx context.Context, preparationID, query string) ([]*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidInstrumentIDToSpan(span, query)

	// TODO: restrict results by preparation ID

	args := []any{
		wrapQueryForILIKE(query),
	}

	rows, err := q.getRows(ctx, q.db, "valid ingredients search", validInstrumentSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	validInstruments, _, _, err := q.scanValidInstruments(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning validInstrument")
	}

	return validInstruments, nil
}

// GetValidInstruments fetches a list of valid instruments from the database that meet a particular filter.
func (q *Querier) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidInstrument], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.ValidInstrument]{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	} else {
		filter = types.DefaultQueryFilter()
	}

	query, args := q.buildListQuery(ctx, validInstrumentsTable, nil, nil, nil, householdOwnershipColumn, validInstrumentsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "valid instruments", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid instruments list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidInstruments(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid instruments")
	}

	return x, nil
}

// GetValidInstrumentsWithIDs fetches a list of valid instruments from the database that meet a particular filter.
func (q *Querier) GetValidInstrumentsWithIDs(ctx context.Context, ids []string) ([]*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	where := squirrel.Eq{"valid_instruments.id": ids}
	query, args := q.buildListQuery(ctx, validInstrumentsTable, nil, nil, where, householdOwnershipColumn, validInstrumentsTableColumns, "", false, nil)

	rows, err := q.getRows(ctx, q.db, "valid instruments", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid instruments id list retrieval query")
	}

	instruments, _, _, err := q.scanValidInstruments(ctx, rows, true)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid instruments")
	}

	return instruments, nil
}

//go:embed generated_queries/valid_instruments/get_needing_indexing.sql
var validInstrumentsNeedingIndexingQuery string

// GetValidInstrumentIDsThatNeedSearchIndexing fetches a list of valid instruments from the database that meet a particular filter.
func (q *Querier) GetValidInstrumentIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	rows, err := q.getRows(ctx, q.db, "valid instruments needing indexing", validInstrumentsNeedingIndexingQuery, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid instruments list retrieval query")
	}

	return q.scanIDs(ctx, rows)
}

//go:embed queries/valid_instruments/create.sql
var validInstrumentCreationQuery string

// CreateValidInstrument creates a valid instrument in the database.
func (q *Querier) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentDatabaseCreationInput) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, input.ID)

	args := []any{
		input.ID,
		input.Name,
		input.PluralName,
		input.Description,
		input.IconPath,
		input.UsableForStorage,
		input.DisplayInSummaryLists,
		input.IncludeInGeneratedInstructions,
		input.IsVessel,
		input.IsExclusivelyVessel,
		input.Slug,
	}

	// create the valid instrument.
	if err := q.performWriteQuery(ctx, q.db, "valid instrument creation", validInstrumentCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid instrument creation query")
	}

	x := &types.ValidInstrument{
		ID:                             input.ID,
		Name:                           input.Name,
		PluralName:                     input.PluralName,
		Description:                    input.Description,
		IconPath:                       input.IconPath,
		UsableForStorage:               input.UsableForStorage,
		Slug:                           input.Slug,
		DisplayInSummaryLists:          input.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
		IsVessel:                       input.IsVessel,
		IsExclusivelyVessel:            input.IsExclusivelyVessel,
		CreatedAt:                      q.currentTime(),
	}

	tracing.AttachValidInstrumentIDToSpan(span, x.ID)
	logger.Info("valid instrument created")

	return x, nil
}

//go:embed queries/valid_instruments/update.sql
var updateValidInstrumentQuery string

// UpdateValidInstrument updates a particular valid instrument.
func (q *Querier) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, updated.ID)
	tracing.AttachValidInstrumentIDToSpan(span, updated.ID)

	args := []any{
		updated.Name,
		updated.PluralName,
		updated.Description,
		updated.IconPath,
		updated.UsableForStorage,
		updated.DisplayInSummaryLists,
		updated.IncludeInGeneratedInstructions,
		updated.IsVessel,
		updated.IsExclusivelyVessel,
		updated.Slug,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid instrument update", updateValidInstrumentQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid instrument")
	}

	logger.Info("valid instrument updated")

	return nil
}

//go:embed queries/valid_instruments/update_last_indexed_at.sql
var updateValidInstrumentLastIndexedAtQuery string

// MarkValidInstrumentAsIndexed updates a particular valid instrument's last_indexed_at value.
func (q *Querier) MarkValidInstrumentAsIndexed(ctx context.Context, validInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	args := []any{
		validInstrumentID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid instrument last_indexed_at", updateValidInstrumentLastIndexedAtQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid instrument as indexed")
	}

	logger.Info("valid instrument marked as indexed")

	return nil
}

//go:embed queries/valid_instruments/archive.sql
var archiveValidInstrumentQuery string

// ArchiveValidInstrument archives a valid instrument from the database by its ID.
func (q *Querier) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	args := []any{
		validInstrumentID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid instrument archive", archiveValidInstrumentQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid instrument")
	}

	logger.Info("valid instrument archived")

	return nil
}
