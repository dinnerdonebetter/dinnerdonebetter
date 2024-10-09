// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IRecipeStepCompletionConditionIngredient {
  recipeStepIngredient: string;
  archivedAt: string;
  belongsToRecipeStepCompletionCondition: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
}

export class RecipeStepCompletionConditionIngredient implements IRecipeStepCompletionConditionIngredient {
  recipeStepIngredient: string;
  archivedAt: string;
  belongsToRecipeStepCompletionCondition: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  constructor(input: Partial<RecipeStepCompletionConditionIngredient> = {}) {
    this.recipeStepIngredient = input.recipeStepIngredient || '';
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStepCompletionCondition = input.belongsToRecipeStepCompletionCondition || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
  }
}
