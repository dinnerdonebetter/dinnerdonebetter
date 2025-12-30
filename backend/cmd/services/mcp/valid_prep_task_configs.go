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
	GetValidPrepTaskConfigInvocation struct {
		ValidPrepTaskConfigID string `jsonschema:"description=The prep task config MealPlanTaskID"`
	}
)

var validPrepTaskConfigsSchema = map[string]any{
	"MealPlanTaskID":              stringField("The MealPlanTaskID of the valid prep task config"),
	"CreatedAt":                   timestampField("When the valid prep task config was created"),
	"LastUpdatedAt":               timestampField("When the valid prep task config was last updated"),
	"ArchivedAt":                  timestampField("When the valid prep task config was soft deleted"),
	"StorageDurationInSeconds":    uint32RangeWithOptionalMaxSchema(),
	"StorageTemperatureInCelsius": optionalFloat32RangeSchema(),
	"StorageType":                 stringField("The type of storage container (e.g., covered, airtight, uncovered)"),
	"StorageInstructions":         stringField("Instructions for how to store the prepped ingredient"),
	"Notes":                       stringField("Additional notes about the prep task config"),
	"Source":                      stringField("The source of this prep task config information"),
	"Preparation":                 objectType(validPreparationsSchema),
	"Ingredient":                  objectType(validIngredientsSchema),
}

var getValidPrepTaskConfigTool = &mcp.Tool{
	Name:        "GetValidPrepTaskConfig",
	Description: "Get a valid prep task config by its MealPlanTaskID. A prep task config defines how long a prepped ingredient can be stored under specific conditions.",
	InputSchema: schemaObject(map[string]any{
		"ValidPrepTaskConfigID": stringField("The MealPlanTaskID of the valid prep task config to get"),
	}),
	OutputSchema: schemaObject(validPrepTaskConfigsSchema),
}

func (h *mcpToolManager) GetValidPrepTaskConfig() mcp.ToolHandlerFor[*GetValidPrepTaskConfigInvocation, *mealplanning.ValidPrepTaskConfig] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPrepTaskConfigInvocation) (*mcp.CallToolResult, *mealplanning.ValidPrepTaskConfig, error) {
		result, err := h.client.GetValidPrepTaskConfig(ctx, &mealplanninggrpc.GetValidPrepTaskConfigRequest{
			ValidPrepTaskConfigId: x.ValidPrepTaskConfigID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPrepTaskConfigToValidPrepTaskConfig(result.Result), nil
	}
}

type (
	GetValidPrepTaskConfigsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetValidPrepTaskConfigsResult struct {
		Results []*mealplanning.ValidPrepTaskConfig
	}
)

var getValidPrepTaskConfigsTool = &mcp.Tool{
	Name:        "GetValidPrepTaskConfigs",
	Description: "Get valid prep task configs with optional filtering. Prep task configs define how long prepped ingredients can be stored.",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validPrepTaskConfigsSchema)),
	}),
}

func (h *mcpToolManager) GetValidPrepTaskConfigs() mcp.ToolHandlerFor[*GetValidPrepTaskConfigsInvocation, *GetValidPrepTaskConfigsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPrepTaskConfigsInvocation) (*mcp.CallToolResult, *GetValidPrepTaskConfigsResult, error) {
		results, err := h.client.GetValidPrepTaskConfigs(ctx, &mealplanninggrpc.GetValidPrepTaskConfigsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidPrepTaskConfigsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidPrepTaskConfigToValidPrepTaskConfig(result))
		}

		return nil, out, nil
	}
}

type (
	GetValidPrepTaskConfigsByIngredientInvocation struct {
		Filter            *filtering.QueryFilter
		ValidIngredientID string `jsonschema:"description=The ingredient MealPlanTaskID to filter by"`
	}
)

var getValidPrepTaskConfigsByIngredientTool = &mcp.Tool{
	Name:        "GetValidPrepTaskConfigsByIngredient",
	Description: "Get valid prep task configs for a specific ingredient. Use this to find storage information for a particular ingredient.",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientID": stringField("The MealPlanTaskID of the ingredient to get prep task configs for"),
		"Filter":            queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validPrepTaskConfigsSchema)),
	}),
}

