package main

import (
	"github.com/cristalhq/builq"
)

const (
	recipePrepTasksTableName = "recipe_prep_tasks"
)

var recipePrepTasksColumns = []string{
	idColumn,
	notesColumn,
	"explicit_storage_instructions",
	"minimum_time_buffer_before_recipe_in_seconds",
	"maximum_time_buffer_before_recipe_in_seconds",
	"storage_type",
	"minimum_storage_temperature_in_celsius",
	"maximum_storage_temperature_in_celsius",
	"belongs_to_recipe",
	nameColumn,
	descriptionColumn,
	"optional",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipePrepTasksQueries() []*Query {
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
