package main

import (
	"github.com/cristalhq/builq"
)

const validMeasurementUnitsTableName = "valid_measurement_units"

var validMeasurementUnitsColumns = []string{
	idColumn,
	"name",
	"description",
	"icon_path",
	"volumetric",
	"universal",
	"metric",
	"imperial",
	"plural_name",
	"slug",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidMeasurementUnitsQueries() []*Query {
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
