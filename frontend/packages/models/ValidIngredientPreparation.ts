// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidPreparation } from './ValidPreparation';

export interface IValidIngredientPreparation {
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
  archivedAt: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
}

export class ValidIngredientPreparation implements IValidIngredientPreparation {
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
  archivedAt: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  constructor(input: Partial<ValidIngredientPreparation> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.preparation = input.preparation || new ValidPreparation();
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.ingredient = input.ingredient || new ValidIngredient();
  }
}
