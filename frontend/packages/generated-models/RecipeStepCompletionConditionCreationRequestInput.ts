// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeStepCompletionConditionCreationRequestInput {
  notes: string;
  optional: boolean;
  belongsToRecipeStep: string;
  ingredientState: string;
  ingredients: number;
}

export class RecipeStepCompletionConditionCreationRequestInput
  implements IRecipeStepCompletionConditionCreationRequestInput
{
  notes: string;
  optional: boolean;
  belongsToRecipeStep: string;
  ingredientState: string;
  ingredients: number;
  constructor(input: Partial<RecipeStepCompletionConditionCreationRequestInput> = {}) {
    this.notes = input.notes = '';
    this.optional = input.optional = false;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.ingredientState = input.ingredientState = '';
    this.ingredients = input.ingredients = 0;
  }
}
