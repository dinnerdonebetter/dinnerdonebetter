package validinstruments

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/encoding/mock"
	mockpublishers "gitlab.com/prixfixe/prixfixe/internal/messagequeue/publishers/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	mockrouting "gitlab.com/prixfixe/prixfixe/internal/routing/mock"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	mocksearch "gitlab.com/prixfixe/prixfixe/internal/search/mock"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"
)

func buildTestService() *service {
	return &service{
		logger:                     logging.NewNoopLogger(),
		validInstrumentDataManager: &mocktypes.ValidInstrumentDataManager{},
		validInstrumentIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:             mockencoding.NewMockEncoderDecoder(),
		search:                     &mocksearch.IndexManager{},
		tracer:                     tracing.NewTracer("test"),
	}
}

func TestProvideValidInstrumentsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			ValidInstrumentIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := Config{
			SearchIndexPath:      "example/path",
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
			&mocktypes.ValidInstrumentDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
				return &mocksearch.IndexManager{}, nil
			},
			rpm,
			pp,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing pre-writes producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := Config{
			SearchIndexPath:      "example/path",
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
			&mocktypes.ValidInstrumentDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
				return &mocksearch.IndexManager{}, nil
			},
			nil,
			pp,
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing pre-updates producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := Config{
			SearchIndexPath:      "example/path",
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
			&mocktypes.ValidInstrumentDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
				return &mocksearch.IndexManager{}, nil
			},
			nil,
			pp,
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing pre-archives producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := Config{
			SearchIndexPath:      "example/path",
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
			&mocktypes.ValidInstrumentDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
				return &mocksearch.IndexManager{}, nil
			},
			nil,
			pp,
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		s, err := ProvideService(
			ctx,
			logging.NewNoopLogger(),
			&cfg,
			&mocktypes.ValidInstrumentDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
				return nil, errors.New("blah")
			},
			mockrouting.NewRouteParamManager(),
			nil,
		)

		assert.Nil(t, s)
		assert.Error(t, err)
	})
}
