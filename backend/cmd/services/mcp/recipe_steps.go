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
	GetRecipeStepInvocation struct {
		RecipeID     string `jsonschema:"description=The recipe ID"`
		RecipeStepID string `jsonschema:"description=The recipe step ID"`
	}
)

var recipeMediaSchema = map[string]any{
	"ID":                  stringField("The ID of the recipe media"),
	"CreatedAt":           timestampField("When the recipe media was created"),
	"LastUpdatedAt":       timestampField("When the recipe media was last updated"),
	"ArchivedAt":          timestampField("When the recipe media was soft deleted"),
	"BelongsToRecipe":     stringField("The ID of the recipe this media belongs to, if any"),
	"BelongsToRecipeStep": stringField("The ID of the recipe step this media belongs to, if any"),
	"MimeType":            stringField("The MIME type of the media"),
	"InternalPath":        stringField("The internal path to the media file"),
	"ExternalPath":        stringField("The external path to the media file"),
	"Index":               uintField("The index of the media"),
}

var recipeStepsSchema = map[string]any{
	"ID":                    stringField("The ID of the recipe step"),
	"CreatedAt":             timestampField("When the recipe step was created"),
	"LastUpdatedAt":         timestampField("When the recipe step was last updated"),
	"ArchivedAt":            timestampField("When the recipe step was soft deleted"),
	"BelongsToRecipe":       stringField("The ID of the recipe this step belongs to"),
	"ConditionExpression":   stringField("The condition expression for this step"),
	"Notes":                 stringField("Notes about the step"),
	"ExplicitInstructions":  stringField("Explicit instructions for this step"),
	"Media":                 arrayType(recipeMediaSchema),
	"Products":              arrayType(recipeStepProductsSchema),
	"Instruments":           arrayType(recipeStepInstrumentsSchema),
	"Vessels":               arrayType(recipeStepVesselsSchema),
	"CompletionConditions": arrayType(recipeStepCompletionConditionsSchema),
	"Ingredients":          arrayType(recipeStepIngredientsSchema),
	"Preparation":           objectType(validPreparationsSchema),
	"Index":                 uintField("The index of the step within the recipe"),
	"Optional":              boolField("Whether this step is optional"),
	"StartTimerAutomatically": boolField("Whether to start a timer automatically for this step"),
	"EstimatedTimeInSeconds":  optionalUint32RangeSchema(),
	"TemperatureInCelsius":    optionalFloat32RangeSchema(),
}

var getRecipeStepTool = &mcp.Tool{
	Name:        "GetRecipeStep",
	Description: "Get a recipe step by it's ID",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":     stringField("The ID of the recipe"),
		"RecipeStepID": stringField("The ID of the recipe step to get"),
	}),
	OutputSchema: schemaObject(recipeStepsSchema),
}

func (h *mcpToolManager) GetRecipeStep() mcp.ToolHandlerFor[*GetRecipeStepInvocation, *mealplanning.RecipeStep] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStep, error) {
		result, err := h.client.GetRecipeStep(ctx, &mealplanninggrpc.GetRecipeStepRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepToRecipeStep(result.Result), nil
	}
}

type (
	GetRecipeStepsInvocation struct {
		RecipeID string
		Filter   *filtering.QueryFilter
	}

	GetRecipeStepsResult struct {
		Results []*mealplanning.RecipeStep
	}
)

var getRecipeStepsTool = &mcp.Tool{
	Name:        "GetRecipeSteps",
	Description: "Get recipe steps with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"RecipeID": stringField("The ID of the recipe"),
		"Filter":    queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(recipeStepsSchema),
	}),
}

func (h *mcpToolManager) GetRecipeSteps() mcp.ToolHandlerFor[*GetRecipeStepsInvocation, *GetRecipeStepsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepsInvocation) (*mcp.CallToolResult, *GetRecipeStepsResult, error) {
		results, err := h.client.GetRecipeSteps(ctx, &mealplanninggrpc.GetRecipeStepsRequest{
			RecipeID: x.RecipeID,
			Filter:   grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetRecipeStepsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCRecipeStepToRecipeStep(result))
		}

		return nil, out, nil
	}
}

type (
	CreateRecipeStepInvocation struct {
		*mealplanning.RecipeStepCreationRequestInput
		RecipeID string `jsonschema:"required,description=The recipe ID"`
	}
)

