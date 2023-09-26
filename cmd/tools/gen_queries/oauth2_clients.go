package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	oauth2ClientsTableName = "oauth2_clients"
)

var oauth2ClientsColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	"client_id",
	"client_secret",
	createdAtColumn,
	archivedAtColumn,
}

func buildOAuth2ClientsQueries() []*Query {
	insertColumns := filterForInsert(oauth2ClientsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveOAuth2Client",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				oauth2ClientsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateOAuth2Client",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				oauth2ClientsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetOAuth2ClientByClientID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
	AND oauth2_clients.client_id = $1;`,
				strings.Join(applyToEach(oauth2ClientsColumns, func(_ int, s string) string {
					return fmt.Sprintf("%s.%s", oauth2ClientsTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetOAuth2ClientByDatabaseID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
	AND oauth2_clients.id = $1;`,
				strings.Join(applyToEach(oauth2ClientsColumns, func(_ int, s string) string {
					return fmt.Sprintf("%s.%s", oauth2ClientsTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetOAuth2Clients",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	(
		SELECT
			COUNT(oauth2_clients.id)
		FROM
			oauth2_clients
		WHERE
			oauth2_clients.archived_at IS NULL
			%s
	) as filtered_count,
	%s
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
	%s
	OFFSET sqlc.narg(query_offset)
	LIMIT sqlc.narg(query_limit);
`,
				strings.Join(applyToEach(oauth2ClientsColumns, func(_ int, s string) string {
					return fmt.Sprintf("%s.%s", oauth2ClientsTableName, s)
				}), ",\n\t"),
				strings.Join(strings.Split(
					buildFilterConditions(
						oauth2ClientsTableName,
						false,
					), "\n"),
					"\n\t\t",
				),
				buildTotalCountSelect(
					usersTableName,
				),
				buildFilterConditions(
					oauth2ClientsTableName,
					false,
				),
			)),
		},
	}
}
