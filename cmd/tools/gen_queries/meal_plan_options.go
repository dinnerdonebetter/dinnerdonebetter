package main

const mealPlanOptionsTableName = "meal_plan_options"

var mealPlanOptionsColumns = []string{
	"id",
	"meal_id",
	"notes",
	"chosen",
	"tiebroken",
	"assigned_cook",
	"assigned_dishwasher",
	"belongs_to_meal_plan_event",
	"meal_scale",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
