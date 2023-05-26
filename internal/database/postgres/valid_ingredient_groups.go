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
	_ types.ValidIngredientGroupDataManager = (*Querier)(nil)

	// validIngredientGroupsTableColumns are the columns for the valid_ingredient_groups table.
	validIngredientGroupsTableColumns = []string{
		"valid_ingredient_groups.id",
		"valid_ingredient_groups.name",
		"valid_ingredient_groups.description",
		"valid_ingredient_groups.slug",
		"valid_ingredient_groups.created_at",
		"valid_ingredient_groups.last_updated_at",
		"valid_ingredient_groups.archived_at",
	}
)

// scanValidIngredientGroup takes a database Scanner (i.e. *sql.Row) and scans the result into a valid ingredient group struct.
func (q *Querier) scanValidIngredientGroup(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidIngredientGroup, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ValidIngredientGroup{}

	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Description,
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

// scanValidIngredientGroups takes some database rows and turns them into a slice of valid ingredients group.
func (q *Querier) scanValidIngredientGroups(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredientGroups []*types.ValidIngredientGroup, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidIngredientGroup(ctx, rows, includeCounts)
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

		validIngredientGroups = append(validIngredientGroups, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validIngredientGroups, filteredCount, totalCount, nil
}

//go:embed queries/valid_ingredient_groups/exists.sql
var validIngredientGroupExistenceQuery string

// ValidIngredientGroupExists fetches whether a valid ingredient group exists from the database.
func (q *Querier) ValidIngredientGroupExists(ctx context.Context, validIngredientGroupID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientGroupID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachValidIngredientGroupIDToSpan(span, validIngredientGroupID)

	args := []any{
		validIngredientGroupID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validIngredientGroupExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient group existence check")
	}

	return result, nil
}

//go:embed queries/valid_ingredient_groups/get_one.sql
var getValidIngredientGroupQuery string

// GetValidIngredientGroup fetches a valid ingredient group from the database.
func (q *Querier) GetValidIngredientGroup(ctx context.Context, validIngredientGroupID string) (*types.ValidIngredientGroup, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientGroupID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachValidIngredientGroupIDToSpan(span, validIngredientGroupID)

	args := []any{
		validIngredientGroupID,
	}

	row := q.getOneRow(ctx, q.db, "valid ingredient group", getValidIngredientGroupQuery, args)

	validIngredientGroup, _, _, err := q.scanValidIngredientGroup(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredient group")
	}

	return validIngredientGroup, nil
}

//go:embed queries/valid_ingredient_groups/search.sql
var validIngredientGroupSearchQuery string

// SearchForValidIngredientGroups fetches a valid ingredient group from the database.
func (q *Querier) SearchForValidIngredientGroups(ctx context.Context, query string, filter *types.QueryFilter) ([]*types.ValidIngredientGroup, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidIngredientGroupIDToSpan(span, query)

	args := []any{
		wrapQueryForILIKE(query),
	}

	rows, err := q.getRows(ctx, q.db, "valid ingredients group", validIngredientGroupSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients group list retrieval query")
	}

	x, _, _, err := q.scanValidIngredientGroups(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredients group")
	}

	return x, nil
}

// GetValidIngredientGroups fetches a list of valid ingredients group from the database that meet a particular filter.
func (q *Querier) GetValidIngredientGroups(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientGroup], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.ValidIngredientGroup]{}
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

	query, args := q.buildListQuery(ctx, "valid_ingredient_groups", nil, nil, nil, householdOwnershipColumn, validIngredientGroupsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "valid ingredients group", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients group list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidIngredientGroups(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid ingredients group")
	}

	return x, nil
}

//go:embed queries/valid_ingredient_groups/create.sql
var validIngredientGroupCreationQuery string

// CreateValidIngredientGroup creates a valid ingredient group in the database.
func (q *Querier) CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupDatabaseCreationInput) (*types.ValidIngredientGroup, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientGroupIDKey, input.ID)

	args := []any{
		input.ID,
		input.Name,
		input.Description,
		input.Slug,
	}

	// create the valid ingredient group.
	if err := q.performWriteQuery(ctx, q.db, "valid ingredient group creation", validIngredientGroupCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient group creation query")
	}

	x := &types.ValidIngredientGroup{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
		CreatedAt:   q.currentTime(),
	}

	tracing.AttachValidIngredientGroupIDToSpan(span, x.ID)
	logger.Info("valid ingredient group created")

	return x, nil
}

//go:embed queries/valid_ingredient_groups/update.sql
var updateValidIngredientGroupQuery string

// UpdateValidIngredientGroup updates a particular valid ingredient group.
func (q *Querier) UpdateValidIngredientGroup(ctx context.Context, updated *types.ValidIngredientGroup) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientGroupIDKey, updated.ID)
	tracing.AttachValidIngredientGroupIDToSpan(span, updated.ID)

	args := []any{
		updated.Name,
		updated.Description,
		updated.Slug,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient group update", updateValidIngredientGroupQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient group")
	}

	logger.Info("valid ingredient group updated")

	return nil
}

//go:embed queries/valid_ingredient_groups/archive.sql
var archiveValidIngredientGroupQuery string

// ArchiveValidIngredientGroup archives a valid ingredient group from the database by its ID.
func (q *Querier) ArchiveValidIngredientGroup(ctx context.Context, validIngredientGroupID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientGroupID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachValidIngredientGroupIDToSpan(span, validIngredientGroupID)

	args := []any{
		validIngredientGroupID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid ingredient group archive", archiveValidIngredientGroupQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient group")
	}

	logger.Info("valid ingredient group archived")

	return nil
}
