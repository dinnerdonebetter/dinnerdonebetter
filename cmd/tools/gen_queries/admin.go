package main

func buildAdminQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "SetUserAccountStatus",
				Type: ExecRowsType,
			},
			Content: formatQuery(
				buildUpdateQuery(
					"users",
					[]string{"user_account_status", "user_account_status_explanation"},
					"",
				),
			),
		},
	}
}
