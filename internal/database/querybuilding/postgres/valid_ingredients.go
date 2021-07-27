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

var _ querybuilding.ValidIngredientSQLQueryBuilder = (*Postgres)(nil)

// BuildValidIngredientExistsQuery constructs a SQL query for checking if a valid ingredient with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildValidIngredientExistsQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.ValidIngredientsTableName).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.IDColumn):         validIngredientID,
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetValidIngredientIDForNameQuery constructs a SQL query for retrieving the ID of a given named valid preparation.
func (b *Postgres) BuildGetValidIngredientIDForNameQuery(ctx context.Context, validIngredientName string) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.IDColumn)).
			From(querybuilding.ValidIngredientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.ValidIngredientsTableNameColumn): validIngredientName,
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.ArchivedOnColumn):                nil,
			}),
	)
}

// BuildSearchForValidIngredientByNameQuery returns a SQL query (and argument) for retrieving a valid ingredient by its name.
func (b *Postgres) BuildSearchForValidIngredientByNameQuery(ctx context.Context, name string) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidIngredientsTableColumns...).
			From(querybuilding.ValidIngredientsTableName).
			Where(squirrel.Expr(
				fmt.Sprintf("%s.%s ILIKE ?", querybuilding.ValidIngredientsTableName, querybuilding.ValidIngredientsTableNameColumn),
				fmt.Sprintf("%s%%", name),
			)).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetValidIngredientQuery constructs a SQL query for fetching a valid ingredient with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetValidIngredientQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidIngredientsTableColumns...).
			From(querybuilding.ValidIngredientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.IDColumn):         validIngredientID,
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetAllValidIngredientsCountQuery returns a query that fetches the total number of valid ingredients in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllValidIngredientsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.ValidIngredientsTableName)).
			From(querybuilding.ValidIngredientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfValidIngredientsQuery returns a query that fetches every valid ingredient in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfValidIngredientsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidIngredientsTableColumns...).
			From(querybuilding.ValidIngredientsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetValidIngredientsQuery builds a SQL query selecting valid ingredients that adhere to a given QueryFilter and belong to a given account,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetValidIngredientsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.ValidIngredientsTableName,
		nil,
		where,
		"",
		querybuilding.ValidIngredientsTableColumns,
		0,
		includeArchived,
		filter,
	)
}

// BuildGetValidIngredientsWithIDsQuery builds a SQL query selecting valid ingredients that belong to a given account,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetValidIngredientsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.IDColumn):         ids,
		fmt.Sprintf("%s.%s", querybuilding.ValidIngredientsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.ValidIngredientsTableColumns...).
		From(querybuilding.ValidIngredientsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.ValidIngredientsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.ValidIngredientsTableName).
			Where(where),
	)
}

// BuildCreateValidIngredientQuery takes a valid ingredient and returns a creation query for that valid ingredient and the relevant arguments.
func (b *Postgres) BuildCreateValidIngredientQuery(ctx context.Context, input *types.ValidIngredientCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.ValidIngredientsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.ValidIngredientsTableNameColumn,
				querybuilding.ValidIngredientsTableVariantColumn,
				querybuilding.ValidIngredientsTableDescriptionColumn,
				querybuilding.ValidIngredientsTableWarningColumn,
				querybuilding.ValidIngredientsTableContainsEggColumn,
				querybuilding.ValidIngredientsTableContainsDairyColumn,
				querybuilding.ValidIngredientsTableContainsPeanutColumn,
				querybuilding.ValidIngredientsTableContainsTreeNutColumn,
				querybuilding.ValidIngredientsTableContainsSoyColumn,
				querybuilding.ValidIngredientsTableContainsWheatColumn,
				querybuilding.ValidIngredientsTableContainsShellfishColumn,
				querybuilding.ValidIngredientsTableContainsSesameColumn,
				querybuilding.ValidIngredientsTableContainsFishColumn,
				querybuilding.ValidIngredientsTableContainsGlutenColumn,
				querybuilding.ValidIngredientsTableAnimalFleshColumn,
				querybuilding.ValidIngredientsTableAnimalDerivedColumn,
				querybuilding.ValidIngredientsTableVolumetricColumn,
				querybuilding.ValidIngredientsTableIconPathColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.Name,
				input.Variant,
				input.Description,
				input.Warning,
				input.ContainsEgg,
				input.ContainsDairy,
				input.ContainsPeanut,
				input.ContainsTreeNut,
				input.ContainsSoy,
				input.ContainsWheat,
				input.ContainsShellfish,
				input.ContainsSesame,
				input.ContainsFish,
				input.ContainsGluten,
				input.AnimalFlesh,
				input.AnimalDerived,
				input.Volumetric,
				input.IconPath,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateValidIngredientQuery takes a valid ingredient and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateValidIngredientQuery(ctx context.Context, input *types.ValidIngredient) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, input.ID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidIngredientsTableName).
			Set(querybuilding.ValidIngredientsTableNameColumn, input.Name).
			Set(querybuilding.ValidIngredientsTableVariantColumn, input.Variant).
			Set(querybuilding.ValidIngredientsTableDescriptionColumn, input.Description).
			Set(querybuilding.ValidIngredientsTableWarningColumn, input.Warning).
			Set(querybuilding.ValidIngredientsTableContainsEggColumn, input.ContainsEgg).
			Set(querybuilding.ValidIngredientsTableContainsDairyColumn, input.ContainsDairy).
			Set(querybuilding.ValidIngredientsTableContainsPeanutColumn, input.ContainsPeanut).
			Set(querybuilding.ValidIngredientsTableContainsTreeNutColumn, input.ContainsTreeNut).
			Set(querybuilding.ValidIngredientsTableContainsSoyColumn, input.ContainsSoy).
			Set(querybuilding.ValidIngredientsTableContainsWheatColumn, input.ContainsWheat).
			Set(querybuilding.ValidIngredientsTableContainsShellfishColumn, input.ContainsShellfish).
			Set(querybuilding.ValidIngredientsTableContainsSesameColumn, input.ContainsSesame).
			Set(querybuilding.ValidIngredientsTableContainsFishColumn, input.ContainsFish).
			Set(querybuilding.ValidIngredientsTableContainsGlutenColumn, input.ContainsGluten).
			Set(querybuilding.ValidIngredientsTableAnimalFleshColumn, input.AnimalFlesh).
			Set(querybuilding.ValidIngredientsTableAnimalDerivedColumn, input.AnimalDerived).
			Set(querybuilding.ValidIngredientsTableVolumetricColumn, input.Volumetric).
			Set(querybuilding.ValidIngredientsTableIconPathColumn, input.IconPath).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         input.ID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildArchiveValidIngredientQuery returns a SQL query which marks a given valid ingredient belonging to a given account as archived.
func (b *Postgres) BuildArchiveValidIngredientQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.ValidIngredientsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         validIngredientID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildGetAuditLogEntriesForValidIngredientQuery constructs a SQL query for fetching audit log entries relating to a valid ingredient with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForValidIngredientQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	validIngredientIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.ValidIngredientAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{validIngredientIDKey: validIngredientID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}
