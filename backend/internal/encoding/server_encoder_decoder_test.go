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

	"github.com/go-yaml/yaml"
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

	testCases := map[string]struct {
		contentType      ContentType
		expectedResponse string
	}{
		"json": {
			contentType:      ContentTypeJSON,
			expectedResponse: `{"name":"name"}` + "\n",
		},
		"xml": {
			contentType:      ContentTypeXML,
			expectedResponse: "<example><name>name</name></example>",
		},
		"toml": {
			contentType:      ContentTypeTOML,
			expectedResponse: `Name = "name"` + "\n",
		},
		"yaml": {
			contentType:      ContentTypeYAML,
			expectedResponse: "name: name\n",
		},
	}

	for testName, tc := range testCases {
		T.Run(testName, func(t *testing.T) {
			t.Parallel()

			ex := &example{Name: "name"}
			encoderDecoder, ok := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), tc.contentType).(*serverEncoderDecoder)
			require.True(t, ok)

			ctx := context.Background()
			res := httptest.NewRecorder()
			res.Header().Set(ContentTypeHeaderKey, ContentTypeToString(tc.contentType))

			encoderDecoder.encodeResponse(ctx, res, ex, http.StatusOK)
			actual := res.Body.String()
			assert.Equal(t, tc.expectedResponse, actual)
		})
	}

	T.Run("emoji", func(t *testing.T) {
		t.Parallel()

		ex := &example{Name: "name"}
		encoderDecoder, ok := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeEmoji).(*serverEncoderDecoder)
		require.True(t, ok)

		ctx := context.Background()
		res := httptest.NewRecorder()
		res.Header().Set(ContentTypeHeaderKey, ContentTypeToString(ContentTypeEmoji))

		encoderDecoder.encodeResponse(ctx, res, ex, http.StatusOK)
		actual := res.Body.String()
		assert.NotEmpty(t, actual)
	})

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

type exampleError struct {
	Error string `json:"error"`
}

