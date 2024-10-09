// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionIngredient } from './RecipeStepCompletionConditionIngredient';
import { ValidIngredientState } from './ValidIngredientState';

export interface IRecipeStepCompletionCondition {
  notes: string;
  createdAt: string;
  ingredients: RecipeStepCompletionConditionIngredient[];
  id: string;
  ingredientState: ValidIngredientState;
  lastUpdatedAt: string;
  optional: boolean;
  archivedAt: string;
  belongsToRecipeStep: string;
}

export class RecipeStepCompletionCondition implements IRecipeStepCompletionCondition {
  notes: string;
  createdAt: string;
  ingredients: RecipeStepCompletionConditionIngredient[];
  id: string;
  ingredientState: ValidIngredientState;
  lastUpdatedAt: string;
  optional: boolean;
  archivedAt: string;
  belongsToRecipeStep: string;
  constructor(input: Partial<RecipeStepCompletionCondition> = {}) {
    this.notes = input.notes || '';
    this.createdAt = input.createdAt || '';
    this.ingredients = input.ingredients || [];
    this.id = input.id || '';
    this.ingredientState = input.ingredientState || new ValidIngredientState();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.optional = input.optional || false;
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
  }
}
