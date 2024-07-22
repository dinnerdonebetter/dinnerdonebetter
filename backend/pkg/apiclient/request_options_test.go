package apiclient

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func arbitraryErrorRequestOption(req *http.Request) error {
	return errors.New("blah")
}

func TestClient_applyRequestOptions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		c, _ := buildSimpleTestClient(t)

		assert.NoError(t, c.applyRequestOptions(req, WithHTTPHeader("things", "stuff")))
	})

	T.Run("with error", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.applyRequestOptions(req, arbitraryErrorRequestOption))
	})
}
