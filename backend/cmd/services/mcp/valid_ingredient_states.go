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
	GetValidIngredientStateInvocation struct {
		ValidIngredientStateID string `jsonschema:"description=The ingredient state MealPlanTaskID"`
	}
)

var validIngredientStatesSchema = map[string]any{
	"MealPlanTaskID": stringField("The MealPlanTaskID of the valid ingredient state"),
	"CreatedAt":      timestampField("When the valid ingredient state was created"),
	"LastUpdatedAt":  timestampField("When the valid ingredient state was last updated"),
	"ArchivedAt":     timestampField("When the valid ingredient state was soft deleted"),
	"Name":           stringField("Name of the ingredient state"),
	"Description":    stringField("Description of the ingredient state"),
	"IconPath":       stringField("The URL for the icon for the item"),
	"Slug":           stringField("An easy-to-use URL slug for the ingredient state"),
	"PastTense":      stringField("The past tense form of the ingredient state name (e.g., 'chopped' for 'chop')"),
	"AttributeType":  stringField("The attribute type of the ingredient state (texture, consistency, temperature, color, appearance, odor, taste, sound, or other)"),
}

var getValidIngredientStateTool = &mcp.Tool{
	Name:        "GetValidIngredientState",
	Description: "Get a valid ingredient state by it's MealPlanTaskID",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientStateID": stringField("The MealPlanTaskID of the valid ingredient state to get"),
	}),
	OutputSchema: schemaObject(validIngredientStatesSchema),
}

func (h *mcpToolManager) GetValidIngredientState() mcp.ToolHandlerFor[*GetValidIngredientStateInvocation, *mealplanning.ValidIngredientState] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidIngredientStateInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredientState, error) {
		result, err := h.client.GetValidIngredientState(ctx, &mealplanninggrpc.GetValidIngredientStateRequest{
			ValidIngredientStateId: x.ValidIngredientStateID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientStateToValidIngredientState(result.Result), nil
	}
}

type (
	SearchValidIngredientStatesInvocation struct {
		Filter           *filtering.QueryFilter
		Query            string `jsonschema_description:"The ingredient state name query"`
		UseSearchService bool   `jsonschema_description:"Whether or not to use a search index or just a database search"`
	}

	SearchValidIngredientStatesResult struct {
		Results []*mealplanning.ValidIngredientState
	}
)

var searchForValidIngredientStatesTool = &mcp.Tool{
	Name:        "SearchForValidIngredientStates",
	Description: "Search for valid ingredient states with optional filtering and query string",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
		"Query": map[string]any{
			"type":        strType,
			"description": "The ingredient state name query",
		},
		"UseSearchService": map[string]any{
			"type":        boolType,
			"description": "Whether or not to use a search index or just a database search",
		},
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validIngredientStatesSchema)),
	}),
}

func (h *mcpToolManager) SearchForValidIngredientStates() mcp.ToolHandlerFor[*SearchValidIngredientStatesInvocation, *SearchValidIngredientStatesResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchValidIngredientStatesInvocation) (*mcp.CallToolResult, *SearchValidIngredientStatesResult, error) {
		results, err := h.client.SearchForValidIngredientStates(ctx, &mealplanninggrpc.SearchForValidIngredientStatesRequest{
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchValidIngredientStatesResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidIngredientStateToValidIngredientState(result))
		}

		return nil, out, nil
	}
}

var validIngredientStateCreationTool = &mcp.Tool{
	Name:        "CreateValidIngredientState",
	Description: "Create a valid ingredient state for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"Name":          stringField("Name of the ingredient state"),
		"Description":   stringField("Description of the ingredient state"),
		"IconPath":      stringField("The URL for the icon for the item"),
		"Slug":          stringField("An easy-to-use URL slug for the ingredient state"),
		"PastTense":     stringField("The past tense form of the ingredient state name (e.g., 'chopped' for 'chop')"),
		"AttributeType": stringField("The attribute type of the ingredient state (texture, consistency, temperature, color, appearance, odor, taste, sound, or other)"),
	}),
	OutputSchema: schemaObject(validIngredientStatesSchema),
}

func (h *mcpToolManager) CreateValidIngredientState() mcp.ToolHandlerFor[*mealplanning.ValidIngredientStateCreationRequestInput, *mealplanning.ValidIngredientState] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidIngredientStateCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidIngredientState, error) {
		result, err := h.client.CreateValidIngredientState(ctx, &mealplanninggrpc.CreateValidIngredientStateRequest{Input: mealplanningconverters.ConvertValidIngredientStateCreationRequestInputToGRPCValidIngredientStateCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientStateToValidIngredientState(result.Result), nil
	}
}

type (
	UpdateValidIngredientStateInvocation struct {
		*mealplanning.ValidIngredientStateUpdateRequestInput
		ValidIngredientStateID string `jsonschema:"required,description=The ingredient state MealPlanTaskID"`
	}
)

var validIngredientStateUpdateTool = &mcp.Tool{
	Name:        "UpdateValidIngredientState",
	Description: "Update a valid ingredient state for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientStateID": stringField("The MealPlanTaskID of the valid ingredient state to update"),
		"Name":                   stringField("Name of the ingredient state"),
		"Description":            stringField("Description of the ingredient state"),
		"IconPath":               stringField("The URL for the icon for the item"),
		"Slug":                   stringField("An easy-to-use URL slug for the ingredient state"),
		"PastTense":              stringField("The past tense form of the ingredient state name (e.g., 'chopped' for 'chop')"),
		"AttributeType":          stringField("The attribute type of the ingredient state (texture, consistency, temperature, color, appearance, odor, taste, sound, or other)"),
	}),
	OutputSchema: schemaObject(validIngredientStatesSchema),
}

func (h *mcpToolManager) UpdateValidIngredientState() mcp.ToolHandlerFor[*UpdateValidIngredientStateInvocation, *mealplanning.ValidIngredientState] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidIngredientStateInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredientState, error) {
		result, err := h.client.UpdateValidIngredientState(ctx, &mealplanninggrpc.UpdateValidIngredientStateRequest{
			ValidIngredientStateId: x.ValidIngredientStateID,
			Input:                  mealplanningconverters.ConvertValidIngredientStateUpdateRequestInputToGRPCValidIngredientStateUpdateRequestInput(x.ValidIngredientStateUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientStateToValidIngredientState(result.Result), nil
	}
}

//
