package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer = "dinnerdonebetter"
)

type (
	JWTSigner interface {
		IssueJWT(ctx context.Context, user *types.User, expiry time.Duration) (string, error)
		ParseJWT(ctx context.Context, token string) (*jwt.Token, error)
	}

	jwtSigner struct {
		tracer     tracing.Tracer
		logger     logging.Logger
		audience   string
		signingKey []byte
	}
)

func NewJWTSigner(logger logging.Logger, tracerProvider tracing.TracerProvider, audience string, signingKey []byte) (JWTSigner, error) {
	s := &jwtSigner{
		audience:   audience,
		signingKey: signingKey,
		logger:     logging.EnsureLogger(logger),
		tracer:     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("jwt_signer")),
	}

	return s, nil
}

// IssueJWT issues a new JSON web token.
func (s *jwtSigner) IssueJWT(ctx context.Context, user *types.User, expiry time.Duration) (string, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	if expiry <= 0 {
		expiry = time.Minute * 10
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(expiry).UTC()),           /* expiration time */
		"nbf": jwt.NewNumericDate(time.Now().Add(-1 * time.Minute).UTC()), /* not before */
		"iat": jwt.NewNumericDate(time.Now().UTC()),                       /* issued at */
		"aud": s.audience,                                                 /* audience, i.e. server address */
		"iss": issuer,                                                     /* issuer */
		"sub": user.ID,                                                    /* subject */
		"jti": identifiers.New(),                                          /* JWT ID */
	})

	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		// this error actually cannot happen with SigningMethodHS256, but it can with other methods.
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses a Token and returns the associated token.
func (s *jwtSigner) ParseJWT(ctx context.Context, token string) (*jwt.Token, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.signingKey, nil
	})
	if err != nil {
		s.logger.Error(err, "parsing JWT")
		return nil, err
	}

	return parsedToken, nil
}
