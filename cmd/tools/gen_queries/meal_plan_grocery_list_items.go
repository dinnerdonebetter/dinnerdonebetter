package main

import (
	"github.com/cristalhq/builq"
)

const mealPlanGroceryListItemsTableName = "meal_plan_grocery_list_items"

var mealPlanGroceryListItemsColumns = []string{
	idColumn,
	"valid_ingredient",
	"valid_measurement_unit",
	"minimum_quantity_needed",
	"maximum_quantity_needed",
	"quantity_purchased",
	"purchased_measurement_unit",
	"purchased_upc",
	"purchase_price",
	"status_explanation",
	"status",
	"belongs_to_meal_plan",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealPlanGroceryListItemsQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}
