package grpc

import (
	"testing"

	uploadedmediamock "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/uploadedmedia/mock"
	uploadedmediasvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/uploaded_media"

	"github.com/stretchr/testify/assert"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing"
	mockuploads "github.com/verygoodsoftwarenotvirus/platform/v2/uploads/mock"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		uploadedMediaManager := &uploadedmediamock.Repository{}

		service := NewService(logger, tracerProvider, uploadedMediaManager, &mockuploads.MockUploadManager{})

		assert.NotNil(t, service)
		assert.Implements(t, (*uploadedmediasvc.UploadedMediaServiceServer)(nil), service)
	})
}
