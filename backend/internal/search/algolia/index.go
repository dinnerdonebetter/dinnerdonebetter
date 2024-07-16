package algolia

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	objectIDKey = "objectID"
	idKey       = "id"
)

var (
	// ErrEmptyQueryProvided indicates an empty query was provided as input.
	ErrEmptyQueryProvided = errors.New("empty search query provided")
)

// Index implements our indexManager interface.
func (m *indexManager[T]) Index(ctx context.Context, id string, value any) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if m.circuitBreaker.CannotProceed() {
		return types.ErrCircuitBroken
	}

	logger := m.logger.WithValue(idKey, id).WithValue("value", value)
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
	newValue[objectIDKey] = newValue[idKey]
	delete(newValue, idKey)

	if _, err = m.client.SaveObject(newValue); err != nil {
		return err
	}

	return nil
}

// Search implements our indexManager interface.
func (m *indexManager[T]) Search(ctx context.Context, query string) ([]*T, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if m.circuitBreaker.CannotProceed() {
		return nil, types.ErrCircuitBroken
	}

	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	logger := m.logger.WithValue(keys.SearchQueryKey, query)

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	res, searchErr := m.client.Search(query)
	if searchErr != nil {
		m.circuitBreaker.Failed()
		return nil, searchErr
	}

	results := []*T{}
	for _, hit := range res.Hits {
		var x *T

		// we make the same assumption here, sort of
		if _, ok := hit[objectIDKey]; ok {
			hit[idKey] = hit[objectIDKey]
			delete(hit, objectIDKey)
		}

		var encodedAsJSON []byte
		encodedAsJSON, err := json.Marshal(hit)
		if err != nil {
			return nil, err
		}

		if unmarshalErr := json.Unmarshal(encodedAsJSON, &x); unmarshalErr != nil {
			return nil, unmarshalErr
		}

		results = append(results, x)
	}

	logger.Debug("search performed")

	m.circuitBreaker.Succeeded()
	return results, nil
}

// Delete implements our indexManager interface.
func (m *indexManager[T]) Delete(ctx context.Context, id string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if m.circuitBreaker.CannotProceed() {
		return types.ErrCircuitBroken
	}

	logger := m.logger.WithValue(idKey, id)

	if _, err := m.client.DeleteObject(id); err != nil {
		m.circuitBreaker.Failed()
		return err
	}

	logger.Debug("removed from index")

	m.circuitBreaker.Succeeded()
	return nil
}

// Wipe implements our indexManager interface.
func (m *indexManager[T]) Wipe(ctx context.Context) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if m.circuitBreaker.CannotProceed() {
		return types.ErrCircuitBroken
	}

	if _, err := m.client.ClearObjects(); err != nil {
		m.circuitBreaker.Failed()
		return err
	}

	m.circuitBreaker.Succeeded()
	return nil
}
