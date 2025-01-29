// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateWebhookTriggerEvent(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/webhooks/%s/trigger_events"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		webhookID := fakes.BuildFakeID()

		data := fakes.BuildFakeWebhookTriggerEvent()
		expected := &types.APIResponse[*types.WebhookTriggerEvent]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeWebhookTriggerEventCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, webhookID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateWebhookTriggerEvent(ctx, webhookID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeWebhookTriggerEventCreationRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateWebhookTriggerEvent(ctx, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		webhookID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeWebhookTriggerEventCreationRequestInput()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateWebhookTriggerEvent(ctx, webhookID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		webhookID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeWebhookTriggerEventCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, webhookID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateWebhookTriggerEvent(ctx, webhookID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
