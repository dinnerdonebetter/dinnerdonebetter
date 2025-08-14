package grpc

import (
	"context"
	"errors"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	contextKey string
)

const (
	authHeaderName = "Authorization"
	tokenPrefix    = "Bearer "

	zuckModeUserHeader    = "X-Zuck-Mode-User"
	zuckModeAccountHeader = "X-Zuck-Mode-Account"

	SessionContextKey contextKey = "session_context"
)

var (
	// ErrUserNotAuthorizedToImpersonateOthers is returned when a user is not authorized to impersonate others.
	ErrUserNotAuthorizedToImpersonateOthers = errors.New("user not authorized to impersonate others")
)

func Unauthenticated(msg string) error {
	return status.Error(codes.Unauthenticated, msg)
}

func (s *serviceImpl) fetchSessionContext(ctx context.Context) *sessions.ContextData {
	sessionContext, ok := ctx.Value(SessionContextKey).(*sessions.ContextData)
	if !ok {
		return nil
	}

	return sessionContext
}

func (s *serviceImpl) determineZuckMode(ctx context.Context, metadata metadata.MD, sessionContextData *sessions.ContextData) (userID, accountID string, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if zuckUserHeaders := metadata.Get(zuckModeUserHeader); len(zuckUserHeaders) > 0 {
		var (
			zuckUserID    = zuckUserHeaders[0]
			zuckAccountID string
		)

		if !sessionContextData.ServiceRolePermissionChecker().CanImpersonateUsers() {
			return "", "", ErrUserNotAuthorizedToImpersonateOthers
		}

		if _, err = s.identityRepository.GetUser(ctx, zuckUserID); err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching user info for zuck mode")
			return "", "", err
		}

		if zuckAccountIDs := metadata.Get(zuckModeAccountHeader); len(zuckAccountIDs) > 0 {
			zuckAccountID = zuckAccountIDs[0]
			accountID, err = s.identityRepository.GetDefaultAccountIDForUser(ctx, zuckUserID)
			if err != nil {
				observability.AcknowledgeError(err, logger, span, "fetching account info for zuck mode")
				return "", "", err
			}
		} else {
			return zuckUserID, zuckAccountID, nil
		}

		return zuckUserID, accountID, nil
	}

	return "", "", nil
}

func (s *serviceImpl) extractSessionContextDataFromOAuth2(ctx context.Context, metadata metadata.MD) (*sessions.ContextData, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	authHeader := metadata.Get("authorization")
	if len(authHeader) == 0 {
		return nil, observability.PrepareAndLogGRPCStatus(status.Error(codes.Unauthenticated, "missing authorization header"), logger, span, codes.Unauthenticated, "missing authorization header")
	}

	/*
		accessToken := strings.TrimPrefix(authHeader[0], tokenPrefix)

		token, err := s.oauth2ClientManager.LoadAccessToken(ctx, accessToken)
		if err != nil {
			return nil, fmt.Errorf("loading access token: %w", err)
		}
		if userID := token.GetUserID(); userID != "" {
			sessionCtxData, err := s.identityRepository.BuildSessionContextDataForUser(ctx, userID)
			if err != nil {
				return nil, observability.PrepareAndLogError(err, logger, span, "fetching user info for cookie")
			}

			zuckUserID, zuckAccountID, zuckErr := s.determineZuckMode(ctx, metadata, sessionCtxData)
			if zuckErr != nil {
				return nil, observability.PrepareAndLogError(zuckErr, logger, span, "fetching user info for zuck mode")
			}

			if zuckUserID != "" {
				sessionCtxData.Requester.UserID = zuckUserID
			}

			if zuckAccountID != "" {
				sessionCtxData.ActiveAccountID = zuckAccountID
				sessionCtxData.AccountPermissions[zuckAccountID] = authorization.NewAccountRolePermissionChecker(authorization.AccountMemberRole.String())
			}

			if sessionCtxData != nil {
				return sessionCtxData, nil
			}
		}
	*/

	return nil, Unauthenticated("invalid OAuth2 token")
}

func (s *serviceImpl) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger := s.logger.WithValue("grpc.method", info.FullMethod)

		switch info.FullMethod {
		// these methods don't require authentication
		case "/mealplanning.EatingService/AdminLoginForToken",
			"/mealplanning.EatingService/CreateUser",
			"/mealplanning.EatingService/Ping",
			"/mealplanning.EatingService/VerifyTOTPSecret",
			"/mealplanning.EatingService/LoginForToken":
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
