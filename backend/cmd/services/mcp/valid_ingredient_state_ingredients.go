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
	GetValidIngredientStateIngredientInvocation struct {
		ValidIngredientStateIngredientID string `jsonschema:"description=The ingredient state ingredient ID"`
	}
)

var validIngredientStateIngredientsSchema = map[string]any{
	"ID":              stringField("The ID of the valid ingredient state ingredient"),
	"CreatedAt":       timestampField("When the valid ingredient state ingredient was created"),
	"LastUpdatedAt":   timestampField("When the valid ingredient state ingredient was last updated"),
	"ArchivedAt":      timestampField("When the valid ingredient state ingredient was soft deleted"),
	"Notes":           stringField("Notes about the ingredient state ingredient"),
	"IngredientState": objectType(validIngredientStatesSchema),
	"Ingredient":      objectType(validIngredientsSchema),
}

var getValidIngredientStateIngredientTool = &mcp.Tool{
	Name:        "GetValidIngredientStateIngredient",
	Description: "Get a valid ingredient state ingredient by it's ID",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientStateIngredientID": stringField("The ID of the valid ingredient state ingredient to get"),
	}),
	OutputSchema: schemaObject(validIngredientStateIngredientsSchema),
}

func (h *mcpToolManager) GetValidIngredientStateIngredient() mcp.ToolHandlerFor[*GetValidIngredientStateIngredientInvocation, *mealplanning.ValidIngredientStateIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidIngredientStateIngredientInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredientStateIngredient, error) {
		result, err := h.client.GetValidIngredientStateIngredient(ctx, &mealplanninggrpc.GetValidIngredientStateIngredientRequest{
			ValidIngredientStateIngredientID: x.ValidIngredientStateIngredientID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientStateIngredientToValidIngredientStateIngredient(result.Result), nil
	}
}

type (
	GetValidIngredientStateIngredientsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetValidIngredientStateIngredientsResult struct {
		Results []*mealplanning.ValidIngredientStateIngredient
	}
)

var getValidIngredientStateIngredientsTool = &mcp.Tool{
	Name:        "GetValidIngredientStateIngredients",
	Description: "Get valid ingredient state ingredients with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(validIngredientStateIngredientsSchema),
	}),
}

func (h *mcpToolManager) GetValidIngredientStateIngredients() mcp.ToolHandlerFor[*GetValidIngredientStateIngredientsInvocation, *GetValidIngredientStateIngredientsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidIngredientStateIngredientsInvocation) (*mcp.CallToolResult, *GetValidIngredientStateIngredientsResult, error) {
		results, err := h.client.GetValidIngredientStateIngredients(ctx, &mealplanninggrpc.GetValidIngredientStateIngredientsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidIngredientStateIngredientsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidIngredientStateIngredientToValidIngredientStateIngredient(result))
		}

		return nil, out, nil
	}
}

var validIngredientStateIngredientCreationTool = &mcp.Tool{
	Name:        "CreateValidIngredientStateIngredient",
	Description: "Create a valid ingredient state ingredient linking an ingredient to an ingredient state.",
	InputSchema: schemaObject(map[string]any{
		"Notes":                  stringField("Notes about the ingredient state ingredient"),
		"ValidIngredientStateID": stringField("The ID of the valid ingredient state"),
		"ValidIngredientID":      stringField("The ID of the valid ingredient"),
	}),
	OutputSchema: schemaObject(validIngredientStateIngredientsSchema),
}

func (h *mcpToolManager) CreateValidIngredientStateIngredient() mcp.ToolHandlerFor[*mealplanning.ValidIngredientStateIngredientCreationRequestInput, *mealplanning.ValidIngredientStateIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidIngredientStateIngredientCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidIngredientStateIngredient, error) {
		result, err := h.client.CreateValidIngredientStateIngredient(ctx, &mealplanninggrpc.CreateValidIngredientStateIngredientRequest{Input: mealplanningconverters.ConvertCreateValidIngredientStateIngredientRequestToGRPCValidIngredientStateIngredientCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientStateIngredientToValidIngredientStateIngredient(result.Result), nil
	}
}

type (
	UpdateValidIngredientStateIngredientInvocation struct {
		*mealplanning.ValidIngredientStateIngredientUpdateRequestInput
		ValidIngredientStateIngredientID string `jsonschema:"required,description=The ingredient state ingredient ID"`
	}
)

var validIngredientStateIngredientUpdateTool = &mcp.Tool{
	Name:        "UpdateValidIngredientStateIngredient",
	Description: "Update a valid ingredient state ingredient.",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientStateIngredientID": stringField("The ID of the valid ingredient state ingredient to update"),
		"Notes":                            stringField("Notes about the ingredient state ingredient"),
		"ValidIngredientStateID":           stringField("The ID of the valid ingredient state"),
		"ValidIngredientID":                stringField("The ID of the valid ingredient"),
	}),
	OutputSchema: schemaObject(validIngredientStateIngredientsSchema),
}

func (h *mcpToolManager) UpdateValidIngredientStateIngredient() mcp.ToolHandlerFor[*UpdateValidIngredientStateIngredientInvocation, *mealplanning.ValidIngredientStateIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidIngredientStateIngredientInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredientStateIngredient, error) {
		result, err := h.client.UpdateValidIngredientStateIngredient(ctx, &mealplanninggrpc.UpdateValidIngredientStateIngredientRequest{
			ValidIngredientStateIngredientID: x.ValidIngredientStateIngredientID,
			Input:                            mealplanningconverters.ConvertValidIngredientStateIngredientUpdateRequestInputToGRPCValidIngredientStateIngredientUpdateRequestInput(x.ValidIngredientStateIngredientUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientStateIngredientToValidIngredientStateIngredient(result.Result), nil
	}
}

//
