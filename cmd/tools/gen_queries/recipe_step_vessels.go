package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeStepVesselsTableName = "recipe_step_vessels"

	recipeStepVesselIDColumn = "recipe_step_vessel_id"
)

var recipeStepVesselsColumns = []string{
	idColumn,
	nameColumn,
	notesColumn,
	belongsToRecipeStepColumn,
	"recipe_step_product_id",
	"valid_vessel_id",
	"vessel_predicate",
	"minimum_quantity",
	"maximum_quantity",
	"unavailable_after_step",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepVesselsQueries() []*Query {
	insertColumns := filterForInsert(recipeStepVesselsColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(filterFromSlice(recipeStepVesselsColumns, validVesselIDColumn), func(i int, s string) string {
			return fmt.Sprintf("%s.%s", recipeStepVesselsTableName, s)
		}),
		mergeColumns(
			applyToEach(filterFromSlice(validVesselsColumns, capacityUnitColumn), func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_vessel_%s", validVesselsTableName, s, s)
			}),
			applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_measurement_unit_%s", validMeasurementUnitsTableName, s, s)
			}),
			10,
		),
		1,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveRecipeStepVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				recipeStepVesselsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToRecipeStepColumn,
				belongsToRecipeStepColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateRecipeStepVessel",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipeStepVesselsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckRecipeStepVesselExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
		JOIN %s ON %s.%s=%s.%s
		JOIN %s ON %s.%s=%s.%s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				recipeStepVesselsTableName, idColumn,
				recipeStepVesselsTableName,
				recipeStepsTableName, recipeStepVesselsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepVesselsTableName, archivedAtColumn,
				recipeStepVesselsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeStepVesselsTableName, idColumn, recipeStepVesselIDColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepVesselsForRecipe",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM recipe_step_vessels
	LEFT JOIN valid_vessels ON recipe_step_vessels.valid_vessel_id=valid_vessels.id
	LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
	JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_vessels.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id);`,
				strings.Join(fullSelectColumns, ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepVessel",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM recipe_step_vessels
	LEFT JOIN valid_vessels ON recipe_step_vessels.valid_vessel_id=valid_vessels.id
	LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
	JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_vessels.archived_at IS NULL
	AND recipe_step_vessels.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
	AND recipe_step_vessels.id = sqlc.arg(recipe_step_vessel_id)
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipe_steps.id = sqlc.arg(recipe_step_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id);`,
				strings.Join(fullSelectColumns, ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepVessels",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	(
		SELECT COUNT(recipe_step_vessels.id)
		FROM recipe_step_vessels
		WHERE
			recipe_step_vessels.archived_at IS NULL
			AND recipe_step_vessels.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
			AND recipe_step_vessels.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
			AND (recipe_step_vessels.last_updated_at IS NULL OR recipe_step_vessels.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
			AND (recipe_step_vessels.last_updated_at IS NULL OR recipe_step_vessels.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
	) as filtered_count,
	(
		SELECT COUNT(recipe_step_vessels.id)
		FROM recipe_step_vessels
		WHERE recipe_step_vessels.archived_at IS NULL
	) as total_count
FROM recipe_step_vessels
	 LEFT JOIN valid_vessels ON recipe_step_vessels.valid_vessel_id=valid_vessels.id
	 LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
	 JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
	 JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_vessels.archived_at IS NULL
	AND recipe_step_vessels.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.id = sqlc.arg(recipe_step_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id)
	AND recipe_step_vessels.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
	AND recipe_step_vessels.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
	AND (recipe_step_vessels.last_updated_at IS NULL OR recipe_step_vessels.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
	AND (recipe_step_vessels.last_updated_at IS NULL OR recipe_step_vessels.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);`,
				strings.Join(fullSelectColumns, ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipeStepVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				recipeStepVesselsTableName,
				strings.Join(applyToEach(filterForUpdate(recipeStepVesselsColumns), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				belongsToRecipeStepColumn, belongsToRecipeStepColumn,
				idColumn, idColumn,
			)),
		},
	}
}
