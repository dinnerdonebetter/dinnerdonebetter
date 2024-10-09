// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeMedia {
  index: number;
  lastUpdatedAt?: string;
  archivedAt?: string;
  belongsToRecipe?: string;
  createdAt: string;
  id: string;
  belongsToRecipeStep?: string;
  externalPath: string;
  internalPath: string;
  mimeType: string;
}

export class RecipeMedia implements IRecipeMedia {
  index: number;
  lastUpdatedAt?: string;
  archivedAt?: string;
  belongsToRecipe?: string;
  createdAt: string;
  id: string;
  belongsToRecipeStep?: string;
  externalPath: string;
  internalPath: string;
  mimeType: string;
  constructor(input: Partial<RecipeMedia> = {}) {
    this.index = input.index = 0;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
    this.belongsToRecipe = input.belongsToRecipe;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.externalPath = input.externalPath = '';
    this.internalPath = input.internalPath = '';
    this.mimeType = input.mimeType = '';
  }
}
