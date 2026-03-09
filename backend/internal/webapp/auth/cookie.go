package auth

import (
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
)

// BuildCookie provides a consistent way of constructing an HTTP cookie for session auth.
// See https://www.calhoun.io/securing-cookies-in-go/
func BuildCookie(cfg *cookies.Config, value string) *http.Cookie {
	expiry := time.Now().Add(cfg.Lifetime)

	return &http.Cookie{
		Name:     cfg.CookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.SecureOnly,
		Expires:  expiry,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(expiry).Seconds()),
	}
}
