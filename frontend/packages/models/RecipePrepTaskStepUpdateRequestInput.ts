// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipePrepTaskStepUpdateRequestInput {
  belongsToRecipeStepTask: string;
  satisfiesRecipeStep: boolean;
  belongsToRecipeStep: string;
}

export class RecipePrepTaskStepUpdateRequestInput implements IRecipePrepTaskStepUpdateRequestInput {
  belongsToRecipeStepTask: string;
  satisfiesRecipeStep: boolean;
  belongsToRecipeStep: string;
  constructor(input: Partial<RecipePrepTaskStepUpdateRequestInput> = {}) {
    this.belongsToRecipeStepTask = input.belongsToRecipeStepTask || '';
    this.satisfiesRecipeStep = input.satisfiesRecipeStep || false;
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
  }
}
