package elasticsearch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: use testcontainers to properly test this: https://golang.testcontainers.org/modules/elasticsearch/

func Test_elasticsearchIsReady(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.True(t, true)
	})
}
