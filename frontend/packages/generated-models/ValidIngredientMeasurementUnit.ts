// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidIngredientMeasurementUnit {
  lastUpdatedAt?: string;
  measurementUnit: ValidMeasurementUnit;
  notes: string;
  allowableQuantity: NumberRangeWithOptionalMax;
  archivedAt?: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
}

export class ValidIngredientMeasurementUnit implements IValidIngredientMeasurementUnit {
  lastUpdatedAt?: string;
  measurementUnit: ValidMeasurementUnit;
  notes: string;
  allowableQuantity: NumberRangeWithOptionalMax;
  archivedAt?: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  constructor(input: Partial<ValidIngredientMeasurementUnit> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.measurementUnit = input.measurementUnit = new ValidMeasurementUnit();
    this.notes = input.notes = '';
    this.allowableQuantity = input.allowableQuantity = { min: 0 };
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.ingredient = input.ingredient = new ValidIngredient();
  }
}
