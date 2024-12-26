package textsearchcfg

import (
	"context"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/search/text"
	"github.com/dinnerdonebetter/backend/internal/search/text/algolia"
	"github.com/dinnerdonebetter/backend/internal/search/text/elasticsearch"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ElasticsearchProvider represents the elasticsearch search index provider.
	ElasticsearchProvider = "elasticsearch"
	// AlgoliaProvider represents the algolia search index provider.
	AlgoliaProvider = "algolia"
)

// Config contains settings regarding search indices.
type Config struct {
	_ struct{} `json:"-"`

	Algolia              *algolia.Config         `envPrefix:"ALGOLIA_"         json:"algolia"`
	Elasticsearch        *elasticsearch.Config   `envPrefix:"ELASTICSEARCH_"   json:"elasticsearch"`
	CircuitBreakerConfig *circuitbreaking.Config `envPrefix:"CIRCUIT_BREAKER_" json:"circuitBreakerConfig"`
	Provider             string                  `env:"PROVIDER"               json:"provider"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Provider, validation.In(ElasticsearchProvider, AlgoliaProvider)),
		validation.Field(&cfg.Algolia, validation.When(cfg.Provider == AlgoliaProvider, validation.Required)),
		validation.Field(&cfg.Elasticsearch, validation.When(cfg.Provider == ElasticsearchProvider, validation.Required)),
	)
}

// ProvideIndex validates a Config struct.
func ProvideIndex[T textsearch.Searchable](ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *Config, indexName string) (textsearch.Index[T], error) {
	circuitBreaker := circuitbreaking.ProvideCircuitBreaker(cfg.CircuitBreakerConfig)

	switch strings.TrimSpace(strings.ToLower(cfg.Provider)) {
	case ElasticsearchProvider:
		return elasticsearch.ProvideIndexManager[T](ctx, logger, tracerProvider, cfg.Elasticsearch, indexName, circuitBreaker)
	case AlgoliaProvider:
		return algolia.ProvideIndexManager[T](ctx, logger, tracerProvider, cfg.Algolia, indexName, circuitBreaker)
	default:
		return &textsearch.NoopIndexManager[T]{}, nil
	}
}
