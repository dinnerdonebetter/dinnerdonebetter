package postgres

import (
	"context"
	"database/sql"
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
		"valid_ingredient_group_members.id",
		"valid_ingredient_group_members.belongs_to_group",
		"valid_ingredient_group_members.valid_ingredient",
		"valid_ingredient_group_members.created_at",
		"valid_ingredient_group_members.archived_at",
	}
)

// scanValidIngredientGroup is a consistent way to turn a *sql.Row into a webhook struct.
func (q *Querier) scanValidIngredientGroup(ctx context.Context, rows database.ResultIterator) (validIngredientGroup *types.ValidIngredientGroup, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	validIngredientGroup = &types.ValidIngredientGroup{}

	for rows.Next() {
		groupMember := &types.ValidIngredientGroupMember{}

		targetVars := []any{
			&validIngredientGroup.ID,
			&validIngredientGroup.Name,
			&validIngredientGroup.Description,
			&validIngredientGroup.Slug,
			&validIngredientGroup.CreatedAt,
			&validIngredientGroup.LastUpdatedAt,
			&validIngredientGroup.ArchivedAt,
			&groupMember.ID,
			&groupMember.BelongsToGroup,
			&groupMember.ValidIngredientID,
			&groupMember.CreatedAt,
			&groupMember.ArchivedAt,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, observability.PrepareError(err, span, "scanning validIngredientGroup")
		}

		validIngredientGroup.Members = append(validIngredientGroup.Members, groupMember)
	}

	if err = rows.Err(); err != nil {
		return nil, observability.PrepareError(err, span, "fetching validIngredientGroup from database")
	}

	if validIngredientGroup.ID == "" {
		return nil, sql.ErrNoRows
	}

	return validIngredientGroup, nil
}

// scanValidIngredientGroups provides a consistent way to turn sql rows into a slice of webhooks.
func (q *Querier) scanValidIngredientGroups(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validIngredientGroups []*types.ValidIngredientGroup, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x := &types.ValidIngredientGroup{}
	for rows.Next() {
		validIngredientGroup := &types.ValidIngredientGroup{}
		groupMember := &types.ValidIngredientGroupMember{}

		var (
			lastUpdatedAt,
			archivedAt sql.NullTime
		)

		targetVars := []any{
			&validIngredientGroup.ID,
			&validIngredientGroup.Name,
			&validIngredientGroup.Description,
			&validIngredientGroup.Slug,
			&validIngredientGroup.CreatedAt,
			&validIngredientGroup.LastUpdatedAt,
			&validIngredientGroup.ArchivedAt,
			&groupMember.ID,
			&groupMember.BelongsToGroup,
			&groupMember.ValidIngredientID,
			&groupMember.CreatedAt,
			&groupMember.ArchivedAt,
		}

		if includeCounts {
			targetVars = append(targetVars, &filteredCount, &totalCount)
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, 0, 0, observability.PrepareError(err, span, "scanning validIngredientGroup")
		}

		if lastUpdatedAt.Valid {
			validIngredientGroup.LastUpdatedAt = &lastUpdatedAt.Time
		}
		if archivedAt.Valid {
			validIngredientGroup.ArchivedAt = &archivedAt.Time
		}

		if x.ID == "" {
			events := x.Members
			x = validIngredientGroup
			x.Members = events
		}

		if x.ID != validIngredientGroup.ID {
			validIngredientGroups = append(validIngredientGroups, x)
			x = validIngredientGroup
		}

		x.Members = append(x.Members, groupMember)
	}

	if x.ID != "" {
		validIngredientGroups = append(validIngredientGroups, x)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "fetching webhook from database")
	}

	if len(validIngredientGroups) == 0 {
		return nil, 0, 0, sql.ErrNoRows
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

	rows, err := q.getRows(ctx, q.db, "valid ingredient group", getValidIngredientGroupQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredients groups from database")
	}

	validIngredientGroup, err := q.scanValidIngredientGroup(ctx, rows)
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

//go:embed queries/valid_ingredient_groups/get_many.sql
var getValidIngredientGroupsQuery string

// GetValidIngredientGroups fetches a list of valid ingredients group from the database that meet a particular filter.
func (q *Querier) GetValidIngredientGroups(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientGroup], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	x = &types.QueryFilteredResult[types.ValidIngredientGroup]{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter.Page != nil {
		x.Page = *filter.Page
	}

	if filter.Limit != nil {
		x.Limit = *filter.Limit
	}

	args := []any{
		filter.QueryOffset(),
		filter.CreatedAfter,
		filter.CreatedBefore,
		filter.UpdatedAfter,
		filter.UpdatedBefore,
	}

	rows, err := q.getRows(ctx, q.db, "valid ingredients group", getValidIngredientGroupsQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhook from database")
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

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "starting transaction")
	}

	args := []any{
		input.ID,
		input.Name,
		input.Description,
		input.Slug,
	}

	// create the valid ingredient group.
	if err = q.performWriteQuery(ctx, tx, "valid ingredient group creation", validIngredientGroupCreationQuery, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient group creation query")
	}

	x := &types.ValidIngredientGroup{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
		CreatedAt:   q.currentTime(),
	}

	for i := range input.Members {
		m := input.Members[i]
		var member *types.ValidIngredientGroupMember
		member, err = q.CreateValidIngredientGroupMember(ctx, tx, x.ID, m)
		if err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient group member")
		}

		x.Members = append(x.Members, member)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	tracing.AttachValidIngredientGroupIDToSpan(span, x.ID)
	logger.WithValue("member_count", len(input.Members)).Info("valid ingredient group created")

	return x, nil
}

//go:embed queries/valid_ingredient_groups/create_group_member.sql
var validIngredientGroupMemberCreationQuery string

// CreateValidIngredientGroupMember creates a valid ingredient group member in the database.
func (q *Querier) CreateValidIngredientGroupMember(ctx context.Context, db database.SQLQueryExecutorAndTransactionManager, groupID string, input *types.ValidIngredientGroupMemberDatabaseCreationInput) (*types.ValidIngredientGroupMember, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidIngredientGroupIDKey, input.ID)

	args := []any{
		input.ID,
		groupID,
		input.ValidIngredientID,
	}

	// create the valid ingredient group.
	if err := q.performWriteQuery(ctx, db, "valid ingredient group member creation", validIngredientGroupMemberCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient group member creation query")
	}

	x := &types.ValidIngredientGroupMember{
		ID:                input.ID,
		BelongsToGroup:    groupID,
		ValidIngredientID: input.ValidIngredientID,
		CreatedAt:         q.currentTime(),
	}

	tracing.AttachValidIngredientGroupIDToSpan(span, x.ID)
	logger.Info("valid ingredient group member created")

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
