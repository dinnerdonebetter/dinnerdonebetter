package main

/*
import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads"
)

func handleUserDataRequest(
	ctx context.Context,
	logger logging.Logger,
	tracer tracing.Tracer,
	uploadManager uploads.UploadManager,
	identityRepo identity.Repository,
	userDataCollectionRequest *identity.UserDataAggregationRequest,
) error {
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	logger = logger.WithValue(keys.UserDataAggregationReportIDKey, userDataCollectionRequest.ReportID)
	tracing.AttachToSpan(span, keys.UserDataAggregationReportIDKey, userDataCollectionRequest.ReportID)
	logger.Info("loaded payload, aggregating data")

	collection, err := identityRepo.AggregateUserData(ctx, userDataCollectionRequest.UserID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "collecting user data")
	}
	collection.ReportID = userDataCollectionRequest.ReportID

	logger.Info("compiled payload, saving")

	logger.Info("establishing upload manager")

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
*/
