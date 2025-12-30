package main

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	mealplanningconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type (
	GetValidIngredientsInvocation struct {
		ValidIngredientID string `jsonschema:"description=The ingredient MealPlanTaskID"`
	}
)

var validIngredientsSchema = map[string]any{
	"MealPlanTaskID":              stringField("The MealPlanTaskID of the valid ingredient"),
	"CreatedAt":                   timestampField("When the valid ingredient was created"),
	"LastUpdatedAt":               timestampField("When the valid ingredient was last updated"),
	"ArchivedAt":                  timestampField("When the valid ingredient was soft deleted"),
	"Name":                        stringField("Name of the ingredient"),
	"Description":                 stringField("Description of the ingredient"),
	"Warning":                     stringField("For things like allergen warnings"),
	"IconPath":                    stringField("The URL for the icon for the item"),
	"ContainsDairy":               boolField("Whether or not the valid ingredient contains dairy"),
	"ContainsPeanut":              boolField("Whether or not the valid ingredient contains peanut"),
	"ContainsTreeNut":             boolField("Whether or not the valid ingredient contains tree nut"),
	"ContainsEgg":                 boolField("Whether or not the valid ingredient contains egg"),
	"ContainsWheat":               boolField("Whether or not the valid ingredient contains wheat"),
	"ContainsShellfish":           boolField("Whether or not the valid ingredient contains shellfish"),
	"ContainsSesame":              boolField("Whether or not the valid ingredient contains sesame"),
	"ContainsFish":                boolField("Whether or not the valid ingredient contains fish"),
	"ContainsGluten":              boolField("Whether or not the valid ingredient contains gluten"),
	"AnimalFlesh":                 boolField("Whether or not the valid ingredient is derived from animal flesh"),
	"IsLiquid":                    boolField("Whether or not the valid ingredient is a liquid"),
	"ContainsSoy":                 boolField("Whether or not the valid ingredient contains soy"),
	"PluralName":                  stringField("The plural name for the ingredient. So for an ingredient named 'onion', this would be 'onions'"),
	"AnimalDerived":               boolField("Whether or not the valid ingredient AnimalDerived"),
	"RestrictToPreparations":      boolField("Whether or not the valid ingredient is restrictToPreparations"),
	"StorageTemperatureInCelsius": optionalFloatRangeSchema(),
	"StorageInstructions":         stringField("Instructions on how to store the item."),
	"Slug":                        stringField("An easy-to-use URL slug for the ingredient"),
	"ContainsAlcohol":             boolField("Whether or not the valid ingredient contains Alcohol"),
	"ShoppingSuggestions":         stringField("Suggestions for the user to keep in mind when shopping for this ingredient"),
	"IsStarch":                    boolField("Whether or not the valid ingredient is a Starch"),
	"IsProtein":                   boolField("Whether or not the valid ingredient is a Protein"),
	"IsGrain":                     boolField("Whether or not the valid ingredient is a Grain"),
	"IsFruit":                     boolField("Whether or not the valid ingredient is a Fruit"),
	"IsSalt":                      boolField("Whether or not the valid ingredient is a Salt"),
	"IsFat":                       boolField("Whether or not the valid ingredient is a Fat"),
	"IsAcid":                      boolField("Whether or not the valid ingredient is a Acid"),
	"IsHeat":                      boolField("Whether or not the valid ingredient is a Heat"),
}

var getValidIngredientTool = &mcp.Tool{
	Name:        "GetValidIngredient",
	Description: "Get a valid ingredient by it's MealPlanTaskID",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientID": stringField("The MealPlanTaskID of the valid ingredient to update"),
	}),
	OutputSchema: schemaObject(validIngredientsSchema),
}

func (h *mcpToolManager) GetValidIngredient() mcp.ToolHandlerFor[*GetValidIngredientsInvocation, *mealplanning.ValidIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidIngredientsInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredient, error) {
		result, err := h.client.GetValidIngredient(ctx, &mealplanninggrpc.GetValidIngredientRequest{
			ValidIngredientId: x.ValidIngredientID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result.Result), nil
	}
}

