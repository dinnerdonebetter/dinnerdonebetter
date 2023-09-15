package main

const oauth2ClientsTableName = "oauth2_clients"

var oauth2ClientsColumns = []string{
	idColumn,
	"name",
	"description",
	"client_id",
	"client_secret",
	createdAtColumn,
	archivedAtColumn,
}

func buildOAuth2ClientsQueries() []*Query {
	return []*Query{
		//
	}
}
