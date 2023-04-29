package email

import (
	"testing"

	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuildGeneratedPasswordResetTokenEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		token := fakes.BuildFakePasswordResetToken()

		actual, err := BuildGeneratedPasswordResetTokenEmail(t.Name(), token, envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestBuildInviteMemberEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		invitation := fakes.BuildFakeHouseholdInvitation()

		actual, err := BuildInviteMemberEmail(invitation, envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestBuildMealPlanCreatedEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		mealPlan := fakes.BuildFakeMealPlan()

		actual, err := BuildMealPlanCreatedEmail(t.Name(), mealPlan, envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Contains(t, actual.HTMLContent, logoURL)
	})
}

func TestBuildPasswordResetTokenRedeemedEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual, err := BuildPasswordResetTokenRedeemedEmail(t.Name(), envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestBuildUsernameReminderEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		user := fakes.BuildFakeUser()

		actual, err := BuildUsernameReminderEmail(t.Name(), user.Username, envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}
