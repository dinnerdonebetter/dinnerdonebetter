package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/search"
)

var _ search.IndexManager = (*indexManager)(nil)

type (
	esClient interface {
		IndexExists(indices ...string) *elastic.IndicesExistsService
		CreateIndex(name string) *elastic.IndicesCreateService
		Search(indices ...string) *elastic.SearchService
		Index() *elastic.IndexService
		DeleteByQuery(indices ...string) *elastic.DeleteByQueryService
	}

	indexManager struct {
		logger       logging.Logger
		tracer       tracing.Tracer
		esclient     esClient
		indexName    string
		searchFields []string
	}

	indexManagerProvider struct {
		esclient esClient
	}
)

// NewIndexManagerProvider instantiates an Elasticsearch client.
func NewIndexManagerProvider(
	logger logging.Logger,
	client *http.Client,
	cfg *search.Config,
) (search.IndexManagerProvider, error) {
	if !elasticsearchIsReady(client, cfg.Address, logger, 50) {
		return nil, errors.New("elasticsearch isn't ready")
	}

	c, err := elastic.NewClient(
		elastic.SetURL(string(cfg.Address)),
		elastic.SetBasicAuth(cfg.Username, cfg.Password),
		elastic.SetHttpClient(client),
		elastic.SetHealthcheck(false),
	)
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
		indexName:    string(name),
	}

	if indexErr := im.ensureIndices(ctx); indexErr != nil {
		return nil, indexErr
	}

	return im, nil
}

func elasticsearchIsReady(
	client *http.Client,
	path search.IndexPath,
	l logging.Logger,
	maxAttempts uint8,
) (ready bool) {
	attemptCount := 0

	logger := l.WithValues(map[string]interface{}{
		"interval": time.Second.String(),
		"address":  path,
	})

	for !ready {
		_, err := elastic.NewClient(
			elastic.SetURL(string(path)),
			elastic.SetHttpClient(client),
		)
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

	indexExists, err := sm.esclient.IndexExists(strings.ToLower(sm.indexName)).Do(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "checking index existence")
	}

	if !indexExists {
		if _, err = sm.esclient.CreateIndex(strings.ToLower(sm.indexName)).Do(ctx); err != nil && !strings.Contains(elastic.ErrorReason(err), "already exists") {
			return observability.PrepareError(err, logger, span, elastic.ErrorReason(err))
		}
	}

	return nil
}
