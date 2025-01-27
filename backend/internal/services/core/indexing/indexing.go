package indexing

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	ErrNilIndexRequest = errors.New("nil index request")
)

func HandleIndexRequest(
	ctx context.Context,
	l logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	searchConfig *textsearchcfg.Config,
	dataManager database.DataManager,
	indexReq *textsearch.IndexRequest,
) error {
	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("core_search_indexer"))
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	if indexReq == nil {
		return observability.PrepareAndLogError(ErrNilIndexRequest, l, span, "handling index requests")
	}

	logger := l.WithValue("index_type_requested", indexReq.IndexType)

	var (
		im                textsearch.IndexManager
		toBeIndexed       any
		err               error
		markAsIndexedFunc func() error
	)

	switch indexReq.IndexType {
	case textsearch.IndexTypeUsers:
		im, err = textsearchcfg.ProvideIndex[UserSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, indexReq.IndexType)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var user *types.User
		user, err = dataManager.GetUser(ctx, indexReq.RowID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting user")
		}

		toBeIndexed = ConvertUserToUserSearchSubset(user)
		markAsIndexedFunc = func() error { return dataManager.MarkUserAsIndexed(ctx, indexReq.RowID) }
	default:
		logger.Info("invalid index type specified, exiting")
		return nil
	}

	if indexReq.Delete {
		if err = im.Delete(ctx, indexReq.RowID); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "deleting data")
		}

		return nil
	} else {
		if err = im.Index(ctx, indexReq.RowID, toBeIndexed); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "indexing data")
		}

		if err = markAsIndexedFunc(); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "marking data as indexed")
		}
	}

	return nil
}
