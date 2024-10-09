// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidPreparation } from './ValidPreparation';

export interface IValidIngredientPreparation {
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt?: string;
  notes: string;
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
}

export class ValidIngredientPreparation implements IValidIngredientPreparation {
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt?: string;
  notes: string;
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
  constructor(input: Partial<ValidIngredientPreparation> = {}) {
    this.id = input.id = '';
    this.ingredient = input.ingredient = new ValidIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.preparation = input.preparation = new ValidPreparation();
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
  }
}
