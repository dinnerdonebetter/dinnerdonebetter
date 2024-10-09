// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipePrepTaskStepCreationRequestInput {
  belongsToRecipeStep: string;
  satisfiesRecipeStep: boolean;
}

export class RecipePrepTaskStepCreationRequestInput implements IRecipePrepTaskStepCreationRequestInput {
  belongsToRecipeStep: string;
  satisfiesRecipeStep: boolean;
  constructor(input: Partial<RecipePrepTaskStepCreationRequestInput> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.satisfiesRecipeStep = input.satisfiesRecipeStep || false;
  }
}
