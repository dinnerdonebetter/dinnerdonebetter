package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	contextKey string
)

const (
	SessionContextKey contextKey = "session_context"
)

func (s *Server) fetchSessionContext(ctx context.Context) *sessions.ContextData {
	sessionContext, ok := ctx.Value(SessionContextKey).(*sessions.ContextData)
	if !ok {
		return nil
	}

	return sessionContext
}

func (s *Server) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger := s.logger.WithValue("grpc.method", info.FullMethod)

		switch info.FullMethod {
		// these methods don't require prior authentication
		case "/eating.EatingService/AdminLoginForToken",
			"/eating.EatingService/LoginForToken":
			logger.Info("skipping authentication for method")
			return handler(ctx, req)
		}

		_, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		//authHeader := md.Get("authorization")
		//if len(authHeader) == 0 {
		//	return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		//}
		//
		//token := strings.TrimPrefix(authHeader[0], "Bearer ")

		//userID, err := s.tokenIssuer.ParseUserIDFromToken(ctx, token)
		//if err != nil {
		//	return nil, status.Error(codes.Unauthenticated, "invalid token")
		//}

		sessionContextData := &sessions.ContextData{
			HouseholdPermissions: nil,
			Requester: sessions.RequesterInfo{
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
				AccountStatus:            "good",
				AccountStatusExplanation: "normal",
				UserID:                   "12345asdf",
				EmailAddress:             "fart@butts.com",
				Username:                 "example",
			},
			ActiveHouseholdID: "household123",
		}

		//sessionContextData, err := s.dataManager.BuildSessionContextDataForUser(ctx, userID)
		//if err != nil {
		//	return nil, status.Error(codes.Internal, "building session context data for user")
		//}

		ctx = context.WithValue(ctx, SessionContextKey, sessionContextData)

		return handler(ctx, req)
	}
}

func (s *Server) LoginForToken(ctx context.Context, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) AdminLoginForToken(ctx context.Context, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	// TODO: validation

	user, err := s.dataManager.GetAdminUserByUsername(ctx, input.Username)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching user by username")
	}

	loginValid, err := s.authenticator.CredentialsAreValid(
		ctx,
		user.HashedPassword,
		input.Password,
		user.TwoFactorSecret,
		input.TOTPToken,
	)
	if err != nil {
		return nil, observability.PrepareError(err, span, "validating login")
	}

	if !loginValid {
		return nil, observability.PrepareError(err, span, "invalid login")
	}

	if loginValid && user.TwoFactorSecretVerifiedAt != nil && input.TOTPToken == "" {
		return nil, observability.PrepareError(err, span, "user with two factor verification active attempted to log in without providing TOTP")
	}

	defaultHouseholdID, err := s.dataManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching user memberships")
	}

	var token string
	token, err = s.tokenIssuer.IssueToken(ctx, user, s.config.Services.Auth.TokenLifetime)
	if err != nil {
		return nil, observability.PrepareError(err, span, "signing token")
	}

	output := &messages.TokenResponse{
		UserID:      user.ID,
		HouseholdID: defaultHouseholdID,
		Token:       token,
	}

	return output, nil
}