type (
	SearchValidIngredientsInvocation struct {
		Filter           *filtering.QueryFilter
		Query            string `jsonschema_description:"The ingredient name query"`
		UseSearchService bool   `jsonschema_description:"Whether or not to use a search index or just a database search"`
	}

	SearchValidIngredientsResult struct {
		Results []*mealplanning.ValidIngredient
	}
)

var searchForValidIngredientsTool = &mcp.Tool{
	Name:        "SearchForValidIngredients",
	Description: "Search for valid ingredients with optional filtering and query string",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
		"Query": map[string]any{
			"type":        strType,
			"description": "The ingredient name query",
		},
		"UseSearchService": map[string]any{
			"type":        boolType,
			"description": "Whether or not to use a search index or just a database search",
		},
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validIngredientsSchema)),
	}),
}

func (h *mcpToolManager) SearchForValidIngredients() mcp.ToolHandlerFor[*SearchValidIngredientsInvocation, *SearchValidIngredientsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchValidIngredientsInvocation) (*mcp.CallToolResult, *SearchValidIngredientsResult, error) {
		results, err := h.client.SearchForValidIngredients(ctx, &mealplanninggrpc.SearchForValidIngredientsRequest{
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchValidIngredientsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result))
		}

		return nil, out, nil
	}
}

var validIngredientCreationTool = &mcp.Tool{
	Name:        "CreateValidIngredient",
	Description: "Create a valid ingredient for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"Name":                        stringField("Name of the ingredient"),
		"Description":                 stringField("Description of the ingredient"),
		"Warning":                     stringField("For things like allergen warnings"),
		"IconPath":                    stringField("The URL for the icon for the item"),
		"ContainsDairy":               boolField("Whether or not the valid ingredient contains dairy"),
		"ContainsPeanut":              boolField("Whether or not the valid ingredient contains peanut"),
		"ContainsTreeNut":             boolField("Whether or not the valid ingredient contains tree nut"),
		"ContainsEgg":                 boolField("Whether or not the valid ingredient contains egg"),
		"ContainsWheat":               boolField("Whether or not the valid ingredient contains wheat"),
		"ContainsShellfish":           boolField("Whether or not the valid ingredient contains shellfish"),
		"ContainsSesame":              boolField("Whether or not the valid ingredient contains sesame"),
		"ContainsFish":                boolField("Whether or not the valid ingredient contains fish"),
		"ContainsGluten":              boolField("Whether or not the valid ingredient contains gluten"),
		"AnimalFlesh":                 boolField("Whether or not the valid ingredient is derived from animal flesh"),
		"IsLiquid":                    boolField("Whether or not the valid ingredient is a liquid"),
		"ContainsSoy":                 boolField("Whether or not the valid ingredient contains soy"),
		"PluralName":                  stringField("The plural name for the ingredient. So for an ingredient named 'onion', this would be 'onions'"),
		"AnimalDerived":               boolField("Whether or not the valid ingredient AnimalDerived"),
		"RestrictToPreparations":      boolField("Whether or not the valid ingredient is restrictToPreparations"),
		"StorageTemperatureInCelsius": optionalFloatRangeSchema(),
		"StorageInstructions":         stringField("Instructions on how to store the item."),
		"Slug":                        stringField("An easy-to-use URL slug for the ingredient"),
		"ContainsAlcohol":             boolField("Whether or not the valid ingredient contains Alcohol"),
		"ShoppingSuggestions":         stringField("Suggestions for the user to keep in mind when shopping for this ingredient"),
		"IsStarch":                    boolField("Whether or not the valid ingredient is a Starch"),
		"IsProtein":                   boolField("Whether or not the valid ingredient is a Protein"),
		"IsGrain":                     boolField("Whether or not the valid ingredient is a Grain"),
		"IsFruit":                     boolField("Whether or not the valid ingredient is a Fruit"),
		"IsSalt":                      boolField("Whether or not the valid ingredient is a Salt"),
		"IsFat":                       boolField("Whether or not the valid ingredient is a Fat"),
		"IsAcid":                      boolField("Whether or not the valid ingredient is a Acid"),
		"IsHeat":                      boolField("Whether or not the valid ingredient is a Heat"),
	}),
	OutputSchema: schemaObject(validIngredientsSchema),
}

