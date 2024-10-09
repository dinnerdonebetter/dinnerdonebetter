// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidIngredientState } from './ValidIngredientState';

export interface IValidIngredientStateIngredient {
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  ingredientState: ValidIngredientState;
  lastUpdatedAt: string;
  notes: string;
  archivedAt: string;
}

export class ValidIngredientStateIngredient implements IValidIngredientStateIngredient {
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  ingredientState: ValidIngredientState;
  lastUpdatedAt: string;
  notes: string;
  archivedAt: string;
  constructor(input: Partial<ValidIngredientStateIngredient> = {}) {
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.ingredient = input.ingredient || new ValidIngredient();
    this.ingredientState = input.ingredientState || new ValidIngredientState();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.archivedAt = input.archivedAt || '';
  }
}
