// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IRecipePrepTaskStep {
   belongsToRecipeStep: string;
 belongsToRecipeStepTask: string;
 id: string;
 satisfiesRecipeStep: boolean;

}

export class RecipePrepTaskStep implements IRecipePrepTaskStep {
   belongsToRecipeStep: string;
 belongsToRecipeStepTask: string;
 id: string;
 satisfiesRecipeStep: boolean;
constructor(input: Partial<RecipePrepTaskStep> = {}) {
	 this.belongsToRecipeStep = input.belongsToRecipeStep || '';
 this.belongsToRecipeStepTask = input.belongsToRecipeStepTask || '';
 this.id = input.id || '';
 this.satisfiesRecipeStep = input.satisfiesRecipeStep || false;
}
}