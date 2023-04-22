package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/pkg/pointers"
	"github.com/prixfixeco/backend/pkg/types"
)

func TestConvert(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		original := types.ValidPreparation{
			Name:           "Name",
			Description:    "Description",
			OnlyForVessels: true,
		}
		updateInput := types.ValidPreparationUpdateRequestInput{
			Slug:      pointers.Pointer("Slug"),
			PastTense: pointers.Pointer("PastTense"),
		}
		updated := types.ValidPreparation{
			Name:           "Name",
			Description:    "Description",
			Slug:           "Slug",
			OnlyForVessels: true,
		}

		require.NoError(t, Merge(&original, &updateInput))

		assert.Equal(t, updated, original)
	})
}
