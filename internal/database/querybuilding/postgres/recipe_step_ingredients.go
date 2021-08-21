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

var _ querybuilding.RecipeStepIngredientSQLQueryBuilder = (*Postgres)(nil)

// BuildRecipeStepIngredientExistsQuery constructs a SQL query for checking if a recipe step ingredient with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildRecipeStepIngredientExistsQuery(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.RecipeStepIngredientsTableName).
			Join(querybuilding.RecipeStepsOnRecipeStepIngredientsJoinClause).
			Join(querybuilding.RecipesOnRecipeStepsJoinClause).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):                                                          recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn):                                                  nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn):                                                      recipeStepID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn):                                              nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.RecipeStepsTableBelongsToRecipeColumn):                         recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.IDColumn):                                            recipeStepIngredientID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.ArchivedOnColumn):                                    nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.RecipeStepIngredientsTableBelongsToRecipeStepColumn): recipeStepID,
			}),
	)
}

// BuildGetRecipeStepIngredientQuery constructs a SQL query for fetching a recipe step ingredient with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetRecipeStepIngredientQuery(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipeStepIngredientsTableColumns...).
			From(querybuilding.RecipeStepIngredientsTableName).
			Join(querybuilding.RecipeStepsOnRecipeStepIngredientsJoinClause).
			Join(querybuilding.RecipesOnRecipeStepsJoinClause).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):                                                          recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn):                                                  nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn):                                                      recipeStepID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn):                                              nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.RecipeStepsTableBelongsToRecipeColumn):                         recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.IDColumn):                                            recipeStepIngredientID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.ArchivedOnColumn):                                    nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.RecipeStepIngredientsTableBelongsToRecipeStepColumn): recipeStepID,
			}),
	)
}

// BuildGetAllRecipeStepIngredientsCountQuery returns a query that fetches the total number of recipe step ingredients in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllRecipeStepIngredientsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.RecipeStepIngredientsTableName)).
			From(querybuilding.RecipeStepIngredientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfRecipeStepIngredientsQuery returns a query that fetches every recipe step ingredient in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfRecipeStepIngredientsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipeStepIngredientsTableColumns...).
			From(querybuilding.RecipeStepIngredientsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetRecipeStepIngredientsQuery builds a SQL query selecting recipe step ingredients that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetRecipeStepIngredientsQuery(ctx context.Context, recipeID, recipeStepID uint64, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	joins := []string{
		querybuilding.RecipeStepsOnRecipeStepIngredientsJoinClause,
		querybuilding.RecipesOnRecipeStepsJoinClause,
	}
	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):                                  recipeID,
		fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn):                          nil,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn):                              recipeStepID,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn):                      nil,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.RecipeStepsTableBelongsToRecipeColumn): recipeID,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.ArchivedOnColumn):            nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.RecipeStepIngredientsTableName,
		joins,
		where,
		querybuilding.RecipeStepIngredientsTableBelongsToRecipeStepColumn,
		querybuilding.RecipeStepIngredientsTableColumns,
		recipeStepID,
		includeArchived,
		filter,
	)
}

// BuildGetRecipeStepIngredientsWithIDsQuery builds a SQL query selecting recipe step ingredients that belong to a given household,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetRecipeStepIngredientsWithIDsQuery(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.IDColumn):                                            ids,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.ArchivedOnColumn):                                    nil,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepIngredientsTableName, querybuilding.RecipeStepIngredientsTableBelongsToRecipeStepColumn): recipeStepID,
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.RecipeStepIngredientsTableColumns...).
		From(querybuilding.RecipeStepIngredientsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipeStepIngredientsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.RecipeStepIngredientsTableName).
			Where(where),
	)
}

// BuildCreateRecipeStepIngredientQuery takes a recipe step ingredient and returns a creation query for that recipe step ingredient and the relevant arguments.
func (b *Postgres) BuildCreateRecipeStepIngredientQuery(ctx context.Context, input *types.RecipeStepIngredientCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.RecipeStepIngredientsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.RecipeStepIngredientsTableIngredientIDColumn,
				querybuilding.RecipeStepIngredientsTableNameColumn,
				querybuilding.RecipeStepIngredientsTableQuantityTypeColumn,
				querybuilding.RecipeStepIngredientsTableQuantityValueColumn,
				querybuilding.RecipeStepIngredientsTableQuantityNotesColumn,
				querybuilding.RecipeStepIngredientsTableProductOfRecipeStepColumn,
				querybuilding.RecipeStepIngredientsTableIngredientNotesColumn,
				querybuilding.RecipeStepIngredientsTableBelongsToRecipeStepColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.IngredientID,
				input.Name,
				input.QuantityType,
				input.QuantityValue,
				input.QuantityNotes,
				input.ProductOfRecipeStep,
				input.IngredientNotes,
				input.BelongsToRecipeStep,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateRecipeStepIngredientQuery takes a recipe step ingredient and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateRecipeStepIngredientQuery(ctx context.Context, input *types.RecipeStepIngredient) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, input.BelongsToRecipeStep)
	tracing.AttachRecipeStepIngredientIDToSpan(span, input.ID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.RecipeStepIngredientsTableName).
			Set(querybuilding.RecipeStepIngredientsTableIngredientIDColumn, input.IngredientID).
			Set(querybuilding.RecipeStepIngredientsTableNameColumn, input.Name).
			Set(querybuilding.RecipeStepIngredientsTableQuantityTypeColumn, input.QuantityType).
			Set(querybuilding.RecipeStepIngredientsTableQuantityValueColumn, input.QuantityValue).
			Set(querybuilding.RecipeStepIngredientsTableQuantityNotesColumn, input.QuantityNotes).
			Set(querybuilding.RecipeStepIngredientsTableProductOfRecipeStepColumn, input.ProductOfRecipeStep).
			Set(querybuilding.RecipeStepIngredientsTableIngredientNotesColumn, input.IngredientNotes).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         input.ID,
				querybuilding.ArchivedOnColumn: nil,
				querybuilding.RecipeStepIngredientsTableBelongsToRecipeStepColumn: input.BelongsToRecipeStep,
			}),
	)
}

// BuildArchiveRecipeStepIngredientQuery returns a SQL query which marks a given recipe step ingredient belonging to a given household as archived.
func (b *Postgres) BuildArchiveRecipeStepIngredientQuery(ctx context.Context, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.RecipeStepIngredientsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:         recipeStepIngredientID,
				querybuilding.ArchivedOnColumn: nil,
				querybuilding.RecipeStepIngredientsTableBelongsToRecipeStepColumn: recipeStepID,
			}),
	)
}

// BuildGetAuditLogEntriesForRecipeStepIngredientQuery constructs a SQL query for fetching audit log entries relating to a recipe step ingredient with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForRecipeStepIngredientQuery(ctx context.Context, recipeStepIngredientID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)

	recipeStepIngredientIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.RecipeStepIngredientAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{recipeStepIngredientIDKey: recipeStepIngredientID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}
