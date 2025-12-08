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
	GetRecipeStepProductInvocation struct {
		RecipeID            string `jsonschema:"description=The recipe ID"`
		RecipeStepID        string `jsonschema:"description=The recipe step ID"`
		RecipeStepProductID string `jsonschema:"description=The recipe step product ID"`
	}
)

var recipeStepProductsSchema = map[string]any{
	"ID":                          stringField("The ID of the recipe step product"),
	"CreatedAt":                   timestampField("When the recipe step product was created"),
	"LastUpdatedAt":               timestampField("When the recipe step product was last updated"),
	"ArchivedAt":                  timestampField("When the recipe step product was soft deleted"),
	"BelongsToRecipeStep":         stringField("The ID of the recipe step this product belongs to"),
	"Name":                        stringField("Name of the product"),
	"Type":                        stringField("The type of product (e.g., 'ingredient', 'waste', 'intermediate')"),
	"QuantityNotes":               stringField("Notes about the quantity"),
	"StorageInstructions":         stringField("Storage instructions for the product"),
	"MeasurementUnit":             objectType(validMeasurementUnitsSchema),
	"Quantity":                    optionalFloatRangeSchema(),
	"StorageTemperatureInCelsius": optionalFloatRangeSchema(),
	"StorageDurationInSeconds":    optionalUint32RangeSchema(),
	"ContainedInVesselIndex":      uintField("The index of the vessel this product is contained in, if any"),
	"Index":                       uintField("The display index/order of this product"),
	"IsWaste":                     boolField("Whether this product is waste"),
	"IsLiquid":                    boolField("Whether this product is a liquid"),
	"Compostable":                 boolField("Whether this product is compostable"),
}

var getRecipeStepProductTool = &mcp.Tool{
	Name:        "GetRecipeStepProduct",
	Description: "Get a recipe step product by it's ID",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":            stringField("The ID of the recipe"),
		"RecipeStepID":        stringField("The ID of the recipe step"),
		"RecipeStepProductID": stringField("The ID of the recipe step product to get"),
	}),
	OutputSchema: schemaObject(recipeStepProductsSchema),
}

func (h *mcpToolManager) GetRecipeStepProduct() mcp.ToolHandlerFor[*GetRecipeStepProductInvocation, *mealplanning.RecipeStepProduct] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepProductInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepProduct, error) {
		result, err := h.client.GetRecipeStepProduct(ctx, &mealplanninggrpc.GetRecipeStepProductRequest{
			RecipeID:            x.RecipeID,
			RecipeStepID:        x.RecipeStepID,
			RecipeStepProductID: x.RecipeStepProductID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepProductToRecipeStepProduct(result.Result), nil
	}
}

type (
	GetRecipeStepProductsInvocation struct {
		Filter       *filtering.QueryFilter
		RecipeID     string
		RecipeStepID string
	}

	GetRecipeStepProductsResult struct {
		Results []*mealplanning.RecipeStepProduct
	}
)

var getRecipeStepProductsTool = &mcp.Tool{
	Name:        "GetRecipeStepProducts",
	Description: "Get recipe step products with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":     stringField("The ID of the recipe"),
		"RecipeStepID": stringField("The ID of the recipe step"),
		"Filter":       queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(recipeStepProductsSchema)),
	}),
}

func (h *mcpToolManager) GetRecipeStepProducts() mcp.ToolHandlerFor[*GetRecipeStepProductsInvocation, *GetRecipeStepProductsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepProductsInvocation) (*mcp.CallToolResult, *GetRecipeStepProductsResult, error) {
		results, err := h.client.GetRecipeStepProducts(ctx, &mealplanninggrpc.GetRecipeStepProductsRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
			Filter:       grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetRecipeStepProductsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCRecipeStepProductToRecipeStepProduct(result))
		}

		return nil, out, nil
	}
}

type (
	CreateRecipeStepProductInvocation struct {
		*mealplanning.RecipeStepProductCreationRequestInput
		RecipeID     string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID string `jsonschema:"required,description=The recipe step ID"`
	}
)

