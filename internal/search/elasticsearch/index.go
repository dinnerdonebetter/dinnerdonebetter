package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/olivere/elastic/v7"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

var (
	// ErrEmptyQueryProvided indicates an empty query was provided as input.
	ErrEmptyQueryProvided = errors.New("empty search query provided")
)

type idContainer struct {
	ID string `json:"id"`
}

// Index implements our IndexManager interface.
func (sm *indexManager) Index(ctx context.Context, id string, value interface{}) error {
	_, span := sm.tracer.StartSpan(ctx)
	defer span.End()

	logger := sm.logger.WithValue("id", id).WithValue("value", value)
	logger.Debug("adding to index")

	if _, err := sm.esclient.Index().Index(sm.indexName).Id(id).BodyJson(value).Do(ctx); err != nil {
		return observability.PrepareError(err, logger, span, "indexing value")
	}

	return nil
}

// search executes search queries.
func (sm *indexManager) search(ctx context.Context, byField, query, householdID string) (ids []string, err error) {
	_, span := sm.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachSearchQueryToSpan(span, query)
	logger := sm.logger.WithValue(keys.SearchQueryKey, query)

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	baseQuery := elastic.NewWildcardQuery(byField, fmt.Sprintf("*%s*", query))

	var q elastic.Query
	if householdID == "" {
		q = baseQuery
	} else {
		householdIDMatchQuery := elastic.NewMatchQuery("householdID", householdID)
		q = elastic.NewBoolQuery().Should(householdIDMatchQuery).Should(baseQuery)
	}

	results, err := sm.esclient.Search().Index(sm.indexName).Query(q).Do(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying elasticsearch")
	}

	resultIDs := []string{}
	for _, hit := range results.Hits.Hits {
		var i *idContainer
		if unmarshalErr := json.Unmarshal(hit.Source, &i); unmarshalErr != nil {
			return nil, observability.PrepareError(err, logger, span, "unmarshalling search result")
		}
		resultIDs = append(resultIDs, i.ID)
	}

	return resultIDs, nil
}

// Search implements our IndexManager interface.
func (sm *indexManager) Search(ctx context.Context, byField, query, householdID string) (ids []string, err error) {
	return sm.search(ctx, byField, query, householdID)
}

// Delete implements our IndexManager interface.
func (sm *indexManager) Delete(ctx context.Context, id string) error {
	_, span := sm.tracer.StartSpan(ctx)
	defer span.End()

	logger := sm.logger.WithValue("id", id)

	q := elastic.NewTermQuery("id", id)
	if _, err := sm.esclient.DeleteByQuery(sm.indexName).Query(q).Do(ctx); err != nil {
		return observability.PrepareError(err, logger, span, "deleting from elasticsearch")
	}

	logger.Debug("removed from index")

	return nil
}
