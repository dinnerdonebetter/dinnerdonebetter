// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IRecipeRating {
   archivedAt: string;
 byUser: string;
 cleanup: number;
 createdAt: string;
 difficulty: number;
 id: string;
 instructions: number;
 lastUpdatedAt: string;
 notes: string;
 overall: number;
 recipeID: string;
 taste: number;

}

export class RecipeRating implements IRecipeRating {
   archivedAt: string;
 byUser: string;
 cleanup: number;
 createdAt: string;
 difficulty: number;
 id: string;
 instructions: number;
 lastUpdatedAt: string;
 notes: string;
 overall: number;
 recipeID: string;
 taste: number;
constructor(input: Partial<RecipeRating> = {}) {
	 this.archivedAt = input.archivedAt || '';
 this.byUser = input.byUser || '';
 this.cleanup = input.cleanup || 0;
 this.createdAt = input.createdAt || '';
 this.difficulty = input.difficulty || 0;
 this.id = input.id || '';
 this.instructions = input.instructions || 0;
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.notes = input.notes || '';
 this.overall = input.overall || 0;
 this.recipeID = input.recipeID || '';
 this.taste = input.taste || 0;
}
}