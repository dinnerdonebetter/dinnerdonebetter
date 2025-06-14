package accounts

import (
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	mockencoding "github.com/dinnerdonebetter/backend/internal/platform/encoding/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	mockrouting "github.com/dinnerdonebetter/backend/internal/platform/routing/mock"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                       logging.NewNoopLogger(),
		accountDataManager:           &mocktypes.AccountDataManagerMock{},
		accountMembershipDataManager: &mocktypes.AccountUserMembershipDataManagerMock{},
		accountIDFetcher:             func(req *http.Request) string { return "" },
		encoderDecoder:               encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		secretGenerator:              random.NewGenerator(nil, nil),
		tracer:                       tracing.NewTracerForTest("test"),
	}
}

func TestProvideAccountsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			AccountIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			UserIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		msgCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		pp := &mockpublishers.PublisherProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logging.NewNoopLogger(),
			&mocktypes.AccountDataManagerMock{},
			&mocktypes.AccountUserMembershipDataManagerMock{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			tracing.NewNoopTracerProvider(),
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			msgCfg,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing publisher", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()

		msgCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		pp := &mockpublishers.PublisherProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, errors.New("blah"))

		s, err := ProvideService(
			logging.NewNoopLogger(),
			&mocktypes.AccountDataManagerMock{},
			&mocktypes.AccountUserMembershipDataManagerMock{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			tracing.NewNoopTracerProvider(),
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			msgCfg,
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})
}
