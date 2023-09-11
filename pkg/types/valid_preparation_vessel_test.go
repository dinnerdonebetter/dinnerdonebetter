package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	"github.com/stretchr/testify/assert"
)

func TestValidPreparationVesselCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationVesselCreationRequestInput{
			Notes:              t.Name(),
			ValidPreparationID: t.Name(),
			ValidVesselID:      t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationVesselCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationVesselUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationVesselUpdateRequestInput{
			Notes:              pointers.Pointer(t.Name()),
			ValidPreparationID: pointers.Pointer(t.Name()),
			ValidVesselID:      pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationVesselUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationVesselCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationVesselCreationRequestInput{
			ValidPreparationID: t.Name(),
			ValidVesselID:      t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationVesselDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationVesselDatabaseCreationInput{
			ID:                 t.Name(),
			ValidPreparationID: t.Name(),
			ValidVesselID:      t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationVesselUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationVesselUpdateRequestInput{
			ValidPreparationID: pointers.Pointer(t.Name()),
			ValidVesselID:      pointers.Pointer(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationVessel_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationVessel{}

		x.Update(&ValidPreparationVesselUpdateRequestInput{
			Notes:              pointers.Pointer(t.Name()),
			ValidPreparationID: pointers.Pointer(t.Name()),
			ValidVesselID:      pointers.Pointer(t.Name()),
		})
	})
}
