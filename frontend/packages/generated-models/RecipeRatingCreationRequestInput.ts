// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeRatingCreationRequestInput {
  overall: number;
  recipeID: string;
  taste: number;
  byUser: string;
  cleanup: number;
  difficulty: number;
  instructions: number;
  notes: string;
}

export class RecipeRatingCreationRequestInput implements IRecipeRatingCreationRequestInput {
  overall: number;
  recipeID: string;
  taste: number;
  byUser: string;
  cleanup: number;
  difficulty: number;
  instructions: number;
  notes: string;
  constructor(input: Partial<RecipeRatingCreationRequestInput> = {}) {
    this.overall = input.overall = 0;
    this.recipeID = input.recipeID = '';
    this.taste = input.taste = 0;
    this.byUser = input.byUser = '';
    this.cleanup = input.cleanup = 0;
    this.difficulty = input.difficulty = 0;
    this.instructions = input.instructions = 0;
    this.notes = input.notes = '';
  }
}