var recipeStepCreationTool = &mcp.Tool{
	Name:        "CreateRecipeStep",
	Description: "Create a recipe step",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                stringField("The ID of the recipe"),
		"PreparationID":           stringField("The ID of the preparation for this step"),
		"Notes":                   stringField("Notes about the step"),
		"ConditionExpression":     stringField("The condition expression for this step"),
		"ExplicitInstructions":    stringField("Explicit instructions for this step"),
		"Index":                   uintField("The index of the step within the recipe"),
		"Optional":                boolField("Whether this step is optional"),
		"StartTimerAutomatically": boolField("Whether to start a timer automatically for this step"),
		"EstimatedTimeInSeconds":   optionalUint32RangeSchema(),
		"TemperatureInCelsius":    optionalFloat32RangeSchema(),
		"Instruments":             arrayType(objectType(map[string]any{
			"InstrumentID":        stringField("The ID of the instrument"),
			"RecipeStepProductID": stringField("The ID of the recipe step product this instrument is associated with, if any"),
			"Quantity":            uint32RangeWithOptionalMaxSchema(),
			"Name":                 stringField("Name of the instrument"),
			"Notes":                stringField("Notes about the instrument"),
			"OptionIndex":         uintField("The option index for this instrument"),
			"PreferenceRank":      uintField("The preference rank for this instrument (0-255)"),
			"Optional":             boolField("Whether this instrument is optional"),
		})),
		"Vessels": arrayType(objectType(map[string]any{
			"VesselID":            stringField("The ID of the vessel"),
			"RecipeStepProductID": stringField("The ID of the recipe step product this vessel is associated with, if any"),
			"Quantity":             uint16RangeWithOptionalMaxSchema(),
			"Name":                 stringField("Name of the vessel"),
			"Notes":                stringField("Notes about the vessel"),
			"VesselPreposition":   stringField("The preposition to use with the vessel (e.g., 'in', 'on', 'over')"),
			"UnavailableAfterStep": boolField("Whether this vessel becomes unavailable after this step"),
		})),
		"Products": arrayType(objectType(map[string]any{
			"Name":                      stringField("Name of the product"),
			"Type":                      stringField("Type of the product"),
			"MeasurementUnitID":         stringField("The ID of the measurement unit"),
			"QuantityNotes":             stringField("Notes about the quantity"),
			"StorageInstructions":       stringField("Storage instructions for the product"),
			"Quantity":                  optionalFloat32RangeSchema(),
			"StorageTemperatureInCelsius": optionalFloat32RangeSchema(),
			"StorageDurationInSeconds":     optionalUint32RangeSchema(),
			"ContainedInVesselIndex":       uintField("The index of the vessel this product is contained in, if any"),
			"Index":                       uintField("The index of the product"),
			"IsWaste":                    boolField("Whether this product is waste"),
			"IsLiquid":                   boolField("Whether this product is a liquid"),
			"Compostable":                boolField("Whether this product is compostable"),
		})),
		"Ingredients": arrayType(objectType(map[string]any{
			"IngredientID":          stringField("The ID of the ingredient"),
			"MeasurementUnitID":     stringField("The ID of the measurement unit"),
			"Quantity":               float32RangeWithOptionalMaxSchema(),
			"Name":                   stringField("Name of the ingredient"),
			"QuantityNotes":         stringField("Notes about the quantity"),
			"IngredientNotes":       stringField("Notes about the ingredient"),
			"RecipeStepProductID":   stringField("The ID of the recipe step product this ingredient is associated with, if any"),
			"ProductOfRecipeID":     stringField("The ID of the recipe that produces this ingredient, if any"),
			"ProductPercentageToUse":  floatField("The percentage of the product to use, if any"),
			"VesselIndex":            uintField("The index of the vessel this ingredient is in, if any"),
			"OptionIndex":            uintField("The option index for this ingredient"),
			"Optional":               boolField("Whether this ingredient is optional"),
			"ToTaste":                boolField("Whether this ingredient is 'to taste'"),
		})),
		"CompletionConditions": arrayType(objectType(map[string]any{
			"IngredientStateID":   stringField("The ID of the ingredient state"),
			"BelongsToRecipeStep": stringField("The ID of the recipe step this completion condition belongs to"),
			"Notes":                stringField("Notes about the completion condition"),
			"Ingredients":          arrayType(objectType(map[string]any{
				"RecipeStepIngredient": uintField("The index of the recipe step ingredient"),
			})),
			"Optional": boolField("Whether this completion condition is optional"),
		})),
	}),
	OutputSchema: schemaObject(recipeStepsSchema),
}

func (h *mcpToolManager) CreateRecipeStep() mcp.ToolHandlerFor[*CreateRecipeStepInvocation, *mealplanning.RecipeStep] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateRecipeStepInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStep, error) {
		result, err := h.client.CreateRecipeStep(ctx, &mealplanninggrpc.CreateRecipeStepRequest{
			RecipeID: x.RecipeID,
			Input:    mealplanningconverters.ConvertRecipeStepCreationRequestInputToGRPCRecipeStepCreationRequestInput(x.RecipeStepCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepToRecipeStep(result.Created), nil
	}
}

type (
	UpdateRecipeStepInvocation struct {
		*mealplanning.RecipeStepUpdateRequestInput
		RecipeID     string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID string `jsonschema:"required,description=The recipe step ID"`
	}
)

var recipeStepUpdateTool = &mcp.Tool{
	Name:        "UpdateRecipeStep",
	Description: "Update a recipe step",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                stringField("The ID of the recipe"),
		"RecipeStepID":            stringField("The ID of the recipe step to update"),
		"PreparationID":           stringField("The ID of the preparation for this step"),
		"Notes":                   stringField("Notes about the step"),
		"ConditionExpression":     stringField("The condition expression for this step"),
		"ExplicitInstructions":    stringField("Explicit instructions for this step"),
		"Index":                   uintField("The index of the step within the recipe"),
		"Optional":                boolField("Whether this step is optional"),
		"StartTimerAutomatically": boolField("Whether to start a timer automatically for this step"),
		"EstimatedTimeInSeconds":  optionalUint32RangeSchema(),
		"TemperatureInCelsius":    optionalFloat32RangeSchema(),
		"BelongsToRecipe":         stringField("The ID of the recipe this step belongs to"),
	}),
	OutputSchema: schemaObject(recipeStepsSchema),
}

func (h *mcpToolManager) UpdateRecipeStep() mcp.ToolHandlerFor[*UpdateRecipeStepInvocation, *mealplanning.RecipeStep] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateRecipeStepInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStep, error) {
		result, err := h.client.UpdateRecipeStep(ctx, &mealplanninggrpc.UpdateRecipeStepRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
			Input:        mealplanningconverters.ConvertRecipeStepUpdateRequestInputToGRPCRecipeStepUpdateRequestInput(x.RecipeStepUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepToRecipeStep(result.Updated), nil
	}
}

//
