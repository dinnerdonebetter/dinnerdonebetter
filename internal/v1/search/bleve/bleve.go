package bleve

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.com/prixfixe/prixfixe/internal/v1/search"
	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	bleve "github.com/blevesearch/bleve"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	base    = 10
	bitSize = 64

	// testingSearchIndexName is an index name that is only valid for testing's sake.
	testingSearchIndexName search.IndexName = "testing"
)

var _ search.IndexManager = (*bleveIndexManager)(nil)

type (
	bleveIndexManager struct {
		index  bleve.Index
		logger logging.Logger
	}
)

// NewBleveIndexManager instantiates a bleve index
func NewBleveIndexManager(path search.IndexPath, name search.IndexName, logger logging.Logger) (search.IndexManager, error) {
	var index bleve.Index

	preexistingIndex, openIndexErr := bleve.Open(string(path))
	switch openIndexErr {
	case nil:
		index = preexistingIndex
	case bleve.ErrorIndexPathDoesNotExist:
		logger.WithValue("path", path).Debug("tried to open existing index, but didn't find it")
		var newIndexErr error

		switch name {
		case testingSearchIndexName:
			index, newIndexErr = bleve.New(string(path), bleve.NewIndexMapping())
			if newIndexErr != nil {
				logger.Error(newIndexErr, "failed to create new index")
				return nil, newIndexErr
			}
		case models.ValidInstrumentsSearchIndexName:
			index, newIndexErr = bleve.New(string(path), buildValidInstrumentMapping())
			if newIndexErr != nil {
				logger.Error(newIndexErr, "failed to create new index")
				return nil, newIndexErr
			}
		case models.ValidIngredientsSearchIndexName:
			index, newIndexErr = bleve.New(string(path), buildValidIngredientMapping())
			if newIndexErr != nil {
				logger.Error(newIndexErr, "failed to create new index")
				return nil, newIndexErr
			}
		case models.ValidPreparationsSearchIndexName:
			index, newIndexErr = bleve.New(string(path), buildValidPreparationMapping())
			if newIndexErr != nil {
				logger.Error(newIndexErr, "failed to create new index")
				return nil, newIndexErr
			}
		default:
			return nil, fmt.Errorf("invalid index name: %q", name)
		}
	default:
		logger.Error(openIndexErr, "failed to open index")
		return nil, openIndexErr
	}

	im := &bleveIndexManager{
		index:  index,
		logger: logger.WithName(fmt.Sprintf("%s_search", name)),
	}

	return im, nil
}

// Index implements our IndexManager interface
func (sm *bleveIndexManager) Index(ctx context.Context, id uint64, value interface{}) error {
	_, span := tracing.StartSpan(ctx, "Index")
	defer span.End()

	sm.logger.WithValue("id", id).Debug("adding to index")
	return sm.index.Index(strconv.FormatUint(id, base), value)
}

// Search implements our IndexManager interface
func (sm *bleveIndexManager) Search(ctx context.Context, query string) (ids []uint64, err error) {
	_, span := tracing.StartSpan(ctx, "Search")
	defer span.End()

	tracing.AttachSearchQueryToSpan(span, query)
	sm.logger.WithValue("search_query", query).Debug("performing search")

	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchResults, err := sm.index.SearchInContext(ctx, searchRequest)
	if err != nil {
		sm.logger.Error(err, "performing search query")
		return nil, err
	}

	out := []uint64{}
	for _, result := range searchResults.Hits {
		x, err := strconv.ParseUint(result.ID, base, bitSize)
		if err != nil {
			// this should literally never happen
			return nil, err
		}
		out = append(out, x)
	}

	return out, nil
}

// Delete implements our IndexManager interface
func (sm *bleveIndexManager) Delete(ctx context.Context, id uint64) error {
	_, span := tracing.StartSpan(ctx, "Delete")
	defer span.End()

	sm.logger.WithValue("id", id).Debug("removing from index")
	return sm.index.Delete(strconv.FormatUint(id, base))
}
