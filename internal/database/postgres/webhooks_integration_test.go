package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
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

	ctx := context.Background()
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
	for i := 0; i < exampleQuantity; i++ {
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
	webhooksByHouseholdAndEvent, err := dbc.GetWebhooksForHouseholdAndEvent(ctx, householdID, types.CustomerEventType(createdWebhooks[0].Events[0].TriggerEvent))
	assert.NoError(t, err)
	assert.NotEmpty(t, webhooksByHouseholdAndEvent)

	// delete
	for _, webhook := range createdWebhooks {
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
