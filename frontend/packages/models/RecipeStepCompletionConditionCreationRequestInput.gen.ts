// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeStepCompletionConditionCreationRequestInput {
  belongsToRecipeStep: string;
  ingredientState: string;
  ingredients: number[];
  notes: string;
  optional: boolean;
}

export class RecipeStepCompletionConditionCreationRequestInput
  implements IRecipeStepCompletionConditionCreationRequestInput
{
  belongsToRecipeStep: string;
  ingredientState: string;
  ingredients: number[];
  notes: string;
  optional: boolean;
  constructor(input: Partial<RecipeStepCompletionConditionCreationRequestInput> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.ingredientState = input.ingredientState || '';
    this.ingredients = input.ingredients || [];
    this.notes = input.notes || '';
    this.optional = input.optional || false;
  }
}
