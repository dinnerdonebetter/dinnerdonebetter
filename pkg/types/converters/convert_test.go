package converters

import (
	"testing"

	"github.com/prixfixeco/backend/internal/pkg/pointers"
	"github.com/prixfixeco/backend/pkg/types"
)

func TestConvert(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		original := types.ValidPreparation{
			Name:        "Name",
			Description: "Description",
		}
		updateInput := types.ValidPreparationUpdateRequestInput{
			Slug: pointers.Pointer("Slug"),
		}
		updated := types.ValidPreparation{
			Name:        "Name",
			Description: "Description",
			Slug:        "Slug",
		}

		_, _, _ = original, updateInput, updated
	})
}
