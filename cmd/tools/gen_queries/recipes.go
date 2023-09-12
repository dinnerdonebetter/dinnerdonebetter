package main

const recipesTableName = "recipes"

var recipesColumns = []string{
	"id",
	"name",
	"source",
	"description",
	"inspired_by_recipe_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	"created_by_user",
	"min_estimated_portions",
	"seal_of_approval",
	"slug",
	"portion_name",
	"plural_portion_name",
	"max_estimated_portions",
	"eligible_for_meals",
	"last_indexed_at",
	"last_validated_at",
	"yields_component_type",
}
