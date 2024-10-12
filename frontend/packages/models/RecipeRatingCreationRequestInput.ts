// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IRecipeRatingCreationRequestInput {
   byUser: string;
 cleanup: number;
 difficulty: number;
 instructions: number;
 notes: string;
 overall: number;
 recipeID: string;
 taste: number;

}

export class RecipeRatingCreationRequestInput implements IRecipeRatingCreationRequestInput {
   byUser: string;
 cleanup: number;
 difficulty: number;
 instructions: number;
 notes: string;
 overall: number;
 recipeID: string;
 taste: number;
constructor(input: Partial<RecipeRatingCreationRequestInput> = {}) {
	 this.byUser = input.byUser || '';
 this.cleanup = input.cleanup || 0;
 this.difficulty = input.difficulty || 0;
 this.instructions = input.instructions || 0;
 this.notes = input.notes || '';
 this.overall = input.overall || 0;
 this.recipeID = input.recipeID || '';
 this.taste = input.taste || 0;
}
}