package main

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
)

func Test_schemaFromType(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := encoding.MustEncodeJSON(map[string]any{
			"Filter": queryFilterSchema(),
			"Query": map[string]any{
				"type":        strType,
				"description": "Search query string to match ingredient names or descriptions",
			},
			"UseSearchService": map[string]any{
				"type":        boolType,
				"description": "Whether to use the search service for more advanced search capabilities",
			},
		})

		actual := encoding.MustEncodeJSON(jsonschema.Reflect(&SearchValidIngredientsInvocation{}))

		assert.Equal(t, expected, actual) // actual.Definitions[0].Properties.MarshalJSON())
	})
}
