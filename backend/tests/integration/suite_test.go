package integration

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated/v2"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/suite"
)

const (
	cookieAuthType = "cookie"
	oauth2AuthType = "oauth2"
)

type testClientWrapper struct {
	user        *types.User
	userClient  *apiclient.Client
	adminClient *apiclient.Client
	authType    string
}

func TestIntegration(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(TestSuite))
}

type TestSuite struct {
	suite.Suite

	ctx  context.Context
	user *types.User
	oauthedClient,
	adminOAuthedClient *apiclient.Client
}

var _ suite.SetupTestSuite = (*TestSuite)(nil)

func (s *TestSuite) SetupTest() {
	t := s.T()
	testName := t.Name()

	ctx, span := tracing.StartCustomSpan(context.Background(), testName)
	defer span.End()

	s.ctx, _ = tracing.StartCustomSpan(ctx, testName)
	s.user, s.oauthedClient = createUserAndClientForTest(s.ctx, t, nil)
	s.adminOAuthedClient = buildAdminCookieAndOAuthedClients(s.ctx, t)
}

func (s *TestSuite) runTest(name string, subtestBuilder func(*testClientWrapper) func()) {
	s.Run(name, subtestBuilder(&testClientWrapper{authType: oauth2AuthType, userClient: s.oauthedClient, adminClient: s.adminOAuthedClient, user: s.user}))
}
