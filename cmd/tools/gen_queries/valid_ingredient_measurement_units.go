package main

import (
	"github.com/cristalhq/builq"
)

const validIngredientMeasurementUnitsTableName = "valid_ingredient_measurement_units"

var validIngredientMeasurementUnitsColumns = []string{
	idColumn,
	"notes",
	"valid_ingredient_id",
	"valid_measurement_unit_id",
	"minimum_allowable_quantity",
	"maximum_allowable_quantity",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientMeasurementUnitsQueries() []*Query {
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