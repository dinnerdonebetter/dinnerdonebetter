// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: webhook_trigger_events.sql

package generated

import (
	"context"
)

const archiveWebhookTriggerEvent = `-- name: ArchiveWebhookTriggerEvent :execrows
UPDATE webhook_trigger_events SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $1
	AND belongs_to_webhook = $2
`

type ArchiveWebhookTriggerEventParams struct {
	ID               string
	BelongsToWebhook string
}

func (q *Queries) ArchiveWebhookTriggerEvent(ctx context.Context, db DBTX, arg *ArchiveWebhookTriggerEventParams) (int64, error) {
	result, err := db.ExecContext(ctx, archiveWebhookTriggerEvent, arg.ID, arg.BelongsToWebhook)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const createWebhookTriggerEvent = `-- name: CreateWebhookTriggerEvent :exec
INSERT INTO webhook_trigger_events (
	id,
	trigger_event,
	belongs_to_webhook
) VALUES (
	$1,
	$2,
	$3
)
`

type CreateWebhookTriggerEventParams struct {
	ID               string
	TriggerEvent     WebhookEvent
	BelongsToWebhook string
}

func (q *Queries) CreateWebhookTriggerEvent(ctx context.Context, db DBTX, arg *CreateWebhookTriggerEventParams) error {
	_, err := db.ExecContext(ctx, createWebhookTriggerEvent, arg.ID, arg.TriggerEvent, arg.BelongsToWebhook)
	return err
}
