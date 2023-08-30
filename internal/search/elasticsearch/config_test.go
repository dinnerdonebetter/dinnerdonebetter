package elasticsearch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_provideElasticsearchClient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		cfg := &Config{}

		esClient, err := cfg.provideElasticsearchClient()
		assert.NoError(t, err)
		assert.NotNil(t, esClient)
	})
}
