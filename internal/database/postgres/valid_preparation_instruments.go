package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	validInstrumentsOnValidPreparationInstrumentsJoinClause  = "valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id"
	validPreparationsOnValidPreparationInstrumentsJoinClause = "valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id"
)

var (
	_ types.ValidPreparationInstrumentDataManager = (*Querier)(nil)

	// fullValidPreparationInstrumentsTableColumns are the columns for the valid_preparation_instruments table.
	fullValidPreparationInstrumentsTableColumns = []string{
		"valid_preparation_instruments.id",
		"valid_preparation_instruments.notes",
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
		"valid_instruments.id",
		"valid_instruments.name",
		"valid_instruments.plural_name",
		"valid_instruments.description",
		"valid_instruments.icon_path",
		"valid_instruments.usable_for_storage",
		"valid_instruments.created_at",
		"valid_instruments.last_updated_at",
		"valid_instruments.archived_at",
		"valid_preparation_instruments.created_at",
		"valid_preparation_instruments.last_updated_at",
		"valid_preparation_instruments.archived_at",
	}
)

// scanValidPreparationInstrument takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient preparation struct.
func (q *Querier) scanValidPreparationInstrument(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidPreparationInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidPreparationInstrument{}

	targetVars := []interface{}{
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
		&x.Instrument.ID,
		&x.Instrument.Name,
		&x.Instrument.PluralName,
		&x.Instrument.Description,
		&x.Instrument.IconPath,
		&x.Instrument.UsableForStorage,
		&x.Instrument.CreatedAt,
		&x.Instrument.LastUpdatedAt,
		&x.Instrument.ArchivedAt,
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

// scanValidPreparationInstruments takes some database rows and turns them into a slice of valid ingredient preparations.
func (q *Querier) scanValidPreparationInstruments(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validPreparationInstruments []*types.ValidPreparationInstrument, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

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
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validPreparationInstruments, filteredCount, totalCount, nil
}

//go:embed queries/valid_preparation_instruments_exists.sql
var validPreparationInstrumentExistenceQuery string

// ValidPreparationInstrumentExists fetches whether a valid ingredient preparation exists from the database.
func (q *Querier) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	args := []interface{}{
		validPreparationInstrumentID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validPreparationInstrumentExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, span, "performing valid ingredient preparation existence check")
	}

	return result, nil
}

//go:embed queries/valid_preparation_instruments_get_one.sql
var getValidPreparationInstrumentQuery string

// GetValidPreparationInstrument fetches a valid ingredient preparation from the database.
func (q *Querier) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	args := []interface{}{
		validPreparationInstrumentID,
	}

	row := q.getOneRow(ctx, q.db, "validPreparationInstrument", getValidPreparationInstrumentQuery, args)

	validPreparationInstrument, _, _, err := q.scanValidPreparationInstrument(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning validPreparationInstrument")
	}

	return validPreparationInstrument, nil
}

// GetValidPreparationInstruments fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.ValidPreparationInstrumentList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidPreparationInstrumentList{}
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
		validInstrumentsOnValidPreparationInstrumentsJoinClause,
		validPreparationsOnValidPreparationInstrumentsJoinClause,
	}

	groupBys := []string{
		"valid_preparations.id",
		"valid_instruments.id",
		"valid_preparation_instruments.id",
	}

	query, args := q.buildListQuery(ctx, "valid_preparation_instruments", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationInstrumentsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "valid preparation instruments", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid ingredient preparations list retrieval query")
	}

	if x.ValidPreparationInstruments, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationInstruments(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid ingredient preparations")
	}

	return x, nil
}

func (q *Querier) buildGetValidPreparationInstrumentsRestrictedByIDsQuery(ctx context.Context, column string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query, args, err := q.sqlBuilder.Select(fullValidPreparationInstrumentsTableColumns...).
		From("valid_preparation_instruments").
		Join(validInstrumentsOnValidPreparationInstrumentsJoinClause).
		Join(validPreparationsOnValidPreparationInstrumentsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("valid_preparation_instruments.%s", column): ids,
			"valid_preparation_instruments.archived_at":             nil,
		}).
		Limit(uint64(limit)).
		ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

