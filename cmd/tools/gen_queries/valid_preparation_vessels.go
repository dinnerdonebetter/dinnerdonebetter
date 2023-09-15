package main

import (
	"github.com/cristalhq/builq"
)

const validPreparationVesselsTableName = "valid_preparation_vessels"

var validPreparationVesselsColumns = []string{
	idColumn,
	"valid_preparation_id",
	"valid_vessel_id",
	"notes",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidPreparationVesselsQueries() []*Query {
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
