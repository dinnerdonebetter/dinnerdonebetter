package main

const mealPlansTableName = "meal_plans"

var mealPlansColumns = []string{
	"id",
	"notes",
	"status",
	"voting_deadline",
	belongsToHouseholdColumn,
	"grocery_list_initialized",
	"tasks_created",
	"election_method",
	"created_by_user",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
