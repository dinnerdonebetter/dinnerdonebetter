// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeStepCompletionConditionCreationRequestInput {
  ingredientState: string;
  ingredients: number;
  notes: string;
  optional: boolean;
  belongsToRecipeStep: string;
}

export class RecipeStepCompletionConditionCreationRequestInput
  implements IRecipeStepCompletionConditionCreationRequestInput
{
  ingredientState: string;
  ingredients: number;
  notes: string;
  optional: boolean;
  belongsToRecipeStep: string;
  constructor(input: Partial<RecipeStepCompletionConditionCreationRequestInput> = {}) {
    this.ingredientState = input.ingredientState = '';
    this.ingredients = input.ingredients = 0;
    this.notes = input.notes = '';
    this.optional = input.optional = false;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
  }
}
