// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IRecipePrepTaskStep {
   belongsToRecipeStepTask: string;
 id: string;
 satisfiesRecipeStep: boolean;
 belongsToRecipeStep: string;

}

export class RecipePrepTaskStep implements IRecipePrepTaskStep {
   belongsToRecipeStepTask: string;
 id: string;
 satisfiesRecipeStep: boolean;
 belongsToRecipeStep: string;
constructor(input: Partial<RecipePrepTaskStep> = {}) {
	 this.belongsToRecipeStepTask = input.belongsToRecipeStepTask = '';
 this.id = input.id = '';
 this.satisfiesRecipeStep = input.satisfiesRecipeStep = false;
 this.belongsToRecipeStep = input.belongsToRecipeStep = '';
}
}