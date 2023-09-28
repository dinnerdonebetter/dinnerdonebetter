package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipeStepProductsTableName = "recipe_step_products"
)

var recipeStepProductsColumns = []string{
	idColumn,
	nameColumn,
	belongsToRecipeStepColumn,
	"quantity_notes",
	"minimum_quantity_value",
	"maximum_quantity_value",
	"measurement_unit",
	"type",
	"compostable",
	"maximum_storage_duration_in_seconds",
	"minimum_storage_temperature_in_celsius",
	"maximum_storage_temperature_in_celsius",
	"storage_instructions",
	"is_liquid",
	"is_waste",
	"index",
	"contained_in_vessel_index",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepProductsQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "",
				Type: "",
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}
