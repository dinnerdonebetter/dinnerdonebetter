package main

import (
	"github.com/cristalhq/builq"
)

const validPreparationsTableName = "valid_preparations"

var validPreparationsColumns = []string{
	idColumn,
	"name",
	"description",
	"icon_path",
	"yields_nothing",
	"restrict_to_ingredients",
	"past_tense",
	"slug",
	"minimum_ingredient_count",
	"maximum_ingredient_count",
	"minimum_instrument_count",
	"maximum_instrument_count",
	"temperature_required",
	"time_estimate_required",
	"condition_expression_required",
	"consumes_vessel",
	"only_for_vessels",
	"minimum_vessel_count",
	"maximum_vessel_count",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidPreparationsQueries() []*Query {
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
