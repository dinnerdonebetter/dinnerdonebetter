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
	recipe_step_vessels.id,
	valid_vessels.id as valid_vessel_id,
	valid_vessels.name as valid_vessel_name,
	valid_vessels.plural_name as valid_vessel_plural_name,
	valid_vessels.description as valid_vessel_description,
	valid_vessels.icon_path as valid_vessel_icon_path,
	valid_vessels.usable_for_storage as valid_vessel_usable_for_storage,
	valid_vessels.slug as valid_vessel_slug,
	valid_vessels.display_in_summary_lists as valid_vessel_display_in_summary_lists,
	valid_vessels.include_in_generated_instructions as valid_vessel_include_in_generated_instructions,
	valid_vessels.capacity as valid_vessel_capacity,
	valid_measurement_units.id as valid_measurement_unit_id,
	valid_measurement_units.name as valid_measurement_unit_name,
	valid_measurement_units.description as valid_measurement_unit_description,
	valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
	valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
	valid_measurement_units.universal as valid_measurement_unit_universal,
	valid_measurement_units.metric as valid_measurement_unit_metric,
	valid_measurement_units.imperial as valid_measurement_unit_imperial,
	valid_measurement_units.slug as valid_measurement_unit_slug,
	valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
	valid_measurement_units.created_at as valid_measurement_unit_created_at,
	valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
	valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
	valid_vessels.width_in_millimeters as valid_vessel_width_in_millimeters,
	valid_vessels.length_in_millimeters as valid_vessel_length_in_millimeters,
	valid_vessels.height_in_millimeters as valid_vessel_height_in_millimeters,
	valid_vessels.shape as valid_vessel_shape,
	valid_vessels.created_at as valid_vessel_created_at,
	valid_vessels.last_updated_at as valid_vessel_last_updated_at,
	valid_vessels.archived_at as valid_vessel_archived_at,
	recipe_step_vessels.name,
	recipe_step_vessels.notes,
	recipe_step_vessels.belongs_to_recipe_step,
	recipe_step_vessels.recipe_step_product_id,
	recipe_step_vessels.vessel_predicate,
	recipe_step_vessels.minimum_quantity,
	recipe_step_vessels.maximum_quantity,
	recipe_step_vessels.unavailable_after_step,
	recipe_step_vessels.created_at,
	recipe_step_vessels.last_updated_at,
	recipe_step_vessels.archived_at
FROM recipe_step_vessels
	LEFT JOIN valid_vessels ON recipe_step_vessels.valid_vessel_id=valid_vessels.id
	LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
	JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_vessels.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepVessel",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	recipe_step_vessels.id,
	valid_vessels.id as valid_vessel_id,
	valid_vessels.name as valid_vessel_name,
	valid_vessels.plural_name as valid_vessel_plural_name,
	valid_vessels.description as valid_vessel_description,
	valid_vessels.icon_path as valid_vessel_icon_path,
	valid_vessels.usable_for_storage as valid_vessel_usable_for_storage,
	valid_vessels.slug as valid_vessel_slug,
	valid_vessels.display_in_summary_lists as valid_vessel_display_in_summary_lists,
	valid_vessels.include_in_generated_instructions as valid_vessel_include_in_generated_instructions,
	valid_vessels.capacity as valid_vessel_capacity,
	valid_measurement_units.id as valid_measurement_unit_id,
	valid_measurement_units.name as valid_measurement_unit_name,
	valid_measurement_units.description as valid_measurement_unit_description,
	valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
	valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
	valid_measurement_units.universal as valid_measurement_unit_universal,
	valid_measurement_units.metric as valid_measurement_unit_metric,
	valid_measurement_units.imperial as valid_measurement_unit_imperial,
	valid_measurement_units.slug as valid_measurement_unit_slug,
	valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
	valid_measurement_units.created_at as valid_measurement_unit_created_at,
	valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
	valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
	valid_vessels.width_in_millimeters as valid_vessel_width_in_millimeters,
	valid_vessels.length_in_millimeters as valid_vessel_length_in_millimeters,
	valid_vessels.height_in_millimeters as valid_vessel_height_in_millimeters,
	valid_vessels.shape as valid_vessel_shape,
	valid_vessels.created_at as valid_vessel_created_at,
	valid_vessels.last_updated_at as valid_vessel_last_updated_at,
	valid_vessels.archived_at as valid_vessel_archived_at,
	recipe_step_vessels.name,
	recipe_step_vessels.notes,
	recipe_step_vessels.belongs_to_recipe_step,
	recipe_step_vessels.recipe_step_product_id,
	recipe_step_vessels.vessel_predicate,
	recipe_step_vessels.minimum_quantity,
	recipe_step_vessels.maximum_quantity,
	recipe_step_vessels.unavailable_after_step,
	recipe_step_vessels.created_at,
	recipe_step_vessels.last_updated_at,
	recipe_step_vessels.archived_at
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
	AND recipes.id = sqlc.arg(recipe_id);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepVessels",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	recipe_step_vessels.id,
	valid_vessels.id as valid_vessel_id,
	valid_vessels.name as valid_vessel_name,
	valid_vessels.plural_name as valid_vessel_plural_name,
	valid_vessels.description as valid_vessel_description,
	valid_vessels.icon_path as valid_vessel_icon_path,
	valid_vessels.usable_for_storage as valid_vessel_usable_for_storage,
	valid_vessels.slug as valid_vessel_slug,
	valid_vessels.display_in_summary_lists as valid_vessel_display_in_summary_lists,
	valid_vessels.include_in_generated_instructions as valid_vessel_include_in_generated_instructions,
	valid_vessels.capacity as valid_vessel_capacity,
	valid_measurement_units.id as valid_measurement_unit_id,
	valid_measurement_units.name as valid_measurement_unit_name,
	valid_measurement_units.description as valid_measurement_unit_description,
	valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
	valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
	valid_measurement_units.universal as valid_measurement_unit_universal,
	valid_measurement_units.metric as valid_measurement_unit_metric,
	valid_measurement_units.imperial as valid_measurement_unit_imperial,
	valid_measurement_units.slug as valid_measurement_unit_slug,
	valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
	valid_measurement_units.created_at as valid_measurement_unit_created_at,
	valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
	valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
	valid_vessels.width_in_millimeters as valid_vessel_width_in_millimeters,
	valid_vessels.length_in_millimeters as valid_vessel_length_in_millimeters,
	valid_vessels.height_in_millimeters as valid_vessel_height_in_millimeters,
	valid_vessels.shape as valid_vessel_shape,
	valid_vessels.created_at as valid_vessel_created_at,
	valid_vessels.last_updated_at as valid_vessel_last_updated_at,
	valid_vessels.archived_at as valid_vessel_archived_at,
	recipe_step_vessels.name,
	recipe_step_vessels.notes,
	recipe_step_vessels.belongs_to_recipe_step,
	recipe_step_vessels.recipe_step_product_id,
	recipe_step_vessels.vessel_predicate,
	recipe_step_vessels.minimum_quantity,
	recipe_step_vessels.maximum_quantity,
	recipe_step_vessels.unavailable_after_step,
	recipe_step_vessels.created_at,
	recipe_step_vessels.last_updated_at,
	recipe_step_vessels.archived_at,
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
	AND recipe_step_vessels.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
	AND recipe_step_vessels.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
	AND (recipe_step_vessels.last_updated_at IS NULL OR recipe_step_vessels.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
	AND (recipe_step_vessels.last_updated_at IS NULL OR recipe_step_vessels.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.id = sqlc.arg(recipe_step_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id)
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);`)),
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
