package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvideHTTPServer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := ProvideStdLibHTTPServer(12345)

		assert.NotNil(t, x)
	})
}
