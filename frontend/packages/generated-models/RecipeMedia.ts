// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IRecipeMedia {
   belongsToRecipe?: string;
 internalPath: string;
 createdAt: string;
 externalPath: string;
 id: string;
 index: number;
 lastUpdatedAt?: string;
 mimeType: string;
 archivedAt?: string;
 belongsToRecipeStep?: string;

}

export class RecipeMedia implements IRecipeMedia {
   belongsToRecipe?: string;
 internalPath: string;
 createdAt: string;
 externalPath: string;
 id: string;
 index: number;
 lastUpdatedAt?: string;
 mimeType: string;
 archivedAt?: string;
 belongsToRecipeStep?: string;
constructor(input: Partial<RecipeMedia> = {}) {
	 this.belongsToRecipe = input.belongsToRecipe;
 this.internalPath = input.internalPath = '';
 this.createdAt = input.createdAt = '';
 this.externalPath = input.externalPath = '';
 this.id = input.id = '';
 this.index = input.index = 0;
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.mimeType = input.mimeType = '';
 this.archivedAt = input.archivedAt;
 this.belongsToRecipeStep = input.belongsToRecipeStep;
}
}