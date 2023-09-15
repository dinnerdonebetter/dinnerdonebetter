package main

const mealPlansTableName = "meal_plans"

var mealPlansColumns = []string{
	idColumn,
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

func buildMealPlansQueries() []*Query {
	return []*Query{
		//
	}
}
