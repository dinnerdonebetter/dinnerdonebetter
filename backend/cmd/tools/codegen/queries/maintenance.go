package main

import "fmt"

func buildMaintenanceQueries(database string) []*Query {
	switch database {
	case postgres:
		const oneDayAgo = `(NOW() - interval '1 day')`
		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "DeleteExpiredOAuth2ClientTokens",
					Type: ExecRowsType,
				},
				Content: fmt.Sprintf(`DELETE FROM %s WHERE %s < %s AND %s < %s AND %s < %s;`, oauth2ClientTokensTableName, codeExpiresAtColumn, oneDayAgo, accessExpiresAtColumn, oneDayAgo, refreshExpiresAtColumn, oneDayAgo),
			},
		}
	default:
		return nil
	}
}
