// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_UpdateUserNotification(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/user_notifications/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userNotificationID := fakes.BuildFakeID()

		data := fakes.BuildFakeUserNotification()
		expected := &types.APIResponse[*types.UserNotification]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeUserNotificationUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, userNotificationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateUserNotification(ctx, userNotificationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid userNotification ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeUserNotificationUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateUserNotification(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userNotificationID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeUserNotificationUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateUserNotification(ctx, userNotificationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userNotificationID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeUserNotificationUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, userNotificationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateUserNotification(ctx, userNotificationID, exampleInput)

		assert.Error(t, err)
	})
}
