package asyncfunc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"go.opentelemetry.io/otel"
)

func CollectAndSaveUserData(
	ctx context.Context,
	logger logging.Logger,
	cfg *config.InstanceConfig,
	userDataCollectionRequest *types.UserDataAggregationRequest,
) error {
	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error(err, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, tracerProvider, &cfg.Database)
	if err != nil {
		cancel()
		return observability.PrepareAndLogError(err, logger, nil, "establishing database connection")
	}
	cancel()
	defer dataManager.Close()

	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("search_indexer_cloud_function"))
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	logger = logger.WithValue(keys.UserDataAggregationReportIDKey, userDataCollectionRequest.ReportID)
	tracing.AttachToSpan(span, keys.UserDataAggregationReportIDKey, userDataCollectionRequest.ReportID)
	logger.Info("loaded payload, aggregating data")

	collection, err := dataManager.AggregateUserData(ctx, userDataCollectionRequest.UserID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "collecting user data")
	}
	collection.ReportID = userDataCollectionRequest.ReportID

	logger.Info("compiled payload, saving")

	storageConfig := &objectstorage.Config{
		GCPConfig:         &objectstorage.GCPConfig{BucketName: "userdata.dinnerdonebetter.dev"},
		BucketPrefix:      "",
		BucketName:        "userdata.dinnerdonebetter.dev",
		UploadFilenameKey: "",
		Provider:          objectstorage.GCPCloudStorageProvider,
	}

	logger.Info("establishing upload manager")

	uploadManager, err := objectstorage.NewUploadManager(ctx, logger, tracerProvider, storageConfig, chi.NewRouteParamManager())
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "creating upload manager")
	}

	collectionBytes, err := json.Marshal(collection)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marshaling collection")
	}

	logger.Info("saving file")

	if err = uploadManager.SaveFile(ctx, fmt.Sprintf("%s.json", userDataCollectionRequest.ReportID), collectionBytes); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "saving collection")
	}

	return nil
}
