package main

import (
	"github.com/cristalhq/builq"
)

const validInstrumentsTableName = "valid_instruments"

var validInstrumentsColumns = []string{
	idColumn,
	"name",
	"description",
	"icon_path",
	"plural_name",
	"usable_for_storage",
	"slug",
	"display_in_summary_lists",
	"include_in_generated_instructions",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidInstrumentsQueries() []*Query {
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
