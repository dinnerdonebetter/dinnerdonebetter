package main

import (
	"github.com/cristalhq/builq"
)

const householdInvitationsTableName = "household_invitations"

var householdInvitationsColumns = []string{
	idColumn,
	"destination_household",
	"to_email",
	"to_user",
	"from_user",
	"status",
	"note",
	"status_note",
	"token",
	"expires_at",
	"to_name",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildHouseholdInvitationsQueries() []*Query {
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
