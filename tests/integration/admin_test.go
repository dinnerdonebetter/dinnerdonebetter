package integration

import (
	"log"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestAdmin_Returns404WhenModifyingUserAccountStatus() {
	s.runForEachClientExcept("should not be possible to ban a user that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			input := fakes.BuildFakeUserAccountStatusUpdateInput()
			input.TargetUserID = nonexistentID

			// Ban user.
			assert.Error(t, testClients.admin.UpdateUserAccountStatus(ctx, input))
		}
	})
}

func (s *TestSuite) TestAdmin_BanningUsers() {
	s.runForEachClientExcept("should be possible to ban users", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var (
				user       *types.User
				userClient *apiclient.Client
			)

			switch testClients.authType {
			case cookieAuthType:
				user, _, userClient, _ = createUserAndClientForTest(ctx, t, nil)
			default:
				log.Panicf("invalid auth type: %q", testClients.authType)
			}

			// Assert that user can access service
			_, initialCheckErr := userClient.GetWebhooks(ctx, nil)
			require.NoError(t, initialCheckErr)

			input := &types.UserAccountStatusUpdateInput{
				TargetUserID: user.ID,
				NewStatus:    string(types.BannedUserAccountStatus),
				Reason:       "testing",
			}

			assert.NoError(t, testClients.admin.UpdateUserAccountStatus(ctx, input))

			// Assert user can no longer access service
			_, subsequentCheckErr := userClient.GetWebhooks(ctx, nil)
			assert.Error(t, subsequentCheckErr)

			// Clean up.
			assert.NoError(t, testClients.admin.ArchiveUser(ctx, user.ID))
		}
	})
}
