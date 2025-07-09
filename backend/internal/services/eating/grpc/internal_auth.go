package eatinggrpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"

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

	zuckModeUserHeader      = "X-Zuck-Mode-User"
	zuckModeHouseholdHeader = "X-Zuck-Mode-Household"

	SessionContextKey contextKey = "session_context"
)

func (s *ServiceImpl) fetchSessionContext(ctx context.Context) *sessions.ContextData {
	sessionContext, ok := ctx.Value(SessionContextKey).(*sessions.ContextData)
	if !ok {
		return nil
	}

	return sessionContext
}

func (s *ServiceImpl) determineZuckMode(ctx context.Context, metadata metadata.MD, sessionContextData *sessions.ContextData) (userID, accountID string, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	/*
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
				accountID, err = s.dataManager.GetDefaultHouseholdIDForUser(ctx, zuckUserID)
				if err != nil {
					observability.AcknowledgeError(err, logger, span, "fetching account info for zuck mode")
					return "", "", err
				}
			} else {
				return zuckUserID, zuckHouseholdID, nil
			}

			return zuckUserID, accountID, nil
		}
	*/

	return "", "", nil
}

func (s *ServiceImpl) extractSessionContextDataFromOAuth2(ctx context.Context, metadata metadata.MD) (*sessions.ContextData, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	/*
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
	*/

	return nil, Unauthenticated("invalid OAuth2 token")
}

func (s *ServiceImpl) AuthInterceptor() grpc.UnaryServerInterceptor {
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
