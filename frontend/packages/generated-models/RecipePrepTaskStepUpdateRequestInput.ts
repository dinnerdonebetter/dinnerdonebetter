// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IRecipePrepTaskStepUpdateRequestInput {
   satisfiesRecipeStep?: boolean;
 belongsToRecipeStep?: string;
 belongsToRecipeStepTask?: string;

}

export class RecipePrepTaskStepUpdateRequestInput implements IRecipePrepTaskStepUpdateRequestInput {
   satisfiesRecipeStep?: boolean;
 belongsToRecipeStep?: string;
 belongsToRecipeStepTask?: string;
constructor(input: Partial<RecipePrepTaskStepUpdateRequestInput> = {}) {
	 this.satisfiesRecipeStep = input.satisfiesRecipeStep;
 this.belongsToRecipeStep = input.belongsToRecipeStep;
 this.belongsToRecipeStepTask = input.belongsToRecipeStepTask;
}
}