// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeStepCompletionConditionIngredient {
  lastUpdatedAt?: string;
  recipeStepIngredient: string;
  archivedAt?: string;
  belongsToRecipeStepCompletionCondition: string;
  createdAt: string;
  id: string;
}

export class RecipeStepCompletionConditionIngredient implements IRecipeStepCompletionConditionIngredient {
  lastUpdatedAt?: string;
  recipeStepIngredient: string;
  archivedAt?: string;
  belongsToRecipeStepCompletionCondition: string;
  createdAt: string;
  id: string;
  constructor(input: Partial<RecipeStepCompletionConditionIngredient> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.recipeStepIngredient = input.recipeStepIngredient = '';
    this.archivedAt = input.archivedAt;
    this.belongsToRecipeStepCompletionCondition = input.belongsToRecipeStepCompletionCondition = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
  }
}
