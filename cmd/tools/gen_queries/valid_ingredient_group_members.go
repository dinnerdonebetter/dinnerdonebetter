package main

const validIngredientGroupMembersTableName = "valid_ingredient_group_members"

var validIngredientGroupMembersColumns = []string{
	idColumn,
	"belongs_to_group",
	"valid_ingredient",
	createdAtColumn,
	archivedAtColumn,
}
