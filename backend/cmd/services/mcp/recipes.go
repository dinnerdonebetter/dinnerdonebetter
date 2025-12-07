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
	GetRecipeInvocation struct {
		RecipeID string `jsonschema:"description=The recipe ID"`
	}
)

var recipesSchema = map[string]any{
	"ID":                  stringField("The ID of the recipe"),
	"CreatedAt":           timestampField("When the recipe was created"),
	"LastUpdatedAt":       timestampField("When the recipe was last updated"),
	"ArchivedAt":          timestampField("When the recipe was soft deleted"),
	"InspiredByRecipeID":  stringField("The ID of the recipe this recipe was inspired by, if any"),
	"Name":                stringField("Name of the recipe"),
	"Description":         stringField("Description of the recipe"),
	"Source":              stringField("Source of the recipe"),
	"Slug":                stringField("An easy-to-use URL slug for the recipe"),
	"CreatedByUser":       stringField("The ID of the user who created the recipe"),
	"PortionName":         stringField("Name for a single portion (e.g., 'serving', 'piece')"),
	"PluralPortionName":   stringField("Plural name for portions (e.g., 'servings', 'pieces')"),
	"YieldsComponentType": stringField("The type of component this recipe yields"),
	"EstimatedPortions":   float32RangeWithOptionalMaxSchema(),
	"SealOfApproval":      boolField("Whether this recipe has the seal of approval"),
	"EligibleForMeals":    boolField("Whether this recipe is eligible for meals"),
	"PrepTasks":           arrayType(schemaObject(recipePrepTasksSchema)),
	"Steps":               arrayType(schemaObject(recipeStepsSchema)),
	"Media":               arrayType(schemaObject(recipeMediaSchema)),
}

var getRecipeTool = &mcp.Tool{
	Name:        "GetRecipe",
	Description: "Get a recipe by it's ID",
	InputSchema: schemaObject(map[string]any{
		"RecipeID": stringField("The ID of the recipe to get"),
	}),
	OutputSchema: schemaObject(recipesSchema),
}

func (h *mcpToolManager) GetRecipe() mcp.ToolHandlerFor[*GetRecipeInvocation, *mealplanning.Recipe] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeInvocation) (*mcp.CallToolResult, *mealplanning.Recipe, error) {
		result, err := h.client.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{
			RecipeID: x.RecipeID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeToRecipe(result.Result), nil
	}
}

type (
	GetRecipesInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetRecipesResult struct {
		Results []*mealplanning.Recipe
	}
)

var getRecipesTool = &mcp.Tool{
	Name:        "GetRecipes",
	Description: "Get recipes with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(recipesSchema)),
	}),
}

func (h *mcpToolManager) GetRecipes() mcp.ToolHandlerFor[*GetRecipesInvocation, *GetRecipesResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipesInvocation) (*mcp.CallToolResult, *GetRecipesResult, error) {
		results, err := h.client.GetRecipes(ctx, &mealplanninggrpc.GetRecipesRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetRecipesResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCRecipeToRecipe(result))
		}

		return nil, out, nil
	}
}

type (
	SearchForRecipesInvocation struct {
		Query            string
		UseSearchService bool
		Filter           *filtering.QueryFilter
	}

	SearchForRecipesResult struct {
		Results []*mealplanning.Recipe
	}
)

var searchForRecipesTool = &mcp.Tool{
	Name:        "SearchForRecipes",
	Description: "Search for recipes with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Query":            stringField("The search query string"),
		"UseSearchService": boolField("Whether to use the search service (if false, uses database search)"),
		"Filter":           queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(recipesSchema)),
	}),
}

func (h *mcpToolManager) SearchForRecipes() mcp.ToolHandlerFor[*SearchForRecipesInvocation, *SearchForRecipesResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchForRecipesInvocation) (*mcp.CallToolResult, *SearchForRecipesResult, error) {
		results, err := h.client.SearchForRecipes(ctx, &mealplanninggrpc.SearchForRecipesRequest{
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchForRecipesResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCRecipeToRecipe(result))
		}

		return nil, out, nil
	}
}

type (
	CreateRecipeInvocation struct {
		*mealplanning.RecipeCreationRequestInput
	}
)

