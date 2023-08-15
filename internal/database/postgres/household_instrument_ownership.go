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
)

var (
	_ types.HouseholdInstrumentOwnershipDataManager = (*Querier)(nil)

	// householdInstrumentOwnershipsTableColumns are the columns for the household_instrument_ownerships table.
	householdInstrumentOwnershipsTableColumns = []string{
		"household_instrument_ownerships.id",
		"household_instrument_ownerships.notes",
		"household_instrument_ownerships.quantity",
		"valid_instruments.id",
		"valid_instruments.name",
		"valid_instruments.plural_name",
		"valid_instruments.description",
		"valid_instruments.icon_path",
		"valid_instruments.usable_for_storage",
		"valid_instruments.display_in_summary_lists",
		"valid_instruments.include_in_generated_instructions",
		"valid_instruments.slug",
		"valid_instruments.created_at",
		"valid_instruments.last_updated_at",
		"valid_instruments.archived_at",
		"household_instrument_ownerships.belongs_to_household",
		"household_instrument_ownerships.created_at",
		"household_instrument_ownerships.last_updated_at",
		"household_instrument_ownerships.archived_at",
	}
)

// scanHouseholdInstrumentOwnership takes a database Scanner (i.e. *sql.Row) and scans the result into a household instrument ownership struct.
func (q *Querier) scanHouseholdInstrumentOwnership(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.HouseholdInstrumentOwnership, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.HouseholdInstrumentOwnership{}

	targetVars := []any{
		&x.ID,
		&x.Notes,
		&x.Quantity,
		&x.Instrument.ID,
		&x.Instrument.Name,
		&x.Instrument.PluralName,
		&x.Instrument.Description,
		&x.Instrument.IconPath,
		&x.Instrument.UsableForStorage,
		&x.Instrument.DisplayInSummaryLists,
		&x.Instrument.IncludeInGeneratedInstructions,
		&x.Instrument.Slug,
		&x.Instrument.CreatedAt,
		&x.Instrument.LastUpdatedAt,
		&x.Instrument.ArchivedAt,
		&x.BelongsToHousehold,
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

// scanHouseholdInstrumentOwnerships takes some database rows and turns them into a slice of household instrument ownerships.
func (q *Querier) scanHouseholdInstrumentOwnerships(ctx context.Context, rows database.ResultIterator, includeCounts bool) (householdInstrumentOwnerships []*types.HouseholdInstrumentOwnership, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanHouseholdInstrumentOwnership(ctx, rows, includeCounts)
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

		householdInstrumentOwnerships = append(householdInstrumentOwnerships, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return householdInstrumentOwnerships, filteredCount, totalCount, nil
}

//go:embed queries/household_instrument_ownerships/exists.sql
var householdInstrumentOwnershipExistenceQuery string

// HouseholdInstrumentOwnershipExists fetches whether a household instrument ownership exists from the database.
func (q *Querier) HouseholdInstrumentOwnershipExists(ctx context.Context, householdInstrumentOwnershipID, householdID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdInstrumentOwnershipID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	tracing.AttachHouseholdInstrumentOwnershipIDToSpan(span, householdInstrumentOwnershipID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []any{
		householdInstrumentOwnershipID,
		householdID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, householdInstrumentOwnershipExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing household instrument ownership existence check")
	}

	return result, nil
}

//go:embed queries/household_instrument_ownerships/get_one.sql
var getHouseholdInstrumentOwnershipQuery string

// GetHouseholdInstrumentOwnership fetches a household instrument ownership from the database.
func (q *Querier) GetHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) (*types.HouseholdInstrumentOwnership, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdInstrumentOwnershipID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	tracing.AttachHouseholdInstrumentOwnershipIDToSpan(span, householdInstrumentOwnershipID)

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []any{
		householdInstrumentOwnershipID,
		householdID,
	}

	row := q.getOneRow(ctx, q.db, "householdInstrumentOwnership", getHouseholdInstrumentOwnershipQuery, args)

	householdInstrumentOwnership, _, _, err := q.scanHouseholdInstrumentOwnership(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning householdInstrumentOwnership")
	}

	return householdInstrumentOwnership, nil
}

// GetHouseholdInstrumentOwnerships fetches a list of household instrument ownerships from the database that meet a particular filter.
func (q *Querier) GetHouseholdInstrumentOwnerships(ctx context.Context, householdID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.HouseholdInstrumentOwnership], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.HouseholdInstrumentOwnership]{
		Pagination: filter.ToPagination(),
	}

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	query, args := q.buildListQuery(ctx, "household_instrument_ownerships", []string{"valid_instruments ON valid_instruments.id = household_instrument_ownerships.valid_instrument_id"}, []string{"household_instrument_ownerships.id", "valid_instruments.id"}, nil, householdOwnershipColumn, householdInstrumentOwnershipsTableColumns, householdID, false, filter)

	rows, err := q.getRows(ctx, q.db, "household instrument ownerships", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing household instrument ownerships list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanHouseholdInstrumentOwnerships(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning household instrument ownerships")
	}

	return x, nil
}

//go:embed queries/household_instrument_ownerships/create.sql
var householdInstrumentOwnershipCreationQuery string

// CreateHouseholdInstrumentOwnership creates a household instrument ownership in the database.
func (q *Querier) CreateHouseholdInstrumentOwnership(ctx context.Context, input *types.HouseholdInstrumentOwnershipDatabaseCreationInput) (*types.HouseholdInstrumentOwnership, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, input.ID)

	args := []any{
		input.ID,
		input.Notes,
		input.Quantity,
		input.ValidInstrumentID,
		input.BelongsToHousehold,
	}

	// create the household instrument ownership.
	if err := q.performWriteQuery(ctx, q.db, "household instrument ownership creation", householdInstrumentOwnershipCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing household instrument ownership creation query")
	}

	x := &types.HouseholdInstrumentOwnership{
		ID:                 input.ID,
		Notes:              input.Notes,
		Quantity:           input.Quantity,
		Instrument:         types.ValidInstrument{ID: input.ValidInstrumentID},
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedAt:          q.currentTime(),
	}

	tracing.AttachHouseholdInstrumentOwnershipIDToSpan(span, x.ID)
	logger.Info("household instrument ownership created")

	return x, nil
}

//go:embed queries/household_instrument_ownerships/update.sql
var updateHouseholdInstrumentOwnershipQuery string

// UpdateHouseholdInstrumentOwnership updates a particular household instrument ownership.
func (q *Querier) UpdateHouseholdInstrumentOwnership(ctx context.Context, updated *types.HouseholdInstrumentOwnership) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, updated.ID)
	tracing.AttachHouseholdInstrumentOwnershipIDToSpan(span, updated.ID)

	args := []any{
		updated.Notes,
		updated.Quantity,
		updated.Instrument.ID,
		updated.ID,
		updated.BelongsToHousehold,
	}

	if err := q.performWriteQuery(ctx, q.db, "household instrument ownership update", updateHouseholdInstrumentOwnershipQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating household instrument ownership")
	}

	logger.Info("household instrument ownership updated")

	return nil
}

// ArchiveHouseholdInstrumentOwnership archives a household instrument ownership from the database by its ID.
func (q *Querier) ArchiveHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdInstrumentOwnershipID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	tracing.AttachHouseholdInstrumentOwnershipIDToSpan(span, householdInstrumentOwnershipID)

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	if err := q.generatedQuerier.ArchiveHouseholdInstrumentOwnership(ctx, q.db, &generated.ArchiveHouseholdInstrumentOwnershipParams{
		ID:                 householdInstrumentOwnershipID,
		BelongsToHousehold: householdID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving household instrument ownership")
	}

	logger.Info("household instrument ownership archived")

	return nil
}
