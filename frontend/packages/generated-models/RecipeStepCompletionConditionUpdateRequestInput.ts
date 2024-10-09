// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeStepCompletionConditionUpdateRequestInput {
  belongsToRecipeStep?: string;
  ingredientState?: string;
  notes?: string;
  optional?: boolean;
}

export class RecipeStepCompletionConditionUpdateRequestInput
  implements IRecipeStepCompletionConditionUpdateRequestInput
{
  belongsToRecipeStep?: string;
  ingredientState?: string;
  notes?: string;
  optional?: boolean;
  constructor(input: Partial<RecipeStepCompletionConditionUpdateRequestInput> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.ingredientState = input.ingredientState;
    this.notes = input.notes;
    this.optional = input.optional;
  }
}
