// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';

export interface IValidMeasurementUnitConversion {
  onlyForIngredient?: ValidIngredient;
  createdAt: string;
  from: ValidMeasurementUnit;
  id: string;
  modifier: number;
  notes: string;
  archivedAt?: string;
  lastUpdatedAt?: string;
  to: ValidMeasurementUnit;
}

export class ValidMeasurementUnitConversion implements IValidMeasurementUnitConversion {
  onlyForIngredient?: ValidIngredient;
  createdAt: string;
  from: ValidMeasurementUnit;
  id: string;
  modifier: number;
  notes: string;
  archivedAt?: string;
  lastUpdatedAt?: string;
  to: ValidMeasurementUnit;
  constructor(input: Partial<ValidMeasurementUnitConversion> = {}) {
    this.onlyForIngredient = input.onlyForIngredient;
    this.createdAt = input.createdAt = '';
    this.from = input.from = new ValidMeasurementUnit();
    this.id = input.id = '';
    this.modifier = input.modifier = 0;
    this.notes = input.notes = '';
    this.archivedAt = input.archivedAt;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.to = input.to = new ValidMeasurementUnit();
  }
}
