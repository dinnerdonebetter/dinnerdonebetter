package capitalism

import (
	"bytes"
	"net/http"
	"testing"

	capitalismmock "github.com/dinnerdonebetter/backend/internal/lib/capitalism/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/random"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCapitalismService_StripeWebhookHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		secret, err := random.GenerateHexEncodedString(ctx, 32)
		require.NoError(t, err)
		require.NotEmpty(t, secret)

		helper := buildTestHelper(t)

		mpm := &capitalismmock.MockPaymentManager{}
		mpm.On("HandleEventWebhook", mock.AnythingOfType("*http.Request")).Return(nil)
		helper.service.paymentManager = mpm

		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.IncomingWebhookHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
	})
}
