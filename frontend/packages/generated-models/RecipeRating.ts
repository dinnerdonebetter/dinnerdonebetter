// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeRating {
  lastUpdatedAt?: string;
  notes: string;
  overall: number;
  recipeID: string;
  archivedAt?: string;
  cleanup: number;
  createdAt: string;
  difficulty: number;
  taste: number;
  byUser: string;
  id: string;
  instructions: number;
}

export class RecipeRating implements IRecipeRating {
  lastUpdatedAt?: string;
  notes: string;
  overall: number;
  recipeID: string;
  archivedAt?: string;
  cleanup: number;
  createdAt: string;
  difficulty: number;
  taste: number;
  byUser: string;
  id: string;
  instructions: number;
  constructor(input: Partial<RecipeRating> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.overall = input.overall = 0;
    this.recipeID = input.recipeID = '';
    this.archivedAt = input.archivedAt;
    this.cleanup = input.cleanup = 0;
    this.createdAt = input.createdAt = '';
    this.difficulty = input.difficulty = 0;
    this.taste = input.taste = 0;
    this.byUser = input.byUser = '';
    this.id = input.id = '';
    this.instructions = input.instructions = 0;
  }
}
