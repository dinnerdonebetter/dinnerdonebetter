package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createWebhookForTest(t *testing.T, ctx context.Context, exampleWebhook *types.Webhook, dbc *Querier) *types.Webhook {
	t.Helper()

	// create
	if exampleWebhook == nil {
		exampleWebhook = fakes.BuildFakeWebhook()
	}
	dbInput := converters.ConvertWebhookToWebhookDatabaseCreationInput(exampleWebhook)

	created, err := dbc.CreateWebhook(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleWebhook.CreatedAt = created.CreatedAt
	for i := range created.Events {
		exampleWebhook.Events[i].CreatedAt = created.Events[i].CreatedAt
	}
	assert.Equal(t, exampleWebhook, created)

	webhook, err := dbc.GetWebhook(ctx, created.ID, created.BelongsToHousehold)
	exampleWebhook.CreatedAt = webhook.CreatedAt
	for i := range created.Events {
		exampleWebhook.Events[i].CreatedAt = webhook.Events[i].CreatedAt
	}

	assert.NoError(t, err)
	assert.Equal(t, webhook, exampleWebhook)

	return created
}

func TestQuerier_Integration_Webhooks(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := t.Context()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, householdID)

	exampleWebhook := fakes.BuildFakeWebhook()
	exampleWebhook.BelongsToHousehold = householdID
	createdWebhooks := []*types.Webhook{}

	// create
	createdWebhooks = append(createdWebhooks, createWebhookForTest(t, ctx, exampleWebhook, dbc))

	// create more
	for i := range exampleQuantity {
		input := fakes.BuildFakeWebhook()
		input.Name = fmt.Sprintf("%s %d", exampleWebhook.Name, i)
		input.BelongsToHousehold = householdID
		createdWebhooks = append(createdWebhooks, createWebhookForTest(t, ctx, input, dbc))
	}

	// fetch as list
	webhooks, err := dbc.GetWebhooks(ctx, householdID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, webhooks.Data)
	assert.Equal(t, len(createdWebhooks), len(webhooks.Data))

	// fetch as list
	webhooksByHouseholdAndEvent, err := dbc.GetWebhooksForHouseholdAndEvent(ctx, householdID, createdWebhooks[0].Events[0].TriggerEvent)
	assert.NoError(t, err)
	assert.NotEmpty(t, webhooksByHouseholdAndEvent)

	createdEvent, err := dbc.AddWebhookTriggerEvent(ctx, householdID, &types.WebhookTriggerEventDatabaseCreationInput{
		ID:               identifiers.New(),
		BelongsToWebhook: createdWebhooks[0].ID,
		TriggerEvent:     types.WebhookArchivedServiceEventType,
	})
	assert.NoError(t, err)
	assert.NotNil(t, createdEvent)

	createdWebhooks[0].Events = append(createdWebhooks[0].Events, createdEvent)

	// delete
	for _, webhook := range createdWebhooks {
		for _, event := range webhook.Events {
			assert.NoError(t, dbc.ArchiveWebhookTriggerEvent(ctx, webhook.ID, event.ID))
		}

		assert.NoError(t, dbc.ArchiveWebhook(ctx, webhook.ID, householdID))

		var exists bool
		exists, err = dbc.WebhookExists(ctx, webhook.ID, householdID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.Webhook
		y, err = dbc.GetWebhook(ctx, webhook.ID, householdID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_WebhookExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		exampleHouseholdID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.WebhookExists(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

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

		ctx := t.Context()
		c, _ := buildTestClient(t)

		actual, err := c.GetWebhook(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleWebhook := fakes.BuildFakeWebhook()

		ctx := t.Context()
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

		ctx := t.Context()
		filter := filtering.DefaultQueryFilter()
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

		ctx := t.Context()
		c, _ := buildTestClient(t)

		actual, err := c.CreateWebhook(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with msising user ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		input := converters.ConvertWebhookToWebhookDatabaseCreationInput(fakes.BuildFakeWebhook())

		actual, err := c.CreateWebhook(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_createWebhookTriggerEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		created, err := c.createWebhookTriggerEvent(ctx, c.db, fakes.BuildFakeID(), nil)
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("with missing household ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildTestClient(t)
		fakes.BuildFakeWebhookTriggerEvent()

		input := converters.ConvertWebhookTriggerEventToWebhookTriggerEventDatabaseCreationInput(fakes.BuildFakeWebhookTriggerEvent())

		created, err := c.createWebhookTriggerEvent(ctx, c.db, "", input)
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestQuerier_ArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveWebhook(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleWebhookID := fakes.BuildFakeID()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveWebhook(ctx, exampleWebhookID, ""))
	})
}

func TestQuerier_ArchiveWebhookTriggerEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveWebhookTriggerEvent(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid webhook trigger event ID", func(t *testing.T) {
		t.Parallel()

		exampleWebhookID := fakes.BuildFakeID()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveWebhookTriggerEvent(ctx, exampleWebhookID, ""))
	})
}
