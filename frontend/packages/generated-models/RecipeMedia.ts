// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeMedia {
  lastUpdatedAt?: string;
  archivedAt?: string;
  externalPath: string;
  createdAt: string;
  id: string;
  index: number;
  internalPath: string;
  mimeType: string;
  belongsToRecipe?: string;
  belongsToRecipeStep?: string;
}

export class RecipeMedia implements IRecipeMedia {
  lastUpdatedAt?: string;
  archivedAt?: string;
  externalPath: string;
  createdAt: string;
  id: string;
  index: number;
  internalPath: string;
  mimeType: string;
  belongsToRecipe?: string;
  belongsToRecipeStep?: string;
  constructor(input: Partial<RecipeMedia> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
    this.externalPath = input.externalPath = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.index = input.index = 0;
    this.internalPath = input.internalPath = '';
    this.mimeType = input.mimeType = '';
    this.belongsToRecipe = input.belongsToRecipe;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
  }
}
