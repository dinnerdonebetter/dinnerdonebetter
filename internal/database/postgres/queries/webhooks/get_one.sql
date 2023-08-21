-- name: GetWebhook :many

SELECT
	webhooks.id as webhook_id,
	webhooks.name as webhook_name,
	webhooks.content_type as webhook_content_type,
	webhooks.url as webhook_url,
	webhooks.method as webhook_method,
	webhook_trigger_events.id as webhook_trigger_event_id,
	webhook_trigger_events.trigger_event as webhook_trigger_event_trigger_event,
	webhook_trigger_events.belongs_to_webhook as webhook_trigger_event_belongs_to_webhook,
	webhook_trigger_events.created_at as webhook_trigger_event_created_at,
	webhook_trigger_events.archived_at as webhook_trigger_event_archived_at,
	webhooks.created_at as webhook_created_at,
	webhooks.last_updated_at as webhook_last_updated_at,
	webhooks.archived_at as webhook_archived_at,
	webhooks.belongs_to_household as webhook_belongs_to_household
FROM webhook_trigger_events
 JOIN webhooks ON webhook_trigger_events.belongs_to_webhook=webhooks.id
WHERE webhooks.archived_at IS NULL
  AND webhook_trigger_events.archived_at IS NULL
  AND webhooks.belongs_to_household = sqlc.arg(household_id)
  AND webhooks.id = sqlc.arg(webhook_id);