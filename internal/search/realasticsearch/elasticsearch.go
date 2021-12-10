package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/search"
)

var _ search.IndexManager = (*indexManager)(nil)

type (
	indexManager struct {
		logger       logging.Logger
		tracer       tracing.Tracer
		esclient     *elasticsearch.Client
		indexName    string
		timeout      time.Duration
		searchFields []string
	}

	indexManagerProvider struct {
		esclient *elasticsearch.Client
	}
)

// NewIndexManagerProvider instantiates an Elasticsearch client.
func NewIndexManagerProvider(
	logger logging.Logger,
	cfg *search.Config,
) (search.IndexManagerProvider, error) {
	if !elasticsearchIsReady(cfg, logger, 10) {
		return nil, errors.New("elasticsearch isn't ready")
	}

	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			string(cfg.Address),
		},
		Username:             cfg.Username,
		Password:             cfg.Password,
		RetryOnStatus:        nil,
		EnableRetryOnTimeout: true,
		MaxRetries:           10,
		Transport:            nil,
		Logger:               nil,
	})
	if err != nil {
		return nil, fmt.Errorf("initializing search client: %w", err)
	}

	im := &indexManagerProvider{esclient: c}

	return im, nil
}

func (m *indexManagerProvider) ProvideIndexManager(ctx context.Context, logger logging.Logger, name search.IndexName, fields ...string) (search.IndexManager, error) {
	im := &indexManager{
		tracer:       tracing.NewTracer("search"),
		logger:       logging.EnsureLogger(logger).WithName(string(name)),
		searchFields: fields,
		esclient:     m.esclient,
		timeout:      30 * time.Second,
		indexName:    string(name),
	}

	if indexErr := im.ensureIndices(ctx); indexErr != nil {
		return nil, indexErr
	}

	return im, nil
}

func elasticsearchIsReady(
	cfg *search.Config,
	l logging.Logger,
	maxAttempts uint8,
) (ready bool) {
	attemptCount := 0

	logger := l.WithValues(map[string]interface{}{
		"interval": time.Second.String(),
		"address":  cfg.Address,
	})

	for !ready {
		_, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{
				string(cfg.Address),
			},
			Username:             cfg.Username,
			Password:             cfg.Password,
			RetryOnStatus:        nil,
			EnableRetryOnTimeout: true,
			MaxRetries:           10,
			Transport:            observability.HTTPClient().Transport,
			Logger:               nil,
		})
		if err != nil {
			logger.WithValue("attempt_count", attemptCount).Debug("ping failed, waiting for elasticsearch")
			time.Sleep(time.Second)

			attemptCount++
			if attemptCount >= int(maxAttempts) {
				break
			}
		} else {
			ready = true
			return ready
		}
	}

	return false
}

func (sm *indexManager) ensureIndices(ctx context.Context) error {
	_, span := sm.tracer.StartSpan(ctx)
	defer span.End()

	logger := sm.logger.WithValue("index", sm.indexName)

	res, err := esapi.IndicesExistsRequest{
		Index:             []string{sm.indexName},
		IgnoreUnavailable: esapi.BoolPtr(false),
		ErrorTrace:        false,
		FilterPath:        nil,
	}.Do(ctx, sm.esclient)
	if err != nil {
		return observability.PrepareError(err, logger, span, "checking index existence")
	}

	var x map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&x)

	if res.StatusCode == http.StatusNotFound {
		if _, err = (esapi.IndicesCreateRequest{Index: strings.ToLower(sm.indexName)}).Do(ctx, sm.esclient); err != nil {
			return observability.PrepareError(err, logger, span, "checking index existence")
		}
	}

	return nil
}
