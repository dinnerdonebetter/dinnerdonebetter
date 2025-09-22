# Recipe Object Documentation

This document describes all the fields of the `Recipe` object and their purposes within the meal planning system.

## If you're new here

The `Recipe` object is the central entity in the meal planning system. It represents a complete cooking recipe that can be used for meal planning, grocery list generation, and cooking guidance. Key concepts to understand:

- **Recipe vs. Ingredient Distinction**: Recipes are multi-step cooking processes, while ingredients are raw materials. The system prevents single-step "recipes" (like "diced onion") but allows proper preparation recipes (like "caramelized onions").
- **Meal Component Types**: Recipes are categorized by what type of meal component they produce (appetizer, main course, dessert, etc.) to help with meal planning.
- **Recipe Relationships**: Recipes can be cloned/inspired by other recipes, maintaining provenance while allowing customization.
- **Scaling Support**: The system supports dynamic scaling of recipes through the frontend interface.
- **Preparation vs. Cooking**: Recipes separate prep tasks (advance preparation) from cooking steps (immediate cooking process).

## URL and Routing Fields

### `Slug` (string)
- **Purpose**: URL-friendly version of the recipe name
- **Usage**: Intended for SEO-friendly routing and frontend interfaces
- **Example**: "chocolate-chip-cookies", "beef-stir-fry"
- **Note**: Currently not implemented - no way to fetch recipes via slug. This field is subject to either enhancement or deprecation, but not really acceptable as-is.

## Metadata Fields

### `Description` (string)
- **Purpose**: Detailed description of the recipe
- **Usage**: Provides context about the dish, its origins, or special characteristics
- **Example**: "A classic American cookie recipe with crispy edges and chewy centers"

### `Source` (string)
- **Purpose**: Attribution or origin of the recipe
- **Usage**: Credits the original creator, cookbook, or website
- **Example**: "Grandma's Cookbook", "Food Network", "AllRecipes.com"

### `CreatedByUser` (string)
- **Purpose**: ID of the user who created this recipe
- **Usage**: Used for ownership tracking and user-specific recipe management

## Timing and Lifecycle Fields

### `CreatedAt` (time.Time)
- **Purpose**: Timestamp when the recipe was first created
- **Usage**: Audit trail, sorting, and display purposes

### `LastUpdatedAt` (*time.Time)
- **Purpose**: Timestamp of the most recent modification
- **Usage**: Tracks when recipe was last edited, can be null if never updated

### `ArchivedAt` (*time.Time)
- **Purpose**: Timestamp when the recipe was archived/deleted
- **Usage**: Soft delete functionality, can be null if recipe is active

## Portion and Yield Fields

### `EstimatedPortions` ([Float32RangeWithOptionalMax](../../internal/platform/types/main.go))
- **Purpose**: Range of how many portions/servings the recipe yields
- **Structure**: 
  - `Min` (float32): Minimum number of portions (required)
  - `Max` (*float32): Maximum number of portions (optional)
- **Usage**: Helps users understand recipe scale and plan for different group sizes
- **Logic**: If `Max` is null, serves exactly `Min` people. If `Max` is present, serves `Min` to `Max` people
- **Example**: Min: 4, Max: 6 (serves 4-6 people) or Min: 4, Max: null (serves exactly 4 people)

### `PortionName` (string)
- **Purpose**: Singular name for what the recipe yields
- **Usage**: Used in UI to describe individual servings
- **Example**: "serving", "cookie", "slice", "bowl"

### `PluralPortionName` (string)
- **Purpose**: Plural name for what the recipe yields
- **Usage**: Used in UI when referring to multiple servings
- **Example**: "servings", "cookies", "slices", "bowls"

## Recipe Classification Fields

### `YieldsComponentType` (string)
- **Purpose**: Categorizes what type of meal component this recipe produces
- **Valid Values**:
  - `"unspecified"` - No specific category
  - `"amuse-bouche"` - Small appetizer
  - `"appetizer"` - Starter course
  - `"soup"` - Soup course
  - `"main"` - Main course
  - `"salad"` - Salad course
  - `"beverage"` - Drink
  - `"side"` - Side dish
  - `"dessert"` - Dessert course
- **Usage**: Helps organize recipes in meal planning and determines how they can be used in meals
- **Note**: These values are defined as constants in [meal.go](../internal/domain/mealplanning/meal.go) and are used throughout the meal planning system

### `EligibleForMeals` (bool)
- **Purpose**: Whether this recipe can be included in new meals
- **Use Case**: Intended as a soft deletion mechanism for recipes currently in use. Allows admins to prevent new usage of a recipe while keeping existing meals functional, giving time to replace the recipe in existing meal plans
- **Note**: Recipes for individual ingredients should be marked as ineligible for meals. **This functionality is largely unimplemented** - the system doesn't currently enforce this flag when creating new meals

## Quality and Approval Fields

### `SealOfApproval` (bool)
- **Purpose**: Indicates if the recipe has been reviewed and approved by service operators
- **Usage**: Intended for highlighting "winner" recipes to new users
- **Note**: Currently unused - this field was added prematurely for future functionality

