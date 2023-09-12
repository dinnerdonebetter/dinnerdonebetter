package main

const mealPlanOptionVotesTableName = "meal_plan_option_votes"

var mealPlanOptionVotesColumns = []string{
	"id",
	"rank",
	"abstain",
	"notes",
	"by_user",
	"belongs_to_meal_plan_option",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
