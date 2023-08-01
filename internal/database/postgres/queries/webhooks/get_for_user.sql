SELECT webhooks.id,
       webhooks.name,
       webhooks.content_type,
       webhooks.url,
       webhooks.method,
       webhook_trigger_events.id,
       webhook_trigger_events.trigger_event,
       webhook_trigger_events.belongs_to_webhook,
       webhook_trigger_events.created_at,
       webhook_trigger_events.archived_at,
       webhooks.created_at,
       webhooks.last_updated_at,
       webhooks.archived_at,
       webhooks.belongs_to_household
  FROM webhooks
       JOIN webhook_trigger_events ON
			webhook_trigger_events.belongs_to_webhook
			= webhooks.id
 WHERE webhook_trigger_events.archived_at IS NULL
       AND webhooks.archived_at IS NULL
       AND webhooks.belongs_to_household = $1
       AND webhooks.id = $2;
