package authentication

import (
	"testing"
	"time"

	mockauthn "github.com/prixfixeco/backend/internal/authentication/mock"
	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/featureflags"
	mockpublishers "github.com/prixfixeco/backend/internal/messagequeue/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/random"
	mocktypes "github.com/prixfixeco/backend/pkg/types/mock"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	logger := logging.NewNoopLogger()
	encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	cfg := &Config{
		Cookies: CookieConfig{
			Name:     DefaultCookieName,
			BlockKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!",
			Domain:   ".whocares.gov",
		},
		PASETO: PASETOConfig{
			Issuer:       "test",
			LocalModeKey: []byte("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
			Lifetime:     time.Hour,
		},
	}

	pp := &mockpublishers.ProducerProvider{}
	pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	s, err := ProvideService(
		logger,
		cfg,
		&mockauthn.Authenticator{},
		&mocktypes.UserDataManager{},
		&mocktypes.APIClientDataManager{},
		&mocktypes.HouseholdUserMembershipDataManager{},
		scs.New(),
		encoderDecoder,
		tracing.NewNoopTracerProvider(),
		pp,
		random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
		&email.MockEmailer{},
		&featureflags.NoopFeatureFlagManager{},
	)
	require.NoError(t, err)

	return s.(*service)
}

func TestProvideService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		logger := logging.NewNoopLogger()
		encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		cfg := &Config{
			Cookies: CookieConfig{
				Name:     DefaultCookieName,
				BlockKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!",
			},
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logger,
			cfg,
			&mockauthn.Authenticator{},
			&mocktypes.UserDataManager{},
			&mocktypes.APIClientDataManager{},
			&mocktypes.HouseholdUserMembershipDataManager{},
			scs.New(),
			encoderDecoder,
			tracing.NewNoopTracerProvider(),
			pp,
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			&email.MockEmailer{},
			&featureflags.NoopFeatureFlagManager{},
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})

	T.Run("with invalid cookie key", func(t *testing.T) {
		t.Parallel()
		logger := logging.NewNoopLogger()
		encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		cfg := &Config{
			Cookies: CookieConfig{
				Name:     DefaultCookieName,
				BlockKey: "BLAHBLAHBLAH",
			},
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

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
			tracing.NewNoopTracerProvider(),
			pp,
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			&email.MockEmailer{},
			&featureflags.NoopFeatureFlagManager{},
		)

		assert.Nil(t, s)
		assert.Error(t, err)
	})
}
