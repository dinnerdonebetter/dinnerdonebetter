package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestHouseholdInstrumentOwnership_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &HouseholdInstrumentOwnership{}
		input := &HouseholdInstrumentOwnershipUpdateRequestInput{}

		fake.Struct(&input)

		x.Update(input)
	})
}

func TestHouseholdInstrumentOwnershipCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
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

		ctx := context.Background()
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

		ctx := context.Background()
		x := &HouseholdInstrumentOwnershipUpdateRequestInput{
			Quantity:          pointers.Pointer[uint16](1),
			ValidInstrumentID: pointers.Pointer(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
