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
	GetRecipeStepVesselInvocation struct {
		RecipeID           string `jsonschema:"description=The recipe ID"`
		RecipeStepID       string `jsonschema:"description=The recipe step ID"`
		RecipeStepVesselID string `jsonschema:"description=The recipe step vessel ID"`
	}
)

var recipeStepVesselsSchema = map[string]any{
	"ID":                   stringField("The ID of the recipe step vessel"),
	"CreatedAt":            timestampField("When the recipe step vessel was created"),
	"LastUpdatedAt":        timestampField("When the recipe step vessel was last updated"),
	"ArchivedAt":           timestampField("When the recipe step vessel was soft deleted"),
	"BelongsToRecipeStep":  stringField("The ID of the recipe step this vessel belongs to"),
	"Name":                 stringField("Name of the vessel"),
	"Notes":                stringField("Notes about the vessel"),
	"Vessel":               objectType(validVesselsSchema),
	"RecipeStepProductID":  stringField("The ID of the recipe step product this vessel is associated with, if any"),
	"Quantity":             uint16RangeWithOptionalMaxSchema(),
	"VesselPreposition":    stringField("The preposition to use with the vessel (e.g., 'in', 'on', 'over')"),
	"UnavailableAfterStep": boolField("Whether this vessel becomes unavailable after this step"),
}

var getRecipeStepVesselTool = &mcp.Tool{
	Name:        "GetRecipeStepVessel",
	Description: "Get a recipe step vessel by it's ID",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":           stringField("The ID of the recipe"),
		"RecipeStepID":       stringField("The ID of the recipe step"),
		"RecipeStepVesselID": stringField("The ID of the recipe step vessel to get"),
	}),
	OutputSchema: schemaObject(recipeStepVesselsSchema),
}

func (h *mcpToolManager) GetRecipeStepVessel() mcp.ToolHandlerFor[*GetRecipeStepVesselInvocation, *mealplanning.RecipeStepVessel] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepVesselInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepVessel, error) {
		result, err := h.client.GetRecipeStepVessel(ctx, &mealplanninggrpc.GetRecipeStepVesselRequest{
			RecipeID:           x.RecipeID,
			RecipeStepID:       x.RecipeStepID,
			RecipeStepVesselID: x.RecipeStepVesselID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepVesselToRecipeStepVessel(result.Result), nil
	}
}

type (
	GetRecipeStepVesselsInvocation struct {
		RecipeID     string
		RecipeStepID string
		Filter       *filtering.QueryFilter
	}

	GetRecipeStepVesselsResult struct {
		Results []*mealplanning.RecipeStepVessel
	}
)

var getRecipeStepVesselsTool = &mcp.Tool{
	Name:        "GetRecipeStepVessels",
	Description: "Get recipe step vessels with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":     stringField("The ID of the recipe"),
		"RecipeStepID": stringField("The ID of the recipe step"),
		"Filter":       queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(recipeStepVesselsSchema)),
	}),
}

func (h *mcpToolManager) GetRecipeStepVessels() mcp.ToolHandlerFor[*GetRecipeStepVesselsInvocation, *GetRecipeStepVesselsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepVesselsInvocation) (*mcp.CallToolResult, *GetRecipeStepVesselsResult, error) {
		results, err := h.client.GetRecipeStepVessels(ctx, &mealplanninggrpc.GetRecipeStepVesselsRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
			Filter:       grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetRecipeStepVesselsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCRecipeStepVesselToRecipeStepVessel(result))
		}

		return nil, out, nil
	}
}

type (
	CreateRecipeStepVesselInvocation struct {
		*mealplanning.RecipeStepVesselCreationRequestInput
		RecipeID     string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID string `jsonschema:"required,description=The recipe step ID"`
	}
)

var recipeStepVesselCreationTool = &mcp.Tool{
	Name:        "CreateRecipeStepVessel",
	Description: "Create a recipe step vessel",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                        stringField("The ID of the recipe"),
		"RecipeStepID":                    stringField("The ID of the recipe step"),
		"VesselID":                        stringField("The ID of the vessel"),
		"RecipeStepProductID":             stringField("The ID of the recipe step product this vessel is associated with, if any"),
		"ProductOfRecipeStepIndex":        uintField("The index of the recipe step that produces this vessel, if any"),
		"ProductOfRecipeStepProductIndex": uintField("The index of the recipe step product that produces this vessel, if any"),
		"Name":                            stringField("Name of the vessel"),
		"Notes":                           stringField("Notes about the vessel"),
		"Quantity":                        uint16RangeWithOptionalMaxSchema(),
		"VesselPreposition":               stringField("The preposition to use with the vessel (e.g., 'in', 'on', 'over')"),
		"UnavailableAfterStep":            boolField("Whether this vessel becomes unavailable after this step"),
	}),
	OutputSchema: schemaObject(recipeStepVesselsSchema),
}

func (h *mcpToolManager) CreateRecipeStepVessel() mcp.ToolHandlerFor[*CreateRecipeStepVesselInvocation, *mealplanning.RecipeStepVessel] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateRecipeStepVesselInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepVessel, error) {
		result, err := h.client.CreateRecipeStepVessel(ctx, &mealplanninggrpc.CreateRecipeStepVesselRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
			Input:        mealplanningconverters.ConvertRecipeStepVesselCreationRequestInputToGRPCRecipeStepVesselCreationRequestInput(x.RecipeStepVesselCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepVesselToRecipeStepVessel(result.Created), nil
	}
}

type (
	UpdateRecipeStepVesselInvocation struct {
		*mealplanning.RecipeStepVesselUpdateRequestInput
		RecipeID           string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID       string `jsonschema:"required,description=The recipe step ID"`
		RecipeStepVesselID string `jsonschema:"required,description=The recipe step vessel ID"`
	}
)

var recipeStepVesselUpdateTool = &mcp.Tool{
	Name:        "UpdateRecipeStepVessel",
	Description: "Update a recipe step vessel",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":             stringField("The ID of the recipe"),
		"RecipeStepID":         stringField("The ID of the recipe step"),
		"RecipeStepVesselID":   stringField("The ID of the recipe step vessel to update"),
		"VesselID":             stringField("The ID of the vessel"),
		"RecipeStepProductID":  stringField("The ID of the recipe step product this vessel is associated with, if any"),
		"Name":                 stringField("Name of the vessel"),
		"Notes":                stringField("Notes about the vessel"),
		"Quantity":             uint16RangeWithOptionalMaxSchema(),
		"VesselPreposition":    stringField("The preposition to use with the vessel (e.g., 'in', 'on', 'over')"),
		"UnavailableAfterStep": boolField("Whether this vessel becomes unavailable after this step"),
	}),
	OutputSchema: schemaObject(recipeStepVesselsSchema),
}

func (h *mcpToolManager) UpdateRecipeStepVessel() mcp.ToolHandlerFor[*UpdateRecipeStepVesselInvocation, *mealplanning.RecipeStepVessel] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateRecipeStepVesselInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepVessel, error) {
		result, err := h.client.UpdateRecipeStepVessel(ctx, &mealplanninggrpc.UpdateRecipeStepVesselRequest{
			RecipeID:           x.RecipeID,
			RecipeStepID:       x.RecipeStepID,
			RecipeStepVesselID: x.RecipeStepVesselID,
			Input:              mealplanningconverters.ConvertRecipeStepVesselUpdateRequestInputToGRPCRecipeStepVesselUpdateRequestInput(x.RecipeStepVesselUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepVesselToRecipeStepVessel(result.Updated), nil
	}
}

//
