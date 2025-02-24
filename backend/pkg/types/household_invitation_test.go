package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHouseholdInvitationCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &HouseholdInvitationCreationRequestInput{
			ToEmail: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestHouseholdInvitationUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &HouseholdInvitationUpdateRequestInput{
			Token: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
