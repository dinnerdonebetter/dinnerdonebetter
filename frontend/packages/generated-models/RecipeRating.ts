// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IRecipeRating {
   taste: number;
 byUser: string;
 instructions: number;
 notes: string;
 difficulty: number;
 id: string;
 lastUpdatedAt?: string;
 overall: number;
 recipeID: string;
 archivedAt?: string;
 cleanup: number;
 createdAt: string;

}

export class RecipeRating implements IRecipeRating {
   taste: number;
 byUser: string;
 instructions: number;
 notes: string;
 difficulty: number;
 id: string;
 lastUpdatedAt?: string;
 overall: number;
 recipeID: string;
 archivedAt?: string;
 cleanup: number;
 createdAt: string;
constructor(input: Partial<RecipeRating> = {}) {
	 this.taste = input.taste = 0;
 this.byUser = input.byUser = '';
 this.instructions = input.instructions = 0;
 this.notes = input.notes = '';
 this.difficulty = input.difficulty = 0;
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.overall = input.overall = 0;
 this.recipeID = input.recipeID = '';
 this.archivedAt = input.archivedAt;
 this.cleanup = input.cleanup = 0;
 this.createdAt = input.createdAt = '';
}
}