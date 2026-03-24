package grpc

import (
	"testing"

	issuereportmock "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/issuereports/mock"
	issuereportssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/issue_reports"

	"github.com/stretchr/testify/assert"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		issueReportsManager := &issuereportmock.Repository{}

		service := NewService(logger, tracerProvider, issueReportsManager)

		assert.NotNil(t, service)
		assert.Implements(t, (*issuereportssvc.IssueReportsServiceServer)(nil), service)
	})
}
