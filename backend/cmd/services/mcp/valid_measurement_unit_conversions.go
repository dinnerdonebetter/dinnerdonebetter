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
	GetValidMeasurementUnitConversionInvocation struct {
		ValidMeasurementUnitConversionID string `jsonschema:"description=The measurement unit conversion ID"`
	}
)

var validMeasurementUnitConversionsSchema = map[string]any{
	"ID":                stringField("The ID of the valid measurement unit conversion"),
	"CreatedAt":         timestampField("When the valid measurement unit conversion was created"),
	"LastUpdatedAt":     timestampField("When the valid measurement unit conversion was last updated"),
	"ArchivedAt":        timestampField("When the valid measurement unit conversion was soft deleted"),
	"Notes":             stringField("Notes about the measurement unit conversion"),
	"Modifier":          floatField("The conversion modifier (multiplier to convert from 'From' unit to 'To' unit)"),
	"From":              objectType(validMeasurementUnitsSchema),
	"To":                objectType(validMeasurementUnitsSchema),
	"OnlyForIngredient": objectType(validIngredientsSchema),
}

var getValidMeasurementUnitConversionTool = &mcp.Tool{
	Name:        "GetValidMeasurementUnitConversion",
	Description: "Get a valid measurement unit conversion by it's ID",
	InputSchema: schemaObject(map[string]any{
		"ValidMeasurementUnitConversionID": stringField("The ID of the valid measurement unit conversion to get"),
	}),
	OutputSchema: schemaObject(validMeasurementUnitConversionsSchema),
}

func (h *mcpToolManager) GetValidMeasurementUnitConversion() mcp.ToolHandlerFor[*GetValidMeasurementUnitConversionInvocation, *mealplanning.ValidMeasurementUnitConversion] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidMeasurementUnitConversionInvocation) (*mcp.CallToolResult, *mealplanning.ValidMeasurementUnitConversion, error) {
		result, err := h.client.GetValidMeasurementUnitConversion(ctx, &mealplanninggrpc.GetValidMeasurementUnitConversionRequest{
			ValidMeasurementUnitConversionID: x.ValidMeasurementUnitConversionID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidMeasurementUnitConversionToValidMeasurementUnitConversion(result.Result), nil
	}
}

type (
	GetValidMeasurementUnitConversionsForUnitInvocation struct {
		Filter                 *filtering.QueryFilter
		ValidMeasurementUnitID string `jsonschema:"description=The measurement unit ID"`
	}

	GetValidMeasurementUnitConversionsForUnitResult struct {
		Results []*mealplanning.ValidMeasurementUnitConversion
	}
)

var getValidMeasurementUnitConversionsForUnitTool = &mcp.Tool{
	Name:        "GetValidMeasurementUnitConversionsForUnit",
	Description: "Get valid measurement unit conversions for a specific measurement unit with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter":                 queryFilterSchema(),
		"ValidMeasurementUnitID": stringField("The ID of the valid measurement unit"),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validMeasurementUnitConversionsSchema)),
	}),
}

func (h *mcpToolManager) GetValidMeasurementUnitConversionsForUnit() mcp.ToolHandlerFor[*GetValidMeasurementUnitConversionsForUnitInvocation, *GetValidMeasurementUnitConversionsForUnitResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidMeasurementUnitConversionsForUnitInvocation) (*mcp.CallToolResult, *GetValidMeasurementUnitConversionsForUnitResult, error) {
		results, err := h.client.GetValidMeasurementUnitConversionsForUnit(ctx, &mealplanninggrpc.GetValidMeasurementUnitConversionsForUnitRequest{
			Filter:                 grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
			ValidMeasurementUnitID: x.ValidMeasurementUnitID,
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidMeasurementUnitConversionsForUnitResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidMeasurementUnitConversionToValidMeasurementUnitConversion(result))
		}

		return nil, out, nil
	}
}

var validMeasurementUnitConversionCreationTool = &mcp.Tool{
	Name:        "CreateValidMeasurementUnitConversion",
	Description: "Create a valid measurement unit conversion.",
	InputSchema: schemaObject(map[string]any{
		"From":              stringField("The ID of the measurement unit to convert from"),
		"To":                stringField("The ID of the measurement unit to convert to"),
		"Modifier":          floatField("The conversion modifier (multiplier to convert from 'From' unit to 'To' unit)"),
		"Notes":             stringField("Notes about the measurement unit conversion"),
		"OnlyForIngredient": stringField("The ID of the ingredient this conversion is specific to (optional, if not provided, conversion is universal)"),
	}),
	OutputSchema: schemaObject(validMeasurementUnitConversionsSchema),
}

func (h *mcpToolManager) CreateValidMeasurementUnitConversion() mcp.ToolHandlerFor[*mealplanning.ValidMeasurementUnitConversionCreationRequestInput, *mealplanning.ValidMeasurementUnitConversion] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidMeasurementUnitConversionCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidMeasurementUnitConversion, error) {
		result, err := h.client.CreateValidMeasurementUnitConversion(ctx, &mealplanninggrpc.CreateValidMeasurementUnitConversionRequest{Input: mealplanningconverters.ConvertCreateValidMeasurementUnitConversionRequestToGRPCValidMeasurementUnitConversionCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidMeasurementUnitConversionToValidMeasurementUnitConversion(result.Result), nil
	}
}

type (
	UpdateValidMeasurementUnitConversionInvocation struct {
		*mealplanning.ValidMeasurementUnitConversionUpdateRequestInput
		ValidMeasurementUnitConversionID string `jsonschema:"required,description=The measurement unit conversion ID"`
	}
)

var validMeasurementUnitConversionUpdateTool = &mcp.Tool{
	Name:        "UpdateValidMeasurementUnitConversion",
	Description: "Update a valid measurement unit conversion.",
	InputSchema: schemaObject(map[string]any{
		"ValidMeasurementUnitConversionID": stringField("The ID of the valid measurement unit conversion to update"),
		"From":                             stringField("The ID of the measurement unit to convert from"),
		"To":                               stringField("The ID of the measurement unit to convert to"),
		"Modifier":                         floatField("The conversion modifier (multiplier to convert from 'From' unit to 'To' unit)"),
		"Notes":                            stringField("Notes about the measurement unit conversion"),
		"OnlyForIngredient":                stringField("The ID of the ingredient this conversion is specific to (optional, if not provided, conversion is universal)"),
	}),
	OutputSchema: schemaObject(validMeasurementUnitConversionsSchema),
}

func (h *mcpToolManager) UpdateValidMeasurementUnitConversion() mcp.ToolHandlerFor[*UpdateValidMeasurementUnitConversionInvocation, *mealplanning.ValidMeasurementUnitConversion] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidMeasurementUnitConversionInvocation) (*mcp.CallToolResult, *mealplanning.ValidMeasurementUnitConversion, error) {
		result, err := h.client.UpdateValidMeasurementUnitConversion(ctx, &mealplanninggrpc.UpdateValidMeasurementUnitConversionRequest{
			ValidMeasurementUnitConversionID: x.ValidMeasurementUnitConversionID,
			Input:                            mealplanningconverters.ConvertValidMeasurementUnitConversionUpdateRequestInputToGRPCValidMeasurementUnitConversionUpdateRequestInput(x.ValidMeasurementUnitConversionUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidMeasurementUnitConversionToValidMeasurementUnitConversion(result.Result), nil
	}
}

//
