package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestQuerier_Migrate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleUser.CreatedOn = exampleCreationTime

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)

		exampleTestUserConfig := &types.TestUserCreationConfig{
			Username:       exampleUser.Username,
			Password:       exampleUser.HashedPassword,
			HashedPassword: exampleUser.HashedPassword,
			IsServiceAdmin: true,
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() uint64 {
			return exampleCreationTime
		}

		// called by c.IsReady()
		db.ExpectPing()

		c.migrateOnce.Do(func() {})

		// expect TestUser to be queried for
		testUserExistenceArgs := []interface{}{exampleTestUserConfig.Username}
		db.ExpectQuery(formatQueryForSQLMock(testUserExistenceQuery)).
			WithArgs(interfaceToDriverValue(testUserExistenceArgs)...).
			WillReturnError(sql.ErrNoRows)

		db.ExpectBegin()

		// expect TestUser to be created
		testUserCreationArgs := []interface{}{
			&idMatcher{},
			exampleTestUserConfig.Username,
			exampleTestUserConfig.HashedPassword,
			defaultTestUserTwoFactorSecret,
			types.GoodStandingHouseholdStatus,
			authorization.ServiceAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(testUserCreationQuery)).
			WithArgs(interfaceToDriverValue(testUserCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleUser.ID))

		// create household for created TestUser
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationArgs := []interface{}{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactEmail,
			householdCreationInput.ContactPhone,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHousehold.ID))

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []interface{}{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			true,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHousehold.ID))

		db.ExpectCommit()

		err := c.Migrate(ctx, 1, exampleTestUserConfig)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with failure executing creation query", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleUser.CreatedOn = exampleCreationTime

		exampleTestUserConfig := &types.TestUserCreationConfig{
			Username:       exampleUser.Username,
			Password:       exampleUser.HashedPassword,
			HashedPassword: exampleUser.HashedPassword,
			IsServiceAdmin: true,
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() uint64 {
			return exampleCreationTime
		}

		// called by c.IsReady()
		db.ExpectPing()

		c.migrateOnce.Do(func() {})

		// expect TestUser to be queried for
		testUserExistenceArgs := []interface{}{exampleTestUserConfig.Username}
		db.ExpectQuery(formatQueryForSQLMock(testUserExistenceQuery)).
			WithArgs(interfaceToDriverValue(testUserExistenceArgs)...).
			WillReturnError(sql.ErrNoRows)

		db.ExpectBegin()

		// expect TestUser to be created
		testUserCreationArgs := []interface{}{
			&idMatcher{},
			exampleTestUserConfig.Username,
			exampleTestUserConfig.HashedPassword,
			defaultTestUserTwoFactorSecret,
			types.GoodStandingHouseholdStatus,
			authorization.ServiceAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(testUserCreationQuery)).
			WithArgs(interfaceToDriverValue(testUserCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		err := c.Migrate(ctx, 1, exampleTestUserConfig)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}
