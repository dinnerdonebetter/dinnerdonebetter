package integration

import (
	"context"
	"os"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/suite"
)

const (
	oauth2AuthType = "oauth2"
)

type testClientWrapper struct {
	user                    *types.User
	grpcClient              service.EatingServiceClient
	userClient, adminClient service.EatingServiceClient
	authType                string
}

func TestIntegration(T *testing.T) {
	T.Parallel()

	s := new(TestSuite)
	s.address = os.Getenv(serviceURLEnvVarKey) // it's this or global variables

	suite.Run(T, s)
}

type TestSuite struct {
	suite.Suite
	ctx                context.Context
	grpcClient         service.EatingServiceClient
	oauthedClient      service.EatingServiceClient
	adminOAuthedClient service.EatingServiceClient
	user               *types.User
	address            string
}

var _ suite.SetupTestSuite = (*TestSuite)(nil)

func (s *TestSuite) SetupTest() {
	t := s.T()
	testName := t.Name()

	ctx, span := tracing.StartCustomSpan(context.Background(), testName)
	defer span.End()

	s.ctx, _ = tracing.StartCustomSpan(ctx, testName)
	s.user, s.oauthedClient = createUserAndClientForTest(s.ctx, t, s.address, nil)
	s.adminOAuthedClient = buildAdminCookieAndOAuthedClients(s.ctx, s.address, t)
}

/*
var _ suite.TearDownAllSuite = (*TestSuite)(nil)

func (s *TestSuite) TearDownSuite() {
	t := s.T()

	res, err := http.Post("http://coordination_server:9999/completed/golang", "application/json", http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}
*/

func (s *TestSuite) runTest(name string, subtestBuilder func(*testClientWrapper) func()) {
	s.Run(name, subtestBuilder(&testClientWrapper{
		authType:    oauth2AuthType,
		userClient:  s.oauthedClient,
		adminClient: s.adminOAuthedClient,
		user:        s.user,
	}))
}
