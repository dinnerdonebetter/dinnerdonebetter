-- name: WebhookExists :one
SELECT EXISTS ( SELECT webhooks.id FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_household = $1 AND webhooks.id = $2 );

-- name: GetWebhook :one
SELECT webhooks.id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_household FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_household = $1 AND webhooks.id = $2;

-- name: GetAllWebhooksCount :one
SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL;

-- name: CreateWebhook :exec
INSERT INTO webhooks (id,name,content_type,url,method,events,data_types,topics,belongs_to_household) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);

-- name: ArchiveWebhook :exec
UPDATE webhooks SET
	last_updated_on = extract(epoch FROM NOW()),
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND belongs_to_household = $1
AND id = $2;
