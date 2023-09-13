package main

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
