package serverimpl

import (
	"context"
	"strings"

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

func (s *Server) AuthInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
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

	return handler(context.WithValue(ctx, SessionContextKey, sessionContextData), req)
}
