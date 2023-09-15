package main

import (
	"github.com/cristalhq/builq"
)

const validPreparationInstrumentsTableName = "valid_preparation_instruments"

var validPreparationInstrumentsColumns = []string{
	idColumn,
	"notes",
	"valid_preparation_id",
	"valid_instrument_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidPreparationInstrumentsQueries() []*Query {
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
