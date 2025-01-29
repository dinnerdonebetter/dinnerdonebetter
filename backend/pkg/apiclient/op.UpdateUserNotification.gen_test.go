// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateUserNotification(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/user_notifications/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userNotificationID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*UserNotification](t)

		expected := &APIResponse[*UserNotification]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*UserNotificationUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, userNotificationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateUserNotification(ctx, userNotificationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid userNotification ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*UserNotificationUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateUserNotification(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userNotificationID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*UserNotificationUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateUserNotification(ctx, userNotificationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userNotificationID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*UserNotificationUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, userNotificationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateUserNotification(ctx, userNotificationID, exampleInput)

		assert.Error(t, err)
	})
}
