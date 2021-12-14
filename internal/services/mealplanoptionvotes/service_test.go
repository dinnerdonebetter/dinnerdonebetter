package mealplanoptionvotes

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/customerdata"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	mockrouting "github.com/prixfixeco/api_server/internal/routing/mock"
	mealplanoptionsservice "github.com/prixfixeco/api_server/internal/services/mealplanoptions"
	mealplansservice "github.com/prixfixeco/api_server/internal/services/mealplans"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
)

func buildTestService() *service {
	return &service{
		logger:                        logging.NewNoopLogger(),
		mealPlanOptionVoteDataManager: &mocktypes.MealPlanOptionVoteDataManager{},
		mealPlanOptionVoteIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:                mockencoding.NewMockEncoderDecoder(),
		tracer:                        tracing.NewTracerForTest("test"),
	}
}

func TestProvideMealPlanOptionVotesService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			mealplansservice.MealPlanIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			mealplanoptionsservice.MealPlanOptionIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			MealPlanOptionVoteIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := Config{
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreUpdatesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreArchivesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			ctx,
			logging.NewNoopLogger(),
			&cfg,
			&mocktypes.MealPlanOptionVoteDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing pre-writes producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := Config{
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			ctx,
			logging.NewNoopLogger(),
			&cfg,
			&mocktypes.MealPlanOptionVoteDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing pre-updates producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := Config{
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreUpdatesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			ctx,
			logging.NewNoopLogger(),
			&cfg,
			&mocktypes.MealPlanOptionVoteDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing pre-archives producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := Config{
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreUpdatesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreArchivesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			ctx,
			logging.NewNoopLogger(),
			&cfg,
			&mocktypes.MealPlanOptionVoteDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
