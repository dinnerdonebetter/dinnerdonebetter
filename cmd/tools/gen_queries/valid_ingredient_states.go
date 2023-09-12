package main

const validIngredientStatesTableName = "valid_ingredient_states"

var validIngredientStatesColumns = []string{
	"id",
	"name",
	"past_tense",
	"slug",
	"description",
	"icon_path",
	"attribute_type",
	"last_indexed_at",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
