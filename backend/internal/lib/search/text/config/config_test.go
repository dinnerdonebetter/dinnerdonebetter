package textsearchcfg

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/search/text/elasticsearch"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		cfg := &Config{
			Provider: ElasticsearchProvider,
			Elasticsearch: &elasticsearch.Config{
				Address: t.Name(),
			},
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