var recipeCreationTool = &mcp.Tool{
	Name:        "CreateRecipe",
	Description: "Create a recipe",
	InputSchema: schemaObject(map[string]any{
		"InspiredByRecipeID":  stringField("The ID of the recipe this recipe was inspired by, if any"),
		"Name":                stringField("Name of the recipe"),
		"Source":              stringField("Source of the recipe"),
		"Description":         stringField("Description of the recipe"),
		"PluralPortionName":   stringField("Plural name for portions (e.g., 'servings', 'pieces')"),
		"PortionName":         stringField("Name for a single portion (e.g., 'serving', 'piece')"),
		"Slug":                stringField("An easy-to-use URL slug for the recipe"),
		"YieldsComponentType": stringField("The type of component this recipe yields"),
		"EstimatedPortions":   float32RangeWithOptionalMaxSchema(),
		"AlsoCreateMeal":      boolField("Whether to also create a meal from this recipe"),
		"SealOfApproval":      boolField("Whether this recipe has the seal of approval"),
		"EligibleForMeals":    boolField("Whether this recipe is eligible for meals"),
		"PrepTasks": arrayType(objectType(map[string]any{
			"Name":                            stringField("Name of the prep task"),
			"Description":                     stringField("Description of the prep task"),
			"Notes":                           stringField("Notes about the prep task"),
			"StorageType":                     stringField("The storage type for the prep task (e.g., 'covered', 'uncovered', 'on a wire rack')"),
			"ExplicitStorageInstructions":     stringField("Explicit storage instructions for the prep task"),
			"StorageTemperatureInCelsius":     optionalFloat32RangeSchema(),
			"TimeBufferBeforeRecipeInSeconds": uint32RangeWithOptionalMaxSchema(),
			"Optional":                        boolField("Whether this prep task is optional"),
			"TaskSteps": arrayType(objectType(map[string]any{
				"BelongsToRecipeStep": stringField("The ID of the recipe step this prep task step belongs to"),
				"SatisfiesRecipeStep": boolField("Whether this prep task step satisfies the recipe step"),
			})),
		})),
		"Steps": arrayType(objectType(map[string]any{
			"PreparationID":           stringField("The ID of the preparation for this step"),
			"Notes":                   stringField("Notes about the step"),
			"ConditionExpression":     stringField("The condition expression for this step"),
			"ExplicitInstructions":    stringField("Explicit instructions for this step"),
			"Index":                   uintField("The index of the step within the recipe"),
			"Optional":                boolField("Whether this step is optional"),
			"StartTimerAutomatically": boolField("Whether to start a timer automatically for this step"),
			"EstimatedTimeInSeconds":  optionalUint32RangeSchema(),
			"TemperatureInCelsius":    optionalFloat32RangeSchema(),
			"Instruments": arrayType(objectType(map[string]any{
				"InstrumentID":        stringField("The ID of the instrument"),
				"RecipeStepProductID": stringField("The ID of the recipe step product this instrument is associated with, if any"),
				"Quantity":            uint32RangeWithOptionalMaxSchema(),
				"Name":                stringField("Name of the instrument"),
				"Notes":               stringField("Notes about the instrument"),
				"OptionIndex":         uintField("The option index for this instrument"),
				"PreferenceRank":      uintField("The preference rank for this instrument (0-255)"),
				"Optional":            boolField("Whether this instrument is optional"),
			})),
			"Vessels": arrayType(objectType(map[string]any{
				"VesselID":             stringField("The ID of the vessel"),
				"RecipeStepProductID":  stringField("The ID of the recipe step product this vessel is associated with, if any"),
				"Quantity":             uint16RangeWithOptionalMaxSchema(),
				"Name":                 stringField("Name of the vessel"),
				"Notes":                stringField("Notes about the vessel"),
				"VesselPreposition":    stringField("The preposition to use with the vessel (e.g., 'in', 'on', 'over')"),
				"UnavailableAfterStep": boolField("Whether this vessel becomes unavailable after this step"),
			})),
			"Products": arrayType(objectType(map[string]any{
				"Name":                        stringField("Name of the product"),
				"Type":                        stringField("Type of the product"),
				"MeasurementUnitID":           stringField("The ID of the measurement unit"),
				"QuantityNotes":               stringField("Notes about the quantity"),
				"StorageInstructions":         stringField("Storage instructions for the product"),
				"Quantity":                    optionalFloat32RangeSchema(),
				"StorageTemperatureInCelsius": optionalFloat32RangeSchema(),
				"StorageDurationInSeconds":    optionalUint32RangeSchema(),
				"ContainedInVesselIndex":      uintField("The index of the vessel this product is contained in, if any"),
				"Index":                       uintField("The index of the product"),
				"IsWaste":                     boolField("Whether this product is waste"),
				"IsLiquid":                    boolField("Whether this product is a liquid"),
				"Compostable":                 boolField("Whether this product is compostable"),
			})),
			"Ingredients": arrayType(objectType(map[string]any{
				"IngredientID":           stringField("The ID of the ingredient"),
				"MeasurementUnitID":      stringField("The ID of the measurement unit"),
				"Quantity":               float32RangeWithOptionalMaxSchema(),
				"Name":                   stringField("Name of the ingredient"),
				"QuantityNotes":          stringField("Notes about the quantity"),
				"IngredientNotes":        stringField("Notes about the ingredient"),
				"RecipeStepProductID":    stringField("The ID of the recipe step product this ingredient is associated with, if any"),
				"ProductOfRecipeID":      stringField("The ID of the recipe that produces this ingredient, if any"),
				"ProductPercentageToUse": floatField("The percentage of the product to use, if any"),
				"VesselIndex":            uintField("The index of the vessel this ingredient is in, if any"),
				"OptionIndex":            uintField("The option index for this ingredient"),
				"Optional":               boolField("Whether this ingredient is optional"),
				"ToTaste":                boolField("Whether this ingredient is 'to taste'"),
			})),
			"CompletionConditions": arrayType(objectType(map[string]any{
				"IngredientStateID":   stringField("The ID of the ingredient state"),
				"BelongsToRecipeStep": stringField("The ID of the recipe step this completion condition belongs to"),
				"Notes":               stringField("Notes about the completion condition"),
				"Ingredients": arrayType(objectType(map[string]any{
					"RecipeStepIngredient": uintField("The index of the recipe step ingredient"),
				})),
				"Optional": boolField("Whether this completion condition is optional"),
			})),
		})),
		"Media": arrayType(objectType(map[string]any{
			"BelongsToRecipe":     stringField("The ID of the recipe this media belongs to, if any"),
			"BelongsToRecipeStep": stringField("The ID of the recipe step this media belongs to, if any"),
			"MimeType":            stringField("The MIME type of the media"),
			"InternalPath":        stringField("The internal path to the media file"),
			"ExternalPath":        stringField("The external path to the media file"),
			"Index":               uintField("The index of the media"),
		})),
	}),
	OutputSchema: schemaObject(recipesSchema),
}

