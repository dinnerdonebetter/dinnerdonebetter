package main

import (
	"github.com/cristalhq/builq"
)

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
		{
			Annotation: QueryAnnotation{
				Name: "",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}
