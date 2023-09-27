package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	householdInvitationsTableName = "household_invitations"
	destinationHouseholdColumn    = "destination_household"
	fromUserColumn                = "from_user"
	toUserColumn                  = "to_user"
	toEmailColumn                 = "to_email"
	tokenColumn                   = "token"
	statusColumn                  = "status"
	statusNoteColumn              = "status_note"
	expiresAtColumn               = "expires_at"
)

var householdInvitationsColumns = []string{
	idColumn,
	fromUserColumn,
	toUserColumn,
	"to_name",
	"note",
	toEmailColumn,
	tokenColumn,
	destinationHouseholdColumn,
	expiresAtColumn,
	statusColumn,
	statusNoteColumn,
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
	%s = sqlc.arg(%s),
	%s = %s
WHERE %s IS NULL
	AND %s = LOWER(sqlc.arg(%s));`,
				householdInvitationsTableName,
				toUserColumn, toUserColumn,
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				toEmailColumn, toEmailColumn,
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
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
);`,
				householdInvitationsTableName, idColumn,
				householdInvitationsTableName,
				householdInvitationsTableName, archivedAtColumn,
				householdInvitationsTableName, idColumn, idColumn,
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
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s > %s
	AND %s.%s = LOWER(sqlc.arg(%s))
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				householdInvitationsTableName,
				householdsTableName, householdInvitationsTableName, destinationHouseholdColumn, householdsTableName, idColumn,
				usersTableName, householdInvitationsTableName, fromUserColumn, usersTableName, idColumn,
				householdInvitationsTableName, archivedAtColumn,
				householdInvitationsTableName, expiresAtColumn, currentTimeExpression,
				householdInvitationsTableName, toEmailColumn, toEmailColumn,
				householdInvitationsTableName, tokenColumn, tokenColumn,
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
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s > %s
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				householdInvitationsTableName,
				householdsTableName, householdInvitationsTableName, destinationHouseholdColumn, householdsTableName, idColumn,
				usersTableName, householdInvitationsTableName, fromUserColumn, usersTableName, idColumn,
				householdInvitationsTableName, archivedAtColumn,
				householdInvitationsTableName, expiresAtColumn, currentTimeExpression,
				householdInvitationsTableName, destinationHouseholdColumn, destinationHouseholdColumn,
				householdInvitationsTableName, idColumn, idColumn,
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
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s > %s
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				householdInvitationsTableName,
				householdsTableName, householdInvitationsTableName, destinationHouseholdColumn, householdsTableName, idColumn,
				usersTableName, householdInvitationsTableName, fromUserColumn, usersTableName, idColumn,
				householdInvitationsTableName, archivedAtColumn,
				householdInvitationsTableName, expiresAtColumn, currentTimeExpression,
				householdInvitationsTableName, tokenColumn, tokenColumn,
				householdInvitationsTableName, idColumn, idColumn,
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
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					householdInvitationsTableName,
					true,
				),
				buildTotalCountSelect(
					householdInvitationsTableName,
				),
				householdInvitationsTableName,

				householdsTableName, householdInvitationsTableName, destinationHouseholdColumn, householdsTableName, idColumn,
				usersTableName, householdInvitationsTableName, fromUserColumn, usersTableName, idColumn,
				householdInvitationsTableName, archivedAtColumn,
				householdInvitationsTableName, fromUserColumn, fromUserColumn,
				householdInvitationsTableName, statusColumn, statusColumn,

				buildFilterConditions(
					householdInvitationsTableName,
					true,
				),
				offsetLimitAddendum,
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
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					householdInvitationsTableName,
					true,
				),
				buildTotalCountSelect(
					householdInvitationsTableName,
				),
				householdInvitationsTableName,
				householdsTableName, householdInvitationsTableName, destinationHouseholdColumn, householdsTableName, idColumn,
				usersTableName, householdInvitationsTableName, fromUserColumn, usersTableName, idColumn,
				householdInvitationsTableName, archivedAtColumn,
				householdInvitationsTableName, toUserColumn, toUserColumn,
				householdInvitationsTableName, statusColumn, statusColumn,
				buildFilterConditions(
					householdInvitationsTableName,
					true,
				),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SetHouseholdInvitationStatus",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = sqlc.arg(%s),
	%s = sqlc.arg(%s),
	%s = %s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				householdInvitationsTableName,
				statusColumn, statusColumn,
				statusNoteColumn, statusNoteColumn,
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn, currentTimeExpression,
				archivedAtColumn,
				idColumn, idColumn,
			)),
		},
	}
}
