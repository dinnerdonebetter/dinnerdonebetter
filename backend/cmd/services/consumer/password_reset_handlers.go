package main

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/components"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/pkg/client"

	g "maragu.dev/gomponents"
)

func (s *ConsumerFrontendServer) ForgotPasswordPage(_ http.ResponseWriter, _ *http.Request) (g.Node, error) {
	return page("Forgot password",
		s.componentRenderer.ForgotPasswordForm(&components.ForgotPasswordFormProps{}),
	), nil
}

func (s *ConsumerFrontendServer) ForgotPasswordSubmission(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	var input *auth.PasswordResetTokenCreationRequestInput
	if err := s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		return s.componentRenderer.ForgotPasswordForm(&components.ForgotPasswordFormProps{
			GeneralError: err.Error(),
		}), nil
	}
	if input == nil {
		return s.componentRenderer.ForgotPasswordForm(&components.ForgotPasswordFormProps{
			GeneralError: "Invalid request",
		}), nil
	}

	var emailError string
	if err := input.ValidateWithContext(ctx); err != nil {
		emailError = fetchErrorString(err, "emailAddress")
	}
	if emailError != "" {
		return s.componentRenderer.ForgotPasswordForm(&components.ForgotPasswordFormProps{
			EmailError: emailError,
		}), nil
	}

	unauthedClient, err := s.buildUnauthenticatedClient()
	if err != nil {
		return s.componentRenderer.ForgotPasswordForm(&components.ForgotPasswordFormProps{
			GeneralError: err.Error(),
		}), nil
	}

	_, err = unauthedClient.RequestPasswordResetToken(ctx, &authsvc.RequestPasswordResetTokenRequest{
		EmailAddress: input.EmailAddress,
	})
	if err != nil {
		return s.componentRenderer.ForgotPasswordForm(&components.ForgotPasswordFormProps{
			GeneralError: err.Error(),
		}), nil
	}

	// Always show success to avoid email enumeration.
	return s.componentRenderer.ForgotPasswordForm(&components.ForgotPasswordFormProps{
		Success: true,
	}), nil
}

func (s *ConsumerFrontendServer) ResetPasswordPage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	token := req.URL.Query().Get("t")
	if token == "" {
		return page("Reset password",
			s.componentRenderer.ResetPasswordForm(&components.ResetPasswordFormProps{
				MissingToken: true,
			}),
		), nil
	}
	return page("Reset password",
		s.componentRenderer.ResetPasswordForm(&components.ResetPasswordFormProps{
			Token: token,
		}),
	), nil
}

type resetPasswordSubmitInput struct {
	Token           string `json:"token"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (s *ConsumerFrontendServer) ResetPasswordSubmission(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	var input resetPasswordSubmitInput
	if err := s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		return s.componentRenderer.ResetPasswordForm(&components.ResetPasswordFormProps{
			Token:        input.Token,
			GeneralError: err.Error(),
		}), nil
	}

	redemptionInput := &auth.PasswordResetTokenRedemptionRequestInput{
		Token:       input.Token,
		NewPassword: input.NewPassword,
	}
	var passwordError, confirmError string
	if err := redemptionInput.ValidateWithContext(ctx); err != nil {
		passwordError = fetchErrorString(err, "newPassword")
	}
	if input.NewPassword != input.ConfirmPassword {
		confirmError = "Passwords do not match."
	}
	if passwordError != "" || confirmError != "" {
		return s.componentRenderer.ResetPasswordForm(&components.ResetPasswordFormProps{
			Token:         input.Token,
			PasswordError: passwordError,
			ConfirmError:  confirmError,
		}), nil
	}

	unauthedClient, err := s.buildUnauthenticatedClient()
	if err != nil {
		return s.componentRenderer.ResetPasswordForm(&components.ResetPasswordFormProps{
			Token:        input.Token,
			GeneralError: err.Error(),
		}), nil
	}

	_, err = unauthedClient.RedeemPasswordResetToken(ctx, &authsvc.RedeemPasswordResetTokenRequest{
		Token:       input.Token,
		NewPassword: input.NewPassword,
	})
	if err != nil {
		return s.componentRenderer.ResetPasswordForm(&components.ResetPasswordFormProps{
			Token:        input.Token,
			GeneralError: err.Error(),
		}), nil
	}

	res.Header().Set("HX-Redirect", "/login?reset=success")
	res.WriteHeader(http.StatusOK)
	return g.El("div"), nil
}

func (s *ConsumerFrontendServer) buildUnauthenticatedClient() (client.Client, error) {
	if s.developingLocally {
		return client.BuildUnauthenticatedGRPCClient(s.config.APIServiceConnection.GRPCAPIServerURL)
	}
	return client.BuildTLSGRPCClient(s.config.APIServiceConnection.GRPCAPIServerURL)
}
