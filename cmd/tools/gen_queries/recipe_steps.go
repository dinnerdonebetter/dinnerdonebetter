package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeStepsTableName = "recipe_steps"

	indexColumn               = "index"
	recipeStepIDColumn        = "recipe_step_id"
	belongsToRecipeStepColumn = "belongs_to_recipe_step"
)

var recipeStepsColumns = []string{
	idColumn,
	indexColumn,
	preparationIDColumn,
	"minimum_estimated_time_in_seconds",
	"maximum_estimated_time_in_seconds",
	"minimum_temperature_in_celsius",
	"maximum_temperature_in_celsius",
	notesColumn,
	"explicit_instructions",
	"condition_expression",
	"optional",
	"start_timer_automatically",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	belongsToRecipeColumn,
}

func buildRecipeStepsQueries() []*Query {
	insertColumns := filterForInsert(recipeStepsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveRecipeStep",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				recipeStepsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToRecipeColumn,
				belongsToRecipeColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateRecipeStep",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipeStepsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckRecipeStepExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
		JOIN %s ON %s.%s=%s.%s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				recipeStepsTableName, idColumn,
				recipeStepsTableName,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStep",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id as valid_preparation_id,
	valid_preparations.name as valid_preparation_name,
	valid_preparations.description as valid_preparation_description,
	valid_preparations.icon_path as valid_preparation_icon_path,
	valid_preparations.yields_nothing as valid_preparation_yields_nothing,
	valid_preparations.restrict_to_ingredients as valid_preparation_restrict_to_ingredients,
	valid_preparations.minimum_ingredient_count as valid_preparation_minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count as valid_preparation_maximum_ingredient_count,
	valid_preparations.minimum_instrument_count as valid_preparation_minimum_instrument_count,
	valid_preparations.maximum_instrument_count as valid_preparation_maximum_instrument_count,
	valid_preparations.temperature_required as valid_preparation_temperature_required,
	valid_preparations.time_estimate_required as valid_preparation_time_estimate_required,
	valid_preparations.condition_expression_required as valid_preparation_condition_expression_required,
	valid_preparations.consumes_vessel as valid_preparation_consumes_vessel,
	valid_preparations.only_for_vessels as valid_preparation_only_for_vessels,
	valid_preparations.minimum_vessel_count as valid_preparation_minimum_vessel_count,
	valid_preparations.maximum_vessel_count as valid_preparation_maximum_vessel_count,
	valid_preparations.slug as valid_preparation_slug,
	valid_preparations.past_tense as valid_preparation_past_tense,
	valid_preparations.created_at as valid_preparation_created_at,
	valid_preparations.last_updated_at as valid_preparation_last_updated_at,
	valid_preparations.archived_at as valid_preparation_archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.condition_expression,
	recipe_steps.optional,
	recipe_steps.start_timer_automatically,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe
FROM recipe_steps
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipe_steps.id = sqlc.arg(recipe_step_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeSteps",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id as valid_preparation_id,
	valid_preparations.name as valid_preparation_name,
	valid_preparations.description as valid_preparation_description,
	valid_preparations.icon_path as valid_preparation_icon_path,
	valid_preparations.yields_nothing as valid_preparation_yields_nothing,
	valid_preparations.restrict_to_ingredients as valid_preparation_restrict_to_ingredients,
	valid_preparations.minimum_ingredient_count as valid_preparation_minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count as valid_preparation_maximum_ingredient_count,
	valid_preparations.minimum_instrument_count as valid_preparation_minimum_instrument_count,
	valid_preparations.maximum_instrument_count as valid_preparation_maximum_instrument_count,
	valid_preparations.temperature_required as valid_preparation_temperature_required,
	valid_preparations.time_estimate_required as valid_preparation_time_estimate_required,
	valid_preparations.condition_expression_required as valid_preparation_condition_expression_required,
	valid_preparations.consumes_vessel as valid_preparation_consumes_vessel,
	valid_preparations.only_for_vessels as valid_preparation_only_for_vessels,
	valid_preparations.minimum_vessel_count as valid_preparation_minimum_vessel_count,
	valid_preparations.maximum_vessel_count as valid_preparation_maximum_vessel_count,
	valid_preparations.slug as valid_preparation_slug,
	valid_preparations.past_tense as valid_preparation_past_tense,
	valid_preparations.created_at as valid_preparation_created_at,
	valid_preparations.last_updated_at as valid_preparation_last_updated_at,
	valid_preparations.archived_at as valid_preparation_archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.condition_expression,
	recipe_steps.optional,
	recipe_steps.start_timer_automatically,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe,
	(
		SELECT
			COUNT(recipe_steps.id)
		FROM
			recipe_steps
		WHERE
			recipe_steps.archived_at IS NULL
			AND recipe_steps.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
			AND recipe_steps.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
			AND (recipe_steps.last_updated_at IS NULL OR recipe_steps.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
			AND (recipe_steps.last_updated_at IS NULL OR recipe_steps.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
	) as filtered_count,
	(
		SELECT COUNT(recipe_steps.id)
		FROM recipe_steps
		WHERE recipe_steps.archived_at IS NULL
	) as total_count
FROM recipe_steps
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_at IS NULL
	AND recipe_steps.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
	AND recipe_steps.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
	AND (recipe_steps.last_updated_at IS NULL OR recipe_steps.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
	AND (recipe_steps.last_updated_at IS NULL OR recipe_steps.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipes.archived_at IS NULL
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepByRecipeID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id as valid_preparation_id,
	valid_preparations.name as valid_preparation_name,
	valid_preparations.description as valid_preparation_description,
	valid_preparations.icon_path as valid_preparation_icon_path,
	valid_preparations.yields_nothing as valid_preparation_yields_nothing,
	valid_preparations.restrict_to_ingredients as valid_preparation_restrict_to_ingredients,
	valid_preparations.minimum_ingredient_count as valid_preparation_minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count as valid_preparation_maximum_ingredient_count,
	valid_preparations.minimum_instrument_count as valid_preparation_minimum_instrument_count,
	valid_preparations.maximum_instrument_count as valid_preparation_maximum_instrument_count,
	valid_preparations.temperature_required as valid_preparation_temperature_required,
	valid_preparations.time_estimate_required as valid_preparation_time_estimate_required,
	valid_preparations.condition_expression_required as valid_preparation_condition_expression_required,
	valid_preparations.consumes_vessel as valid_preparation_consumes_vessel,
	valid_preparations.only_for_vessels as valid_preparation_only_for_vessels,
	valid_preparations.minimum_vessel_count as valid_preparation_minimum_vessel_count,
	valid_preparations.maximum_vessel_count as valid_preparation_maximum_vessel_count,
	valid_preparations.slug as valid_preparation_slug,
	valid_preparations.past_tense as valid_preparation_past_tense,
	valid_preparations.created_at as valid_preparation_created_at,
	valid_preparations.last_updated_at as valid_preparation_last_updated_at,
	valid_preparations.archived_at as valid_preparation_archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.condition_expression,
	recipe_steps.optional,
	recipe_steps.start_timer_automatically,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe
FROM recipe_steps
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_at IS NULL
	AND recipe_steps.id = sqlc.arg(id);`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipeStep",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				recipeStepsTableName,
				strings.Join(applyToEach(filterForUpdate(recipeStepsColumns, belongsToRecipeColumn), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				belongsToRecipeColumn, belongsToRecipeColumn,
				idColumn, idColumn,
			)),
		},
	}
}
