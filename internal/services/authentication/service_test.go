package authentication

import (
	"testing"
	"time"

	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"

	"go.opentelemetry.io/otel/trace"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mockauthn "github.com/prixfixeco/api_server/internal/authentication/mock"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	logger := logging.NewNoopLogger()
	encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	cfg := &Config{
		Cookies: CookieConfig{
			Name:     DefaultCookieName,
			BlockKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!",
			Domain:   ".prixfixe.dev",
		},
		PASETO: PASETOConfig{
			Issuer:       "test",
			LocalModeKey: []byte("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
			Lifetime:     time.Hour,
		},
	}

	pp := &mockpublishers.ProducerProvider{}
	pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	s, err := ProvideService(
		logger,
		cfg,
		&mockauthn.Authenticator{},
		&mocktypes.UserDataManager{},
		&mocktypes.APIClientDataManager{},
		&mocktypes.HouseholdUserMembershipDataManager{},
		scs.New(),
		encoderDecoder,
		trace.NewNoopTracerProvider(),
		pp,
	)
	require.NoError(t, err)

	return s.(*service)
}

func TestProvideService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		logger := logging.NewNoopLogger()
		encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		cfg := &Config{
			Cookies: CookieConfig{
				Name:     DefaultCookieName,
				BlockKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!",
			},
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logger,
			cfg,
			&mockauthn.Authenticator{},
			&mocktypes.UserDataManager{},
			&mocktypes.APIClientDataManager{},
			&mocktypes.HouseholdUserMembershipDataManager{},
			scs.New(),
			encoderDecoder,
			trace.NewNoopTracerProvider(),
			pp,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})

	T.Run("with invalid cookie key", func(t *testing.T) {
		t.Parallel()
		logger := logging.NewNoopLogger()
		encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		cfg := &Config{
			Cookies: CookieConfig{
				Name:     DefaultCookieName,
				BlockKey: "BLAHBLAHBLAH",
			},
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logger,
			&Config{
				Cookies: CookieConfig{
					Name:     DefaultCookieName,
					BlockKey: "BLAHBLAHBLAH",
				},
			},
			&mockauthn.Authenticator{},
			&mocktypes.UserDataManager{},
			&mocktypes.APIClientDataManager{},
			&mocktypes.HouseholdUserMembershipDataManager{},
			scs.New(),
			encoderDecoder,
			trace.NewNoopTracerProvider(),
			pp,
		)

		assert.Nil(t, s)
		assert.Error(t, err)
	})
}
