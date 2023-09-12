package main

const mealPlanEventsTableName = "meal_plan_events"

var mealPlanEventsColumns = []string{
	"id",
	"notes",
	"starts_at",
	"ends_at",
	"meal_name",
	"belongs_to_meal_plan",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