var recipeStepProductCreationTool = &mcp.Tool{
	Name:        "CreateRecipeStepProduct",
	Description: "Create a recipe step product",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                    stringField("The ID of the recipe"),
		"RecipeStepID":                stringField("The ID of the recipe step"),
		"Name":                        stringField("Name of the product"),
		"Type":                        stringField("The type of product (e.g., 'ingredient', 'waste', 'intermediate')"),
		"MeasurementUnitID":           stringField("The ID of the measurement unit, if any"),
		"QuantityNotes":               stringField("Notes about the quantity"),
		"StorageInstructions":         stringField("Storage instructions for the product"),
		"Quantity":                    optionalFloatRangeSchema(),
		"StorageTemperatureInCelsius": optionalFloatRangeSchema(),
		"StorageDurationInSeconds":    optionalUint32RangeSchema(),
		"ContainedInVesselIndex":      uintField("The index of the vessel this product is contained in, if any"),
		"Index":                       uintField("The display index/order of this product"),
		"IsWaste":                     boolField("Whether this product is waste"),
		"IsLiquid":                    boolField("Whether this product is a liquid"),
		"Compostable":                 boolField("Whether this product is compostable"),
	}),
	OutputSchema: schemaObject(recipeStepProductsSchema),
}

func (h *mcpToolManager) CreateRecipeStepProduct() mcp.ToolHandlerFor[*CreateRecipeStepProductInvocation, *mealplanning.RecipeStepProduct] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateRecipeStepProductInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepProduct, error) {
		result, err := h.client.CreateRecipeStepProduct(ctx, &mealplanninggrpc.CreateRecipeStepProductRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
			Input:        mealplanningconverters.ConvertRecipeStepProductCreationRequestInputToGRPCRecipeStepProductCreationRequestInput(x.RecipeStepProductCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepProductToRecipeStepProduct(result.Created), nil
	}
}

type (
	UpdateRecipeStepProductInvocation struct {
		*mealplanning.RecipeStepProductUpdateRequestInput
		RecipeID            string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID        string `jsonschema:"required,description=The recipe step ID"`
		RecipeStepProductID string `jsonschema:"required,description=The recipe step product ID"`
	}
)

var recipeStepProductUpdateTool = &mcp.Tool{
	Name:        "UpdateRecipeStepProduct",
	Description: "Update a recipe step product",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                    stringField("The ID of the recipe"),
		"RecipeStepID":                stringField("The ID of the recipe step"),
		"RecipeStepProductID":         stringField("The ID of the recipe step product to update"),
		"Name":                        stringField("Name of the product"),
		"Type":                        stringField("The type of product (e.g., 'ingredient', 'waste', 'intermediate')"),
		"MeasurementUnitID":           stringField("The ID of the measurement unit, if any"),
		"QuantityNotes":               stringField("Notes about the quantity"),
		"StorageInstructions":         stringField("Storage instructions for the product"),
		"Quantity":                    optionalFloatRangeSchema(),
		"StorageTemperatureInCelsius": optionalFloatRangeSchema(),
		"StorageDurationInSeconds":    optionalUint32RangeSchema(),
		"ContainedInVesselIndex":      uintField("The index of the vessel this product is contained in, if any"),
		"Index":                       uintField("The display index/order of this product"),
		"IsWaste":                     boolField("Whether this product is waste"),
		"IsLiquid":                    boolField("Whether this product is a liquid"),
		"Compostable":                 boolField("Whether this product is compostable"),
	}),
	OutputSchema: schemaObject(recipeStepProductsSchema),
}

func (h *mcpToolManager) UpdateRecipeStepProduct() mcp.ToolHandlerFor[*UpdateRecipeStepProductInvocation, *mealplanning.RecipeStepProduct] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateRecipeStepProductInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepProduct, error) {
		result, err := h.client.UpdateRecipeStepProduct(ctx, &mealplanninggrpc.UpdateRecipeStepProductRequest{
			RecipeID:            x.RecipeID,
			RecipeStepID:        x.RecipeStepID,
			RecipeStepProductID: x.RecipeStepProductID,
			Input:               mealplanningconverters.ConvertRecipeStepProductUpdateRequestInputToGRPCRecipeStepProductUpdateRequestInput(x.RecipeStepProductUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepProductToRecipeStepProduct(result.Updated), nil
	}
}

//
