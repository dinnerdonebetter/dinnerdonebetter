// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeStepCompletionConditionIngredient } from './RecipeStepCompletionConditionIngredient';
import { ValidIngredientState } from './ValidIngredientState';

export interface IRecipeStepCompletionCondition {
  createdAt: string;
  id: string;
  ingredientState: ValidIngredientState;
  ingredients: RecipeStepCompletionConditionIngredient;
  lastUpdatedAt?: string;
  notes: string;
  archivedAt?: string;
  belongsToRecipeStep: string;
  optional: boolean;
}

export class RecipeStepCompletionCondition implements IRecipeStepCompletionCondition {
  createdAt: string;
  id: string;
  ingredientState: ValidIngredientState;
  ingredients: RecipeStepCompletionConditionIngredient;
  lastUpdatedAt?: string;
  notes: string;
  archivedAt?: string;
  belongsToRecipeStep: string;
  optional: boolean;
  constructor(input: Partial<RecipeStepCompletionCondition> = {}) {
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.ingredientState = input.ingredientState = new ValidIngredientState();
    this.ingredients = input.ingredients = new RecipeStepCompletionConditionIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.archivedAt = input.archivedAt;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.optional = input.optional = false;
  }
}
