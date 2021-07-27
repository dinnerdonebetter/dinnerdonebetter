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

var _ querybuilding.RecipeStepSQLQueryBuilder = (*Postgres)(nil)

// BuildRecipeStepExistsQuery constructs a SQL query for checking if a recipe step with a given ID belong to a user with a given ID exists.
func (b *Postgres) BuildRecipeStepExistsQuery(ctx context.Context, recipeID, recipeStepID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn)).
			Prefix(querybuilding.ExistencePrefix).
			From(querybuilding.RecipeStepsTableName).
			Join(querybuilding.RecipesOnRecipeStepsJoinClause).
			Suffix(querybuilding.ExistenceSuffix).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):                                  recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn):                          nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn):                              recipeStepID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn):                      nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.RecipeStepsTableBelongsToRecipeColumn): recipeID,
			}),
	)
}

// BuildGetRecipeStepQuery constructs a SQL query for fetching a recipe step with a given ID belong to a user with a given ID.
func (b *Postgres) BuildGetRecipeStepQuery(ctx context.Context, recipeID, recipeStepID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipeStepsTableColumns...).
			From(querybuilding.RecipeStepsTableName).
			Join(querybuilding.RecipesOnRecipeStepsJoinClause).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):                                  recipeID,
				fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn):                          nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn):                              recipeStepID,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn):                      nil,
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.RecipeStepsTableBelongsToRecipeColumn): recipeID,
			}),
	)
}

// BuildGetAllRecipeStepsCountQuery returns a query that fetches the total number of recipe steps in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (b *Postgres) BuildGetAllRecipeStepsCountQuery(ctx context.Context) string {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQueryOnly(
		span,
		b.sqlBuilder.Select(fmt.Sprintf(columnCountQueryTemplate, querybuilding.RecipeStepsTableName)).
			From(querybuilding.RecipeStepsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn): nil,
			}),
	)
}

// BuildGetBatchOfRecipeStepsQuery returns a query that fetches every recipe step in the database within a bucketed range.
func (b *Postgres) BuildGetBatchOfRecipeStepsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipeStepsTableColumns...).
			From(querybuilding.RecipeStepsTableName).
			Where(squirrel.Gt{
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn): beginID,
			}).
			Where(squirrel.Lt{
				fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn): endID,
			}),
	)
}

// BuildGetRecipeStepsQuery builds a SQL query selecting recipe steps that adhere to a given QueryFilter and belong to a given account,
// and returns both the query and the relevant args to pass to the query executor.
func (b *Postgres) BuildGetRecipeStepsQuery(ctx context.Context, recipeID uint64, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	joins := []string{
		querybuilding.RecipesOnRecipeStepsJoinClause,
	}
	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.IDColumn):             recipeID,
		fmt.Sprintf("%s.%s", querybuilding.RecipesTableName, querybuilding.ArchivedOnColumn):     nil,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn): nil,
	}

	return b.buildListQuery(
		ctx,
		querybuilding.RecipeStepsTableName,
		joins,
		where,
		querybuilding.RecipeStepsTableBelongsToRecipeColumn,
		querybuilding.RecipeStepsTableColumns,
		recipeID,
		includeArchived,
		filter,
	)
}

// BuildGetRecipeStepsWithIDsQuery builds a SQL query selecting recipe steps that belong to a given account,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (b *Postgres) BuildGetRecipeStepsWithIDsQuery(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	where := squirrel.Eq{
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.IDColumn):                              ids,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.ArchivedOnColumn):                      nil,
		fmt.Sprintf("%s.%s", querybuilding.RecipeStepsTableName, querybuilding.RecipeStepsTableBelongsToRecipeColumn): recipeID,
	}

	subqueryBuilder := b.sqlBuilder.Select(querybuilding.RecipeStepsTableColumns...).
		From(querybuilding.RecipeStepsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.RecipeStepsTableColumns...).
			FromSelect(subqueryBuilder, querybuilding.RecipeStepsTableName).
			Where(where),
	)
}

