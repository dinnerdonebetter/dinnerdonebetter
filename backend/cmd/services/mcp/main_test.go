package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_schemaFromType(T *testing.T) {
	T.Parallel()

	T.Run("SearchValidIngredientsInvocation", func(t *testing.T) {
		t.Parallel()

		expected := map[string]any{
			"$schema": jsonSchemaVersion,
			"type":    objType,
			"properties": map[string]any{
				"Filter": queryFilterSchema(),
				"Query": map[string]any{
					"type":        strType,
					"description": "The ingredient name query",
				},
				"UseSearchService": map[string]any{
					"type":        boolType,
					"description": "Whether or not to use a search index or just a database search",
				},
			},
		}
		actual := schemaForType(SearchValidIngredientsInvocation{})

		assert.Equal(t, expected, actual)
	})

	T.Run("SearchValidIngredientsResult", func(t *testing.T) {
		t.Parallel()

		expected := map[string]any{
			"$schema": jsonSchemaVersion,
			"type":    objType,
			"items": map[string]any{
				"type": arrType,
				"name": "Results",
				"properties": map[string]any{
					"createdAt": map[string]any{
						"type":   strType,
						"format": dtFmt,
					},
					"lastUpdatedAt": map[string]any{
						"type":   []any{strType, "null"},
						"format": dtFmt,
					},
					"archivedAt": map[string]any{
						"type":   []any{strType, "null"},
						"format": dtFmt,
					},
					"storageTemperatureInCelsius": map[string]any{
						"type": objType,
						"properties": map[string]any{
							"min": map[string]any{
								"type": []any{"number", "null"},
							},
							"max": map[string]any{
								"type": []any{"number", "null"},
							},
						},
					},
					"iconPath": map[string]any{
						"type": strType,
					},
					"warning": map[string]any{
						"type": strType,
					},
					"pluralName": map[string]any{
						"type": strType,
					},
					"storageInstructions": map[string]any{
						"type": strType,
					},
					"name": map[string]any{
						"type": strType,
					},
					"id": map[string]any{
						"type": strType,
					},
					"description": map[string]any{
						"type": strType,
					},
					"slug": map[string]any{
						"type": strType,
					},
					"shoppingSuggestions": map[string]any{
						"type": strType,
					},
					"containsShellfish": map[string]any{
						"type": boolType,
					},
					"isLiquid": map[string]any{
						"type": boolType,
					},
					"containsPeanut": map[string]any{
						"type": boolType,
					},
					"containsTreeNut": map[string]any{
						"type": boolType,
					},
					"containsEgg": map[string]any{
						"type": boolType,
					},
					"containsWheat": map[string]any{
						"type": boolType,
					},
					"containsSoy": map[string]any{
						"type": boolType,
					},
					"animalDerived": map[string]any{
						"type": boolType,
					},
					"restrictToPreparations": map[string]any{
						"type": boolType,
					},
					"containsSesame": map[string]any{
						"type": boolType,
					},
					"containsFish": map[string]any{
						"type": boolType,
					},
					"containsGluten": map[string]any{
						"type": boolType,
					},
					"containsDairy": map[string]any{
						"type": boolType,
					},
					"containsAlcohol": map[string]any{
						"type": boolType,
					},
					"animalFlesh": map[string]any{
						"type": boolType,
					},
					"isStarch": map[string]any{
						"type": boolType,
					},
					"isProtein": map[string]any{
						"type": boolType,
					},
					"isGrain": map[string]any{
						"type": boolType,
					},
					"isFruit": map[string]any{
						"type": boolType,
					},
					"isSalt": map[string]any{
						"type": boolType,
					},
					"isFat": map[string]any{
						"type": boolType,
					},
					"isAcid": map[string]any{
						"type": boolType,
					},
					"isHeat": map[string]any{
						"type": boolType,
					},
				},
			},
		}
		actual := schemaForType(SearchValidIngredientsResult{})

		assert.Equal(t, expected, actual)
	})

	T.Run("UpdateValidIngredientInvocation", func(t *testing.T) {
		t.Parallel()

		// Test with pointer (as used in actual code)
		actual := schemaForType(&UpdateValidIngredientInvocation{})
		assert.NotNil(t, actual)
		assert.Equal(t, jsonSchemaVersion, actual["$schema"])
		assert.Equal(t, objType, actual["type"])

		// Check that Input property is properly resolved (not just a $ref)
		props, ok := actual["properties"].(map[string]any)
		assert.True(t, ok, "should have properties")
		inputProp, ok := props["Input"].(map[string]any)
		assert.True(t, ok, "should have Input property")

		// Input should be an object type with properties, not just a $ref
		// Check that it doesn't have a $ref (which would make it a JSON blob)
		_, hasRef := inputProp["$ref"]
		assert.False(t, hasRef, "Input should not have a $ref, it should be expanded")

		// Check that it doesn't have allOf/oneOf (which would make it a JSON blob)
		_, hasAllOf := inputProp["allOf"]
		assert.False(t, hasAllOf, "Input should not have allOf")
		_, hasOneOf := inputProp["oneOf"]
		assert.False(t, hasOneOf, "Input should not have oneOf")

		// Check that it has properties expanded
		inputProps, ok := inputProp["properties"].(map[string]any)
		assert.True(t, ok, "Input should have properties map")
		assert.Greater(t, len(inputProps), 0, "Input should have at least one property")

		// Verify some expected properties exist
		_, hasName := inputProps["name"]
		_, hasDescription := inputProps["description"]
		assert.True(t, hasName || hasDescription, "Input should have name or description property")

		// Debug: print the actual schema to see what we're generating
		t.Logf("Input property schema: %+v", inputProp)
	})
}
