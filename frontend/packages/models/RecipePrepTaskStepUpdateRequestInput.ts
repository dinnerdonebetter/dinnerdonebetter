// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipePrepTaskStepUpdateRequestInput {
  belongsToRecipeStep: string;
  belongsToRecipeStepTask: string;
  satisfiesRecipeStep: boolean;
}

export class RecipePrepTaskStepUpdateRequestInput implements IRecipePrepTaskStepUpdateRequestInput {
  belongsToRecipeStep: string;
  belongsToRecipeStepTask: string;
  satisfiesRecipeStep: boolean;
  constructor(input: Partial<RecipePrepTaskStepUpdateRequestInput> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.belongsToRecipeStepTask = input.belongsToRecipeStepTask || '';
    this.satisfiesRecipeStep = input.satisfiesRecipeStep || false;
  }
}
