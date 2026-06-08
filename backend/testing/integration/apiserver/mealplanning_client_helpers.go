package integration

// Domain: mealplanning
//
// This file contains test-infrastructure helpers that return client.MealPlanningClient.
// When swapping the meal planning domain out, replace this file with one that returns
// your own domain client type.

import (
	"context"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/pkg/client"

	"github.com/stretchr/testify/require"
)

var (
	// adminClient is the pre-authed service-admin client used across all integration tests.
	// It implements client.MealPlanningClient so meal-planning RPCs can be called directly.
	adminClient client.MealPlanningClient
)

// buildUnauthenticatedGRPCClientForTest returns an unauthenticated MealPlanningClient
// pointed at the test gRPC server.
func buildUnauthenticatedGRPCClientForTest(t *testing.T) client.MealPlanningClient {
	t.Helper()

	c, err := client.BuildUnauthenticatedMealPlanningGRPCClient(fmt.Sprintf(":%d", apiServiceConfig.GRPCServer.Port))
	require.NoError(t, err)

	return c
}

// buildAuthedGRPCClient exchanges the given JWT for an OAuth2 token and returns an
// authenticated MealPlanningClient.
func buildAuthedGRPCClient(ctx context.Context, token string) (client.MealPlanningClient, error) {
	c, err := localdev.BuildInsecureOAuthedMealPlanningGRPCClient(
		ctx,
		createdClientID,
		createdClientSecret,
		httpTestServerAddress,
		fmt.Sprintf(":%d", apiServiceConfig.GRPCServer.Port),
		token,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// buildAuthedGRPCClientWithBearerToken builds a MealPlanningClient that sends the JWT
// directly as a Bearer token. Use this when the token has an account_id claim (e.g. from
// LoginForToken with DesiredAccountId) so that GetAuthStatus uses that account as active.
func buildAuthedGRPCClientWithBearerToken(token string) (client.MealPlanningClient, error) {
	return client.BuildUnauthenticatedMealPlanningGRPCClientWithBearerToken(
		fmt.Sprintf(":%d", apiServiceConfig.GRPCServer.Port),
		token,
	)
}

// createClientForUser fetches a login token for the user and returns an authenticated
// MealPlanningClient.
func createClientForUser(ctx context.Context, user *identity.User) (client.MealPlanningClient, error) {
	token, err := fetchLoginTokenForUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("fetching token for user %s: %w", user.Username, err)
	}

	oauthedClient, err := buildAuthedGRPCClient(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("building oauthed client: %w", err)
	}

	return oauthedClient, nil
}

// createUserAndClientForTest registers a fresh user and returns an authenticated
// MealPlanningClient for that user.
func createUserAndClientForTest(t *testing.T) (*identity.User, client.MealPlanningClient) {
	t.Helper()

	return createUserAndClientForTestWithRegistrationInput(t, buildUserRegistrationInputForTest(t))
}

// createUserAndClientForTestWithRegistrationInput registers a user with the given input
// and returns an authenticated MealPlanningClient for that user.
func createUserAndClientForTestWithRegistrationInput(t *testing.T, input *identity.UserRegistrationInput) (*identity.User, client.MealPlanningClient) {
	t.Helper()

	ctx := t.Context()

	user := createServiceUserForTest(t, true, input)
	oauthedClient, err := buildAuthedGRPCClient(ctx, fetchLoginTokenForUserForTest(t, user))
	require.NoError(t, err)

	return user, oauthedClient
}
