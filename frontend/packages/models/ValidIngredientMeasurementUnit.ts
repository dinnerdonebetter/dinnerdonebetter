// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidIngredientMeasurementUnit {
  allowableQuantity: NumberRangeWithOptionalMax;
  archivedAt: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  measurementUnit: ValidMeasurementUnit;
  notes: string;
}

export class ValidIngredientMeasurementUnit implements IValidIngredientMeasurementUnit {
  allowableQuantity: NumberRangeWithOptionalMax;
  archivedAt: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  measurementUnit: ValidMeasurementUnit;
  notes: string;
  constructor(input: Partial<ValidIngredientMeasurementUnit> = {}) {
    this.allowableQuantity = input.allowableQuantity || { min: 0 };
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.ingredient = input.ingredient || new ValidIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.measurementUnit = input.measurementUnit || new ValidMeasurementUnit();
    this.notes = input.notes || '';
  }
}
