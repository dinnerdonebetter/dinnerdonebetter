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
	GetValidIngredientMeasurementUnitInvocation struct {
		ValidIngredientMeasurementUnitID string `jsonschema:"description=The ingredient measurement unit MealPlanTaskID"`
	}
)

var validIngredientMeasurementUnitsSchema = map[string]any{
	"MealPlanTaskID":    stringField("The MealPlanTaskID of the valid ingredient measurement unit"),
	"CreatedAt":         timestampField("When the valid ingredient measurement unit was created"),
	"LastUpdatedAt":     timestampField("When the valid ingredient measurement unit was last updated"),
	"ArchivedAt":        timestampField("When the valid ingredient measurement unit was soft deleted"),
	"Notes":             stringField("Notes about the ingredient measurement unit"),
	"AllowableQuantity": float32RangeWithOptionalMaxSchema(),
	"MeasurementUnit":   objectType(validMeasurementUnitsSchema),
	"Ingredient":        objectType(validIngredientsSchema),
}

var getValidIngredientMeasurementUnitTool = &mcp.Tool{
	Name:        "GetValidIngredientMeasurementUnit",
	Description: "Get a valid ingredient measurement unit by it's MealPlanTaskID",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientMeasurementUnitID": stringField("The MealPlanTaskID of the valid ingredient measurement unit to get"),
	}),
	OutputSchema: schemaObject(validIngredientMeasurementUnitsSchema),
}

func (h *mcpToolManager) GetValidIngredientMeasurementUnit() mcp.ToolHandlerFor[*GetValidIngredientMeasurementUnitInvocation, *mealplanning.ValidIngredientMeasurementUnit] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidIngredientMeasurementUnitInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredientMeasurementUnit, error) {
		result, err := h.client.GetValidIngredientMeasurementUnit(ctx, &mealplanninggrpc.GetValidIngredientMeasurementUnitRequest{
			ValidIngredientMeasurementUnitId: x.ValidIngredientMeasurementUnitID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientMeasurementUnitToValidIngredientMeasurementUnit(result.Result), nil
	}
}

type (
	GetValidIngredientMeasurementUnitsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetValidIngredientMeasurementUnitsResult struct {
		Results []*mealplanning.ValidIngredientMeasurementUnit
	}
)

var getValidIngredientMeasurementUnitsTool = &mcp.Tool{
	Name:        "GetValidIngredientMeasurementUnits",
	Description: "Get valid ingredient measurement units with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validIngredientMeasurementUnitsSchema)),
	}),
}

func (h *mcpToolManager) GetValidIngredientMeasurementUnits() mcp.ToolHandlerFor[*GetValidIngredientMeasurementUnitsInvocation, *GetValidIngredientMeasurementUnitsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidIngredientMeasurementUnitsInvocation) (*mcp.CallToolResult, *GetValidIngredientMeasurementUnitsResult, error) {
		results, err := h.client.GetValidIngredientMeasurementUnits(ctx, &mealplanninggrpc.GetValidIngredientMeasurementUnitsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidIngredientMeasurementUnitsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidIngredientMeasurementUnitToValidIngredientMeasurementUnit(result))
		}

		return nil, out, nil
	}
}

var validIngredientMeasurementUnitCreationTool = &mcp.Tool{
	Name:        "CreateValidIngredientMeasurementUnit",
	Description: "Create a valid ingredient measurement unit linking an ingredient to a measurement unit.",
	InputSchema: schemaObject(map[string]any{
		"Notes":                  stringField("Notes about the ingredient measurement unit"),
		"ValidMeasurementUnitID": stringField("The MealPlanTaskID of the valid measurement unit"),
		"ValidIngredientID":      stringField("The MealPlanTaskID of the valid ingredient"),
		"AllowableQuantity":      float32RangeWithOptionalMaxSchema(),
	}),
	OutputSchema: schemaObject(validIngredientMeasurementUnitsSchema),
}

func (h *mcpToolManager) CreateValidIngredientMeasurementUnit() mcp.ToolHandlerFor[*mealplanning.ValidIngredientMeasurementUnitCreationRequestInput, *mealplanning.ValidIngredientMeasurementUnit] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidIngredientMeasurementUnitCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidIngredientMeasurementUnit, error) {
		result, err := h.client.CreateValidIngredientMeasurementUnit(ctx, &mealplanninggrpc.CreateValidIngredientMeasurementUnitRequest{Input: mealplanningconverters.ConvertCreateValidIngredientMeasurementUnitRequestToGRPCValidIngredientMeasurementUnitCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientMeasurementUnitToValidIngredientMeasurementUnit(result.Result), nil
	}
}

type (
	UpdateValidIngredientMeasurementUnitInvocation struct {
		*mealplanning.ValidIngredientMeasurementUnitUpdateRequestInput
		ValidIngredientMeasurementUnitID string `jsonschema:"required,description=The ingredient measurement unit MealPlanTaskID"`
	}
)

var validIngredientMeasurementUnitUpdateTool = &mcp.Tool{
	Name:        "UpdateValidIngredientMeasurementUnit",
	Description: "Update a valid ingredient measurement unit.",
	InputSchema: schemaObject(map[string]any{
		"ValidIngredientMeasurementUnitID": stringField("The MealPlanTaskID of the valid ingredient measurement unit to update"),
		"Notes":                            stringField("Notes about the ingredient measurement unit"),
		"ValidMeasurementUnitID":           stringField("The MealPlanTaskID of the valid measurement unit"),
		"ValidIngredientID":                stringField("The MealPlanTaskID of the valid ingredient"),
		"AllowableQuantity":                float32RangeWithOptionalMaxSchema(),
	}),
	OutputSchema: schemaObject(validIngredientMeasurementUnitsSchema),
}

func (h *mcpToolManager) UpdateValidIngredientMeasurementUnit() mcp.ToolHandlerFor[*UpdateValidIngredientMeasurementUnitInvocation, *mealplanning.ValidIngredientMeasurementUnit] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidIngredientMeasurementUnitInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredientMeasurementUnit, error) {
		result, err := h.client.UpdateValidIngredientMeasurementUnit(ctx, &mealplanninggrpc.UpdateValidIngredientMeasurementUnitRequest{
			ValidIngredientMeasurementUnitId: x.ValidIngredientMeasurementUnitID,
			Input:                            mealplanningconverters.ConvertValidIngredientMeasurementUnitUpdateRequestInputToGRPCValidIngredientMeasurementUnitUpdateRequestInput(x.ValidIngredientMeasurementUnitUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientMeasurementUnitToValidIngredientMeasurementUnit(result.Result), nil
	}
}

//