## Recipe Inspiration

### `InspiredByRecipeID` (*string)
- **Purpose**: References another recipe that inspired this one
- **Usage**: Enables easy recipe cloning and modification while maintaining provenance
- **Note**: If null or blank, the recipe is considered "original". Used when users want to make minor adjustments to existing recipes (e.g., reducing garlic in a soup recipe)

## Content Fields

### `Steps` ([]*[RecipeStep](../internal/domain/mealplanning/recipe_step.go))
- **Purpose**: Ordered list of cooking steps for the recipe
- **Usage**: Contains the actual cooking instructions and process
- **Note**: Must have at least 2 steps to prevent overly granular "recipes" (like "diced onion") while allowing proper prepared ingredient recipes (like "caramelized onions")
- **Step Requirements**: Each step must have at least one instrument OR vessel (not necessarily both)
  - **Vessels**: Identify what the preparation happens in (cutting board for chopping, pot for boiling)
  - **Instruments**: Tools used for the preparation (knife for chopping, spoon for stirring)
  - **Special Cases**: Some steps may only have instruments (cleaning a knife) or only vessels (preheating a pan)

### `PrepTasks` ([]*[RecipePrepTask](../internal/domain/mealplanning/recipe_prep_tasks.go))
- **Purpose**: Preparation tasks that can be done ahead of time
- **Usage**: Helps with meal planning by identifying advance preparation opportunities
- **Note**: These are separate from the main cooking steps

### `Media` ([]*[RecipeMedia](../internal/domain/mealplanning/recipe_media.go))
- **Purpose**: Images, videos, or other media associated with the recipe
- **Usage**: Can be attached to both the recipe level and individual steps
- **Note**: Currently a nascent concept. Planned to abstract to a general `UploadedMedia` type supporting images and videos. Recipe-level media might include finished product photos or full video tutorials, while step-level media provides preparation guidance

## Recipe Scaling

The system supports dynamic recipe scaling through the frontend interface. Users can adjust the scale of a recipe (e.g., 2x, 0.5x) which automatically multiplies ingredient quantities in all steps. For example, if a step calls for 1 clove of garlic and the recipe is scaled to 2x, it will display 2 cloves of garlic.

**Future Enhancement**: The system may support step-level scaling modifiers, allowing certain ingredients to scale at different rates (e.g., garlic scaling at 75% of the overall scale).

## Search and Discovery

Recipe search is currently implemented through a `RecipeSearchSubset` type that provides name-based searching. This is a basic implementation that may be enhanced in the future.

## Additional Notes

- **Recipe Relationships**: Recipes can be assigned as components into meals, but meal planning functionality is documented separately
- **Portion Ranges**: The min/max portion ranges serve as guidelines rather than strict requirements, reflecting the belief that cooking is inherently imprecise
- **Step Validation**: Each recipe step has its own validation logic enforced through typing (e.g., Preparation field is required)

## Common Patterns and Usage

### Creating a New Recipe
1. Set basic metadata (name, description, source)
2. Define portion information (estimated portions, portion names)
3. Set meal component type and eligibility
4. Add at least 2 cooking steps with ingredients and instruments and/or vessels
5. Optionally add prep tasks for advance preparation
6. Add media attachments if available

### Recipe Cloning Workflow
1. User finds a recipe they like
2. System creates a copy with `InspiredByRecipeID` pointing to original
3. `CreatedByUser` field is changed to the new user
4. All non-foreign-key IDs are regenerated (step IDs, product IDs, etc.)
5. Preparation IDs remain the same (they reference shared preparation methods)
6. User can modify any aspect of the cloned recipe
7. Original recipe remains unchanged

### Validation Rules
- **Minimum Steps**: Recipes must have at least 2 steps
- **Required Fields**: Name, slug, estimated portions, portion names, and yields component type are required
- **Step Requirements**: Each step must have at least one instrument OR vessel (not necessarily both), and a preparation method
- **Component Type**: Must be one of the predefined meal component types

### Data Flow
- Recipes are created through the API using [`RecipeCreationRequestInput`](../internal/domain/mealplanning/recipe.go)
- Database operations use [`RecipeDatabaseCreationInput`](../internal/domain/mealplanning/recipe.go)
- Updates use [`RecipeUpdateRequestInput`](../internal/domain/mealplanning/recipe.go)
- The system supports soft deletes via the `ArchivedAt` field

## Common Gotchas and Edge Cases

### Recipe vs. Ingredient Confusion
- **Problem**: New developers might try to create single-step "recipes" for basic ingredients
- **Solution**: The 2-step minimum prevents this. Use the ingredient system for raw materials
- **Example**: Don't create a recipe for "1 cup flour" - that's an ingredient. But "caramelized onions" is a valid recipe

### Slug Implementation Status
- **Problem**: Slug field exists but isn't functional
- **Solution**: Don't rely on slug-based routing until it's implemented
- **Workaround**: Use recipe ID for all routing needs
