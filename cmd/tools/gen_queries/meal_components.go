package main

const mealComponentsTableName = "meal_components"

var mealComponentsColumns = []string{
	idColumn,
	"meal_id",
	"recipe_id",
	"meal_component_type",
	"recipe_scale",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealComponentsQueries() []*Query {
	return []*Query{
		//
	}
}
