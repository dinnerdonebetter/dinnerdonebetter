package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
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

	b, err := json.Marshal(value)
	if err != nil {
		println("fart")
	}

	res, err := esapi.IndexRequest{
		Index:               sm.indexName,
		DocumentID:          id,
		Body:                bytes.NewReader(b),
		Timeout:             sm.timeout,
		Version:             nil,
		VersionType:         "",
		WaitForActiveShards: "",
		Pretty:              false,
		Human:               false,
		ErrorTrace:          false,
		FilterPath:          nil,
		Header:              nil,
	}.Do(ctx, sm.esclient)
	if err != nil {
		return observability.PrepareError(err, span, "indexing value")
	}

	if res.StatusCode != http.StatusOK {
		println("")
	}

	return nil
}

type matchCondition struct {
	Query string `json:"query"`
}

type matchQuery map[string]matchCondition

type wildcardCondition struct {
	Value string `json:"value"`
}

type wildcardQuery map[string]wildcardCondition

type condition struct {
	Match    matchQuery     `json:"match,omitempty"`
	Wildcard *wildcardQuery `json:"wildcard,omitempty"`
}

type should struct {
	Should []condition `json:"should"`
}

type queryContainer struct {
	Bool should `json:"bool"`
}

type searchQuery struct {
	Query queryContainer `json:"query"`
}

type esHit struct {
	ID         string          `json:"_id"`
	Source     json.RawMessage `json:"_source"`
	Highlights json.RawMessage `json:"highlight"`
	Sort       []interface{}   `json:"sort"`
}

type esResponse struct {
	Hits struct {
		Hits  []*esHit
		Total struct{ Value int }
	}
	Took int
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

	resultIDs := []string{}
	q := searchQuery{
		Query: queryContainer{
			Bool: should{
				Should: []condition{},
			},
		},
	}

	if householdID != "" {
		q.Query.Bool.Should = append(q.Query.Bool.Should, condition{
			Match: map[string]matchCondition{
				"householdID": {Query: householdID},
			},
		})
	}

	q.Query.Bool.Should = append(q.Query.Bool.Should, condition{
		Wildcard: &wildcardQuery{
			byField: wildcardCondition{
				Value: fmt.Sprintf("*%s*", query),
			},
		},
	})

	queryBody, err := json.Marshal(q)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encodign search query")
	}

	res, err := sm.esclient.Search(
		sm.esclient.Search.WithIndex(sm.indexName),
		sm.esclient.Search.WithBody(bytes.NewReader(queryBody)),
	)
	defer func() {
		if res != nil {
			if err = res.Body.Close(); err != nil {
				observability.AcknowledgeError(err, logger, span, "closing response body")
			}
		}
	}()

	if err != nil {
		return nil, observability.PrepareError(err, span, "querying elasticsearch successfully")
	}

	if res.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, observability.PrepareError(err, span, "invalid response from elasticsearch")
		}

		err = errors.New(strings.Join(res.Warnings(), ", "))
		return nil, observability.PrepareError(err, span, "querying elasticsearch")
	}

	var r esResponse
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, observability.PrepareError(err, span, "decoding response")
	}

	for _, hit := range r.Hits.Hits {
		var c *idContainer
		if err = json.Unmarshal(hit.Source, &c); err != nil {
			return nil, observability.PrepareError(err, span, "decoding response")
		}
		resultIDs = append(resultIDs, c.ID)
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

	_, err := esapi.DeleteRequest{
		Index:      sm.indexName,
		DocumentID: id,
	}.Do(ctx, sm.esclient)
	if err != nil {
		return observability.PrepareError(err, span, "deleting from elasticsearch")
	}

	logger.Debug("removed from index")

	return nil
}
