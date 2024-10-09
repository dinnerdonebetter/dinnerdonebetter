// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionIngredient } from './RecipeStepCompletionConditionIngredient';
import { ValidIngredientState } from './ValidIngredientState';

export interface IRecipeStepCompletionCondition {
  optional: boolean;
  archivedAt?: string;
  belongsToRecipeStep: string;
  ingredientState: ValidIngredientState;
  notes: string;
  createdAt: string;
  id: string;
  ingredients: RecipeStepCompletionConditionIngredient;
  lastUpdatedAt?: string;
}

export class RecipeStepCompletionCondition implements IRecipeStepCompletionCondition {
  optional: boolean;
  archivedAt?: string;
  belongsToRecipeStep: string;
  ingredientState: ValidIngredientState;
  notes: string;
  createdAt: string;
  id: string;
  ingredients: RecipeStepCompletionConditionIngredient;
  lastUpdatedAt?: string;
  constructor(input: Partial<RecipeStepCompletionCondition> = {}) {
    this.optional = input.optional = false;
    this.archivedAt = input.archivedAt;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.ingredientState = input.ingredientState = new ValidIngredientState();
    this.notes = input.notes = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.ingredients = input.ingredients = new RecipeStepCompletionConditionIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
