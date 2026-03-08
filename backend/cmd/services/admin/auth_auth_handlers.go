package main

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/webappauth"
	"github.com/dinnerdonebetter/backend/pkg/client"

	g "maragu.dev/gomponents"
)

func (s *AdminFrontendServer) LoginPage(_ http.ResponseWriter, _ *http.Request) (g.Node, error) {
	return page("Login",
		s.componentRenderer.LoginForm(&components.LoginFormProps{}),
	), nil
}

func (s *AdminFrontendServer) LoginSubmission(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	var loginInput *auth.UserLoginInput
	if err := s.encoder.DecodeRequest(ctx, req, &loginInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		return s.componentRenderer.LoginForm(&components.LoginFormProps{GeneralError: err.Error()}), nil
	}

	var (
		usernameError, passwordError, totpError string
	)

	if err := loginInput.ValidateWithContext(ctx); err != nil {
		usernameError = fetchErrorString(err, "username")
		passwordError = fetchErrorString(err, "password")
		totpError = fetchErrorString(err, "totpToken")
	}

	if usernameError != "" || passwordError != "" || totpError != "" {
		return s.componentRenderer.LoginForm(&components.LoginFormProps{
			UsernameError: usernameError,
			PasswordError: passwordError,
			TOTPError:     totpError,
		}), nil
	}

	unauthedClient, err := s.buildUnauthedGRPCClient(ctx)
	if err != nil {
		return s.componentRenderer.LoginForm(&components.LoginFormProps{
			GeneralError: err.Error(),
		}), nil
	}

	tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
		Input: &authsvc.UserLoginInput{
			Username:  loginInput.Username,
			Password:  loginInput.Password,
			TotpToken: loginInput.TOTPToken,
		},
	})
	if err != nil {
		return s.componentRenderer.LoginForm(&components.LoginFormProps{
			GeneralError: err.Error(),
		}), nil
	}

	encodedCookie, err := s.cookieManager.Encode(ctx, s.config.Cookies.CookieName, &webappauth.AuthPayload{AccessToken: tokenRes.Result.AccessToken})
	if err != nil {
		return s.componentRenderer.LoginForm(&components.LoginFormProps{
			GeneralError: err.Error(),
		}), nil
	}

	http.SetCookie(res, webappauth.BuildCookie(&s.config.Cookies, encodedCookie))

	// Redirect to home page after successful login using HTMX redirect
	// This causes HTMX to do a full page redirect instead of swapping content
	res.Header().Set("HX-Redirect", "/")
	res.WriteHeader(http.StatusOK)

	return g.El("div"), nil
}

// buildUnauthedGRPCClient builds an unauthenticated gRPC client for auth calls.
func (s *AdminFrontendServer) buildUnauthedGRPCClient(_ context.Context) (client.Client, error) {
	if s.developingLocally {
		return client.BuildUnauthenticatedGRPCClient(s.config.APIServiceConnection.GRPCAPIServerURL)
	}
	return client.BuildTLSGRPCClient(s.config.APIServiceConnection.GRPCAPIServerURL)
}
