-- name: GetWebhookTriggerEventsForWebhook :many
SELECT
    webhook_trigger_events.id,
    webhook_trigger_events.trigger_event,
    webhook_trigger_events.belongs_to_webhook,
    webhook_trigger_events.created_at,
    webhook_trigger_events.archived_at
FROM webhook_trigger_events
    JOIN webhooks ON webhook_trigger_events.belongs_to_webhook=webhooks.id
WHERE webhook_trigger_events.archived_at IS NULL
    AND webhooks.archived_at IS NULL
    AND webhooks.id = $1
    AND webhook_trigger_events.belongs_to_webhook = $1;
