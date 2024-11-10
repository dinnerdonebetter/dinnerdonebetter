package compression

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_compressor_CompressBytes(T *testing.T) {
	T.Parallel()

	T.Run("zstandard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		comp := &compressor{algo: algoZstd}

		dt, err := time.Parse(time.DateTime, time.DateTime)
		require.NoError(t, err)

		x := &types.Webhook{
			CreatedAt:          dt,
			Name:               "testing",
			URL:                "https://whatever.gov",
			Method:             http.MethodPost,
			ID:                 "blah-blah-blah",
			BelongsToHousehold: "something",
			ContentType:        "application/json",
			Events:             []*types.WebhookTriggerEvent{},
		}

		encoder := encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		expected := "KLUv_QQAjQUAsownIUAL3AC7zVqQSRJj2JxtUcpeGL5HRbWY0WOmmqqZh3ixAgARDAv416uPHn35NRzLrFNhzbk8Xt1D0qu9rKxWht_xqvHS53ABJ29Zbro4urGyLY5X7OgJgsHB8YpXK3P0h5eARh_tXA7Hgc91r-Nb4xx9WdZP-3hlcpVL6JftqDVnH7aUJo-uoEgokFCBJIlkIMkYSClxfBMo7dRGOUQFAJ5hOIJSDOewuCJ_Z1BlfrNA6Q=="
		compressed, err := comp.CompressBytes(encoder.MustEncodeJSON(ctx, x))
		assert.NoError(t, err)
		actual := base64.URLEncoding.EncodeToString(compressed)

		assert.Equal(t, expected, actual)
	})

	T.Run("s2", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		comp := &compressor{algo: algoS2}

		dt, err := time.Parse(time.DateTime, time.DateTime)
		require.NoError(t, err)

		x := &types.Webhook{
			CreatedAt:          dt,
			Name:               "testing",
			URL:                "https://whatever.gov",
			Method:             http.MethodPost,
			ID:                 "blah-blah-blah",
			BelongsToHousehold: "something",
			ContentType:        "application/json",
			Events:             []*types.WebhookTriggerEvent{},
		}

		encoder := encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		expected := "_wYAAFMyc1R3TwHyAACmH-pBeyJjcmVhdGVkQXQiOiIyMDA2LTAxLTAyVDE1OjA0OjA1WiIsImFyY2hpdmVkQXQiOm51bGwsImxhc3RVcGRhdGVkQXQiOm51bGwsIm5hbWUiOiJ0ZXN0aW5nIiwidXJsIjoiaHR0cHM6Ly93aGF0ZXZlci5nb3YiLCJtZXRob2QiOiJQT1NUIiwiaWQiOiJibGFoLWJsYWgtYmxhaCIsImJlbG9uZ3NUb0hvdXNlaG9sZCI6InNvbWV0aGluZyIsImNvbnRlbnRUeXBlIjoiYXBwbGljYXRpb24vanNvbiIsImV2ZW50cyI6W119Cg=="
		compressed, err := comp.CompressBytes(encoder.MustEncodeJSON(ctx, x))
		assert.NoError(t, err)
		actual := base64.URLEncoding.EncodeToString(compressed)

		assert.Equal(t, expected, actual)
	})
}

func Test_compressor_DecompressBytes(T *testing.T) {
	T.Parallel()

	algorithms := []algo{
		algoZstd,
		algoS2,
	}

	for _, a := range algorithms {
		T.Run("zstandard", func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			comp := &compressor{algo: a}

			dt, err := time.Parse(time.DateTime, time.DateTime)
			require.NoError(t, err)

			x := &types.Webhook{
				CreatedAt:          dt,
				Name:               "testing",
				URL:                "https://whatever.gov",
				Method:             http.MethodPost,
				ID:                 "blah-blah-blah",
				BelongsToHousehold: "something",
				ContentType:        "application/json",
				Events:             []*types.WebhookTriggerEvent{},
			}

			encoder := encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

			compressed, err := comp.CompressBytes(encoder.MustEncodeJSON(ctx, x))
			assert.NoError(t, err)

			decompressed, err := comp.DecompressBytes(compressed)
			assert.NoError(t, err)

			var y *types.Webhook
			require.NoError(t, encoder.DecodeBytes(ctx, decompressed, &y))

			assert.Equal(t, x, y)
		})
	}
}
