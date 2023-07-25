package v2

import (
	"context"
	"database/sql"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/doug-martin/goqu/v9"
	"github.com/jinzhu/copier"
)

const (
	webhooksTableName = "webhooks"
)

type (
	// Webhook represents a webhook listener, an endpoint to send an HTTP request to upon an event.
	Webhook struct {
		_ struct{}

		CreatedAt          time.Time              `db:"created_at"             goqu:"skipinsert"`
		LastUpdatedAt      *time.Time             `db:"last_updated_at"        goqu:"skipinsert"`
		ArchivedAt         *time.Time             `db:"archived_at"            goqu:"skipinsert"`
		Name               string                 `db:"name"`
		URL                string                 `db:"url"`
		Method             string                 `db:"method"`
		ID                 string                 `db:"id"                     goqu:"skipupdate"`
		BelongsToHousehold string                 `db:"belongs_to_household"`
		ContentType        string                 `db:"content_type"`
		Events             []*WebhookTriggerEvent `db:"webhook_trigger_events"`
	}

	// WebhookTriggerEvent represents a webhook trigger event.
	WebhookTriggerEvent struct {
		_ struct{}

		CreatedAt        time.Time  `db:"created_at"         goqu:"skipinsert"`
		ArchivedAt       *time.Time `db:"archived_at"        goqu:"skipinsert"`
		ID               string     `db:"id"                 goqu:"skipupdate"`
		BelongsToWebhook string     `db:"belongs_to_webhook"`
		TriggerEvent     string     `db:"trigger_event"`
	}
)

// CreateWebhook gets a webhook from the database.
func (c *DatabaseClient) CreateWebhook(ctx context.Context, input *Webhook) (*types.Webhook, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Insert(webhooksTableName).Rows(
		input,
	)

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "creating webhook")
	}

	var output types.Webhook
	if err := copier.Copy(&output, input); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetWebhook gets a webhook from the database.
func (c *DatabaseClient) GetWebhook(ctx context.Context, webhookID string) (*types.Webhook, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	x := &Webhook{}
	q := c.xdb.From(webhooksTableName).Where(goqu.Ex{
		idColumn:         webhookID,
		archivedAtColumn: nil,
	})

	found, err := q.ScanStructContext(ctx, x)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, sql.ErrNoRows
	}

	var output types.Webhook
	if err = copier.Copy(&output, x); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetWebhooks gets a webhook from the database.
func (c *DatabaseClient) GetWebhooks(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Webhook], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	q := c.xdb.From(webhooksTableName).Where(goqu.Ex{
		archivedAtColumn: nil,
	})
	q = queryFilterToGoqu(q, filter)

	var x []Webhook
	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	output := &types.QueryFilteredResult[types.Webhook]{
		Data:       []*types.Webhook{},
		Pagination: filter.ToPagination(),
	}
	for _, y := range x {
		var z types.Webhook
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output.Data = append(output.Data, &z)
	}

	return output, nil
}

// UpdateWebhook gets a webhook from the database.
func (c *DatabaseClient) UpdateWebhook(ctx context.Context, input *types.Webhook) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	var updateInput Webhook
	if err := copier.Copy(&updateInput, input); err != nil {
		return observability.PrepareError(err, span, "copying input to output")
	}

	q := c.xdb.Update(webhooksTableName).Set(
		updateInput,
	).Set(goqu.Ex{lastUpdatedAtColumn: goqu.L("NOW()")})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "updating webhook")
	}

	return nil
}

// ArchiveWebhook gets a webhook from the database.
func (c *DatabaseClient) ArchiveWebhook(ctx context.Context, webhookID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Update(webhooksTableName).
		Set(goqu.Record{archivedAtColumn: goqu.L("NOW()")}).
		Where(goqu.Ex{idColumn: webhookID})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "archiving webhook")
	}

	return nil
}
