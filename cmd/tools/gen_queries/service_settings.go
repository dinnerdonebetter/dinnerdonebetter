package main

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
