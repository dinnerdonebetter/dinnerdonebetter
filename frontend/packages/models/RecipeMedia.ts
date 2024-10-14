// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeMedia {
  archivedAt: string;
  belongsToRecipe: string;
  belongsToRecipeStep: string;
  createdAt: string;
  externalPath: string;
  id: string;
  index: number;
  internalPath: string;
  lastUpdatedAt: string;
  mimeType: string;
}

export class RecipeMedia implements IRecipeMedia {
  archivedAt: string;
  belongsToRecipe: string;
  belongsToRecipeStep: string;
  createdAt: string;
  externalPath: string;
  id: string;
  index: number;
  internalPath: string;
  lastUpdatedAt: string;
  mimeType: string;
  constructor(input: Partial<RecipeMedia> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipe = input.belongsToRecipe || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.createdAt = input.createdAt || '';
    this.externalPath = input.externalPath || '';
    this.id = input.id || '';
    this.index = input.index || 0;
    this.internalPath = input.internalPath || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.mimeType = input.mimeType || '';
  }
}
