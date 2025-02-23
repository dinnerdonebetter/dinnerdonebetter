package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAccountStatusUpdateInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &UserAccountStatusUpdateInput{
			NewStatus:    string(GoodStandingUserAccountStatus),
			Reason:       t.Name(),
			TargetUserID: "123",
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
