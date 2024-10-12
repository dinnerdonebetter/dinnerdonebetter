// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IRecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput {
   recipeStepIngredient: string;

}

export class RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput implements IRecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput {
   recipeStepIngredient: string;
constructor(input: Partial<RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput> = {}) {
	 this.recipeStepIngredient = input.recipeStepIngredient || '';
}
}