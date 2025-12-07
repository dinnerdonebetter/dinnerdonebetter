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
	GetRecipeStepCompletionConditionInvocation struct {
		RecipeID                        string `jsonschema:"description=The recipe ID"`
		RecipeStepID                    string `jsonschema:"description=The recipe step ID"`
		RecipeStepCompletionConditionID string `jsonschema:"description=The recipe step completion condition ID"`
	}
)

var recipeStepCompletionConditionIngredientSchema = map[string]any{
	"ID":                                     stringField("The ID of the recipe step completion condition ingredient"),
	"CreatedAt":                              timestampField("When the recipe step completion condition ingredient was created"),
	"LastUpdatedAt":                          timestampField("When the recipe step completion condition ingredient was last updated"),
	"ArchivedAt":                             timestampField("When the recipe step completion condition ingredient was soft deleted"),
	"BelongsToRecipeStepCompletionCondition": stringField("The ID of the recipe step completion condition this ingredient belongs to"),
	"RecipeStepIngredient":                   stringField("The ID of the recipe step ingredient"),
}

var recipeStepCompletionConditionsSchema = map[string]any{
	"ID":                  stringField("The ID of the recipe step completion condition"),
	"CreatedAt":           timestampField("When the recipe step completion condition was created"),
	"LastUpdatedAt":       timestampField("When the recipe step completion condition was last updated"),
	"ArchivedAt":          timestampField("When the recipe step completion condition was soft deleted"),
	"BelongsToRecipeStep": stringField("The ID of the recipe step this completion condition belongs to"),
	"IngredientState":     objectType(validIngredientStatesSchema),
	"Notes":               stringField("Notes about the completion condition"),
	"Ingredients":         arrayType(schemaObject(recipeStepCompletionConditionIngredientSchema)),
	"Optional":            boolField("Whether this completion condition is optional"),
}

var getRecipeStepCompletionConditionTool = &mcp.Tool{
	Name:        "GetRecipeStepCompletionCondition",
	Description: "Get a recipe step completion condition by it's ID",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                        stringField("The ID of the recipe"),
		"RecipeStepID":                    stringField("The ID of the recipe step"),
		"RecipeStepCompletionConditionID": stringField("The ID of the recipe step completion condition to get"),
	}),
	OutputSchema: schemaObject(recipeStepCompletionConditionsSchema),
}

func (h *mcpToolManager) GetRecipeStepCompletionCondition() mcp.ToolHandlerFor[*GetRecipeStepCompletionConditionInvocation, *mealplanning.RecipeStepCompletionCondition] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepCompletionConditionInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepCompletionCondition, error) {
		result, err := h.client.GetRecipeStepCompletionCondition(ctx, &mealplanninggrpc.GetRecipeStepCompletionConditionRequest{
			RecipeID:                        x.RecipeID,
			RecipeStepID:                    x.RecipeStepID,
			RecipeStepCompletionConditionID: x.RecipeStepCompletionConditionID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepCompletionConditionToRecipeStepCompletionCondition(result.Result), nil
	}
}

type (
	GetRecipeStepCompletionConditionsInvocation struct {
		RecipeID     string
		RecipeStepID string
		Filter       *filtering.QueryFilter
	}

	GetRecipeStepCompletionConditionsResult struct {
		Results []*mealplanning.RecipeStepCompletionCondition
	}
)

var getRecipeStepCompletionConditionsTool = &mcp.Tool{
	Name:        "GetRecipeStepCompletionConditions",
	Description: "Get recipe step completion conditions with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":     stringField("The ID of the recipe"),
		"RecipeStepID": stringField("The ID of the recipe step"),
		"Filter":       queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(recipeStepCompletionConditionsSchema)),
	}),
}

func (h *mcpToolManager) GetRecipeStepCompletionConditions() mcp.ToolHandlerFor[*GetRecipeStepCompletionConditionsInvocation, *GetRecipeStepCompletionConditionsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepCompletionConditionsInvocation) (*mcp.CallToolResult, *GetRecipeStepCompletionConditionsResult, error) {
		results, err := h.client.GetRecipeStepCompletionConditions(ctx, &mealplanninggrpc.GetRecipeStepCompletionConditionsRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
			Filter:       grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetRecipeStepCompletionConditionsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCRecipeStepCompletionConditionToRecipeStepCompletionCondition(result))
		}

		return nil, out, nil
	}
}

