package types

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidPreparationInstrumentCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrumentCreationRequestInput{
			Notes:              fake.LoremIpsumSentence(exampleQuantity),
			ValidPreparationID: fake.LoremIpsumSentence(exampleQuantity),
			ValidInstrumentID:  fake.LoremIpsumSentence(exampleQuantity),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrumentCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationInstrumentUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrumentUpdateRequestInput{
			Notes:              pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			ValidPreparationID: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			ValidInstrumentID:  pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrumentUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationInstrumentCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationInstrumentCreationRequestInput{
			ValidPreparationID: t.Name(),
			ValidInstrumentID:  t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationInstrumentDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationInstrumentDatabaseCreationInput{
			ID:                 t.Name(),
			ValidPreparationID: t.Name(),
			ValidInstrumentID:  t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationInstrumentUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationInstrumentUpdateRequestInput{
			ValidPreparationID: pointers.String(t.Name()),
			ValidInstrumentID:  pointers.String(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrument{}

		x.Update(&ValidPreparationInstrumentUpdateRequestInput{
			Notes:              pointers.String(t.Name()),
			ValidPreparationID: pointers.String(t.Name()),
			ValidInstrumentID:  pointers.String(t.Name()),
		})
	})
}
