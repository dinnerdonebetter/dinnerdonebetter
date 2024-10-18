// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeStepCompletionConditionIngredient {
  archivedAt: string;
  belongsToRecipeStepCompletionCondition: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  recipeStepIngredient: string;
}

export class RecipeStepCompletionConditionIngredient implements IRecipeStepCompletionConditionIngredient {
  archivedAt: string;
  belongsToRecipeStepCompletionCondition: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  recipeStepIngredient: string;
  constructor(input: Partial<RecipeStepCompletionConditionIngredient> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStepCompletionCondition = input.belongsToRecipeStepCompletionCondition || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.recipeStepIngredient = input.recipeStepIngredient || '';
  }
}
