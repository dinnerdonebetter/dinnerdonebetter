package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	userIngredientPreferencesTableName = "user_ingredient_preferences"
)

var userIngredientPreferencesColumns = []string{
	idColumn,
	"ingredient",
	"rating",
	notesColumn,
	"allergy",
	belongsToUserColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildUserIngredientPreferencesQueries() []*Query {
	insertColumns := filterForInsert(userIngredientPreferencesColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveUserIngredientPreference",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				userIngredientPreferencesTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToUserColumn,
				belongsToUserColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateUserIngredientPreference",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				userIngredientPreferencesTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckUserIngredientPreferenceExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
);`,
				userIngredientPreferencesTableName, idColumn,
				userIngredientPreferencesTableName,
				userIngredientPreferencesTableName, archivedAtColumn,
				userIngredientPreferencesTableName, idColumn, idColumn,
				userIngredientPreferencesTableName, belongsToUserColumn, belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserIngredientPreferencesForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	user_ingredient_preferences.id,
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
	user_ingredient_preferences.rating,
	user_ingredient_preferences.notes,
	user_ingredient_preferences.allergy,
	user_ingredient_preferences.created_at,
	user_ingredient_preferences.last_updated_at,
	user_ingredient_preferences.archived_at,
	user_ingredient_preferences.belongs_to_user,
	(
		SELECT
			COUNT(user_ingredient_preferences.id)
		FROM
			user_ingredient_preferences
		WHERE
			user_ingredient_preferences.archived_at IS NULL
			AND user_ingredient_preferences.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
			AND user_ingredient_preferences.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
			AND (
				user_ingredient_preferences.last_updated_at IS NULL
				OR user_ingredient_preferences.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
			)
			AND (
				user_ingredient_preferences.last_updated_at IS NULL
				OR user_ingredient_preferences.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
			)
	) AS filtered_count,
	(
		SELECT
			COUNT(user_ingredient_preferences.id)
		FROM
			user_ingredient_preferences
		WHERE
			user_ingredient_preferences.archived_at IS NULL
	) AS total_count
FROM user_ingredient_preferences
	JOIN valid_ingredients ON valid_ingredients.id = user_ingredient_preferences.ingredient
WHERE user_ingredient_preferences.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND user_ingredient_preferences.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
	AND user_ingredient_preferences.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
	AND (
		user_ingredient_preferences.last_updated_at IS NULL
		OR user_ingredient_preferences.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
	)
	AND (
		user_ingredient_preferences.last_updated_at IS NULL
		OR user_ingredient_preferences.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
	)
	AND user_ingredient_preferences.belongs_to_user = sqlc.arg(user_id)
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserIngredientPreference",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	user_ingredient_preferences.id,
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
	user_ingredient_preferences.rating,
	user_ingredient_preferences.notes,
	user_ingredient_preferences.allergy,
	user_ingredient_preferences.created_at,
	user_ingredient_preferences.last_updated_at,
	user_ingredient_preferences.archived_at,
	user_ingredient_preferences.belongs_to_user
FROM user_ingredient_preferences
	JOIN valid_ingredients ON valid_ingredients.id = user_ingredient_preferences.ingredient
WHERE user_ingredient_preferences.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND user_ingredient_preferences.id = sqlc.arg(user_ingredient_preference_id)
	AND user_ingredient_preferences.belongs_to_user = sqlc.arg(user_id);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserIngredientPreference",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				userIngredientPreferencesTableName,
				strings.Join(applyToEach(filterForUpdate(userIngredientPreferencesColumns, belongsToUserColumn), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				belongsToUserColumn, belongsToUserColumn,
				idColumn, idColumn,
			)),
		},
	}
}
