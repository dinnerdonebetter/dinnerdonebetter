package encoding

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type example struct {
	Name string `json:"name" xml:"name"`
}

type broken struct {
	Name json.Number `json:"name" xml:"name"`
}

func init() {
	gob.Register(&example{})
	gob.Register(&broken{})
}

func TestServerEncoderDecoder_encodeResponse(T *testing.T) {
	T.Parallel()

	T.Run("defaults to JSON", func(t *testing.T) {
		t.Parallel()
		expectation := "name"
		ex := &example{Name: expectation}
		encoderDecoder, ok := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON).(*serverEncoderDecoder)
		require.True(t, ok)

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder.encodeResponse(ctx, res, ex, http.StatusOK)
		assert.Equal(t, res.Body.String(), fmt.Sprintf("{%q:%q}\n", "name", ex.Name))
	})

	T.Run("as XML", func(t *testing.T) {
		t.Parallel()
		expectation := "name"
		ex := &example{Name: expectation}
		encoderDecoder, ok := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON).(*serverEncoderDecoder)
		require.True(t, ok)

		ctx := context.Background()
		res := httptest.NewRecorder()
		res.Header().Set(ContentTypeHeaderKey, contentTypeXML)

		encoderDecoder.encodeResponse(ctx, res, ex, http.StatusOK)
		assert.Equal(t, fmt.Sprintf("<example><name>%s</name></example>", expectation), res.Body.String())
	})

	T.Run("as Emoji", func(t *testing.T) {
		t.Parallel()
		expectation := "name"
		ex := &example{Name: expectation}
		encoderDecoder, ok := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeEmoji).(*serverEncoderDecoder)
		require.True(t, ok)

		ctx := context.Background()
		res := httptest.NewRecorder()
		res.Header().Set(ContentTypeHeaderKey, contentTypeEmoji)

		encoderDecoder.encodeResponse(ctx, res, ex, http.StatusOK)
		assert.NotEmpty(t, res.Body.String())
	})

	T.Run("with broken structure", func(t *testing.T) {
		t.Parallel()
		expectation := "name"
		ex := &broken{Name: json.Number(expectation)}
		encoderDecoder, ok := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON).(*serverEncoderDecoder)
		require.True(t, ok)

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder.encodeResponse(ctx, res, ex, http.StatusOK)
		assert.Empty(t, res.Body.String())
	})
}

func TestServerEncoderDecoder_EncodeErrorResponse(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		exampleMessage := "something went awry"
		exampleCode := http.StatusBadRequest

		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder.EncodeErrorResponse(ctx, res, exampleMessage, exampleCode)
		assert.Equal(t, "\"something went awry\"\n", res.Body.String())
		assert.Equal(t, exampleCode, res.Code, "expected status code to match")
	})

	T.Run("as XML", func(t *testing.T) {
		t.Parallel()
		exampleMessage := "something went awry"
		exampleCode := http.StatusBadRequest

		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		ctx := context.Background()
		res := httptest.NewRecorder()
		res.Header().Set(ContentTypeHeaderKey, contentTypeXML)

		encoderDecoder.EncodeErrorResponse(ctx, res, exampleMessage, exampleCode)
		assert.Equal(t, "<string>something went awry</string>", res.Body.String())
		assert.Equal(t, exampleCode, res.Code, "expected status code to match")
	})

	T.Run("as Emoji", func(t *testing.T) {
		t.Parallel()
		exampleMessage := "something went awry"
		exampleCode := http.StatusBadRequest

		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		ctx := context.Background()
		res := httptest.NewRecorder()
		res.Header().Set(ContentTypeHeaderKey, contentTypeEmoji)

		encoderDecoder.EncodeErrorResponse(ctx, res, exampleMessage, exampleCode)
		assert.NotEmpty(t, res.Body.String())
		assert.Equal(t, exampleCode, res.Code, "expected status code to match")
	})
}

func TestServerEncoderDecoder_EncodeInvalidInputResponse(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)
		encoderDecoder.EncodeInvalidInputResponse(ctx, res)

		expectedCode := http.StatusBadRequest
		assert.EqualValues(t, expectedCode, res.Code, "expected code to be %d, got %d instead", expectedCode, res.Code)
	})
}

func TestServerEncoderDecoder_EncodeNotFoundResponse(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)
		encoderDecoder.EncodeNotFoundResponse(ctx, res)

		expectedCode := http.StatusNotFound
		assert.EqualValues(t, expectedCode, res.Code, "expected code to be %d, got %d instead", expectedCode, res.Code)
	})
}

func TestServerEncoderDecoder_EncodeUnspecifiedInternalServerErrorResponse(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)
		encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)

		expectedCode := http.StatusInternalServerError
		assert.EqualValues(t, expectedCode, res.Code, "expected code to be %d, got %d instead", expectedCode, res.Code)
	})
}

func TestServerEncoderDecoder_EncodeUnauthorizedResponse(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)
		encoderDecoder.EncodeUnauthorizedResponse(ctx, res)

		expectedCode := http.StatusUnauthorized
		assert.EqualValues(t, expectedCode, res.Code, "expected code to be %d, got %d instead", expectedCode, res.Code)
	})
}

