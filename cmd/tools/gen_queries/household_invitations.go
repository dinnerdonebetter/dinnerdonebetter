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
				Name: "AttachHouseholdInvitationsToUserID",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateHouseholdInvitation",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckHouseholdInvitationExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdInvitationByEmailAndToken",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdInvitationByHouseholdAndID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdInvitationByTokenAndID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetPendingInvitesFromUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetPendingInvitesForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SetHouseholdInvitationStatus",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}
