// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidIngredientState } from './ValidIngredientState';

export interface IValidIngredientStateIngredient {
  archivedAt?: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  ingredientState: ValidIngredientState;
  lastUpdatedAt?: string;
  notes: string;
}

export class ValidIngredientStateIngredient implements IValidIngredientStateIngredient {
  archivedAt?: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  ingredientState: ValidIngredientState;
  lastUpdatedAt?: string;
  notes: string;
  constructor(input: Partial<ValidIngredientStateIngredient> = {}) {
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.ingredient = input.ingredient = new ValidIngredient();
    this.ingredientState = input.ingredientState = new ValidIngredientState();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
  }
}
