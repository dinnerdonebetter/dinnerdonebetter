package main

const serviceSettingsTableName = "service_settings"

var serviceSettingsColumns = []string{
	"id",
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
