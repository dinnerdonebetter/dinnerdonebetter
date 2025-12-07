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
	GetRecipePrepTaskInvocation struct {
		RecipeID         string `jsonschema:"description=The recipe ID"`
		RecipePrepTaskID string `jsonschema:"description=The recipe prep task ID"`
	}
)

var recipePrepTaskStepSchema = map[string]any{
	"ID":                      stringField("The ID of the recipe prep task step"),
	"BelongsToRecipeStep":     stringField("The ID of the recipe step this prep task step belongs to"),
	"BelongsToRecipePrepTask": stringField("The ID of the recipe prep task this step belongs to"),
	"SatisfiesRecipeStep":     boolField("Whether this prep task step satisfies the recipe step"),
}

var recipePrepTasksSchema = map[string]any{
	"ID":                              stringField("The ID of the recipe prep task"),
	"CreatedAt":                       timestampField("When the recipe prep task was created"),
	"LastUpdatedAt":                   timestampField("When the recipe prep task was last updated"),
	"ArchivedAt":                      timestampField("When the recipe prep task was soft deleted"),
	"BelongsToRecipe":                 stringField("The ID of the recipe this prep task belongs to"),
	"Name":                            stringField("Name of the prep task"),
	"Description":                     stringField("Description of the prep task"),
	"Notes":                           stringField("Notes about the prep task"),
	"StorageType":                     stringField("The storage type for the prep task (e.g., 'covered', 'uncovered', 'on a wire rack')"),
	"ExplicitStorageInstructions":     stringField("Explicit storage instructions for the prep task"),
	"StorageTemperatureInCelsius":     optionalFloatRangeSchema(),
	"TimeBufferBeforeRecipeInSeconds": uint32RangeWithOptionalMaxSchema(),
	"Optional":                        boolField("Whether this prep task is optional"),
	"TaskSteps":                       arrayType(schemaObject(recipePrepTaskStepSchema)),
}

var getRecipePrepTaskTool = &mcp.Tool{
	Name:        "GetRecipePrepTask",
	Description: "Get a recipe prep task by it's ID",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":         stringField("The ID of the recipe"),
		"RecipePrepTaskID": stringField("The ID of the recipe prep task to get"),
	}),
	OutputSchema: schemaObject(recipePrepTasksSchema),
}

func (h *mcpToolManager) GetRecipePrepTask() mcp.ToolHandlerFor[*GetRecipePrepTaskInvocation, *mealplanning.RecipePrepTask] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipePrepTaskInvocation) (*mcp.CallToolResult, *mealplanning.RecipePrepTask, error) {
		result, err := h.client.GetRecipePrepTask(ctx, &mealplanninggrpc.GetRecipePrepTaskRequest{
			RecipeID:         x.RecipeID,
			RecipePrepTaskID: x.RecipePrepTaskID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipePrepTaskToRecipePrepTask(result.Result), nil
	}
}

type (
	GetRecipePrepTasksInvocation struct {
		Filter   *filtering.QueryFilter
		RecipeID string
	}

	GetRecipePrepTasksResult struct {
		Results []*mealplanning.RecipePrepTask
	}
)

var getRecipePrepTasksTool = &mcp.Tool{
	Name:        "GetRecipePrepTasks",
	Description: "Get recipe prep tasks with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"RecipeID": stringField("The ID of the recipe"),
		"Filter":   queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(recipePrepTasksSchema)),
	}),
}

func (h *mcpToolManager) GetRecipePrepTasks() mcp.ToolHandlerFor[*GetRecipePrepTasksInvocation, *GetRecipePrepTasksResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipePrepTasksInvocation) (*mcp.CallToolResult, *GetRecipePrepTasksResult, error) {
		results, err := h.client.GetRecipePrepTasks(ctx, &mealplanninggrpc.GetRecipePrepTasksRequest{
			RecipeID: x.RecipeID,
			Filter:   grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetRecipePrepTasksResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCRecipePrepTaskToRecipePrepTask(result))
		}

		return nil, out, nil
	}
}

