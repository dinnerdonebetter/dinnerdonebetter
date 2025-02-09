package serverimpl

import (
	"context"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			"/eating.EatingService/CreateUser",
			"/eating.EatingService/Ping",
			"/eating.EatingService/VerifyTOTPSecret",
			"/eating.EatingService/LoginForToken":
			logger.Info("skipping authentication for method")
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, Unauthenticated("missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")

		userID, err := s.tokenIssuer.ParseUserIDFromToken(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		sessionContextData, err := s.dataManager.BuildSessionContextDataForUser(ctx, userID)
		if err != nil {
			return nil, status.Error(codes.Internal, "building session context data for user")
		}

		ctx = context.WithValue(ctx, SessionContextKey, sessionContextData)

		return handler(ctx, req)
	}
}

func (s *Server) ExchangeToken(ctx context.Context, request *messages.ExchangeTokenRequest) (*messages.ExchangeTokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	newToken, err := s.authManager.ExchangeTokenForUser(ctx, request.RefreshToken)
	if err != nil {
		return nil, Unauthenticated("invalid token")
	}

	output := &messages.ExchangeTokenResponse{
		UserID:       newToken.UserID,
		HouseholdID:  newToken.HouseholdID,
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiresUTC:   timestamppb.New(newToken.ExpiresUTC),
	}

	return output, nil
}

func (s *Server) LoginForToken(ctx context.Context, request *messages.LoginForTokenRequest) (*messages.LoginForTokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.Clone()

	tr, err := s.loginForToken(ctx, false, request.Input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating login")
	}

	return &messages.LoginForTokenResponse{Result: tr}, nil
}

func (s *Server) AdminLoginForToken(ctx context.Context, request *messages.AdminLoginForTokenRequest) (*messages.AdminLoginForTokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.Clone()

	tr, err := s.loginForToken(ctx, true, request.Input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating login")
	}

	return &messages.AdminLoginForTokenResponse{Result: tr}, nil
}

func (s *Server) loginForToken(ctx context.Context, admin bool, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	tokenResponse, err := s.authManager.ProcessLogin(ctx, admin, &authentication.UserLoginInput{
		Username:  input.Username,
		Password:  input.Password,
		TOTPToken: input.TOTPToken,
	})
	if err != nil {
		return nil, observability.PrepareError(err, span, "processing login")
	}

	output := &messages.TokenResponse{
		UserID:       tokenResponse.UserID,
		HouseholdID:  tokenResponse.HouseholdID,
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresUTC:   timestamppb.New(tokenResponse.ExpiresUTC),
	}

	return output, nil
}

func (s *Server) CheckPermissions(ctx context.Context, request *messages.CheckPermissionsRequest) (*messages.CheckPermissionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetAuthStatus(ctx context.Context, request *messages.GetAuthStatusRequest) (*messages.GetAuthStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
