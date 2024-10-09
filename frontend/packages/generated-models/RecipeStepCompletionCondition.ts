// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionIngredient } from './RecipeStepCompletionConditionIngredient';
import { ValidIngredientState } from './ValidIngredientState';

export interface IRecipeStepCompletionCondition {
  archivedAt?: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  ingredients: RecipeStepCompletionConditionIngredient;
  lastUpdatedAt?: string;
  ingredientState: ValidIngredientState;
  notes: string;
  optional: boolean;
}

export class RecipeStepCompletionCondition implements IRecipeStepCompletionCondition {
  archivedAt?: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  ingredients: RecipeStepCompletionConditionIngredient;
  lastUpdatedAt?: string;
  ingredientState: ValidIngredientState;
  notes: string;
  optional: boolean;
  constructor(input: Partial<RecipeStepCompletionCondition> = {}) {
    this.archivedAt = input.archivedAt;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.ingredients = input.ingredients = new RecipeStepCompletionConditionIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.ingredientState = input.ingredientState = new ValidIngredientState();
    this.notes = input.notes = '';
    this.optional = input.optional = false;
  }
}
