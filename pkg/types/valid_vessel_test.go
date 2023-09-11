package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidVessel_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidVessel{
			CapacityUnit: &ValidMeasurementUnit{},
		}
		input := &ValidVesselUpdateRequestInput{}

		fake.Struct(&input)
		input.CapacityUnitID = pointers.Pointer(t.Name())
		input.UsableForStorage = pointers.Pointer(true)
		input.DisplayInSummaryLists = pointers.Pointer(true)
		input.IncludeInGeneratedInstructions = pointers.Pointer(true)

		x.Update(input)
	})
}

func TestValidVesselCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselCreationRequestInput{
			Name:           t.Name(),
			Description:    t.Name(),
			IconPath:       t.Name(),
			Capacity:       exampleQuantity,
			CapacityUnitID: pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidVesselDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselDatabaseCreationInput{
			ID:             t.Name(),
			Name:           t.Name(),
			CapacityUnitID: pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidVesselUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselUpdateRequestInput{
			Name:        pointers.Pointer(t.Name()),
			Description: pointers.Pointer(t.Name()),
			IconPath:    pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
