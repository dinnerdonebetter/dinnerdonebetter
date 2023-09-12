package main

const oauth2ClientsTableName = "oauth2_clients"

var oauth2ClientsColumns = []string{
	"id",
	"name",
	"description",
	"client_id",
	"client_secret",
	createdAtColumn,
	archivedAtColumn,
}
