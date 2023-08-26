-- name: CreateWebhookTriggerEvent :exec

INSERT INTO webhook_trigger_events (
	id,
	trigger_event,
	belongs_to_webhook
) VALUES ($1, $2, $3);