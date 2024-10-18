// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionIngredient } from './RecipeStepCompletionConditionIngredient';
import { ValidIngredientState } from './ValidIngredientState';

export interface IRecipeStepCompletionCondition {
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  ingredientState: ValidIngredientState;
  ingredients: RecipeStepCompletionConditionIngredient[];
  lastUpdatedAt: string;
  notes: string;
  optional: boolean;
}

export class RecipeStepCompletionCondition implements IRecipeStepCompletionCondition {
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  ingredientState: ValidIngredientState;
  ingredients: RecipeStepCompletionConditionIngredient[];
  lastUpdatedAt: string;
  notes: string;
  optional: boolean;
  constructor(input: Partial<RecipeStepCompletionCondition> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.ingredientState = input.ingredientState || new ValidIngredientState();
    this.ingredients = input.ingredients || [];
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.optional = input.optional || false;
  }
}
