// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeRating {
  createdAt: string;
  difficulty: number;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  byUser: string;
  cleanup: number;
  instructions: number;
  notes: string;
  overall: number;
  recipeID: string;
  taste: number;
}

export class RecipeRating implements IRecipeRating {
  createdAt: string;
  difficulty: number;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  byUser: string;
  cleanup: number;
  instructions: number;
  notes: string;
  overall: number;
  recipeID: string;
  taste: number;
  constructor(input: Partial<RecipeRating> = {}) {
    this.createdAt = input.createdAt = '';
    this.difficulty = input.difficulty = 0;
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
    this.byUser = input.byUser = '';
    this.cleanup = input.cleanup = 0;
    this.instructions = input.instructions = 0;
    this.notes = input.notes = '';
    this.overall = input.overall = 0;
    this.recipeID = input.recipeID = '';
    this.taste = input.taste = 0;
  }
}
