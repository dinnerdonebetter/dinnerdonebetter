package main

import (
	"fmt"
	"strings"

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
	insertColumns := filterForInsert(oauth2ClientsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveOAuth2Client",
				Type: ":execrows",
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;`,
				oauth2ClientsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateOAuth2Client",
				Type: ":exec",
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
				Type: ":one",
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
				Type: ":one",
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
				Type: ":many",
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
          AND oauth2_clients.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND oauth2_clients.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    ) as filtered_count,
    (
        SELECT
            COUNT(oauth2_clients.id)
        FROM
            oauth2_clients
        WHERE
            oauth2_clients.archived_at IS NULL
    ) as total_count
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
    AND oauth2_clients.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND oauth2_clients.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);
`,
				strings.Join(applyToEach(oauth2ClientsColumns, func(_ int, s string) string {
					return fmt.Sprintf("%s.%s", oauth2ClientsTableName, s)
				}), ",\n\t"),
			)),
		},
	}
}
