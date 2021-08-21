package postgres

import (
	"context"
	"fmt"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/Masterminds/squirrel"
)

var _ querybuilding.InvitationSQLQueryBuilder = (*Postgres)(nil)

// BuildInvitationExistsQuery constructs a SQL query for checking if an invitation with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildInvitationExistsQuery(ctx context.Context, invitationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachInvitationIDToSpan(span, invitationID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.InvitationsTableName).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.IDColumn):         invitationID,
				fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetInvitationQuery constructs a SQL query for fetching an invitation with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetInvitationQuery(ctx context.Context, invitationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachInvitationIDToSpan(span, invitationID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.InvitationsTableColumns...).
			From(querybuilding.InvitationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.IDColumn):         invitationID,
				fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetAllInvitationsCountQuery returns a query that fetches the total number of invitations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllInvitationsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.InvitationsTableName)).
			From(querybuilding.InvitationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfInvitationsQuery returns a query that fetches every invitation in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfInvitationsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.InvitationsTableColumns...).
			From(querybuilding.InvitationsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetInvitationsQuery builds a SQL query selecting invitations that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetInvitationsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.InvitationsTableName,
		nil,
		where,
		"",
		querybuilding.InvitationsTableColumns,
		0,
		includeArchived,
		filter,
	)
}

// BuildGetInvitationsWithIDsQuery builds a SQL query selecting invitations that belong to a given household,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetInvitationsWithIDsQuery(ctx context.Context, householdID uint64, limit uint8, ids []uint64, restrictToHousehold bool) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.IDColumn):         ids,
		fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	if restrictToHousehold {
		where[fmt.Sprintf("%s.%s", querybuilding.InvitationsTableName, querybuilding.InvitationsTableHouseholdOwnershipColumn)] = householdID
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.InvitationsTableColumns...).
		From(querybuilding.InvitationsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.InvitationsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.InvitationsTableName).
			Where(where),
	)
}

// BuildCreateInvitationQuery takes an invitation and returns a creation query for that invitation and the relevant arguments.
func (b *Postgres) BuildCreateInvitationQuery(ctx context.Context, input *types.InvitationCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.InvitationsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.InvitationsTableCodeColumn,
				querybuilding.InvitationsTableConsumedColumn,
				querybuilding.InvitationsTableHouseholdOwnershipColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.Code,
				input.Consumed,
				input.BelongsToHousehold,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateInvitationQuery takes an invitation and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateInvitationQuery(ctx context.Context, input *types.Invitation) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachInvitationIDToSpan(span, input.ID)
	tracing.AttachHouseholdIDToSpan(span, input.BelongsToHousehold)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.InvitationsTableName).
			Set(querybuilding.InvitationsTableCodeColumn, input.Code).
			Set(querybuilding.InvitationsTableConsumedColumn, input.Consumed).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                                 input.ID,
				querybuilding.ArchivedOnColumn:                         nil,
				querybuilding.InvitationsTableHouseholdOwnershipColumn: input.BelongsToHousehold,
			}),
	)
}

// BuildArchiveInvitationQuery returns a SQL query which marks a given invitation belonging to a given household as archived.
func (b *Postgres) BuildArchiveInvitationQuery(ctx context.Context, invitationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachInvitationIDToSpan(span, invitationID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.InvitationsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         invitationID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildGetAuditLogEntriesForInvitationQuery constructs a SQL query for fetching audit log entries relating to an invitation with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForInvitationQuery(ctx context.Context, invitationID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachInvitationIDToSpan(span, invitationID)

	invitationIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.InvitationAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{invitationIDKey: invitationID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}
