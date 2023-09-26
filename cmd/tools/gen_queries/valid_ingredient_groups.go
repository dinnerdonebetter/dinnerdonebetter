package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validIngredientGroupsTableName = "valid_ingredient_groups"
)

var validIngredientGroupsColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	slugColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientGroupsQueries() []*Query {
	groupInsertColumns := filterForInsert(validIngredientGroupsColumns)
	groupMemberInsertColumns := filterForInsert(validIngredientGroupMembersColumns)

	fullMemberSelectColumns := mergeColumns(
		applyToEach(filterFromSlice(validIngredientGroupMembersColumns, "valid_ingredient"), func(i int, s string) string {
			return fmt.Sprintf("%s.%s", validIngredientGroupMembersTableName, s)
		}),
		applyToEach(validIngredientsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
		}),
		2,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidIngredientGroup",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validIngredientGroupsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidIngredientGroupMember",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND belongs_to_group = sqlc.arg(group_id);`,
				validIngredientGroupMembersTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidIngredientGroup",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validIngredientGroupsTableName,
				strings.Join(groupInsertColumns, ",\n\t"),
				strings.Join(applyToEach(groupInsertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidIngredientGroupMember",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validIngredientGroupMembersTableName,
				strings.Join(groupMemberInsertColumns, ",\n\t"),
				strings.Join(applyToEach(groupMemberInsertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidIngredientGroupExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				validIngredientGroupsTableName,
				validIngredientGroupsTableName,
				validIngredientGroupsTableName,
				archivedAtColumn,
				validIngredientGroupsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientGroups",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
    %s,
    %s
FROM %s
WHERE
	%s.%s IS NULL
	%s
GROUP BY %s.%s
ORDER BY %s.%s
%s;`,
				strings.Join(applyToEach(validIngredientGroupsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientGroupsTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(
					validIngredientGroupsTableName,
					true,
				),
				buildTotalCountSelect(
					validIngredientGroupsTableName,
				),
				validIngredientGroupsTableName,
				validIngredientGroupsTableName,
				archivedAtColumn,
				buildFilterConditions(
					validIngredientGroupsTableName,
					true,
				),
				validIngredientGroupsTableName,
				idColumn,
				validIngredientGroupsTableName,
				idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientGroupMembers",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN valid_ingredient_groups ON valid_ingredient_groups.id = valid_ingredient_group_members.belongs_to_group
	JOIN valid_ingredients ON valid_ingredients.id = valid_ingredient_group_members.valid_ingredient
WHERE 
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.belongs_to_group = sqlc.arg(group_id);`,
				strings.Join(fullMemberSelectColumns, ",\n\t"),
				validIngredientGroupMembersTableName,
				validIngredientGroupsTableName,
				archivedAtColumn,
				validIngredientGroupMembersTableName,
				archivedAtColumn,
				validIngredientGroupMembersTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientGroup",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(validIngredientGroupsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientGroupsTableName, s)
				}), ",\n\t"),
				validIngredientGroupsTableName,
				validIngredientGroupsTableName,
				archivedAtColumn,
				validIngredientGroupsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchForValidIngredientGroups",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
    %s,
    %s
FROM %s
WHERE
	%s.%s IS NULL
	AND valid_ingredient_groups.name %s
	%s
GROUP BY %s.%s
ORDER BY %s.%s
%s;`,
				strings.Join(applyToEach(validIngredientGroupsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientGroupsTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(
					validIngredientGroupsTableName,
					true,
				),
				buildTotalCountSelect(
					validIngredientGroupsTableName,
				),
				validIngredientGroupsTableName,
				validIngredientGroupsTableName,
				archivedAtColumn,
				"ILIKE '%' || sqlc.arg(name)::text || '%'",
				buildFilterConditions(
					validIngredientGroupsTableName,
					true,
				),
				validIngredientGroupsTableName,
				idColumn,
				validIngredientGroupsTableName,
				idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientGroupsWithIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[]);`,
				strings.Join(applyToEach(validIngredientGroupsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientGroupsTableName, s)
				}), ",\n\t"),
				validIngredientGroupsTableName,
				validIngredientGroupsTableName,
				archivedAtColumn,
				validIngredientGroupsTableName,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidIngredientGroup",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
    AND %s = sqlc.arg(%s);`,
				validIngredientGroupsTableName,
				strings.Join(applyToEach(filterForUpdate(validIngredientGroupsColumns), func(i int, s string) string {
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
