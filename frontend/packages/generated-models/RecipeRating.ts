// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeRating {
  cleanup: number;
  instructions: number;
  notes: string;
  taste: number;
  id: string;
  lastUpdatedAt?: string;
  overall: number;
  recipeID: string;
  archivedAt?: string;
  byUser: string;
  createdAt: string;
  difficulty: number;
}

export class RecipeRating implements IRecipeRating {
  cleanup: number;
  instructions: number;
  notes: string;
  taste: number;
  id: string;
  lastUpdatedAt?: string;
  overall: number;
  recipeID: string;
  archivedAt?: string;
  byUser: string;
  createdAt: string;
  difficulty: number;
  constructor(input: Partial<RecipeRating> = {}) {
    this.cleanup = input.cleanup = 0;
    this.instructions = input.instructions = 0;
    this.notes = input.notes = '';
    this.taste = input.taste = 0;
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.overall = input.overall = 0;
    this.recipeID = input.recipeID = '';
    this.archivedAt = input.archivedAt;
    this.byUser = input.byUser = '';
    this.createdAt = input.createdAt = '';
    this.difficulty = input.difficulty = 0;
  }
}
