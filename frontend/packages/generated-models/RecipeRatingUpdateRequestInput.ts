// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeRatingUpdateRequestInput {
  notes?: string;
  overall?: number;
  recipeID?: string;
  taste?: number;
  byUser?: string;
  cleanup?: number;
  difficulty?: number;
  instructions?: number;
}

export class RecipeRatingUpdateRequestInput implements IRecipeRatingUpdateRequestInput {
  notes?: string;
  overall?: number;
  recipeID?: string;
  taste?: number;
  byUser?: string;
  cleanup?: number;
  difficulty?: number;
  instructions?: number;
  constructor(input: Partial<RecipeRatingUpdateRequestInput> = {}) {
    this.notes = input.notes;
    this.overall = input.overall;
    this.recipeID = input.recipeID;
    this.taste = input.taste;
    this.byUser = input.byUser;
    this.cleanup = input.cleanup;
    this.difficulty = input.difficulty;
    this.instructions = input.instructions;
  }
}
