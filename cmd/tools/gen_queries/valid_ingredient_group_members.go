package main

const validIngredientGroupMembersTableName = "valid_ingredient_group_members"

var validIngredientGroupMembersColumns = []string{
	"id",
	"belongs_to_group",
	"valid_ingredient",
	createdAtColumn,
	archivedAtColumn,
}
