package algolia

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

var (
	// ErrEmptyQueryProvided indicates an empty query was provided as input.
	ErrEmptyQueryProvided = errors.New("empty search query provided")
)

// Index implements our IndexManager interface.
func (m *IndexManager[T]) Index(ctx context.Context, id string, value any) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithValue("id", id).WithValue("value", value)
	logger.Debug("adding to index")

	jsonEncoded, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var newValue map[string]any
	if unmarshalErr := json.Unmarshal(jsonEncoded, &newValue); unmarshalErr != nil {
		return unmarshalErr
	}

	// we make a huge, albeit safe assumption here.
	newValue["objectID"] = newValue["id"]
	delete(newValue, "id")

	if _, err = m.client.SaveObject(newValue); err != nil {
		return err
	}

	return nil
}

// Search implements our IndexManager interface.
func (m *IndexManager[T]) Search(ctx context.Context, query string) ([]*T, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachSearchQueryToSpan(span, query)
	logger := m.logger.WithValue(keys.SearchQueryKey, query)

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	res, err := m.client.Search(query)
	if err != nil {
		return nil, err
	}

	results := []*T{}
	for _, hit := range res.Hits {
		var x *T

		var encodedAsJSON []byte
		encodedAsJSON, err = json.Marshal(hit)
		if err != nil {
			return nil, err
		}

		if unmarshalErr := json.Unmarshal(encodedAsJSON, &x); unmarshalErr != nil {
			return nil, unmarshalErr
		}

		results = append(results, x)
	}

	logger.Debug("search performed")

	return results, nil
}

// Delete implements our IndexManager interface.
func (m *IndexManager[T]) Delete(ctx context.Context, id string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithValue("id", id)

	if _, err := m.client.DeleteObject(id); err != nil {
		return err
	}

	logger.Debug("removed from index")

	return nil
}

// Wipe implements our IndexManager interface.
func (m *IndexManager[T]) Wipe(ctx context.Context) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if _, err := m.client.ClearObjects(); err != nil {
		return err
	}

	return nil
}
