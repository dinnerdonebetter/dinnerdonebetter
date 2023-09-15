package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const householdsTableName = "households"

var householdsColumns = []string{
	idColumn,
	"name",
	"billing_status",
	"contact_phone",
	"payment_processor_customer_id",
	"subscription_plan_id",
	"belongs_to_user",
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
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveHousehold",
				Type: ExecRowsType,
			},
			Content: "",
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
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdByIDWithMemberships",
				Type: ManyType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdsForUser",
				Type: ManyType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateHousehold",
				Type: ExecRowsType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateHouseholdWebhookEncryptionKey",
				Type: ExecRowsType,
			},
			Content: "",
		},
	}
}