func (h *mcpToolManager) CreateValidIngredient() mcp.ToolHandlerFor[*mealplanning.ValidIngredientCreationRequestInput, *mealplanning.ValidIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidIngredientCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidIngredient, error) {
		result, err := h.client.CreateValidIngredient(ctx, &mealplanninggrpc.CreateValidIngredientRequest{Input: mealplanningconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result.Result), nil
	}
}

type (
	UpdateValidIngredientInvocation struct {
		*mealplanning.ValidIngredientUpdateRequestInput
		ValidIngredientID string `jsonschema:"required,description=The ingredient MealPlanTaskID"`
	}
)

var validIngredientUpdateTool = &mcp.Tool{
	Name:        "UpdateValidIngredient",
	Description: "Update a valid ingredient for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientID":           stringField("The MealPlanTaskID of the valid ingredient to update"),
		"Name":                        stringField("Name of the ingredient"),
		"Description":                 stringField("Description of the ingredient"),
		"Warning":                     stringField("For things like allergen warnings"),
		"IconPath":                    stringField("The URL for the icon for the item"),
		"ContainsDairy":               boolField("Whether or not the valid ingredient contains dairy"),
		"ContainsPeanut":              boolField("Whether or not the valid ingredient contains peanut"),
		"ContainsTreeNut":             boolField("Whether or not the valid ingredient contains tree nut"),
		"ContainsEgg":                 boolField("Whether or not the valid ingredient contains egg"),
		"ContainsWheat":               boolField("Whether or not the valid ingredient contains wheat"),
		"ContainsShellfish":           boolField("Whether or not the valid ingredient contains shellfish"),
		"ContainsSesame":              boolField("Whether or not the valid ingredient contains sesame"),
		"ContainsFish":                boolField("Whether or not the valid ingredient contains fish"),
		"ContainsGluten":              boolField("Whether or not the valid ingredient contains gluten"),
		"AnimalFlesh":                 boolField("Whether or not the valid ingredient is derived from animal flesh"),
		"IsLiquid":                    boolField("Whether or not the valid ingredient is a liquid"),
		"ContainsSoy":                 boolField("Whether or not the valid ingredient contains soy"),
		"PluralName":                  stringField("The plural name for the ingredient. So for an ingredient named 'onion', this would be 'onions'"),
		"AnimalDerived":               boolField("Whether or not the valid ingredient AnimalDerived"),
		"RestrictToPreparations":      boolField("Whether or not the valid ingredient is restrictToPreparations"),
		"StorageTemperatureInCelsius": optionalFloatRangeSchema(),
		"StorageInstructions":         stringField("Instructions on how to store the item."),
		"Slug":                        stringField("An easy-to-use URL slug for the ingredient"),
		"ContainsAlcohol":             boolField("Whether or not the valid ingredient contains Alcohol"),
		"ShoppingSuggestions":         stringField("Suggestions for the user to keep in mind when shopping for this ingredient"),
		"IsStarch":                    boolField("Whether or not the valid ingredient is a Starch"),
		"IsProtein":                   boolField("Whether or not the valid ingredient is a Protein"),
		"IsGrain":                     boolField("Whether or not the valid ingredient is a Grain"),
		"IsFruit":                     boolField("Whether or not the valid ingredient is a Fruit"),
		"IsSalt":                      boolField("Whether or not the valid ingredient is a Salt"),
		"IsFat":                       boolField("Whether or not the valid ingredient is a Fat"),
		"IsAcid":                      boolField("Whether or not the valid ingredient is a Acid"),
		"IsHeat":                      boolField("Whether or not the valid ingredient is a Heat"),
	}),
	OutputSchema: schemaObject(validIngredientsSchema),
}

func (h *mcpToolManager) UpdateValidIngredient() mcp.ToolHandlerFor[*UpdateValidIngredientInvocation, *mealplanning.ValidIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidIngredientInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredient, error) {
		result, err := h.client.UpdateValidIngredient(ctx, &mealplanninggrpc.UpdateValidIngredientRequest{
			ValidIngredientId: x.ValidIngredientID,
			Input:             mealplanningconverters.ConvertValidIngredientUpdateRequestInputToGRPCValidIngredientUpdateRequestInput(x.ValidIngredientUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result.Result), nil
	}
}
