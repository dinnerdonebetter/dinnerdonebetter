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
	GetValidIngredientPreparationInvocation struct {
		ValidIngredientPreparationID string `jsonschema:"description=The ingredient preparation ID"`
	}
)

var validIngredientPreparationsSchema = map[string]any{
	"ID":            stringField("The ID of the valid ingredient preparation"),
	"CreatedAt":     timestampField("When the valid ingredient preparation was created"),
	"LastUpdatedAt": timestampField("When the valid ingredient preparation was last updated"),
	"ArchivedAt":    timestampField("When the valid ingredient preparation was soft deleted"),
	"Notes":         stringField("Notes about the ingredient preparation"),
	"Preparation":   objectType(validPreparationsSchema),
	"Ingredient":    objectType(validIngredientsSchema),
}

var getValidIngredientPreparationTool = &mcp.Tool{
	Name:        "GetValidIngredientPreparation",
	Description: "Get a valid ingredient preparation by it's ID",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientPreparationID": stringField("The ID of the valid ingredient preparation to get"),
	}),
	OutputSchema: schemaObject(validIngredientPreparationsSchema),
}

func (h *mcpToolManager) GetValidIngredientPreparation() mcp.ToolHandlerFor[*GetValidIngredientPreparationInvocation, *mealplanning.ValidIngredientPreparation] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidIngredientPreparationInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredientPreparation, error) {
		result, err := h.client.GetValidIngredientPreparation(ctx, &mealplanninggrpc.GetValidIngredientPreparationRequest{
			ValidIngredientPreparationID: x.ValidIngredientPreparationID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientPreparationToValidIngredientPreparation(result.Result), nil
	}
}

type (
	GetValidIngredientPreparationsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetValidIngredientPreparationsResult struct {
		Results []*mealplanning.ValidIngredientPreparation
	}
)

var getValidIngredientPreparationsTool = &mcp.Tool{
	Name:        "GetValidIngredientPreparations",
	Description: "Get valid ingredient preparations with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validIngredientPreparationsSchema)),
	}),
}

func (h *mcpToolManager) GetValidIngredientPreparations() mcp.ToolHandlerFor[*GetValidIngredientPreparationsInvocation, *GetValidIngredientPreparationsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidIngredientPreparationsInvocation) (*mcp.CallToolResult, *GetValidIngredientPreparationsResult, error) {
		results, err := h.client.GetValidIngredientPreparations(ctx, &mealplanninggrpc.GetValidIngredientPreparationsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidIngredientPreparationsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidIngredientPreparationToValidIngredientPreparation(result))
		}

		return nil, out, nil
	}
}

var validIngredientPreparationCreationTool = &mcp.Tool{
	Name:        "CreateValidIngredientPreparation",
	Description: "Create a valid ingredient preparation linking an ingredient to a preparation.",
	InputSchema: schemaObject(map[string]any{
		"Notes":              stringField("Notes about the ingredient preparation"),
		"ValidPreparationID": stringField("The ID of the valid preparation"),
		"ValidIngredientID":  stringField("The ID of the valid ingredient"),
	}),
	OutputSchema: schemaObject(validIngredientPreparationsSchema),
}

func (h *mcpToolManager) CreateValidIngredientPreparation() mcp.ToolHandlerFor[*mealplanning.ValidIngredientPreparationCreationRequestInput, *mealplanning.ValidIngredientPreparation] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidIngredientPreparationCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidIngredientPreparation, error) {
		result, err := h.client.CreateValidIngredientPreparation(ctx, &mealplanninggrpc.CreateValidIngredientPreparationRequest{Input: mealplanningconverters.ConvertCreateValidIngredientPreparationRequestToGRPCValidIngredientPreparationCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientPreparationToValidIngredientPreparation(result.Result), nil
	}
}

type (
	UpdateValidIngredientPreparationInvocation struct {
		*mealplanning.ValidIngredientPreparationUpdateRequestInput
		ValidIngredientPreparationID string `jsonschema:"required,description=The ingredient preparation ID"`
	}
)

var validIngredientPreparationUpdateTool = &mcp.Tool{
	Name:        "UpdateValidIngredientPreparation",
	Description: "Update a valid ingredient preparation.",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientPreparationID": stringField("The ID of the valid ingredient preparation to update"),
		"Notes":                        stringField("Notes about the ingredient preparation"),
		"ValidPreparationID":           stringField("The ID of the valid preparation"),
		"ValidIngredientID":            stringField("The ID of the valid ingredient"),
	}),
	OutputSchema: schemaObject(validIngredientPreparationsSchema),
}

func (h *mcpToolManager) UpdateValidIngredientPreparation() mcp.ToolHandlerFor[*UpdateValidIngredientPreparationInvocation, *mealplanning.ValidIngredientPreparation] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidIngredientPreparationInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredientPreparation, error) {
		result, err := h.client.UpdateValidIngredientPreparation(ctx, &mealplanninggrpc.UpdateValidIngredientPreparationRequest{
			ValidIngredientPreparationID: x.ValidIngredientPreparationID,
			Input:                        mealplanningconverters.ConvertValidIngredientPreparationUpdateRequestInputToGRPCValidIngredientPreparationUpdateRequestInput(x.ValidIngredientPreparationUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientPreparationToValidIngredientPreparation(result.Result), nil
	}
}

//
