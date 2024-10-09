// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeStepCompletionConditionIngredient {
  id: string;
  lastUpdatedAt?: string;
  recipeStepIngredient: string;
  archivedAt?: string;
  belongsToRecipeStepCompletionCondition: string;
  createdAt: string;
}

export class RecipeStepCompletionConditionIngredient implements IRecipeStepCompletionConditionIngredient {
  id: string;
  lastUpdatedAt?: string;
  recipeStepIngredient: string;
  archivedAt?: string;
  belongsToRecipeStepCompletionCondition: string;
  createdAt: string;
  constructor(input: Partial<RecipeStepCompletionConditionIngredient> = {}) {
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.recipeStepIngredient = input.recipeStepIngredient = '';
    this.archivedAt = input.archivedAt;
    this.belongsToRecipeStepCompletionCondition = input.belongsToRecipeStepCompletionCondition = '';
    this.createdAt = input.createdAt = '';
  }
}
