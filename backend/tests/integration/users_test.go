package integration

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	testutils2 "github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkUserCreationEquality(t *testing.T, expected *types.UserRegistrationInput, actual *types.UserCreationResponse) {
	t.Helper()

	assert.NotZero(t, actual.CreatedUserID)
	assert.Equal(t, expected.FirstName, actual.FirstName)
	assert.Equal(t, expected.LastName, actual.LastName)
	assert.Equal(t, expected.Username, actual.Username)
	assert.NotEmpty(t, actual.TwoFactorSecret)
	assert.NotZero(t, actual.CreatedAt)
}

func checkUserEquality(t *testing.T, expected, actual *types.User) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.NotZero(t, actual.CreatedAt)
	assert.Nil(t, actual.ArchivedAt)
}

func (s *TestSuite) TestUsers_Creating() {
	s.runTest("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create userClient.
			exampleUserInput := fakes.BuildFakeUserCreationInput()
			createdUser, err := testClients.userClient.CreateUser(ctx, exampleUserInput)
			requireNotNilAndNoProblems(t, createdUser, err)

			// Assert userClient equality.
			checkUserCreationEquality(t, exampleUserInput, createdUser)

			assert.NoError(t, testClients.adminClient.ArchiveUser(ctx, createdUser.CreatedUserID))
		}
	})

	s.runTest("should return 400 for duplicate userClient registration", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create userClient.
			exampleUserInput := fakes.BuildFakeUserCreationInput()
			createdUser, err := testClients.userClient.CreateUser(ctx, exampleUserInput)
			requireNotNilAndNoProblems(t, createdUser, err)

			// attempt to create userClient again.
			_, err = testClients.userClient.CreateUser(ctx, exampleUserInput)
			require.Error(t, err)

			// Assert userClient equality.
			checkUserCreationEquality(t, exampleUserInput, createdUser)

			assert.NoError(t, testClients.adminClient.ArchiveUser(ctx, createdUser.CreatedUserID))
		}
	})
}

func (s *TestSuite) TestUsers_Reading() {
	s.runTest("should be able to be read", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, _ := createUserAndClientForTest(ctx, t, nil)

			actual, err := testClients.adminClient.GetUser(ctx, user.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert userClient equality.
			checkUserEquality(t, user, actual)

			// Clean up.
			assert.NoError(t, testClients.adminClient.ArchiveUser(ctx, actual.ID))
		}
	})

	s.runTest("should return an error when trying to read a userClient that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.adminClient.GetUser(ctx, nonexistentID)
			assert.Nil(t, actual)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestUsers_PermissionsChecking() {
	s.runTest("should be able to check users permissions", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, _ := createUserAndClientForTest(ctx, t, nil)

			input := &types.UserPermissionsRequestInput{Permissions: []string{authorization.ReadWebhooksPermission.ID()}}

			permissions, err := testClients.userClient.CheckPermissions(ctx, input)
			requireNotNilAndNoProblems(t, permissions, err)

			for _, status := range permissions.Permissions {
				assert.True(t, status)
			}

			// Clean up.
			assert.NoError(t, testClients.adminClient.ArchiveUser(ctx, user.ID))
		}
	})
}

func (s *TestSuite) TestUsers_Deleting() {
	s.runTest("should be able to delete your data", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			require.NoError(t, testClients.userClient.DestroyAllUserData(ctx))

			// Clean up.
			u, err := testClients.adminClient.GetUser(ctx, testClients.user.ID)
			assert.Nil(t, u)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestUsers_Searching_OnlyAccessibleToAdmins() {
	s.runTest("it should only be accessible to admins", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Search For userClient.
			actual, err := testClients.userClient.SearchForUsers(ctx, s.user.Username, nil)
			assert.Nil(t, actual)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestUsers_Searching() {
	s.runTest("it should return be searchable", func(testClients *testClientWrapper) func() {
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
				user, err := testutils2.CreateServiceUser(ctx, urlToUse, in)
				require.NoError(t, err)
				createdUserIDs = append(createdUserIDs, user.ID)
			}

			// execute search
			actual, err := testClients.adminClient.SearchForUsers(ctx, exampleUsername, nil)
			assert.NoError(t, err)
			assert.NotEmpty(t, actual)

			// ensure results look how we expect them to look
			for _, result := range actual.Data {
				assert.True(t, strings.HasPrefix(result.Username, exampleUsername))
			}

			// clean up
			for _, id := range createdUserIDs {
				require.NoError(t, testClients.adminClient.ArchiveUser(ctx, id))
			}
		}
	})
}

func (s *TestSuite) TestUsers_Archiving_Returns404ForNonexistentUser() {
	s.runTest("should error when archiving a non-existent userClient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.adminClient.ArchiveUser(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestUsers_Archiving() {
	s.runTest("should be able to be archived", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create userClient.
			exampleUserInput := fakes.BuildFakeUserCreationInput()
			createdUser, err := testClients.userClient.CreateUser(ctx, exampleUserInput)
			assert.NoError(t, err)
			assert.NotNil(t, createdUser)

			if createdUser == nil || err != nil {
				t.FailNow()
			}

			// Execute.
			assert.NoError(t, testClients.adminClient.ArchiveUser(ctx, createdUser.CreatedUserID))
		}
	})
}

func (s *TestSuite) TestUsers_AvatarManagement() {
	s.runTest("should be able to upload an avatar", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, avatar := testutils2.BuildArbitraryImagePNGBytes(256)

			encoded := base64.RawStdEncoding.EncodeToString(avatar)

			_, err := testClients.userClient.UploadUserAvatar(ctx, &types.AvatarUpdateInput{Base64EncodedData: encoded})
			require.NoError(t, err)

			// Assert userClient equality.
			user, err := testClients.adminClient.GetUser(ctx, s.user.ID)
			requireNotNilAndNoProblems(t, user, err)

			assert.NotEmpty(t, user.AvatarSrc)

			assert.NoError(t, testClients.adminClient.ArchiveUser(ctx, s.user.ID))
		}
	})
}
