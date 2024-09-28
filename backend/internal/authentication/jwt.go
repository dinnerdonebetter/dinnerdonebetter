package authentication

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type (
	JWTSigner interface {
		IssueJWT() (string, error)
		ParseJWT(token string) (*jwt.Token, error)
	}

	jwtSigner struct {
		signingKey []byte
	}
)

func NewJWTSigner(signingKey []byte) (JWTSigner, error) {
	s := &jwtSigner{
		signingKey: signingKey,
	}

	return s, nil
}

// IssueJWT issues a new JSON web token.
func (s *jwtSigner) IssueJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})

	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		// this error actually cannot happen with SigningMethodHS256.
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses a JWT and returns the associated token.
func (s *jwtSigner) ParseJWT(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	return parsedToken, nil
}
