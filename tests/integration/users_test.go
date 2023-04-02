package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/prixfixeco/backend/internal/authorization"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	testutils "github.com/prixfixeco/backend/tests/utils"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkUserCreationEquality(t *testing.T, expected *types.UserRegistrationInput, actual *types.UserCreationResponse) {
	t.Helper()

	assert.NotZero(t, actual.CreatedUserID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.NotEmpty(t, actual.TwoFactorSecret)
	assert.NotZero(t, actual.CreatedAt)
}

func checkUserEquality(t *testing.T, expected, actual *types.User) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.NotZero(t, actual.CreatedAt)
	assert.Nil(t, actual.LastUpdatedAt)
	assert.Nil(t, actual.ArchivedAt)
}

func (s *TestSuite) TestUsers_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create user.
			exampleUserInput := fakes.BuildFakeUserCreationInput()
			createdUser, err := testClients.user.CreateUser(ctx, exampleUserInput)
			requireNotNilAndNoProblems(t, createdUser, err)

			// Assert user equality.
			checkUserCreationEquality(t, exampleUserInput, createdUser)

			assert.NoError(t, testClients.admin.ArchiveUser(ctx, createdUser.CreatedUserID))
		}
	})

	s.runForEachClientExcept("should return 400 for duplicate user registration", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create user.
			exampleUserInput := fakes.BuildFakeUserCreationInput()
			createdUser, err := testClients.user.CreateUser(ctx, exampleUserInput)
			requireNotNilAndNoProblems(t, createdUser, err)

			// attempt to create user again.
			_, err = testClients.user.CreateUser(ctx, exampleUserInput)
			require.Error(t, err)

			// Assert user equality.
			checkUserCreationEquality(t, exampleUserInput, createdUser)

			assert.NoError(t, testClients.admin.ArchiveUser(ctx, createdUser.CreatedUserID))
		}
	})
}

func (s *TestSuite) TestUsers_Reading_Returns404ForNonexistentUser() {
	s.runForEachClientExcept("should return an error when trying to read a user that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.admin.GetUser(ctx, nonexistentID)
			assert.Nil(t, actual)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestUsers_Reading() {
	s.runForEachClientExcept("should be able to be read", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, _, _, _ := createUserAndClientForTest(ctx, t, nil)

			actual, err := testClients.admin.GetUser(ctx, user.ID)
			if err != nil {
				t.Logf("error encountered trying to fetch user %q: %v\n", user.Username, err)
			}
			requireNotNilAndNoProblems(t, actual, err)

			// Assert user equality.
			checkUserEquality(t, user, actual)

			// Clean up.
			assert.NoError(t, testClients.admin.ArchiveUser(ctx, actual.ID))
		}
	})
}

func (s *TestSuite) TestUsers_PermissionsChecking() {
	s.runForEachClientExcept("should be able to check users permissions", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, _, _, _ := createUserAndClientForTest(ctx, t, nil)

			permissions, err := testClients.user.CheckUserPermissions(ctx, authorization.ReadWebhooksPermission.ID())
			if err != nil {
				t.Logf("error encountered trying to fetch user %q: %v\n", user.Username, err)
			}
			requireNotNilAndNoProblems(t, permissions, err)

			for _, status := range permissions.Permissions {
				assert.True(t, status)
			}

			// Clean up.
			assert.NoError(t, testClients.admin.ArchiveUser(ctx, user.ID))
		}
	})
}

func (s *TestSuite) TestUsers_Searching_ReturnsEmptyWhenSearchingForUsernameThatIsNotPresentInTheDatabase() {
	s.runForEachClientExcept("it should return empty slice when searching for a username that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.admin.SearchForUsersByUsername(ctx, "   this is a really long string that contains characters unlikely to yield any real results   ")
			assert.Nil(t, actual)
			assert.NoError(t, err)
		}
	})
}

func (s *TestSuite) TestUsers_Searching_OnlyAccessibleToAdmins() {
	s.runForEachClientExcept("it should only be accessible to admins", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Search For user.
			actual, err := testClients.user.SearchForUsersByUsername(ctx, s.user.Username)
			assert.Nil(t, actual)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestUsers_Searching() {
	s.runForEachClientExcept("it should return be searchable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleUsername := fakes.BuildFakeUser().Username

			// create users
			createdUserIDs := []string{}
			for i := 0; i < 5; i++ {
				in := &types.UserRegistrationInput{
					EmailAddress: gofakeit.Email(),
					Username:     fmt.Sprintf("%s%d", exampleUsername, i),
					Password:     gofakeit.Password(true, true, true, true, false, 64),
				}
				user, err := testutils.CreateServiceUser(ctx, urlToUse, in)
				require.NoError(t, err)
				createdUserIDs = append(createdUserIDs, user.ID)
			}

			// execute search
			actual, err := testClients.admin.SearchForUsersByUsername(ctx, exampleUsername)
			assert.NoError(t, err)
			assert.NotEmpty(t, actual)

			// ensure results look how we expect them to look
			for _, result := range actual {
				assert.True(t, strings.HasPrefix(result.Username, exampleUsername))
			}

			// clean up
			for _, id := range createdUserIDs {
				require.NoError(t, testClients.admin.ArchiveUser(ctx, id))
			}
		}
	})
}

func (s *TestSuite) TestUsers_Archiving_Returns404ForNonexistentUser() {
	s.runForEachClientExcept("should fail to archive a non-existent user", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.admin.ArchiveUser(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestUsers_Archiving() {
	s.runForEachClientExcept("should be able to be archived", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create user.
			exampleUserInput := fakes.BuildFakeUserCreationInput()
			createdUser, err := testClients.user.CreateUser(ctx, exampleUserInput)
			assert.NoError(t, err)
			assert.NotNil(t, createdUser)

			if createdUser == nil || err != nil {
				t.Log("something has gone awry, user returned is nil")
				t.FailNow()
			}

			// Execute.
			assert.NoError(t, testClients.admin.ArchiveUser(ctx, createdUser.CreatedUserID))
		}
	})
}

func (s *TestSuite) TestUsers_AvatarManagement() {
	s.runForEachClientExcept("should be able to upload an avatar", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, avatar := testutils.BuildArbitraryImagePNGBytes(256)

			require.NoError(t, testClients.user.UploadNewAvatar(ctx, avatar, "png"))

			// Assert user equality.
			user, err := testClients.admin.GetUser(ctx, s.user.ID)
			requireNotNilAndNoProblems(t, user, err)

			assert.NotEmpty(t, user.AvatarSrc)

			assert.NoError(t, testClients.admin.ArchiveUser(ctx, s.user.ID))
		}
	}, pasetoAuthType)
}
