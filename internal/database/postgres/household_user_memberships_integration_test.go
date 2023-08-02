package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_Integration_HouseholdUserMemberships(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleUser := createUserForTest(t, ctx, nil, dbc)
	households, err := dbc.getHouseholdsForUser(ctx, dbc.db, exampleUser.ID, false, nil)
	assert.NoError(t, err)
	assert.Len(t, households.Data, 1)

	household := fakes.BuildFakeHousehold()
	household.BelongsToUser = exampleUser.ID
	exampleHousehold := createHouseholdForTest(t, ctx, household, dbc)

	memberUserIDs := []string{exampleUser.ID}

	for i := 0; i < exampleQuantity; i++ {
		newMember := createUserForTest(t, ctx, nil, dbc)
		assert.NoError(t, dbc.addUserToHousehold(ctx, dbc.db, &types.HouseholdUserMembershipDatabaseCreationInput{
			ID:            identifiers.New(),
			Reason:        "testing",
			UserID:        newMember.ID,
			HouseholdID:   exampleHousehold.ID,
			HouseholdRole: "household_member",
		}))
		memberUserIDs = append(memberUserIDs, newMember.ID)
	}

	household, err = dbc.GetHousehold(ctx, exampleHousehold.ID)
	assert.NoError(t, err)

	householdMemberUserIDs := []string{}
	for _, member := range household.Members {
		householdMemberUserIDs = append(householdMemberUserIDs, member.BelongsToUser.ID)
	}

	assert.Subset(t, memberUserIDs, householdMemberUserIDs)

	isMember, err := dbc.UserIsMemberOfHousehold(ctx, memberUserIDs[0], exampleHousehold.ID)
	assert.NoError(t, err)
	assert.True(t, isMember)

	assert.NoError(t, dbc.MarkHouseholdAsUserDefault(ctx, memberUserIDs[1], exampleHousehold.ID))
	defaultHouseholdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, memberUserIDs[1])
	assert.NoError(t, err)
	assert.Equal(t, exampleHousehold.ID, defaultHouseholdID)

	sessionCtxData, err := dbc.BuildSessionContextDataForUser(ctx, memberUserIDs[1])
	assert.NoError(t, err)
	assert.NotNil(t, sessionCtxData)
	assert.Equal(t, exampleHousehold.ID, sessionCtxData.ActiveHouseholdID)

	assert.NoError(t, dbc.RemoveUserFromHousehold(ctx, memberUserIDs[len(memberUserIDs)-1], exampleHousehold.ID))

	assert.NoError(t, dbc.TransferHouseholdOwnership(ctx, exampleHousehold.ID, &types.HouseholdOwnershipTransferInput{
		Reason:       "testing",
		CurrentOwner: exampleUser.ID,
		NewOwner:     memberUserIDs[1],
	}))

	assert.NoError(t, dbc.ModifyUserPermissions(ctx, memberUserIDs[0], exampleHousehold.ID, &types.ModifyUserPermissionsInput{
		Reason:  "testing",
		NewRole: "household_admin",
	}))
}
