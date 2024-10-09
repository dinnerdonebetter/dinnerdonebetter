// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeRating {
  createdAt: string;
  instructions: number;
  notes: string;
  overall: number;
  archivedAt: string;
  byUser: string;
  cleanup: number;
  recipeID: string;
  taste: number;
  difficulty: number;
  id: string;
  lastUpdatedAt: string;
}

export class RecipeRating implements IRecipeRating {
  createdAt: string;
  instructions: number;
  notes: string;
  overall: number;
  archivedAt: string;
  byUser: string;
  cleanup: number;
  recipeID: string;
  taste: number;
  difficulty: number;
  id: string;
  lastUpdatedAt: string;
  constructor(input: Partial<RecipeRating> = {}) {
    this.createdAt = input.createdAt || '';
    this.instructions = input.instructions || 0;
    this.notes = input.notes || '';
    this.overall = input.overall || 0;
    this.archivedAt = input.archivedAt || '';
    this.byUser = input.byUser || '';
    this.cleanup = input.cleanup || 0;
    this.recipeID = input.recipeID || '';
    this.taste = input.taste || 0;
    this.difficulty = input.difficulty || 0;
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
  }
}
