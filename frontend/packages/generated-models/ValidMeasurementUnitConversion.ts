// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';

export interface IValidMeasurementUnitConversion {
  modifier: number;
  to: ValidMeasurementUnit;
  onlyForIngredient?: ValidIngredient;
  archivedAt?: string;
  createdAt: string;
  from: ValidMeasurementUnit;
  id: string;
  lastUpdatedAt?: string;
  notes: string;
}

export class ValidMeasurementUnitConversion implements IValidMeasurementUnitConversion {
  modifier: number;
  to: ValidMeasurementUnit;
  onlyForIngredient?: ValidIngredient;
  archivedAt?: string;
  createdAt: string;
  from: ValidMeasurementUnit;
  id: string;
  lastUpdatedAt?: string;
  notes: string;
  constructor(input: Partial<ValidMeasurementUnitConversion> = {}) {
    this.modifier = input.modifier = 0;
    this.to = input.to = new ValidMeasurementUnit();
    this.onlyForIngredient = input.onlyForIngredient;
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.from = input.from = new ValidMeasurementUnit();
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
  }
}
