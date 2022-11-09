package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/search"
)

var _ search.IndexManager = (*indexManager)(nil)

type (
	indexManager struct {
		logger       logging.Logger
		tracer       tracing.Tracer
		esclient     *elasticsearch.Client
		indexName    string
		searchFields []string
		timeout      time.Duration
	}

	indexManagerProvider struct {
		esclient       *elasticsearch.Client
		tracerProvider tracing.TracerProvider
		address        string
	}
)

// NewIndexManagerProvider instantiates an Elasticsearch client.
func NewIndexManagerProvider(
	ctx context.Context,
	logger logging.Logger,
	cfg *search.Config,
	tracerProvider tracing.TracerProvider,
) (search.IndexManagerProvider, error) {
	if !elasticsearchIsReady(ctx, cfg, logger, 10) {
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

	im := &indexManagerProvider{
		esclient:       c,
		tracerProvider: tracerProvider,
		address:        string(cfg.Address),
	}

	return im, nil
}

func (m *indexManagerProvider) ProvideIndexManager(ctx context.Context, logger logging.Logger, name search.IndexName, fields ...string) (search.IndexManager, error) {
	im := &indexManager{
		tracer:       tracing.NewTracer(m.tracerProvider.Tracer(fmt.Sprintf("search_%s", name))),
		logger:       logging.EnsureLogger(logger).WithName(string(name)).WithValue("address", m.address),
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
	ctx context.Context,
	cfg *search.Config,
	l logging.Logger,
	maxAttempts uint8,
) (ready bool) {
	attemptCount := 0

	logger := l.WithValues(map[string]interface{}{
		"interval": time.Second.String(),
		"address":  cfg.Address,
	})

	logger.Debug("checking if elasticsearch is ready")

	for !ready {
		c, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{
				string(cfg.Address),
			},
			Username:             cfg.Username,
			Password:             cfg.Password,
			DiscoverNodesOnStart: true,
			RetryOnStatus:        nil,
			EnableRetryOnTimeout: true,
			MaxRetries:           50,
			RetryBackoff:         func(attempt int) time.Duration { return time.Second },
			Transport:            observability.HTTPClient().Transport,
			Logger:               nil,
		})
		if err != nil {
			logger.WithValue("attempt_count", attemptCount).Debug("client setup failed, waiting for elasticsearch")
			time.Sleep(time.Second)

			attemptCount++
			if attemptCount >= int(maxAttempts) {
				break
			}
		}

		if res, infoReqErr := (esapi.InfoRequest{}).Do(ctx, c); infoReqErr != nil && !res.IsError() {
			logger.WithValue("attempt_count", attemptCount).Debug("ping failed, waiting for elasticsearch")
			time.Sleep(time.Second)

			attemptCount++
			if attemptCount >= int(maxAttempts) {
				break
			}
		} else {
			ready = true
			logger.Debug("elasticsearch is ready")
			return ready
		}
	}

	logger.Debug("elasticsearch is ready")

	return false
}

func (sm *indexManager) ensureIndices(ctx context.Context) error {
	_, span := sm.tracer.StartSpan(ctx)
	defer span.End()

	res, err := esapi.IndicesExistsRequest{
		Index:             []string{sm.indexName},
		IgnoreUnavailable: esapi.BoolPtr(false),
		ErrorTrace:        false,
		FilterPath:        nil,
	}.Do(ctx, sm.esclient)
	if err != nil {
		return observability.PrepareError(err, span, "checking index existence successfully")
	}

	if res.StatusCode == http.StatusNotFound {
		if _, err = (esapi.IndicesCreateRequest{Index: strings.ToLower(sm.indexName)}).Do(ctx, sm.esclient); err != nil {
			return observability.PrepareError(err, span, "checking index existence")
		}
	}

	return nil
}
