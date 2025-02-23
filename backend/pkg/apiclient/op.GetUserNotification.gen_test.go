// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetUserNotification(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/user_notifications/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userNotificationID := fake.BuildFakeID()

		data := &UserNotification{}
		expected := &APIResponse[*UserNotification]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, userNotificationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetUserNotification(ctx, userNotificationID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid userNotification ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetUserNotification(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userNotificationID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetUserNotification(ctx, userNotificationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		userNotificationID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, userNotificationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetUserNotification(ctx, userNotificationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
