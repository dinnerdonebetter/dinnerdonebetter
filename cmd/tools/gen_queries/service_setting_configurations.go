package main

import (
	"github.com/cristalhq/builq"
)

const serviceSettingConfigurationsTableName = "service_setting_configurations"

var serviceSettingConfigurationsColumns = []string{
	idColumn,
	"value",
	"notes",
	"service_setting_id",
	"belongs_to_user",
	belongsToHouseholdColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildServiceSettingConfigurationQueries() []*Query {
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
