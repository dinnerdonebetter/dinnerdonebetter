package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	mealPlanTasksTableName = "meal_plan_tasks"

	mealPlanTaskIDColumn                = "meal_plan_task_id"
	mealPlanTaskCompletedAtColumn       = "completed_at"
	mealPlanTaskStatusExplanationColumn = "status_explanation"
	mealPlanTaskStatusColumn            = "status"
)

var mealPlanTasksColumns = []string{
	idColumn,
	mealPlanTaskStatusColumn,
	mealPlanTaskStatusExplanationColumn,
	"creation_explanation",
	"belongs_to_meal_plan_option",
	"belongs_to_recipe_prep_task",
	mealPlanTaskCompletedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	"assigned_to_user",
}

func buildMealPlanTasksQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(mealPlanTasksColumns, mealPlanTaskCompletedAtColumn)

		fullSelectColumns := mergeColumns(
			applyToEach(mealPlanTasksColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s", mealPlanTasksTableName, s)
			}),
			append(
				applyToEach(mealPlanOptionsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s as meal_plan_option_%s", mealPlanOptionsTableName, s, s)
				}),
				append(
					applyToEach(recipePrepTasksColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s as prep_task_%s", recipePrepTasksTableName, s, s)
					}),
					applyToEach(recipePrepTaskStepsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s as prep_task_step_%s", recipePrepTaskStepsTableName, s, s)
					})...,
				)...,
			),
			1,
		)

		applyToEach(recipeStepsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s", s, s)
		})

		completeSelectColumns := mergeColumns(
			applyToEach(mealPlanTasksColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s", mealPlanTasksTableName, s)
			}),
			append(
				applyToEach(mealPlanOptionsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s as meal_plan_option_%s", mealPlanOptionsTableName, s, s)
				}),

				mergeColumns(
					applyToEach(recipeStepsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s as recipe_step_%s", recipeStepsTableName, s, s)
					}),
					applyToEach(validPreparationsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s as valid_preparation_%s", validPreparationsTableName, s, s)
					}),
					2,
				)...,
			),
			1,
		)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ChangeMealPlanTaskStatus",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = sqlc.arg(%s),
	%s = sqlc.arg(%s),
	%s = sqlc.arg(%s)
WHERE %s = sqlc.arg(%s);`,
					mealPlanTasksTableName,
					mealPlanTaskCompletedAtColumn, mealPlanTaskCompletedAtColumn,
					mealPlanTaskStatusExplanationColumn, mealPlanTaskStatusExplanationColumn,
					mealPlanTaskStatusColumn, mealPlanTaskStatusColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateMealPlanTask",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					mealPlanTasksTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckMealPlanTaskExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
		FULL OUTER JOIN %s ON %s.%s=%s.%s
		FULL OUTER JOIN %s ON %s.%s=%s.%s
		FULL OUTER JOIN %s ON %s.%s=%s.%s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
					mealPlanTasksTableName, idColumn,
					mealPlanTasksTableName,
					mealPlanOptionsTableName, mealPlanTasksTableName, belongsToMealPlanOptionColumn, mealPlanOptionsTableName, idColumn,
					mealPlanEventsTableName, mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventsTableName, idColumn,
					mealPlansTableName, mealPlanEventsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
					mealPlanTasksTableName, mealPlanTaskCompletedAtColumn,
					mealPlansTableName, idColumn, mealPlanIDColumn,
					mealPlansTableName, archivedAtColumn,
					mealPlanTasksTableName, idColumn, mealPlanTaskIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetMealPlanTask",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					mealPlanTasksTableName,
					mealPlanOptionsTableName, mealPlanTasksTableName, belongsToMealPlanOptionColumn, mealPlanOptionsTableName, idColumn,
					mealPlanEventsTableName, mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventsTableName, idColumn,
					mealPlansTableName, mealPlanEventsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
					mealsTableName, mealPlanOptionsTableName, mealIDColumn, mealsTableName, idColumn,
					recipePrepTasksTableName, mealPlanTasksTableName, belongsToRecipePrepTaskColumn, recipePrepTasksTableName, idColumn,
					recipePrepTaskStepsTableName, recipePrepTaskStepsTableName, belongsToRecipePrepTaskColumn, recipePrepTasksTableName, idColumn,
					recipeStepsTableName, recipePrepTaskStepsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					mealPlanOptionsTableName, archivedAtColumn,
					mealPlanEventsTableName, archivedAtColumn,
					mealPlansTableName, archivedAtColumn,
					mealsTableName, archivedAtColumn,
					recipeStepsTableName, archivedAtColumn,
					mealPlanTasksTableName, idColumn, mealPlanTaskIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ListAllMealPlanTasksByMealPlan",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					mealPlanTasksTableName,
					mealPlanOptionsTableName, mealPlanTasksTableName, belongsToMealPlanOptionColumn, mealPlanOptionsTableName, idColumn,
					mealPlanEventsTableName, mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventsTableName, idColumn,
					mealPlansTableName, mealPlanEventsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
					mealsTableName, mealPlanOptionsTableName, mealIDColumn, mealsTableName, idColumn,
					recipePrepTasksTableName, mealPlanTasksTableName, belongsToRecipePrepTaskColumn, recipePrepTasksTableName, idColumn,
					recipePrepTaskStepsTableName, recipePrepTaskStepsTableName, belongsToRecipePrepTaskColumn, recipePrepTasksTableName, idColumn,
					recipeStepsTableName, recipePrepTaskStepsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					mealPlanOptionsTableName, archivedAtColumn,
					mealPlanEventsTableName, archivedAtColumn,
					mealPlansTableName, archivedAtColumn,
					mealsTableName, archivedAtColumn,
					recipeStepsTableName, archivedAtColumn,
					mealPlansTableName, idColumn, mealPlanIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ListIncompleteMealPlanTasksByMealPlanOption",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	 FULL OUTER JOIN %s ON %s.%s=%s.%s
	 FULL OUTER JOIN %s ON %s.%s=%s.%s
	 FULL OUTER JOIN %s ON %s.%s=%s.%s
	 JOIN %s ON %s.%s=%s.%s
	 JOIN %s ON %s.%s=%s.%s
WHERE %s.%s = sqlc.arg(%s)
AND %s.%s IS NULL;`,
					strings.Join(completeSelectColumns, ",\n\t"),
					mealPlanTasksTableName,
					mealPlanOptionsTableName, mealPlanTasksTableName, belongsToMealPlanOptionColumn, mealPlanOptionsTableName, idColumn,
					mealPlansTableName, mealPlanOptionsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
					mealsTableName, mealPlanOptionsTableName, mealIDColumn, mealsTableName, idColumn,
					recipeStepsTableName, mealPlanTasksTableName, satisfiesRecipeStepColumn, recipeStepsTableName, idColumn,
					validPreparationsTableName, recipeStepsTableName, preparationIDColumn, validPreparationsTableName, idColumn,
					mealPlanTasksTableName, belongsToMealPlanOptionColumn, belongsToMealPlanOptionColumn,
					mealPlanTasksTableName, mealPlanTaskCompletedAtColumn,
				)),
			},
		}
	default:
		return nil
	}
}
