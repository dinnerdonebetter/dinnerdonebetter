package serverimpl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/authorization"
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
	authHeaderName = "Authorization"
	tokenPrefix    = "Bearer "

	zuckModeUserHeader      = "X-Zuck-Mode-User"
	zuckModeHouseholdHeader = "X-Zuck-Mode-Household"

	SessionContextKey contextKey = "session_context"
)

var (
	// ErrUserNotAuthorizedToImpersonateOthers is returned when a user is not authorized to impersonate others.
	ErrUserNotAuthorizedToImpersonateOthers = errors.New("user not authorized to impersonate others")
)

func (s *Server) fetchSessionContext(ctx context.Context) *sessions.ContextData {
	sessionContext, ok := ctx.Value(SessionContextKey).(*sessions.ContextData)
	if !ok {
		return nil
	}

	return sessionContext
}

func (s *Server) determineZuckMode(ctx context.Context, metadata metadata.MD, sessionContextData *sessions.ContextData) (userID, householdID string, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if zuckUserHeaders := metadata.Get(zuckModeUserHeader); len(zuckUserHeaders) > 0 {
		var (
			zuckUserID      = zuckUserHeaders[0]
			zuckHouseholdID string
		)

		if !sessionContextData.ServiceRolePermissionChecker().CanImpersonateUsers() {
			return "", "", ErrUserNotAuthorizedToImpersonateOthers
		}

		if _, err = s.dataManager.GetUser(ctx, zuckUserID); err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching user info for zuck mode")
			return "", "", err
		}

		if zuckHouseholdIDs := metadata.Get(zuckModeHouseholdHeader); len(zuckHouseholdIDs) > 0 {
			zuckHouseholdID = zuckHouseholdIDs[0]
			householdID, err = s.dataManager.GetDefaultHouseholdIDForUser(ctx, zuckUserID)
			if err != nil {
				observability.AcknowledgeError(err, logger, span, "fetching household info for zuck mode")
				return "", "", err
			}
		} else {
			return zuckUserID, zuckHouseholdID, nil
		}

		return zuckUserID, householdID, nil
	}

	return "", "", nil
}

func (s *Server) extractSessionContextDataFromOAuth2(ctx context.Context, metadata metadata.MD) (*sessions.ContextData, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	authHeader := metadata.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	accessToken := strings.TrimPrefix(authHeader[0], tokenPrefix)

	token, err := s.oauth2Server.Manager.LoadAccessToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("loading access token: %w", err)
	}

	if userID := token.GetUserID(); userID != "" {
		sessionCtxData, err := s.dataManager.BuildSessionContextDataForUser(ctx, userID)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "fetching user info for cookie")
		}

		zuckUserID, zuckHouseholdID, zuckErr := s.determineZuckMode(ctx, metadata, sessionCtxData)
		if zuckErr != nil {
			return nil, observability.PrepareAndLogError(zuckErr, logger, span, "fetching user info for zuck mode")
		}

		if zuckUserID != "" {
			sessionCtxData.Requester.UserID = zuckUserID
		}

		if zuckHouseholdID != "" {
			sessionCtxData.ActiveHouseholdID = zuckHouseholdID
			sessionCtxData.HouseholdPermissions[zuckHouseholdID] = authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdMemberRole.String())
		}

		if sessionCtxData != nil {
			return sessionCtxData, nil
		}
	}

	return nil, Unauthenticated("invalid OAuth2 token")
}

func (s *Server) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger := s.logger.WithValue("grpc.method", info.FullMethod)

		switch info.FullMethod {
		// these methods don't require authentication
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

		authHeader := md.Get(authHeaderName)
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		sessionContextData, err := s.extractSessionContextDataFromOAuth2(ctx, md)
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

func (s *Server) GetSelf(ctx context.Context, request *messages.GetSelfRequest) (*messages.GetSelfResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