func (q *Querier) buildGetValidPreparationInstrumentsWithPreparationIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	return q.buildGetValidPreparationInstrumentsRestrictedByIDsQuery(ctx, "valid_preparation_id", limit, ids)
}

// GetValidPreparationInstrumentsForPreparation fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) GetValidPreparationInstrumentsForPreparation(ctx context.Context, preparationID string, filter *types.QueryFilter) (x *types.ValidPreparationInstrumentList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationInstrumentIDToSpan(span, preparationID)

	x = &types.ValidPreparationInstrumentList{
		Pagination: types.Pagination{
			Limit: 20,
		},
	}
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
	query, args := q.buildGetValidPreparationInstrumentsWithPreparationIDsQuery(ctx, x.Limit, []string{preparationID})

	rows, err := q.performReadQuery(ctx, q.db, "valid preparation instruments for preparation", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid ingredient preparations list retrieval query")
	}

	if x.ValidPreparationInstruments, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationInstruments(ctx, rows, false); err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid ingredient preparations")
	}

	return x, nil
}

func (q *Querier) buildGetValidPreparationInstrumentsWithInstrumentIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	return q.buildGetValidPreparationInstrumentsRestrictedByIDsQuery(ctx, "valid_instrument_id", limit, ids)
}

// GetValidPreparationInstrumentsForInstrument fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) GetValidPreparationInstrumentsForInstrument(ctx context.Context, instrumentID string, filter *types.QueryFilter) (x *types.ValidPreparationInstrumentList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if instrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationInstrumentIDToSpan(span, instrumentID)

	x = &types.ValidPreparationInstrumentList{
		Pagination: types.Pagination{
			Limit: 20,
		},
	}
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
	query, args := q.buildGetValidPreparationInstrumentsWithInstrumentIDsQuery(ctx, x.Limit, []string{instrumentID})

	rows, err := q.performReadQuery(ctx, q.db, "valid preparation instruments for instrument", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid ingredient preparations list retrieval query")
	}

	if x.ValidPreparationInstruments, x.FilteredCount, x.TotalCount, err = q.scanValidPreparationInstruments(ctx, rows, false); err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid ingredient preparations")
	}

	return x, nil
}

//go:embed queries/valid_preparation_instruments_create.sql
var validPreparationInstrumentCreationQuery string

// CreateValidPreparationInstrument creates a valid ingredient preparation in the database.
func (q *Querier) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentDatabaseCreationInput) (*types.ValidPreparationInstrument, error) {
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

	// create the valid ingredient preparation.
	if err := q.performWriteQuery(ctx, q.db, "valid ingredient preparation creation", validPreparationInstrumentCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient preparation creation query")
	}

	x := &types.ValidPreparationInstrument{
		ID:          input.ID,
		Notes:       input.Notes,
		Preparation: types.ValidPreparation{ID: input.ValidPreparationID},
		Instrument:  types.ValidInstrument{ID: input.ValidInstrumentID},
		CreatedAt:   q.currentTime(),
	}

	tracing.AttachValidPreparationInstrumentIDToSpan(span, x.ID)
	logger.Info("valid ingredient preparation created")

	return x, nil
}

//go:embed queries/valid_preparation_instruments_update.sql
var updateValidPreparationInstrumentQuery string

// UpdateValidPreparationInstrument updates a particular valid ingredient preparation.
func (q *Querier) UpdateValidPreparationInstrument(ctx context.Context, updated *types.ValidPreparationInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, updated.ID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Notes,
		updated.Preparation.ID,
		updated.Instrument.ID,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient preparation update", updateValidPreparationInstrumentQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation")
	}

	logger.Info("valid ingredient preparation updated")

	return nil
}

//go:embed queries/valid_preparation_instruments_archive.sql
var archiveValidPreparationInstrumentQuery string

// ArchiveValidPreparationInstrument archives a valid ingredient preparation from the database by its ID.
func (q *Querier) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)

	args := []interface{}{
		validPreparationInstrumentID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient preparation archive", archiveValidPreparationInstrumentQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation")
	}

	logger.Info("valid ingredient preparation archived")

	return nil
}
