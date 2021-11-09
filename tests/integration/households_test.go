package integration

import (
	"github.com/brianvoe/gofakeit/v5"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkHouseholdEquality(t *testing.T, expected, actual *types.Household) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected BucketName for household %s to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestHouseholds_Creating() {
	s.runForEachClientExcept("should be possible to create households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, createdHousehold)

			// Clean up.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Listing() {
	s.runForEachClientExcept("should be possible to list households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create households.
			var expected []*types.Household
			for i := 0; i < 5; i++ {
				// Create household.
				exampleHousehold := fakes.BuildFakeHousehold()
				exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
				createdHousehold, householdCreationErr := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
				requireNotNilAndNoProblems(t, createdHousehold, householdCreationErr)

				expected = append(expected, createdHousehold)
			}

			// Assert household list equality.
			actual, err := testClients.main.GetHouseholds(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Households),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Households),
			)

			// Clean up.
			for _, createdHousehold := range actual.Households {
				assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
			}
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading_Returns404ForNonexistentHousehold() {
	s.runForEachClientExcept("should not be possible to read a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Attempt to fetch nonexistent household.
			_, err := testClients.main.GetHousehold(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading() {
	s.runForEachClientExcept("should be possible to read a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Fetch household.
			actual, err := testClients.main.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)

			// Clean up household.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Updating_Returns404ForNonexistentHousehold() {
	s.runForEachClientExcept("should not be possible to update a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHousehold.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateHousehold(ctx, exampleHousehold))
		}
	})
}

// convertHouseholdToHouseholdUpdateInput creates a householdUpdateInput struct from a household.
func convertHouseholdToHouseholdUpdateInput(x *types.Household) *types.HouseholdUpdateInput {
	return &types.HouseholdUpdateInput{
		Name:          x.Name,
		BelongsToUser: x.BelongsToUser,
	}
}

func (s *TestSuite) TestHouseholds_Updating() {
	s.runForEachClientExcept("should be possible to update a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Change household.
			createdHousehold.Update(convertHouseholdToHouseholdUpdateInput(exampleHousehold))
			assert.NoError(t, testClients.main.UpdateHousehold(ctx, createdHousehold))

			// Fetch household.
			actual, err := testClients.main.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up household.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Archiving_Returns404ForNonexistentHousehold() {
	s.runForEachClientExcept("should not be possible to archive a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveHousehold(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Archiving() {
	s.runForEachClientExcept("should be possible to archive a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(exampleHousehold)
			createdHousehold, err := testClients.main.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Clean up household.
			assert.NoError(t, testClients.main.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingPreExistentUser() {
	s.runForCookieClient("should be possible to change members of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.main.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhookID, err := testClients.main.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			n := <-notificationsChan
			assert.Equal(t, n.DataType, types.WebhookDataType)
			require.NotNil(t, n.Webhook)
			checkWebhookEquality(t, exampleWebhook, n.Webhook)

			createdWebhook, err := testClients.main.GetWebhook(ctx, createdWebhookID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			t.Logf("creating user to invite")
			u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

			t.Logf("inviting user")
			invitationID, err := testClients.main.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
				FromUser:             s.user.ID,
				Note:                 t.Name(),
				ToEmail:              u.EmailAddress,
				DestinationHousehold: relevantHouseholdID,
			})
			require.NoError(t, err)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.HouseholdInvitationDataType)

			t.Logf("accepting invitation")
			err = c.AcceptHouseholdInvitation(ctx, relevantHouseholdID, invitationID, t.Name())
			require.NoError(t, err)

			t.Logf("fetching households")
			households, err := c.GetHouseholds(ctx, nil)

			var found bool
			for _, household := range households.Households {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)
			require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingUserWhoSignsUpIndependently() {
	s.runForCookieClient("should be possible to change members of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.main.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhookID, err := testClients.main.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			n := <-notificationsChan
			assert.Equal(t, n.DataType, types.WebhookDataType)
			require.NotNil(t, n.Webhook)
			checkWebhookEquality(t, exampleWebhook, n.Webhook)

			createdWebhook, err := testClients.main.GetWebhook(ctx, createdWebhookID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			t.Logf("inviting user")
			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				FromUser:             s.user.ID,
				Note:                 t.Name(),
				ToEmail:              gofakeit.Email(),
				DestinationHousehold: relevantHouseholdID,
			}
			invitationID, err := testClients.main.InviteUserToHousehold(ctx, inviteReq)
			require.NoError(t, err)

			t.Logf("creating user to invite")
			_, _, c, _ := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, true, 64),
			})

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.HouseholdInvitationDataType)

			t.Logf("accepting invitation")
			err = c.AcceptHouseholdInvitation(ctx, relevantHouseholdID, invitationID, t.Name())
			require.NoError(t, err)

			t.Logf("fetching households")
			households, err := c.GetHouseholds(ctx, nil)

			var found bool
			for _, household := range households.Households {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)
			require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingNewUser() {
	// create inviting user

	// create user to invite

	// invite user to household

	// register new user with token

	// assert that the invited user is a member of the household
}

func (s *TestSuite) TestHouseholds_ChangingMemberships() {
	//s.runForCookieClient("should be possible to change members of a household", func(testClients *testClientWrapper) func() {
	//	return func() {
	//		t := s.T()
	//
	//		ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
	//		defer span.End()
	//
	//		const userCount = 1
	//
	//		currentStatus, statusErr := testClients.main.UserStatus(s.ctx)
	//		requireNotNilAndNoProblems(t, currentStatus, statusErr)
	//		t.Logf("initial household is %s; initial user ID is %s", currentStatus.ActiveHousehold, s.user.ID)
	//
	//		// fetch household data
	//		householdCreationInput := &types.HouseholdCreationRequestInput{
	//			Name: fakes.BuildFakeHousehold().Name,
	//		}
	//		household, householdCreationErr := testClients.main.CreateHousehold(ctx, householdCreationInput)
	//		require.NoError(t, householdCreationErr)
	//		require.NotNil(t, household)
	//
	//		t.Logf("created household %s", household.ID)
	//
	//		require.NoError(t, testClients.main.SwitchActiveHousehold(ctx, household.ID))
	//
	//		t.Logf("switched main test client active household to %s, creating webhook", household.ID)
	//
	//		stopChan := make(chan bool, 1)
	//		notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
	//		require.NotNil(t, notificationsChan)
	//		require.NoError(t, err)
	//
	//		// Create webhook.
	//		exampleWebhook := fakes.BuildFakeWebhook()
	//		exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
	//		createdWebhookID, err := testClients.main.CreateWebhook(ctx, exampleWebhookInput)
	//		require.NoError(t, err)
	//
	//		n := <-notificationsChan
	//		assert.Equal(t, n.DataType, types.WebhookDataType)
	//		require.NotNil(t, n.Webhook)
	//		checkWebhookEquality(t, exampleWebhook, n.Webhook)
	//
	//		createdWebhook, err := testClients.main.GetWebhook(ctx, createdWebhookID)
	//		requireNotNilAndNoProblems(t, createdWebhook, err)
	//		require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)
	//
	//		t.Logf("created webhook %s for household %s", createdWebhook.ID, createdWebhook.BelongsToHousehold)
	//
	//		// create dummy users
	//		users := []*types.User{}
	//		clients := []*httpclient.Client{}
	//
	//		// create users
	//		for i := 0; i < userCount; i++ {
	//			u, _, c, _ := createUserAndClientForTest(ctx, t)
	//			users = append(users, u)
	//			clients = append(clients, c)
	//
	//			currentStatus, statusErr = c.UserStatus(s.ctx)
	//			requireNotNilAndNoProblems(t, currentStatus, statusErr)
	//			t.Logf("created user user %q with household %s", u.ID, currentStatus.ActiveHousehold)
	//		}
	//
	//		// check that each user cannot see the unreachable webhook
	//		for i := 0; i < userCount; i++ {
	//			t.Logf("checking that user %q CANNOT see webhook %s belonging to household %s", users[i].ID, createdWebhook.ID, createdWebhook.BelongsToHousehold)
	//			webhook, err := clients[i].GetWebhook(ctx, createdWebhook.ID)
	//			require.Nil(t, webhook)
	//			require.Error(t, err)
	//		}
	//
	//		// add them to the household
	//		for i := 0; i < userCount; i++ {
	//			t.Logf("adding user %q to household %s", users[i].ID, household.ID)
	//			require.NoError(t, testClients.main.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
	//				ToEmail:              users[i].EmailAddress,
	//				DestinationHousehold: household.ID,
	//				Note:                 t.Name(),
	//			}))
	//
	//			n := <-notificationsChan
	//			assert.Equal(t, n.DataType, types.UserMembershipDataType)
	//			t.Logf("added user %q to household %s", users[i].ID, household.ID)
	//
	//			t.Logf("setting user %q's client to household %s", users[i].ID, household.ID)
	//			require.NoError(t, clients[i].SwitchActiveHousehold(ctx, household.ID))
	//
	//			currentStatus, statusErr = clients[i].UserStatus(s.ctx)
	//			requireNotNilAndNoProblems(t, currentStatus, statusErr)
	//			require.Equal(t, currentStatus.ActiveHousehold, household.ID)
	//			t.Logf("set user %q's current active household to %s", users[i].ID, household.ID)
	//		}
	//
	//		// grant all permissions
	//		for i := 0; i < userCount; i++ {
	//			input := &types.ModifyUserPermissionsInput{
	//				Reason:   t.Name(),
	//				NewRoles: []string{authorization.HouseholdAdminRole.String()},
	//			}
	//			require.NoError(t, testClients.main.ModifyMemberPermissions(ctx, household.ID, users[i].ID, input))
	//		}
	//
	//		// check that each user can see the webhook
	//		for i := 0; i < userCount; i++ {
	//			t.Logf("checking if user %q CAN now see webhook %s belonging to household %s", users[i].ID, createdWebhook.ID, createdWebhook.BelongsToHousehold)
	//			webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
	//			requireNotNilAndNoProblems(t, webhook, webhookRetrievalError)
	//		}
	//
	//		// remove users from household
	//		for i := 0; i < userCount; i++ {
	//			require.NoError(t, testClients.main.RemoveUserFromHousehold(ctx, household.ID, users[i].ID))
	//		}
	//
	//		// check that each user cannot see the webhook
	//		for i := 0; i < userCount; i++ {
	//			webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
	//			require.Nil(t, webhook)
	//			require.Error(t, webhookRetrievalError)
	//		}
	//
	//		// Clean up.
	//		require.NoError(t, testClients.main.ArchiveWebhook(ctx, createdWebhook.ID))
	//
	//		for i := 0; i < userCount; i++ {
	//			require.NoError(t, testClients.admin.ArchiveUser(ctx, users[i].ID))
	//		}
	//	}
	//})
	//
	//s.runForPASETOClient("should be possible to change members of a household", func(testClients *testClientWrapper) func() {
	//	return func() {
	//		t := s.T()
	//
	//		ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
	//		defer span.End()
	//
	//		const userCount = 1
	//
	//		currentStatus, statusErr := testClients.main.UserStatus(s.ctx)
	//		requireNotNilAndNoProblems(t, currentStatus, statusErr)
	//		t.Logf("initial household is %s; initial user ID is %s", currentStatus.ActiveHousehold, s.user.ID)
	//
	//		// fetch household data
	//		householdCreationInput := &types.HouseholdCreationRequestInput{
	//			Name: fakes.BuildFakeHousehold().Name,
	//		}
	//		household, householdCreationErr := testClients.main.CreateHousehold(ctx, householdCreationInput)
	//		require.NoError(t, householdCreationErr)
	//		require.NotNil(t, household)
	//
	//		t.Logf("created household %s", household.ID)
	//
	//		require.NoError(t, testClients.main.SwitchActiveHousehold(ctx, household.ID))
	//
	//		t.Logf("switched main test client active household to %s, creating webhook", household.ID)
	//
	//		// Create webhook.
	//		exampleWebhook := fakes.BuildFakeWebhook()
	//		exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
	//		createdWebhookID, err := testClients.main.CreateWebhook(ctx, exampleWebhookInput)
	//		require.NoError(t, err)
	//
	//		var createdWebhook *types.Webhook
	//		checkFunc := func() bool {
	//			createdWebhook, err = testClients.main.GetWebhook(ctx, createdWebhookID)
	//			return assert.NotNil(t, createdWebhook) && assert.NoError(t, err)
	//		}
	//		assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
	//
	//		// assert webhook equality
	//		requireNotNilAndNoProblems(t, createdWebhook, err)
	//		require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)
	//
	//		t.Logf("created webhook %s for household %s", createdWebhook.ID, createdWebhook.BelongsToHousehold)
	//
	//		// create dummy users
	//		users := []*types.User{}
	//		clients := []*httpclient.Client{}
	//
	//		// create users
	//		for i := 0; i < userCount; i++ {
	//			u, _, c, _ := createUserAndClientForTest(ctx, t)
	//			users = append(users, u)
	//			clients = append(clients, c)
	//
	//			currentStatus, statusErr = c.UserStatus(s.ctx)
	//			requireNotNilAndNoProblems(t, currentStatus, statusErr)
	//			t.Logf("created user user %q with household %s", u.ID, currentStatus.ActiveHousehold)
	//		}
	//
	//		// check that each user cannot see the unreachable webhook
	//		for i := 0; i < userCount; i++ {
	//			t.Logf("checking that user %q CANNOT see webhook %s belonging to household %s", users[i].ID, createdWebhook.ID, createdWebhook.BelongsToHousehold)
	//			webhook, err := clients[i].GetWebhook(ctx, createdWebhook.ID)
	//			require.Nil(t, webhook)
	//			require.Error(t, err)
	//		}
	//
	//		// add them to the household
	//		for i := 0; i < userCount; i++ {
	//			t.Logf("adding user %q to household %s", users[i].ID, household.ID)
	//			require.NoError(t, testClients.main.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
	//				ToEmail:              users[i].EmailAddress,
	//				DestinationHousehold: household.ID,
	//				Note:                 t.Name(),
	//			}))
	//			t.Logf("added user %q to household %s", users[i].ID, household.ID)
	//
	//			t.Logf("setting user %q's client to household %s", users[i].ID, household.ID)
	//			checkFunc = func() bool {
	//				return assert.NoError(t, clients[i].SwitchActiveHousehold(ctx, household.ID))
	//			}
	//			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
	//
	//			currentStatus, statusErr = clients[i].UserStatus(s.ctx)
	//			requireNotNilAndNoProblems(t, currentStatus, statusErr)
	//			require.Equal(t, currentStatus.ActiveHousehold, household.ID)
	//			t.Logf("set user %q's current active household to %s", users[i].ID, household.ID)
	//		}
	//
	//		// grant all permissions
	//		for i := 0; i < userCount; i++ {
	//			input := &types.ModifyUserPermissionsInput{
	//				Reason:   t.Name(),
	//				NewRoles: []string{authorization.HouseholdAdminRole.String()},
	//			}
	//			require.NoError(t, testClients.main.ModifyMemberPermissions(ctx, household.ID, users[i].ID, input))
	//		}
	//
	//		// check that each user can see the webhook
	//		for i := 0; i < userCount; i++ {
	//			t.Logf("checking if user %q CAN now see webhook %s belonging to household %s", users[i].ID, createdWebhook.ID, createdWebhook.BelongsToHousehold)
	//			webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
	//			requireNotNilAndNoProblems(t, webhook, webhookRetrievalError)
	//		}
	//
	//		// remove users from household
	//		for i := 0; i < userCount; i++ {
	//			require.NoError(t, testClients.main.RemoveUserFromHousehold(ctx, household.ID, users[i].ID))
	//		}
	//
	//		// check that each user cannot see the webhook
	//		for i := 0; i < userCount; i++ {
	//			webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
	//			require.Nil(t, webhook)
	//			require.Error(t, webhookRetrievalError)
	//		}
	//
	//		// Clean up.
	//		require.NoError(t, testClients.main.ArchiveWebhook(ctx, createdWebhook.ID))
	//
	//		for i := 0; i < userCount; i++ {
	//			require.NoError(t, testClients.admin.ArchiveUser(ctx, users[i].ID))
	//		}
	//	}
	//})
}

func (s *TestSuite) TestHouseholds_OwnershipTransfer() {
	s.runForCookieClient("should be possible to transfer ownership of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create users
			futureOwner, _, _, futureOwnerClient := createUserAndClientForTest(ctx, t, nil)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationRequestInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.main.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			t.Logf("created household %s", household.ID)

			require.NoError(t, testClients.main.SwitchActiveHousehold(ctx, household.ID))

			t.Logf("switched to active household: %s", household.ID)

			// create a webhook
			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhookID, err := testClients.main.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			n := <-notificationsChan
			assert.Equal(t, n.DataType, types.WebhookDataType)
			require.NotNil(t, n.Webhook)
			checkWebhookEquality(t, exampleWebhook, n.Webhook)

			createdWebhook, err := testClients.main.GetWebhook(ctx, createdWebhookID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			t.Logf("created webhook %s belonging to household %s", createdWebhook.ID, createdWebhook.BelongsToHousehold)
			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			// check that user cannot see the webhook
			webhook, err := futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// add them to the household
			require.NoError(t, testClients.main.TransferHouseholdOwnership(ctx, household.ID, &types.HouseholdOwnershipTransferInput{
				Reason:       t.Name(),
				CurrentOwner: household.BelongsToUser,
				NewOwner:     futureOwner.ID,
			}))

			t.Logf("transferred household %s from user %s to user %s", household.ID, household.BelongsToUser, futureOwner.ID)

			require.NoError(t, futureOwnerClient.SwitchActiveHousehold(ctx, household.ID))

			// check that user can see the webhook
			webhook, err = futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)

			// check that old user cannot see the webhook
			webhook, err = testClients.main.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// check that new owner can delete the webhook
			require.NoError(t, futureOwnerClient.ArchiveWebhook(ctx, createdWebhook.ID))

			// Clean up.
			require.Error(t, testClients.main.ArchiveWebhook(ctx, createdWebhook.ID))
			require.NoError(t, testClients.admin.ArchiveUser(ctx, futureOwner.ID))
		}
	})

	s.runForPASETOClient("should be possible to transfer ownership of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create users
			futureOwner, _, _, futureOwnerClient := createUserAndClientForTest(ctx, t, nil)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationRequestInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.main.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			t.Logf("created household %s", household.ID)

			require.NoError(t, testClients.main.SwitchActiveHousehold(ctx, household.ID))

			t.Logf("switched to active household: %s", household.ID)

			// create a webhook
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			createdWebhookID, err := testClients.main.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			var createdWebhook *types.Webhook
			checkFunc := func() bool {
				createdWebhook, err = testClients.main.GetWebhook(ctx, createdWebhookID)
				return assert.NotNil(t, createdWebhook) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, createdWebhook, err)

			t.Logf("created webhook %s belonging to household %s", createdWebhook.ID, createdWebhook.BelongsToHousehold)
			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			// check that user cannot see the webhook
			webhook, err := futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// add them to the household
			require.NoError(t, testClients.main.TransferHouseholdOwnership(ctx, household.ID, &types.HouseholdOwnershipTransferInput{
				Reason:       t.Name(),
				CurrentOwner: household.BelongsToUser,
				NewOwner:     futureOwner.ID,
			}))

			t.Logf("transferred household %s from user %s to user %s", household.ID, household.BelongsToUser, futureOwner.ID)

			require.NoError(t, futureOwnerClient.SwitchActiveHousehold(ctx, household.ID))

			// check that user can see the webhook
			webhook, err = futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)

			// check that old user cannot see the webhook
			webhook, err = testClients.main.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// check that new owner can delete the webhook
			require.NoError(t, futureOwnerClient.ArchiveWebhook(ctx, createdWebhook.ID))

			// Clean up.
			require.Error(t, testClients.main.ArchiveWebhook(ctx, createdWebhook.ID))
			require.NoError(t, testClients.admin.ArchiveUser(ctx, futureOwner.ID))
		}
	})
}
