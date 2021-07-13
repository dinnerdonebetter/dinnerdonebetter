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

var _ querybuilding.RecipeStepProductSQLQueryBuilder = (*Postgres)(nil)

// BuildRecipeStepProductExistsQuery constructs a SQL query for checking if a recipe step product with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildRecipeStepProductExistsQuery(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.RecipeStepProductsTableName).
			Join(querybuilding.RecipeStepsOnRecipeStepProductsJoinClause).
			Join(querybuilding.RecipesOnRecipeStepsJoinClause).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):                                                    recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn):                                            nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn):                                                recipeStepID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn):                                        nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.RecipeStepsTableBelongsToRecipeColumn):                   recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.IDColumn):                                         recipeStepProductID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.ArchivedOnColumn):                                 nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.RecipeStepProductsTableBelongsToRecipeStepColumn): recipeStepID,
			}),
	)
}

// BuildGetRecipeStepProductQuery constructs a SQL query for fetching a recipe step product with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetRecipeStepProductQuery(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipeStepProductsTableColumns...).
			From(querybuilding.RecipeStepProductsTableName).
			Join(querybuilding.RecipeStepsOnRecipeStepProductsJoinClause).
			Join(querybuilding.RecipesOnRecipeStepsJoinClause).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):                                                    recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn):                                            nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn):                                                recipeStepID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn):                                        nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.RecipeStepsTableBelongsToRecipeColumn):                   recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.IDColumn):                                         recipeStepProductID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.ArchivedOnColumn):                                 nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.RecipeStepProductsTableBelongsToRecipeStepColumn): recipeStepID,
			}),
	)
}

// BuildGetAllRecipeStepProductsCountQuery returns a query that fetches the total number of recipe step products in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllRecipeStepProductsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.RecipeStepProductsTableName)).
			From(querybuilding.RecipeStepProductsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfRecipeStepProductsQuery returns a query that fetches every recipe step product in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfRecipeStepProductsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipeStepProductsTableColumns...).
			From(querybuilding.RecipeStepProductsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetRecipeStepProductsQuery builds a SQL query selecting recipe step products that adhere to a given QueryFilter and belong to a given account,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetRecipeStepProductsQuery(ctx context.Context, recipeID, recipeStepID uint64, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	joins := []string{
		querybuilding.RecipeStepsOnRecipeStepProductsJoinClause,
		querybuilding.RecipesOnRecipeStepsJoinClause,
	}
	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):                                  recipeID,
		fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn):                          nil,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn):                              recipeStepID,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn):                      nil,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.RecipeStepsTableBelongsToRecipeColumn): recipeID,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.ArchivedOnColumn):               nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.RecipeStepProductsTableName,
		joins,
		where,
		querybuilding.RecipeStepProductsTableBelongsToRecipeStepColumn,
		querybuilding.RecipeStepProductsTableColumns,
		recipeStepID,
		includeArchived,
		filter,
	)
}

// BuildGetRecipeStepProductsWithIDsQuery builds a SQL query selecting recipe step products that belong to a given account,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetRecipeStepProductsWithIDsQuery(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.IDColumn):                                         ids,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.ArchivedOnColumn):                                 nil,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepProductsTableName, querybuilding.RecipeStepProductsTableBelongsToRecipeStepColumn): recipeStepID,
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.RecipeStepProductsTableColumns...).
		From(querybuilding.RecipeStepProductsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipeStepProductsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.RecipeStepProductsTableName).
			Where(where),
	)
}

// BuildCreateRecipeStepProductQuery takes a recipe step product and returns a creation query for that recipe step product and the relevant arguments.
func (b *Postgres) BuildCreateRecipeStepProductQuery(ctx context.Context, input *types.RecipeStepProductCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.RecipeStepProductsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.RecipeStepProductsTableNameColumn,
				querybuilding.RecipeStepProductsTableQuantityTypeColumn,
				querybuilding.RecipeStepProductsTableQuantityValueColumn,
				querybuilding.RecipeStepProductsTableQuantityNotesColumn,
				querybuilding.RecipeStepProductsTableRecipeStepIDColumn,
				querybuilding.RecipeStepProductsTableBelongsToRecipeStepColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.Name,
				input.QuantityType,
				input.QuantityValue,
				input.QuantityNotes,
				input.RecipeStepID,
				input.BelongsToRecipeStep,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateRecipeStepProductQuery takes a recipe step product and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateRecipeStepProductQuery(ctx context.Context, input *types.RecipeStepProduct) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, input.BelongsToRecipeStep)
	tracing.AttachRecipeStepProductIDToSpan(span, input.ID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.RecipeStepProductsTableName).
			Set(querybuilding.RecipeStepProductsTableNameColumn, input.Name).
			Set(querybuilding.RecipeStepProductsTableQuantityTypeColumn, input.QuantityType).
			Set(querybuilding.RecipeStepProductsTableQuantityValueColumn, input.QuantityValue).
			Set(querybuilding.RecipeStepProductsTableQuantityNotesColumn, input.QuantityNotes).
			Set(querybuilding.RecipeStepProductsTableRecipeStepIDColumn, input.RecipeStepID).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                                         input.ID,
				querybuilding.ArchivedOnColumn:                                 nil,
				querybuilding.RecipeStepProductsTableBelongsToRecipeStepColumn: input.BelongsToRecipeStep,
			}),
	)
}

// BuildArchiveRecipeStepProductQuery returns a SQL query which marks a given recipe step product belonging to a given account as archived.
func (b *Postgres) BuildArchiveRecipeStepProductQuery(ctx context.Context, recipeStepID, recipeStepProductID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.RecipeStepProductsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                                         recipeStepProductID,
				querybuilding.ArchivedOnColumn:                                 nil,
				querybuilding.RecipeStepProductsTableBelongsToRecipeStepColumn: recipeStepID,
			}),
	)
}

// BuildGetAuditLogEntriesForRecipeStepProductQuery constructs a SQL query for fetching audit log entries relating to a recipe step product with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForRecipeStepProductQuery(ctx context.Context, recipeStepProductID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	recipeStepProductIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.RecipeStepProductAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{recipeStepProductIDKey: recipeStepProductID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}
