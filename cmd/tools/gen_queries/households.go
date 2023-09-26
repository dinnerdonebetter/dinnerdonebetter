package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	householdsTableName = "households"
)

var householdsColumns = []string{
	idColumn,
	nameColumn,
	"billing_status",
	"contact_phone",
	"payment_processor_customer_id",
	"subscription_plan_id",
	belongsToUserColumn,
	"time_zone",
	"address_line_1",
	"address_line_2",
	"city",
	"state",
	"zip_code",
	"country",
	"latitude",
	"longitude",
	"last_payment_provider_sync_occurred_at",
	"webhook_hmac_secret",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildHouseholdsQueries() []*Query {
	insertColumns := filterForInsert(householdsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "AddToHouseholdDuringCreation",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO household_user_memberships (
	%s
) VALUES (
	%s
);`,
				strings.Join(filterForInsert(householdUserMembershipsColumns, "default_household"), ",\n\t"),
				strings.Join(applyToEach(filterForInsert(householdUserMembershipsColumns, "default_household"), func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveHousehold",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdsTableName,
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToUserColumn,
				belongsToUserColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateHousehold",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				householdsTableName,
				strings.Join(filterForInsert(
					insertColumns,
					"time_zone",
					"payment_processor_customer_id",
					"last_payment_provider_sync_occurred_at",
					"subscription_plan_id",
				), ",\n\t"),
				strings.Join(applyToEach(filterForInsert(
					insertColumns,
					"time_zone",
					"payment_processor_customer_id",
					"last_payment_provider_sync_occurred_at",
					"subscription_plan_id",
				), func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdByIDWithMemberships",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
	JOIN users ON household_user_memberships.belongs_to_user = users.id
WHERE households.archived_at IS NULL
	AND household_user_memberships.archived_at IS NULL
	AND households.id = sqlc.arg(id);`,
				strings.Join(append(
					append(
						applyToEach(householdsColumns, func(_ int, s string) string {
							return fmt.Sprintf("%s.%s", householdsTableName, s)
						}),
						applyToEach(usersColumns, func(_ int, s string) string {
							return fmt.Sprintf("%s.%s as user_%s", usersTableName, s, s)
						})...,
					),
					applyToEach(householdUserMembershipsColumns, func(_ int, s string) string {
						return fmt.Sprintf("%s.%s as membership_%s", householdUserMembershipsTableName, s, s)
					})...,
				), ",\n\t")),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdsForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	(
		SELECT
			COUNT(households.id)
		FROM
			households
		    JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
        WHERE households.archived_at IS NULL
            AND household_user_memberships.belongs_to_user = sqlc.arg(user_id)%s
	) as filtered_count,
    %s
FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
    JOIN users ON household_user_memberships.belongs_to_user = users.id
WHERE households.archived_at IS NULL
    AND household_user_memberships.archived_at IS NULL
    AND household_user_memberships.belongs_to_user = sqlc.arg(user_id)
	%s
	LIMIT sqlc.narg(query_limit)
	OFFSET sqlc.narg(query_offset);`,
				strings.Join(applyToEach(householdsColumns, func(_ int, s string) string {
					return fmt.Sprintf("%s.%s", householdsTableName, s)
				}), ",\n\t"),
				strings.Join(applyToEach(strings.Split(buildFilterConditions(
					householdsTableName,
					true,
				), "\n"), func(i int, s string) string {
					if i == 0 {
						return fmt.Sprintf("\n\t\t\t%s", s)
					}
					return fmt.Sprintf("\n\t\t%s", s)
				}), ""),
				buildTotalCountSelect(
					householdsTableName,
				),
				buildFilterConditions(
					householdsTableName,
					true,
				),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateHousehold",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdsTableName,
				strings.Join(
					applyToEach(
						filterForUpdate(
							householdsColumns,
							"billing_status",
							"payment_processor_customer_id",
							"subscription_plan_id",
							belongsToUserColumn,
							"time_zone",
							"last_payment_provider_sync_occurred_at",
							"webhook_hmac_secret",
						),
						func(_ int, s string) string {
							return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
						},
					),
					",\n\t",
				),
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToUserColumn,
				belongsToUserColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateHouseholdWebhookEncryptionKey",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s
SET
    webhook_hmac_secret = sqlc.arg(webhook_hmac_secret),
    %s = %s
WHERE %s IS NULL
    AND %s = sqlc.arg(%s)
    AND %s = sqlc.arg(%s);`,
				householdsTableName,
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToUserColumn,
				belongsToUserColumn,
				idColumn,
				idColumn,
			)),
		},
	}
}
