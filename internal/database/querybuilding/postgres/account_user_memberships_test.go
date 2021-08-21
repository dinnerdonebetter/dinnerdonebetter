package postgres

import (
	"context"
	"strings"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestPostgres_BuildGetDefaultHouseholdIDForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()

		expectedQuery := "SELECT households.id FROM households JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id WHERE household_user_memberships.belongs_to_user = $1 AND household_user_memberships.default_household = $2"
		expectedArgs := []interface{}{
			exampleUser.ID,
			true,
		}
		actualQuery, actualArgs := q.BuildGetDefaultHouseholdIDForUserQuery(ctx, exampleUser.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildUserIsMemberOfHouseholdQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		expectedQuery := "SELECT EXISTS ( SELECT household_user_memberships.id FROM household_user_memberships WHERE household_user_memberships.archived_on IS NULL AND household_user_memberships.belongs_to_household = $1 AND household_user_memberships.belongs_to_user = $2 )"
		expectedArgs := []interface{}{
			exampleHousehold.ID,
			exampleUser.ID,
		}
		actualQuery, actualArgs := q.BuildUserIsMemberOfHouseholdQuery(ctx, exampleUser.ID, exampleHousehold.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildAddUserToHouseholdQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := &types.AddUserToHouseholdInput{
			UserID:         exampleUser.ID,
			HouseholdID:    exampleHousehold.ID,
			Reason:         t.Name(),
			HouseholdRoles: []string{authorization.HouseholdMemberRole.String()},
		}

		expectedQuery := "INSERT INTO household_user_memberships (belongs_to_user,belongs_to_household,household_roles) VALUES ($1,$2,$3)"
		expectedArgs := []interface{}{
			exampleInput.UserID,
			exampleHousehold.ID,
			strings.Join(exampleInput.HouseholdRoles, householdMemberRolesSeparator),
		}
		actualQuery, actualArgs := q.BuildAddUserToHouseholdQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildRemoveUserFromHouseholdQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		expectedQuery := "DELETE FROM household_user_memberships WHERE household_user_memberships.archived_on IS NULL AND household_user_memberships.belongs_to_household = $1 AND household_user_memberships.belongs_to_user = $2"
		expectedArgs := []interface{}{
			exampleHousehold.ID,
			exampleUser.ID,
		}
		actualQuery, actualArgs := q.BuildRemoveUserFromHouseholdQuery(ctx, exampleUser.ID, exampleHousehold.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveHouseholdMembershipsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()

		expectedQuery := "UPDATE household_user_memberships SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1"
		expectedArgs := []interface{}{
			exampleUser.ID,
		}
		actualQuery, actualArgs := q.BuildArchiveHouseholdMembershipsForUserQuery(ctx, exampleUser.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateMembershipForNewUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		expectedQuery := "INSERT INTO household_user_memberships (belongs_to_user,belongs_to_household,default_household,household_roles) VALUES ($1,$2,$3,$4)"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleHousehold.ID,
			true,
			authorization.HouseholdAdminRole.String(),
		}
		actualQuery, actualArgs := q.BuildCreateMembershipForNewUserQuery(ctx, exampleUser.ID, exampleHousehold.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetHouseholdMembershipsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()

		expectedQuery := "SELECT household_user_memberships.id, household_user_memberships.belongs_to_user, household_user_memberships.belongs_to_household, household_user_memberships.household_roles, household_user_memberships.default_household, household_user_memberships.created_on, household_user_memberships.last_updated_on, household_user_memberships.archived_on FROM household_user_memberships JOIN households ON households.id = household_user_memberships.belongs_to_household WHERE household_user_memberships.archived_on IS NULL AND household_user_memberships.belongs_to_user = $1"
		expectedArgs := []interface{}{
			exampleUser.ID,
		}
		actualQuery, actualArgs := q.BuildGetHouseholdMembershipsForUserQuery(ctx, exampleUser.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildMarkHouseholdAsUserDefaultQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		expectedQuery := "UPDATE household_user_memberships SET default_household = (belongs_to_user = $1 AND belongs_to_household = $2) WHERE archived_on IS NULL AND belongs_to_user = $3"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleHousehold.ID,
			exampleUser.ID,
		}
		actualQuery, actualArgs := q.BuildMarkHouseholdAsUserDefaultQuery(ctx, exampleUser.ID, exampleHousehold.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildModifyUserPermissionsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleRoles := []string{authorization.HouseholdMemberRole.String()}
		exampleHousehold := fakes.BuildFakeHousehold()

		expectedQuery := "UPDATE household_user_memberships SET household_roles = $1 WHERE belongs_to_household = $2 AND belongs_to_user = $3"
		expectedArgs := []interface{}{
			strings.Join(exampleRoles, householdMemberRolesSeparator),
			exampleHousehold.ID,
			exampleUser.ID,
		}
		actualQuery, actualArgs := q.BuildModifyUserPermissionsQuery(ctx, exampleUser.ID, exampleHousehold.ID, exampleRoles)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildTransferHouseholdOwnershipQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleOldOwner := fakes.BuildFakeUser()
		exampleNewOwner := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		expectedQuery := "UPDATE households SET belongs_to_user = $1 WHERE archived_on IS NULL AND belongs_to_user = $2 AND id = $3"
		expectedArgs := []interface{}{
			exampleNewOwner.ID,
			exampleOldOwner.ID,
			exampleHousehold.ID,
		}
		actualQuery, actualArgs := q.BuildTransferHouseholdOwnershipQuery(ctx, exampleOldOwner.ID, exampleNewOwner.ID, exampleHousehold.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildTransferHouseholdMembershipsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleOldOwner := fakes.BuildFakeUser()
		exampleNewOwner := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		expectedQuery := "UPDATE household_user_memberships SET belongs_to_user = $1 WHERE archived_on IS NULL AND belongs_to_household = $2 AND belongs_to_user = $3"
		expectedArgs := []interface{}{
			exampleNewOwner.ID,
			exampleHousehold.ID,
			exampleOldOwner.ID,
		}
		actualQuery, actualArgs := q.BuildTransferHouseholdMembershipsQuery(ctx, exampleOldOwner.ID, exampleNewOwner.ID, exampleHousehold.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
