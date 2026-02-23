package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer       = "dinner-done-better"
	accountIDKey = "account_id"
)

type (
	signer struct {
		tracer     tracing.Tracer
		logger     logging.Logger
		audience   string
		signingKey []byte
	}
)

func NewJWTSigner(logger logging.Logger, tracerProvider tracing.TracerProvider, audience string, signingKey []byte) (tokens.Issuer, error) {
	s := &signer{
		audience:   audience,
		signingKey: signingKey,
		logger:     logging.EnsureLogger(logger),
		tracer:     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("jwt_signer")),
	}

	return s, nil
}

// IssueToken issues a new JSON web token.
func (s *signer) IssueToken(ctx context.Context, user tokens.User, expiry time.Duration) (string, error) {
	return s.IssueTokenWithAccount(ctx, user, expiry, "")
}

// IssueTokenWithAccount issues a new JSON web token, optionally including an account ID claim.
func (s *signer) IssueTokenWithAccount(ctx context.Context, user tokens.User, expiry time.Duration, accountID string) (string, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	if expiry <= 0 {
		expiry = time.Minute * 10
	}

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(expiry).UTC()),           /* expiration time */
		"nbf": jwt.NewNumericDate(time.Now().Add(-1 * time.Minute).UTC()), /* not before */
		"iat": jwt.NewNumericDate(time.Now().UTC()),                       /* issued at */
		"aud": s.audience,                                                 /* audience, i.e. server address */
		"iss": issuer,                                                     /* issuer */
		"sub": user.GetID(),                                               /* subject */
		"jti": identifiers.New(),                                          /* JWT ID */
	}
	if accountID != "" {
		claims[accountIDKey] = accountID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		// this error actually cannot happen with SigningMethodHS256, but it can with other methods.
		return "", err
	}

	return tokenString, nil
}

// ParseUserIDFromToken parses a AccessToken and returns the associated user ID.
func (s *signer) ParseUserIDFromToken(ctx context.Context, token string) (string, error) {
	userID, _, err := s.ParseUserIDAndAccountIDFromToken(ctx, token)
	return userID, err
}

// ParseUserIDAndAccountIDFromToken parses a JWT and returns the user ID and optional account ID.
func (s *signer) ParseUserIDAndAccountIDFromToken(ctx context.Context, tokenString string) (userID, accountID string, err error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.signingKey, nil
	})
	if err != nil {
		s.logger.Error("parsing JWT", err)
		return "", "", err
	}

	subject, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return "", "", err
	}

	// Extract optional account_id claim
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if aid, ok2 := claims[accountIDKey].(string); ok2 && aid != "" {
			return subject, aid, nil
		}
	}

	return subject, "", nil
}
