package types

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/pointers"

	"github.com/stretchr/testify/assert"
)

func TestHousehold_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Household{}
		name := t.Name()

		x.Update(&HouseholdUpdateRequestInput{
			Name:          pointers.String(name),
			ContactEmail:  pointers.String(name),
			ContactPhone:  pointers.String(name),
			TimeZone:      pointers.String(name),
			BelongsToUser: name,
		})
	})
}

func TestHouseholdCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &HouseholdCreationRequestInput{
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestHouseholdUpdateInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		name := t.Name()

		x := &HouseholdUpdateRequestInput{
			Name: &name,
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestHouseholdCreationInputForNewUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, HouseholdCreationInputForNewUser(&User{}))
	})
}