type (
	CreateRecipeStepCompletionConditionInvocation struct {
		*mealplanning.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput
		RecipeID     string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID string `jsonschema:"required,description=The recipe step ID"`
	}
)

var recipeStepCompletionConditionCreationTool = &mcp.Tool{
	Name:        "CreateRecipeStepCompletionCondition",
	Description: "Create a recipe step completion condition",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":            stringField("The ID of the recipe"),
		"RecipeStepID":        stringField("The ID of the recipe step"),
		"IngredientStateID":   stringField("The ID of the ingredient state"),
		"BelongsToRecipeStep": stringField("The ID of the recipe step this completion condition belongs to"),
		"Notes":               stringField("Notes about the completion condition"),
		"Ingredients": arrayType(objectType(map[string]any{
			"RecipeStepIngredient": stringField("The ID of the recipe step ingredient"),
		})),
		"Optional": boolField("Whether this completion condition is optional"),
	}),
	OutputSchema: schemaObject(recipeStepCompletionConditionsSchema),
}

func (h *mcpToolManager) CreateRecipeStepCompletionCondition() mcp.ToolHandlerFor[*CreateRecipeStepCompletionConditionInvocation, *mealplanning.RecipeStepCompletionCondition] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateRecipeStepCompletionConditionInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepCompletionCondition, error) {
		result, err := h.client.CreateRecipeStepCompletionCondition(ctx, &mealplanninggrpc.CreateRecipeStepCompletionConditionRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
			Input:        mealplanningconverters.ConvertRecipeStepCompletionConditionForExistingRecipeCreationRequestInputToGRPCRecipeStepCompletionConditionForExistingRecipeCreationRequestInput(x.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepCompletionConditionToRecipeStepCompletionCondition(result.Created), nil
	}
}

type (
	UpdateRecipeStepCompletionConditionInvocation struct {
		*mealplanning.RecipeStepCompletionConditionUpdateRequestInput
		RecipeID                        string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID                    string `jsonschema:"required,description=The recipe step ID"`
		RecipeStepCompletionConditionID string `jsonschema:"required,description=The recipe step completion condition ID"`
	}
)

var recipeStepCompletionConditionUpdateTool = &mcp.Tool{
	Name:        "UpdateRecipeStepCompletionCondition",
	Description: "Update a recipe step completion condition",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                        stringField("The ID of the recipe"),
		"RecipeStepID":                    stringField("The ID of the recipe step"),
		"RecipeStepCompletionConditionID": stringField("The ID of the recipe step completion condition to update"),
		"IngredientStateID":               stringField("The ID of the ingredient state"),
		"BelongsToRecipeStep":             stringField("The ID of the recipe step this completion condition belongs to"),
		"Notes":                           stringField("Notes about the completion condition"),
		"Optional":                        boolField("Whether this completion condition is optional"),
	}),
	OutputSchema: schemaObject(recipeStepCompletionConditionsSchema),
}

func (h *mcpToolManager) UpdateRecipeStepCompletionCondition() mcp.ToolHandlerFor[*UpdateRecipeStepCompletionConditionInvocation, *mealplanning.RecipeStepCompletionCondition] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateRecipeStepCompletionConditionInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepCompletionCondition, error) {
		result, err := h.client.UpdateRecipeStepCompletionCondition(ctx, &mealplanninggrpc.UpdateRecipeStepCompletionConditionRequest{
			RecipeID:                        x.RecipeID,
			RecipeStepID:                    x.RecipeStepID,
			RecipeStepCompletionConditionID: x.RecipeStepCompletionConditionID,
			Input:                           mealplanningconverters.ConvertRecipeStepCompletionConditionUpdateRequestInputToGRPCRecipeStepCompletionConditionUpdateRequestInput(x.RecipeStepCompletionConditionUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepCompletionConditionToRecipeStepCompletionCondition(result.Updated), nil
	}
}

//
