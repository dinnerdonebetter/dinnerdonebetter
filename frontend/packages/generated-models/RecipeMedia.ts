// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeMedia {
  mimeType: string;
  archivedAt?: string;
  belongsToRecipe?: string;
  createdAt: string;
  externalPath: string;
  id: string;
  internalPath: string;
  lastUpdatedAt?: string;
  belongsToRecipeStep?: string;
  index: number;
}

export class RecipeMedia implements IRecipeMedia {
  mimeType: string;
  archivedAt?: string;
  belongsToRecipe?: string;
  createdAt: string;
  externalPath: string;
  id: string;
  internalPath: string;
  lastUpdatedAt?: string;
  belongsToRecipeStep?: string;
  index: number;
  constructor(input: Partial<RecipeMedia> = {}) {
    this.mimeType = input.mimeType = '';
    this.archivedAt = input.archivedAt;
    this.belongsToRecipe = input.belongsToRecipe;
    this.createdAt = input.createdAt = '';
    this.externalPath = input.externalPath = '';
    this.id = input.id = '';
    this.internalPath = input.internalPath = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.index = input.index = 0;
  }
}