func (h *mcpToolManager) GetValidPrepTaskConfigsByIngredient() mcp.ToolHandlerFor[*GetValidPrepTaskConfigsByIngredientInvocation, *GetValidPrepTaskConfigsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPrepTaskConfigsByIngredientInvocation) (*mcp.CallToolResult, *GetValidPrepTaskConfigsResult, error) {
		results, err := h.client.GetValidPrepTaskConfigsByIngredient(ctx, &mealplanninggrpc.GetValidPrepTaskConfigsByIngredientRequest{
			ValidIngredientId: x.ValidIngredientID,
			Filter:            grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidPrepTaskConfigsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidPrepTaskConfigToValidPrepTaskConfig(result))
		}

		return nil, out, nil
	}
}

type (
	GetValidPrepTaskConfigsByPreparationInvocation struct {
		Filter             *filtering.QueryFilter
		ValidPreparationID string `jsonschema:"description=The preparation MealPlanTaskID to filter by"`
	}
)

var getValidPrepTaskConfigsByPreparationTool = &mcp.Tool{
	Name:        "GetValidPrepTaskConfigsByPreparation",
	Description: "Get valid prep task configs for a specific preparation method. Use this to find storage information for ingredients prepared a certain way.",
	InputSchema: schemaObject(map[string]any{
		"ValidPreparationID": stringField("The MealPlanTaskID of the preparation to get prep task configs for"),
		"Filter":             queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validPrepTaskConfigsSchema)),
	}),
}

func (h *mcpToolManager) GetValidPrepTaskConfigsByPreparation() mcp.ToolHandlerFor[*GetValidPrepTaskConfigsByPreparationInvocation, *GetValidPrepTaskConfigsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPrepTaskConfigsByPreparationInvocation) (*mcp.CallToolResult, *GetValidPrepTaskConfigsResult, error) {
		results, err := h.client.GetValidPrepTaskConfigsByPreparation(ctx, &mealplanninggrpc.GetValidPrepTaskConfigsByPreparationRequest{
			ValidPreparationId: x.ValidPreparationID,
			Filter:             grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidPrepTaskConfigsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidPrepTaskConfigToValidPrepTaskConfig(result))
		}

		return nil, out, nil
	}
}

type (
	GetValidPrepTaskConfigsByIngredientAndPreparationInvocation struct {
		Filter             *filtering.QueryFilter
		ValidIngredientID  string `jsonschema:"description=The ingredient MealPlanTaskID to filter by"`
		ValidPreparationID string `jsonschema:"description=The preparation MealPlanTaskID to filter by"`
	}
)

var getValidPrepTaskConfigsByIngredientAndPreparationTool = &mcp.Tool{
	Name:        "GetValidPrepTaskConfigsByIngredientAndPreparation",
	Description: "Get valid prep task configs for a specific ingredient and preparation combination. Use this to find exactly how long a specific prepped ingredient (e.g., diced onions) can be stored.",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientID":  stringField("The MealPlanTaskID of the ingredient"),
		"ValidPreparationID": stringField("The MealPlanTaskID of the preparation"),
		"Filter":             queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validPrepTaskConfigsSchema)),
	}),
}

func (h *mcpToolManager) GetValidPrepTaskConfigsByIngredientAndPreparation() mcp.ToolHandlerFor[*GetValidPrepTaskConfigsByIngredientAndPreparationInvocation, *GetValidPrepTaskConfigsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPrepTaskConfigsByIngredientAndPreparationInvocation) (*mcp.CallToolResult, *GetValidPrepTaskConfigsResult, error) {
		results, err := h.client.GetValidPrepTaskConfigsByIngredientAndPreparation(ctx, &mealplanninggrpc.GetValidPrepTaskConfigsByIngredientAndPreparationRequest{
			ValidIngredientId:  x.ValidIngredientID,
			ValidPreparationId: x.ValidPreparationID,
			Filter:             grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidPrepTaskConfigsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidPrepTaskConfigToValidPrepTaskConfig(result))
		}

		return nil, out, nil
	}
}

