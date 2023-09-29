package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	mealPlanGroceryListItemsTableName = "meal_plan_grocery_list_items"

	mealPlanGroceryListItemIDColumn = "meal_plan_grocery_list_item_id"
)

var mealPlanGroceryListItemsColumns = []string{
	idColumn,
	"belongs_to_meal_plan",
	"valid_ingredient",
	"valid_measurement_unit",
	"minimum_quantity_needed",
	"maximum_quantity_needed",
	"quantity_purchased",
	"purchased_measurement_unit",
	"purchased_upc",
	"purchase_price",
	"status_explanation",
	"status",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealPlanGroceryListItemsQueries() []*Query {
	insertColumns := filterForInsert(mealPlanGroceryListItemsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveMealPlanGroceryListItem",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				mealPlanGroceryListItemsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateMealPlanGroceryListItem",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				mealPlanGroceryListItemsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckMealPlanGroceryListItemExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
);`,
				mealPlanGroceryListItemsTableName, idColumn,
				mealPlanGroceryListItemsTableName,
				mealPlanGroceryListItemsTableName, archivedAtColumn,
				mealPlanGroceryListItemsTableName, idColumn, mealPlanGroceryListItemIDColumn,
				mealPlanGroceryListItemsTableName, belongsToMealPlanColumn, mealPlanIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetMealPlanGroceryListItemsForMealPlan",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	meal_plan_grocery_list_items.id,
	meal_plan_grocery_list_items.belongs_to_meal_plan,
	valid_ingredients.id as valid_ingredient_id,
	valid_ingredients.name as valid_ingredient_name,
	valid_ingredients.description as valid_ingredient_description,
	valid_ingredients.warning as valid_ingredient_warning,
	valid_ingredients.contains_egg as valid_ingredient_contains_egg,
	valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
	valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
	valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
	valid_ingredients.contains_soy as valid_ingredient_contains_soy,
	valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
	valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
	valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
	valid_ingredients.contains_fish as valid_ingredient_contains_fish,
	valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
	valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
	valid_ingredients.volumetric as valid_ingredient_volumetric,
	valid_ingredients.is_liquid as valid_ingredient_is_liquid,
	valid_ingredients.icon_path as valid_ingredient_icon_path,
	valid_ingredients.animal_derived as valid_ingredient_animal_derived,
	valid_ingredients.plural_name as valid_ingredient_plural_name,
	valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
	valid_ingredients.slug as valid_ingredient_slug,
	valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
	valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
	valid_ingredients.is_starch as valid_ingredient_is_starch,
	valid_ingredients.is_protein as valid_ingredient_is_protein,
	valid_ingredients.is_grain as valid_ingredient_is_grain,
	valid_ingredients.is_fruit as valid_ingredient_is_fruit,
	valid_ingredients.is_salt as valid_ingredient_is_salt,
	valid_ingredients.is_fat as valid_ingredient_is_fat,
	valid_ingredients.is_acid as valid_ingredient_is_acid,
	valid_ingredients.is_heat as valid_ingredient_is_heat,
	valid_ingredients.created_at as valid_ingredient_created_at,
	valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
	valid_ingredients.archived_at as valid_ingredient_archived_at,
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
	meal_plan_grocery_list_items.minimum_quantity_needed,
	meal_plan_grocery_list_items.maximum_quantity_needed,
	meal_plan_grocery_list_items.quantity_purchased,
	meal_plan_grocery_list_items.purchased_measurement_unit,
	meal_plan_grocery_list_items.purchased_upc,
	meal_plan_grocery_list_items.purchase_price,
	meal_plan_grocery_list_items.status_explanation,
	meal_plan_grocery_list_items.status,
	meal_plan_grocery_list_items.created_at,
	meal_plan_grocery_list_items.last_updated_at,
	meal_plan_grocery_list_items.archived_at
FROM meal_plan_grocery_list_items
	JOIN meal_plans ON meal_plan_grocery_list_items.belongs_to_meal_plan=meal_plans.id
	JOIN valid_ingredients ON meal_plan_grocery_list_items.valid_ingredient=valid_ingredients.id
	JOIN valid_measurement_units ON meal_plan_grocery_list_items.valid_measurement_unit=valid_measurement_units.id
WHERE meal_plan_grocery_list_items.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND meal_plan_grocery_list_items.belongs_to_meal_plan = sqlc.arg(meal_plan_id)
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = sqlc.arg(meal_plan_id)
GROUP BY meal_plan_grocery_list_items.id,
	valid_ingredients.id,
	valid_measurement_units.id,
	meal_plans.id
ORDER BY meal_plan_grocery_list_items.id;`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetMealPlanGroceryListItem",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	meal_plan_grocery_list_items.id,
	meal_plan_grocery_list_items.belongs_to_meal_plan,
	valid_ingredients.id as valid_ingredient_id,
	valid_ingredients.name as valid_ingredient_name,
	valid_ingredients.description as valid_ingredient_description,
	valid_ingredients.warning as valid_ingredient_warning,
	valid_ingredients.contains_egg as valid_ingredient_contains_egg,
	valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
	valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
	valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
	valid_ingredients.contains_soy as valid_ingredient_contains_soy,
	valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
	valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
	valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
	valid_ingredients.contains_fish as valid_ingredient_contains_fish,
	valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
	valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
	valid_ingredients.volumetric as valid_ingredient_volumetric,
	valid_ingredients.is_liquid as valid_ingredient_is_liquid,
	valid_ingredients.icon_path as valid_ingredient_icon_path,
	valid_ingredients.animal_derived as valid_ingredient_animal_derived,
	valid_ingredients.plural_name as valid_ingredient_plural_name,
	valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
	valid_ingredients.slug as valid_ingredient_slug,
	valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
	valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
	valid_ingredients.is_starch as valid_ingredient_is_starch,
	valid_ingredients.is_protein as valid_ingredient_is_protein,
	valid_ingredients.is_grain as valid_ingredient_is_grain,
	valid_ingredients.is_fruit as valid_ingredient_is_fruit,
	valid_ingredients.is_salt as valid_ingredient_is_salt,
	valid_ingredients.is_fat as valid_ingredient_is_fat,
	valid_ingredients.is_acid as valid_ingredient_is_acid,
	valid_ingredients.is_heat as valid_ingredient_is_heat,
	valid_ingredients.created_at as valid_ingredient_created_at,
	valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
	valid_ingredients.archived_at as valid_ingredient_archived_at,
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
	meal_plan_grocery_list_items.minimum_quantity_needed,
	meal_plan_grocery_list_items.maximum_quantity_needed,
	meal_plan_grocery_list_items.quantity_purchased,
	meal_plan_grocery_list_items.purchased_measurement_unit,
	meal_plan_grocery_list_items.purchased_upc,
	meal_plan_grocery_list_items.purchase_price,
	meal_plan_grocery_list_items.status_explanation,
	meal_plan_grocery_list_items.status,
	meal_plan_grocery_list_items.created_at,
	meal_plan_grocery_list_items.last_updated_at,
	meal_plan_grocery_list_items.archived_at
FROM meal_plan_grocery_list_items
	JOIN meal_plans ON meal_plan_grocery_list_items.belongs_to_meal_plan=meal_plans.id
	JOIN valid_ingredients ON meal_plan_grocery_list_items.valid_ingredient=valid_ingredients.id
	JOIN valid_measurement_units ON meal_plan_grocery_list_items.valid_measurement_unit=valid_measurement_units.id
WHERE meal_plan_grocery_list_items.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND meal_plan_grocery_list_items.id = sqlc.arg(meal_plan_grocery_list_item_id)
	AND meal_plan_grocery_list_items.belongs_to_meal_plan = sqlc.arg(meal_plan_id);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateMealPlanGroceryListItem",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				mealPlanGroceryListItemsTableName,
				strings.Join(applyToEach(filterForUpdate(mealPlanGroceryListItemsColumns, belongsToUserColumn), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				idColumn, idColumn,
			)),
		},
	}
}