func TestServerEncoderDecoder_EncodeInvalidPermissionsResponse(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)
		encoderDecoder.EncodeInvalidPermissionsResponse(ctx, res)

		expectedCode := http.StatusForbidden
		assert.EqualValues(t, expectedCode, res.Code, "expected code to be %d, got %d instead", expectedCode, res.Code)
	})
}

func TestServerEncoderDecoder_MustEncodeJSON(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		expected := `{"name":"TestServerEncoderDecoder_MustEncodeJSON/standard"}
`
		actual := string(encoderDecoder.MustEncodeJSON(ctx, &example{Name: t.Name()}))

		assert.Equal(t, expected, actual)
	})

	T.Run("with panic", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		defer func() {
			assert.NotNil(t, recover())
		}()

		encoderDecoder.MustEncodeJSON(ctx, &broken{Name: json.Number(t.Name())})
	})
}

func TestServerEncoderDecoder_MustEncode(T *testing.T) {
	T.Parallel()

	T.Run("with JSON", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		expected := `{"name":"TestServerEncoderDecoder_MustEncode/with_JSON"}
`
		actual := string(encoderDecoder.MustEncode(ctx, &example{Name: t.Name()}))

		assert.Equal(t, expected, actual)
	})

	T.Run("with XML", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeXML)

		expected := `<example><name>TestServerEncoderDecoder_MustEncode/with_XML</name></example>`
		actual := string(encoderDecoder.MustEncode(ctx, &example{Name: t.Name()}))

		assert.Equal(t, expected, actual)
	})

	T.Run("with Emoji", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeEmoji)

		actual := string(encoderDecoder.MustEncode(ctx, &example{Name: t.Name()}))
		assert.NotEmpty(t, actual)
	})

	T.Run("with broken struct", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		encoderDecoder, ok := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON).(*serverEncoderDecoder)
		require.True(t, ok)

		defer func() {
			assert.NotNil(t, recover())
		}()

		encoderDecoder.MustEncode(ctx, &broken{Name: json.Number(t.Name())})
	})
}

func TestServerEncoderDecoder_RespondWithData(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		expectation := "name"
		ex := &example{Name: expectation}
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder.RespondWithData(ctx, res, ex)
		assert.Equal(t, res.Body.String(), fmt.Sprintf("{%q:%q}\n", "name", ex.Name))
	})

	T.Run("as XML", func(t *testing.T) {
		t.Parallel()
		expectation := "name"
		ex := &example{Name: expectation}
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		ctx := context.Background()
		res := httptest.NewRecorder()
		res.Header().Set(ContentTypeHeaderKey, contentTypeXML)

		encoderDecoder.RespondWithData(ctx, res, ex)
		assert.Equal(t, fmt.Sprintf("<example><name>%s</name></example>", expectation), res.Body.String())
	})

	T.Run("as Emoji", func(t *testing.T) {
		t.Parallel()
		expectation := "name"
		ex := &example{Name: expectation}
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		ctx := context.Background()
		res := httptest.NewRecorder()
		res.Header().Set(ContentTypeHeaderKey, contentTypeEmoji)

		encoderDecoder.RespondWithData(ctx, res, ex)
		assert.NotEmpty(t, res.Body.String())
	})
}

func TestServerEncoderDecoder_EncodeResponseWithStatus(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		expectation := "name"
		ex := &example{Name: expectation}
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		ctx := context.Background()
		res := httptest.NewRecorder()

		expected := 666
		encoderDecoder.EncodeResponseWithStatus(ctx, res, ex, expected)

		assert.Equal(t, expected, res.Code, "expected code to be %d, but got %d", expected, res.Code)
		assert.Equal(t, res.Body.String(), fmt.Sprintf("{%q:%q}\n", "name", ex.Name))
	})
}

func TestServerEncoderDecoder_DecodeRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		expectation := "name"
		e := &example{Name: expectation}
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		bs, err := json.Marshal(e)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"https://whatever.whocares.gov",
			bytes.NewReader(bs),
		)
		require.NoError(t, err)
		req.Header.Set(ContentTypeHeaderKey, contentTypeJSON)

		var x example
		assert.NoError(t, encoderDecoder.DecodeRequest(ctx, req, &x))
		assert.Equal(t, x.Name, e.Name)
	})

	T.Run("as XML", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		expectation := "name"
		e := &example{Name: expectation}
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		bs, err := xml.Marshal(e)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"https://whatever.whocares.gov",
			bytes.NewReader(bs),
		)
		require.NoError(t, err)
		req.Header.Set(ContentTypeHeaderKey, contentTypeXML)

		var x example
		assert.NoError(t, encoderDecoder.DecodeRequest(ctx, req, &x))
		assert.Equal(t, x.Name, e.Name)
	})

	T.Run("as Emoji", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		expectation := "name"
		e := &example{Name: expectation}
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeJSON)

		bs, err := marshalEmoji(e)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"https://whatever.whocares.gov",
			bytes.NewReader(bs),
		)
		require.NoError(t, err)
		req.Header.Set(ContentTypeHeaderKey, contentTypeEmoji)

		var x example
		assert.NoError(t, encoderDecoder.DecodeRequest(ctx, req, &x))
		assert.Equal(t, x.Name, e.Name)
	})
}
