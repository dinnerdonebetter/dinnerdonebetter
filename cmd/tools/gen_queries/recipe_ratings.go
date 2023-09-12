package main

const recipeRatingsTableName = "recipe_ratings"

var recipeRatingsColumns = []string{
	"id",
	"recipe_id",
	"taste",
	"difficulty",
	"cleanup",
	"instructions",
	"overall",
	"notes",
	"by_user",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
