package auditlogentries

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	mockencoding "github.com/dinnerdonebetter/backend/internal/lib/encoding/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	mockrouting "github.com/dinnerdonebetter/backend/internal/lib/routing/mock"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                   logging.NewNoopLogger(),
		auditLogEntryDataManager: &mocktypes.AuditLogEntryDataManagerMock{},
		auditLogEntryIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:           encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		tracer:                   tracing.NewTracerForTest("test"),
	}
}

func TestProvideAuditLogEntriesService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			AuditLogEntryIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		s, err := ProvideService(
			logger,
			&mocktypes.AuditLogEntryDataManagerMock{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