type (
	CreateRecipePrepTaskInvocation struct {
		*mealplanning.RecipePrepTaskCreationRequestInput
		RecipeID string `jsonschema:"required,description=The recipe ID"`
	}
)

var recipePrepTaskCreationTool = &mcp.Tool{
	Name:        "CreateRecipePrepTask",
	Description: "Create a recipe prep task",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                        stringField("The ID of the recipe"),
		"Name":                            stringField("Name of the prep task"),
		"Description":                     stringField("Description of the prep task"),
		"Notes":                           stringField("Notes about the prep task"),
		"StorageType":                     stringField("The storage type for the prep task (e.g., 'covered', 'uncovered', 'on a wire rack')"),
		"ExplicitStorageInstructions":     stringField("Explicit storage instructions for the prep task"),
		"StorageTemperatureInCelsius":     optionalFloatRangeSchema(),
		"TimeBufferBeforeRecipeInSeconds": uint32RangeWithOptionalMaxSchema(),
		"Optional":                        boolField("Whether this prep task is optional"),
		"RecipeSteps":                     arrayType(schemaObject(recipePrepTaskStepSchema)),
	}),
	OutputSchema: schemaObject(recipePrepTasksSchema),
}

func (h *mcpToolManager) CreateRecipePrepTask() mcp.ToolHandlerFor[*CreateRecipePrepTaskInvocation, *mealplanning.RecipePrepTask] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateRecipePrepTaskInvocation) (*mcp.CallToolResult, *mealplanning.RecipePrepTask, error) {
		result, err := h.client.CreateRecipePrepTask(ctx, &mealplanninggrpc.CreateRecipePrepTaskRequest{
			RecipeID: x.RecipeID,
			Input:    mealplanningconverters.ConvertRecipePrepTaskCreationRequestInputToGRPCRecipePrepTaskCreationRequestInput(x.RecipePrepTaskCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipePrepTaskToRecipePrepTask(result.Created), nil
	}
}

type (
	UpdateRecipePrepTaskInvocation struct {
		*mealplanning.RecipePrepTaskUpdateRequestInput
		RecipeID         string `jsonschema:"required,description=The recipe ID"`
		RecipePrepTaskID string `jsonschema:"required,description=The recipe prep task ID"`
	}
)

var recipePrepTaskUpdateTool = &mcp.Tool{
	Name:        "UpdateRecipePrepTask",
	Description: "Update a recipe prep task",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                        stringField("The ID of the recipe"),
		"RecipePrepTaskID":                stringField("The ID of the recipe prep task to update"),
		"Name":                            stringField("Name of the prep task"),
		"Description":                     stringField("Description of the prep task"),
		"Notes":                           stringField("Notes about the prep task"),
		"StorageType":                     stringField("The storage type for the prep task (e.g., 'covered', 'uncovered', 'on a wire rack')"),
		"ExplicitStorageInstructions":     stringField("Explicit storage instructions for the prep task"),
		"StorageTemperatureInCelsius":     optionalFloatRangeSchema(),
		"TimeBufferBeforeRecipeInSeconds": uint32RangeWithOptionalMaxSchema(),
		"Optional":                        boolField("Whether this prep task is optional"),
		"RecipeSteps":                     arrayType(schemaObject(recipePrepTaskStepSchema)),
	}),
	OutputSchema: schemaObject(recipePrepTasksSchema),
}

func (h *mcpToolManager) UpdateRecipePrepTask() mcp.ToolHandlerFor[*UpdateRecipePrepTaskInvocation, *mealplanning.RecipePrepTask] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateRecipePrepTaskInvocation) (*mcp.CallToolResult, *mealplanning.RecipePrepTask, error) {
		result, err := h.client.UpdateRecipePrepTask(ctx, &mealplanninggrpc.UpdateRecipePrepTaskRequest{
			RecipeID:         x.RecipeID,
			RecipePrepTaskID: x.RecipePrepTaskID,
			Input:            mealplanningconverters.ConvertRecipePrepTaskUpdateRequestInputToGRPCRecipePrepTaskUpdateRequestInput(x.RecipePrepTaskUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipePrepTaskToRecipePrepTask(result.Updated), nil
	}
}

//
