package main

const householdsTableName = "households"

var householdsColumns = []string{
	"id",
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