// BuildCreateRecipeStepQuery takes a recipe step and returns a creation query for that recipe step and the relevant arguments.
func (b *Postgres) BuildCreateRecipeStepQuery(ctx context.Context, input *types.RecipeStepCreationInput) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	return b.buildQuery(
		span,
		b.sqlBuilder.Insert(querybuilding.RecipeStepsTableName).
			Columns(
				querybuilding.ExternalIDColumn,
				querybuilding.RecipeStepsTableIndexColumn,
				querybuilding.RecipeStepsTablePreparationIDColumn,
				querybuilding.RecipeStepsTablePrerequisiteStepColumn,
				querybuilding.RecipeStepsTableMinEstimatedTimeInSecondsColumn,
				querybuilding.RecipeStepsTableMaxEstimatedTimeInSecondsColumn,
				querybuilding.RecipeStepsTableTemperatureInCelsiusColumn,
				querybuilding.RecipeStepsTableNotesColumn,
				querybuilding.RecipeStepsTableWhyColumn,
				querybuilding.RecipeStepsTableBelongsToRecipeColumn,
			).
			Values(
				b.externalIDGenerator.NewExternalID(),
				input.Index,
				input.PreparationID,
				input.PrerequisiteStep,
				input.MinEstimatedTimeInSeconds,
				input.MaxEstimatedTimeInSeconds,
				input.TemperatureInCelsius,
				input.Notes,
				input.Why,
				input.BelongsToRecipe,
			).
			Suffix(fmt.Sprintf("RETURNING %s", querybuilding.IDColumn)),
	)
}

// BuildUpdateRecipeStepQuery takes a recipe step and returns an update SQL query, with the relevant query parameters.
func (b *Postgres) BuildUpdateRecipeStepQuery(ctx context.Context, input *types.RecipeStep) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, input.BelongsToRecipe)
	tracing.AttachRecipeStepIDToSpan(span, input.ID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.RecipeStepsTableName).
			Set(querybuilding.RecipeStepsTableIndexColumn, input.Index).
			Set(querybuilding.RecipeStepsTablePreparationIDColumn, input.PreparationID).
			Set(querybuilding.RecipeStepsTablePrerequisiteStepColumn, input.PrerequisiteStep).
			Set(querybuilding.RecipeStepsTableMinEstimatedTimeInSecondsColumn, input.MinEstimatedTimeInSeconds).
			Set(querybuilding.RecipeStepsTableMaxEstimatedTimeInSecondsColumn, input.MaxEstimatedTimeInSeconds).
			Set(querybuilding.RecipeStepsTableTemperatureInCelsiusColumn, input.TemperatureInCelsius).
			Set(querybuilding.RecipeStepsTableNotesColumn, input.Notes).
			Set(querybuilding.RecipeStepsTableWhyColumn, input.Why).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                              input.ID,
				querybuilding.ArchivedOnColumn:                      nil,
				querybuilding.RecipeStepsTableBelongsToRecipeColumn: input.BelongsToRecipe,
			}),
	)
}

// BuildArchiveRecipeStepQuery returns a SQL query which marks a given recipe step belonging to a given account as archived.
func (b *Postgres) BuildArchiveRecipeStepQuery(ctx context.Context, recipeID, recipeStepID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeIDToSpan(span, recipeID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	return b.buildQuery(
		span,
		b.sqlBuilder.Update(querybuilding.RecipeStepsTableName).
			Set(querybuilding.LastUpdatedOnColumn, currentUnixTimeQuery).
			Set(querybuilding.ArchivedOnColumn, currentUnixTimeQuery).
			Where(squirrel.Eq{
				querybuilding.IDColumn:                              recipeStepID,
				querybuilding.ArchivedOnColumn:                      nil,
				querybuilding.RecipeStepsTableBelongsToRecipeColumn: recipeID,
			}),
	)
}

// BuildGetAuditLogEntriesForRecipeStepQuery constructs a SQL query for fetching audit log entries relating to a recipe step with a given ID.
func (b *Postgres) BuildGetAuditLogEntriesForRecipeStepQuery(ctx context.Context, recipeStepID uint64) (query string, args []interface{}) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	recipeStepIDKey := fmt.Sprintf(
		jsonPluckQuery,
		querybuilding.AuditLogEntriesTableName,
		querybuilding.AuditLogEntriesTableContextColumn,
		audit.RecipeStepAssignmentKey,
	)

	return b.buildQuery(
		span,
		b.sqlBuilder.Select(querybuilding.AuditLogEntriesTableColumns...).
			From(querybuilding.AuditLogEntriesTableName).
			Where(squirrel.Eq{recipeStepIDKey: recipeStepID}).
			OrderBy(fmt.Sprintf("%s.%s", querybuilding.AuditLogEntriesTableName, querybuilding.CreatedOnColumn)),
	)
}
