package types

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"

	"github.com/stretchr/testify/assert"
)

func TestTransferHouseholdOwnershipInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &HouseholdOwnershipTransferInput{
			CurrentOwner: "123",
			NewOwner:     "321",
			Reason:       t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestModifyUserPermissionsInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &ModifyUserPermissionsInput{
			NewRole: authorization.HouseholdMemberRole.String(),
			Reason:  t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
