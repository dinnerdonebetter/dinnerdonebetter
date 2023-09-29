package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeStepProductsTableName = "recipe_step_products"

	recipeStepProductIDColumn = "recipe_step_product_id"
)

var recipeStepProductsColumns = []string{
	idColumn,
	nameColumn,
	"type",
	"measurement_unit",
	"minimum_quantity_value",
	"maximum_quantity_value",
	"quantity_notes",
	"compostable",
	"maximum_storage_duration_in_seconds",
	"minimum_storage_temperature_in_celsius",
	"maximum_storage_temperature_in_celsius",
	"storage_instructions",
	"is_liquid",
	"is_waste",
	"index",
	"contained_in_vessel_index",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	belongsToRecipeStepColumn,
}

func buildRecipeStepProductsQueries() []*Query {
	insertColumns := filterForInsert(recipeStepProductsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveRecipeStepProduct",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				recipeStepProductsTableName,
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
				Name: "CreateRecipeStepProduct",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipeStepProductsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckRecipeStepProductExistence",
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
				recipeStepProductsTableName, idColumn,
				recipeStepProductsTableName,
				recipeStepsTableName, recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepProductsTableName, archivedAtColumn,
				recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeStepProductsTableName, idColumn, recipeStepProductIDColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepProductsForRecipe",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	recipe_step_products.id,
	recipe_step_products.name,
	recipe_step_products.type,
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
	recipe_step_products.minimum_quantity_value,
	recipe_step_products.maximum_quantity_value,
	recipe_step_products.quantity_notes,
	recipe_step_products.compostable,
	recipe_step_products.maximum_storage_duration_in_seconds,
	recipe_step_products.minimum_storage_temperature_in_celsius,
	recipe_step_products.maximum_storage_temperature_in_celsius,
	recipe_step_products.storage_instructions,
	recipe_step_products.is_liquid,
	recipe_step_products.is_waste,
	recipe_step_products.index,
	recipe_step_products.contained_in_vessel_index,
	recipe_step_products.created_at,
	recipe_step_products.last_updated_at,
	recipe_step_products.archived_at,
	recipe_step_products.belongs_to_recipe_step
FROM recipe_step_products
	JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	LEFT JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepProducts",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	recipe_step_products.id,
	recipe_step_products.name,
	recipe_step_products.type,
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
	recipe_step_products.minimum_quantity_value,
	recipe_step_products.maximum_quantity_value,
	recipe_step_products.quantity_notes,
	recipe_step_products.compostable,
	recipe_step_products.maximum_storage_duration_in_seconds,
	recipe_step_products.minimum_storage_temperature_in_celsius,
	recipe_step_products.maximum_storage_temperature_in_celsius,
	recipe_step_products.storage_instructions,
	recipe_step_products.is_liquid,
	recipe_step_products.is_waste,
	recipe_step_products.index,
	recipe_step_products.contained_in_vessel_index,
	recipe_step_products.created_at,
	recipe_step_products.last_updated_at,
	recipe_step_products.archived_at,
	recipe_step_products.belongs_to_recipe_step,
	(
		SELECT COUNT(recipe_step_products.id)
		FROM recipe_step_products
		WHERE
			recipe_step_products.archived_at IS NULL
			AND recipe_step_products.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
			AND recipe_step_products.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
			AND recipe_step_products.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
			AND (
				recipe_step_products.last_updated_at IS NULL
				OR recipe_step_products.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
			)
			AND (
				recipe_step_products.last_updated_at IS NULL
				OR recipe_step_products.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
			)
	) AS filtered_count,
	(
		SELECT COUNT(recipe_step_products.id)
		FROM recipe_step_products
		WHERE recipe_step_products.archived_at IS NULL
			AND recipe_step_products.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
	) AS total_count
FROM recipe_step_products
	JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	LEFT JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_at IS NULL
	AND recipe_step_products.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.id = sqlc.arg(recipe_step_id)
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id)
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepProduct",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	recipe_step_products.id,
	recipe_step_products.name,
	recipe_step_products.type,
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
	recipe_step_products.minimum_quantity_value,
	recipe_step_products.maximum_quantity_value,
	recipe_step_products.quantity_notes,
	recipe_step_products.compostable,
	recipe_step_products.maximum_storage_duration_in_seconds,
	recipe_step_products.minimum_storage_temperature_in_celsius,
	recipe_step_products.maximum_storage_temperature_in_celsius,
	recipe_step_products.storage_instructions,
	recipe_step_products.is_liquid,
	recipe_step_products.is_waste,
	recipe_step_products.index,
	recipe_step_products.contained_in_vessel_index,
	recipe_step_products.created_at,
	recipe_step_products.last_updated_at,
	recipe_step_products.archived_at,
	recipe_step_products.belongs_to_recipe_step
FROM recipe_step_products
	JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	LEFT JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_at IS NULL
	AND recipe_step_products.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
	AND recipe_step_products.id = sqlc.arg(recipe_step_product_id)
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipe_steps.id = sqlc.arg(recipe_step_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipeStepProduct",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				recipeStepProductsTableName,
				strings.Join(applyToEach(filterForUpdate(recipeStepProductsColumns, belongsToRecipeStepColumn), func(i int, s string) string {
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
