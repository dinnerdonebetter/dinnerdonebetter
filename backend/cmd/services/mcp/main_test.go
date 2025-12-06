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
					"description": "Search query string to match ingredient names or descriptions",
				},
				"UseSearchService": map[string]any{
					"type":        boolType,
					"description": "Whether to use the search service for more advanced search capabilities",
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
}
