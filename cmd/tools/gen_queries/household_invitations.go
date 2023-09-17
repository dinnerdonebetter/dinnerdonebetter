package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const householdInvitationsTableName = "household_invitations"

var householdInvitationsColumns = []string{
	idColumn,
	"from_user",
	"to_user",
	"to_name",
	"note",
	"to_email",
	"token",
	"destination_household",
	"expires_at",
	"status",
	"status_note",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildHouseholdInvitationsQueries() []*Query {
	insertColumns := filterForInsert(householdInvitationsColumns,
		"status",
		"status_note",
	)

	fullSelectColumns := mergeColumns(mergeColumns(
		applyToEach(householdInvitationsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s", householdInvitationsTableName, s)
		}),
		applyToEach(usersColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as user_%s", usersTableName, s, s)
		}),
		3,
	),
		applyToEach(householdsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as household_%s", householdsTableName, s, s)
		}),
		1,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "AttachHouseholdInvitationsToUserID",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	to_user = sqlc.arg(user_id),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND to_email = LOWER(sqlc.arg(email_address));`,
				householdInvitationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateHouseholdInvitation",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				householdInvitationsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckHouseholdInvitationExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.id
	FROM %s
	WHERE %s.archived_at IS NULL
	AND %s.id = $1
);`,
				householdInvitationsTableName,
				householdInvitationsTableName,
				householdInvitationsTableName,
				householdInvitationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdInvitationByEmailAndToken",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN households ON household_invitations.destination_household = households.id
	JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
	AND household_invitations.expires_at > NOW()
	AND household_invitations.to_email = LOWER(sqlc.arg(email_address))
	AND household_invitations.token = sqlc.arg(token);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				householdInvitationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdInvitationByHouseholdAndID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN households ON household_invitations.destination_household = households.id
	JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
	AND household_invitations.expires_at > NOW()
	AND household_invitations.expires_at > NOW()
	AND household_invitations.destination_household = $1
	AND household_invitations.id = $2;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				householdInvitationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdInvitationByTokenAndID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN households ON household_invitations.destination_household = households.id
	JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
	AND household_invitations.expires_at > NOW()
	AND household_invitations.token = $1
	AND household_invitations.id = $2;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				householdInvitationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetPendingInvitesFromUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN households ON household_invitations.destination_household = households.id
	JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
	AND household_invitations.from_user = sqlc.arg(user_id)
	AND household_invitations.status = sqlc.arg(status)
	%s
	OFFSET sqlc.narg(query_offset)
	LIMIT sqlc.narg(query_limit);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					householdInvitationsTableName,
					true,
				),
				buildTotalCountSelect(
					householdInvitationsTableName,
				),
				householdInvitationsTableName,
				buildFilterConditions(
					householdInvitationsTableName,
					true,
				),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetPendingInvitesForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN households ON household_invitations.destination_household = households.id
	JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
	AND household_invitations.to_user = sqlc.arg(user_id)
	AND household_invitations.status = sqlc.arg(status)
	%s
	OFFSET sqlc.narg(query_offset)
	LIMIT sqlc.narg(query_limit);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					householdInvitationsTableName,
					true,
				),
				buildTotalCountSelect(
					householdInvitationsTableName,
				),
				householdInvitationsTableName,
				buildFilterConditions(
					householdInvitationsTableName,
					true,
				),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SetHouseholdInvitationStatus",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	status = sqlc.arg(status),
	status_note = sqlc.arg(status_note),
	last_updated_at = NOW(),
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);`,
				householdInvitationsTableName,
			)),
		},
	}
}
