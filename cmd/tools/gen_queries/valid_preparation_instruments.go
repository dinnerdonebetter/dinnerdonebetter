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
				Name: "ArchiveValidPreparationInstrument",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidPreparationInstrument",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidPreparationInstrumentExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstrumentsForInstrument",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstrumentsForPreparation",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstruments",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstrument",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ValidPreparationInstrumentPairIsValid",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidPreparationInstrument",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}
