package email

import (
	"testing"
	"time"

	identityfakes "github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	mealplanningfakes "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	"github.com/dinnerdonebetter/backend/internal/platform/email"

	"github.com/stretchr/testify/assert"
)

func TestBuildMealPlanCreatedEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		user := identityfakes.BuildFakeUser()
		user.EmailAddressVerifiedAt = new(time.Now())
		mealPlan := mealplanningfakes.BuildFakeMealPlan()

		actual, err := BuildMealPlanCreatedEmail(user, mealPlan, &email.EnvironmentConfig{})
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Contains(t, actual.HTMLContent, logoURL)
	})
}