func TestServerEncoderDecoder_EncodeErrorResponse(T *testing.T) {
	T.Parallel()

	const (
		exampleMessage = "something went awry"
		defaultTraceID = "00000000000000000000000000000000"
	)

	testCases := map[string]struct {
		contentType      ContentType
		expectedResponse string
	}{
		"json": {
			contentType:      ContentTypeJSON,
			expectedResponse: fmt.Sprintf(`{"error":{"message":%q,"code":""},"details":{"currentHouseholdID":"","traceID":%q}}`, exampleMessage, defaultTraceID) + "\n",
		},
		"xml": {
			contentType:      ContentTypeXML,
			expectedResponse: fmt.Sprintf("<APIResponse><Error><Message>%s</Message><Code></Code></Error><Details><CurrentHouseholdID></CurrentHouseholdID><TraceID>%s</TraceID></Details></APIResponse>", exampleMessage, defaultTraceID),
		},
		"toml": {
			contentType: ContentTypeTOML,
			expectedResponse: fmt.Sprintf(`[Error]
  Message = %q
  Code = ""

[Details]
  CurrentHouseholdID = ""
  TraceID = %q
`, exampleMessage, defaultTraceID),
		},
		"yaml": {
			contentType: ContentTypeYAML,
			expectedResponse: fmt.Sprintf(`data: null
pagination: null
error:
  message: %s
  code: ""
details:
  currenthouseholdid: ""
  traceid: %q
`, exampleMessage, defaultTraceID),
		},
	}

	for name, tc := range testCases {
		T.Run(name, func(t *testing.T) {
			t.Parallel()

			expectedCode := http.StatusBadRequest
			encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), tc.contentType)

			ctx := context.Background()
			res := httptest.NewRecorder()

			encoderDecoder.EncodeErrorResponse(ctx, res, exampleMessage, expectedCode)
			assert.Equal(t, tc.expectedResponse, res.Body.String())
			assert.Equal(t, expectedCode, res.Code, "expected status code to match")
		})
	}

	T.Run("emoji", func(t *testing.T) {
		t.Parallel()

		expectedCode := http.StatusBadRequest
		encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ContentTypeEmoji)

		ctx := context.Background()
		res := httptest.NewRecorder()

		encoderDecoder.EncodeErrorResponse(ctx, res, exampleMessage, expectedCode)
		assert.NotEmpty(t, res.Body.String())
		assert.Equal(t, expectedCode, res.Code, "expected status code to match")
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

	testCases := map[string]struct {
		contentType ContentType
		expected    string
	}{
		"json": {
			contentType: ContentTypeJSON,
			expected:    `{"name":"TestServerEncoderDecoder_MustEncode/json"}` + "\n",
		},
		"xml": {
			contentType: ContentTypeXML,
			expected:    "<example><name>TestServerEncoderDecoder_MustEncode/xml</name></example>",
		},
		"toml": {
			contentType: ContentTypeTOML,
			expected:    "Name = \"TestServerEncoderDecoder_MustEncode/toml\"\n",
		},
		"yaml": {
			contentType: ContentTypeYAML,
			expected:    "name: TestServerEncoderDecoder_MustEncode/yaml\n",
		},
	}

	for name, tc := range testCases {
		T.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), tc.contentType)

			actual := string(encoderDecoder.MustEncode(ctx, &example{Name: t.Name()}))

			assert.Equal(t, tc.expected, actual)
		})
	}

	T.Run("emoji", func(t *testing.T) {
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

	T.Run("as JSON", func(t *testing.T) {
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

	testCases := map[string]struct {
		contentType ContentType
		marshaller  func(v interface{}) ([]byte, error)
		expected    string
	}{
		"json": {
			contentType: ContentTypeJSON,
			expected:    `{"name":"name"}`,
			marshaller:  json.Marshal,
		},
		"xml": {
			contentType: ContentTypeXML,
			expected:    `<example><name>name</name></example>`,
			marshaller:  xml.Marshal,
		},
		"toml": {
			contentType: ContentTypeTOML,
			expected:    `<example><name>name</name></example>`,
			marshaller:  tomlMarshalFunc,
		},
		"yaml": {
			contentType: ContentTypeYAML,
			expected:    `<example><name>name</name></example>`,
			marshaller:  yaml.Marshal,
		},
		"emoji": {
			contentType: ContentTypeEmoji,
			expected:    `<example><name>name</name></example>`,
			marshaller:  marshalEmoji,
		},
	}

	e := &example{Name: "name"}

	for name, tc := range testCases {
		T.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), tc.contentType)

			bs, err := tc.marshaller(e)
			require.NoError(t, err)

			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodGet,
				"https://whatever.whocares.gov",
				bytes.NewReader(bs),
			)
			require.NoError(t, err)
			req.Header.Set(ContentTypeHeaderKey, ContentTypeToString(tc.contentType))

			var x example
			assert.NoError(t, encoderDecoder.DecodeRequest(ctx, req, &x))
			assert.Equal(t, x.Name, e.Name)
		})
	}
}

func Test_serverEncoderDecoder_DecodeBytes(T *testing.T) {
	T.Parallel()

	goodDataTestCases := map[string]struct {
		contentType ContentType
		data        []byte
	}{
		"json": {
			data:        []byte(`{"name":"name"}`),
			contentType: ContentTypeJSON,
		},
		"xml": {
			data:        []byte(`<example><name>name</name></example>`),
			contentType: ContentTypeXML,
		},
		"toml": {
			data:        []byte(`name = "name"`),
			contentType: ContentTypeTOML,
		},
		"yaml": {
			data:        []byte(`name: "name"`),
			contentType: ContentTypeYAML,
		},
		"emoji": {
			data:        []byte("🍃🧁🌆🙍☔🌾🐯🦮💆🚂🚕🏏🧔✊🀄🏏☔🌊🥈🐾👥♓🙌🀄🀄🍧🦖📓♿😱🦨🐶🀄☕\n"),
			contentType: ContentTypeEmoji,
		},
	}
	goodDataExpectation := &example{Name: "name"}

	for name, tc := range goodDataTestCases {
		T.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			encoderDecoder := ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), tc.contentType)

			var dest *example
			assert.NoError(t, encoderDecoder.DecodeBytes(ctx, tc.data, &dest))

			assert.Equal(t, goodDataExpectation, dest)
		})
	}
}
