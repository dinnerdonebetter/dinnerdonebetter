package main

import (
	"github.com/cristalhq/builq"
)

const serviceSettingsTableName = "service_settings"

var serviceSettingsColumns = []string{
	idColumn,
	"name",
	"type",
	"description",
	"default_value",
	"enumeration",
	"admins_only",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildServiceSettingQueries() []*Query {
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
