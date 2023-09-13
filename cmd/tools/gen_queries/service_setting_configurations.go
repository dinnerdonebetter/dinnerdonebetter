package main

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
