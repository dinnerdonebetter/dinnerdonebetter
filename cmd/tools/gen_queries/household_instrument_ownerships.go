package main

import (
	"github.com/cristalhq/builq"
)

const householdInstrumentOwnershipsTableName = "household_instrument_ownerships"

var householdInstrumentOwnershipsColumns = []string{
	idColumn,
	"notes",
	"quantity",
	"valid_instrument_id",
	belongsToHouseholdColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildHouseholdInstrumentOwnershipQueries() []*Query {
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
