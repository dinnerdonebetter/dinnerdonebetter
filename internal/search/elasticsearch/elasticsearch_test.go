package elasticsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: use testcontainers to properly test this: https://golang.testcontainers.org/modules/elasticsearch/

func Test_elasticsearchIsReady(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.True(t, true)
	})
}