func (h *mcpToolManager) CreateRecipe() mcp.ToolHandlerFor[*CreateRecipeInvocation, *mealplanning.Recipe] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateRecipeInvocation) (*mcp.CallToolResult, *mealplanning.Recipe, error) {
		result, err := h.client.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{
			Input: mealplanningconverters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(x.RecipeCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeToRecipe(result.Created), nil
	}
}

type (
	UpdateRecipeInvocation struct {
		*mealplanning.RecipeUpdateRequestInput
		RecipeID string `jsonschema:"required,description=The recipe ID"`
	}
)

var recipeUpdateTool = &mcp.Tool{
	Name:        "UpdateRecipe",
	Description: "Update a recipe",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":            stringField("The ID of the recipe to update"),
		"Name":                stringField("Name of the recipe"),
		"Source":              stringField("Source of the recipe"),
		"Description":         stringField("Description of the recipe"),
		"InspiredByRecipeID":  stringField("The ID of the recipe this recipe was inspired by, if any"),
		"Slug":                stringField("An easy-to-use URL slug for the recipe"),
		"PortionName":         stringField("Name for a single portion (e.g., 'serving', 'piece')"),
		"PluralPortionName":   stringField("Plural name for portions (e.g., 'servings', 'pieces')"),
		"YieldsComponentType": stringField("The type of component this recipe yields"),
		"EstimatedPortions":   float32RangeWithOptionalMaxSchema(),
		"SealOfApproval":      boolField("Whether this recipe has the seal of approval"),
		"EligibleForMeals":    boolField("Whether this recipe is eligible for meals"),
	}),
	OutputSchema: schemaObject(recipesSchema),
}

func (h *mcpToolManager) UpdateRecipe() mcp.ToolHandlerFor[*UpdateRecipeInvocation, *mealplanning.Recipe] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateRecipeInvocation) (*mcp.CallToolResult, *mealplanning.Recipe, error) {
		result, err := h.client.UpdateRecipe(ctx, &mealplanninggrpc.UpdateRecipeRequest{
			RecipeID: x.RecipeID,
			Input:    mealplanningconverters.ConvertRecipeUpdateRequestInputToGRPCRecipeUpdateRequestInput(x.RecipeUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeToRecipe(result.Updated), nil
	}
}

//
