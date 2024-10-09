// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IRecipeRatingUpdateRequestInput {
   difficulty?: number;
 instructions?: number;
 notes?: string;
 overall?: number;
 recipeID?: string;
 taste?: number;
 byUser?: string;
 cleanup?: number;

}

export class RecipeRatingUpdateRequestInput implements IRecipeRatingUpdateRequestInput {
   difficulty?: number;
 instructions?: number;
 notes?: string;
 overall?: number;
 recipeID?: string;
 taste?: number;
 byUser?: string;
 cleanup?: number;
constructor(input: Partial<RecipeRatingUpdateRequestInput> = {}) {
	 this.difficulty = input.difficulty;
 this.instructions = input.instructions;
 this.notes = input.notes;
 this.overall = input.overall;
 this.recipeID = input.recipeID;
 this.taste = input.taste;
 this.byUser = input.byUser;
 this.cleanup = input.cleanup;
}
}