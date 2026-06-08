package client

import (
	"crypto/tls"
	"fmt"

	mealplanninggrpc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// MealPlanningClient extends the base Client with domain-specific MealPlanningService RPCs.
// Use BuildMealPlanningClient (or its convenience wrappers) instead of BuildClient when the
// caller needs to call meal-planning endpoints.
type MealPlanningClient interface {
	Client
	mealplanninggrpc.MealPlanningServiceClient
}

type mealPlanningClient struct {
	*client
	mealplanninggrpc.MealPlanningServiceClient
}

// BuildMealPlanningClient builds a MealPlanningClient from the given gRPC server address
// and dial options. The returned value satisfies both Client and
// mealplanninggrpc.MealPlanningServiceClient.
func BuildMealPlanningClient(grpcServerAddress string, opts ...grpc.DialOption) (MealPlanningClient, error) {
	conn, err := grpc.NewClient(grpcServerAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("building grpc client: %w", err)
	}

	return &mealPlanningClient{
		client:                    newClientFromConn(conn),
		MealPlanningServiceClient: mealplanninggrpc.NewMealPlanningServiceClient(conn),
	}, nil
}

// BuildUnauthenticatedMealPlanningGRPCClient connects without TLS or auth tokens.
// Use only for plaintext backends (e.g. kubectl port-forward).
func BuildUnauthenticatedMealPlanningGRPCClient(grpcServerAddr string, opts ...grpc.DialOption) (MealPlanningClient, error) {
	return BuildMealPlanningClient(grpcServerAddr, append([]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, opts...)...)
}

// BuildTLSMealPlanningGRPCClient connects with TLS but no auth tokens.
func BuildTLSMealPlanningGRPCClient(grpcServerAddr string, opts ...grpc.DialOption) (MealPlanningClient, error) {
	return BuildMealPlanningClient(grpcServerAddr, append([]grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))}, opts...)...)
}

// BuildUnauthenticatedMealPlanningGRPCClientWithBearerToken connects with a Bearer token (e.g. JWT from LoginForToken).
func BuildUnauthenticatedMealPlanningGRPCClientWithBearerToken(grpcServerAddr, token string) (MealPlanningClient, error) {
	return BuildMealPlanningClient(grpcServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		WithBearerTokenCredentials(token),
	)
}
