package integration

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/services/webhooks/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createWebhookForTest(t *testing.T, testClient client.Client) *webhooks.Webhook {
	t.Helper()
	ctx := t.Context()

	exampleWebhookInput := &webhooks.WebhookCreationRequestInput{
		ContentType: "application/json",
		Method:      http.MethodPost,
		Name:        t.Name(),
		URL:         "https://whatever.gov",
		Events:      []string{webhooks.WebhookCreatedTriggerEvent},
	}

	input := converters.ConvertWebhookCreationRequestInputToGRPCWebhookCreationRequestInput(exampleWebhookInput)

	createdWebhook, err := testClient.CreateWebhook(ctx, &webhookssvc.CreateWebhookRequest{Input: input})
	require.NoError(t, err)

	retrievedWebhook, err := testClient.GetWebhook(ctx, &webhookssvc.GetWebhookRequest{WebhookID: createdWebhook.Created.ID})
	require.NoError(t, err)
	require.NotNil(t, retrievedWebhook)

	return converters.ConvertGRPCWebhookToWebhook(retrievedWebhook.Result)
}

func TestWebhooks_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		_, testClient := createUserAndClientForTest(t)
		createWebhookForTest(t, testClient)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.CreateWebhook(ctx, &webhookssvc.CreateWebhookRequest{})
		require.Error(t, err)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)

		exampleWebhookInput := &webhooks.WebhookCreationRequestInput{
			ContentType: "application/whatever",
			Method:      "UNRECOGNIZED",
			Name:        t.Name(),
			URL:         "invalid protocol :\\ neato.ai",
			Events:      []string{},
		}

		input := converters.ConvertWebhookCreationRequestInputToGRPCWebhookCreationRequestInput(exampleWebhookInput)

		_, err := testClient.CreateWebhook(ctx, &webhookssvc.CreateWebhookRequest{Input: input})
		assert.Error(t, err)
	})
}

func TestWebhooks_Reading(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		createdWebhook := createWebhookForTest(t, testClient)

		retrieved, err := testClient.GetWebhook(ctx, &webhookssvc.GetWebhookRequest{WebhookID: createdWebhook.ID})
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
	})

	T.Run("nonexistent ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)

		retrieved, err := testClient.GetWebhook(ctx, &webhookssvc.GetWebhookRequest{WebhookID: nonexistentID})
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)
		_, err := c.GetWebhook(ctx, &webhookssvc.GetWebhookRequest{})
		assert.Error(t, err)
	})
}

func TestWebhooks_Listing(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)

		createdWebhooks := []*webhooks.Webhook{}
		for range exampleQuantity {
			createdWebhooks = append(createdWebhooks, createWebhookForTest(t, testClient))
		}

		results, err := testClient.GetWebhooks(ctx, &webhookssvc.GetWebhooksRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdWebhooks))
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)
		_, err := c.GetWebhooks(ctx, &webhookssvc.GetWebhooksRequest{})
		assert.Error(t, err)
	})
}

func TestWebhooks_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		createdWebhook := createWebhookForTest(t, testClient)

		_, err := testClient.ArchiveWebhook(ctx, &webhookssvc.ArchiveWebhookRequest{WebhookID: createdWebhook.ID})
		assert.NoError(t, err)
	})

	T.Run("nonexistentID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		createWebhookForTest(t, testClient)

		_, err := testClient.ArchiveWebhook(ctx, &webhookssvc.ArchiveWebhookRequest{WebhookID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)
		_, err := c.ArchiveWebhook(ctx, &webhookssvc.ArchiveWebhookRequest{})
		assert.Error(t, err)
	})
}

func TestWebhookTriggerEvents_Adding(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		createdWebhook := createWebhookForTest(t, testClient)

		_, err := testClient.AddWebhookTriggerEvent(ctx, &webhookssvc.AddWebhookTriggerEventRequest{
			WebhookID: createdWebhook.ID,
			Input: &webhookssvc.WebhookTriggerEventCreationRequestInput{
				BelongsToWebhook: createdWebhook.ID,
				TriggerEvent:     webhooks.WebhookArchivedTriggerEvent,
			},
		})
		assert.NoError(t, err)
	})

	T.Run("nonexistentID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		createWebhookForTest(t, testClient)

		_, err := testClient.AddWebhookTriggerEvent(ctx, &webhookssvc.AddWebhookTriggerEventRequest{WebhookID: nonexistentID, Input: &webhookssvc.WebhookTriggerEventCreationRequestInput{}})
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)
		_, err := c.AddWebhookTriggerEvent(ctx, &webhookssvc.AddWebhookTriggerEventRequest{})
		assert.Error(t, err)
	})
}

func TestWebhookTriggerEvents_Removing(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		createdWebhook := createWebhookForTest(t, testClient)

		createdTriggerEvent, err := testClient.AddWebhookTriggerEvent(ctx, &webhookssvc.AddWebhookTriggerEventRequest{
			WebhookID: createdWebhook.ID,
			Input: &webhookssvc.WebhookTriggerEventCreationRequestInput{
				BelongsToWebhook: createdWebhook.ID,
				TriggerEvent:     webhooks.WebhookArchivedTriggerEvent,
			},
		})
		require.NoError(t, err)

		_, err = testClient.ArchiveWebhookTriggerEvent(ctx, &webhookssvc.ArchiveWebhookTriggerEventRequest{
			WebhookID:             createdWebhook.ID,
			WebhookTriggerEventID: createdTriggerEvent.Created.ID,
		})
		assert.NoError(t, err)
	})

	T.Run("nonexistentID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		createWebhookForTest(t, testClient)

		_, err := testClient.ArchiveWebhookTriggerEvent(ctx, &webhookssvc.ArchiveWebhookTriggerEventRequest{WebhookID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)
		_, err := c.ArchiveWebhookTriggerEvent(ctx, &webhookssvc.ArchiveWebhookTriggerEventRequest{})
		assert.Error(t, err)
	})
}
