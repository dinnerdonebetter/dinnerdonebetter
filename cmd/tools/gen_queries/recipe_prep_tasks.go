package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipePrepTasksTableName      = "recipe_prep_tasks"
	belongsToRecipePrepTaskColumn = "belongs_to_recipe_prep_task"
)

var recipePrepTasksColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	notesColumn,
	"optional",
	"explicit_storage_instructions",
	"minimum_time_buffer_before_recipe_in_seconds",
	"maximum_time_buffer_before_recipe_in_seconds",
	"storage_type",
	"minimum_storage_temperature_in_celsius",
	"maximum_storage_temperature_in_celsius",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	belongsToRecipeColumn,
}

func buildRecipePrepTasksQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(recipePrepTasksColumns)

		fullSelectColumns := append(
			applyToEach(recipePrepTasksColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s", recipePrepTasksTableName, s)
			}),
			applyToEach(recipePrepTaskStepsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as task_step_%s", recipePrepTaskStepsTableName, s, s)
			})...,
		)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveRecipePrepTask",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					recipePrepTasksTableName,
					archivedAtColumn,
					currentTimeExpression,
					archivedAtColumn,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateRecipePrepTask",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					recipePrepTasksTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckRecipePrepTaskExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
		JOIN %s ON %s.%s=%s.%s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(recipe_id)
		AND %s.%s = sqlc.arg(recipe_prep_task_id)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(recipe_id)
);`,
					recipePrepTasksTableName, idColumn,
					recipePrepTasksTableName,
					recipesTableName, recipePrepTasksTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					recipePrepTasksTableName, archivedAtColumn,
					recipePrepTasksTableName, belongsToRecipeColumn,
					recipePrepTasksTableName, idColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipePrepTask",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(recipe_prep_task_id);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					recipePrepTasksTableName,
					recipePrepTaskStepsTableName, recipePrepTasksTableName, idColumn, recipePrepTaskStepsTableName, belongsToRecipePrepTaskColumn,
					recipePrepTasksTableName, archivedAtColumn,
					recipePrepTasksTableName, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ListAllRecipePrepTasksByRecipe",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(recipe_id);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					recipePrepTasksTableName,
					recipePrepTaskStepsTableName, recipePrepTaskStepsTableName, belongsToRecipePrepTaskColumn, recipePrepTasksTableName, idColumn,
					recipeStepsTableName, recipePrepTaskStepsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					recipesTableName, recipePrepTasksTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					recipePrepTasksTableName, archivedAtColumn,
					recipeStepsTableName, archivedAtColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateRecipePrepTask",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					recipePrepTasksTableName,
					strings.Join(applyToEach(filterForUpdate(recipePrepTasksColumns), func(i int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
		}
	default:
		return nil
	}
}
