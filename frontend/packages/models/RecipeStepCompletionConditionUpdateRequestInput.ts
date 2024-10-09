// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeStepCompletionConditionUpdateRequestInput {
  notes: string;
  optional: boolean;
  belongsToRecipeStep: string;
  ingredientState: string;
}

export class RecipeStepCompletionConditionUpdateRequestInput
  implements IRecipeStepCompletionConditionUpdateRequestInput
{
  notes: string;
  optional: boolean;
  belongsToRecipeStep: string;
  ingredientState: string;
  constructor(input: Partial<RecipeStepCompletionConditionUpdateRequestInput> = {}) {
    this.notes = input.notes || '';
    this.optional = input.optional || false;
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.ingredientState = input.ingredientState || '';
  }
}
