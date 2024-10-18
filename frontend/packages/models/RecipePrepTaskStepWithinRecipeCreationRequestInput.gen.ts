// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipePrepTaskStepWithinRecipeCreationRequestInput {
  belongsToRecipeStepIndex: number;
  satisfiesRecipeStep: boolean;
}

export class RecipePrepTaskStepWithinRecipeCreationRequestInput
  implements IRecipePrepTaskStepWithinRecipeCreationRequestInput
{
  belongsToRecipeStepIndex: number;
  satisfiesRecipeStep: boolean;
  constructor(input: Partial<RecipePrepTaskStepWithinRecipeCreationRequestInput> = {}) {
    this.belongsToRecipeStepIndex = input.belongsToRecipeStepIndex || 0;
    this.satisfiesRecipeStep = input.satisfiesRecipeStep || false;
  }
}
