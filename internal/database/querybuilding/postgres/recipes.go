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

var _ querybuilding.RecipeSQLQueryBuilder = (*Postgres)(nil)

// BuildRecipeExistsQuery constructs a SQL query for checking if a recipe with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildRecipeExistsQuery(ctx context.Context, recipeID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.RecipesTableName).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):         recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetRecipeQuery constructs a SQL query for fetching a recipe with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetRecipeQuery(ctx context.Context, recipeID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipesTableColumns...).
			From(querybuilding.RecipesTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):         recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetAllRecipesCountQuery returns a query that fetches the total number of recipes in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllRecipesCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.RecipesTableName)).
			From(querybuilding.RecipesTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfRecipesQuery returns a query that fetches every recipe in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfRecipesQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipesTableColumns...).
			From(querybuilding.RecipesTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetRecipesQuery builds a SQL query selecting recipes that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetRecipesQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn): nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.RecipesTableName,
		nil,
		where,
		"",
		querybuilding.RecipesTableColumns,
		0,
		includeArchived,
		filter,
	)
}

// BuildGetRecipesWithIDsQuery builds a SQL query selecting recipes that belong to a given household,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetRecipesWithIDsQuery(ctx context.Context, householdID uint64, limit uint8, ids []uint64, restrictToHousehold bool) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):         ids,
		fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn): nil,
	}

	if restrictToHousehold {
		where[fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.RecipesTableHouseholdOwnershipColumn)] = householdID
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.RecipesTableColumns...).
		From(querybuilding.RecipesTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipesTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.RecipesTableName).
			Where(where),
	)
}

// BuildCreateRecipeQuery takes a recipe and returns a creation query for that recipe and the relevant arguments.
func (b *Postgres) BuildCreateRecipeQuery(ctx context.Context, input *types.RecipeCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.RecipesTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.RecipesTableNameColumn,
				querybuilding.RecipesTableSourceColumn,
				querybuilding.RecipesTableDescriptionColumn,
				querybuilding.RecipesTableInspiredByRecipeIDColumn,
				querybuilding.RecipesTableHouseholdOwnershipColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.Name,
				input.Source,
				input.Description,
				input.InspiredByRecipeID,
				input.BelongsToHousehold,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateRecipeQuery takes a recipe and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateRecipeQuery(ctx context.Context, input *types.Recipe) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, input.ID)
	tracing.AttachHouseholdIDToSpan(span, input.BelongsToHousehold)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.RecipesTableName).
			Set(querybuilding.RecipesTableNameColumn, input.Name).
			Set(querybuilding.RecipesTableSourceColumn, input.Source).
			Set(querybuilding.RecipesTableDescriptionColumn, input.Description).
			Set(querybuilding.RecipesTableInspiredByRecipeIDColumn, input.InspiredByRecipeID).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                             input.ID,
				querybuilding.ArchivedOnColumn:                     nil,
				querybuilding.RecipesTableHouseholdOwnershipColumn: input.BelongsToHousehold,
			}),
	)
}

// BuildArchiveRecipeQuery returns a SQL query which marks a given recipe belonging to a given household as archived.
func (b *Postgres) BuildArchiveRecipeQuery(ctx context.Context, recipeID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.RecipesTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         recipeID,
				querybuilding.ArchivedOnColumn: nil,
			}),
	)
}

// BuildGetAuditLogEntriesForRecipeQuery constructs a SQL query for fetching audit log entries relating to a recipe with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForRecipeQuery(ctx context.Context, recipeID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)

	recipeIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.RecipeAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{recipeIDKey: recipeID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}
