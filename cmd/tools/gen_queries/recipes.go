package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipesTableName = "recipes"

	belongsToRecipeColumn = "belongs_to_recipe"
	recipeIDColumn        = "recipe_id"
	lastValidatedAtColumn = "last_validated_at"
)

var recipesColumns = []string{
	idColumn,
	nameColumn,
	slugColumn,
	"source",
	descriptionColumn,
	"inspired_by_recipe_id",
	"min_estimated_portions",
	"max_estimated_portions",
	"portion_name",
	"plural_portion_name",
	"seal_of_approval",
	"eligible_for_meals",
	"yields_component_type",
	lastIndexedAtColumn,
	lastValidatedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	createdByUserColumn,
}

func buildRecipesQueries() []*Query {
	insertColumns := filterForInsert(recipesColumns, lastValidatedAtColumn)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveRecipe",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				recipesTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				createdByUserColumn,
				createdByUserColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateRecipe",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipesTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckRecipeExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.%s
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				recipesTableName, idColumn,
				recipesTableName,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeByID",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeByIDAndAuthorID",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipes",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "RecipeSearch",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipesNeedingIndexing",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT %s.%s
FROM %s
WHERE %s.%s IS NULL
	AND (
		%s.%s IS NULL
		OR %s.%s < %s - '24 hours'::INTERVAL
	);`,
				recipesTableName, idColumn,
				recipesTableName,
				recipesTableName, archivedAtColumn,
				recipesTableName, lastIndexedAtColumn,
				recipesTableName, lastIndexedAtColumn, currentTimeExpression,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeIDsForMeal",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT %s.%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
GROUP BY %s.%s
ORDER BY %s.%s;`,
				recipesTableName, idColumn,
				recipesTableName,
				mealComponentsTableName, mealComponentsTableName, recipeIDColumn, recipesTableName, idColumn,
				mealsTableName, mealComponentsTableName, mealIDColumn, mealsTableName, idColumn,
				recipesTableName, archivedAtColumn,
				mealsTableName, idColumn, mealIDColumn,
				recipesTableName, idColumn,
				recipesTableName, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipe",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
    AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				recipesTableName,
				strings.Join(applyToEach(filterForUpdate(recipesColumns, lastValidatedAtColumn, createdByUserColumn), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				createdByUserColumn, createdByUserColumn,
				idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipeLastIndexedAt",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s = sqlc.arg(%s) AND %s IS NULL;`,
				recipesTableName,
				lastIndexedAtColumn,
				currentTimeExpression,
				idColumn,
				idColumn,
				archivedAtColumn,
			)),
		},
	}
}
