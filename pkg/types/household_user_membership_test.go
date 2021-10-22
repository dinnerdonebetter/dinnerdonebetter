package types

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
)

func TestAddUserToHouseholdInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &AddUserToHouseholdInput{
			UserID: "123",
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestTransferHouseholdOwnershipInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
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

		ctx := context.Background()
		x := &ModifyUserPermissionsInput{
			NewRoles: []string{authorization.HouseholdMemberRole.String()},
			Reason:   t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
