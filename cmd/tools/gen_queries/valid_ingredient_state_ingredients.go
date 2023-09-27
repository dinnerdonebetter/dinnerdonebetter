package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validIngredientStateIngredientsTableName = "valid_ingredient_state_ingredients"
)

var validIngredientStateIngredientsColumns = []string{
	idColumn,
	notesColumn,
	"valid_ingredient_state",
	"valid_ingredient",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientStateIngredientsQueries() []*Query {
	insertColumns := filterForInsert(validIngredientStateIngredientsColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(filterFromSlice(validIngredientStateIngredientsColumns, "valid_ingredient_id", "valid_measurement_unit_id"), func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_ingredient_state_ingredient_%s", validIngredientStateIngredientsTableName, s, s)
		}),
		append(
			applyToEach(validIngredientStatesColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_state_%s", validIngredientStatesTableName, s, s)
			}),
			applyToEach(validIngredientsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
			})...),
		2,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidIngredientStateIngredient",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validIngredientStateIngredientsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidIngredientStateIngredient",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				validIngredientStateIngredientsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidIngredientStateIngredientExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				validIngredientStateIngredientsTableName, idColumn,
				validIngredientStateIngredientsTableName,
				validIngredientStateIngredientsTableName, archivedAtColumn,
				validIngredientStateIngredientsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredientsForIngredient",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN valid_ingredients ON valid_ingredient_state_ingredients.valid_ingredient = valid_ingredients.id
	JOIN valid_ingredient_states ON valid_ingredient_state_ingredients.valid_ingredient_state = valid_ingredient_states.id
WHERE
	valid_ingredient_state_ingredients.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_state_ingredients.valid_ingredient = sqlc.arg(valid_ingredient)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					validIngredientStateIngredientsTableName,
					true,
				),
				buildTotalCountSelect(
					validIngredientStateIngredientsTableName,
				),
				validIngredientStateIngredientsTableName,
				buildFilterConditions(validIngredientStateIngredientsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredientsForIngredientState",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN valid_ingredients ON valid_ingredient_state_ingredients.valid_ingredient = valid_ingredients.id
	JOIN valid_ingredient_states ON valid_ingredient_state_ingredients.valid_ingredient_state = valid_ingredient_states.id
WHERE
	valid_ingredient_state_ingredients.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_state_ingredients.valid_ingredient_state = sqlc.arg(valid_ingredient_state)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					validIngredientStateIngredientsTableName,
					true,
				),
				buildTotalCountSelect(
					validIngredientStateIngredientsTableName,
				),
				validIngredientStateIngredientsTableName,
				buildFilterConditions(validIngredientStateIngredientsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredients",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN valid_ingredients ON valid_ingredient_state_ingredients.valid_ingredient = valid_ingredients.id
	JOIN valid_ingredient_states ON valid_ingredient_state_ingredients.valid_ingredient_state = valid_ingredient_states.id
WHERE
	valid_ingredient_state_ingredients.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_ingredient_states.archived_at IS NULL
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					validIngredientStateIngredientsTableName,
					true,
				),
				buildTotalCountSelect(
					validIngredientStateIngredientsTableName,
				),
				validIngredientStateIngredientsTableName,
				buildFilterConditions(validIngredientStateIngredientsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredient",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN valid_ingredients ON valid_ingredient_state_ingredients.valid_ingredient = valid_ingredients.id
	JOIN valid_ingredient_states ON valid_ingredient_state_ingredients.valid_ingredient_state = valid_ingredient_states.id
WHERE
	valid_ingredient_state_ingredients.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_state_ingredients.id = sqlc.arg(id);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validIngredientStateIngredientsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredientsWithIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN valid_ingredients ON valid_ingredient_state_ingredients.valid_ingredient = valid_ingredients.id
	JOIN valid_ingredient_states ON valid_ingredient_state_ingredients.valid_ingredient_state = valid_ingredient_states.id
WHERE
	valid_ingredient_state_ingredients.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_state_ingredients.id = ANY(sqlc.arg(ids)::text[]);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validIngredientStateIngredientsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidityOfValidIngredientStateIngredientPair",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s
	FROM %s
	WHERE valid_ingredient = sqlc.arg(valid_ingredient)
	AND valid_ingredient_state = sqlc.arg(valid_ingredient_state)
	AND %s IS NULL
);`,
				idColumn,
				validIngredientStateIngredientsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidIngredientStateIngredient",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validIngredientStateIngredientsTableName,
				strings.Join(applyToEach(filterForUpdate(validIngredientStateIngredientsColumns), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
	}
}
