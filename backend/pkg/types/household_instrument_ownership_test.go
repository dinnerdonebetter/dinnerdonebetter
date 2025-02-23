package types

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestHouseholdInstrumentOwnership_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &HouseholdInstrumentOwnership{}
		input := &HouseholdInstrumentOwnershipUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))

		x.Update(input)
	})
}

func TestHouseholdInstrumentOwnershipCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &HouseholdInstrumentOwnershipCreationRequestInput{
			Quantity:          1,
			ValidInstrumentID: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestHouseholdInstrumentOwnershipDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &HouseholdInstrumentOwnershipDatabaseCreationInput{
			ID:                t.Name(),
			Quantity:          1,
			ValidInstrumentID: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestHouseholdInstrumentOwnershipUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &HouseholdInstrumentOwnershipUpdateRequestInput{
			Quantity:          pointer.To[uint16](1),
			ValidInstrumentID: pointer.To(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
