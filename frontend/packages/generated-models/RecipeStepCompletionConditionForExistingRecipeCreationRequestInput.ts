// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput } from './RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput';

export interface IRecipeStepCompletionConditionForExistingRecipeCreationRequestInput {
  optional: boolean;
  belongsToRecipeStep: string;
  ingredientStateID: string;
  ingredients: RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput;
  notes: string;
}

export class RecipeStepCompletionConditionForExistingRecipeCreationRequestInput
  implements IRecipeStepCompletionConditionForExistingRecipeCreationRequestInput
{
  optional: boolean;
  belongsToRecipeStep: string;
  ingredientStateID: string;
  ingredients: RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput;
  notes: string;
  constructor(input: Partial<RecipeStepCompletionConditionForExistingRecipeCreationRequestInput> = {}) {
    this.optional = input.optional = false;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.ingredientStateID = input.ingredientStateID = '';
    this.ingredients = input.ingredients =
      new RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput();
    this.notes = input.notes = '';
  }
}
