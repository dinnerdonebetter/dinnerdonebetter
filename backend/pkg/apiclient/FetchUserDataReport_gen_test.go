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

func TestClient_FetchUserDataReport(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/data_privacy/reports/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userDataAggregationReportID := fakes.BuildFakeID()

		data := fakes.BuildFakeUserDataCollection()
		data.User.TwoFactorSecret = ""
		data.User.HashedPassword = ""
		data.User.TwoFactorSecretVerifiedAt = nil
		for i := range data.Households {
			data.Households[i].WebhookEncryptionKey = ""
		}
		for i := range data.SentInvites {
			data.SentInvites[i].DestinationHousehold.WebhookEncryptionKey = ""
			data.SentInvites[i].FromUser.TwoFactorSecret = ""
			data.SentInvites[i].FromUser.HashedPassword = ""
			data.SentInvites[i].FromUser.TwoFactorSecretVerifiedAt = nil
		}
		for i := range data.ReceivedInvites {
			data.ReceivedInvites[i].DestinationHousehold.WebhookEncryptionKey = ""
			data.ReceivedInvites[i].FromUser.TwoFactorSecret = ""
			data.ReceivedInvites[i].FromUser.HashedPassword = ""
			data.ReceivedInvites[i].FromUser.TwoFactorSecretVerifiedAt = nil
		}

		expected := &types.APIResponse[*types.UserDataCollection]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, userDataAggregationReportID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.FetchUserDataReport(ctx, userDataAggregationReportID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid userDataAggregationReport ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.FetchUserDataReport(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userDataAggregationReportID := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.FetchUserDataReport(ctx, userDataAggregationReportID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userDataAggregationReportID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, userDataAggregationReportID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.FetchUserDataReport(ctx, userDataAggregationReportID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
