package postgres

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	// webhooksTableColumns are the columns for the webhooks table.
	webhooksTableColumns = []string{
		"webhooks.id",
		"webhooks.name",
		"webhooks.content_type",
		"webhooks.url",
		"webhooks.method",
		"webhook_trigger_events.id",
		"webhook_trigger_events.trigger_event",
		"webhook_trigger_events.belongs_to_webhook",
		"webhook_trigger_events.created_at",
		"webhook_trigger_events.archived_at",
		"webhooks.created_at",
		"webhooks.last_updated_at",
		"webhooks.archived_at",
		"webhooks.belongs_to_household",
	}
)

func buildMockRowsFromWebhooks(includeCounts bool, filteredCount uint64, webhooks ...*types.Webhook) *sqlmock.Rows {
	columns := webhooksTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range webhooks {
		for _, e := range w.Events {
			rowValues := []driver.Value{
				w.ID,
				w.Name,
				w.ContentType,
				w.URL,
				w.Method,
				e.ID,
				e.TriggerEvent,
				e.BelongsToWebhook,
				e.CreatedAt,
				e.ArchivedAt,
				w.CreatedAt,
				w.LastUpdatedAt,
				w.ArchivedAt,
				w.BelongsToHousehold,
			}

			if includeCounts {
				rowValues = append(rowValues, filteredCount, len(webhooks))
			}

			exampleRows.AddRow(rowValues...)
		}
	}

	return exampleRows
}

func buildErroneousMockRowsFromWebhooks(includeCounts bool, filteredCount uint64, webhooks ...*types.Webhook) *sqlmock.Rows {
	columns := webhooksTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range webhooks {
		for _, e := range w.Events {
			rowValues := []driver.Value{
				w.ArchivedAt,
				e.ID,
				e.TriggerEvent,
				e.BelongsToWebhook,
				e.CreatedAt,
				e.ArchivedAt,
				w.ID,
				w.Name,
				w.ContentType,
				w.URL,
				w.Method,
				w.CreatedAt,
				w.LastUpdatedAt,
				w.BelongsToHousehold,
			}

			if includeCounts {
				rowValues = append(rowValues, filteredCount, len(webhooks))
			}

			exampleRows.AddRow(rowValues...)
		}
	}

	return exampleRows
}

func TestQuerier_WebhookExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.WebhookExists(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleWebhookID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		actual, err := c.WebhookExists(ctx, exampleWebhookID, "")
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetWebhook(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleWebhook := fakes.BuildFakeWebhook()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetWebhook(ctx, exampleWebhook.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		c, _ := buildTestClient(t)

		actual, err := c.GetWebhooks(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateWebhook(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_createWebhookTriggerEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		created, err := c.createWebhookTriggerEvent(ctx, c.db, nil)
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestQuerier_ArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveWebhook(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleWebhookID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveWebhook(ctx, exampleWebhookID, ""))
	})
}
