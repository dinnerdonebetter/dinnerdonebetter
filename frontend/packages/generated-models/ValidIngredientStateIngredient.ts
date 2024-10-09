// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidIngredientState } from './ValidIngredientState';

export interface IValidIngredientStateIngredient {
  id: string;
  ingredient: ValidIngredient;
  ingredientState: ValidIngredientState;
  lastUpdatedAt?: string;
  notes: string;
  archivedAt?: string;
  createdAt: string;
}

export class ValidIngredientStateIngredient implements IValidIngredientStateIngredient {
  id: string;
  ingredient: ValidIngredient;
  ingredientState: ValidIngredientState;
  lastUpdatedAt?: string;
  notes: string;
  archivedAt?: string;
  createdAt: string;
  constructor(input: Partial<ValidIngredientStateIngredient> = {}) {
    this.id = input.id = '';
    this.ingredient = input.ingredient = new ValidIngredient();
    this.ingredientState = input.ingredientState = new ValidIngredientState();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
  }
}
