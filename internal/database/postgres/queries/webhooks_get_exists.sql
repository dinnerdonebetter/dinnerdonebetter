SELECT EXISTS ( SELECT webhooks.id FROM webhooks WHERE webhooks.archived_at IS NULL AND webhooks.belongs_to_household = $1 AND webhooks.id = $2 );
