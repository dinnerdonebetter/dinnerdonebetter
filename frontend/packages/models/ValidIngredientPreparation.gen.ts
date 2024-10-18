// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient.gen';
import { ValidPreparation } from './ValidPreparation.gen';

export interface IValidIngredientPreparation {
  archivedAt: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
}

export class ValidIngredientPreparation implements IValidIngredientPreparation {
  archivedAt: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
  constructor(input: Partial<ValidIngredientPreparation> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.ingredient = input.ingredient || new ValidIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.preparation = input.preparation || new ValidPreparation();
  }
}
