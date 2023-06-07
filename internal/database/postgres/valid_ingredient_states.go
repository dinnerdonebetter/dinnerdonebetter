package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.ValidIngredientStateDataManager = (*Querier)(nil)

	// validIngredientStatesTableColumns are the columns for the valid_ingredient_states table.
	validIngredientStatesTableColumns = []string{
		"valid_ingredient_states.id",
		"valid_ingredient_states.name",
		"valid_ingredient_states.description",
		"valid_ingredient_states.icon_path",
		"valid_ingredient_states.slug",
		"valid_ingredient_states.past_tense",
		"valid_ingredient_states.attribute_type",
		"valid_ingredient_states.created_at",
		"valid_ingredient_states.last_updated_at",
		"valid_ingredient_states.archived_at",
	}
)

// scanValidIngredientState takes a database Scanner (i.e. *sql.Row) and scans the result into a valid preparation struct.
func (q *Querier) scanValidIngredientState(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredientState, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidIngredientState{}

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Description,
		&x.IconPath,
		&x.Slug,
		&x.PastTense,
		&x.AttributeType,
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

// scanValidIngredientStates takes some database rows and turns them into a slice of valid preparations.
func (q *Querier) scanValidIngredientStates(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredientStates []*types.ValidIngredientState, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidIngredientState(ctx, rows, includeCounts)
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

		validIngredientStates = append(validIngredientStates, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validIngredientStates, filteredCount, totalCount, nil
}

//go:embed queries/valid_ingredient_states/exists.sql
var validIngredientStateExistenceQuery string

// ValidIngredientStateExists fetches whether a valid preparation exists from the database.
func (q *Querier) ValidIngredientStateExists(ctx context.Context, validIngredientStateID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachValidIngredientStateIDToSpan(span, validIngredientStateID)

	args := []any{
		validIngredientStateID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validIngredientStateExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid preparation existence check")
	}

	return result, nil
}

//go:embed queries/valid_ingredient_states/get_one.sql
var getValidIngredientStateQuery string

// GetValidIngredientState fetches a valid preparation from the database.
func (q *Querier) GetValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachValidIngredientStateIDToSpan(span, validIngredientStateID)

	args := []any{
		validIngredientStateID,
	}

	row := q.getOneRow(ctx, q.db, "validIngredientState", getValidIngredientStateQuery, args)

	validIngredientState, _, _, err := q.scanValidIngredientState(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning validIngredientState")
	}

	return validIngredientState, nil
}

//go:embed queries/valid_ingredient_states/search.sql
var validIngredientStateSearchQuery string

// SearchForValidIngredientStates fetches a valid preparation from the database.
func (q *Querier) SearchForValidIngredientStates(ctx context.Context, query string) ([]*types.ValidIngredientState, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidIngredientStateIDToSpan(span, query)

	args := []any{
		wrapQueryForILIKE(query),
	}

	rows, err := q.getRows(ctx, q.db, "valid preparations", validIngredientStateSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparations list retrieval query")
	}

	x, _, _, err := q.scanValidIngredientStates(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid preparations")
	}

	return x, nil
}

// GetValidIngredientStates fetches a list of valid preparations from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStates(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientState], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.ValidIngredientState]{}
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

	query, args := q.buildListQuery(ctx, "valid_ingredient_states", nil, nil, nil, householdOwnershipColumn, validIngredientStatesTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "valid ingredient states", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparations list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientStates(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid preparations")
	}

	return x, nil
}

//go:embed queries/valid_ingredient_states/create.sql
var validIngredientStateCreationQuery string

// CreateValidIngredientState creates a valid preparation in the database.
func (q *Querier) CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateDatabaseCreationInput) (*types.ValidIngredientState, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientStateIDKey, input.ID)

	args := []any{
		input.ID,
		input.Name,
		input.Description,
		input.IconPath,
		input.PastTense,
		input.Slug,
		input.AttributeType,
	}

	// create the valid preparation.
	if err := q.performWriteQuery(ctx, q.db, "valid preparation creation", validIngredientStateCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid preparation creation query")
	}

	x := &types.ValidIngredientState{
		ID:            input.ID,
		Name:          input.Name,
		Description:   input.Description,
		IconPath:      input.IconPath,
		Slug:          input.Slug,
		PastTense:     input.PastTense,
		AttributeType: input.AttributeType,
		CreatedAt:     q.currentTime(),
	}

	tracing.AttachValidIngredientStateIDToSpan(span, x.ID)
	logger.Info("valid preparation created")

	return x, nil
}

//go:embed queries/valid_ingredient_states/update.sql
var updateValidIngredientStateQuery string

// UpdateValidIngredientState updates a particular valid preparation.
func (q *Querier) UpdateValidIngredientState(ctx context.Context, updated *types.ValidIngredientState) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientStateIDKey, updated.ID)
	tracing.AttachValidIngredientStateIDToSpan(span, updated.ID)

	args := []any{
		updated.Name,
		updated.Description,
		updated.IconPath,
		updated.Slug,
		updated.PastTense,
		updated.AttributeType,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation update", updateValidIngredientStateQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation updated")

	return nil
}

//go:embed queries/valid_ingredient_states/update_last_indexed_at.sql
var updateValidIngredientStateLastIndexedAtQuery string

// MarkValidIngredientStateAsIndexed updates a particular valid ingredient state's last_indexed_at value.
func (q *Querier) MarkValidIngredientStateAsIndexed(ctx context.Context, validIngredientStateID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachValidIngredientStateIDToSpan(span, validIngredientStateID)

	args := []any{
		validIngredientStateID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient state last_indexed_at", updateValidIngredientStateLastIndexedAtQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid ingredient state as indexed")
	}

	logger.Info("valid ingredient state marked as indexed")

	return nil
}

//go:embed queries/valid_ingredient_states/archive.sql
var archiveValidIngredientStateQuery string

// ArchiveValidIngredientState archives a valid preparation from the database by its ID.
func (q *Querier) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachValidIngredientStateIDToSpan(span, validIngredientStateID)

	args := []any{
		validIngredientStateID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid preparation archive", archiveValidIngredientStateQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation archived")

	return nil
}
