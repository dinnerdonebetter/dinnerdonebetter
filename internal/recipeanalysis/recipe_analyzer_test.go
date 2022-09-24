package recipeanalysis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/simple"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

func TestRecipeGrapher_makeGraphForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.SkipNow()

		g := &recipeAnalyzer{
			tracer: tracing.NewTracerForTest(t.Name()),
		}

		ctx := context.Background()
		r := &types.Recipe{
			Steps: []*types.RecipeStep{
				{},
			},
		}

		expected := &simple.DirectedGraph{}

		actual, err := g.makeGraphForRecipe(ctx, r)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
