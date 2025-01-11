package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer = "dinner-done-better"
)

type (
	jwtSigner struct {
		tracer     tracing.Tracer
		logger     logging.Logger
		audience   string
		signingKey []byte
	}
)

func NewJWTSigner(logger logging.Logger, tracerProvider tracing.TracerProvider, audience string, signingKey []byte) (tokens.Issuer, error) {
	s := &jwtSigner{
		audience:   audience,
		signingKey: signingKey,
		logger:     logging.EnsureLogger(logger),
		tracer:     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("jwt_signer")),
	}

	return s, nil
}

// IssueToken issues a new JSON web token.
func (s *jwtSigner) IssueToken(ctx context.Context, user *types.User, expiry time.Duration) (string, error) {
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

// ParseUserIDFromToken parses a Token and returns the associated user ID.
func (s *jwtSigner) ParseUserIDFromToken(ctx context.Context, token string) (string, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.signingKey, nil
	})
	if err != nil {
		s.logger.Error("parsing JWT", err)
		return "", err
	}

	subject, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return subject, nil
}
