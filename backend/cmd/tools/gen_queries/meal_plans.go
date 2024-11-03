package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	mealPlansTableName = "meal_plans"

	mealPlanIDColumn                     = "meal_plan_id"
	mealPlanStatusColumn                 = "status"
	mealPlanVotingDeadlineColumn         = "voting_deadline"
	mealPlanGroceryListInitializedColumn = "grocery_list_initialized"
	mealPlanTasksCreatedColumn           = "tasks_created"
	electionMethodColumn                 = "election_method"
)

var mealPlansColumns = []string{
	idColumn,
	notesColumn,
	mealPlanStatusColumn,
	mealPlanVotingDeadlineColumn,
	mealPlanGroceryListInitializedColumn,
	mealPlanTasksCreatedColumn,
	electionMethodColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	belongsToHouseholdColumn,
	createdByUserColumn,
}

func buildMealPlansQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(mealPlansColumns, mealPlanGroceryListInitializedColumn, mealPlanTasksCreatedColumn, electionMethodColumn)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveMealPlan",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
					mealPlansTableName,
					archivedAtColumn,
					currentTimeExpression,
					archivedAtColumn,
					belongsToHouseholdColumn,
					belongsToHouseholdColumn,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateMealPlan",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					mealPlansTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckMealPlanExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
);`,
					mealPlansTableName, idColumn,
					mealPlansTableName,
					mealPlansTableName, archivedAtColumn,
					mealPlansTableName, idColumn, mealPlanIDColumn,
					mealPlansTableName, belongsToHouseholdColumn, belongsToHouseholdColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "FinalizeMealPlan",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(
					`UPDATE %s SET %s = sqlc.arg(%s) WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					mealPlansTableName,
					mealPlanStatusColumn,
					mealPlanStatusColumn,
					archivedAtColumn,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetExpiredAndUnresolvedMealPlans",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = 'awaiting_votes'
	AND %s < %s
GROUP BY %s.%s
ORDER BY %s.%s;`,
					strings.Join(applyToEach(mealPlansColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", mealPlansTableName, s)
					}), ",\n\t"),
					mealPlansTableName,
					mealPlansTableName, archivedAtColumn,
					mealPlansTableName, mealPlanStatusColumn,
					mealPlanVotingDeadlineColumn, currentTimeExpression,
					mealPlansTableName, idColumn,
					mealPlansTableName, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetFinalizedMealPlansForPlanning",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s.%s as meal_plan_id,
	%s.%s as meal_plan_option_id,
	%s.%s as meal_id,
	%s.%s as meal_plan_event_id,
	%s.%s as %s
FROM
	%s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s = 'finalized'
	AND %s.%s IS TRUE
	AND %s.%s IS FALSE
GROUP BY
	%s.%s,
	%s.%s,
	%s.%s,
	%s.%s,
	%s.%s
ORDER BY
	%s.%s;`,
					mealPlansTableName, idColumn,
					mealPlanOptionsTableName, idColumn,
					mealsTableName, idColumn,
					mealPlanEventsTableName, idColumn,
					mealComponentsTableName, recipeIDColumn, recipeIDColumn,
					mealPlanOptionsTableName,
					mealPlanEventsTableName, mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventsTableName, idColumn,
					mealPlansTableName, mealPlanEventsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
					mealComponentsTableName, mealPlanOptionsTableName, mealIDColumn, mealComponentsTableName, mealIDColumn,
					mealsTableName, mealPlanOptionsTableName, mealIDColumn, mealsTableName, idColumn,
					mealPlansTableName, archivedAtColumn,
					mealPlansTableName, mealPlanStatusColumn,
					mealPlanOptionsTableName, mealPlanOptionsChosenColumn,
					mealPlansTableName, mealPlanTasksCreatedColumn,
					mealPlansTableName, idColumn,
					mealPlanOptionsTableName, idColumn,
					mealsTableName, idColumn,
					mealPlanEventsTableName, idColumn,
					mealComponentsTableName, recipeIDColumn,
					mealPlansTableName, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetFinalizedMealPlansWithoutGroceryListInit",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s.%s,
	%s.%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = 'finalized'
	AND %s.%s IS FALSE;`,
					mealPlansTableName, idColumn,
					mealPlansTableName, belongsToHouseholdColumn,
					mealPlansTableName,
					mealPlansTableName, archivedAtColumn,
					mealPlansTableName, mealPlanStatusColumn,
					mealPlansTableName, mealPlanGroceryListInitializedColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetMealPlan",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
  AND %s.%s = sqlc.arg(%s)
  AND %s.%s = sqlc.arg(%s);`,
					strings.Join(applyToEach(mealPlansColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", mealPlansTableName, s)
					}), ",\n\t"),
					mealPlansTableName,
					mealPlansTableName, archivedAtColumn,
					mealPlansTableName, idColumn, idColumn,
					mealPlansTableName, belongsToHouseholdColumn, belongsToHouseholdColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetMealPlansForHousehold",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s.%s IS NULL
	%s
%s;`,
					strings.Join(applyToEach(mealPlansColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", mealPlansTableName, s)
					}), ",\n\t"),
					buildFilterCountSelect(mealPlansTableName, true, true, []string{}, fmt.Sprintf("%s.%s = sqlc.arg(%s)", mealPlansTableName, belongsToHouseholdColumn, belongsToHouseholdColumn)),
					buildTotalCountSelect(mealPlansTableName, true, []string{}, fmt.Sprintf("%s.%s = sqlc.arg(%s)", mealPlansTableName, belongsToHouseholdColumn, belongsToHouseholdColumn)),
					mealPlansTableName,
					mealPlansTableName, archivedAtColumn,
					buildFilterConditions(
						mealPlansTableName,
						true,
						true,
						nil,
						fmt.Sprintf("%s.%s = sqlc.arg(%s)", mealPlansTableName, belongsToHouseholdColumn, belongsToHouseholdColumn),
					),
					offsetLimitAddendum,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetMealPlanPastVotingDeadline",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(meal_plan_id)
	AND %s.%s = sqlc.arg(household_id)
	AND %s.%s = 'awaiting_votes'
	AND %s > %s.%s;`,
					strings.Join(applyToEach(mealPlansColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", mealPlansTableName, s)
					}), ",\n\t"),
					mealPlansTableName,
					mealPlansTableName, archivedAtColumn,
					mealPlansTableName, idColumn,
					mealPlansTableName, belongsToHouseholdColumn,
					mealPlansTableName, mealPlanStatusColumn,
					currentTimeExpression, mealPlansTableName, mealPlanVotingDeadlineColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "MarkMealPlanAsGroceryListInitialized",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = TRUE,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					mealPlansTableName,
					mealPlanGroceryListInitializedColumn,
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "MarkMealPlanAsPrepTasksCreated",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = TRUE,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					mealPlansTableName,
					mealPlanTasksCreatedColumn,
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateMealPlan",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
					mealPlansTableName,
					strings.Join(applyToEach(filterForUpdate(mealPlansColumns, mealPlanGroceryListInitializedColumn, mealPlanTasksCreatedColumn, electionMethodColumn, belongsToHouseholdColumn, createdByUserColumn), func(i int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					belongsToHouseholdColumn, belongsToHouseholdColumn,
					idColumn, idColumn,
				)),
			},
		}
	default:
		return nil
	}
}
