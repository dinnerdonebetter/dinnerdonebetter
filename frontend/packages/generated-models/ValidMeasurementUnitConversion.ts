// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';

export interface IValidMeasurementUnitConversion {
  archivedAt?: string;
  id: string;
  notes: string;
  onlyForIngredient?: ValidIngredient;
  createdAt: string;
  from: ValidMeasurementUnit;
  lastUpdatedAt?: string;
  modifier: number;
  to: ValidMeasurementUnit;
}

export class ValidMeasurementUnitConversion implements IValidMeasurementUnitConversion {
  archivedAt?: string;
  id: string;
  notes: string;
  onlyForIngredient?: ValidIngredient;
  createdAt: string;
  from: ValidMeasurementUnit;
  lastUpdatedAt?: string;
  modifier: number;
  to: ValidMeasurementUnit;
  constructor(input: Partial<ValidMeasurementUnitConversion> = {}) {
    this.archivedAt = input.archivedAt;
    this.id = input.id = '';
    this.notes = input.notes = '';
    this.onlyForIngredient = input.onlyForIngredient;
    this.createdAt = input.createdAt = '';
    this.from = input.from = new ValidMeasurementUnit();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.modifier = input.modifier = 0;
    this.to = input.to = new ValidMeasurementUnit();
  }
}