var validPrepTaskConfigCreationTool = &mcp.Tool{
	Name:        "CreateValidPrepTaskConfig",
	Description: "Create a valid prep task config defining how long a prepped ingredient can be stored. For example, diced onions can be stored for 72 hours in an airtight container at refrigerator temperature.",
	InputSchema: schemaObject(map[string]any{
		"StorageDurationInSeconds":    uint32RangeWithOptionalMaxSchema(),
		"StorageTemperatureInCelsius": optionalFloat32RangeSchema(),
		"StorageType":                 stringField("The type of storage container (e.g., covered, airtight, uncovered)"),
		"StorageInstructions":         stringField("Instructions for how to store the prepped ingredient"),
		"Notes":                       stringField("Additional notes about the prep task config"),
		"Source":                      stringField("The source of this prep task config information"),
		"ValidPreparationID":          stringField("The MealPlanTaskID of the valid preparation (required)"),
		"ValidIngredientID":           stringField("The MealPlanTaskID of the valid ingredient (required)"),
	}),
	OutputSchema: schemaObject(validPrepTaskConfigsSchema),
}

func (h *mcpToolManager) CreateValidPrepTaskConfig() mcp.ToolHandlerFor[*mealplanning.ValidPrepTaskConfigCreationRequestInput, *mealplanning.ValidPrepTaskConfig] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidPrepTaskConfigCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidPrepTaskConfig, error) {
		result, err := h.client.CreateValidPrepTaskConfig(ctx, &mealplanninggrpc.CreateValidPrepTaskConfigRequest{
			Input: mealplanningconverters.ConvertValidPrepTaskConfigCreationRequestInputToGRPCValidPrepTaskConfigCreationRequestInput(x),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPrepTaskConfigToValidPrepTaskConfig(result.Result), nil
	}
}

type (
	UpdateValidPrepTaskConfigInvocation struct {
		*mealplanning.ValidPrepTaskConfigUpdateRequestInput
		ValidPrepTaskConfigID string `jsonschema:"required,description=The prep task config MealPlanTaskID"`
	}
)

var validPrepTaskConfigUpdateTool = &mcp.Tool{
	Name:        "UpdateValidPrepTaskConfig",
	Description: "Update a valid prep task config.",
	InputSchema: schemaObject(map[string]any{
		"ValidPrepTaskConfigID":       stringField("The MealPlanTaskID of the valid prep task config to update"),
		"StorageDurationInSeconds":    uint32RangeWithOptionalMaxSchema(),
		"StorageTemperatureInCelsius": optionalFloat32RangeSchema(),
		"StorageType":                 stringField("The type of storage container"),
		"StorageInstructions":         stringField("Instructions for how to store the prepped ingredient"),
		"Notes":                       stringField("Additional notes about the prep task config"),
		"Source":                      stringField("The source of this prep task config information"),
		"ValidPreparationID":          stringField("The MealPlanTaskID of the valid preparation"),
		"ValidIngredientID":           stringField("The MealPlanTaskID of the valid ingredient"),
	}),
	OutputSchema: schemaObject(validPrepTaskConfigsSchema),
}

func (h *mcpToolManager) UpdateValidPrepTaskConfig() mcp.ToolHandlerFor[*UpdateValidPrepTaskConfigInvocation, *mealplanning.ValidPrepTaskConfig] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidPrepTaskConfigInvocation) (*mcp.CallToolResult, *mealplanning.ValidPrepTaskConfig, error) {
		result, err := h.client.UpdateValidPrepTaskConfig(ctx, &mealplanninggrpc.UpdateValidPrepTaskConfigRequest{
			ValidPrepTaskConfigId: x.ValidPrepTaskConfigID,
			Input:                 mealplanningconverters.ConvertValidPrepTaskConfigUpdateRequestInputToGRPCValidPrepTaskConfigUpdateRequestInput(x.ValidPrepTaskConfigUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPrepTaskConfigToValidPrepTaskConfig(result.Result), nil
	}
}
