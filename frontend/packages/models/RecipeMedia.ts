// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeMedia {
  belongsToRecipeStep: string;
  createdAt: string;
  externalPath: string;
  id: string;
  archivedAt: string;
  belongsToRecipe: string;
  index: number;
  internalPath: string;
  lastUpdatedAt: string;
  mimeType: string;
}

export class RecipeMedia implements IRecipeMedia {
  belongsToRecipeStep: string;
  createdAt: string;
  externalPath: string;
  id: string;
  archivedAt: string;
  belongsToRecipe: string;
  index: number;
  internalPath: string;
  lastUpdatedAt: string;
  mimeType: string;
  constructor(input: Partial<RecipeMedia> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.createdAt = input.createdAt || '';
    this.externalPath = input.externalPath || '';
    this.id = input.id || '';
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.index = input.index || 0;
    this.internalPath = input.internalPath || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.mimeType = input.mimeType || '';
  }
}
